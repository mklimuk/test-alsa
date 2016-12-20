package audio

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/mklimuk/test-alsa/config"
	"github.com/mklimuk/websocket"
)

const sampleSizeBytes = 2
const textMessage = 1
const defaultSampleRate = 22050
const defaultChannels = 1

var introEndMsg []byte
var introStartMsg []byte

func init() {
	introStartMsg, _ = json.Marshal(SignallingMsg{"playback:intro:start", ""})
	introEndMsg, _ = json.Marshal(SignallingMsg{"playback:intro:end", ""})
}

//SignallingMsg is a message exchanged with the peer during websocket playback
type SignallingMsg struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

//Playback is the interface responsible for communication with audio device
type Playback interface {
	DeviceBusy() (bool, int)
	PlaybackContext() *StreamContext
	PlayFromWsConnection(c websocket.Connection) error
}

//StreamContext contains information about currently playing stream
type StreamContext struct {
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	Volume      int    `json:"volume"`
	Type        string `json:"type"`
	PlayIntro   bool   `json:"playIntro"`
	SampleRate  int    `json:"sampleRate"`
	Channels    int    `json:"channels"`
	BufferSize  int    `json:"bufferSize"`
	bytesRead   int
	framesWrote int
}

type play struct {
	connMutex  sync.Mutex
	context    *StreamContext
	connection websocket.Connection
	introFile  string
	bufParams  *BufferParams
	factory    DeviceFactory
}

//New is the playback interface constructor
func New(conf *config.AudioConf, factory DeviceFactory, introFile string) Playback {
	p := play{
		factory:   factory,
		bufParams: &BufferParams{BufferFrames: conf.DeviceBuffer, PeriodFrames: conf.PeriodFrames, Periods: conf.Periods},
		introFile: introFile,
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "ws.audio-endpoint.audio", "method": "New", "introFile": p.introFile, "bufParams": fmt.Sprintf("%+v", &(p.bufParams))}).
			Debug("Playback configuration")
	}
	return &p
}

func (p *play) PlaybackContext() *StreamContext {
	return p.context
}

func (p *play) DeviceBusy() (bool, int) {
	p.connMutex.Lock()
	defer p.connMutex.Unlock()
	if p.context == nil {
		return false, 0
	}
	return true, p.context.Priority
}

func (p *play) PlayFromWsConnection(c websocket.Connection) error {
	p.connMutex.Lock()
	defer p.connMutex.Unlock()
	//wait for the context
	var mt int
	var message []byte
	var err error
	if mt, message, err = c.ReadMessage(); err != nil {
		log.WithFields(log.Fields{"logger": "ws.audio-endpoint.audio", "method": "PlayFromWsConnection"}).
			WithError(err).Error("Websocket read error encountered")
		return err
	}

	if mt != textMessage {
		return errors.New("Wrong message type. First message should be a text one")
	}

	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "ws.audio-endpoint.audio", "method": "PlayFromWsConnection", "message": message}).
			Debug("Received message")
	}
	var context = new(StreamContext)
	if err = json.Unmarshal(message, context); err != nil {
		log.WithFields(log.Fields{"logger": "ws.audio-endpoint.audio", "method": "PlayFromWsConnection"}).
			WithError(err).Error("Could not decode stream context")
		return err
	}

	//if we try to register stream without checking if it is legitimate we get rejected
	if p.context != nil && p.context.Priority > context.Priority {
		return errors.New("Device busy")
	}
	//if there is an existing connection we have to stop it
	if p.connection != nil {
		log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "registerStream"}).
			Info("Stopping currently open connection")
		p.connection.Close("Higher priority transmission")
	}
	p.connection = c
	p.context = context
	//we continue in a separate goroutine
	go p.doPlayFromWsConnection()
	return nil
}

func (p *play) doPlayFromWsConnection() {
	defer p.cleanup()
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "doPlayFromWsConnection"}).
			Debug("Starting write loop")
	}
	//we initialize output for signalling
	out := make(chan []byte)
	go p.connection.WriteLoop(out)

	var err error
	//play 'dong' if requested
	if p.context.PlayIntro == true {
		out <- introStartMsg
		if err = p.playFile(p.introFile); err != nil {
			log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "doPlayFromWsConnection"}).
				WithError(err).Error("Could not play dong")
			p.connection.Close("Could not play dong")
			return
		}
		out <- introEndMsg
	}

	//initialize playback device
	var dev PlaybackDevice
	if dev, err = p.factory.New(p.context.SampleRate, p.context.Channels, p.bufParams); err != nil {
		p.connection.Close("Could not initialize dev")
	}

	var buf []byte
	buf16 := make([]int16, p.context.BufferSize)

	//now we can start reading data
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "PlayFromWsConnection"}).
			Debug("Starting read loop")
	}
	go p.connection.ReadLoop()
	bin, _ := p.connection.In()
	ctrl := p.connection.Control()

	var ok bool
	var wrote int

	for {
		select {
		case buf, ok = <-bin:
			if !ok {
				if log.GetLevel() >= log.InfoLevel {
					log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "PlayFromWsConnection", "Connection": p.connection.ID()}).
						Info("Binary input channel is closed; aborting read loop")
				}
				return
			}
			convertBuffers(buf, buf16)
			var err error
			if wrote, err = dev.Write(buf16); err != nil {
				log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "PlayFromWsConnection"}).
					WithError(err).Error("Could not write buffer content to device")
				return
			}
			if log.GetLevel() >= log.DebugLevel {
				log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "PlayFromWsConnection", "wroteFrames": wrote, "readBytes": len(buf)}).
					Debug("Wrote read buffer to device")
			}
			p.context.framesWrote += wrote
			p.context.bytesRead += len(buf)
		case <-ctrl:
			return
		}
	}
}

func (p *play) playFile(filepath string) error {
	var f *os.File
	var err error
	if f, err = os.Open(filepath); err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	//initialize the device and buffers
	var dev PlaybackDevice
	if dev, err = p.factory.New(defaultSampleRate, defaultChannels, p.bufParams); err != nil {
		return err
	}

	buf := make([]byte, 8192)
	buf16 := make([]int16, 8192/sampleSizeBytes)

	var read int
	for {
		if read, err = r.Read(buf); err != nil && err != io.EOF {
			return err
		}
		if read == 0 {
			break
		}
		convertBuffers(buf, buf16)
		if _, err = dev.Write(buf16); err != nil {
			log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "playFile"}).
				WithError(err).Error("Could not write buffer content to device")
			return err
		}
	}
	return nil
}

func (p *play) cleanup() {
	if log.GetLevel() >= log.InfoLevel {
		log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "cleanup", "bytesRead": p.context.bytesRead, "framesWrote": p.context.framesWrote}).
			Info("Audio device read, write summary")
	}
	p.connection.Close("")
	p.connection = nil
	p.context = nil
}

func convertBuffers(buf []byte, buf16 []int16) {
	for i := 0; i < len(buf16); i++ {
		// for little endian
		buf16[i] = int16(binary.LittleEndian.Uint16(buf[i*2 : (i+1)*2]))
	}
}

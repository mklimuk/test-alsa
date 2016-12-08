package audio

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/mklimuk/test-alsa/config"
	"github.com/mklimuk/websocket"
)

const sampleSizeBytes = 2
const textMessage = 1

//Playback is the interface responsible for communication with audio device
type Playback interface {
	DeviceBusy() (bool, int)
	BufferSize() int
	PlaybackContext() *StreamContext
	PlayFromWsConnection(c websocket.Connection) error
	Close()
}

//StreamContext contains information about currently playing stream
type StreamContext struct {
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	Volume      int    `json:"volume"`
	Type        string `json:"type"`
	PlayIntro   bool   `json:"playIntro"`
	bytesRead   int
	framesWrote int
}

type play struct {
	bufParams  *BufferParams
	connMutex  sync.Mutex
	context    *StreamContext
	connection websocket.Connection
	introFile  string
	bufferSize int
	buf        []byte
	buf16      []int16
	device     PlaybackDevice
}

//New is the playback interface constructor
func New(conf *config.AudioConf, dev PlaybackDevice, introFile string) Playback {
	p := play{
		bufParams:  &BufferParams{BufferFrames: conf.DeviceBuffer, PeriodFrames: conf.PeriodFrames, Periods: conf.Periods},
		bufferSize: conf.ReadBuffer,
		device:     dev,
		buf16:      make([]int16, conf.ReadBuffer/sampleSizeBytes),
		buf:        make([]byte, conf.ReadBuffer),
		introFile:  introFile,
	}
	return &p
}

func (p *play) BufferSize() int {
	return p.bufferSize
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
		log.WithFields(log.Fields{"logger": "ws.audio-endpoint.audio", "method": "PlayFromWsConnection", "message": string(message)}).
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
		p.connection.Close()
	}
	p.connection = c
	p.context = context
	//we continue in a separate goroutine
	go p.doPlayFromWsConnection()
	return nil
}

func (p *play) doPlayFromWsConnection() {
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "PlayFromWsConnection"}).
			Debug("Starting read loop")
	}
	go p.connection.ReadLoop()
	bin, _ := p.connection.In()
	ctrl := p.connection.Control()

	if p.context.PlayIntro == true {
		p.playFile(p.introFile)
	}
	var ok bool
	var wrote int
	for {
		select {
		case p.buf, ok = <-bin:
			if !ok {
				if log.GetLevel() >= log.InfoLevel {
					log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "PlayFromWsConnection", "Connection": p.connection.ID()}).
						Info("Binary input channel is closed; aborting read loop")
				}
				return
			}
			convertBuffers(p.buf, p.buf16)
			var err error
			if wrote, err = p.device.Write(p.buf16); err != nil {
				log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "PlayFromWsConnection"}).
					WithError(err).Error("Could not write buffer content to device")
				return
			}
			p.context.framesWrote += wrote
			p.context.bytesRead += len(p.buf)
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

	var read int
	for {
		if read, err = r.Read(p.buf); err != nil && err != io.EOF {
			return err
		}
		if read == 0 {
			break
		}
		convertBuffers(p.buf, p.buf16)
		if _, err = p.device.Write(p.buf16); err != nil {
			log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "playFile"}).
				WithError(err).Error("Could not write buffer content to device")
			return err
		}
	}
	return nil
}

func (p *play) Close() {
	p.device.Close()
}

func convertBuffers(buf []byte, buf16 []int16) {
	for i := 0; i < len(buf16); i++ {
		// for little endian
		buf16[i] = int16(binary.LittleEndian.Uint16(buf[i*2 : (i+1)*2]))
	}
}

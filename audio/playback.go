package audio

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/mklimuk/test-alsa/config"
	"github.com/mklimuk/websocket"
)

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
	PlayFromWsConnection(c websocket.Connection)
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
	connMutex         sync.Mutex
	context           *StreamContext
	connection        websocket.Connection
	introFile         string
	samplesBufferSize int
	bufParams         *BufferParams
	factory           DeviceFactory
	dev               PlaybackDevice
}

//New is the playback interface constructor
func New(conf *config.AudioConf, factory DeviceFactory, introFile string) Playback {
	p := play{
		factory:           factory,
		bufParams:         &BufferParams{BufferFrames: conf.DeviceBuffer, PeriodFrames: conf.PeriodFrames, Periods: conf.Periods},
		introFile:         introFile,
		samplesBufferSize: conf.ReadBuffer,
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

/*PlayFromWsConnection streams audio data from a websocket connection into an audio device.
The first message received must be a text message containing stream context (see StreamContext type) in a JSON format.
*/
func (p *play) PlayFromWsConnection(c websocket.Connection) {
	p.connMutex.Lock()
	defer p.connMutex.Unlock()

	//the first message contains information about the stream context
	var err error
	var mt int
	var message []byte
	if mt, message, err = c.ReadMessage(); err != nil {
		log.WithFields(log.Fields{"logger": "ws.audio-endpoint.audio", "method": "PlayFromWsConnection"}).
			WithError(err).Error("Websocket read error encountered")
		c.CloseWithReason(websocket.CloseProtocolError, "Could not read from connection")
		return
	}
	if mt != textMessage {
		c.CloseWithReason(websocket.CloseInvalidFramePayloadData, "Wrong message type; first message should be a text one")
		return
	}

	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "ws.audio-endpoint.audio", "method": "PlayFromWsConnection", "message": message}).
			Debug("Received stream context message")
	}
	var context = new(StreamContext)
	if err = json.Unmarshal(message, context); err != nil {
		log.WithFields(log.Fields{"logger": "ws.audio-endpoint.audio", "method": "PlayFromWsConnection"}).
			WithError(err).Error("Could not decode stream context")
		c.CloseWithReason(websocket.CloseInvalidFramePayloadData, "Could not decode stream context")
		return
	}

	//stream with lower priority will get rejected
	if p.context != nil && p.context.Priority > context.Priority {
		c.CloseWithReason(websocket.CloseTryAgainLater, "Device busy")
		return
	}

	//if there is an existing connection we have to stop it
	if p.connection != nil {
		log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "registerStream"}).
			Info("Stopping currently open connection")
		p.connection.CloseWithCode(websocket.CloseGoingAway)
	}
	p.connection = c
	p.context = context
	//we continue in a separate goroutine
	go p.doPlayFromWsConnection()
	return
}

func (p *play) doPlayFromWsConnection() {
	defer p.cleanup()

	var err error
	//play intro (ding-dong) if requested
	if p.context.PlayIntro == true {
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "doPlayFromWsConnection", "Connection": p.connection.ID()}).
				Debug("Playing intro file")
		}
		p.connection.WriteMessage(websocket.TextMessage, introStartMsg)
		if err = p.PlayFile(p.introFile); err != nil {
			log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "doPlayFromWsConnection", "Connection": p.connection.ID()}).
				WithError(err).Warn("Could not play intro file")
			msg, _ := json.Marshal(SignallingMsg{"playback:intro:warn", "Could not play intro"})
			p.connection.WriteMessage(websocket.TextMessage, msg)
			//we continue anyway
		}
		p.connection.WriteMessage(websocket.TextMessage, introEndMsg)
	}

	//initialize playback device
	if p.dev, err = p.factory.New(p.context.SampleRate, p.context.Channels, p.bufParams); err != nil {
		p.connection.CloseWithReason(websocket.CloseInternalServerErr, "Could not initialize audio device")
		return
	}
	devbuf := make(chan []int16, p.samplesBufferSize)
	var deverr chan error
	defer close(devbuf)

	//prepare connection read buffers
	var buf []byte
	buf16 := make([]int16, p.context.BufferSize)

	//start the connection read routine
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "doPlayFromWsConnection", "Connection": p.connection.ID()}).
			Debug("Starting read loop")
	}
	go p.connection.ReadLoop()
	bin, _ := p.connection.In()

	var ok bool
	var writing bool
	for {
		select {
		case buf, ok = <-bin: //binary audio data from the websocket
			if !ok {
				if log.GetLevel() >= log.InfoLevel {
					log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "doPlayFromWsConnection", "Connection": p.connection.ID()}).
						Info("Binary input channel is closed; aborting read loop")
				}
				return
			}
			if log.GetLevel() >= log.DebugLevel {
				log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "doPlayFromWsConnection", "readBytes": len(buf)}).
					Debug("Read bytes from connection")
			}
			p.context.bytesRead += len(buf)

			//convert to int16 and push to the output buffer
			convertBuffers(buf, buf16)
			devbuf <- buf16

			//if the buffer is full we start sending audio to the audio device
			if !writing && len(devbuf) == cap(devbuf) {
				if log.GetLevel() >= log.InfoLevel {
					log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "doPlayFromWsConnection", "Connection": p.connection.ID()}).
						Info("Starting audio device write routine")
				}
				deverr = p.dev.WriteAsync(devbuf)
				writing = true
			}
		case <-p.connection.Control(): //connection control chanel
			if log.GetLevel() >= log.DebugLevel {
				log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "doPlayFromWsConnection"}).
					Debug("Received connection close signal")
			}
			return
		case err = <-deverr: //errors from audio device
			log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "doPlayFromWsConnection"}).
				WithError(err).Error("Could not write buffer content to device")
			p.connection.CloseWithReason(websocket.CloseInternalServerErr, "Could not write buffer content to audio device")
			return
		}
	}
}

//PlayFile sends contents of the file represented by 'filepath' to Alsa audio device.
//This method uses default sample rate and channels number.
func (p *play) PlayFile(filepath string) error {
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
	defer dev.Close()
	dev.WriteSync(r)
	return nil
}

func (p *play) cleanup() {
	var wrote int
	if p.dev != nil {
		wrote = p.dev.FramesWrote()
		p.dev.Close()
	}
	if log.GetLevel() >= log.InfoLevel {
		log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "cleanup", "bytesRead": p.context.bytesRead, "framesWrote": wrote}).
			Info("Audio device read, write summary")
	}
	p.connection.CloseWithCode(websocket.CloseNormalClosure)
	p.connection = nil
	p.context = nil
}

func convertBuffers(buf []byte, buf16 []int16) {
	for i := 0; i < len(buf16); i++ {
		// for little endian
		buf16[i] = int16(binary.LittleEndian.Uint16(buf[i*2 : (i+1)*2]))
	}
}

package tts

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	t "github.com/mklimuk/husar/tts"

	log "github.com/Sirupsen/logrus"
)

//DefaultAudioLength is the default length that should be safe enough not to create playback conflicts
const DefaultAudioLength int = 90

//Connector is used to dialog with TTS service in order to generate audio files from text.
type Connector interface {
	GetAudio(id string, text string) (respBody []byte, length int, err error)
}

//NewConnector is Connector's constructor
func NewConnector(ttsURL string) Connector {
	c := connector{}
	cl := http.Client{}
	c.client = &cl
	c.ttsURL = ttsURL
	return Connector(&c)
}

type connector struct {
	client *http.Client
	ttsURL string
}

func (c *connector) GetAudio(id string, text string) (respBody []byte, length int, err error) {
	gen := t.GenerateRequest{
		ID:   id,
		Time: time.Now(),
		Text: text,
	}
	//context logger
	clog := log.WithFields(log.Fields{"logger": "tts.connector", "method": "GetAudio", "id": id})

	if log.GetLevel() >= log.DebugLevel {
		clog.WithField("text", text).
			Debug("Sending TTS generate speech request.")
	}

	var jsonReq []byte
	if jsonReq, err = json.Marshal(gen); err != nil {
		clog.WithError(err).Error("Could not marshal request to json.")
		return
	}

	req, err := http.NewRequest("POST", c.ttsURL, bytes.NewBuffer(jsonReq))

	resp, err := c.client.Do(req)
	if err != nil {
		clog.WithError(err).Error("Error while performing HTTP request.")
		length = -1
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		clog.WithField("status", resp.StatusCode).Error("TTS service responded with an error.")
		err = errors.New("TTS service responded with an error")
		//TODO parse error
		return
	}

	var l string
	if l = resp.Header.Get("Audio-Len"); len(l) == 0 {
		clog.Error("Unknown audio file length. Assuming default.")
		length = DefaultAudioLength
	} else {
		if length, err = strconv.Atoi(l); err != nil {
			clog.WithField("length", l).WithError(err).Error("Could not parse audio length header")
		}
	}

	if respBody, err = ioutil.ReadAll(resp.Body); err != nil {
		clog.WithError(err).Error("Could not read response body.")
		return
	}

	if log.GetLevel() >= log.DebugLevel {
		clog.Debug("Received TTS response.")
	}
	return
}

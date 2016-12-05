package tts

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const wavExtension string = "wav"
const oggExtension string = "ogg"

//Speech is an interface to communication with the TTS service
type Speech interface {
	GenerateSpeech(gr GenerateRequest) (string, int, error)
}

type speech struct {
	tmpPath      string
	ttsEngine    string
	oggConverter string
	tempo        string
	dongPath     string
}

//NewSpeech speech constructor
func NewSpeech(tmpPath string, ttsEngine string, oggConverter string, tempo string, dongPath string) Speech {
	tmp := tmpPath
	if !strings.HasSuffix(tmpPath, "/") {
		tmp += "/"
	}
	s := speech{tmp, ttsEngine, oggConverter, tempo, dongPath}
	return &s
}

func (s *speech) GenerateSpeech(gr GenerateRequest) (string, int, error) {

	clog := log.WithFields(log.Fields{"logger": "tts.speech", "method": "GenerateSpeech", "requestId": gr.ID})
	log.Infoln("Log level: ", log.GetLevel())
	if log.GetLevel() >= log.InfoLevel {
		clog.WithFields(log.Fields{"requestTime": gr.Time, "text": gr.Text}).
			Info("Generating speech.")
	}

	wavPath := filepath.FromSlash(fmt.Sprintf("%s%s.%s", s.tmpPath, gr.ID, wavExtension))
	oggPath := filepath.FromSlash(fmt.Sprintf("%s%s.%s", s.tmpPath, gr.ID, oggExtension))

	if log.GetLevel() >= log.DebugLevel {
		clog.WithFields(log.Fields{"wav": wavPath, "ogg": oggPath}).
			Debug("File paths.")
	}

	//	generating file
	var out []byte
	var err error
	if out, err = exec.Command(s.ttsEngine, gr.Text, wavPath, s.tempo, s.dongPath).CombinedOutput(); err != nil {
		clog.WithField("out", string(out)).Error("Synthesizer output.")
		return "", -1, err
	}
	// we will remove the file afterwards
	defer os.Remove(wavPath)

	if log.GetLevel() >= log.DebugLevel {
		clog.WithField("out", string(out)).Debug("Synthesizer output.")
	}

	// parse headers
	var l int
	if l, err = audioLength(wavPath); err != nil {
		clog.WithField("l", l).WithError(err).
			Warn("Could not read length from file.")
	}

	//	converting to ogg
	if out, err = exec.Command(s.oggConverter, wavPath, oggPath).CombinedOutput(); err != nil {
		clog.WithField("out", string(out)).Error("Ogg converter output.")
		return "", -1, err
	}

	if log.GetLevel() >= log.DebugLevel {
		clog.WithField("out", string(out)).Debug("Ogg converter output.")
	}

	return oggPath, l, err
}

func audioLength(path string) (length int, err error) {
	// open wav file
	var file *os.File

	if file, err = os.Open(path); err != nil {
		log.WithFields(log.Fields{"logger": "tts.api", "method": "audioLength"}).
			WithError(err).Error("Could not open file")
		return -1, err
	}
	defer file.Close()

	// bytes per second is byte 29-32 of the Header
	// bps = (Sample Rate * BitsPerSample * Channels) / 8
	bytesPerSec := make([]byte, 4)
	if _, err = file.ReadAt(bytesPerSec, 28); err != nil {
		log.WithFields(log.Fields{"logger": "tts.api", "method": "audioLength"}).
			WithError(err).Error("Error reading bps")
		return -1, err
	}
	var bps int32
	if err = binary.Read(bytes.NewReader(bytesPerSec), binary.LittleEndian, &bps); err != nil {
		log.WithFields(log.Fields{"logger": "tts.api", "method": "audioLength"}).
			WithError(err).Error("Error decoding bps")
		return -1, err
	}

	// data size should be byte 41-44 but it seems to be bytes 43-46
	dataSize := make([]byte, 4)
	if _, err = file.ReadAt(dataSize, 42); err != nil {
		log.WithFields(log.Fields{"logger": "tts.api", "method": "audioLength"}).
			WithError(err).Error("Error reading data size")
		return -1, err
	}
	log.Debugln(dataSize)
	var size int32
	if err = binary.Read(bytes.NewReader(dataSize), binary.LittleEndian, &size); err != nil {
		log.WithFields(log.Fields{"logger": "tts.api", "method": "audioLength"}).
			WithError(err).Error("Error decoding data size")
		return -1, err
	}
	return int(size / bps), nil
}

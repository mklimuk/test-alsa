package audio

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/mklimuk/husar/client/tts"

	log "github.com/Sirupsen/logrus"
)

const perms os.FileMode = 0644

// AudioExtension defines audio files extension used by the system
const AudioExtension string = "ogg"

// LenExtension defines length file used by the system
const LenExtension string = "len"

//DefaultDuration is default audio duration in seconds
const DefaultDuration = 90

/*Catalog handles saving and reading audio files from hdd */
type Catalog interface {
	GetID(text string) string
	GetPath(ID string) (string, error)
	Generate(text string) (ID string, exists bool, duration int, approxLen bool, err error)
	Get(ID string) (audio []byte, err error)
}

//NewCatalog is a constructor for Catalog
func NewCatalog(path string, t *tts.Connector) Catalog {
	c := cat{t: t, catalogPath: path}
	return Catalog(&c)
}

type cat struct {
	catalogPath string
	t           *tts.Connector
}

func (c *cat) GetID(text string) string {
	hash := sha1.Sum([]byte(text))
	id := hex.EncodeToString(hash[:])
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "lcs.audio", "method": "GetID", "id": id, "text": text}).
			Debug("Generated id.")
	}
	return id
}

func (c *cat) Generate(text string) (ID string, exists bool, duration int, approxDur bool, err error) {
	ID = c.GetID(text)
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "lcs.audio", "method": "Generate", "id": ID, "text": text}).
			Debug("Getting file from catalog.")
	}
	var au []byte
	if _, err = os.Stat(fmt.Sprintf("%s/%s.%s", c.catalogPath, ID, AudioExtension)); os.IsNotExist(err) {
		exists = false
		log.WithFields(log.Fields{"logger": "lcs.audio", "method": "Generate", "id": ID}).
			Info("File not found in the catalog. Calling TTS service.")

		if au, duration, err = (*c.t).GetAudio(ID, text); err != nil {
			log.WithFields(log.Fields{"logger": "lcs.audio", "method": "Generate", "id": ID}).
				WithError(err).Error("Could not get audio from TTS.")
			return ID, exists, -1, true, err
		}

		d := strconv.Itoa(duration)
		ioutil.WriteFile(fmt.Sprintf("%s/%s.%s", c.catalogPath, ID, AudioExtension), au, perms)
		ioutil.WriteFile(fmt.Sprintf("%s/%s.%s", c.catalogPath, ID, LenExtension), []byte(d), perms)
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "lcs.audio", "method": "Generate", "id": ID, "exists": exists, "duration": duration, "approxDur": approxDur}).
				Debug("Returning generated file properties.")
		}
		return ID, exists, duration, false, err
	} else if err != nil {
		log.WithFields(log.Fields{"logger": "lcs.audio", "method": "Generate", "id": ID}).
			WithError(err).Error("Unknown filesystem error.")
		return ID, exists, -1, true, err
	}
	exists = true
	duration, err = c.readLen(ID)
	approxDur = err != nil
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "lcs.audio", "method": "Generate", "id": ID, "exists": exists, "duration": duration, "approxDur": approxDur}).
			Debug("Returning generated file properties.")
	}
	return ID, exists, duration, approxDur, err
}

func (c *cat) readLen(ID string) (length int, err error) {
	var l []byte
	if l, err = ioutil.ReadFile(fmt.Sprintf("%s/%s.%s", c.catalogPath, ID, LenExtension)); err != nil {
		log.WithFields(log.Fields{"logger": "lcs.audio", "method": "readLen", "id": ID}).
			WithError(err).Error("Filesystem error while reading length file. We will assume default length.")
		return DefaultDuration, err
	}
	length, _ = strconv.Atoi(string(l))
	return length, err
}

func (c *cat) Get(ID string) (audio []byte, err error) {
	var au []byte
	if au, err = ioutil.ReadFile(fmt.Sprintf("%s/%s.%s", c.catalogPath, ID, AudioExtension)); os.IsNotExist(err) {
		log.WithFields(log.Fields{"logger": "lcs.audio", "method": "Generate", "id": ID}).
			Info("File not found in the catalog.")
		return audio, err
	} else if err != nil {
		log.WithFields(log.Fields{"logger": "lcs.audio", "method": "Generate", "id": ID}).
			WithError(err).Error("Unknown filesystem error.")
		return audio, err
	}
	return au, err
}

func (c *cat) GetPath(ID string) (string, error) {
	path := fmt.Sprintf("%s/%s.%s", c.catalogPath, ID, AudioExtension)
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "lcs.audio", "method": "GetPath", "path": path}).
			Debug("Checking path.")
	}
	var err error
	if _, err = os.Stat(path); os.IsNotExist(err) {
		log.WithFields(log.Fields{"logger": "lcs.audio", "method": "GetPath", "id": ID}).
			Info("File not found in the catalog. Could not generate path.")
		return "", err
	}
	if err != nil {
		return "", err
	}
	return path, nil
}

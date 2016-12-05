package audio

import (
	"github.com/mklimuk/husar/device"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

//Controller is responsible for interactions with audio endpoints
type Controller interface {
	Upload(audio []byte, title string, devices []*device.Device, id string, concurrent int) ([]string, []int, []error)
	Play(IPs []*device.Device, id string, title string, timestamp time.Time) ([]string, []int, []error)
	SetVolume(devices []*device.Device, volume int) ([]string, []int, []error)
	UploadAndPlay(audio []byte, title string, devices []*device.Device, id string, concurrent int, timestamp time.Time)
}

//NewController is a Controller's constructor
func NewController() Controller {
	c := con{}
	cl := http.Client{}
	c.client = &cl
	return Controller(&c)
}

type con struct {
	client *http.Client
}

func (c *con) UploadAndPlay(audio []byte, title string, devices []*device.Device, id string, concurrent int, timestamp time.Time) {
	clog := log.WithFields(log.Fields{"logger": "lcs.audio", "method": "UploadAndPlay"})
	if log.GetLevel() >= log.DebugLevel {
		clog.WithFields(log.Fields{"id": id, "devices": devices}).Debug("Uploading audio for playback.")
	}
	c.Upload(audio, title, devices, id, concurrent)
	if log.GetLevel() >= log.DebugLevel {
		clog.WithFields(log.Fields{"id": id, "devices": devices}).Debug("Playing audio.")
	}
	c.Play(devices, id, title, timestamp)
}

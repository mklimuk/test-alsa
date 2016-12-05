package device

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
)

//Zone represents a playback zone
type Zone interface {
	HasDevices() bool
	Run()
	Stop()
	GetDevices() []*Device
	AddDevice(d *Device)
	RemoveDevice(deviceID string) error
}

//NewZone is Zone constructor
func NewZone(ID string) Zone {
	z := zone{}
	z.devices = []*Device{}
	return Zone(&z)
}

type zone struct {
	devices       []*Device
	updateControl chan bool
	refreshFreq   time.Duration
}

func (z *zone) HasDevices() bool {
	return len(z.devices) > 0
}

func (z *zone) AddDevice(d *Device) {
	z.devices = append(z.devices, d)
}

func (z *zone) RemoveDevice(deviceID string) error {
	for i, d := range z.devices {
		if d.ID == deviceID {
			z.devices = append(z.devices[:i], z.devices[i+1:]...)
			return nil
		}
	}
	return errors.New("Device not found")
}

func (z *zone) Run() {
	z.updateControl = make(chan bool)
	go func() {
		for {
			select {
			case <-z.updateControl:
				return
			case <-time.After(z.refreshFreq):
				z.checkDevices()
			}
		}
	}()
}

func (z *zone) Stop() {
	z.updateControl <- true
}

func (z *zone) GetDevices() []*Device {
	return z.devices
}

func (z *zone) checkDevices() {
	clog := log.WithFields(log.Fields{"logger": "lcs.device", "method": "checkDevices"})
	clog.Info("Checking bound devices")
	for _, d := range z.devices {
		if d.BindTimeout > 0 {
			d.BindTimeout--
		}
		if d.Bind && d.BindTimeout == 0 {
			clog.WithField("device", d.ID).Info("Bind timeout reached for device")
			d.Bind = false
		}

	}

}

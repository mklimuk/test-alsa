package device

import (
	"errors"
	"fmt"

	"github.com/mklimuk/husar/event"

	log "github.com/Sirupsen/logrus"
)

/*Registry keeps information about registered devices
and available playback zones */
type Registry interface {
	CreateZone(zoneID string)
	GetZone(zoneID string) *Zone
	Bind(device *Device)
	AddToZone(zoneID string, deviceID string) error
	RemoveFromZone(zoneID string, deviceID string) error
	Unregister(deviceID string) error
	GetDevices(zoneID string) ([]*Device, error)
	GetDevice(deviceID string) (*Device, error)
}

//NewRegistry is a Registry constructor
func NewRegistry(b *event.Bus, timeout int) Registry {
	r := reg{timeout: timeout}
	r.zones = make(map[string]*Zone)
	r.devices = make(map[string]*Device)
	// initialize event bus listeners
	b.Subscribe(event.DeviceBind, r.Bind)
	b.Subscribe(event.DeviceAddToZone, r.AddToZone)
	return Registry(&r)
}

type reg struct {
	zones   map[string]*Zone
	devices map[string]*Device
	timeout int
}

func (r *reg) CreateZone(zoneID string) {
	z := NewZone(zoneID)
	r.zones[zoneID] = &z
}

func (r *reg) GetZone(zoneID string) *Zone {
	return r.zones[zoneID]
}

func (r *reg) Bind(dev *Device) {
	if dev == nil {
		log.WithFields(log.Fields{"logger": "lcs.device.registry", "method": "Bind"}).
			Warn("Received empty device pointer. Ignoring...")
		return
	}
	d := r.devices[dev.ID]
	var isNew bool
	if isNew = d == nil; isNew {
		log.WithFields(log.Fields{"logger": "lcs.device.registry", "method": "Bind"}).
			Info("Binding new device.")
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "lcs.device.registry", "method": "Bind", "device": fmt.Sprintf("%+v", dev)}).
				Debug("Device info")
		}
		d = dev
	}
	d.Live = true
	d.Bind = true
	d.BindTimeout = r.timeout
	r.devices[d.ID] = d
	if isNew {
		log.WithFields(log.Fields{"logger": "lcs.device.registry", "method": "Bind", "device": d.ID}).
			Info("Adding device to defined zones.")
		var err error
		for _, z := range d.Zones {
			if err = r.AddToZone(z, dev.ID); err != nil {
				log.WithFields(log.Fields{"logger": "lcs.device.registry", "method": "Bind", "device": d.ID, "zone": z}).
					WithError(err).Error("Could not add device to zone")
			}
		}
	}
}

func (r *reg) Unregister(deviceID string) error {
	var d *Device
	if d = r.devices[deviceID]; d == nil {
		return errors.New("Device not found")
	}
	log.WithFields(log.Fields{"logger": "lcs.device.registry", "method": "Unregister", "device": deviceID}).
		Info("Unregistering device")
	for _, zid := range d.Zones {
		r.RemoveFromZone(zid, deviceID)
	}
	delete(r.devices, deviceID)
	return nil
}

func (r *reg) AddToZone(zoneID string, deviceID string) error {
	var d *Device
	if d = r.devices[deviceID]; d == nil {
		return errors.New("Device not found")
	}
	d.Zones = append(d.Zones, zoneID)
	log.WithFields(log.Fields{"logger": "lcs.device.registry", "method": "AddToZone", "zone": zoneID, "device": deviceID}).
		Info("Adding device to zone")
	(*r.zones[zoneID]).AddDevice(d)
	return nil
}

func (r *reg) RemoveFromZone(zoneID string, deviceID string) error {
	z := r.zones[zoneID]
	if z == nil {
		return errors.New("Zone not found")
	}
	log.WithFields(log.Fields{"logger": "lcs.device.registry", "method": "AddToZone", "zone": zoneID, "device": deviceID}).
		Info("Removing device from zone")
	return (*z).RemoveDevice(deviceID)
}

func (r *reg) GetDevices(zoneID string) ([]*Device, error) {
	z := r.zones[zoneID]
	if z == nil {
		return nil, errors.New("Zone not found")
	}
	return (*z).GetDevices(), nil
}

func (r *reg) GetDevice(deviceID string) (*Device, error) {
	var d *Device
	if d = r.devices[deviceID]; d == nil {
		return nil, errors.New("Device not found")
	}
	return d, nil
}

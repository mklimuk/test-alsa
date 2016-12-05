package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/mklimuk/husar/annon"
	"github.com/mklimuk/husar/audio"
	aud "github.com/mklimuk/husar/client/audio"
	"github.com/mklimuk/husar/device"
	"github.com/mklimuk/husar/event"
	"github.com/mklimuk/husar/msg"
	"github.com/mklimuk/husar/queue"

	log "github.com/Sirupsen/logrus"
)

// Playback is a service handling playback activity
type Playback interface {
}

type playback struct {
	conn     *aud.Controller
	dev      *device.Registry
	cat      *audio.Catalog
	an       *annon.Store
	b        *event.Bus
	q        *queue.Manager
	busyLock sync.Mutex
	busy     map[string]bool
}

//NewPlayback is a playback service constructor
func NewPlayback(d *device.Registry, cat *audio.Catalog, b *event.Bus, an *annon.Store, conn *aud.Controller, q *queue.Manager) Playback {
	p := playback{dev: d, cat: cat, b: b, an: an, conn: conn, q: q, busy: map[string]bool{}}
	b.Subscribe(event.AudioPlay, p.Play)
	b.Subscribe(event.PlaybackTrigger, p.Trigger)
	b.Subscribe(event.SetVolume, p.setVolume)
	return Playback(&p)
}

func (s *playback) setVolume(req string) error {
	clog := log.WithFields(log.Fields{"logger": "lcs.playback", "method": "setVolume"})
	m := msg.ParseVolumeMessage(req)
	var err error
	var d []*device.Device
	if d, err = (*s.dev).GetDevices(m.QueueID); err != nil {
		clog.WithError(err).Error("Could not get devices. Aborting.")
		return err
	}
	if len(d) > 0 {
		clog.Info("Found devices. Setting volume.")
	}
	(*s.conn).SetVolume(d, m.Volume)
	(*s.q).SetVolume(m.QueueID, m.Volume)
	go s.b.Publish(event.SetVolumeAck, event.SetVolumeAck, req)
	return nil
}

func (s *playback) Trigger(req string) error {
	m := msg.ParseTriggerMessage(req)
	clog := log.WithFields(log.Fields{"logger": "lcs.playback", "method": "Trigger", "msg": fmt.Sprintf("%+v", m)})
	if s.busy[m.QueueID] {
		clog.WithField("queue", m.QueueID).Info("The queue is currently busy. Ignoring playback request.")
		return nil
	}
	var err error
	var e *queue.Event
	clog.Info("Getting event")
	if e, err = (*s.q).GetEvent(m.QueueID, m.EventID); err != nil || e == nil {
		clog.WithError(err).Error("Error getting event")
		return err
	}
	clog.Info("Muting first queue event if required")
	(*s.q).MuteFirstIfRequired(m.QueueID, time.Now().Add(*e.Duration))
	if err = s.Play(m.QueueID, e); err != nil {
		clog.WithError(err).Error("Error playing sound")
		return err
	}
	return nil
}

func (s *playback) Play(zoneID string, ev *queue.Event) error {
	clog := log.WithFields(log.Fields{"logger": "lcs.playback", "method": "Play"})
	var err error
	var d []*device.Device
	if d, err = (*s.dev).GetDevices(zoneID); err != nil {
		clog.WithError(err).Error("Could not get devices. Aborting.")
		return err
	}
	if len(d) == 0 {
		clog.Warn("Found no devices. Ignoring playback request.")
		return nil
	}
	var a *annon.Announcement
	if a, err = (*s.an).Get(ev.AnnonID); err != nil {
		clog.WithField("annon", ev.AnnonID).WithError(err).
			Error("Could not get announcement. Aborting.")
		return err
	}

	if a.Audio == nil {
		// this is our last chance to generate the audio
		var id string
		var len int
		var approx bool
		if id, _, len, approx, err = (*s.cat).Generate(a.Text.TtsText); err != nil {
			clog.WithField("announcement", a.ID).WithError(err).Error("Could not generate speech.")
		}
		a.Audio = &annon.Audio{FileID: id, Duration: len, ApproxLen: approx}
		(*s.an).Save(a)
	}

	var au []byte
	if au, err = (*s.cat).Get(a.Audio.FileID); err != nil {
		clog.WithField("fileID", a.Audio.FileID).WithError(err).
			Error("Could not get audio content. Aborting.")
		return err
	}
	if s.busy[zoneID] {
		clog.WithField("queue", zoneID).Info("The queue is currently busy. Ignoring playback request.")
		return nil
	}

	// we build the playback event message before calling the goroutine
	// to be able to return an error in case of problems.
	m := msg.IDMessage{ID: ev.ID}
	var res []byte
	if res, err = json.Marshal(m); err != nil {
		clog.WithError(err).Error("Error marshalling event")
		return err
	}

	// parallel playback events triggers
	go func() {

		s.busyLock.Lock()
		s.busy[zoneID] = true
		s.busyLock.Unlock()

		// publishing playback start to the bus
		clog.Info("Triggering playback:start event")
		s.b.Publish(event.PlaybackStart, event.PlaybackStart, string(res))

		// sleeping for the duration of playback
		time.Sleep(*ev.Duration)
		s.busyLock.Lock()
		s.busy[zoneID] = false
		s.busyLock.Unlock()
		clog.Info("Triggering playback:end event")
		s.b.Publish(event.PlaybackEnd, event.PlaybackEnd, string(res))
	}()

	clog.Info("Sending audio to devices and triggering playback")
	(*s.conn).UploadAndPlay(au, a.ID, d, ev.ID, len(d), *ev.PlaybackStart)

	return nil
}

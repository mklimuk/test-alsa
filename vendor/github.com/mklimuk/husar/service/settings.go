package service

import (
	"encoding/json"

	"github.com/mklimuk/husar/event"
	"github.com/mklimuk/husar/msg"
	"github.com/mklimuk/husar/train"

	log "github.com/Sirupsen/logrus"
)

func (s *timetable) updateSettings(req string) {
	clog := log.WithFields(log.Fields{"logger": "lcs.timetable", "method": "updateSettings"})
	var m *msg.AudioSettingsUpdate
	var err error
	if m, err = msg.ParseSettingsUpdate(req); m == nil {
		clog.WithError(err).Error("Could not parse message")
		return
	}
	var t *train.Train
	if t, err = s.store.Get(m.TrainID); err != nil {
		clog.WithError(err).Error("Could not load train from the database")
		return
	}
	// we get settings for this train
	var set *[]train.Settings
	if set, err = s.store.GetAllSettings(m.TrainID); err != nil {
		clog.WithError(err).Error("Could not load train from the database")
		return
	}
	if len(*set) > 0 {
		t.Settings = &(*set)[0]
	}

	if t.Settings == nil {
		initSettings(t)
	}
	if t.Settings.Audio == nil && m.Settings != nil {
		t.Settings.Audio = &m.Settings
	} else {
		for k, v := range m.Settings {
			(*t.Settings.Audio)[k] = v
		}
	}
	if t.Settings.Lang == nil && m.Lang != nil {
		t.Settings.Lang = &m.Lang
	} else {
		for l, v := range m.Lang {
			current := (*t.Settings.Lang)[l]
			current.Enabled = v.Enabled
			clog.WithField("enabled", current.Enabled).Debug("Enabled")
			if v.I18n != nil {
				if current.I18n == nil {
					current.I18n = &map[train.AnnonOption]string{}
				}
				for opt, val := range *v.I18n {
					(*current.I18n)[opt] = val
				}
			}
			(*t.Settings.Lang)[l] = current
		}
	}
	if t.Settings.Overrides == nil && m.Overrides != nil {
		t.Settings.Overrides = &m.Overrides
	} else {
		for k, v := range m.Overrides {
			(*t.Settings.Overrides)[k] = v
		}
	}

	clog.WithField("settings", t.Settings.ID).Info("Saving settings")
	if _, err = s.store.SaveSettings(t.Settings); err != nil {
		clog.WithError(err).Error("Error saving settings to the database")
	}
	var res []byte
	if res, err = json.Marshal(processTrain(*t)); err != nil {
		clog.WithError(err).Error("Could not parse update message")
	}
	go s.b.Publish(event.TrainUpdate, event.TrainUpdate, string(res))
}

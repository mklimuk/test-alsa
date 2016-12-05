package sync

import (
	"github.com/mklimuk/husar/train"

	log "github.com/Sirupsen/logrus"
)

/*
Settings runs an update loop for train settings change events that will come essentially from
user interaction with the front end application. It calls train announcements update whenever a change is detected.
*/
func (s *syncservice) settings() {
	s.settingsSync = make(chan bool)
	go func() {
		errorCount := 10
		for {
			select {
			case <-s.settingsSync:
				log.WithField("logger", "sync.settings").Infoln("Received quit signal. Exiting.")
				return
			default:
				errorCount = s.processSettingsChange(errorCount)
			}
		}
	}()
}

func (s *syncservice) processSettingsChange(errorCount int) int {
	var next, old *train.Settings
	var err error
	if next, old, err = s.tr.NextSettingsChange(); err != nil {
		log.WithFields(log.Fields{
			"logger": "sync.settings",
			"error":  err,
		}).Error("Error while obtaining new settings changes.")
		// 10 consecutive errors will cause the loop to exit
		errorCount--
		if errorCount == 0 {
			log.WithFields(log.Fields{
				"logger": "sync.settings",
			}).Errorln("Limit number of consecutive errors reached. Exiting.")
			s.settingsSync <- true
		}
		return errorCount
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{
			"logger": "sync.settings",
			"next":   next,
			"old":    old,
		}).Debug("New settings change received.")
	}
	// errors are gobbled by the service itself as we cannot do much about them
	s.a.ProcessSettingsChange(next, old)
	return 10
}

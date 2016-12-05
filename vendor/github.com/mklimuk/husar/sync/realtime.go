package sync

import (
	"github.com/mklimuk/husar/train"

	log "github.com/Sirupsen/logrus"
)

/*
realtime runs an update loop for realtime timetable events.
It updates announcements for a given train whenever a change of realtime state is detected.
*/
func (s *syncservice) realtime() {
	log.WithField("logger", "sync.realtime").Infoln("Enabling realtime synchronization.")
	s.liveSync = make(chan bool)
	go func() {
		errorCount := 10
		for {
			select {
			case <-s.liveSync:
				log.WithField("logger", "sync.realtime").Info("Received quit signal. Exiting.")
				return
			default:
				errorCount = s.processRealtimeChange(errorCount)
			}
		}
	}()
}

func (s *syncservice) processRealtimeChange(errorCount int) int {
	var next, old *train.Realtime
	var err error
	if next, old, err = s.tr.NextRealtimeChange(); err != nil {
		log.WithFields(log.Fields{
			"logger": "sync.realtime",
			"error":  err,
		}).Error("Error while obtaining new realtime changes.")
		// 10 consecutive errors will cause the loop to exit
		errorCount--
		if errorCount == 0 {
			log.WithFields(log.Fields{
				"logger": "sync.realtime",
			}).Errorln("Limit number of consecutive errors reached. Exiting.")
			s.liveSync <- true
		}
		return errorCount
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{
			"logger": "sync.realtime",
			"next":   next,
			"old":    old,
		}).Debug("New realtime change received.")
	}
	// errors are gobbled by the service itself as we cannot do much about them
	s.a.ProcessRealtimeChange(next, old)
	return 10
}

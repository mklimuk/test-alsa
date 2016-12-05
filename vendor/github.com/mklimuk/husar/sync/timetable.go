package sync

import (
	"time"

	"github.com/mklimuk/husar/config"

	log "github.com/Sirupsen/logrus"
)

/*
timetable runs an update loop for scheduled trains that generates announcements
for a window in time defined by windowStart and windowEnd (understood as relative time
in minutes starting from time.Now())
*/
func (s *syncservice) timetable() {
	s.scheduledSync = make(chan bool)
	// first iteration is triggered instantly
	s.processScheduledTrains()
	go func() {
		for {
			select {
			case <-s.scheduledSync:
				log.WithField("logger", "sync.timetable").Infoln("Received quit signal. Exiting.")
				return
			case <-time.After(s.windowSize):
				s.processScheduledTrains()
			}
		}
	}()
}

func (s *syncservice) processScheduledTrains() {
	// generate for a period of time in the future
	start := time.Now().In(config.Timezone).Add(s.windowAhead)
	end := start.Add(s.windowSize)
	log.WithFields(log.Fields{
		"logger": "sync.timetable",
		"start":  start,
		"end":    end,
	}).Info("Generating announcements for scheduled trains.")
	s.a.GenerateForPeriod(&start, &end)
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{
			"logger": "sync.timetable",
		}).Debug("Entering sleep interval.")
	}
}

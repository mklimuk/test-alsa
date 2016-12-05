package sync

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mklimuk/husar/annon"
	"github.com/mklimuk/husar/config"
	"github.com/mklimuk/husar/event"

	log "github.com/Sirupsen/logrus"
)

/*
Announcements runs an update loop for announcements timetable event.
It calls playback queue updates.
*/
func (s *syncservice) EnableAnnouncementsSync() {
	now := time.Now().In(config.Timezone)
	end := now.Add(s.windowSize)
	s.currentTail = end
	log.WithFields(log.Fields{"logger": "lcs.sync", "method": "loadInitial"}).
		Info("Loading initial announcements.")
	(*s).loadInitial(&now, &end)
	s.annonSync = make(chan bool)
	go func() {
		errorCount := 10
		for {
			select {
			case <-s.annonSync:
				log.WithField("logger", "sync.annons").Infoln("Received quit signal. Exiting.")
				return
			default:
				errorCount = s.processChange(errorCount)
			}
		}
	}()
}

func (s *syncservice) loadInitial(start *time.Time, end *time.Time) {
	var annons []annon.Announcement
	var err error
	if annons, err = s.an.GetBetween(start, end); err != nil {
		log.WithFields(log.Fields{"logger": "lcs.sync", "method": "loadInitial"}).
			WithError(err).Error("Could not load initial announcements. Exiting.")
		return
	}

	for _, a := range annons {
		an := a
		//if new add to queue
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "lcs.sync", "method": "loadInitial", "next": fmt.Sprintf("%+v", a)}).
				Debug("Publishing add to event bus.")
		}
		go s.b.Publish(event.QueueEnqueue, strconv.Itoa(a.StationID), &an)
	}
}

func (s *syncservice) ExtendQueueTail() {
	s.tailSync = make(chan bool)
	// first iteration is triggered instantly
	go func() {
		for {
			select {
			case <-s.tailSync:
				log.WithFields(log.Fields{"logger": "lcs.sync", "method": "ExtendQueueTail"}).Info("Received quit signal. Exiting.")
				return
			case <-time.After(s.tailSize):
				s.getTail()
			}
		}
	}()
}

func (s *syncservice) getTail() {
	newTail := s.currentTail.Add(s.tailSize)
	var annons []annon.Announcement
	var err error
	if annons, err = s.an.GetBetween(&s.currentTail, &newTail); err != nil {
		log.WithFields(log.Fields{"logger": "lcs.sync", "method": "getTail"}).
			WithError(err).Error("Could not read tail announcements. Exiting.")
		return
	}

	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "lcs.sync", "method": "getTail", "len": len(annons), "newTail": newTail}).
			Debug("Processing tail announcements.")
	}
	s.currentTail = newTail

	for _, a := range annons {
		an := a
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "lcs.sync", "method": "getTail", "annon": a.ID}).
				Debug("Publishing announcement.")
		}
		s.b.Publish(event.QueueUpdate, strconv.Itoa(a.StationID), &an, &an)
	}
}

func (s *syncservice) processChange(errorCount int) int {
	var next, old *annon.Announcement
	var err error
	if next, old, err = s.an.NextChange(); err != nil {
		log.WithFields(log.Fields{"logger": "sync.annon"}).
			WithError(err).Error("Error while obtaining new realtime changes.")
		// 10 consecutive errors will cause the loop to exit
		errorCount--
		if errorCount == 0 {
			log.WithFields(log.Fields{"logger": "sync.annon"}).
				Error("Limit number of consecutive errors reached. Exiting.")
			s.annonSync <- true
		}
		return errorCount
	}
	if log.GetLevel() >= log.InfoLevel {
		log.WithFields(log.Fields{"logger": "sync.annon"}).
			Info("New announcement change received.")
	}
	errorCount = 10
	var station string
	if next != nil {
		station = strconv.Itoa(next.StationID)
	} else {
		station = strconv.Itoa(old.StationID)
	}

	if next == nil {
		//if delete delete from queue
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "sync.annon", "next": next, "old": old}).
				Debug("Publishing delete to event bus.")
		}
		go s.b.Publish(event.QueueDeleteAnnon, station, old.ID)
		return errorCount
	}
	if old == nil {
		//if new add to queue
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "sync.annon", "next": next, "old": old}).
				Debug("Publishing add to event bus.")
		}
		go s.b.Publish(event.QueueEnqueue, station, next)
		return errorCount
	}
	//if update update in queue
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "sync.annon", "next": next, "old": old}).
			Debug("Publishing update to event bus.")
	}
	go s.b.Publish(event.QueueUpdate, station, next, old)
	return errorCount
}

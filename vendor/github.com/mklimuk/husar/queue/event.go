package queue

import (
	"fmt"
	"time"

	"github.com/mklimuk/husar/annon"

	log "github.com/Sirupsen/logrus"
)

//Event represents playback event on the timeline
type Event struct {
	ID            string         `json:"id"`
	StartTime     *time.Time     `json:"startTime"`
	EndTime       *time.Time     `json:"endTime"`
	PlaybackStart *time.Time     `json:"playbackStart"`
	PlaybackEnd   *time.Time     `json:"playbackEnd"`
	Duration      *time.Duration `json:"duration"`
	Priority      annon.Priority `json:"priority"`
	AnnonType     string         `json:"type"`
	TrainID       string         `json:"trainId"`
	AnnonID       string         `json:"annonId"`
	Lang          string         `json:"lang"`
	Text          string         `json:"text"`
	Mute          bool           `json:"mute"`
	Autoplay      bool           `json:"autoplay"`
}

//AdjustPlaybackTime adjusts event's playback time to align it with preceding event
func (e *Event) AdjustPlaybackTime(previous *Event, gap *time.Duration) {
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "playback.event", "method": "AdjustPlaybackTime", "event": e.ID}).
			Debug("Adjusting playback for event.")
	}
	if (*e.PlaybackStart).Before(*previous.PlaybackEnd) || (*e.PlaybackStart).Equal(*previous.PlaybackEnd) {
		s := (*previous.PlaybackEnd).Add(*gap)
		n := s.Add(*e.Duration)
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "playback.event", "method": "AdjustPlaybackTime", "event": e.ID, "start": s, "end": n}).
				Debug("Adjusted.")
		}
		e.PlaybackStart = &s
		e.PlaybackEnd = &n
	}
}

func compare(current *Event, candidate *Event) (insert bool, conflict bool) {
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "playback.queue", "method": "compare", "candidate": fmt.Sprintf("%+v", *candidate), "current": fmt.Sprintf("%+v", *current)}).
			Debug("Comparing elements.")
	}
	var over bool
	var difference int
	// check if they overlap
	if over, difference = overlap(current, candidate); !over {
		// if not check which one go first
		insert = difference > 0
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "playback.queue", "method": "compare", "candidate": candidate.ID, "current": current.ID, "insert": insert, "overlap": over}).
				Debug("No overlap.")
		}
		return insert, over
	}
	// if they do overlap compare their priorities if equal the one that starts first go first
	if current.Priority == candidate.Priority {
		// if they should start at the same time we play the shorter first
		if difference == 0 {
			insert = *candidate.Duration < *current.Duration
		} else {
			insert = difference > 0
		}
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "playback.queue", "method": "compare", "candidate": candidate.ID, "current": current.ID, "insert": insert, "overlap": over}).
				Debug("Overlap. Elements have equal priority.")
		}
		return insert, over
	}
	// higher priority we do not insert the new one
	if current.Priority < candidate.Priority {
		insert = false
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "playback.queue", "method": "compare", "candidate": candidate.ID, "current": current.ID, "insert": insert, "overlap": over}).
				Debug("Overlap. Current element has higher priority.")
		}
		return insert, over
	}
	diff := (*current.EndTime).Sub(*candidate.StartTime).Seconds()
	insert = diff > 15
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "playback.queue", "method": "compare", "candidate": candidate.ID, "current": current.ID, "insert": insert, "overlap": over, "diff": diff}).
			Debug("Overlap. Current element has lower priority.")
	}
	//if they would overlap by less than 15s we don't change the order
	return insert, over
}

func overlap(current *Event, candidate *Event) (overlap bool, difference int) {
	end := candidate.EndTime
	start := candidate.StartTime
	difference = int((*current.StartTime).Sub(*start).Seconds())
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "playback.queue", "method": "overlap", "candidate": candidate.ID, "current": current.ID, "difference": difference}).
			Debug("Difference between events' start.")
	}
	if difference == 0 {
		overlap = true
		return
	}

	if (*start).Before(*current.StartTime) && (*end).After(*current.StartTime) {
		overlap = true
		return
	}
	if (*start).Before(*current.EndTime) && (*end).After(*current.StartTime) {
		overlap = true
		return
	}
	overlap = false
	return
}

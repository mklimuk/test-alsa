package queue

import (
	"fmt"
	"sync"
	"time"

	"github.com/mklimuk/husar/annon"
	"github.com/mklimuk/husar/audio"

	log "github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
)

//Queue interface represents playback queue
type Queue interface {
	Add(a *annon.Announcement)
	Update(next *annon.Announcement, old *annon.Announcement)
	GetForAnnon(annonID string) (events []*Event)
	Remove(eventID string)
	RemoveForAnnon(annonID string) (removed []*Event)
	GetEvents() *[]*Event
	GetEvent(ID string) *Event
	Build(a *[]*annon.Announcement)
}

//New is a queue constructor
func New(ID string, gap *time.Duration, memory int, playHandler func(string, *Event)) Queue {
	q := queue{
		ID:     ID,
		gap:    gap,
		mute:   false,
		volume: 100,
		play: playback{
			handler: playHandler,
			memory:  memory,
		},
	}
	q.events = []*Event{}
	return Queue(&q)
}

type queue struct {
	ID     string
	events []*Event
	gap    *time.Duration
	lock   sync.Mutex
	play   playback
	volume int
	mute   bool
}

type playback struct {
	next    *Event
	timer   *chan bool
	handler func(string, *Event)
	memory  int
}

func (q *queue) GetEvents() *[]*Event {
	return &q.events
}

func (q *queue) GetEvent(ID string) *Event {
	for _, e := range q.events {
		if ID == e.ID {
			return e
		}
	}
	return nil
}

func (q *queue) Add(a *annon.Announcement) {
	events := q.buildEvents(a)
	for _, e := range events {
		q.addToQueue(e)
	}
}

func (q *queue) buildEvents(a *annon.Announcement) (events []*Event) {
	var d time.Duration
	if a.Audio == nil {
		log.WithFields(log.Fields{"logger": "playback.queue", "method": "buildEvents", "annon": a.ID}).
			Info("Audio info missing for announcement. Using default")
		d = time.Second * audio.DefaultDuration
	} else {
		d = time.Second * time.Duration(a.Audio.Duration)
	}
	for _, t := range a.Time {
		start := t
		end := t.Add(d)
		e := Event{
			ID:            uuid.NewV4().String(),
			AnnonID:       a.ID,
			TrainID:       a.TrainID,
			Lang:          a.Lang,
			AnnonType:     a.Type,
			Text:          a.Text.HTMLText,
			StartTime:     &start,
			PlaybackStart: &start,
			EndTime:       &end,
			PlaybackEnd:   &end,
			Duration:      &d,
			Autoplay:      a.Autoplay,
		}
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "playback.queue", "method": "buildEvents", "annon": a.ID, "event": fmt.Sprintf("%+v", e)}).
				Debug("Generated event for announcement.")
		}
		events = append(events, &e)
	}
	return events
}

func (q *queue) Update(next *annon.Announcement, old *annon.Announcement) {

	// get current events for this announcement
	removed := q.RemoveForAnnon(next.ID)
	events := q.buildEvents(next)

	for i, e := range events {
		// update properties
		if i < len(removed) && removed[i] != nil {
			e.Mute = removed[i].Mute
			e.Autoplay = removed[i].Autoplay
		}
		// add to queue
		q.addToQueue(e)
	}
}

func (q *queue) Remove(eventID string) {
	l := len(q.events)
	for i := 0; i <= l; i++ {
		if q.events[i].ID == eventID {
			copy(q.events[i:], q.events[i+1:])
			q.events[l-1] = nil
			q.events = q.events[:l-1]
			return
		}
	}
}

func (q *queue) RemoveForAnnon(annonID string) []*Event {
	// we use range to operate on a copy
	events := q.GetForAnnon(annonID)
	// copy IDs to avoid null pointers
	IDs := []string{}
	for _, e := range events {
		IDs = append(IDs, e.ID)
	}
	for _, id := range IDs {
		q.Remove(id)
	}
	return events
}

func (q *queue) GetForAnnon(annonID string) (events []*Event) {
	for _, e := range q.events {
		if e.AnnonID == annonID {
			events = append(events, e)
		}
	}
	return events
}

func (q *queue) Build(annons *[]*annon.Announcement) {
	for _, a := range *annons {
		q.Add(a)
	}
}

func (q *queue) addToQueue(element *Event) {
	// of course when there is nothing on the list it is fairly simple
	if len(q.events) == 0 {
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "playback.queue", "method": "addToQueue", "candidate": fmt.Sprintf("%+v", *element)}).
				Debug("Adding first element to the list")
		}
		q.events = append(q.events, element)
		return
	}
	var candidate, tmp *Event

	candidate = element
	l := len(q.events)
	for i := 0; i <= l; i++ {
		if i == len(q.events) {
			if log.GetLevel() >= log.DebugLevel {
				log.WithFields(log.Fields{"logger": "playback.queue", "method": "addToQueue", "candidate": candidate.ID}).
					Debug("Appending element to the end of the list.")
			}
			q.events = append(q.events, candidate)
			return
		}
		// checks if there is a playback time conflict and decides if we should insert the new element or not
		if insert, conflict := compare(q.events[i], candidate); insert {
			if !conflict {
				if log.GetLevel() >= log.DebugLevel {
					log.WithFields(log.Fields{"logger": "playback.queue", "method": "addToQueue", "candidate": candidate.ID, "current": q.events[i].ID}).
						Debug("No conflict. Appending element in front of the current one.")
				}
				q.events = append(q.events, nil)
				copy(q.events[i+1:], q.events[i:])
				q.events[i] = candidate
				return
			}
			// if there is a conflict it means we need to take current element out of the slice
			if log.GetLevel() >= log.DebugLevel {
				log.WithFields(log.Fields{"logger": "playback.queue", "method": "addToQueue", "candidate": candidate.ID, "current": q.events[i].ID}).
					Debug("Conflict detected. Replacing the current element with the candidate and continuing the loop.")
			}
			tmp = candidate
			candidate = q.events[i]
			q.events[i] = tmp
		}
	}
}

func (q *queue) updateNextPlayback(now time.Time) {
	var next *Event
	// we cleanup old events and then we take the first one that is in the future
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "playback.queue", "method": "updateNextPlayback"}).
			Debug("Cleaning up old announcements.")
	}
	q.cleanupOld(&now)
	next = getFirstFuture(q.events, &now)

	if next != nil && log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "playback.queue", "method": "updateNextPlayback", "next": fmt.Sprintf("%+v", *next)}).
			Debug("First playable event in the future.")
	}
	if next != q.play.next {
		if log.GetLevel() >= log.InfoLevel && next != nil {
			log.WithFields(log.Fields{"logger": "playback.queue", "method": "updateNextPlayback", "next": fmt.Sprintf("%+v", *next)}).
				Info("Changing next playback event.")
		}
		var err error
		if err = q.setNextPlayback(next, now); err != nil {
			log.WithFields(log.Fields{"logger": "playback.queue", "method": "updateNextPlayback", "next": fmt.Sprintf("%+v", *next)}).
				WithError(err).Error("Could not update playback item.")
		}
	}
}

func (q *queue) cleanupOld(now *time.Time) {
	if len(q.events) <= q.play.memory {
		return
	}
	for inThePast(q.events[q.play.memory].PlaybackEnd, now) {
		q.events[0] = nil
		q.events = q.events[1:]
		if len(q.events) <= q.play.memory {
			return
		}
	}

}

func getFirstFuture(events []*Event, now *time.Time) *Event {
	for _, eve := range events {
		if !inThePast(eve.PlaybackStart, now) && eve.Autoplay {
			return eve
		}
	}
	return nil
}

func (q *queue) setNextPlayback(e *Event, now time.Time) error {
	if e == nil {
		return nil
	}
	if q.play.timer != nil {
		// stopping current timer
		close(*q.play.timer)
	}
	q.play.next = e
	t := make(chan bool)
	q.play.timer = &t
	d := e.PlaybackStart.Sub(now)
	if d < 0 {
		log.WithFields(log.Fields{"logger": "lcs.queue", "method": "setNextPlayback", "duration": d}).
			Error("Trying to schedule playback in the past. Returning error.")
		return fmt.Errorf("Trying to schedule event in the past.")
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "lcs.queue", "method": "setNextPlayback", "duration": d, "startTime": e.PlaybackStart, "event": e.ID}).
			Debug("Playback timer set.")
	}
	// launch the timer
	go func() {
		select {
		case <-time.After(d):
			log.WithFields(log.Fields{"logger": "lcs.queue", "method": "setNextPlayback", "event": q.play.next}).
				Info("Playback timer triggered.")
			q.play.handler(q.ID, q.play.next)
		case <-(*q.play.timer):
			log.WithFields(log.Fields{"logger": "lcs.queue", "method": "setNextPlayback", "event": q.play.next}).
				Info("Playback timer cancelled.")
		}
	}()
	return nil
}

func (q *queue) adjustPlayback() {
	// We established the order. Now it's time to adjust time.
	for i := 0; i < len(q.events); i++ {
		if i > 0 {
			q.events[i].AdjustPlaybackTime(q.events[i-1], q.gap)
		}
	}
}

func inThePast(t *time.Time, now *time.Time) bool {
	return now.After(*t) || now.Equal(*t)
}

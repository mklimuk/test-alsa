package queue

import (
	"time"

	"github.com/mklimuk/husar/event"

	log "github.com/Sirupsen/logrus"
)

const (
	mute   bool = true
	unmute bool = false
)

func (r *manager) ToggleMute(req string) {
	msg := parseQueueMessage(req)
	q, res := r.toggleMute(msg.QueueID, msg.EventID)
	//TODO refactor
	if res {
		r.publishContent(event.QueueChange, q)
	}
}

func (r *manager) toggleMute(queueID string, eventID string) (*queue, bool) {
	var q *queue
	if q = r.queues[queueID]; q == nil {
		log.WithFields(log.Fields{"logger": "lcs.queue", "method": "Enqueue", "id": queueID}).
			Error("Queue not found.")
		return nil, false
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	for _, e := range q.events {
		if eventID == e.ID {
			e.Mute = !e.Mute
			return q, true
		}
	}
	return q, false
}

func (r *manager) MuteAll(req string) {
	r.setMute(req, mute)
}

func (r *manager) UnmuteAll(req string) {
	r.setMute(req, unmute)
}

func (r *manager) setMute(req string, state bool) {
	msg := parseQueueMessage(req)
	var q *queue
	if q = r.queues[msg.QueueID]; q == nil {
		log.WithFields(log.Fields{"logger": "lcs.queue", "method": "Enqueue", "id": msg.QueueID}).
			Error("Queue not found.")
		return
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	q.mute = state
	// TODO refactor into some kind of update event
	r.publishContent(event.QueueChange, q)
}

func (r *manager) MuteFirstIfRequired(ID string, end time.Time) error {
	q := r.GetQueue(ID).(*queue)
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.play.next != nil && q.play.next.PlaybackStart.Before(end) {
		q.play.next.Mute = true
	}
	r.publishContent(event.QueueChange, q)
	return nil
}

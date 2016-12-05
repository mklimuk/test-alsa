package queue

import (
	"fmt"
	"time"

	"github.com/mklimuk/husar/event"

	log "github.com/Sirupsen/logrus"
)

func (r *manager) Playback(queueID string, ev *Event) {
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "lcs.queue", "method": "Playback", "queueID": queueID, "event": fmt.Sprintf("%+v", ev)}).
			Debug("Triggering playback.")
	}
	var q *queue
	var ok bool
	if q, ok = r.queues[queueID]; !ok {
		log.WithFields(log.Fields{"logger": "lcs.queue", "method": "Playback", "queueID": queueID}).
			Error("Queue not found.")
		return
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	// if mute was not set we play the audio
	if ev.Autoplay && !ev.Mute && !q.mute {
		r.b.Publish(event.AudioPlay, queueID, ev)
	}
	ev.Autoplay = false
	log.WithFields(log.Fields{"logger": "lcs.queue", "method": "Playback", "queueID": queueID}).
		Info("Updating next playback event.")
	q.updateNextPlayback(time.Now())
	r.publishContent(event.QueueChange, q)
}

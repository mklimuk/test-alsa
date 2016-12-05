package queue

import (
	"encoding/json"
	"time"

	"github.com/mklimuk/husar/annon"
	"github.com/mklimuk/husar/event"

	log "github.com/Sirupsen/logrus"
)

//Manager manages a map of playback queues
type Manager interface {
	CreateQueue(ID string) *Queue
	RemoveQueue(ID string)
	GetQueue(ID string) Queue
	GetEvent(ID string, eventID string) (*Event, error)
	MuteFirstIfRequired(ID string, end time.Time) error
	Enqueue(ID string, toAdd *annon.Announcement)
	Update(ID string, next *annon.Announcement, old *annon.Announcement)
	DeleteAnnon(ID string, annonID string)
	SetVolume(ID string, volume int)
}

//NewManager is a queue manager constructor
func NewManager(gap time.Duration, queueMemory int, b *event.Bus) Manager {
	r := manager{gap: gap, b: b, queueMemory: queueMemory}
	r.queues = make(map[string]*queue)
	b.Subscribe(event.QueueEnqueue, r.Enqueue)
	b.Subscribe(event.QueueUpdate, r.Update)
	b.Subscribe(event.QueueDeleteAnnon, r.DeleteAnnon)
	b.Subscribe(event.GetQueue, r.PublishQueueContent)
	b.Subscribe(event.PlaybackToggleMute, r.ToggleMute)
	b.Subscribe(event.PlaybackMuteAll, r.MuteAll)
	b.Subscribe(event.PlaybackUnmuteAll, r.UnmuteAll)
	return Manager(&r)
}

type manager struct {
	queues      map[string]*queue
	gap         time.Duration
	queueMemory int
	b           *event.Bus
}

func (r *manager) CreateQueue(ID string) *Queue {
	log.WithFields(log.Fields{"logger": "lcs.queue", "method": "CreateQueue", "id": ID}).
		Info("Creating queue.")
	q := New(ID, &r.gap, r.queueMemory, r.Playback)
	r.queues[ID] = q.(*queue)
	return &q
}

func (r *manager) RemoveQueue(ID string) {
	log.WithFields(log.Fields{"logger": "lcs.queue", "method": "RemoveQueue", "id": ID}).
		Info("Removing queue.")
	delete(r.queues, ID)
}

func (r *manager) GetQueue(ID string) Queue {
	if log.GetLevel() >= log.InfoLevel {
		log.WithFields(log.Fields{"logger": "lcs.queue", "method": "GetQueue", "id": ID}).
			Info("Returning queue.")
	}
	return Queue(r.queues[ID])
}

func (r *manager) GetEvent(ID string, eventID string) (*Event, error) {
	q := r.GetQueue(ID)
	return q.GetEvent(eventID), nil
}

func (r *manager) Enqueue(ID string, toAdd *annon.Announcement) {
	var q *queue
	if q = r.queues[ID]; q == nil {
		log.WithFields(log.Fields{"logger": "lcs.queue", "method": "Enqueue", "id": ID}).
			Error("Queue not found.")
		return
	}
	log.WithFields(log.Fields{"logger": "lcs.queue", "method": "Enqueue", "id": ID, "annon": toAdd.ID}).
		Info("Adding announcement to the queue.")
	(*q).lock.Lock()
	defer (*q).lock.Unlock()
	(*q).Add(toAdd)
	(*q).adjustPlayback()
	(*q).updateNextPlayback(time.Now())
	r.publishContent(event.QueueChange, q)
}

func (r *manager) Update(ID string, next *annon.Announcement, old *annon.Announcement) {
	var q *queue
	var ok bool
	if q, ok = r.queues[ID]; !ok {
		log.WithFields(log.Fields{"logger": "lcs.queue", "method": "Update", "id": ID}).
			Error("Queue not found.")
		return
	}
	log.WithFields(log.Fields{"logger": "lcs.queue", "method": "Update", "id": ID, "annon": next.ID}).
		Info("Updating announcement in the queue.")
	(*q).lock.Lock()
	defer (*q).lock.Unlock()
	(*q).Update(next, old)
	(*q).adjustPlayback()
	(*q).updateNextPlayback(time.Now())
	r.publishContent(event.QueueChange, q)
}

func (r *manager) DeleteAnnon(ID string, annonID string) {
	var q *queue
	var ok bool
	if q, ok = r.queues[ID]; !ok {
		log.WithFields(log.Fields{"logger": "lcs.queue", "method": "Enqueue", "id": ID}).
			Error("Queue not found.")
		return
	}
	log.WithFields(log.Fields{"logger": "lcs.queue", "method": "DeleteAnnon", "id": ID, "annon": annonID}).
		Info("Deleting announcement from the queue.")
	(*q).lock.Lock()
	defer (*q).lock.Unlock()
	(*q).RemoveForAnnon(annonID)
	(*q).adjustPlayback()
	(*q).updateNextPlayback(time.Now())
	r.publishContent(event.QueueChange, q)
}

func (r *manager) publishContent(eventType event.Type, q *queue) {
	var err error
	var msg []byte
	if msg, err = json.Marshal(queueContent{QueueID: q.ID, Events: (*q).GetEvents(), NextEvent: q.play.next, Mute: q.mute}); err != nil {
		log.WithFields(log.Fields{"logger": "lcs.queue", "method": "publishContent", "id": q.ID}).
			WithError(err).Error("Could not marshall queue contents.")
	}
	if log.GetLevel() >= log.InfoLevel {
		log.WithFields(log.Fields{"logger": "lcs.queue", "method": "publishContent", "id": q.ID}).
			Info("Publishing queue content.")
	}
	r.b.Publish(eventType, eventType, string(msg))
}

func (r *manager) PublishQueueContent(ID string) {
	var q *queue
	var ok bool
	if q, ok = r.queues[ID]; !ok {
		log.WithFields(log.Fields{"logger": "lcs.queue", "method": "Update", "id": ID}).
			Error("Queue not found.")
		return
	}
	r.publishContent(event.QueueContent, q)
}

func (r *manager) SetVolume(ID string, volume int) {
	var q *queue
	var ok bool
	if q, ok = r.queues[ID]; !ok {
		log.WithFields(log.Fields{"logger": "lcs.queue", "method": "SetVolume", "id": ID}).
			Error("Queue not found.")
		return
	}
	if log.GetLevel() >= log.InfoLevel {
		log.WithFields(log.Fields{"logger": "lcs.queue", "method": "SetVolume", "id": ID}).
			Info("Setting volume on the queue.")
	}
	q.volume = volume
}

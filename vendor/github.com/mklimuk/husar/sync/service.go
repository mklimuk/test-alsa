package sync

import (
	"time"

	"github.com/mklimuk/husar/annon"
	"github.com/mklimuk/husar/event"
	"github.com/mklimuk/husar/service"
	"github.com/mklimuk/husar/train"
)

//Service exposes synchronization mechanism
type Service interface {
	EnableAnnouncementsSync()
	ExtendQueueTail()
	EnableTimetableSync()
}

//NewService is a Service listener
func NewService(t service.Timetable, a service.Announcement, tr train.Store, an annon.Store, b *event.Bus, windowSize time.Duration, windowAhead time.Duration, tailSize time.Duration) Service {
	syn := syncservice{
		b:           b,
		an:          an,
		tr:          tr,
		a:           a,
		t:           t,
		windowSize:  windowSize,
		windowAhead: windowAhead,
		tailSize:    tailSize,
		syncEnabled: false}
	return Service(&syn)
}

type syncservice struct {
	scheduledSync chan bool
	liveSync      chan bool
	settingsSync  chan bool
	annonSync     chan bool
	tailSync      chan bool
	syncEnabled   bool
	t             service.Timetable
	a             service.Announcement
	tr            train.Store
	an            annon.Store
	windowSize    time.Duration
	windowAhead   time.Duration
	tailSize      time.Duration
	b             *event.Bus
	currentTail   time.Time
}

// enables all synchronization features
func (s *syncservice) EnableTimetableSync() {
	if !s.syncEnabled {
		s.timetable()
		s.realtime()
		s.settings()
		s.syncEnabled = true
	}
}

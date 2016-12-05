package queue

import (
	"testing"
	"time"

	"github.com/mklimuk/husar/annon"
	"github.com/mklimuk/husar/config"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type QueueTestSuite struct {
	suite.Suite
	g time.Duration
}

func (suite *QueueTestSuite) SetupSuite() {
	suite.g = time.Duration(10) * time.Second
	log.SetLevel(log.DebugLevel)
}

func (suite *QueueTestSuite) TestAddToEmpty() {
	q := queue{gap: &suite.g}
	q.play.memory = 2
	start := time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone)
	end := time.Date(2016, time.April, 10, 15, 31, 0, 0, config.Timezone)
	d := time.Duration(25) * time.Second
	e := Event{ID: "abcd", StartTime: &start, EndTime: &end, Duration: &d, AnnonID: "test", Priority: annon.P1}
	q.addToQueue(&e)
	assert.Len(suite.T(), q.events, 1)
}

func (suite *QueueTestSuite) TestOverlap() {
	o, d := overlap(buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 22, 0, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"),
		buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 22, 0, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"))
	assert.True(suite.T(), o)
	assert.Equal(suite.T(), 0, d)
	o, _ = overlap(buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 22, 0, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"),
		buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 30, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 22, 0, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"))
	assert.True(suite.T(), o)
	o, _ = overlap(buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 22, 0, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"),
		buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 30, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"))
	assert.True(suite.T(), o)
	o, _ = overlap(buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 22, 0, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"),
		buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 10, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 30, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"))
	assert.True(suite.T(), o)
	o, _ = overlap(buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 20, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 30, 0, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"),
		buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 22, 0, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"))
	assert.True(suite.T(), o)
}

func (suite *QueueTestSuite) TestAddNoConflict() {
	q := queue{gap: &suite.g}
	q.play.memory = 2
	q.addToQueue(buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 22, 0, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"))
	q.addToQueue(buildEvent("b", time.Date(2016, time.April, 10, 15, 10, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 10, 30, 0, config.Timezone), time.Duration(30)*time.Second, annon.P2, "test"))
	q.addToQueue(buildEvent("c", time.Date(2016, time.April, 10, 15, 45, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 45, 15, 0, config.Timezone), time.Duration(15)*time.Second, annon.P1, "test"))
	assert.Len(suite.T(), q.events, 3)
	assert.Equal(suite.T(), "b", q.events[0].ID)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 10, 0, 0, config.Timezone), *q.events[0].PlaybackStart)
	assert.Equal(suite.T(), "a", q.events[1].ID)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), *q.events[1].PlaybackStart)
	assert.Equal(suite.T(), "c", q.events[2].ID)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 45, 0, 0, config.Timezone), *q.events[2].PlaybackStart)

}

func (suite *QueueTestSuite) TestAddConflictEqualPriority() {
	q := queue{gap: &suite.g}
	q.play.memory = 2
	q.addToQueue(buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 22, 0, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"))
	q.addToQueue(buildEvent("b", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 30, 0, config.Timezone), time.Duration(30)*time.Second, annon.P1, "test"))
	q.addToQueue(buildEvent("c", time.Date(2016, time.April, 10, 15, 21, 15, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 30, 0, config.Timezone), time.Duration(15)*time.Second, annon.P1, "test"))
	q.adjustPlayback()
	assert.Len(suite.T(), q.events, 3)
	assert.Equal(suite.T(), "b", q.events[0].ID)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), *q.events[0].PlaybackStart)
	assert.Equal(suite.T(), "a", q.events[1].ID)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 21, 40, 0, config.Timezone), *q.events[1].PlaybackStart)
	assert.Equal(suite.T(), "c", q.events[2].ID)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 22, 50, 0, config.Timezone), *q.events[2].PlaybackStart)
}

func (suite *QueueTestSuite) TestAddConflictDifferentPriority() {
	q := queue{gap: &suite.g}
	q.play.memory = 2
	q.addToQueue(buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 22, 0, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"))
	q.addToQueue(buildEvent("b", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 30, 0, config.Timezone), time.Duration(30)*time.Second, annon.P2, "test"))
	q.addToQueue(buildEvent("c", time.Date(2016, time.April, 10, 15, 21, 15, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 30, 0, config.Timezone), time.Duration(15)*time.Second, annon.P1, "test"))
	q.adjustPlayback()
	assert.Len(suite.T(), q.events, 3)
	assert.Equal(suite.T(), "a", q.events[0].ID)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), *q.events[0].PlaybackStart)
	assert.Equal(suite.T(), "b", q.events[1].ID)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 22, 10, 0, config.Timezone), *q.events[1].PlaybackStart)
	// because playback difference for c if it stays after b is less than 15s
	assert.Equal(suite.T(), "c", q.events[2].ID)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 22, 50, 0, config.Timezone), *q.events[2].PlaybackStart)
	// different order same result
	q = queue{gap: &suite.g}
	q.play.memory = 2
	q.addToQueue(buildEvent("c", time.Date(2016, time.April, 10, 15, 21, 15, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 30, 0, config.Timezone), time.Duration(15)*time.Second, annon.P1, "test"))
	q.addToQueue(buildEvent("b", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 30, 0, config.Timezone), time.Duration(30)*time.Second, annon.P2, "test"))
	q.addToQueue(buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 22, 0, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"))
	q.adjustPlayback()
	assert.Len(suite.T(), q.events, 3)
	assert.Equal(suite.T(), "a", q.events[0].ID)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), *q.events[0].PlaybackStart)
	assert.Equal(suite.T(), "b", q.events[1].ID)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 22, 10, 0, config.Timezone), *q.events[1].PlaybackStart)
	// because playback difference for c if it stays after b is less than 15s
	assert.Equal(suite.T(), "c", q.events[2].ID)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 22, 50, 0, config.Timezone), *q.events[2].PlaybackStart)
}

func (suite *QueueTestSuite) TestNextPlayback() {
	q := queue{gap: &suite.g}
	q.play.memory = 2
	q.events = []*Event{
		buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 1, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"),
		buildEvent("c", time.Date(2016, time.April, 10, 15, 21, 15, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 16, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"),
		buildEvent("b", time.Date(2016, time.April, 10, 15, 21, 25, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 26, 0, config.Timezone), time.Duration(60)*time.Second, annon.P2, "test"),
	}
	q.updateNextPlayback(time.Now())
	assert.Len(suite.T(), q.events, 2)
	assert.Nil(suite.T(), q.play.next)
	assert.Nil(suite.T(), q.play.timer)
	q.events = []*Event{
		buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 1, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"),
		buildEvent("c", time.Date(2016, time.April, 10, 15, 21, 15, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 16, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test"),
		buildEvent("b", time.Date(2016, time.April, 10, 15, 21, 25, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 26, 0, config.Timezone), time.Duration(60)*time.Second, annon.P2, "test"),
	}
	q.updateNextPlayback(time.Date(2016, time.April, 10, 15, 0, 0, 0, config.Timezone))
	assert.Len(suite.T(), q.events, 3)
	assert.Equal(suite.T(), "a", q.play.next.ID)
	assert.NotNil(suite.T(), q.play.timer)
	q.events = []*Event{
		buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 1, 0, config.Timezone), time.Duration(1)*time.Second, annon.P1, "test"),
		buildEvent("c", time.Date(2016, time.April, 10, 15, 21, 15, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 16, 0, config.Timezone), time.Duration(1)*time.Second, annon.P1, "test"),
		buildEvent("b", time.Date(2016, time.April, 10, 15, 21, 25, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 26, 0, config.Timezone), time.Duration(1)*time.Second, annon.P2, "test"),
		buildEvent("d", time.Date(2016, time.April, 10, 15, 24, 25, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 24, 26, 0, config.Timezone), time.Duration(1)*time.Second, annon.P2, "test"),
		buildEvent("e", time.Date(2016, time.April, 10, 15, 26, 25, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 26, 26, 0, config.Timezone), time.Duration(1)*time.Second, annon.P2, "test"),
		buildEvent("f", time.Date(2016, time.April, 10, 15, 28, 25, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 28, 26, 0, config.Timezone), time.Duration(1)*time.Second, annon.P2, "test"),
		buildEvent("g", time.Date(2016, time.April, 10, 15, 30, 25, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 30, 26, 0, config.Timezone), time.Duration(1)*time.Second, annon.P2, "test"),
		buildEvent("h", time.Date(2016, time.April, 10, 15, 34, 25, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 34, 26, 0, config.Timezone), time.Duration(1)*time.Second, annon.P2, "test"),
	}
	q.updateNextPlayback(time.Date(2016, time.April, 10, 15, 22, 0, 0, config.Timezone))
	assert.Len(suite.T(), q.events, 7)
	assert.Equal(suite.T(), "d", q.play.next.ID)
	assert.NotNil(suite.T(), q.play.timer)

}

func (suite *QueueTestSuite) TestTriggerPlayback() {
	q := queue{gap: &suite.g, ID: "test"}
	q.play.memory = 2
	triggered := false
	q.play.handler = func(id string, e *Event) {
		triggered = true
		assert.Equal(suite.T(), "test", id)
		assert.Equal(suite.T(), "d", e.ID)
	}
	q.events = []*Event{
		buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 1, 0, config.Timezone), time.Duration(1)*time.Second, annon.P1, "test"),
		buildEvent("c", time.Date(2016, time.April, 10, 15, 21, 15, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 16, 0, config.Timezone), time.Duration(1)*time.Second, annon.P1, "test"),
		buildEvent("b", time.Date(2016, time.April, 10, 15, 21, 25, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 26, 0, config.Timezone), time.Duration(1)*time.Second, annon.P2, "test"),
		buildEvent("d", time.Date(2016, time.April, 10, 15, 24, 25, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 24, 26, 0, config.Timezone), time.Duration(1)*time.Second, annon.P2, "test"),
		buildEvent("e", time.Date(2016, time.April, 10, 15, 26, 25, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 26, 26, 0, config.Timezone), time.Duration(1)*time.Second, annon.P2, "test"),
	}
	q.updateNextPlayback(time.Date(2016, time.April, 10, 15, 24, 23, 0, config.Timezone))
	assert.Len(suite.T(), q.events, 4)
	assert.Equal(suite.T(), "d", q.play.next.ID)
	assert.NotNil(suite.T(), q.play.timer)
	time.Sleep(time.Duration(4) * time.Second)
	assert.True(suite.T(), triggered)
}

func (suite *QueueTestSuite) TestReplacePlayback() {
	q := queue{gap: &suite.g, ID: "test"}
	q.play.memory = 2
	triggered := false
	q.play.handler = func(id string, e *Event) {
		triggered = true
		assert.Equal(suite.T(), "test", id)
		assert.Equal(suite.T(), "b", e.ID)
	}
	q.events = []*Event{
		buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 1, 0, config.Timezone), time.Duration(1)*time.Second, annon.P1, "test"),
		buildEvent("c", time.Date(2016, time.April, 10, 15, 21, 15, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 16, 0, config.Timezone), time.Duration(1)*time.Second, annon.P1, "test"),
	}

	q.updateNextPlayback(time.Date(2016, time.April, 10, 15, 15, 0, 0, config.Timezone))
	assert.Len(suite.T(), q.events, 2)
	assert.Equal(suite.T(), "a", q.play.next.ID)
	assert.NotNil(suite.T(), q.play.timer)

	q.addToQueue(buildEvent("b", time.Date(2016, time.April, 10, 15, 15, 20, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 15, 25, 0, config.Timezone), time.Duration(5)*time.Second, annon.P1, "test"))
	q.adjustPlayback()
	q.updateNextPlayback(time.Date(2016, time.April, 10, 15, 15, 17, 0, config.Timezone))
	assert.Len(suite.T(), q.events, 3)
	assert.Equal(suite.T(), "b", q.play.next.ID)
	assert.NotNil(suite.T(), q.play.timer)
	time.Sleep(time.Duration(4) * time.Second)
	assert.True(suite.T(), triggered)

}

func (suite *QueueTestSuite) TestTriggerTwo() {
	q := queue{gap: &suite.g, ID: "test"}
	q.play.memory = 2
	triggered := 0
	q.play.handler = func(id string, e *Event) {
		triggered++
		q.updateNextPlayback(time.Date(2016, time.April, 10, 15, 21, 2, 0, config.Timezone))
	}
	q.events = []*Event{
		buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 1, 0, config.Timezone), time.Duration(1)*time.Second, annon.P1, "test"),
		buildEvent("c", time.Date(2016, time.April, 10, 15, 21, 3, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 4, 0, config.Timezone), time.Duration(1)*time.Second, annon.P1, "test"),
		buildEvent("b", time.Date(2016, time.April, 10, 15, 21, 25, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 26, 0, config.Timezone), time.Duration(1)*time.Second, annon.P2, "test"),
	}

	q.updateNextPlayback(time.Date(2016, time.April, 10, 15, 20, 59, 0, config.Timezone))
	assert.Len(suite.T(), q.events, 3)
	assert.Equal(suite.T(), "a", q.play.next.ID)
	assert.NotNil(suite.T(), q.play.timer)
	time.Sleep(time.Duration(5) * time.Second)
	assert.Equal(suite.T(), 2, triggered)

}

func (suite *QueueTestSuite) TestBreak() {
	q := queue{gap: &suite.g, ID: "test"}
	q.play.memory = 2
	triggered := 0
	q.play.handler = func(id string, e *Event) {
		triggered++
	}
	q.events = []*Event{
		buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 1, 0, config.Timezone), time.Duration(1)*time.Second, annon.P1, "test"),
		buildEvent("c", time.Date(2016, time.April, 10, 15, 21, 3, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 4, 0, config.Timezone), time.Duration(1)*time.Second, annon.P1, "test"),
		buildEvent("b", time.Date(2016, time.April, 10, 15, 21, 25, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 26, 0, config.Timezone), time.Duration(1)*time.Second, annon.P2, "test"),
	}

	q.updateNextPlayback(time.Date(2016, time.April, 10, 15, 20, 59, 0, config.Timezone))
	assert.Len(suite.T(), q.events, 3)
	assert.Equal(suite.T(), "a", q.play.next.ID)
	assert.NotNil(suite.T(), q.play.timer)
	close(*q.play.timer)
	time.Sleep(time.Duration(2) * time.Second)
	assert.Equal(suite.T(), 0, triggered)

}

func (suite *QueueTestSuite) TestAddForAnnon() {
	a := annon.Announcement{
		Time: []time.Time{
			time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone),
			time.Date(2016, time.April, 10, 15, 23, 0, 0, config.Timezone),
		},
		Text: &annon.Text{
			HTMLText: "abcd",
		},
		Audio: &annon.Audio{
			Duration: 10,
		},
	}
	q := queue{gap: &suite.g}
	q.play.memory = 2
	q.Add(&a)
	q.adjustPlayback()
	assert.Len(suite.T(), q.events, 2)
	assert.Equal(suite.T(), a.Time[0], *q.events[0].StartTime)
	assert.Equal(suite.T(), a.Time[0], *q.events[0].PlaybackStart)
	assert.Equal(suite.T(), a.Time[1], *q.events[1].StartTime)
	assert.Equal(suite.T(), a.Time[1], *q.events[1].PlaybackStart)
}

func (suite *QueueTestSuite) TestRemoveForAnnon() {
	q := queue{gap: &suite.g}
	q.play.memory = 2
	q.addToQueue(buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 22, 0, 0, config.Timezone), time.Duration(60)*time.Second, annon.P1, "test2"))
	q.addToQueue(buildEvent("b", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 30, 0, config.Timezone), time.Duration(30)*time.Second, annon.P1, "test"))
	q.addToQueue(buildEvent("c", time.Date(2016, time.April, 10, 15, 21, 15, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 30, 0, config.Timezone), time.Duration(15)*time.Second, annon.P1, "test2"))
	q.addToQueue(buildEvent("d", time.Date(2016, time.April, 10, 15, 22, 15, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 22, 30, 0, config.Timezone), time.Duration(15)*time.Second, annon.P1, "test"))
	q.addToQueue(buildEvent("e", time.Date(2016, time.April, 10, 15, 23, 15, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 23, 30, 0, config.Timezone), time.Duration(15)*time.Second, annon.P1, "test2"))
	q.addToQueue(buildEvent("f", time.Date(2016, time.April, 10, 15, 23, 45, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 24, 30, 0, config.Timezone), time.Duration(45)*time.Second, annon.P1, "test"))
	res := q.RemoveForAnnon("test2")
	assert.Len(suite.T(), q.events, 3)
	assert.Len(suite.T(), res, 3)
	assert.NotNil(suite.T(), res[0])
}

func (suite *QueueTestSuite) TestFirstFuture() {
	q := queue{gap: &suite.g}
	q.play.memory = 2
	q.addToQueue(buildEvent("a", time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 21, 45, 0, config.Timezone), time.Duration(45)*time.Second, annon.P1, "test2"))
	q.addToQueue(buildEvent("b", time.Date(2016, time.April, 10, 15, 22, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 22, 30, 0, config.Timezone), time.Duration(30)*time.Second, annon.P1, "test"))
	q.addToQueue(buildEvent("c", time.Date(2016, time.April, 10, 15, 22, 25, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 22, 40, 0, config.Timezone), time.Duration(15)*time.Second, annon.P1, "test2"))
	q.addToQueue(buildEvent("d", time.Date(2016, time.April, 10, 15, 22, 45, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 23, 0, 0, config.Timezone), time.Duration(15)*time.Second, annon.P1, "test"))
	q.addToQueue(buildEvent("e", time.Date(2016, time.April, 10, 15, 23, 15, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 23, 30, 0, config.Timezone), time.Duration(15)*time.Second, annon.P1, "test2"))
	q.addToQueue(buildEvent("f", time.Date(2016, time.April, 10, 15, 23, 45, 0, config.Timezone), time.Date(2016, time.April, 10, 15, 24, 30, 0, config.Timezone), time.Duration(45)*time.Second, annon.P1, "test"))
	now := time.Date(2016, time.April, 10, 15, 22, 20, 0, config.Timezone)
	e := getFirstFuture(q.events, &now)
	assert.Equal(suite.T(), "c", e.ID)
	now = time.Date(2016, time.April, 10, 15, 22, 44, 0, config.Timezone)
	e = getFirstFuture(q.events, &now)
	assert.Equal(suite.T(), "d", e.ID)
	now = time.Date(2016, time.April, 10, 15, 20, 20, 0, config.Timezone)
	e = getFirstFuture(q.events, &now)
	assert.Equal(suite.T(), "a", e.ID)
	now = time.Date(2016, time.April, 10, 15, 23, 15, 0, config.Timezone)
	e = getFirstFuture(q.events, &now)
	assert.Equal(suite.T(), "f", e.ID)
}

func buildEvent(ID string, s time.Time, e time.Time, d time.Duration, p annon.Priority, annonID string) *Event {
	ev := Event{ID: ID, StartTime: &s, EndTime: &e, PlaybackStart: &s, PlaybackEnd: &e, Duration: &d, AnnonID: annonID, Priority: p, Autoplay: true}
	return &ev
}

func TestQueueTestSuite(t *testing.T) {
	suite.Run(t, new(QueueTestSuite))
}

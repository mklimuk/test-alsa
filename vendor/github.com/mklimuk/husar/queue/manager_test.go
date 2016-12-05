package queue

import (
	"testing"
	"time"

	"github.com/mklimuk/husar/annon"
	"github.com/mklimuk/husar/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ManagerTestSuite struct {
	suite.Suite
	g time.Duration
}

func (suite *ManagerTestSuite) SetupSuite() {
	suite.g = time.Duration(10) * time.Second
}

func (suite *ManagerTestSuite) TestAddToEmpty() {
	q := queue{gap: &suite.g}
	q.play.memory = 2
	start := time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone)
	end := time.Date(2016, time.April, 10, 15, 31, 0, 0, config.Timezone)
	d := time.Duration(25) * time.Second
	e := Event{ID: "abcd", StartTime: &start, EndTime: &end, Duration: &d, AnnonID: "test", Priority: annon.P1}
	q.addToQueue(&e)
	assert.Len(suite.T(), q.events, 1)
}

func (suite *ManagerTestSuite) TestTriggerPlayback() {
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

func TestManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ManagerTestSuite))
}

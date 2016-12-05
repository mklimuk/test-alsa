package queue

import (
	"testing"
	"time"

	"github.com/mklimuk/husar/config"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAdjust(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	start := time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone)
	end := time.Date(2016, time.April, 10, 15, 31, 0, 0, config.Timezone)
	d := time.Duration(30) * time.Second
	g := time.Duration(5) * time.Second
	e := Event{StartTime: &start, EndTime: &start, Duration: &d, PlaybackStart: &start, PlaybackEnd: &start}
	e2 := Event{PlaybackEnd: &end}
	e.AdjustPlaybackTime(&e2, &g)
	assert.Equal(t, time.Date(2016, time.April, 10, 15, 31, 5, 0, config.Timezone), *e.PlaybackStart)
	assert.Equal(t, time.Date(2016, time.April, 10, 15, 31, 35, 0, config.Timezone), *e.PlaybackEnd)
	end = time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone)
	e2.PlaybackEnd = &end
	e.AdjustPlaybackTime(&e2, &g)
	assert.Equal(t, time.Date(2016, time.April, 10, 15, 31, 5, 0, config.Timezone), *e.PlaybackStart)
	assert.Equal(t, time.Date(2016, time.April, 10, 15, 31, 35, 0, config.Timezone), *e.PlaybackEnd)
}

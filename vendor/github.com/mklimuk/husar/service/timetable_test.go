package service

import (
	"testing"
	"time"

	"github.com/mklimuk/husar/config"
	"github.com/mklimuk/husar/msg"
	"github.com/mklimuk/husar/test"
	"github.com/mklimuk/husar/train"

	"github.com/mklimuk/husar/event"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const trainID string = "abcdef12345"

type TimetableTestSuite struct {
	suite.Suite
	b *event.Bus
}

func (suite *TimetableTestSuite) SetupSuite() {
	suite.b = event.New()
	log.SetLevel(log.DebugLevel)
}

func (suite *TimetableTestSuite) TestGetTimetable() {
	st := test.TrainStoreMock{}
	st.On("GetBetweenWithRealtime", mock.Anything, mock.Anything).Return([]train.Train{
		train.Train{ID: "test"},
	}, []string{"test"}, nil).On("GetAllSettings", []string{"test"}).Return(&[]train.Settings{
		train.Settings{ID: "test"},
	}, nil)
	NewTimetable(suite.b, &st)
	passed := false
	suite.b.Subscribe(event.TimetableContent, func(title event.Type, body string) {
		log.Info("Received notification")
		passed = true
		assert.Equal(suite.T(), event.TimetableContent, title)
		assert.NotEmpty(suite.T(), body)
	})
	suite.b.Publish(event.GetTimetable, `{"stationId":48355,"from":"2016-05-29T15:03:22+02:00","to":"2016-05-29T17:03:22+02:00"}`)
	suite.b.Publish(event.GetTimetable, `{"stationId":48355,"from":"2016-05-29T15:03:22+02:00","to":"2016-05-29T17:03:22+02:00"}`)
	time.Sleep(time.Second)
	assert.True(suite.T(), passed)
}

func (suite *TimetableTestSuite) TestUpdateDelay() {
	t := train.Train{
		Arrival: &train.TimetableEvent{
			Time: time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone),
		},
	}
	del := 10
	initSettings(&t)
	updateDelay(&t, &del)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), t.Settings.Arrival.Time)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), t.Arrival.Time)
	assert.Equal(suite.T(), 10, t.Settings.Delay)
}

func (suite *TimetableTestSuite) TestUpdateEvent() {
	t := &train.TimetableEvent{
		Track:    10,
		Platform: "III",
	}
	tr := 5
	e := &msg.Event{Track: &tr}
	updateEvent(t, e)
	assert.Equal(suite.T(), "III", t.Platform)
	assert.Equal(suite.T(), 5, t.Track)
}

func (suite *TimetableTestSuite) TestInitSettings() {
	arrival := time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone)
	departure := time.Date(2016, time.April, 10, 15, 31, 0, 0, config.Timezone)
	t := train.Train{
		Arrival: &train.TimetableEvent{
			Time:     arrival,
			Platform: "II",
			Track:    5,
		},
		Departure: &train.TimetableEvent{
			Time:     departure,
			Platform: "I",
			Track:    2,
		},
	}
	initSettings(&t)
	assert.Equal(suite.T(), t.Arrival.Platform, t.Settings.Arrival.Platform)
	assert.Equal(suite.T(), t.Arrival.Time, t.Settings.Arrival.Time)
	assert.Equal(suite.T(), t.Arrival.Track, t.Settings.Arrival.Track)
	assert.Equal(suite.T(), t.Departure.Platform, t.Settings.Departure.Platform)
	assert.Equal(suite.T(), t.Departure.Time, t.Settings.Departure.Time)
	assert.Equal(suite.T(), t.Departure.Track, t.Settings.Departure.Track)
	t.Delay = 10
	initSettings(&t)
	assert.Equal(suite.T(), 10, t.Settings.Delay)
}

func (suite *TimetableTestSuite) TestSearch() {
	st := test.TrainStoreMock{}
	st.On("Search", "abcd").Return([]*train.Train{
		&train.Train{TrainID: "abcd"},
		&train.Train{TrainID: "efgh"},
		&train.Train{TrainID: "abcd"},
	}, nil)
	s := timetable{store: &st}
	res, err := s.SearchTrain("abcd", true)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), res, 2)
}

func (suite *TimetableTestSuite) TestProcessing() {
	arrival := time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone)
	departure := time.Date(2016, time.April, 10, 15, 31, 0, 0, config.Timezone)
	t := train.Train{
		Arrival: &train.TimetableEvent{
			Time:     arrival,
			Platform: "II",
			Track:    5,
		},
		Departure: &train.TimetableEvent{
			Time:     departure,
			Platform: "I",
			Track:    2,
		},
	}
	res := processTrain(t)
	assert.Equal(suite.T(), "II", res.Arrival.Platform)
	assert.Equal(suite.T(), arrival, res.Arrival.Time)
	assert.Equal(suite.T(), arrival, res.Arrival.ExpectedTime)
	assert.Equal(suite.T(), 2, res.Departure.Track)
	assert.Equal(suite.T(), departure, res.Departure.Time)
	assert.Equal(suite.T(), departure, res.Departure.ExpectedTime)
	initSettings(&t)
	t.Settings.Delay = 20
	t.Settings.Arrival.Platform = "IV"
	t.Settings.Mode = train.Manual
	res = processTrain(t)
	assert.Equal(suite.T(), "IV", res.Arrival.Platform)
	assert.Equal(suite.T(), arrival, res.Arrival.Time)
	assert.Equal(suite.T(), arrival.Add(20*time.Minute), res.Arrival.ExpectedTime)
	assert.Equal(suite.T(), 2, res.Departure.Track)
	assert.Equal(suite.T(), departure, res.Departure.Time)
	assert.Equal(suite.T(), departure.Add(20*time.Minute), res.Departure.ExpectedTime)
	assert.Equal(suite.T(), 20, res.Delay)
}

func TestTimetableTestSuite(t *testing.T) {
	suite.Run(t, new(TimetableTestSuite))
}

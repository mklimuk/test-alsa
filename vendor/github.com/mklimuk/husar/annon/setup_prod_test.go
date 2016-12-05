package annon

import (
	"testing"
	"time"

	"github.com/mklimuk/husar/config"
	"github.com/mklimuk/husar/test"
	"github.com/mklimuk/husar/train"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SetupProducerTestSuite struct {
	suite.Suite
	a       setup
	refTime time.Time
}

func (suite *SetupProducerTestSuite) SetupSuite() {
	suite.a = setup{}
	suite.refTime = time.Date(2016, time.April, 10, 0, 0, 0, 0, config.Timezone)
}

func (suite *SetupProducerTestSuite) SetupTest() {

}

func (suite *SetupProducerTestSuite) TestEnding() {
	t := test.GetTrain("ending")
	req, upd, err := suite.a.Required(t, nil, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), req)
	assert.False(suite.T(), upd)
	req, upd, err = suite.a.Required(t, t, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), req)
	assert.False(suite.T(), upd)
	t = test.GetTrain("passThrough")
	req, upd, err = suite.a.Required(t, t, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), req)
	assert.False(suite.T(), upd)
}

func (suite *SetupProducerTestSuite) TestNewTrain() {
	t := test.GetTrain("starting")
	req, upd, err := suite.a.Required(t, nil, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), req)
	assert.True(suite.T(), upd)
}

func (suite *SetupProducerTestSuite) TestNoChanges() {
	t := test.GetTrain("starting")
	req, upd, err := suite.a.Required(t, t, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), req)
	assert.False(suite.T(), upd)
}

func (suite *SetupProducerTestSuite) TestDeleteTrain() {
	req, upd, err := suite.a.Required(nil, &train.Train{}, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), req)
	assert.True(suite.T(), upd)
}

func (suite *SetupProducerTestSuite) TestFirstRealtimeDelay() {
	t1 := test.GetTrain("starting")
	t2 := test.GetTrain("starting")
	t1.Delay = 5
	req, upd, err := suite.a.Required(t1, t2, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), req)
	assert.True(suite.T(), upd)
}

func (suite *SetupProducerTestSuite) TestChangeDelay() {
	t1 := test.GetTrain("starting")
	t2 := test.GetTrain("starting")
	t2.Delay = 5
	t1.Delay = 10
	req, upd, err := suite.a.Required(t1, t2, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), req)
	//delay change from 5 to 10
	assert.True(suite.T(), upd)
}

func (suite *SetupProducerTestSuite) TestRemoveDelay() {
	t1 := test.GetTrain("starting")
	t2 := test.GetTrain("starting")
	t2.Delay = 5
	req, upd, err := suite.a.Required(t1, t2, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), req)
	assert.True(suite.T(), upd)
}

func (suite *SetupProducerTestSuite) TestSwitchToManualMode() {
	//if the delay did not change we don't have to update
	t1 := test.GetTrain("starting")
	t2 := test.GetTrain("starting")
	t1.Settings = &train.Settings{
		Mode: train.Manual,
	}
	req, upd, err := suite.a.Required(t1, t2, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), req)
	assert.True(suite.T(), upd)
}

func (suite *SetupProducerTestSuite) TestUpdateManual() {
	t1 := test.GetTrain("starting")
	t2 := test.GetTrain("starting")
	t1.Settings = &train.Settings{
		Mode: train.Manual,
		Departure: &train.TimetableEvent{
			Platform: "I",
			Track:    5,
		},
	}
	req, upd, err := suite.a.Required(t1, t2, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), req)
	assert.True(suite.T(), upd)
	t2.Settings = &train.Settings{
		Mode: train.Manual,
		Departure: &train.TimetableEvent{
			Platform: "II",
			Track:    1,
		},
	}
	req, upd, err = suite.a.Required(t1, t2, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), req)
	assert.True(suite.T(), upd)
}

func (suite *SetupProducerTestSuite) TestManualDelay() {
	// new manual settings
	t1 := test.GetTrain("starting")
	t2 := test.GetTrain("starting")
	t1.Settings = &train.Settings{
		Mode:  train.Manual,
		Delay: 10,
	}
	req, upd, err := suite.a.Required(t1, t2, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), req)
	assert.True(suite.T(), upd)
	// old and new manual with delay change
	t2.Settings = &train.Settings{
		Mode: train.Manual,
	}
	req, upd, err = suite.a.Required(t1, t2, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), req)
	assert.True(suite.T(), upd)
	// old and new manual with delay change no impact
	t2.Settings.Delay = 15
	req, upd, err = suite.a.Required(t1, t2, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), req)
	assert.True(suite.T(), upd)
	// remove manual delay
	t1.Settings.Delay = 0
	req, upd, err = suite.a.Required(t1, t2, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), req)
	assert.True(suite.T(), upd)
}

func (suite *SetupProducerTestSuite) TestSwitchToAutoMode() {
	t1 := test.GetTrain("starting")
	t2 := test.GetTrain("starting")
	t1.Settings = &train.Settings{
		Mode:  train.Manual,
		Delay: 10,
	}
	t2.Settings = &train.Settings{
		Mode:  train.Auto,
		Delay: 10,
	}
	req, upd, err := suite.a.Required(t1, t2, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), req)
	// delay change from 10 to 0
	assert.True(suite.T(), upd)
	// if the delay did not change no update required
	t2.Delay = 10
	req, upd, err = suite.a.Required(t1, t2, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), req)
	// delay stays at 10
	assert.False(suite.T(), upd)
}

func (suite *SetupProducerTestSuite) TestSubroute() {
	t1 := test.GetTrain("starting")
	t1.Route.CurrentStationOnSubroute = true
	req, upd, err := suite.a.Required(t1, nil, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), req)
	assert.False(suite.T(), upd)
}

func (suite *SetupProducerTestSuite) TestPlatformChangeInAutoMode() {
	t1 := test.GetTrain("starting")
	t2 := test.GetTrain("starting")
	t1.Settings = &train.Settings{
		Mode: train.Manual,
		Departure: &train.TimetableEvent{
			Platform: "I",
			Track:    5,
		},
	}
	req, upd, err := suite.a.Required(t1, t2, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), req)
	assert.True(suite.T(), upd)
	t1.Settings = &train.Settings{
		Mode: train.Auto,
		Departure: &train.TimetableEvent{
			Platform: "II",
			Track:    1,
		},
	}
	t2.Settings = &train.Settings{
		Mode: train.Auto,
		Departure: &train.TimetableEvent{
			Platform: "II",
			Track:    2,
		},
	}
	req, upd, err = suite.a.Required(t1, t2, &suite.refTime)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), req)
	assert.True(suite.T(), upd)
}

func TestSetupProducerTestSuite(t *testing.T) {
	suite.Run(t, new(SetupProducerTestSuite))
}

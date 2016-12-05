package sync

import (
	"errors"
	"testing"
	"time"

	"github.com/mklimuk/husar/service"
	"github.com/mklimuk/husar/test"
	"github.com/mklimuk/husar/train"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type SettingsSyncTestSuite struct {
	suite.Suite
}

func (suite *SettingsSyncTestSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)
}

func (suite *SettingsSyncTestSuite) TestBreakOnErrors() {
	trainStore := &shared.TrainStoreMock{}
	trainStore.On("NextSettingsChange").Return(nil, nil, errors.New("just testing")).Times(10)
	st := train.Store(trainStore)
	sy := service{t: st}
	sy.settings()
	time.Sleep(time.Duration(500) * time.Millisecond)
	trainStore.AssertCalled(suite.T(), "NextSettingsChange")
}

func (suite *SettingsSyncTestSuite) TestSync() {
	trainStore := &shared.TrainStoreMock{}
	ser := &test.ServiceMock{}
	res := &train.Settings{}
	ser.On("ProcessSettingsChange", res, res).Return()
	trainStore.On("NextSettingsChange").Return(res, res, nil)
	st := train.Store(trainStore)
	s := service.Service(ser)
	sy := service{s: s, t: st}
	sy.settings()
	time.Sleep(time.Duration(500) * time.Millisecond)
	sy.settingsSync <- true
	time.Sleep(time.Duration(500) * time.Millisecond)
	trainStore.AssertCalled(suite.T(), "NextSettingsChange")
	ser.AssertCalled(suite.T(), "ProcessSettingsChange", res, res)
}

func TestSettingsSyncTestSuite(t *testing.T) {
	suite.Run(t, new(SettingsSyncTestSuite))
}

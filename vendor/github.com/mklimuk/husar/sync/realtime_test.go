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

type RealtimeSyncTestSuite struct {
	suite.Suite
}

func (suite *RealtimeSyncTestSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)
}

func (suite *RealtimeSyncTestSuite) TestBreakOnErrors() {
	trainStore := &shared.TrainStoreMock{}
	trainStore.On("NextRealtimeChange").Return(nil, nil, errors.New("just testing")).Times(10)
	st := train.Store(trainStore)
	sy := service{t: st}
	sy.realtime()
	time.Sleep(time.Duration(500) * time.Millisecond)
	trainStore.AssertCalled(suite.T(), "NextRealtimeChange")
}

func (suite *RealtimeSyncTestSuite) TestSync() {
	trainStore := &shared.TrainStoreMock{}
	ser := &test.ServiceMock{}
	real := &train.Realtime{}
	ser.On("ProcessRealtimeChange", real, real).Return()
	trainStore.On("NextRealtimeChange").Return(real, real, nil)
	st := train.Store(trainStore)
	s := service.Service(ser)
	sy := service{s: s, t: st}
	sy.realtime()
	time.Sleep(time.Duration(100) * time.Millisecond)
	sy.liveSync <- true
	time.Sleep(time.Duration(100) * time.Millisecond)
	trainStore.AssertCalled(suite.T(), "NextRealtimeChange")
	ser.AssertCalled(suite.T(), "ProcessRealtimeChange", real, real)
}

func TestRealtimeSyncTestSuite(t *testing.T) {
	suite.Run(t, new(RealtimeSyncTestSuite))
}

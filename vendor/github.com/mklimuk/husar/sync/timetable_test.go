package sync

import (
	"testing"
	"time"

	"github.com/mklimuk/husar/service"
	"github.com/mklimuk/husar/test"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TimetableSyncTestSuite struct {
	suite.Suite
}

func (suite *TimetableSyncTestSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)
}

func (suite *TimetableSyncTestSuite) TestSync() {
	ser := &test.ServiceMock{}
	ser.On("GenerateForPeriod", mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time")).Return()
	s := service.Service(ser)
	sy := service{s: s, interval: time.Duration(1) * time.Second, windowStart: 5, windowLength: 10}
	sy.timetable()
	time.Sleep(time.Duration(1500) * time.Millisecond)
	sy.scheduledSync <- true
	time.Sleep(time.Duration(500) * time.Millisecond)
	ser.AssertNumberOfCalls(suite.T(), "GenerateForPeriod", 2)
}

func TestTimetableSyncTestSuite(t *testing.T) {
	suite.Run(t, new(TimetableSyncTestSuite))
}

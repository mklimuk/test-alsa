package audio

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DeviceTestSuite struct {
	suite.Suite
}

func (suite *DeviceTestSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)
}

func (suite *DeviceTestSuite) TearDownSuite() {

}

func (suite *DeviceTestSuite) TestConstructor() {
	r := &RawDeviceMock{}
	d := NewPlaybackDevice(r, 1024)
	dev := d.(*dev)
	a := assert.New(suite.T())
	a.NotNil(dev.raw)
	a.Equal(1024, dev.bufferSize)
	a.Nil(dev.errors)
	a.Nil(dev.ctrl)
	a.Equal(0, dev.framesWrote)
}

func (suite *DeviceTestSuite) TestSyncPlaybackError() {
	r := &RawDeviceMock{}
	b := buffers.Byte
}

func (suite *DeviceTestSuite) TestConvertBuffers() {

}

func (suite *DeviceTestSuite) TestPlaybackOngoing() {

}

func (suite *DeviceTestSuite) TestRegularPlayback() {

}

func (suite *DeviceTestSuite) TestPlaybackInterrupt() {

}

func TestDeviceTestSuite(t *testing.T) {
	suite.Run(t, new(DeviceTestSuite))
}

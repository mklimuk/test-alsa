package audio

import (
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/mklimuk/test-alsa/config"
	"github.com/mklimuk/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PlaybackTestSuite struct {
	suite.Suite
}

func (suite *PlaybackTestSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)
}

func (suite *PlaybackTestSuite) TearDownSuite() {

}

func (suite *PlaybackTestSuite) TestConstructor() {
	p := New(&config.AudioConf{DeviceBuffer: 1024, PeriodFrames: 512, Periods: 2}, &FactoryMock{}, "").(*play)
	a := assert.New(suite.T())
	a.NotNil(p.bufParams)
	a.Equal(1024, p.bufParams.BufferFrames)
	a.Equal(512, p.bufParams.PeriodFrames)
	a.Equal(2, p.bufParams.Periods)
}

func (suite *PlaybackTestSuite) TestDeviceBusy() {
	pl := play{context: &StreamContext{Priority: 3}}
	b, p := pl.DeviceBusy()
	a := assert.New(suite.T())
	a.True(b)
	a.Equal(3, p)
	pl.context = nil
	b, p = pl.DeviceBusy()
	a.False(b)
	a.Equal(0, p)
}

func (suite *PlaybackTestSuite) TestConvertBuffers() {
	buf := []byte{0x0A, 0x00}
	buf16 := make([]int16, 1)
	convertBuffers(buf, buf16)
	assert.Equal(suite.T(), int16(0x000A), buf16[0])
}

func (suite *PlaybackTestSuite) TestPlaybackOngoing() {
	c := websocket.ConnectionMock{}
	c.On("ReadMessage").Return(textMessage, []byte(`{"priority": 2}`), nil)
	pl := play{context: &StreamContext{Priority: 3}}
	err := pl.PlayFromWsConnection(&c)
	assert.Error(suite.T(), err)
}

func (suite *PlaybackTestSuite) TestRegularPlayback() {
	d := &DeviceMock{}
	f := &FactoryMock{}
	f.On("New", mock.Anything, mock.Anything, mock.Anything).Return(d, nil)
	d.On("Write", mock.Anything).Return(2, nil)
	c := &websocket.ConnectionMock{}
	c.On("ReadMessage").Return(textMessage, []byte(`{"priority": 2}`), nil)
	c.On("WriteLoop", mock.Anything).After(time.Duration(2 * time.Second)).Return()
	c.On("ReadLoop").After(time.Duration(2 * time.Second)).Return()
	c.On("CloseWithCode", 1000).Return().Once()
	ctrl := make(chan bool)
	bin := make(chan []byte)
	str := make(chan string)
	c.On("Control").Return(ctrl)
	c.On("In").Return(bin, str)
	p := New(&config.AudioConf{DeviceBuffer: 2, PeriodFrames: 1, Periods: 2}, f, "").(*play)
	err := p.PlayFromWsConnection(c)
	assert.NoError(suite.T(), err)
	time.Sleep(time.Duration(1 * time.Second))
	bin <- []byte{0x0A, 0x00, 0x01, 0x02}
	ctrl <- true
	d.AssertExpectations(suite.T())
	c.AssertExpectations(suite.T())
	assert.Nil(suite.T(), p.context)
}

func (suite *PlaybackTestSuite) TestPlaybackInterrupt() {

}

func TestPlaybackTestSuite(t *testing.T) {
	suite.Run(t, new(PlaybackTestSuite))
}

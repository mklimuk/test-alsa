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
	c.On("CloseWithReason", websocket.CloseTryAgainLater, mock.AnythingOfType("string")).Return().Once()
	pl := play{context: &StreamContext{Priority: 3}}
	pl.PlayFromWsConnection(&c)
	c.AssertExpectations(suite.T())
}

func (suite *PlaybackTestSuite) TestInterruptBeforeBufferFull() {
	d := &DeviceMock{}
	d.On("Close").Return().Once()
	d.On("FramesWrote").Return(0)
	f := &FactoryMock{}
	f.On("New", mock.Anything, mock.Anything, mock.Anything).Return(d, nil).Once()
	c := &websocket.ConnectionMock{}
	c.On("ID").Return("ABCD")
	c.On("ReadMessage").Return(textMessage, []byte(`{"priority": 2}`), nil).Once()
	c.On("ReadLoop").After(time.Duration(2 * time.Second)).Return().Once()
	c.On("CloseWithCode", websocket.CloseNormalClosure).Return()
	ctrl := make(chan bool)
	bin := make(chan []byte)
	str := make(chan string)
	c.On("Control").Return(ctrl)
	c.On("In").Return(bin, str).Once()
	p := New(&config.AudioConf{DeviceBuffer: 2, PeriodFrames: 1, Periods: 2, ReadBuffer: 6}, f, "").(*play)
	p.PlayFromWsConnection(c)
	time.Sleep(time.Duration(100 * time.Millisecond))
	bin <- []byte{0x0A, 0x00, 0x01, 0x02}
	ctrl <- true
	time.Sleep(time.Duration(100 * time.Millisecond))
	d.AssertNotCalled(suite.T(), "WriteAsync", mock.AnythingOfType("chan []int16"))
	c.AssertExpectations(suite.T())
	assert.Nil(suite.T(), p.context)
}

func (suite *PlaybackTestSuite) TestDevicePlayback() {
	d := &DeviceMock{}
	d.On("Close").Return().Once()
	d.On("FramesWrote").Return(10)
	e := make(chan error)
	d.On("WriteAsync", mock.AnythingOfType("chan []int16")).Run(func(args mock.Arguments) {
		in := args.Get(0).(chan []int16)
		for i := 0; i < 2; i++ {
			<-in
		}
	}).Return(e)
	f := &FactoryMock{}
	f.On("New", mock.Anything, mock.Anything, mock.Anything).Return(d, nil).Once()
	c := &websocket.ConnectionMock{}
	c.On("ID").Return("ABCD")
	c.On("ReadMessage").Return(textMessage, []byte(`{"priority": 2}`), nil).Once()
	c.On("ReadLoop").After(time.Duration(2 * time.Second)).Return().Once()
	c.On("CloseWithCode", websocket.CloseNormalClosure).Return()
	ctrl := make(chan bool)
	bin := make(chan []byte)
	str := make(chan string)
	c.On("Control").Return(ctrl)
	c.On("In").Return(bin, str).Once()
	p := New(&config.AudioConf{DeviceBuffer: 2, PeriodFrames: 1, Periods: 2, ReadBuffer: 2}, f, "").(*play)
	p.PlayFromWsConnection(c)
	bin <- []byte{0x00, 0x01}
	bin <- []byte{0x02, 0x03}
	bin <- []byte{0x04, 0x05}
	bin <- []byte{0x06, 0x07}
	time.Sleep(time.Duration(100 * time.Millisecond))
	ctrl <- true
	time.Sleep(time.Duration(10 * time.Millisecond))
	c.AssertExpectations(suite.T())
	assert.Nil(suite.T(), p.context)
}

func (suite *PlaybackTestSuite) TestPlaybackInterrupt() {

}

func TestPlaybackTestSuite(t *testing.T) {
	suite.Run(t, new(PlaybackTestSuite))
}

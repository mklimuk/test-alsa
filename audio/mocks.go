package audio

import (
	"github.com/mklimuk/websocket"
	"github.com/stretchr/testify/mock"
)

//DeviceMock is a mock of RegistryClient interface
type DeviceMock struct {
	mock.Mock
}

//Write is a mocked method
func (m *DeviceMock) Write(buffer interface{}) (samples int, err error) {
	args := m.Called(buffer)
	return args.Int(0), args.Error(1)
}

//Close is a mocked method
func (m *DeviceMock) Close() {
	m.Called()
}

//PlaybackMock is a mock of the Playback interface
type PlaybackMock struct {
	mock.Mock
}

//DeviceBusy is a mocked method
func (p *PlaybackMock) DeviceBusy() (bool, int) {
	args := p.Called()
	return args.Bool(0), args.Int(1)
}

//BufferSize is a mocked method
func (p *PlaybackMock) BufferSize() int {
	args := p.Called()
	return args.Int(0)
}

//PlaybackContext is a mocked method
func (p *PlaybackMock) PlaybackContext() *StreamContext {
	args := p.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*StreamContext)
}

//PlayFromWsConnection is a mocked method
func (p *PlaybackMock) PlayFromWsConnection(c websocket.Connection) error {
	args := p.Called(c)
	return args.Error(0)
}

//Close is a mocked method
func (p *PlaybackMock) Close() {
	p.Called()
}

//FactoryMock is a mock of the DeficeFactory interface
type FactoryMock struct {
	mock.Mock
}

//New is a mocked method
func (f *FactoryMock) New(sampleRate int, channels int, bp *BufferParams) (PlaybackDevice, error) {
	args := f.Called(sampleRate, channels, bp)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(PlaybackDevice), args.Error(1)
}

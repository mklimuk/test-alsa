package audio

import "github.com/stretchr/testify/mock"

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

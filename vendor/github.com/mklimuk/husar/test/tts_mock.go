package test

import "github.com/stretchr/testify/mock"

//TTSMock is a mock of AudioCatalog
type TTSMock struct {
	mock.Mock
}

//GetAudio is a mocked method
func (m *TTSMock) GetAudio(id string, text string) (respBody []byte, length int, err error) {
	args := m.Called(id, text)
	return args.Get(0).([]byte), args.Int(1), args.Error(2)
}

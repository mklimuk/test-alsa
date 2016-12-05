package test

import "github.com/stretchr/testify/mock"

//AudioCatalogMock is a mock of AudioCatalog
type AudioCatalogMock struct {
	mock.Mock
}

//Get is a mocked method
func (m *AudioCatalogMock) Get(ID string) (audio []byte, err error) {
	args := m.Called(ID)
	return args.Get(0).([]byte), args.Error(1)
}

//GetPath is a mocked method
func (m *AudioCatalogMock) GetPath(ID string) (path string, err error) {
	args := m.Called(ID)
	return args.Get(0).(string), args.Error(1)
}

//GetID is a mocked method
func (m *AudioCatalogMock) GetID(text string) string {
	args := m.Called(text)
	return args.String(0)
}

//Generate is a mocked method
func (m *AudioCatalogMock) Generate(text string) (ID string, exists bool, duration int, approxLen bool, err error) {
	args := m.Called(text)
	return args.String(0), args.Bool(1), args.Int(2), args.Bool(3), args.Error(4)
}

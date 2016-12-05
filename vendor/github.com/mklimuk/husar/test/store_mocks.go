package test

import (
	"time"

	"github.com/mklimuk/husar/train"

	"github.com/stretchr/testify/mock"
)

//TrainStoreMock is a mocked train.Store interface
type TrainStoreMock struct {
	mock.Mock
}

//Get is a mocked method
func (m *TrainStoreMock) Get(id string) (*train.Train, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*train.Train), args.Error(1)
}

//NextSettingsChange is a mocked method
func (m *TrainStoreMock) NextSettingsChange() (*train.Settings, *train.Settings, error) {
	args := m.Called()
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil, args.Error(2)
		}
		return nil, args.Get(1).(*train.Settings), args.Error(2)
	}
	if args.Get(1) == nil {
		return args.Get(0).(*train.Settings), nil, args.Error(2)
	}
	return args.Get(0).(*train.Settings), args.Get(1).(*train.Settings), args.Error(2)
}

//GetRealtime is a mocked method
func (m *TrainStoreMock) GetRealtime(id string) (*train.Realtime, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*train.Realtime), args.Error(1)
}

//Search is a mocked method
func (m *TrainStoreMock) Search(query string) (t []*train.Train, err error) {
	args := m.Called(query)
	return args.Get(0).([]*train.Train), args.Error(1)
}

//SaveSettings is a mocked method
func (m *TrainStoreMock) SaveSettings(s *train.Settings) (*train.Settings, error) {
	args := m.Called(s)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*train.Settings), args.Error(1)
}

//NextRealtimeChange is a mocked method
func (m *TrainStoreMock) NextRealtimeChange() (*train.Realtime, *train.Realtime, error) {
	args := m.Called()
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil, args.Error(2)
		}
		return nil, args.Get(1).(*train.Realtime), args.Error(2)
	}
	if args.Get(1) == nil {
		return args.Get(0).(*train.Realtime), nil, args.Error(2)
	}
	return args.Get(0).(*train.Realtime), args.Get(1).(*train.Realtime), args.Error(2)
}

//SaveAllRealtime is a mocked method
func (m *TrainStoreMock) SaveAllRealtime(re *[]train.Realtime) error {
	args := m.Called(re)
	return args.Error(0)
}

//GetBetween is a mocked method
func (m *TrainStoreMock) GetBetween(start, end *time.Time) (trains []train.Train, err error) {
	args := m.Called(start, end)
	return args.Get(0).([]train.Train), args.Error(1)
}

//GetBetweenWithRealtime is a mocked method
func (m *TrainStoreMock) GetBetweenWithRealtime(start, end *time.Time) (res []train.Train, IDs []string, err error) {
	args := m.Called(start, end)
	return args.Get(0).([]train.Train), args.Get(1).([]string), args.Error(2)
}

//GetAllSettings is a mocked method
func (m *TrainStoreMock) GetAllSettings(IDs ...string) (*[]train.Settings, error) {
	args := m.Called(IDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]train.Settings), args.Error(1)
}

//GetRealtimeBetween is a mocked method
func (m *TrainStoreMock) GetRealtimeBetween(start, end *time.Time) (set []train.Train, err error) {
	args := m.Called(start, end)
	return args.Get(0).([]train.Train), args.Error(1)
}

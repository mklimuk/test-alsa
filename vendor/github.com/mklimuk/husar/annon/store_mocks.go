package annon

import (
	"time"

	"github.com/stretchr/testify/mock"
)

//AnnonStoreMock is a mock of AnnonStore
type AnnonStoreMock struct {
	mock.Mock
}

//Get is a mocked method
func (m *AnnonStoreMock) Get(id string) (*Announcement, error) {
	args := m.Called(id)
	return args.Get(0).(*Announcement), args.Error(1)
}

//Save is a mocked method
func (m *AnnonStoreMock) Save(a *Announcement) (*Announcement, error) {
	args := m.Called(a)
	return args.Get(0).(*Announcement), args.Error(1)
}

//Delete is a mocked method
func (m *AnnonStoreMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

//DeleteForTrain is a mocked method
func (m *AnnonStoreMock) DeleteForTrain(trainID string) error {
	args := m.Called(trainID)
	return args.Error(0)
}

//SaveAll is a mocked method
func (m *AnnonStoreMock) SaveAll(a *[]*Announcement) error {
	args := m.Called(a)
	return args.Error(0)
}

//GetForTrain is a mocked method
func (m *AnnonStoreMock) GetForTrain(trainID string) (*[]*Announcement, error) {
	args := m.Called(trainID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]*Announcement), args.Error(1)
}

//NextChange is a mocked method
func (m *AnnonStoreMock) NextChange() (*Announcement, *Announcement, error) {
	args := m.Called()
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil, args.Error(2)
		}
		return nil, args.Get(1).(*Announcement), args.Error(2)
	}
	if args.Get(1) == nil {
		return args.Get(0).(*Announcement), nil, args.Error(2)
	}
	return args.Get(0).(*Announcement), args.Get(1).(*Announcement), args.Error(2)
}

//GetBetween is a mocked method
func (m *AnnonStoreMock) GetBetween(start, end *time.Time) (annons []Announcement, err error) {
	args := m.Called(start, end)
	return args.Get(0).([]Announcement), args.Error(1)
}

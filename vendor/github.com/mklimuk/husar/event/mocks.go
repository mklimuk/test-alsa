package event

import "github.com/stretchr/testify/mock"

//BusMock is the event bus interface mock
type BusMock struct {
	mock.Mock
}

//Subscribe is a mocked method
func (b *BusMock) Subscribe(topic Type, fn interface{}) error {
	args := b.Called(topic, fn)
	return args.Error(0)
}

//HasCallback is a mocked method
func (b *BusMock) HasCallback(topic Type) bool {
	args := b.Called(topic)
	return args.Bool(0)
}

//Unsubscribe is a mocked method
func (b *BusMock) Unsubscribe(topic Type, handler interface{}) error {
	args := b.Called(topic, handler)
	return args.Error(0)
}

//Publish is a mocked method
func (b *BusMock) Publish(topic Type, args ...interface{}) {
	b.Called(topic, args)
}

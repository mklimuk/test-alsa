package event

import (
	"fmt"
	"reflect"
)

//Bus represents event bus that calls handlers on events
type Bus interface {
	Subscribe(topic Type, fn interface{}) error
	HasCallback(topic Type) bool
	Unsubscribe(topic Type, handler interface{}) error
	Publish(topic Type, args ...interface{})
}

//bus - box for handlers and callbacks.
type bus struct {
	handlers map[Type][]*eventHandler
}

// Type represents event types
type Type string

type eventHandler struct {
	callBack reflect.Value
}

// New returns new bus with empty handlers.
func New() Bus {
	return Bus(&bus{
		make(map[Type][]*eventHandler),
	})
}

// doSubscribe handles the subscription logic and is utilized by the public Subscribe functions
func (bus *bus) doSubscribe(topic Type, fn interface{}, handler *eventHandler) error {
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		return fmt.Errorf("%s is not of type reflect.Func", reflect.TypeOf(fn).Kind())
	}
	bus.handlers[topic] = append(bus.handlers[topic], handler)
	return nil
}

// Subscribe subscribes to a topic.
// Returns error if `fn` is not a function.
func (bus *bus) Subscribe(topic Type, fn interface{}) error {
	return bus.doSubscribe(topic, fn, &eventHandler{
		reflect.ValueOf(fn),
	})
}

// HasCallback returns true if exists any callback subscribed to the topic.
func (bus *bus) HasCallback(topic Type) bool {
	_, ok := bus.handlers[topic]
	if ok {
		return len(bus.handlers[topic]) > 0
	}
	return false
}

// Unsubscribe removes callback defined for a topic.
// Returns error if there are no callbacks subscribed to the topic.
func (bus *bus) Unsubscribe(topic Type, handler interface{}) error {
	if _, ok := bus.handlers[topic]; ok && len(bus.handlers[topic]) > 0 {
		bus.removeHandler(topic, reflect.ValueOf(handler))
		return nil
	}
	return fmt.Errorf("topic %s doesn't exist", topic)
}

// Publish executes callback defined for a topic. Any additional argument will be transferred to the callback.
func (bus *bus) Publish(topic Type, args ...interface{}) {
	if handlers, ok := bus.handlers[topic]; ok {
		for _, handler := range handlers {
			bus.doPublish(handler, topic, args...)
		}
	}
}

func (bus *bus) doPublish(handler *eventHandler, topic Type, args ...interface{}) {
	passedArguments := bus.setUpPublish(topic, args...)
	handler.callBack.Call(passedArguments)
}

func (bus *bus) findHandlerIdx(topic Type, callback reflect.Value) int {
	if _, ok := bus.handlers[topic]; ok {
		for idx, handler := range bus.handlers[topic] {
			if handler.callBack == callback {
				return idx
			}
		}
	}
	return -1
}

func (bus *bus) removeHandler(topic Type, callback reflect.Value) {
	i := bus.findHandlerIdx(topic, callback)
	if i >= 0 {
		bus.handlers[topic] = append(bus.handlers[topic][:i], bus.handlers[topic][i+1:]...)
	}
}

func (bus *bus) setUpPublish(topic Type, args ...interface{}) (arguments []reflect.Value) {
	for _, arg := range args {
		arguments = append(arguments, reflect.ValueOf(arg))
	}
	return arguments
}

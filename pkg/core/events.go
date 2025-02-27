package core

import (
	"errors"

	"github.com/gookit/event"
)

// Event

type Event[T any] struct {
	event   event.Event
	EventId string
	Data    T
}

func (e *Event[T]) StopPropagation() {
	e.event.Abort(true)
}

func newEvent[T any](event event.Event, data T) Event[T] {
	return Event[T]{
		event:   event,
		EventId: event.Name(),
		Data:    data,
	}
}

// Event Manager

var em = event.NewManager("core")

func OnEvent[T any](pattern string, cb func(Event[T]) error) {
	em.On(pattern, event.ListenerFunc(func(evt event.Event) error {
		casted, ok := evt.Data()["data"].(T)
		if !ok {
			return errors.New("Event data is not of observed type " + pattern)
		}
		return cb(newEvent(evt, casted))
	}))
}

func EmitEvent[T any](eventId string, data T) {
	em.Fire(eventId, event.M{"data": data})
}

func EmitEventAsync[T any](eventId string, data T) {
	em.Async(eventId, event.M{"data": data})
}

func CloseEvents() {
	em.CloseWait()
}

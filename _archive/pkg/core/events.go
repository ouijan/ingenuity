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
		data := evt.Data()
		if data == nil {
			Log.Error("Event data is nil")
			return errors.New("Event data is nil")
		}
		value, ok := data["data"]
		if !ok {
			Log.Error("Event data does not contain 'data' key")
			return errors.New("Event data does not contain 'data' key")
		}
		casted, ok := value.(T)
		if !ok {
			Log.Error("Event data is not of observed type " + pattern)
			return errors.New("Event data is not of observed type " + pattern)
		}
		return cb(newEvent(evt, casted))
	}))
}

func OnEventCh[T any](pattern string, ch chan<- T) {
	OnEvent(pattern, func(evt Event[T]) error {
		ch <- evt.Data
		return nil
	})
}

func EmitEvent[T any](eventId string, data T) {
	em.Async(eventId, event.M{"data": data})
}

func EmitEventSync[T any](eventId string, data T) error {
	err, _ := em.Fire(eventId, event.M{"data": data})
	return err
}

func CloseEvents() {
	em.CloseWait()
}

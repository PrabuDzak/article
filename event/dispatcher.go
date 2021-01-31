package event

import (
	"context"
)

// SubscribeFunc represent function called when event dispatced
type SubscribeFunc func(ctx context.Context, event Event) error

// Dispatcher represent event dispatch
type Dispatcher interface {
	AddSubscriber(ctx context.Context, event Event, fn SubscribeFunc) error
	Dispatch(ctx context.Context, event Event) error
}

var d Dispatcher = &noOpDispatcher{}

// SetDispatcher set global event dispatcher
func SetDispatcher(dispatcher Dispatcher) {
	d = dispatcher
}

// Dispatch send an event to global dispatcher
func Dispatch(ctx context.Context, event Event) error {
	return d.Dispatch(ctx, event)
}

// AddSubscriber register a new event subscriber to global dispatcher
func AddSubscriber(ctx context.Context, event Event, fn SubscribeFunc) error {
	return d.AddSubscriber(ctx, event, fn)
}

type noOpDispatcher struct{}

func (n *noOpDispatcher) AddSubscriber(ctx context.Context, event Event, fn SubscribeFunc) error {
	return nil
}

func (n *noOpDispatcher) Dispatch(ctx context.Context, event Event) error {
	return nil
}

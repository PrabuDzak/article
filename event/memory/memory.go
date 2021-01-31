package memory

import (
	"context"
	"sync"

	"github.com/prabudzak/article/event"
)

type Dispatcher struct {
	subsriberMap map[string][]event.SubscribeFunc
	eventChan    chan event.Event

	processor int
	mutex     sync.Mutex
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		subsriberMap: make(map[string][]event.SubscribeFunc),
		eventChan:    make(chan event.Event),
		processor:    1,
	}
}

func (d *Dispatcher) AddSubscriber(ctx context.Context, e event.Event, fn event.SubscribeFunc) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, ok := d.subsriberMap[e.String()]
	if !ok {
		d.subsriberMap[e.String()] = []event.SubscribeFunc{}
	}

	d.subsriberMap[e.String()] = append(d.subsriberMap[e.String()], fn)
	return nil
}

func (d *Dispatcher) Dispatch(ctx context.Context, e event.Event) error {
	d.eventChan <- e
	return nil
}

func (d *Dispatcher) Start() {
	for i := 0; i < d.processor; i++ {
		go d.start()
	}
}

func (d *Dispatcher) start() {
	for e := range d.eventChan {
		subs, ok := d.subsriberMap[e.String()]
		if !ok {
			continue
		}

		ctx := context.Background()
		for _, sub := range subs {
			go sub(ctx, e)
		}
	}
}

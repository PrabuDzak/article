package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/prabudzak/article/event"
	"github.com/prabudzak/article/event/memory"
	"github.com/stretchr/testify/assert"
)

type eventIncreaseCount struct{}

func (e eventIncreaseCount) String() string {
	return "event_increase_count"
}

type eventSetTitle struct {
	title string
}

func (e eventSetTitle) String() string {
	return "event_set_title"
}

type subscriber struct {
	triggerCount int
	title        string
}

func (s *subscriber) increaseCount(ctx context.Context, e event.Event) error {
	s.triggerCount++
	return nil
}

func (s *subscriber) setTitle(ctx context.Context, e event.Event) error {
	message := e.(eventSetTitle)
	s.title = message.title
	return nil
}

func TestMemoryDispatch(t *testing.T) {
	ctx := context.Background()
	sub := &subscriber{}

	dispatcher := memory.NewDispatcher()
	dispatcher.Start()

	dispatcher.AddSubscriber(ctx, eventIncreaseCount{}, sub.increaseCount)

	dispatcher.Dispatch(ctx, eventIncreaseCount{})
	dispatcher.Dispatch(ctx, eventIncreaseCount{})
	dispatcher.Dispatch(ctx, eventIncreaseCount{})
	dispatcher.Dispatch(ctx, eventSetTitle{title: "event with no subscriber"})

	// wait until subscriber done processing
	time.Sleep(4 * time.Millisecond)
	assert.Equal(t, 3, sub.triggerCount)

	otherSub := &subscriber{}
	dispatcher.AddSubscriber(ctx, eventIncreaseCount{}, otherSub.increaseCount)

	dispatcher.Dispatch(ctx, eventIncreaseCount{})

	// wait until subscriber done processing
	time.Sleep(4 * time.Millisecond)
	assert.Equal(t, 4, sub.triggerCount)
	assert.Equal(t, 1, otherSub.triggerCount)
}

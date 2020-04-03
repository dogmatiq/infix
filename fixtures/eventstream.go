package fixtures

import (
	"context"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/infix/eventstream"
)

// EventStreamStub is a test implementation of the eventstream.Stream interface.
type EventStreamStub struct {
	eventstream.Stream

	ApplicationFunc func() configkit.Identity
	EventTypesFunc  func(context.Context) (message.TypeCollection, error)
	OpenFunc        func(context.Context, eventstream.Offset, message.TypeCollection) (eventstream.Cursor, error)
}

// Application returns the identity of the application that owns the stream.
func (s *EventStreamStub) Application() configkit.Identity {
	if s.ApplicationFunc != nil {
		return s.ApplicationFunc()
	}

	if s.Stream != nil {
		return s.Stream.Application()
	}

	return configkit.Identity{}
}

// EventTypes returns the set of event types that may appear on the stream.
func (s *EventStreamStub) EventTypes(ctx context.Context) (message.TypeCollection, error) {
	if s.EventTypesFunc != nil {
		return s.EventTypesFunc(ctx)
	}

	if s.Stream != nil {
		return s.Stream.EventTypes(ctx)
	}

	return nil, nil
}

// Open returns a cursor that reads events from the stream.
func (s *EventStreamStub) Open(
	ctx context.Context,
	offset eventstream.Offset,
	types message.TypeCollection,
) (eventstream.Cursor, error) {
	if s.OpenFunc != nil {
		return s.OpenFunc(ctx, offset, types)
	}

	if s.Stream != nil {
		return s.Stream.Open(ctx, offset, types)
	}

	return nil, nil
}

// EventStreamHandlerStub is a test implementation of the eventstream.Handler
// interface.
type EventStreamHandlerStub struct {
	eventstream.Handler

	NextOffsetFunc  func(context.Context, configkit.Identity) (eventstream.Offset, error)
	HandleEventFunc func(context.Context, eventstream.Offset, *eventstream.Event) error
}

// NextOffset returns the offset of the next event to be consumed from a
// specific application's event stream.
func (h *EventStreamHandlerStub) NextOffset(
	ctx context.Context,
	id configkit.Identity,
) (eventstream.Offset, error) {
	if h.NextOffsetFunc != nil {
		return h.NextOffsetFunc(ctx, id)
	}

	if h.Handler != nil {
		return h.Handler.NextOffset(ctx, id)
	}

	return 0, nil
}

// HandleEvent handles an event obtained from the event stream.
func (h *EventStreamHandlerStub) HandleEvent(
	ctx context.Context,
	o eventstream.Offset,
	ev *eventstream.Event,
) error {
	if h.HandleEventFunc != nil {
		return h.HandleEventFunc(ctx, o, ev)
	}

	if h.Handler != nil {
		return h.Handler.HandleEvent(ctx, o, ev)
	}

	return nil
}

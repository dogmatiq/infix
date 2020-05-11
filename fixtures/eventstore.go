package fixtures

import (
	"context"

	"github.com/dogmatiq/infix/persistence/subsystem/eventstore"
)

// EventStoreRepositoryStub is a test implementation of the
// eventstore.Repository interface.
type EventStoreRepositoryStub struct {
	eventstore.Repository

	QueryEventsFunc            func(context.Context, eventstore.Query) (eventstore.Result, error)
	LoadEventsForAggregateFunc func(context.Context, string, string, string) (eventstore.Result, error)
}

// LoadEventsForAggregate loads the events for the aggregate with the given
// key and id.
func (r *EventStoreRepositoryStub) LoadEventsForAggregate(
	ctx context.Context,
	hk, id, d string,
) (eventstore.Result, error) {
	if r.LoadEventsForAggregateFunc != nil {
		return r.LoadEventsForAggregateFunc(ctx, hk, id, d)
	}

	if r.Repository != nil {
		return r.Repository.LoadEventsForAggregate(ctx, hk, id, d)
	}

	return nil, nil
}

// QueryEvents queries events in the repository.
func (r *EventStoreRepositoryStub) QueryEvents(ctx context.Context, q eventstore.Query) (eventstore.Result, error) {
	if r.QueryEventsFunc != nil {
		return r.QueryEventsFunc(ctx, q)
	}

	if r.Repository != nil {
		return r.Repository.QueryEvents(ctx, q)
	}

	return nil, nil
}

// EventStoreResultStub is a test implementation of the eventstore.Result
// interface.
type EventStoreResultStub struct {
	eventstore.Result

	NextFunc  func(context.Context) (*eventstore.Item, bool, error)
	CloseFunc func() error
}

// Next returns the next event in the result.
func (r *EventStoreResultStub) Next(ctx context.Context) (*eventstore.Item, bool, error) {
	if r.NextFunc != nil {
		return r.NextFunc(ctx)
	}

	if r.Result != nil {
		return r.Result.Next(ctx)
	}

	return nil, false, nil
}

// Close closes the cursor.
func (r *EventStoreResultStub) Close() error {
	if r.CloseFunc != nil {
		return r.CloseFunc()
	}

	if r.Result != nil {
		return r.Result.Close()
	}

	return nil
}

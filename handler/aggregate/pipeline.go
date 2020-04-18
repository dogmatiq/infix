package aggregate

import (
	"context"
	"fmt"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/infix/draftspecs/envelopespec"
	"github.com/dogmatiq/infix/parcel"
	"github.com/dogmatiq/infix/persistence/subsystem/aggregatestore"
	"github.com/dogmatiq/infix/persistence/subsystem/eventstore"
	"github.com/dogmatiq/infix/pipeline"
	"github.com/dogmatiq/marshalkit"
)

// Sink is a pipeline sink that coordinates the handling of messages by a
// dogma.AggregateMessageHandler.
//
// The Accept() method conforms to the pipeline.Sink() signature.
type Sink struct {
	// Identity is the handler's identity.
	Identity *envelopespec.Identity

	// Handler is the aggregate message handler that implements the
	// application-specific message handling logic.
	Handler dogma.AggregateMessageHandler

	// AggregateStore is the repository used to load an aggregate instance's
	// revisions and snapshots.
	AggregateStore aggregatestore.Repository

	// EventStore is the repository used to load an aggregate instance's
	// historical events.
	EventStore eventstore.Repository

	// Marshaler is used to marshal/unmarshal aggregate snapshots and historical
	// events,
	Marshaler marshalkit.ValueMarshaler

	// Packer is used to create new parcels for events recorded by the
	// handler.
	Packer *parcel.Packer

	// Logger is the target for log messages produced within the handler.
	// If it is nil, logging.DefaultLogger is used.
	Logger logging.Logger
}

// Accept handles a message using s.Handler.
func (s *Sink) Accept(
	ctx context.Context,
	req pipeline.Request,
	res *pipeline.Response,
) error {
	p, err := req.Parcel()
	if err != nil {
		return err
	}

	id := s.Handler.RouteCommandToInstance(p.Message)
	if id == "" {
		panic(fmt.Sprintf(
			"the '%s' aggregate message handler attempted to route a %T command to an empty instance ID",
			s.Identity.Name,
			p.Message,
		))
	}

	rev, err := s.AggregateStore.LoadRevision(ctx, s.Identity.Key, id)
	if err != nil {
		return err
	}

	root := s.Handler.New()
	if root == nil {
		panic(fmt.Sprintf(
			"the '%s' aggregate message handler returned a nil root from New()",
			s.Identity.Name,
		))
	}

	if rev > 0 {
		if err := s.load(ctx, root, id); err != nil {
			return err
		}
	}

	sc := &scope{
		cause:   p,
		packer:  s.Packer,
		handler: s.Identity,
		logger:  s.Logger,

		id:     id,
		root:   root,
		exists: rev > 0,
	}

	s.Handler.HandleCommand(sc, p.Message)

	if len(sc.events) == 0 {
		if sc.created {
			panic(fmt.Sprintf(
				"the '%s' aggregate message handler created the '%s' instance without recording an event while handling a %T command",
				s.Identity.Name,
				id,
				p.Message,
			))
		}

		if sc.destroyed {
			panic(fmt.Sprintf(
				"the '%s' aggregate message handler destroyed the '%s' instance without recording an event while handling a %T command",
				s.Identity.Name,
				id,
				p.Message,
			))
		}

		return nil
	}

	tx, err := req.Tx(ctx)
	if err != nil {
		return err
	}

	for _, p := range sc.events {
		if _, err := res.RecordEvent(ctx, tx, p); err != nil {
			return err
		}
	}

	if !sc.exists {
		// TODO: https://github.com/dogmatiq/infix/issues/171
		panic("not implemented")
	}

	return tx.IncrementAggregateRevision(ctx, s.Identity.Key, id, rev)
}

// load applies an aggregate instance's historical events to the root in order
// to reproduce the current state.
func (s *Sink) load(
	ctx context.Context,
	root dogma.AggregateRoot,
	id string,
) error {
	q := eventstore.Query{
		// TODO: https://github.com/dogmatiq/dogma/issues/113
		// How do we best configure the filter to deal with events that the
		// aggregate once produced, but no longer does?
		Filter:              nil,
		AggregateHandlerKey: s.Identity.Key,
		AggregateInstanceID: id,
	}

	res, err := s.EventStore.QueryEvents(ctx, q)
	if err != nil {
		return err
	}

	for {
		i, ok, err := res.Next(ctx)
		if !ok || err != nil {
			return err
		}

		p, err := parcel.FromEnvelope(s.Marshaler, i.Envelope)
		if err != nil {
			return err
		}

		root.ApplyEvent(p.Message)
	}
}

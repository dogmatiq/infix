package integration

import (
	"context"
	"time"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/dodeca/logging"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/infix/draftspecs/envelopespec"
	"github.com/dogmatiq/infix/internal/mlog"
	"github.com/dogmatiq/infix/parcel"
	"github.com/dogmatiq/infix/pipeline"
	"github.com/dogmatiq/linger"
)

// Sink is a pipeline sink that coordinates the handling of messages by a
// dogma.IntegrationMessageHandler.
//
// The Accept() method conforms to the pipeline.Sink() signature.
type Sink struct {
	// Identity is the handler's identity.
	Identity *envelopespec.Identity

	// Handler is the integration message handler that implements the
	// application-specific message handling logic.
	Handler dogma.IntegrationMessageHandler

	// Packer is used to create new parcels for events recorded by the
	// handler.
	Packer *parcel.Packer

	// DefaultTimeout is the timeout to apply when handling the message if the
	// handler does not provide a timeout hint.
	DefaultTimeout time.Duration

	// Logger is the target for log messages produced within the handler.
	// If it is nil, logging.DefaultLogger is used.
	Logger logging.Logger
}

// Accept handles a message using s.Handler.
func (s *Sink) Accept(
	ctx context.Context,
	req pipeline.Request,
	res *pipeline.Response,
) (err error) {
	defer mlog.LogHandlerResult(
		s.Logger,
		req.Envelope(),
		s.Identity,
		configkit.IntegrationHandlerType,
		&err,
		"",
	)

	p, err := req.Parcel()
	if err != nil {
		return err
	}

	sc := &scope{
		cause:   p,
		packer:  s.Packer,
		handler: s.Identity,
		logger:  s.Logger,
	}

	if err := s.handle(ctx, sc, p.Message); err != nil {
		return err
	}

	tx, err := req.Tx(ctx)
	if err != nil {
		return err
	}

	for _, p := range sc.events {
		if _, err := res.RecordEventX(ctx, tx, p); err != nil {
			return err
		}
	}

	return nil
}

// handle invokes the handler with a timeout based on its timeout hint.
func (s *Sink) handle(ctx context.Context, sc *scope, m dogma.Message) error {
	ctx, cancel := linger.ContextWithTimeout(
		ctx,
		s.Handler.TimeoutHint(m),
		s.DefaultTimeout,
	)
	defer cancel()

	return s.Handler.HandleCommand(ctx, sc, m)
}

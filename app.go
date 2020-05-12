package infix

import (
	"context"
	"fmt"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dodeca/logging"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/infix/draftspecs/envelopespec"
	"github.com/dogmatiq/infix/eventstream"
	"github.com/dogmatiq/infix/eventstream/memorystream"
	"github.com/dogmatiq/infix/eventstream/persistedstream"
	"github.com/dogmatiq/infix/handler/aggregate"
	"github.com/dogmatiq/infix/handler/cache"
	"github.com/dogmatiq/infix/internal/x/loggingx"
	"github.com/dogmatiq/infix/parcel"
	"github.com/dogmatiq/infix/persistence"
	"github.com/dogmatiq/infix/pipeline"
	"github.com/dogmatiq/infix/queue"
)

type app struct {
	Config      configkit.RichApplication
	DataStore   persistence.DataStore
	EventCache  *memorystream.Stream
	EventStream *persistedstream.Stream
	Queue       *queue.Queue
	Pipeline    pipeline.Pipeline
	Logger      logging.Logger
}

// initApp initializes the engine to handle the app represented by cfg.
func (e *Engine) initApp(
	ctx context.Context,
	cfg configkit.RichApplication,
) error {
	logging.Log(
		e.logger,
		"hosting @%s application, identity key is %s",
		cfg.Identity().Name,
		cfg.Identity().Key,
	)

	id := cfg.Identity()

	ds, err := e.dataStores.Get(ctx, id.Key)
	if err != nil {
		return fmt.Errorf(
			"unable to open data-store for %s: %w",
			cfg.Identity(),
			err,
		)
	}

	c, err := e.newEventCache(ctx, cfg, ds)
	if err != nil {
		return err
	}

	q := e.newQueue(cfg, ds)
	s := e.newEventStream(cfg, ds, c)
	x := e.newCommandExecutor(cfg, q)
	l := loggingx.WithPrefix(
		e.opts.Logger,
		"@%s  ",
		cfg.Identity().Name,
	)
	p := e.newPipeline(cfg, ds, q, c, l)

	a := &app{
		Config:      cfg,
		DataStore:   ds,
		EventCache:  c,
		EventStream: s,
		Queue:       q,
		Pipeline:    p,
		Logger:      l,
	}

	if e.apps == nil {
		e.apps = map[string]*app{}
		e.executors = map[message.Type]dogma.CommandExecutor{}
	}

	e.apps[id.Key] = a

	for mt := range x.Packer.Produced {
		e.executors[mt] = x
	}

	return nil
}

// newQueue returns a new queue for a specific app.
func (e *Engine) newQueue(
	cfg configkit.RichApplication,
	ds persistence.DataStore,
) *queue.Queue {
	return &queue.Queue{
		DataStore: ds,
		Marshaler: e.opts.Marshaler,
		// TODO: https://github.com/dogmatiq/infix/issues/102
		// Make buffer size configurable.
		BufferSize: 0,
		Logger: loggingx.WithPrefix(
			e.logger,
			"[queue@%s] ",
			cfg.Identity().Name,
		),
	}
}

// newEventCache returns a new memory stream used as the event cache.
func (e *Engine) newEventCache(
	ctx context.Context,
	cfg configkit.RichApplication,
	ds persistence.DataStore,
) (*memorystream.Stream, error) {
	r := ds.EventStoreRepository()

	next, err := r.NextEventOffset(ctx)
	if err != nil {
		return nil, err
	}

	return &memorystream.Stream{
		App:         cfg.Identity(),
		FirstOffset: next,
		// TODO: https://github.com/dogmatiq/infix/issues/226
		// Make buffer size configurable.
		BufferSize: 0,
		Types: cfg.
			MessageTypes().
			Produced.
			FilterByRole(message.EventRole),
	}, nil
}

// newEventStream returns a new event stream for a specific app.
func (e *Engine) newEventStream(
	cfg configkit.RichApplication,
	ds persistence.DataStore,
	cache eventstream.Stream,
) *persistedstream.Stream {
	return &persistedstream.Stream{
		App:        cfg.Identity(),
		Repository: ds.EventStoreRepository(),
		Marshaler:  e.opts.Marshaler,
		Cache:      cache,
		// TODO: https://github.com/dogmatiq/infix/issues/76
		// Make pre-fetch buffer size configurable.
		PreFetch: 10,
		Types: cfg.
			MessageTypes().
			Produced.
			FilterByRole(message.EventRole),
	}
}

// newCommandExecutor returns a dogma.CommandExecutor for a specific app.
func (e *Engine) newCommandExecutor(
	cfg configkit.RichApplication,
	q *queue.Queue,
) *queue.CommandExecutor {
	return &queue.CommandExecutor{
		Queue: q,
		Packer: &parcel.Packer{
			Application: envelopespec.MarshalIdentity(cfg.Identity()),
			Marshaler:   e.opts.Marshaler,
			Produced: cfg.
				MessageTypes().
				Consumed.
				FilterByRole(message.CommandRole),
		},
	}
}

// newPipeline returns a new pipeline for a specific app.
func (e *Engine) newPipeline(
	cfg configkit.RichApplication,
	ds persistence.DataStore,
	q *queue.Queue,
	c *memorystream.Stream,
	l logging.Logger,
) pipeline.Pipeline {
	rf := &routeFactory{
		opts:         e.opts,
		appLogger:    l,
		engineLogger: e.logger,
		loader: &aggregate.Loader{
			AggregateStore: ds.AggregateStoreRepository(),
			EventStore:     ds.EventStoreRepository(),
			Marshaler:      e.opts.Marshaler,
		},
	}

	if err := cfg.AcceptRichVisitor(context.Background(), rf); err != nil {
		panic(err)
	}

	return pipeline.New(
		pipeline.NewSequence(
			pipeline.Observe(
				func(r pipeline.Result, err error) {
					if err == nil {
						c.Add(r.Events)
						for _, m := range r.QueueMessages {
							q.Track(m)
						}
					}
				},
			),
			pipeline.Acknowledge(
				e.opts.MessageBackoff,
				l,
			),
			pipeline.RouteByType(rf.routes),
		),
	)
}

// routeFactory is a configkit.RichVisitor that constructs the messaging
// pipeline.
type routeFactory struct {
	opts         *engineOptions
	loader       *aggregate.Loader
	engineLogger logging.Logger
	appLogger    logging.Logger

	app    *envelopespec.Identity
	routes map[message.Type]pipeline.Stage
}

func (f *routeFactory) VisitRichApplication(ctx context.Context, cfg configkit.RichApplication) error {
	f.app = envelopespec.MarshalIdentity(cfg.Identity())
	f.routes = map[message.Type]pipeline.Stage{}
	return cfg.RichHandlers().AcceptRichVisitor(ctx, f)
}

func (f *routeFactory) VisitRichAggregate(_ context.Context, cfg configkit.RichAggregate) error {
	s := &aggregate.Sink{
		Identity: envelopespec.MarshalIdentity(cfg.Identity()),
		Handler:  cfg.Handler(),
		Loader:   f.loader,
		Cache: &cache.Cache{
			// TODO: https://github.com/dogmatiq/infix/issues/193
			// Make TTL configurable.
			Logger: loggingx.WithPrefix(
				f.engineLogger,
				"[cache %s@%s] ",
				f.app.Name,
				cfg.Identity().Name,
			),
		},
		Packer: &parcel.Packer{
			Application: f.app,
			Marshaler:   f.opts.Marshaler,
			Produced:    cfg.MessageTypes().Produced,
			Consumed:    cfg.MessageTypes().Consumed,
		},
		LoadTimeout: f.opts.MessageTimeout,
		Logger:      f.appLogger,
	}

	for mt := range cfg.MessageTypes().Consumed {
		f.routes[mt] = pipeline.Terminate(s.Accept)
	}

	return nil
}

func (f *routeFactory) VisitRichProcess(_ context.Context, cfg configkit.RichProcess) error {
	return nil
}

func (f *routeFactory) VisitRichIntegration(_ context.Context, cfg configkit.RichIntegration) error {
	// a := &integration.Adaptor{
	// 	Identity:       envelopespec.MarshalIdentity(cfg.Identity()),
	// 	Handler:        cfg.Handler(),
	// 	DefaultTimeout: f.opts.MessageTimeout,
	// 	Packer: &parcel.Packer{
	// 		Application: f.app,
	// 		Marshaler:   f.opts.Marshaler,
	// 		Produced:    cfg.MessageTypes().Produced,
	// 		Consumed:    cfg.MessageTypes().Consumed,
	// 	},
	// 	Logger: f.appLogger,
	// }

	// for mt := range cfg.MessageTypes().Consumed {
	// f.routes[mt] = pipeline.Terminate(s.Accept)
	// }

	return nil
}

func (f *routeFactory) VisitRichProjection(_ context.Context, cfg configkit.RichProjection) error {
	return nil
}

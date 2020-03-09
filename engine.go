package infix

import (
	"context"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/api/discovery"
	"github.com/dogmatiq/dogma"
	"golang.org/x/sync/errgroup"
)

// Engine hosts a Dogma application.
type Engine struct {
	configs  []configkit.RichApplication
	opts     *engineOptions
	observer discovery.ApplicationObserverSet
}

// New returns a new engine that hosts the given application.
func New(app dogma.Application, options ...EngineOption) *Engine {
	cfg := configkit.FromApplication(app)

	return &Engine{
		configs: []configkit.RichApplication{cfg},
		opts:    resolveOptions(cfg, options),
	}
}

// Run hosts the given application until ctx is canceled or an error occurs.
func (e *Engine) Run(ctx context.Context) (err error) {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error { return e.serveAPI(ctx) })
	g.Go(func() error { return e.discover(ctx) })

	for _, cfg := range e.configs {
		cfg := cfg // capture loop variable
		g.Go(func() error { return e.hostApplication(ctx, cfg) })
	}

	err = g.Wait()

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return err
}

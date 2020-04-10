package infix_test

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/dogmatiq/infix"
	"github.com/dogmatiq/infix/persistence/provider/memory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ dogma.CommandExecutor = (*Engine)(nil)

var _ = Describe("type Engine", func() {
	var app *Application

	BeforeEach(func() {
		app = &Application{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app-name>", "<app-key>")
			},
		}
	})

	Describe("func New()", func() {
		It("allows the app to be provided as the first parameter", func() {
			Expect(func() {
				New(app)
			}).NotTo(Panic())
		})

		It("allows the app to be provided via the WithApplication() option", func() {
			Expect(func() {
				New(nil, WithApplication(app))
			}).NotTo(Panic())
		})

		It("panics if no apps are provided", func() {
			Expect(func() {
				New(nil)
			}).To(Panic())
		})
	})

	Describe("func Run()", func() {
		var (
			ctx    context.Context
			cancel context.CancelFunc
		)

		BeforeEach(func() {
			ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
		})

		AfterEach(func() {
			cancel()
		})

		It("returns an error if the context is canceled before calling", func() {
			cancel()

			err := Run(ctx, app)
			Expect(err).To(MatchError(context.Canceled))
		})

		It("returns an error if the context is canceled while running", func() {
			go func() {
				time.Sleep(10 * time.Millisecond)
				cancel()
			}()

			err := Run(
				ctx,
				app,
				WithPersistence(&memory.Provider{}), // avoid default BoltDB location
			)
			Expect(err).To(MatchError(context.Canceled))
		})
	})
})

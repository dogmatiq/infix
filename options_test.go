package infix

import (
	"time"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/dogmatiq/infix/retry"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogma/fixtures"
	"github.com/dogmatiq/marshalkit/codec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type EngineOption", func() {
	var cfg configkit.RichApplication

	BeforeSuite(func() {
		cfg = configkit.FromApplication(&Application{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app-name>", "<app-key>")
			},
		})
	})

	Describe("func ListenAddress()", func() {
		It("sets the listener address", func() {
			opts := resolveOptions(cfg, []EngineOption{
				ListenAddress("localhost:1234"),
			})

			Expect(opts.ListenAddress).To(Equal("localhost:1234"))
		})

		It("uses the default if the address is empty", func() {
			opts := resolveOptions(cfg, []EngineOption{
				ListenAddress(""),
			})

			Expect(opts.ListenAddress).To(Equal(DefaultListenAddress))
		})

		It("panics if the address is invalid", func() {
			Expect(func() {
				ListenAddress("missing-port")
			}).To(Panic())
		})

		It("panics if the post is an unknown service name", func() {
			Expect(func() {
				ListenAddress("host:xxx")
			}).To(Panic())
		})
	})

	Describe("func MessageTimeout()", func() {
		It("sets the message timeout", func() {
			opts := resolveOptions(cfg, []EngineOption{
				MessageTimeout(10 * time.Minute),
			})

			Expect(opts.MessageTimeout).To(Equal(10 * time.Minute))
		})

		It("uses the default if the duration is zero", func() {
			opts := resolveOptions(cfg, []EngineOption{
				MessageTimeout(0),
			})

			Expect(opts.MessageTimeout).To(Equal(DefaultMessageTimeout))
		})

		It("panics if the duration is less than zero", func() {
			Expect(func() {
				MessageTimeout(-1)
			}).To(Panic())
		})
	})

	Describe("func RetryPolicy()", func() {
		It("sets the retry policy", func() {
			p := &retry.ExponentialBackoff{}

			opts := resolveOptions(cfg, []EngineOption{
				RetryPolicy(p),
			})

			Expect(opts.RetryPolicy).To(BeIdenticalTo(p))
		})

		It("uses the default if the policy is nil", func() {
			opts := resolveOptions(cfg, []EngineOption{
				RetryPolicy(nil),
			})

			Expect(opts.RetryPolicy).To(Equal(DefaultRetryPolicy))
		})
	})

	Describe("func Marshaler()", func() {
		It("sets the marshaler", func() {
			m := &codec.Marshaler{}

			opts := resolveOptions(cfg, []EngineOption{
				Marshaler(m),
			})

			Expect(opts.Marshaler).To(BeIdenticalTo(m))
		})

		It("constructs a default if the marshaler is nil", func() {
			opts := resolveOptions(cfg, []EngineOption{
				Marshaler(nil),
			})

			Expect(opts.Marshaler).To(Equal(NewDefaultMarshaler(cfg)))
		})
	})

	Describe("func Logger()", func() {
		It("sets the logger", func() {
			opts := resolveOptions(cfg, []EngineOption{
				Logger(logging.DebugLogger),
			})

			Expect(opts.Logger).To(BeIdenticalTo(logging.DebugLogger))
		})

		It("uses the default if the logger is nil", func() {
			opts := resolveOptions(cfg, []EngineOption{
				Logger(nil),
			})

			Expect(opts.Logger).To(Equal(DefaultLogger))
		})
	})

	Describe("func Meter()", func() {
		It("sets the mter", func() {
			t := &metric.NoopMeter{}

			opts := resolveOptions(cfg, []EngineOption{
				Meter(t),
			})

			Expect(opts.Meter).To(BeIdenticalTo(t))
		})

		It("uses the default if the meter is nil", func() {
			opts := resolveOptions(cfg, []EngineOption{
				Meter(nil),
			})

			Expect(opts.Meter).To(Equal(DefaultMeter))
		})
	})

	Describe("func Tracer()", func() {
		It("sets the tracer", func() {
			t := &trace.NoopTracer{}

			opts := resolveOptions(cfg, []EngineOption{
				Tracer(t),
			})

			Expect(opts.Tracer).To(BeIdenticalTo(t))
		})

		It("uses the default if the tracer is nil", func() {
			opts := resolveOptions(cfg, []EngineOption{
				Tracer(nil),
			})

			Expect(opts.Tracer).To(Equal(DefaultTracer))
		})
	})
})

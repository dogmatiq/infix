package integration_test

import (
	"context"
	"time"

	. "github.com/dogmatiq/configkit/fixtures"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dodeca/logging"
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogma/fixtures"
	"github.com/dogmatiq/interopspec/envelopespec"
	. "github.com/dogmatiq/marshalkit/fixtures"
	. "github.com/dogmatiq/verity/fixtures"
	. "github.com/dogmatiq/verity/handler/integration"
	"github.com/dogmatiq/verity/parcel"
	. "github.com/jmalloc/gomegax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Adaptor", func() {
	var (
		ctx      context.Context
		upstream *IntegrationMessageHandler
		packer   *parcel.Packer
		logger   *logging.BufferedLogger
		work     *UnitOfWorkStub
		cause    parcel.Parcel
		adaptor  *Adaptor
	)

	BeforeEach(func() {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
		DeferCleanup(cancel)

		upstream = &IntegrationMessageHandler{
			ConfigureFunc: func(c dogma.IntegrationConfigurer) {
				c.Identity("<integration-name>", "27fb3936-6f88-4873-8c56-e6a1d01f027a")
				c.Routes(
					dogma.HandlesCommand[MessageC](),
					dogma.RecordsEvent[MessageE](),
				)
			},
		}

		packer = NewPacker(
			message.TypeRoles{
				MessageCType: message.CommandRole,
				MessageEType: message.EventRole,
			},
		)

		logger = &logging.BufferedLogger{}

		work = &UnitOfWorkStub{}

		cause = NewParcel("<consume>", MessageC1)

		adaptor = &Adaptor{
			Identity: &envelopespec.Identity{
				Name: "<integration-name>",
				Key:  "27fb3936-6f88-4873-8c56-e6a1d01f027a",
			},
			Handler: upstream,
			Packer:  packer,
			Logger:  logger,
		}
	})

	Describe("func HandleMessage()", func() {
		It("forwards the message to the handler", func() {
			called := false
			upstream.HandleCommandFunc = func(
				_ context.Context,
				_ dogma.IntegrationCommandScope,
				m dogma.Command,
			) error {
				called = true
				Expect(m).To(Equal(MessageC1))
				return nil
			}

			err := adaptor.HandleMessage(ctx, work, cause)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(called).To(BeTrue())
		})

		When("an event is recorded", func() {
			BeforeEach(func() {
				upstream.HandleCommandFunc = func(
					_ context.Context,
					s dogma.IntegrationCommandScope,
					_ dogma.Command,
				) error {
					s.RecordEvent(MessageE1)
					return nil
				}

				err := adaptor.HandleMessage(ctx, work, cause)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("adds the event to the unit-of-work", func() {
				Expect(work.Events).To(EqualX(
					[]parcel.Parcel{
						{
							Envelope: &envelopespec.Envelope{
								MessageId:         "0",
								CausationId:       "<consume>",
								CorrelationId:     "<correlation>",
								SourceApplication: packer.Application,
								SourceHandler:     adaptor.Identity,
								CreatedAt:         "2000-01-01T00:00:00Z",
								Description:       "{E1}",
								PortableName:      MessageEPortableName,
								MediaType:         MessageE1Packet.MediaType,
								Data:              MessageE1Packet.Data,
							},
							Message:   MessageE1,
							CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					},
				))
			})

			It("logs about the event", func() {
				Expect(logger.Messages()).To(ContainElement(
					logging.BufferedLogMessage{
						Message: "= 0  ∵ <consume>  ⋲ <correlation>  ▲    MessageE ● {E1}",
					},
				))
			})
		})

		When("a message is logged via the scope", func() {
			BeforeEach(func() {
				upstream.HandleCommandFunc = func(
					_ context.Context,
					s dogma.IntegrationCommandScope,
					_ dogma.Command,
				) error {
					s.Log("format %s", "<value>")
					return nil
				}

				err := adaptor.HandleMessage(ctx, work, cause)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("logs using the standard format", func() {
				Expect(logger.Messages()).To(ContainElement(
					logging.BufferedLogMessage{
						Message: "= <consume>  ∵ <cause>  ⋲ <correlation>  ▼    MessageC ● format <value>",
					},
				))
			})
		})
	})
})

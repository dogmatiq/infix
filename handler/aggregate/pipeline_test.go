package aggregate_test

import (
	"context"
	"errors"
	"time"

	. "github.com/dogmatiq/configkit/fixtures"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dodeca/logging"
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogma/fixtures"
	"github.com/dogmatiq/infix/draftspecs/envelopespec"
	. "github.com/dogmatiq/infix/fixtures"
	. "github.com/dogmatiq/infix/handler/aggregate"
	"github.com/dogmatiq/infix/handler/cache"
	"github.com/dogmatiq/infix/parcel"
	"github.com/dogmatiq/infix/persistence"
	"github.com/dogmatiq/infix/persistence/subsystem/aggregatestore"
	"github.com/dogmatiq/infix/persistence/subsystem/eventstore"
	"github.com/dogmatiq/infix/pipeline"
	. "github.com/dogmatiq/marshalkit/fixtures"
	. "github.com/jmalloc/gomegax"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Sink", func() {
	var (
		ctx           context.Context
		cancel        context.CancelFunc
		tx            *TransactionStub
		dataStore     *DataStoreStub
		aggregateRepo *AggregateStoreRepositoryStub
		eventRepo     *EventStoreRepositoryStub
		req           *PipelineRequestStub
		res           *pipeline.Response
		handler       *AggregateMessageHandler
		packer        *parcel.Packer
		logger        *logging.BufferedLogger
		sink          *Sink
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)

		dataStore = NewDataStoreStub()
		aggregateRepo = dataStore.AggregateStoreRepository().(*AggregateStoreRepositoryStub)
		eventRepo = dataStore.EventStoreRepository().(*EventStoreRepositoryStub)

		req, tx = NewPipelineRequestStub(
			NewParcel("<consume>", MessageC1),
			dataStore,
		)
		res = &pipeline.Response{}

		handler = &AggregateMessageHandler{
			ConfigureFunc: func(c dogma.AggregateConfigurer) {
				c.Identity("<aggregate-name>", "<aggregate-key>")
				c.ConsumesCommandType(MessageC{})
				c.ProducesEventType(MessageE{})
			},
			RouteCommandToInstanceFunc: func(m dogma.Message) string {
				return "<instance>"
			},
		}

		packer = NewPacker(
			message.TypeRoles{
				MessageCType: message.CommandRole,
				MessageEType: message.EventRole,
			},
		)

		logger = &logging.BufferedLogger{}

		handler.NewFunc = func() dogma.AggregateRoot {
			return &AggregateRoot{
				Value: &[]dogma.Message{},
				ApplyEventFunc: func(m dogma.Message, v interface{}) {
					p := v.(*[]dogma.Message)
					*p = append(*p, m)
				},
			}
		}

		sink = &Sink{
			Identity: &envelopespec.Identity{
				Name: "<aggregate-name>",
				Key:  "<aggregate-key>",
			},
			Handler: handler,
			Loader: &Loader{
				AggregateStore: aggregateRepo,
				EventStore:     eventRepo,
				Marshaler:      Marshaler,
			},
			Cache:       &cache.Cache{},
			Packer:      packer,
			LoadTimeout: 1 * time.Second,
			Logger:      logger,
		}
	})

	AfterEach(func() {
		if req != nil {
			req.Close()
		}

		if dataStore != nil {
			dataStore.Close()
		}

		cancel()
	})

	Describe("func Accept()", func() {
		It("forwards the message to the handler", func() {
			called := false
			handler.HandleCommandFunc = func(
				_ dogma.AggregateCommandScope,
				m dogma.Message,
			) {
				called = true
				Expect(m).To(Equal(MessageC1))
			}

			err := sink.Accept(ctx, req, res)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(called).To(BeTrue())
		})

		It("returns an error if the message cannot be unpacked", func() {
			req.ParcelFunc = func() (*parcel.Parcel, error) {
				return nil, errors.New("<error>")
			}

			err := sink.Accept(ctx, req, res)
			Expect(err).To(MatchError("<error>"))
		})

		It("returns an error if the instance can not be loaded", func() {
			aggregateRepo.LoadMetaDataFunc = func(
				context.Context,
				string,
				string,
			) (*aggregatestore.MetaData, error) {
				return nil, errors.New("<error>")
			}

			err := sink.Accept(ctx, req, res)
			Expect(err).To(MatchError("<error>"))
		})

		It("panics if the handler routes the message to an empty instance ID", func() {
			handler.RouteCommandToInstanceFunc = func(dogma.Message) string {
				return ""
			}

			Expect(func() {
				err := sink.Accept(ctx, req, res)
				Expect(err).ShouldNot(HaveOccurred())
			}).To(PanicWith("the '<aggregate-name>' aggregate message handler attempted to route a fixtures.MessageC command to an empty instance ID"))
		})

		It("panics if the handler returns a nil root", func() {
			handler.NewFunc = func() dogma.AggregateRoot {
				return nil
			}

			Expect(func() {
				err := sink.Accept(ctx, req, res)
				Expect(err).ShouldNot(HaveOccurred())
			}).To(PanicWith("the '<aggregate-name>' aggregate message handler returned a nil root from New()"))
		})

		It("returns an error if the deadline is exceeded while acquiring the cache record", func() {
			blockReq, _ := NewPipelineRequestStub(
				NewParcel("<blocking>", MessageC1),
				dataStore,
			)
			blockRes := &pipeline.Response{}
			defer blockReq.Close()

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			barrier := make(chan struct{})

			handler.HandleCommandFunc = func(
				dogma.AggregateCommandScope,
				dogma.Message,
			) {
				close(barrier) // let the test proceed, we now hold the lock on the record
				<-ctx.Done()   // don't unlock until the test assertions are complete
			}

			go sink.Accept(ctx, blockReq, blockRes)

			select {
			case <-barrier:
			case <-ctx.Done():
				Expect(ctx.Err()).ShouldNot(HaveOccurred())
			}

			ctx, cancel = context.WithTimeout(ctx, 20*time.Millisecond)
			defer cancel()

			err := sink.Accept(ctx, blockReq, blockRes)
			Expect(err).To(Equal(context.DeadlineExceeded))
		})

		When("the instance does not exist", func() {
			It("can be created", func() {
				handler.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					ok := s.Create()
					Expect(ok).To(BeTrue())

					s.RecordEvent(MessageE1)
				}

				err := sink.Accept(ctx, req, res)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("panics if the instance is created without recording an event", func() {
				handler.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					s.Create()
				}

				Expect(func() {
					err := sink.Accept(ctx, req, res)
					Expect(err).ShouldNot(HaveOccurred())
				}).To(PanicWith("the '<aggregate-name>' aggregate message handler created the '<instance>' instance without recording an event while handling a fixtures.MessageC command"))
			})
		})

		When("the instance exists", func() {
			BeforeEach(func() {
				createReq, _ := NewPipelineRequestStub(
					NewParcel("<created>", MessageC1),
					dataStore,
				)
				createRes := &pipeline.Response{}
				defer createReq.Close()

				handler.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					handler.HandleCommandFunc = nil
					s.Create()
					s.RecordEvent(MessageE1)
					s.RecordEvent(MessageE2)
				}

				err := sink.Accept(ctx, createReq, createRes)
				Expect(err).ShouldNot(HaveOccurred())

				_, err = createReq.Ack(ctx, nil)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("causes Create() to return false", func() {
				handler.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					ok := s.Create()
					Expect(ok).To(BeFalse())
				}

				err := sink.Accept(ctx, req, res)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("provides a root with the correct state", func() {
				handler.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					r := s.Root().(*AggregateRoot)
					Expect(r.Value).To(Equal(
						&[]dogma.Message{
							MessageE1,
							MessageE2,
						},
					))
				}

				err := sink.Accept(ctx, req, res)
				Expect(err).ShouldNot(HaveOccurred())
			})

			When("when the instance is subsequently destroyed", func() {
				BeforeEach(func() {
					destroyReq, _ := NewPipelineRequestStub(
						NewParcel("<create>", MessageC1),
						dataStore,
					)
					destroyRes := &pipeline.Response{}
					defer destroyReq.Close()

					handler.HandleCommandFunc = func(
						s dogma.AggregateCommandScope,
						_ dogma.Message,
					) {
						handler.HandleCommandFunc = nil
						s.RecordEvent(MessageE{Value: "<destroyed>"})
						s.Destroy()
					}

					err := sink.Accept(ctx, destroyReq, destroyRes)
					Expect(err).ShouldNot(HaveOccurred())

					_, err = destroyReq.Ack(ctx, nil)
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("resets the root state", func() {
					handler.HandleCommandFunc = func(
						s dogma.AggregateCommandScope,
						_ dogma.Message,
					) {
						// Create() should now return true once again.
						Expect(s.Create()).To(BeTrue())

						// The root itself should also have been reset to
						// as-new. For our test root that means the internal
						// slice of historical messages should be empty.
						r := s.Root().(*AggregateRoot)
						Expect(r.Value).To(Equal(
							&[]dogma.Message{},
						))

						// As per the Dogma API specification, we must record an
						// event whenever we call Create(). This is done after
						// the assertion that the root is empty otherwise we
						// would see this event in the state.
						s.RecordEvent(MessageE{Value: "<recreated>"})
					}

					err := sink.Accept(ctx, req, res)
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			It("panics if the instance is destroyed without recording an event", func() {
				handler.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					s.Destroy()
				}

				Expect(func() {
					err := sink.Accept(ctx, req, res)
					Expect(err).ShouldNot(HaveOccurred())
				}).To(PanicWith("the '<aggregate-name>' aggregate message handler destroyed the '<instance>' instance without recording an event while handling a fixtures.MessageC command"))
			})
		})

		When("events are recorded", func() {
			BeforeEach(func() {
				handler.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					s.Create()
					s.RecordEvent(MessageE1)
					s.RecordEvent(MessageE2)
				}
			})

			It("saves the recorded events", func() {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				err := sink.Accept(ctx, req, res)
				Expect(err).ShouldNot(HaveOccurred())

				_, err = req.Ack(ctx, nil)
				Expect(err).ShouldNot(HaveOccurred())

				res, err := eventRepo.QueryEvents(ctx, eventstore.Query{})
				Expect(err).ShouldNot(HaveOccurred())
				defer res.Close()

				i, ok, err := res.Next(ctx)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(ok).To(BeTrue())
				Expect(i).To(EqualX(
					&eventstore.Item{
						Offset: 0,
						Envelope: &envelopespec.Envelope{
							MetaData: &envelopespec.MetaData{
								MessageId:     "0",
								CausationId:   "<consume>",
								CorrelationId: "<correlation>",
								Source: &envelopespec.Source{
									Application: packer.Application,
									Handler:     sink.Identity,
									InstanceId:  "<instance>",
								},
								CreatedAt:   "2000-01-01T00:00:00Z",
								Description: "{E1}",
							},
							PortableName: MessageEPortableName,
							MediaType:    MessageE1Packet.MediaType,
							Data:         MessageE1Packet.Data,
						},
					},
				))

				i, ok, err = res.Next(ctx)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(ok).To(BeTrue())
				Expect(i).To(EqualX(
					&eventstore.Item{
						Offset: 1,
						Envelope: &envelopespec.Envelope{
							MetaData: &envelopespec.MetaData{
								MessageId:     "1",
								CausationId:   "<consume>",
								CorrelationId: "<correlation>",
								Source: &envelopespec.Source{
									Application: packer.Application,
									Handler:     sink.Identity,
									InstanceId:  "<instance>",
								},
								CreatedAt:   "2000-01-01T00:00:01Z",
								Description: "{E2}",
							},
							PortableName: MessageEPortableName,
							MediaType:    MessageE2Packet.MediaType,
							Data:         MessageE2Packet.Data,
						},
					},
				))
			})

			It("updates the instance's meta-data", func() {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				err := sink.Accept(ctx, req, res)
				Expect(err).ShouldNot(HaveOccurred())

				_, err = req.Ack(ctx, nil)
				Expect(err).ShouldNot(HaveOccurred())

				md, err := aggregateRepo.LoadMetaData(ctx, "<aggregate-key>", "<instance>")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(md).To(Equal(
					&aggregatestore.MetaData{
						HandlerKey:     "<aggregate-key>",
						InstanceID:     "<instance>",
						InstanceExists: true,
						Revision:       1,
					},
				))
			})

			It("returns an error if the transaction cannot be started", func() {
				req.TxFunc = func(
					context.Context,
				) (persistence.ManagedTransaction, error) {
					return nil, errors.New("<error>")
				}

				err := sink.Accept(ctx, req, res)
				Expect(err).To(MatchError("<error>"))
			})

			It("returns an error if an event can not be recorded", func() {
				tx.SaveEventFunc = func(
					context.Context,
					*envelopespec.Envelope,
				) (uint64, error) {
					return 0, errors.New("<error>")
				}

				err := sink.Accept(ctx, req, res)
				Expect(err).To(MatchError("<error>"))
			})

			It("returns an error if the meta-data can not be saved", func() {
				tx.SaveAggregateMetaDataFunc = func(
					context.Context,
					*aggregatestore.MetaData,
				) error {
					return errors.New("<error>")
				}

				err := sink.Accept(ctx, req, res)
				Expect(err).To(MatchError("<error>"))
			})
		})

		When("no events are recorded", func() {
			It("does not start a transaction", func() {
				req.TxFunc = func(
					context.Context,
				) (persistence.ManagedTransaction, error) {
					Fail("unexpected call")
					return nil, nil
				}

				err := sink.Accept(ctx, req, res)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})

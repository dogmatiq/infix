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
	"github.com/dogmatiq/infix/handler"
	. "github.com/dogmatiq/infix/handler/aggregate"
	"github.com/dogmatiq/infix/parcel"
	"github.com/dogmatiq/infix/persistence"
	"github.com/dogmatiq/infix/persistence/subsystem/aggregatestore"
	. "github.com/dogmatiq/marshalkit/fixtures"
	. "github.com/jmalloc/gomegax"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Sink", func() {
	var (
		ctx           context.Context
		cancel        context.CancelFunc
		dataStore     *DataStoreStub
		aggregateRepo *AggregateStoreRepositoryStub
		eventRepo     *EventStoreRepositoryStub
		upstream      *AggregateMessageHandler
		packer        *parcel.Packer
		logger        *logging.BufferedLogger
		work          *handler.UnitOfWork
		cause         *parcel.Parcel
		adaptor       *Adaptor
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)

		dataStore = NewDataStoreStub()
		aggregateRepo = dataStore.AggregateStoreRepository().(*AggregateStoreRepositoryStub)
		eventRepo = dataStore.EventStoreRepository().(*EventStoreRepositoryStub)

		upstream = &AggregateMessageHandler{
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

		upstream.NewFunc = func() dogma.AggregateRoot {
			return &AggregateRoot{
				Value: &[]dogma.Message{},
				ApplyEventFunc: func(m dogma.Message, v interface{}) {
					p := v.(*[]dogma.Message)
					*p = append(*p, m)
				},
			}
		}

		work = &handler.UnitOfWork{}

		cause = NewParcel("<consume>", MessageC1)

		adaptor = &Adaptor{
			Identity: &envelopespec.Identity{
				Name: "<aggregate-name>",
				Key:  "<aggregate-key>",
			},
			Handler: upstream,
			Loader: &Loader{
				AggregateStore: aggregateRepo,
				EventStore:     eventRepo,
				Marshaler:      Marshaler,
			},
			Packer:      packer,
			LoadTimeout: 1 * time.Second,
			Logger:      logger,
		}
	})

	AfterEach(func() {
		if dataStore != nil {
			dataStore.Close()
		}

		cancel()
	})

	Describe("func HandleMessage()", func() {
		It("forwards the message to the handler", func() {
			called := false
			upstream.HandleCommandFunc = func(
				_ dogma.AggregateCommandScope,
				m dogma.Message,
			) {
				called = true
				Expect(m).To(Equal(MessageC1))
			}

			err := adaptor.HandleMessage(ctx, work, cause)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(called).To(BeTrue())
		})

		It("makes the instance ID available via the scope ", func() {
			upstream.HandleCommandFunc = func(
				s dogma.AggregateCommandScope,
				_ dogma.Message,
			) {
				Expect(s.InstanceID()).To(Equal("<instance>"))
			}

			err := adaptor.HandleMessage(ctx, work, cause)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns an error if the instance can not be loaded", func() {
			aggregateRepo.LoadMetaDataFunc = func(
				context.Context,
				string,
				string,
			) (*aggregatestore.MetaData, error) {
				return nil, errors.New("<error>")
			}

			err := adaptor.HandleMessage(ctx, work, cause)
			Expect(err).To(MatchError("<error>"))
		})

		It("panics if the handler routes the message to an empty instance ID", func() {
			upstream.RouteCommandToInstanceFunc = func(dogma.Message) string {
				return ""
			}

			Expect(func() {
				err := adaptor.HandleMessage(ctx, work, cause)
				Expect(err).ShouldNot(HaveOccurred())
			}).To(PanicWith("*fixtures.AggregateMessageHandler.RouteCommandToInstance() returned an empty instance ID while routing a fixtures.MessageC command"))
		})

		It("panics if the handler returns a nil root", func() {
			upstream.NewFunc = func() dogma.AggregateRoot {
				return nil
			}

			Expect(func() {
				err := adaptor.HandleMessage(ctx, work, cause)
				Expect(err).ShouldNot(HaveOccurred())
			}).To(PanicWith("*fixtures.AggregateMessageHandler.New() returned nil"))
		})

		Context("when an event is recorded", func() {
			BeforeEach(func() {
			})

			It("adds the event to the unit-of-work", func() {
				upstream.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					s.Create()
					s.RecordEvent(MessageE1)
				}

				err := adaptor.HandleMessage(context.Background(), work, cause)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(work.Events).To(EqualX(
					[]*parcel.Parcel{
						&parcel.Parcel{
							Envelope: &envelopespec.Envelope{
								MetaData: &envelopespec.MetaData{
									MessageId:     "0",
									CausationId:   "<consume>",
									CorrelationId: "<correlation>",
									Source: &envelopespec.Source{
										Application: packer.Application,
										Handler:     adaptor.Identity,
										InstanceId:  "<instance>",
									},
									CreatedAt:   "2000-01-01T00:00:00Z",
									Description: "{E1}",
								},
								PortableName: MessageEPortableName,
								MediaType:    MessageE1Packet.MediaType,
								Data:         MessageE1Packet.Data,
							},
							Message:   MessageE1,
							CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					},
				))
			})

			It("applies the event to the aggregate root", func() {
				upstream.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					s.Create()
					s.RecordEvent(MessageE1)

					r := s.Root().(*AggregateRoot)
					Expect(r.Value).To(Equal(
						&[]dogma.Message{
							MessageE1,
						},
					))
				}

				err := adaptor.HandleMessage(context.Background(), work, cause)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("panics if the instance does not exist", func() {
				upstream.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					s.RecordEvent(MessageE1)
				}

				Expect(func() {
					err := adaptor.HandleMessage(context.Background(), work, cause)
					Expect(err).ShouldNot(HaveOccurred())
				}).To(PanicWith("can not record event against non-existent instance"))
			})

			It("logs about the event", func() {
				upstream.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					s.Create()
					s.RecordEvent(MessageE1)
				}

				err := adaptor.HandleMessage(context.Background(), work, cause)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(logger.Messages()).To(ContainElement(
					logging.BufferedLogMessage{
						Message: "= 0  ∵ <consume>  ⋲ <correlation>  ▲    MessageE ● {E1}",
					},
				))
			})
		})

		Context("when a message is logged via the scope", func() {
			BeforeEach(func() {
				upstream.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					s.Log("format %s", "<value>")
				}

				err := adaptor.HandleMessage(context.Background(), work, cause)
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

		When("the cache record is locked", func() {
			BeforeEach(func() {
				// Handle a message. The cache record should remain locked
				// until the observers in the unit-of-work are notified.
				err := adaptor.HandleMessage(ctx, work, cause)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("times-out attempting to acquire another lock", func() {
				ctx, cancel = context.WithTimeout(ctx, 20*time.Millisecond)
				defer cancel()

				// Because we haven't notified the observers to release the cache
				// record, this invocation should timeout.
				err := adaptor.HandleMessage(ctx, &handler.UnitOfWork{}, cause)
				Expect(err).To(Equal(context.DeadlineExceeded))
			})

			It("releases the cache record when the unit of work is complete", func() {
				handler.NotifyObservers(work, handler.Result{}, nil)

				err := adaptor.HandleMessage(ctx, &handler.UnitOfWork{}, cause)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("retains the cache record if the unit-of-work is successful", func() {
				handler.NotifyObservers(work, handler.Result{}, nil)

				rec, err := adaptor.Cache.Acquire(ctx, "<instance>")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(rec.Instance).NotTo(BeNil())
			})

			It("does not retain the cache record if the unit-of-work is unsuccessful", func() {
				handler.NotifyObservers(work, handler.Result{}, errors.New("<error>"))

				rec, err := adaptor.Cache.Acquire(ctx, "<instance>")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(rec.Instance).To(BeNil())
			})
		})

		When("the instance does not exist", func() {
			When("the instance is created", func() {
				It("causes Create() to return true", func() {
					upstream.HandleCommandFunc = func(
						s dogma.AggregateCommandScope,
						_ dogma.Message,
					) {
						ok := s.Create()
						Expect(ok).To(BeTrue())

						s.RecordEvent(MessageE1)
					}

					err := adaptor.HandleMessage(ctx, work, cause)
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("updates the aggregate meta-data", func() {
					upstream.HandleCommandFunc = func(
						s dogma.AggregateCommandScope,
						_ dogma.Message,
					) {
						s.Create()
						s.RecordEvent(MessageE1)
					}

					err := adaptor.HandleMessage(ctx, work, cause)
					Expect(err).ShouldNot(HaveOccurred())

					Expect(work.Batch).To(Equal(
						persistence.Batch{
							persistence.SaveAggregateMetaData{
								MetaData: aggregatestore.MetaData{
									HandlerKey:     "<aggregate-key>",
									InstanceID:     "<instance>",
									Revision:       0,
									InstanceExists: true,
								},
							},
						},
					))
				})

				It("panics if no event is recorded", func() {
					upstream.HandleCommandFunc = func(
						s dogma.AggregateCommandScope,
						_ dogma.Message,
					) {
						s.Create()
					}

					Expect(func() {
						err := adaptor.HandleMessage(ctx, work, cause)
						Expect(err).ShouldNot(HaveOccurred())
					}).To(PanicWith("*fixtures.AggregateMessageHandler.HandleEvent() created the '<instance>' instance without recording an event while handling a fixtures.MessageC command"))
				})
			})

			It("panics if the root is accessed", func() {
				upstream.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					s.Root()
				}

				Expect(func() {
					err := adaptor.HandleMessage(ctx, work, cause)
					Expect(err).ShouldNot(HaveOccurred())
				}).To(PanicWith("can not access aggregate root of non-existent instance"))
			})

			It("panics if the instance is destroyed", func() {
				upstream.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					s.Destroy()
				}

				Expect(func() {
					err := adaptor.HandleMessage(ctx, work, cause)
					Expect(err).ShouldNot(HaveOccurred())
				}).To(PanicWith("can not destroy non-existent instance"))
			})
		})

		When("the instance already exists", func() {
			BeforeEach(func() {
				upstream.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					s.Create()
					s.RecordEvent(MessageE1)
					s.RecordEvent(MessageE2)
				}

				w := &handler.UnitOfWork{}
				err := adaptor.HandleMessage(ctx, w, cause)
				Expect(err).ShouldNot(HaveOccurred())

				r := &handler.Resolver{
					QueueEvents: message.TypesOf(),
					UnitOfWork:  w,
				}

				pr, err := dataStore.Persist(ctx, r.ResolveBatch())
				Expect(err).ShouldNot(HaveOccurred())

				res := r.ResolveResult(pr)
				handler.NotifyObservers(w, res, nil)

				upstream.HandleCommandFunc = nil
			})

			It("causes Create() to return false", func() {
				upstream.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					ok := s.Create()
					Expect(ok).To(BeFalse())
				}

				err := adaptor.HandleMessage(ctx, work, cause)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("provides a root with the correct state", func() {
				upstream.HandleCommandFunc = func(
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

				err := adaptor.HandleMessage(ctx, work, cause)
				Expect(err).ShouldNot(HaveOccurred())
			})

			When("the instance is destroyed", func() {
				It("resets the root state within the same scope", func() {
					upstream.HandleCommandFunc = func(
						s dogma.AggregateCommandScope,
						_ dogma.Message,
					) {
						s.RecordEvent(MessageE3)
						s.Destroy()

						s.Create()
						r := s.Root().(*AggregateRoot)
						Expect(r.Value).To(Equal(
							&[]dogma.Message{},
						))
					}

					err := adaptor.HandleMessage(ctx, work, cause)
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("causes create to return true again", func() {
					upstream.HandleCommandFunc = func(
						s dogma.AggregateCommandScope,
						_ dogma.Message,
					) {
						s.RecordEvent(MessageE3)
						s.Destroy()

						ok := s.Create()
						Expect(ok).To(BeTrue())
					}

					err := adaptor.HandleMessage(ctx, work, cause)
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("updates the aggregate meta-data", func() {
					upstream.HandleCommandFunc = func(
						s dogma.AggregateCommandScope,
						_ dogma.Message,
					) {
						s.RecordEvent(MessageE3)
						s.Destroy()
					}

					err := adaptor.HandleMessage(ctx, work, cause)
					Expect(err).ShouldNot(HaveOccurred())

					Expect(work.Batch).To(Equal(
						persistence.Batch{
							persistence.SaveAggregateMetaData{
								MetaData: aggregatestore.MetaData{
									HandlerKey:      "<aggregate-key>",
									InstanceID:      "<instance>",
									Revision:        1,
									InstanceExists:  false,
									LastDestroyedBy: "2", // deterministic ID from the packer
								},
							},
						},
					))
				})

				It("panics if no event is recorded", func() {
					upstream.HandleCommandFunc = func(
						s dogma.AggregateCommandScope,
						_ dogma.Message,
					) {
						s.Destroy()
					}

					Expect(func() {
						err := adaptor.HandleMessage(ctx, work, cause)
						Expect(err).ShouldNot(HaveOccurred())
					}).To(PanicWith("*fixtures.AggregateMessageHandler.HandleEvent() destroyed the '<instance>' instance without recording an event while handling a fixtures.MessageC command"))
				})
			})
		})

		When("when the instance has been destroyed", func() {
			BeforeEach(func() {
				upstream.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					s.Create()
					s.RecordEvent(MessageE1)
					s.Destroy()
				}

				w := &handler.UnitOfWork{}
				err := adaptor.HandleMessage(ctx, w, cause)
				Expect(err).ShouldNot(HaveOccurred())

				r := &handler.Resolver{
					QueueEvents: message.TypesOf(),
					UnitOfWork:  w,
				}

				pr, err := dataStore.Persist(ctx, r.ResolveBatch())
				Expect(err).ShouldNot(HaveOccurred())

				res := r.ResolveResult(pr)
				handler.NotifyObservers(w, res, nil)

				upstream.HandleCommandFunc = nil
			})

			It("resets the root state", func() {
				upstream.HandleCommandFunc = func(
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

				err := adaptor.HandleMessage(ctx, work, cause)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})

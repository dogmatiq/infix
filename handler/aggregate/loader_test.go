package aggregate_test

import (
	"context"
	"errors"
	"time"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/dogmatiq/infix/fixtures"
	. "github.com/dogmatiq/infix/handler/aggregate"
	"github.com/dogmatiq/infix/persistence"
	"github.com/dogmatiq/infix/persistence/subsystem/aggregatestore"
	"github.com/dogmatiq/infix/persistence/subsystem/eventstore"
	"github.com/dogmatiq/marshalkit/codec"
	. "github.com/dogmatiq/marshalkit/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Loader", func() {
	var (
		ctx           context.Context
		cancel        context.CancelFunc
		dataStore     *DataStoreStub
		aggregateRepo *AggregateStoreRepositoryStub
		eventRepo     *EventStoreRepositoryStub
		base          *AggregateRoot
		loader        *Loader
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)

		dataStore = NewDataStoreStub()
		aggregateRepo = dataStore.AggregateStoreRepository().(*AggregateStoreRepositoryStub)
		eventRepo = dataStore.EventStoreRepository().(*EventStoreRepositoryStub)

		base = &AggregateRoot{
			Value: &[]dogma.Message{},
			ApplyEventFunc: func(m dogma.Message, v interface{}) {
				p := v.(*[]dogma.Message)
				*p = append(*p, m)
			},
		}

		loader = &Loader{
			AggregateStore: aggregateRepo,
			EventStore:     eventRepo,
			Marshaler:      Marshaler,
		}
	})

	AfterEach(func() {
		if dataStore != nil {
			dataStore.Close()
		}

		cancel()
	})

	Describe("func Load()", func() {
		It("returns an error if the meta-data can not be loaded", func() {
			aggregateRepo.LoadMetaDataFunc = func(
				context.Context,
				string,
				string,
			) (*aggregatestore.MetaData, error) {
				return nil, errors.New("<error>")
			}

			_, err := loader.Load(ctx, "<handler-key>", "<instance>", base)
			Expect(err).To(MatchError("<error>"))
		})

		When("the instance has never existed", func() {
			It("returns an instance with a new meta-data value and the base root", func() {
				inst, err := loader.Load(ctx, "<handler-key>", "<instance>", base)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(inst).To(Equal(
					&Instance{
						MetaData: aggregatestore.MetaData{
							HandlerKey: "<handler-key>",
							InstanceID: "<instance>",
						},
						Root: base,
					},
				))
			})

			It("does not attempt to load events", func() {
				eventRepo.QueryEventsFunc = func(
					context.Context,
					eventstore.Query,
				) (eventstore.Result, error) {
					return nil, errors.New("<error>")
				}

				_, err := loader.Load(ctx, "<handler-key>", "<instance>", base)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("the instance exists", func() {
			BeforeEach(func() {
				_, err := dataStore.Persist(
					ctx,
					persistence.Batch{
						persistence.SaveEvent{
							Envelope: NewEnvelope("<event-0>", MessageE1),
						},
						persistence.SaveEvent{
							Envelope: NewEnvelope("<event-1>", MessageE2),
						},
						persistence.SaveAggregateMetaData{
							MetaData: aggregatestore.MetaData{
								HandlerKey:     "<handler-key>",
								InstanceID:     "<instance>",
								InstanceExists: true,
							},
						},
					},
				)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns an instance with the persisted meta-data and the base root", func() {
				inst, err := loader.Load(ctx, "<handler-key>", "<instance>", base)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(inst).To(Equal(
					&Instance{
						MetaData: aggregatestore.MetaData{
							HandlerKey:     "<handler-key>",
							InstanceID:     "<instance>",
							Revision:       1,
							InstanceExists: true,
						},
						Root: base,
					},
				))
			})

			It("applies historical events to the base root", func() {
				_, err := loader.Load(ctx, "<handler-key>", "<instance>", base)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(base.Value).To(Equal(
					&[]dogma.Message{
						MessageE1,
						MessageE2,
					},
				))
			})

			It("returns an error if the events can not be loaded", func() {
				eventRepo.LoadEventsBySourceFunc = func(
					context.Context,
					string,
					string,
					string,
				) (eventstore.Result, error) {
					return nil, errors.New("<error>")
				}

				_, err := loader.Load(ctx, "<handler-key>", "<instance>", base)
				Expect(err).To(MatchError("<error>"))
			})

			It("returns an error if one of the historical events can not be unmarshaled", func() {
				loader.Marshaler = &codec.Marshaler{} // an empty marshaler cannot unmarshal anything
				_, err := loader.Load(ctx, "<handler-key>", "<instance>", base)
				Expect(err).To(MatchError("no codecs support the 'application/json' media-type"))
			})

			It("does not attempt to load events for a stateless aggregate", func() {
				eventRepo.QueryEventsFunc = func(
					context.Context,
					eventstore.Query,
				) (eventstore.Result, error) {
					return nil, errors.New("<error>")
				}

				_, err := loader.Load(ctx, "<handler-key>", "<instance>", dogma.StatelessAggregateRoot)
				Expect(err).ShouldNot(HaveOccurred())
			})

			When("the instance has been destroyed", func() {
				BeforeEach(func() {
					_, err := dataStore.Persist(
						ctx,
						persistence.Batch{
							persistence.SaveAggregateMetaData{
								MetaData: aggregatestore.MetaData{
									HandlerKey:      "<handler-key>",
									InstanceID:      "<instance>",
									Revision:        1,
									InstanceExists:  false,
									LastDestroyedBy: "<event-1>",
								},
							},
						},
					)
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("does not attempt to load events", func() {
					eventRepo.QueryEventsFunc = func(
						context.Context,
						eventstore.Query,
					) (eventstore.Result, error) {
						return nil, errors.New("<error>")
					}

					_, err := loader.Load(ctx, "<handler-key>", "<instance>", base)
					Expect(err).ShouldNot(HaveOccurred())
				})

				When("the instance is subsequently recreated", func() {
					BeforeEach(func() {
						_, err := dataStore.Persist(
							ctx,
							persistence.Batch{
								persistence.SaveEvent{
									Envelope: NewEnvelope("<event-2>", MessageE3),
								},
								persistence.SaveAggregateMetaData{
									MetaData: aggregatestore.MetaData{
										HandlerKey:      "<handler-key>",
										InstanceID:      "<instance>",
										Revision:        2,
										InstanceExists:  true,
										LastDestroyedBy: "<event-1>",
									},
								},
							},
						)
						Expect(err).ShouldNot(HaveOccurred())
					})

					It("only applies events that were recorded after the destruction", func() {
						_, err := loader.Load(ctx, "<handler-key>", "<instance>", base)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(base.Value).To(Equal(
							&[]dogma.Message{
								MessageE3,
							},
						))
					})
				})
			})
		})
	})
})

package aggregate_test

import (
	"context"
	"errors"
	"github.com/dogmatiq/marshalkit"
	"time"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogma/fixtures"
	"github.com/dogmatiq/marshalkit/codec"
	. "github.com/dogmatiq/marshalkit/fixtures"
	. "github.com/dogmatiq/verity/fixtures"
	. "github.com/dogmatiq/verity/handler/aggregate"
	"github.com/dogmatiq/verity/persistence"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Loader", func() {
	var (
		ctx       context.Context
		cancel    context.CancelFunc
		dataStore *DataStoreStub
		base      *AggregateRoot
		loader    *Loader
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)

		dataStore = NewDataStoreStub()

		base = &AggregateRoot{}

		loader = &Loader{
			AggregateRepository: dataStore,
			EventRepository:     dataStore,
			Marshaler:           Marshaler,
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
			dataStore.LoadAggregateMetaDataFunc = func(
				context.Context,
				string,
				string,
			) (persistence.AggregateMetaData, error) {
				return persistence.AggregateMetaData{}, errors.New("<error>")
			}

			_, err := loader.Load(ctx, "<handler-key>", "<instance>", base)
			Expect(err).To(MatchError("<error>"))
		})

		It("returns an error if the snapshot can not be loaded", func() {
			dataStore.LoadAggregateSnapshotFunc = func(
				context.Context,
				string,
				string,
			) (persistence.AggregateSnapshot, bool, error) {
				return persistence.AggregateSnapshot{}, false, errors.New("<error>")
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
						AggregateMetaData: persistence.AggregateMetaData{
							HandlerKey: "<handler-key>",
							InstanceID: "<instance>",
						},
						Root: base,
					},
				))
			})

			It("does not return a snapshot", func() {
				inst, err := loader.Load(ctx, "<handler-key>", "<instance>", base)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(inst.BarrierEventID).To(BeEmpty()) // is set when snapshot is present
			})

			It("does not attempt to load events", func() {
				dataStore.LoadEventsBySourceFunc = func(
					context.Context,
					string,
					string,
					string,
				) (persistence.EventResult, error) {
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
							MetaData: persistence.AggregateMetaData{
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
						AggregateMetaData: persistence.AggregateMetaData{
							HandlerKey:     "<handler-key>",
							InstanceID:     "<instance>",
							Revision:       1,
							InstanceExists: true,
						},
						Root: base,
					},
				))
			})

			It("returns an instance with the snapshot base root when snapshot is present", func() {
				ssBase := &AggregateRoot{
					AppliedEvents: []dogma.Message{
						MessageE1,
						MessageE2,
						MessageE3,
					},
				}

				var packet marshalkit.Packet
				var err error
				packet, err = Marshaler.Marshal(ssBase)
				Expect(err).ShouldNot(HaveOccurred())

				_, err = dataStore.Persist(
					ctx,
					persistence.Batch{
						persistence.SaveAggregateSnapshot{
							Snapshot: persistence.AggregateSnapshot{
								HandlerKey: "<handler-key>",
								InstanceID: "<instance>",
								Version:    "<event-1>",
								Packet:     packet,
							},
						},
					},
				)
				Expect(err).ShouldNot(HaveOccurred())

				inst, err := loader.Load(ctx, "<handler-key>", "<instance>", base)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(inst).To(Equal(
					&Instance{
						AggregateMetaData: persistence.AggregateMetaData{
							HandlerKey:     "<handler-key>",
							InstanceID:     "<instance>",
							Revision:       1,
							InstanceExists: true,
							BarrierEventID: "<event-1>",
						},
						Root: &AggregateRoot{
							AppliedEvents: []dogma.Message{
								map[string]interface{}{"Value": "E1"},
								map[string]interface{}{"Value": "E2"},
								map[string]interface{}{"Value": "E3"},
							},
						},
					},
				))
			})

			It("applies historical events to the base root", func() {
				_, err := loader.Load(ctx, "<handler-key>", "<instance>", base)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(base.AppliedEvents).To(Equal(
					[]dogma.Message{
						MessageE1,
						MessageE2,
					},
				))
			})

			It("returns an error if the events can not be loaded", func() {
				dataStore.LoadEventsBySourceFunc = func(
					context.Context,
					string,
					string,
					string,
				) (persistence.EventResult, error) {
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

			When("the instance has been destroyed", func() {
				BeforeEach(func() {
					_, err := dataStore.Persist(
						ctx,
						persistence.Batch{
							persistence.SaveAggregateMetaData{
								MetaData: persistence.AggregateMetaData{
									HandlerKey:     "<handler-key>",
									InstanceID:     "<instance>",
									Revision:       1,
									InstanceExists: false,
									LastEventID:    "<event-1>",
									BarrierEventID: "<event-1>",
								},
							},
						},
					)
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("does not attempt to load events", func() {
					dataStore.LoadEventsBySourceFunc = func(
						context.Context,
						string,
						string,
						string,
					) (persistence.EventResult, error) {
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
									MetaData: persistence.AggregateMetaData{
										HandlerKey:     "<handler-key>",
										InstanceID:     "<instance>",
										Revision:       2,
										InstanceExists: true,
										LastEventID:    "<event-2>",
										BarrierEventID: "<event-1>",
									},
								},
							},
						)
						Expect(err).ShouldNot(HaveOccurred())
					})

					It("only applies events that were recorded after the destruction", func() {
						_, err := loader.Load(ctx, "<handler-key>", "<instance>", base)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(base.AppliedEvents).To(Equal(
							[]dogma.Message{
								MessageE3,
							},
						))
					})
				})
			})
		})
	})
})

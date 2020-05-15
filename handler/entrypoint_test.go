package handler_test

import (
	"context"
	"errors"
	"time"

	"github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/dogma/fixtures"
	"github.com/dogmatiq/infix/eventstream"
	. "github.com/dogmatiq/infix/fixtures"
	. "github.com/dogmatiq/infix/handler"
	"github.com/dogmatiq/infix/parcel"
	"github.com/dogmatiq/infix/persistence"
	"github.com/dogmatiq/infix/persistence/subsystem/eventstore"
	"github.com/dogmatiq/infix/persistence/subsystem/queuestore"
	"github.com/dogmatiq/infix/queue"
	. "github.com/jmalloc/gomegax"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type EntryPoint", func() {
	var (
		handler    *HandlerStub
		cause      *parcel.Parcel
		ack        *AcknowledgerStub
		entryPoint *EntryPoint
	)

	BeforeEach(func() {
		handler = &HandlerStub{}

		cause = NewParcel("<consume>", MessageC1)

		ack = &AcknowledgerStub{}

		entryPoint = &EntryPoint{
			QueueEvents: message.TypesOf(MessageQ{}),
			Handler:     handler,
		}
	})

	Describe("func HandleMessage()", func() {
		It("forwards to the handler", func() {
			called := false
			handler.HandleMessageFunc = func(
				_ context.Context,
				_ *UnitOfWork,
				p *parcel.Parcel,
			) error {
				called = true
				Expect(p).To(Equal(cause))
				return nil
			}

			err := entryPoint.HandleMessage(context.Background(), ack, cause)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(called).To(BeTrue())
		})

		When("the message is handled successfully", func() {
			var (
				command        *parcel.Parcel
				timeout        *parcel.Parcel
				unqueuedEvent  *parcel.Parcel
				queuedEvent    *parcel.Parcel
				oper           persistence.Operation
				expectedResult Result
			)

			BeforeEach(func() {
				command = NewParcel("<command>", MessageC1)

				timeout = NewParcel(
					"<timeout>",
					MessageT1,
					time.Now(),
					time.Now().Add(1*time.Hour),
				)

				unqueuedEvent = NewParcel("<unqueued-event>", MessageU1)

				queuedEvent = NewParcel("<queued-event>", MessageQ1)

				oper = persistence.SaveOffset{
					ApplicationKey: "<app-key>",
				}

				handler.HandleMessageFunc = func(
					ctx context.Context,
					w *UnitOfWork,
					p *parcel.Parcel,
				) error {
					w.ExecuteCommand(command)
					w.ScheduleTimeout(timeout)
					w.RecordEvent(unqueuedEvent)
					w.RecordEvent(queuedEvent)
					w.Do(oper)
					return nil
				}

				ack.AckFunc = func(
					context.Context,
					persistence.Batch,
				) (persistence.Result, error) {
					return persistence.Result{
						EventStoreItems: []*eventstore.Item{
							{
								Offset:   0,
								Envelope: unqueuedEvent.Envelope,
							},
							{
								Offset:   1,
								Envelope: queuedEvent.Envelope,
							},
						},
					}, nil
				}

				expectedResult = Result{
					Queued: []queue.Message{
						{
							Parcel: command,
							Item: &queuestore.Item{
								Revision:      1,
								NextAttemptAt: command.CreatedAt,
								Envelope:      command.Envelope,
							},
						},
						{
							Parcel: timeout,
							Item: &queuestore.Item{
								Revision:      1,
								NextAttemptAt: timeout.ScheduledFor,
								Envelope:      timeout.Envelope,
							},
						},
						{
							Parcel: queuedEvent,
							Item: &queuestore.Item{
								Revision:      1,
								NextAttemptAt: queuedEvent.CreatedAt,
								Envelope:      queuedEvent.Envelope,
							},
						},
					},
					Events: []eventstream.Event{
						{
							Offset: 0,
							Parcel: unqueuedEvent,
						},
						{
							Offset: 1,
							Parcel: queuedEvent,
						},
					},
				}
			})

			It("acknowledges the message with the batch from the unit-of-work", func() {
				called := false
				ack.AckFunc = func(
					_ context.Context,
					b persistence.Batch,
				) (persistence.Result, error) {
					called = true
					Expect(b).To(EqualX(
						persistence.Batch{
							// commands ...
							persistence.SaveQueueItem{
								Item: queuestore.Item{
									NextAttemptAt: command.CreatedAt,
									Envelope:      command.Envelope,
								},
							},
							// timeouts ...
							persistence.SaveQueueItem{
								Item: queuestore.Item{
									NextAttemptAt: timeout.ScheduledFor,
									Envelope:      timeout.Envelope,
								},
							},
							// events ...
							persistence.SaveEvent{
								Envelope: unqueuedEvent.Envelope,
							},
							persistence.SaveEvent{
								Envelope: queuedEvent.Envelope,
							},
							persistence.SaveQueueItem{
								Item: queuestore.Item{
									NextAttemptAt: queuedEvent.CreatedAt,
									Envelope:      queuedEvent.Envelope,
								},
							},
							// other ...
							oper,
						},
					))

					return persistence.Result{}, nil
				}

				err := entryPoint.HandleMessage(context.Background(), ack, cause)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(called).To(BeTrue())
			})

			It("notifies the entry point's observers", func() {
				called := false
				entryPoint.Observers = append(
					entryPoint.Observers,
					func(res Result, err error) {
						called = true
						Expect(err).ShouldNot(HaveOccurred())
						Expect(res).To(EqualX(expectedResult))
					},
				)

				err := entryPoint.HandleMessage(context.Background(), ack, cause)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(called).To(BeTrue())
			})

			It("notifies the unit-of-work's observers", func() {
				called := false
				next := handler.HandleMessageFunc
				handler.HandleMessageFunc = func(
					ctx context.Context,
					w *UnitOfWork,
					p *parcel.Parcel,
				) error {
					w.Observe(func(res Result, err error) {
						called = true
						Expect(err).ShouldNot(HaveOccurred())
						Expect(res).To(EqualX(expectedResult))
					})
					return next(ctx, w, p)
				}

				err := entryPoint.HandleMessage(context.Background(), ack, cause)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(called).To(BeTrue())
			})
		})

		When("the handler returns an error", func() {
			BeforeEach(func() {
				handler.HandleMessageFunc = func(
					context.Context,
					*UnitOfWork,
					*parcel.Parcel,
				) error {
					return errors.New("<error>")
				}
			})

			It("notifies the entry point's observers", func() {
				called := false
				entryPoint.Observers = append(
					entryPoint.Observers,
					func(_ Result, err error) {
						called = true
						Expect(err).To(MatchError("<error>"))
					},
				)

				err := entryPoint.HandleMessage(context.Background(), ack, cause)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(called).To(BeTrue())
			})

			It("notifies the unit-of-work's observers", func() {
				called := false
				handler.HandleMessageFunc = func(
					_ context.Context,
					w *UnitOfWork,
					_ *parcel.Parcel,
				) error {
					w.Observe(func(_ Result, err error) {
						called = true
						Expect(err).To(MatchError("<error>"))
					})

					return errors.New("<error>")
				}

				err := entryPoint.HandleMessage(context.Background(), ack, cause)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(called).To(BeTrue())
			})

			It("negatively acknowledges the message", func() {
				called := false
				ack.NackFunc = func(
					_ context.Context,
					err error,
				) error {
					called = true
					Expect(err).To(MatchError("<error>"))
					return nil
				}

				err := entryPoint.HandleMessage(context.Background(), ack, cause)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(called).To(BeTrue())
			})

			It("does not acknowledge the message", func() {
				ack.AckFunc = func(
					context.Context,
					persistence.Batch,
				) (persistence.Result, error) {
					Fail("unexpected call")
					return persistence.Result{}, nil
				}

				err := entryPoint.HandleMessage(context.Background(), ack, cause)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns an error if Nack() fails", func() {
				ack.NackFunc = func(
					context.Context,
					error,
				) error {
					return errors.New("<nack error>")
				}

				err := entryPoint.HandleMessage(context.Background(), ack, cause)
				Expect(err).To(MatchError("<nack error>"))
			})
		})

		When("the acknowledger returns an error", func() {
			BeforeEach(func() {
				ack.AckFunc = func(
					context.Context,
					persistence.Batch,
				) (persistence.Result, error) {
					return persistence.Result{}, errors.New("<error>")
				}
			})

			It("notifies the entry point's observers", func() {
				called := false
				entryPoint.Observers = append(
					entryPoint.Observers,
					func(_ Result, err error) {
						called = true
						Expect(err).To(MatchError("<error>"))
					},
				)

				err := entryPoint.HandleMessage(context.Background(), ack, cause)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(called).To(BeTrue())
			})

			It("notifies the unit-of-work's observers", func() {
				called := false
				handler.HandleMessageFunc = func(
					_ context.Context,
					w *UnitOfWork,
					_ *parcel.Parcel,
				) error {
					w.Observe(func(_ Result, err error) {
						called = true
						Expect(err).To(MatchError("<error>"))
					})

					return errors.New("<error>")
				}

				err := entryPoint.HandleMessage(context.Background(), ack, cause)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(called).To(BeTrue())
			})

			It("negatively acknowledges the message", func() {
				called := false
				ack.NackFunc = func(
					_ context.Context,
					err error,
				) error {
					called = true
					Expect(err).To(MatchError("<error>"))
					return nil
				}

				err := entryPoint.HandleMessage(context.Background(), ack, cause)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(called).To(BeTrue())
			})
		})
	})
})

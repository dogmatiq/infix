package pipeline_test

import (
	"context"
	"time"

	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/dogmatiq/infix/fixtures"
	"github.com/dogmatiq/infix/parcel"
	"github.com/dogmatiq/infix/persistence"
	"github.com/dogmatiq/infix/persistence/subsystem/queuestore"
	"github.com/dogmatiq/infix/pipeline"
	. "github.com/dogmatiq/infix/pipeline"
	"github.com/dogmatiq/infix/queue"
	. "github.com/dogmatiq/marshalkit/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/sync/semaphore"
)

var _ = Describe("type QueueSource", func() {
	var (
		ctx       context.Context
		cancel    context.CancelFunc
		dataStore persistence.DataStore
		source    *QueueSource
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)

		dataStore = NewDataStoreStub()

		source = &QueueSource{
			Queue: &queue.Queue{
				DataStore: dataStore,
				Marshaler: Marshaler,
			},
			Semaphore: semaphore.NewWeighted(1),
		}

		p := NewParcel("<id>", MessageA1)
		i := &queuestore.Item{
			NextAttemptAt: time.Now(),
			Envelope:      p.Envelope,
		}

		_, err := persistence.WithTransaction(
			ctx,
			dataStore,
			func(tx persistence.ManagedTransaction) error {
				return tx.SaveMessageToQueue(ctx, i)
			},
		)
		Expect(err).ShouldNot(HaveOccurred())

		i.Revision++

		err = source.Queue.Track(ctx, p, i)
		Expect(err).ShouldNot(HaveOccurred())
	})

	JustBeforeEach(func() {
		go source.Queue.Run(ctx)
	})

	AfterEach(func() {
		if dataStore != nil {
			dataStore.Close()
		}

		cancel()
	})

	Describe("func Run()", func() {
		It("returns an error if the context is canceled", func() {
			cancel()

			err := source.Run(ctx)
			Expect(err).To(Equal(context.Canceled))
		})

		It("passes the request to the pipeline", func() {
			source.Pipeline = func(
				ctx context.Context,
				req pipeline.Request,
			) error {
				defer GinkgoRecover()
				defer cancel()
				Expect(req.Envelope().MetaData.MessageId).To(Equal("<id>"))
				return nil
			}

			err := source.Run(ctx)
			Expect(err).To(Equal(context.Canceled))
		})

		It("returns an error if the context is canceled while waiting for the sempahore", func() {
			err := source.Semaphore.Acquire(ctx, 1)
			Expect(err).ShouldNot(HaveOccurred())
			defer source.Semaphore.Release(1)

			go func() {
				time.Sleep(100 * time.Millisecond)
				cancel()
			}()

			err = source.Run(ctx)
			Expect(err).To(Equal(context.Canceled))
		})
	})
})

var _ = Describe("func TrackWithQueue()", func() {
	var (
		ctx       context.Context
		cancel    context.CancelFunc
		dataStore persistence.DataStore
		mqueue    *queue.Queue
		observer  pipeline.QueueObserver
		pcl       *parcel.Parcel
		item      *queuestore.Item
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)

		pcl = NewParcel("<id>", MessageA1)
		dataStore = NewDataStoreStub()

		mqueue = &queue.Queue{
			DataStore: dataStore,
			Marshaler: Marshaler,
		}

		observer = TrackWithQueue(mqueue)

		item = &queuestore.Item{
			NextAttemptAt: time.Now(),
			Envelope:      pcl.Envelope,
		}

		_, err := persistence.WithTransaction(
			ctx,
			dataStore,
			func(tx persistence.ManagedTransaction) error {
				return tx.SaveMessageToQueue(ctx, item)
			},
		)
		Expect(err).ShouldNot(HaveOccurred())
		item.Revision++
	})

	AfterEach(func() {
		if dataStore != nil {
			dataStore.Close()
		}

		cancel()
	})

	It("tracks messages when they are enqueued", func() {
		err := observer(
			ctx,
			[]*parcel.Parcel{pcl},
			[]*queuestore.Item{item},
		)
		Expect(err).ShouldNot(HaveOccurred())

		go mqueue.Run(ctx)
		req, err := mqueue.Pop(ctx)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(req.Envelope().MetaData.MessageId).To(Equal("<id>"))
		req.Close()
	})

	It("returns an error if the context deadline is exceeded", func() {
		// It's an implementation detail, but the internal channel used to start
		// tracking is buffered at the same size as the overall buffer size
		// limit.
		//
		// We can't set it to zero, because that will fallback to the default.
		// We also can't start the queue, otherwise it'll start reading from
		// this channel and nothing will block.
		//
		// Instead, we set it to one, and "fill" the channel with a request to
		// ensure that it will block.
		mqueue.BufferSize = 1
		err := mqueue.Track(ctx, pcl, item)
		Expect(err).ShouldNot(HaveOccurred())

		// Setup a short deadline for the test.
		ctx, cancel := context.WithTimeout(ctx, 5*time.Millisecond)
		defer cancel()

		err = observer(
			ctx,
			[]*parcel.Parcel{pcl},
			[]*queuestore.Item{item},
		)
		Expect(err).To(Equal(context.DeadlineExceeded))
	})
})

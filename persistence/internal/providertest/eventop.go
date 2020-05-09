package providertest

import (
	"sort"
	"sync"

	dogmafixtures "github.com/dogmatiq/dogma/fixtures"
	"github.com/dogmatiq/infix/draftspecs/envelopespec"
	infixfixtures "github.com/dogmatiq/infix/fixtures"
	"github.com/dogmatiq/infix/persistence"
	"github.com/dogmatiq/infix/persistence/internal/providertest/common"
	"github.com/dogmatiq/infix/persistence/subsystem/eventstore"
	"github.com/jmalloc/gomegax"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

// declareEventOperationTests declares a functional test-suite for
// persistence operations related to events.
func declareEventOperationTests(tc *common.TestContext) {
	ginkgo.Context("event operations", func() {
		var (
			dataStore  persistence.DataStore
			repository eventstore.Repository
			tearDown   func()

			env0, env1, env2 *envelopespec.Envelope
		)

		ginkgo.BeforeEach(func() {
			dataStore, tearDown = tc.SetupDataStore()
			repository = dataStore.EventStoreRepository()

			env0 = infixfixtures.NewEnvelope("<message-0>", dogmafixtures.MessageA1)
			env1 = infixfixtures.NewEnvelope("<message-1>", dogmafixtures.MessageB1)
			env2 = infixfixtures.NewEnvelope("<message-2>", dogmafixtures.MessageC1)
		})

		ginkgo.AfterEach(func() {
			tearDown()
		})

		ginkgo.Describe("type persistence.SaveEvent", func() {
			ginkgo.It("saves the event to the store", func() {
				persist(
					tc.Context,
					dataStore,
					persistence.SaveEvent{
						Envelope: env0,
					},
					persistence.SaveEvent{
						Envelope: env1,
					},
				)

				items := queryEvents(tc.Context, repository, eventstore.Query{})
				gomega.Expect(items).To(gomegax.EqualX(
					[]*eventstore.Item{
						{
							Offset:   0,
							Envelope: env0,
						},
						{
							Offset:   1,
							Envelope: env1,
						},
					},
				))
			})

			ginkgo.It("has a corresponding item in the batch result", func() {
				res := persist(
					tc.Context,
					dataStore,
					persistence.SaveEvent{
						Envelope: env0,
					},
					persistence.SaveEvent{
						Envelope: env1,
					},
				)

				gomega.Expect(res.EventStoreItems).To(
					gomega.ConsistOf(
						gomegax.EqualX(
							&eventstore.Item{
								Offset:   0,
								Envelope: env0,
							},
						),
						gomegax.EqualX(
							&eventstore.Item{
								Offset:   1,
								Envelope: env1,
							},
						),
					),
				)
			})

			ginkgo.It("serializes operations from concurrent persist calls", func() {
				var (
					g      sync.WaitGroup
					m      sync.Mutex
					expect []*eventstore.Item
				)

				fn := func(env *envelopespec.Envelope) {
					defer ginkgo.GinkgoRecover()
					defer g.Done()

					res := persist(
						tc.Context,
						dataStore,
						persistence.SaveEvent{
							Envelope: env,
						},
					)

					m.Lock()
					expect = append(expect, res.EventStoreItems...)
					m.Unlock()
				}

				g.Add(3)
				go fn(env0)
				go fn(env1)
				go fn(env2)
				g.Wait()

				// Sort the expected slice, as the appends in each goroutine
				// could be out-of-sync with the saves.
				sort.Slice(
					expect,
					func(i, j int) bool {
						return expect[i].Offset < expect[j].Offset
					},
				)

				items := queryEvents(tc.Context, repository, eventstore.Query{})
				gomega.Expect(items).To(gomegax.EqualX(expect))
			})
		})
	})
}

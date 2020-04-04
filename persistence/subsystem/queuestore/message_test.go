package queuestore_test

import (
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/dogmatiq/infix/fixtures"
	. "github.com/dogmatiq/infix/persistence/subsystem/queuestore"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Message", func() {
	Describe("func ID()", func() {
		It("returns the ID from the envelope", func() {
			m := &Message{
				Envelope: NewEnvelopeProto("<id>", MessageA1),
			}

			Expect(m.ID()).To(Equal("<id>"))
		})
	})
})

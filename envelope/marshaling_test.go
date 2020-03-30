package envelope_test

import (
	"time"

	"github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/dogma/fixtures"
	"github.com/dogmatiq/infix/draftspecs/envelopespec"
	"github.com/dogmatiq/infix/envelope"
	. "github.com/dogmatiq/infix/envelope"
	. "github.com/dogmatiq/infix/fixtures"
	. "github.com/dogmatiq/marshalkit/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Marshal()", func() {
	var in *envelope.Envelope

	BeforeEach(func() {
		in = NewEnvelope(
			"<id>",
			MessageA1,
			time.Now(),
			time.Now().Add(1*time.Hour),
		)
	})

	It("marshals to protobuf", func() {
		out, err := Marshal(in)
		Expect(err).ShouldNot(HaveOccurred())

		createdAt := in.CreatedAt.Format(time.RFC3339Nano)
		scheduledFor := in.ScheduledFor.Format(time.RFC3339Nano)

		Expect(out).To(Equal(&envelopespec.Envelope{
			MetaData: &envelopespec.MetaData{
				MessageId:     "<id>",
				CausationId:   "<cause>",
				CorrelationId: "<correlation>",
				Source: &envelopespec.Source{
					Application: &envelopespec.Identity{Name: "<app-name>", Key: "<app-key>"},
					Handler:     &envelopespec.Identity{Name: "<handler-name>", Key: "<handler-key>"},
					InstanceId:  "<instance>",
				},
				CreatedAt:    createdAt,
				ScheduledFor: scheduledFor,
			},
			PortableName: MessageAPortableName,
			MediaType:    MessageA1Packet.MediaType,
			Data:         MessageA1Packet.Data,
		}))
	})

	It("marshals a zero scheduled-for time as an empty string", func() {
		in.MetaData.ScheduledFor = time.Time{}

		out, err := Marshal(in)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out.MetaData.ScheduledFor).To(BeEmpty())
	})

	It("returns an error if the media-type cannot be parsed", func() {
		in.Packet.MediaType = "<malformed>"

		_, err := Marshal(in)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func MarshalMany()", func() {
	It("marshals multiple envelopes", func() {
		in1 := NewEnvelope("<id-2>", MessageA1)
		in2 := NewEnvelope("<id-2>", MessageA2)

		out, err := MarshalMany(
			[]*Envelope{
				in1,
				in2,
			},
		)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(out).To(Equal(
			[]*envelopespec.Envelope{
				MustMarshal(in1),
				MustMarshal(in2),
			},
		))
	})

	It("returns an error if marshaling fails", func() {
		in := NewEnvelope("<id>", MessageA1)
		in.Packet.MediaType = "<malformed>"

		_, err := MarshalMany(
			[]*Envelope{in},
		)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("func MustMarshal()", func() {
	It("marshals the envelope", func() {
		in := NewEnvelope("<id>", MessageA1)

		expect, err := Marshal(in)
		Expect(err).ShouldNot(HaveOccurred())

		out := MustMarshal(in)
		Expect(out).To(Equal(expect))
	})

	It("panics if marshaling fails", func() {
		in := NewEnvelope("<id>", MessageA1)
		in.Packet.MediaType = "<malformed>"

		Expect(func() {
			MustMarshal(in)
		}).To(Panic())
	})
})

var _ = Describe("func MustMarshalMany()", func() {
	It("marshals multiple envelopes", func() {
		in1 := NewEnvelope("<id-2>", MessageA1)
		in2 := NewEnvelope("<id-2>", MessageA2)
		out := MustMarshalMany(
			[]*Envelope{
				in1,
				in2,
			},
		)

		Expect(out).To(Equal(
			[]*envelopespec.Envelope{
				MustMarshal(in1),
				MustMarshal(in2),
			},
		))
	})

	It("panics if marshaling fails", func() {
		in := NewEnvelope("<id>", MessageA1)
		in.Packet.MediaType = "<malformed>"

		Expect(func() {
			MustMarshalMany(
				[]*Envelope{in},
			)
		}).To(Panic())
	})
})

var _ = Describe("func Unmarshal()", func() {
	var (
		createdAt, scheduledFor time.Time
		in                      *envelopespec.Envelope
	)

	BeforeEach(func() {
		createdAt = time.Now()
		scheduledFor = time.Now().Add(1 * time.Hour)

		in = MustMarshal(NewEnvelope(
			"<id>",
			MessageA1,
			createdAt,
			scheduledFor,
		))
	})

	It("unmarshals from protobuf", func() {
		var out envelope.Envelope
		err := Unmarshal(Marshaler, in, &out)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(out.CreatedAt).To(BeTemporally("==", createdAt))
		Expect(out.ScheduledFor).To(BeTemporally("==", scheduledFor))

		// clear values for comparison below
		out.CreatedAt = time.Time{}
		out.ScheduledFor = time.Time{}

		Expect(out).To(Equal(envelope.Envelope{
			MetaData: envelope.MetaData{
				MessageID:     "<id>",
				CausationID:   "<cause>",
				CorrelationID: "<correlation>",
				Source: envelope.Source{
					Application: configkit.MustNewIdentity("<app-name>", "<app-key>"),
					Handler:     configkit.MustNewIdentity("<handler-name>", "<handler-key>"),
					InstanceID:  "<instance>",
				},
			},
			Message: MessageA1,
			Packet:  MessageA1Packet,
		}))
	})

	It("returns an error if the created-at time can not be unmarshaled", func() {
		in.MetaData.CreatedAt = "not-a-valid-time"

		var out envelope.Envelope
		err := Unmarshal(Marshaler, in, &out)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the scheduled-for time can not be unmarshaled", func() {
		in.MetaData.ScheduledFor = "not-a-valid-time"

		var out envelope.Envelope
		err := Unmarshal(Marshaler, in, &out)
		Expect(err).Should(HaveOccurred())
	})

	It("does not return an error if the scheduled-for time is empty", func() {
		in.MetaData.ScheduledFor = ""

		var out envelope.Envelope
		err := Unmarshal(Marshaler, in, &out)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out.MetaData.ScheduledFor.IsZero()).To(BeTrue())
	})

	It("returns an error if the source is not valid", func() {
		in.MetaData.Source.Handler = nil

		var out envelope.Envelope
		err := Unmarshal(Marshaler, in, &out)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the meta-data is not valid", func() {
		in.MetaData.CreatedAt = ""

		var out envelope.Envelope
		err := Unmarshal(Marshaler, in, &out)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the marshaler fails", func() {
		in.MediaType = "<unknown>"

		var out envelope.Envelope
		err := Unmarshal(Marshaler, in, &out)
		Expect(err).Should(HaveOccurred())
	})
})

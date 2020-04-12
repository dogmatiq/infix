package envelope

import (
	"fmt"
	"time"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
	"github.com/google/uuid"
)

// Packer puts messages into envelopes.
type Packer struct {
	// Application is the identity of this application.
	Application configkit.Identity

	// Roles is a map of message type to role, used to validate the messages
	// that are being packed.
	Roles message.TypeRoles

	// GenerateID is a function used to generate new message IDs. If it is nil,
	// a UUID is generated.
	GenerateID func() string

	// Now is a function used to get the current time. If it is nil, time.Now()
	// is used.
	Now func() time.Time
}

// PackCommand returns a new command envelope containing the given message.
func (p *Packer) PackCommand(m dogma.Message) *Envelope {
	p.checkRole(m, message.CommandRole)
	return p.new(m)
}

// PackEvent returns a new event envelope containing the given message.
func (p *Packer) PackEvent(m dogma.Message) *Envelope {
	p.checkRole(m, message.EventRole)
	return p.new(m)
}

// PackChildCommand returns a new command envelope containing the given message
// and configured as a child of cause.
func (p *Packer) PackChildCommand(
	cause *Envelope,
	m dogma.Message,
	handler configkit.Identity,
	instanceID string,
) *Envelope {
	p.checkRole(cause.Message, message.EventRole, message.TimeoutRole)
	p.checkRole(m, message.CommandRole)

	return p.newChild(
		cause,
		m,
		handler,
		instanceID,
	)
}

// PackChildEvent returns a new event envelope containing the given message and
// configured as a child of cause.
func (p *Packer) PackChildEvent(
	cause *Envelope,
	m dogma.Message,
	handler configkit.Identity,
	instanceID string,
) *Envelope {
	p.checkRole(cause.Message, message.CommandRole)
	p.checkRole(m, message.EventRole)

	return p.newChild(
		cause,
		m,
		handler,
		instanceID,
	)
}

// PackChildTimeout returns a new timeout envelope containing the given message
// and configured as a child of cause.
func (p *Packer) PackChildTimeout(
	cause *Envelope,
	m dogma.Message,
	t time.Time,
	handler configkit.Identity,
	instanceID string,
) *Envelope {
	p.checkRole(cause.Message, message.EventRole, message.TimeoutRole)
	p.checkRole(m, message.TimeoutRole)

	env := p.newChild(
		cause,
		m,
		handler,
		instanceID,
	)

	env.ScheduledFor = t

	return env
}

// new returns an envelope containing the given message.
func (p *Packer) new(m dogma.Message) *Envelope {
	id := p.generateID()

	return &Envelope{
		MetaData{
			id,
			id,
			id,
			Source{
				Application: p.Application,
			},
			p.now(),
			time.Time{},
		},
		m,
	}
}

// newChild returns an envelope containing the given message, which was a
// produced as a result of handling a causal message.
func (p *Packer) newChild(
	cause *Envelope,
	m dogma.Message,
	handler configkit.Identity,
	instanceID string,
) *Envelope {
	env := p.new(m)

	env.CausationID = cause.MessageID
	env.CorrelationID = cause.CorrelationID
	env.Source.Handler = handler
	env.Source.InstanceID = instanceID

	return env
}

// now returns the current time.
func (p *Packer) now() time.Time {
	now := p.Now
	if now == nil {
		now = time.Now
	}

	return now()
}

// generateID generates a new message ID.
func (p *Packer) generateID() string {
	if p.GenerateID != nil {
		return p.GenerateID()
	}

	return uuid.New().String()
}

// checkRolepatnics if mt does not fill one of the the given roles.
func (p *Packer) checkRole(m dogma.Message, roles ...message.Role) {
	mt := message.TypeOf(m)
	x, ok := p.Roles[mt]

	if !ok {
		panic(fmt.Sprintf("%s is not a recognised message type", mt))
	}

	x.MustBe(roles...)
}

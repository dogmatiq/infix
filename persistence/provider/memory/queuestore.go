package memory

import (
	"context"
	"sort"

	"github.com/dogmatiq/infix/persistence/subsystem/queuestore"
)

// SaveMessageToQueue persists a message to the application's message queue.
//
// If the message is already on the queue its meta-data is updated.
//
// m.Revision must be the revision of the message as currently persisted,
// otherwise an optimistic concurrency conflict has occurred, the message
// is not saved and ErrConflict is returned.
func (t *transaction) SaveMessageToQueue(
	ctx context.Context,
	m *queuestore.Message,
) error {
	if err := t.begin(ctx); err != nil {
		return err
	}

	id := m.Envelope.GetMetaData().GetMessageId()
	committed := t.ds.db.queue.uniq[id]
	uncommitted, changed := t.uncommitted.queue[id]

	effective := committed

	if changed {
		effective = uncommitted
	}

	var rev queuestore.Revision
	if effective != nil {
		rev = effective.Revision
	}

	if m.Revision != rev {
		return queuestore.ErrConflict
	}

	if uncommitted == nil {
		uncommitted = cloneQueueMessage(m)
	} else {
		uncommitted.NextAttemptAt = m.NextAttemptAt
	}

	uncommitted.Revision++

	if t.uncommitted.queue == nil {
		t.uncommitted.queue = map[string]*queuestore.Message{}
	}
	t.uncommitted.queue[id] = uncommitted

	return nil
}

// RemoveMessageFromQueue removes a specific message from the application's
// message queue.
//
// m.Revision must be the revision of the message as currently persisted,
// otherwise an optimistic concurrency conflict has occurred, the message
// remains on the queue and ErrConflict is returned.
func (t *transaction) RemoveMessageFromQueue(
	ctx context.Context,
	m *queuestore.Message,
) error {
	if err := t.begin(ctx); err != nil {
		return err
	}

	id := m.Envelope.GetMetaData().GetMessageId()
	committed, exists := t.ds.db.queue.uniq[id]
	uncommitted, changed := t.uncommitted.queue[id]

	effective := committed
	if changed {
		effective = uncommitted
	}

	if effective == nil || m.Revision != effective.Revision {
		return queuestore.ErrConflict
	}

	if !exists {
		// The message has been saved in this transaction, is now being deleted,
		// and never existed before this transaction.
		delete(t.uncommitted.queue, id)
		return nil
	}

	if t.uncommitted.queue == nil {
		t.uncommitted.queue = map[string]*queuestore.Message{}
	}
	t.uncommitted.queue[id] = nil

	return nil
}

// commitQueue commits staged queue items to the database.
func (t *transaction) commitQueue() {
	q := &t.ds.db.queue

	needsSort := false
	hasDeletes := false

	for id, uncommitted := range t.uncommitted.queue {
		committed, exists := q.uniq[id]

		if !exists {
			t.commitQueueInsert(id, uncommitted)
		} else if uncommitted != nil {
			if t.commitQueueUpdate(id, committed, uncommitted) {
				needsSort = true
			}
		} else {
			committed.Revision = 0 // mark for deletion
			needsSort = true
			hasDeletes = true
		}
	}

	if needsSort {
		sort.Slice(
			q.order,
			func(i, j int) bool {
				a := q.order[i]
				b := q.order[j]

				// If a message has a revision of 0 that means it's being
				// deleted, always sort it AFTER anything that's being kept.
				if a.Revision == 0 || b.Revision == 0 {
					return a.Revision > b.Revision
				}

				return a.NextAttemptAt.Before(b.NextAttemptAt)
			},
		)

		if hasDeletes {
			t.commitQueueRemovals()
		}
	}
}

// commitQueueInsert inserts a new message into the queue as part of a commit.
func (t *transaction) commitQueueInsert(
	id string,
	uncommitted *queuestore.Message,
) {
	q := &t.ds.db.queue

	// Add the message to the unique index.
	if q.uniq == nil {
		q.uniq = map[string]*queuestore.Message{}
	}
	q.uniq[id] = uncommitted

	// Find the index where we'll insert our event. It's the index of the
	// first message that has a NextAttemptedAt greater than m's.
	index := sort.Search(
		len(q.order),
		func(i int) bool {
			return uncommitted.NextAttemptAt.Before(
				q.order[i].NextAttemptAt,
			)
		},
	)

	// Expand the size of the queue.
	q.order = append(q.order, nil)

	// Shift messages further back to make space for m.
	copy(q.order[index+1:], q.order[index:])

	// Insert m at the index.
	q.order[index] = uncommitted
}

// commitQueueUpdate updates an existing message as part of a commit.
func (t *transaction) commitQueueUpdate(
	id string,
	committed *queuestore.Message,
	uncommitted *queuestore.Message,
) bool {
	committed.Revision = uncommitted.Revision

	if uncommitted.NextAttemptAt.Equal(committed.NextAttemptAt) {
		// Bail early and return false if the next-attempt time has not change
		// to avoid sorting the queue unnecessarily.
		return false
	}

	committed.NextAttemptAt = uncommitted.NextAttemptAt

	return true
}

// commitQueueRemovals removes existing messages as part of a commit.
func (t *transaction) commitQueueRemovals() {
	q := &t.ds.db.queue

	i := len(q.order) - 1

	for i >= 0 {
		m := q.order[i]

		if m.Revision > 0 {
			// All deleted messages have a revision of 0 and have been sorted to
			// the end, so once we find a real revision, we're done.
			break
		}

		q.order[i] = nil // prevent memory leak
		delete(q.uniq, m.ID())

		i--
	}

	// Trim the slice to the new length.
	q.order = q.order[:i+1]
}

// queueStoreRepository is an implementation of queuestore.Repository that
// stores queued messages in memory.
type queueStoreRepository struct {
	db *database
}

// LoadQueueMessages loads the next n messages from the queue.
func (r *queueStoreRepository) LoadQueueMessages(
	ctx context.Context,
	n int,
) ([]*queuestore.Message, error) {
	if err := r.db.RLock(ctx); err != nil {
		return nil, err
	}
	defer r.db.RUnlock()

	max := len(r.db.queue.order)
	if n > max {
		n = max
	}

	result := make([]*queuestore.Message, n)

	for i, m := range r.db.queue.order[:n] {
		result[i] = cloneQueueMessage(m)
	}

	return result, nil
}

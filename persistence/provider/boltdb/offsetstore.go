package boltdb

import (
	"context"

	"github.com/dogmatiq/infix/eventstream"
	"github.com/dogmatiq/infix/internal/x/bboltx"
	"github.com/dogmatiq/infix/persistence/subsystem/offsetstore"
	"go.etcd.io/bbolt"
)

var (
	// offsetStoreBucketKey is the key for the bucket at the root of the
	// offsetstore.
	offsetStoreBucketKey = []byte("offsetstore")
)

// SaveOffset persists the "next" offset to be consumed for a specific
// application.
func (t *transaction) SaveOffset(
	ctx context.Context,
	ak string,
	c, n eventstream.Offset,
) (err error) {
	defer bboltx.Recover(&err)

	if err := t.begin(ctx); err != nil {
		return err
	}

	bak := []byte(ak)
	store := bboltx.CreateBucketIfNotExists(
		t.actual,
		offsetStoreBucketKey,
	)

	o := unmarshalOffsetStoreOffset(
		store.Get(bak),
	)

	if c != o {
		return offsetstore.ErrConflict
	}

	bboltx.Put(store, bak, marshalOffsetStoreOffset(n))

	return nil
}

// offsetStoreRepository is an implementation of offsetstore.Repository that
// stores the event stream offset associated with a specific application in a
// BoltDB database.
type offsetStoreRepository struct {
	db *database
}

// LoadOffset loads the offset associated with a specific application.
func (r *offsetStoreRepository) LoadOffset(
	ctx context.Context,
	ak string,
) (o eventstream.Offset, err error) {
	defer bboltx.Recover(&err)

	r.db.View(
		ctx,
		func(tx *bbolt.Tx) {
			if store, exists := bboltx.TryBucket(
				tx,
				offsetStoreBucketKey,
			); exists {
				bak := []byte(ak)

				o = unmarshalOffsetStoreOffset(
					store.Get(bak),
				)
			}
		},
	)

	return
}

// marshalOffsetStoreOffset marshals an event stream offset to its binary
// representation.
func marshalOffsetStoreOffset(offset eventstream.Offset) []byte {
	return marshalUint64(uint64(offset))
}

// unmarshalOffsetStoreOffset unmarshals an event stream offset from its binary
// representation.
func unmarshalOffsetStoreOffset(data []byte) eventstream.Offset {
	return eventstream.Offset(unmarshalUint64(data))
}

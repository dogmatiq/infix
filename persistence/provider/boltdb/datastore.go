package boltdb

import (
	"context"
	"sync"

	"github.com/dogmatiq/infix/internal/x/bboltx"
	"github.com/dogmatiq/infix/persistence"
	"github.com/dogmatiq/infix/persistence/subsystem/aggregatestore"
	"github.com/dogmatiq/infix/persistence/subsystem/eventstore"
	"github.com/dogmatiq/infix/persistence/subsystem/offsetstore"
	"github.com/dogmatiq/infix/persistence/subsystem/queuestore"
	"go.etcd.io/bbolt"
)

// dataStore is an implementation of persistence.DataStore for BoltDB.
type dataStore struct {
	db     *bbolt.DB
	appKey []byte

	m       sync.RWMutex
	release func(string) error
}

// AggregateStoreRepository returns application's aggregate store repository.
func (ds *dataStore) AggregateStoreRepository() aggregatestore.Repository {
	return ds
}

// EventStoreRepository returns the application's event store repository.
func (ds *dataStore) EventStoreRepository() eventstore.Repository {
	return ds
}

// OffsetStoreRepository returns the application's event store repository.
func (ds *dataStore) OffsetStoreRepository() offsetstore.Repository {
	return ds
}

// QueueStoreRepository returns the application's queue store repository.
func (ds *dataStore) QueueStoreRepository() queuestore.Repository {
	return ds
}

// Persist commits a batch of operations atomically.
//
// If any one of the operations causes an optimistic concurrency conflict
// the entire batch is aborted and a ConflictError is returned.
func (ds *dataStore) Persist(
	ctx context.Context,
	b persistence.Batch,
) (_ persistence.Result, err error) {
	b.MustValidate()

	defer bboltx.Recover(&err)

	ds.m.RLock()
	defer ds.m.RUnlock()

	if ds.release == nil {
		return persistence.Result{}, persistence.ErrDataStoreClosed
	}

	c := &committer{}

	bboltx.Update(
		ds.db,
		func(tx *bbolt.Tx) {
			c.root = bboltx.CreateBucketIfNotExists(tx, ds.appKey)
			bboltx.Must(b.AcceptVisitor(ctx, c))
		},
	)

	return c.result, nil
}

// Close closes the data store.
//
// Closing a data-store causes any future calls to Persist() to return
// ErrDataStoreClosed.
//
// The behavior read operations on a closed data-store is
// implementation-defined.
//
// In general use it is expected that all pending calls to Persist() will
// have finished before a data-store is closed. Close() may block until any
// in-flight calls to Persist() return, or may prevent any such calls from
// succeeding.
func (ds *dataStore) Close() error {
	ds.m.Lock()
	defer ds.m.Unlock()

	if ds.release == nil {
		return persistence.ErrDataStoreClosed
	}

	r := ds.release
	ds.db = nil
	ds.release = nil

	return r(string(ds.appKey))
}

// committer is an implementation of persitence.OperationVisitor that
// applies operations to the database.
//
// It is expected that the operations have already been validated using
// validator.
type committer struct {
	root   *bbolt.Bucket
	result persistence.Result
}

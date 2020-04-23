package offsetstore

import (
	"context"
	"errors"

	"github.com/dogmatiq/infix/eventstream"
)

// ErrConflict is returned by transaction operations when a persisted
// offset can not be updated because the supplied "current" offset is
// out of date.
var ErrConflict = errors.New("an optimistic concurrency conflict occured while persisting to the offset store")

// Transaction defines the primitive persistence operations for manipulating
// the application event stream offset.
type Transaction interface {
	// SaveOffset persists the "next" offset to be consumed for a specific
	// application.
	//
	// ak is the application's identity key.
	//
	// c must be the offset currently associated with ak, otherwise an optimistic
	// concurrency conflict has occurred, the offset is not saved and ErrConflict
	// is returned.
	SaveOffset(
		ctx context.Context,
		ak string,
		c, n eventstream.Offset,
	) error
}

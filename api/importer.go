package api

import (
	"context"
	"time"

	"github.com/gernest/roaring"
	"github.com/gernest/sql3/dax"
)

type Importer interface {
	StartTransaction(ctx context.Context, id string, timeout time.Duration, exclusive bool, requestTimeout time.Duration) (*Transaction, error)
	FinishTransaction(ctx context.Context, id string) (*Transaction, error)
	CreateTableKeys(ctx context.Context, tid dax.TableID, keys ...string) (map[string]uint64, error)
	CreateFieldKeys(ctx context.Context, tid dax.TableID, fname dax.FieldName, keys ...string) (map[string]uint64, error)
	ImportRoaringBitmap(ctx context.Context, tid dax.TableID, fld *dax.Field, shard uint64, views map[string]*roaring.Bitmap, clear bool) error
	ImportRoaringShard(ctx context.Context, tid dax.TableID, shard uint64, request *ImportRoaringShardRequest) error
	EncodeImportValues(ctx context.Context, tid dax.TableID, fld *dax.Field, shard uint64, vals []int64, ids []uint64, clear bool) (path string, data []byte, err error)
	EncodeImport(ctx context.Context, tid dax.TableID, fld *dax.Field, shard uint64, vals, ids []uint64, clear bool) (path string, data []byte, err error)
	DoImport(ctx context.Context, tid dax.TableID, fld *dax.Field, shard uint64, path string, data []byte) error
}

// Transaction contains information related to a block of work that
// needs to be tracked and spans multiple API calls.
type Transaction struct {
	// ID is an arbitrary string identifier. All transactions must have a unique ID.
	ID string `json:"id"`

	// Active notes whether an exclusive transaction is active, or
	// still pending (if other active transactions exist). All
	// non-exclusive transactions are always active.
	Active bool `json:"active"`

	// Exclusive is set to true for transactions which can only become active when no other
	// transactions exist.
	Exclusive bool `json:"exclusive"`

	// Timeout is the minimum idle time for which this transaction should continue to exist.
	Timeout time.Duration `json:"timeout"`

	// CreatedAt is the timestamp at which the transaction was created. This supports
	// the case of listing transactions in a useful order.
	CreatedAt time.Time `json:"createdAt"`

	// Deadline is calculated from Timeout. TODO reset deadline each time there is activity
	// on the transaction. (we can't do this until there is some method of associating a
	// request/call with a transaction)
	Deadline time.Time `json:"deadline"`

	// Stats track statistics for the transaction. Not yet used.
	Stats TransactionStats `json:"stats"`
}
type TransactionStats struct{}

// ImportRoaringShardRequest is the request for the shard
// transactional endpoint.
type ImportRoaringShardRequest struct {
	// Has this request already been forwarded to all replicas? If
	// Remote=false, then the handling server is responsible for
	// ensuring this request is sent to all repliacs before returning
	// a successful response to the client.
	Remote bool
	Views  []RoaringUpdate

	// SuppressLog requests we not write to the write log. Typically
	// that would be because this request is being replayed from a
	// write log.
	SuppressLog bool
}

// RoaringUpdate represents the bits to clear and then set in a particular view.
type RoaringUpdate struct {
	Field string
	View  string

	// Clear is a roaring encoded bitmatrix of bits to clear. For
	// mutex or int-like fields, only the first row is looked at and
	// the bits in that row are cleared from every row.
	Clear []byte

	// Set is the roaring encoded bitmatrix of bits to set. If this is
	// a mutex or int-like field, we'll assume the first shard width
	// of containers is the exists row and we will first clear all
	// bits in those columns and then set
	Set []byte

	// ClearRecords, when true, denotes that Clear should be
	// interpreted as a single row which will be subtracted from every
	// row in this view.
	ClearRecords bool
}

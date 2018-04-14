package storage

import (
	"time"

	api "github.com/elixirhealth/key/pkg/keyapi"
	bstorage "github.com/elixirhealth/service-base/pkg/server/storage"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

const (
	// DefaultType is the default storage type.
	DefaultType = bstorage.Memory

	// DefaultMaxBatchSize is the maximum size of a batch of public keys.
	DefaultMaxBatchSize = 64

	// MaxEntityKeyTypeKeys indicates the maximum number of public keys an entity can have for
	// a given key type.
	MaxEntityKeyTypeKeys = 256

	// DefaultQueryTimeout is the default timeout for DataStore queries.
	DefaultQueryTimeout = 1 * time.Second
)

var (
	// ErrMaxBatchSizeExceeded indicates when the number of public keys an in an add or get
	// request ot the storer exceeds the maximum size.
	ErrMaxBatchSizeExceeded = errors.New("number of public keys in request exceeeds max " +
		"batch size")
)

// Storer manages public key details.
type Storer interface {
	AddPublicKeys(pkds []*api.PublicKeyDetail) error
	GetPublicKeys(pks [][]byte) ([]*api.PublicKeyDetail, error)
	GetEntityPublicKeys(entityID string, kt api.KeyType) ([]*api.PublicKeyDetail, error)
	CountEntityPublicKeys(entityID string, kt api.KeyType) (int, error)
	Close() error
}

// Parameters defines the parameters of the Storer.
type Parameters struct {
	Type                  bstorage.Type
	MaxBatchSize          uint
	AddQueryTimeout       time.Duration
	GetQueryTimeout       time.Duration
	GetEntityQueryTimeout time.Duration
}

// NewDefaultParameters returns a *Parameters object with default values.
func NewDefaultParameters() *Parameters {
	return &Parameters{
		Type:                  DefaultType,
		MaxBatchSize:          DefaultMaxBatchSize,
		AddQueryTimeout:       DefaultQueryTimeout,
		GetQueryTimeout:       DefaultQueryTimeout,
		GetEntityQueryTimeout: DefaultQueryTimeout,
	}
}

// MarshalLogObject writes the parameters to the given object encoder.
func (p *Parameters) MarshalLogObject(oe zapcore.ObjectEncoder) error {
	oe.AddString(logType, p.Type.String())
	oe.AddDuration(logAddQueryTimeout, p.AddQueryTimeout)
	oe.AddDuration(logGetQueryTimeout, p.GetQueryTimeout)
	return nil
}

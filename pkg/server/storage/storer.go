package storage

import (
	"errors"
	"time"

	api "github.com/elxirhealth/key/pkg/keyapi"
	bstorage "github.com/elxirhealth/service-base/pkg/server/storage"
	"go.uber.org/zap/zapcore"
)

const (
	// DefaultType is the default storage type.
	DefaultType = bstorage.Memory

	// DefaultMaxBatchSize is the maximum size of a batch of public keys.
	DefaultMaxBatchSize = 64

	// DefaultQueryTimeout is the default timeout for DataStore queries.
	DefaultQueryTimeout = 1 * time.Second
)

var (
	// ErrNoSuchPublicKey indicates when details for a requested public key do not exist.
	ErrNoSuchPublicKey = errors.New("not details found for given public key")
)

// Storer manages public key details.
type Storer interface {
	AddPublicKeys(pkds []*api.PublicKeyDetail) error
	GetPublicKeys(pks [][]byte) ([]*api.PublicKeyDetail, error)
	GetEntityPublicKeys(entityID string) ([]*api.PublicKeyDetail, error)
	CountEntityPublicKeys(entityID string, kt api.KeyType) (int, error)
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

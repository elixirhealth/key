package storage

import (
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

	// DefaultAddQueryTimeout is the timeout for DataStore queries associated with a Store Add
	// method.
	DefaultAddQueryTimeout = 1 * time.Second

	// DefaultGetQueryTimeout is the timeout for DataStore queries associated with a Store Get
	// method.
	DefaultGetQueryTimeout = 1 * time.Second
)

// Storer manages public key details.
type Storer interface {
	AddPublicKeys(pkds []*api.PublicKeyDetail) error
	GetPublicKeys(pks [][]byte) ([]*api.PublicKeyDetail, error)
}

// Parameters defines the parameters of the Storer.
type Parameters struct {
	Type            bstorage.Type
	MaxBatchSize    uint
	AddQueryTimeout time.Duration
	GetQueryTimeout time.Duration
}

// NewDefaultParameters returns a *Parameters object with default values.
func NewDefaultParameters() *Parameters {
	return &Parameters{
		Type:            DefaultType,
		MaxBatchSize:    DefaultMaxBatchSize,
		AddQueryTimeout: DefaultAddQueryTimeout,
		GetQueryTimeout: DefaultGetQueryTimeout,
	}
}

// MarshalLogObject writes the parameters to the given object encoder.
func (p *Parameters) MarshalLogObject(oe zapcore.ObjectEncoder) error {
	oe.AddString(logType, p.Type.String())
	return nil
}

package storage

import (
	api "github.com/elxirhealth/key/pkg/keyapi"
	bstorage "github.com/elxirhealth/service-base/pkg/server/storage"
	"go.uber.org/zap/zapcore"
)

var (
	// DefaultType is the default storage type.
	DefaultType = bstorage.Memory
)

// Storer manages public key details.
type Storer interface {
	AddPublicKeys(details []*api.PublicKeyDetails) error
	GetPublicKeys(publicKeys [][]byte) ([]*api.PublicKeyDetails, error)
}

// Parameters defines the parameters of the Storer.
type Parameters struct {
	Type bstorage.Type
}

// NewDefaultParameters returns a *Parameters object with default values.
func NewDefaultParameters() *Parameters {
	return &Parameters{
		Type: DefaultType,
	}
}

// MarshalLogObject writes the parameters to the given object encoder.
func (p *Parameters) MarshalLogObject(oe zapcore.ObjectEncoder) error {
	oe.AddString(logType, p.Type.String())
	return nil
}

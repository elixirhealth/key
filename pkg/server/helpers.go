package server

import (
	"errors"

	api "github.com/elxirhealth/key/pkg/keyapi"
	"github.com/elxirhealth/key/pkg/server/storage"
	bstorage "github.com/elxirhealth/service-base/pkg/server/storage"
	"go.uber.org/zap"
)

var (
	// ErrInvalidStorageType indicates when a storage type is not expected.
	ErrInvalidStorageType = errors.New("invalid storage type")
)

func getStorer(config *Config, logger *zap.Logger) (storage.Storer, error) {
	switch config.Storage.Type {
	case bstorage.Memory:
		return storage.NewMemoryStorer(config.Storage, logger), nil
	case bstorage.DataStore:
		return storage.NewDatastore(config.GCPProjectID, config.Storage, logger)
	default:
		return nil, ErrInvalidStorageType
	}
}

func getPublicKeyDetails(rq *api.AddPublicKeysRequest) []*api.PublicKeyDetail {
	pkds := make([]*api.PublicKeyDetail, len(rq.PublicKeys))
	for i, pk := range rq.PublicKeys {
		pkds[i] = &api.PublicKeyDetail{
			PublicKey: pk,
			EntityId:  rq.EntityId,
			KeyType:   rq.KeyType,
		}
	}
	return pkds
}

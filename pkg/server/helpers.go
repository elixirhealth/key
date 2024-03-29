package server

import (
	"errors"

	api "github.com/elixirhealth/key/pkg/keyapi"
	"github.com/elixirhealth/key/pkg/server/storage"
	"github.com/elixirhealth/key/pkg/server/storage/datastore"
	"github.com/elixirhealth/key/pkg/server/storage/memory"
	"github.com/elixirhealth/key/pkg/server/storage/postgres"
	bstorage "github.com/elixirhealth/service-base/pkg/server/storage"
	"go.uber.org/zap"
)

var (
	// ErrInvalidStorageType indicates when a storage type is not expected.
	ErrInvalidStorageType = errors.New("invalid storage type")
)

func getStorer(config *Config, logger *zap.Logger) (storage.Storer, error) {
	switch config.Storage.Type {
	case bstorage.Memory:
		return memory.New(config.Storage, logger), nil
	case bstorage.Postgres:
		return postgres.New(config.DBUrl, config.Storage, logger)
	case bstorage.DataStore:
		return datastore.New(config.GCPProjectID, config.Storage, logger)
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

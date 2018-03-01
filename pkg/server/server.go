package server

import (
	"github.com/elxirhealth/key/pkg/keyapi"
	"github.com/elxirhealth/key/pkg/server/storage"
	"github.com/elxirhealth/service-base/pkg/server"
	"golang.org/x/net/context"
)

// Key implements the KeyServer interface.
type Key struct {
	*server.BaseServer
	config *Config

	storer storage.Storer
	// TODO maybe add other things here
}

// newKey creates a new KeyServer from the given config.
func newKey(config *Config) (*Key, error) {
	baseServer := server.NewBaseServer(config.BaseConfig)
	storer, err := getStorer(config, baseServer.Logger)
	if err != nil {
		return nil, err
	}
	// TODO maybe add other init

	return &Key{
		BaseServer: baseServer,
		config:     config,
		storer:     storer,
		// TODO maybe add other things
	}, nil
}

func (k *Key) AddPublicKeys(
	ctx context.Context, rq *keyapi.AddPublicKeysRequest,
) (*keyapi.AddPublicKeysResponse, error) {
	panic("implement me")
}

func (k *Key) GetPublicKeys(
	ctx context.Context, rq *keyapi.GetPublicKeysRequest,
) (*keyapi.GetPublicKeysResponse, error) {
	panic("implement me")
}

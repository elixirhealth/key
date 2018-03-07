package server

import (
	"math/rand"
	"time"

	api "github.com/elxirhealth/key/pkg/keyapi"
	"github.com/elxirhealth/key/pkg/server/storage"
	"github.com/elxirhealth/service-base/pkg/server"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

var (
	// ErrTooManyActivePublicKeys indicates when adding more public keys would bring the total
	// number of active PKs abvoe the maximum allowed.
	ErrTooManyActivePublicKeys = errors.New("too many active public keys for the entity and " +
		"key type")
)

// Key implements the KeyServer interface.
type Key struct {
	*server.BaseServer
	config *Config

	storer storage.Storer
}

// newKey creates a new KeyServer from the given config.
func newKey(config *Config) (*Key, error) {
	baseServer := server.NewBaseServer(config.BaseConfig)
	storer, err := getStorer(config, baseServer.Logger)
	if err != nil {
		return nil, err
	}
	return &Key{
		BaseServer: baseServer,
		config:     config,
		storer:     storer,
	}, nil
}

// AddPublicKeys adds a set of public keys associated with a given entity.
func (k *Key) AddPublicKeys(
	ctx context.Context, rq *api.AddPublicKeysRequest,
) (*api.AddPublicKeysResponse, error) {
	k.Logger.Debug("received add public keys request", logAddPublicKeysRq(rq)...)
	if err := api.ValidateAddPublicKeysRequest(rq); err != nil {
		return nil, err
	}
	if n, err := k.storer.CountEntityPublicKeys(rq.EntityId, rq.KeyType); err != nil {
		return nil, err
	} else if n+len(rq.PublicKeys) > storage.MaxEntityKeyTypeKeys {
		return nil, ErrTooManyActivePublicKeys
	}
	pkds := getPublicKeyDetails(rq)
	if err := k.storer.AddPublicKeys(pkds); err != nil {
		return nil, err
	}
	k.Logger.Info("added public keys", logAddPublicKeysRq(rq)...)
	return &api.AddPublicKeysResponse{}, nil
}

// GetPublicKeys gets the details (including their associated entity IDs) for a given set of public
// keys.
func (k *Key) GetPublicKeys(
	ctx context.Context, rq *api.GetPublicKeysRequest,
) (*api.GetPublicKeysResponse, error) {
	k.Logger.Debug("received get public keys request", zap.Int(logNKeys, len(rq.PublicKeys)))
	if err := api.ValidateGetPublicKeysRequest(rq); err != nil {
		return nil, err
	}
	pkds, err := k.storer.GetPublicKeys(rq.PublicKeys)
	if err != nil {
		return nil, err
	}
	k.Logger.Info("got public keys", zap.Int(logNKeys, len(pkds)))
	return &api.GetPublicKeysResponse{
		PublicKeyDetails: pkds,
	}, nil
}

// SamplePublicKeys returns a sample of public keys of the given entity.
func (k *Key) SamplePublicKeys(
	ctx context.Context, rq *api.SamplePublicKeysRequest,
) (*api.SamplePublicKeysResponse, error) {
	k.Logger.Debug("received sample public keys request", logSamplePublicKeysRq(rq)...)
	if err := api.ValidateSamplePublicKeysRequest(rq); err != nil {
		return nil, err
	}
	allPKDs, err := k.storer.GetEntityPublicKeys(rq.OfEntityId)
	if err != nil {
		return nil, err
	}
	orderKey := []byte(rq.RequesterEntityId)
	topOrdered := getOrderedLimit(allPKDs, orderKey, api.MaxSamplePublicKeysSize)
	rng := rand.New(rand.NewSource(int64(time.Now().Nanosecond()))) // good enough
	topSampled := sampleWithoutReplacement(topOrdered, rng, int(rq.NPublicKeys))
	k.Logger.Info("sampled public keys", logSamplePublicKeysRq(rq)...)
	return &api.SamplePublicKeysResponse{
		PublicKeyDetails: topSampled,
	}, nil
}

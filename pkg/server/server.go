package server

import (
	"math/rand"
	"time"

	api "github.com/elixirhealth/key/pkg/keyapi"
	"github.com/elixirhealth/key/pkg/server/storage"
	"github.com/elixirhealth/service-base/pkg/server"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrTooManyActivePublicKeys indicates when adding more public keys would bring the total
	// number of active PKs abvoe the maximum allowed.
	ErrTooManyActivePublicKeys = status.Error(codes.FailedPrecondition,
		"too many active public keys for the entity and key type")

	// ErrInternal represents an internal error (e.g., with storage or dependency service call).
	ErrInternal = status.Error(codes.Internal, "internal error")
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
		k.Logger.Info("add public keys request invalid", zap.String(logErr, err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if n, err := k.storer.CountEntityPublicKeys(rq.EntityId, rq.KeyType); err != nil {
		k.Logger.Error("storer count entity public keys error", zap.Error(err))
		return nil, ErrInternal
	} else if n+len(rq.PublicKeys) > storage.MaxEntityKeyTypeKeys {
		return nil, ErrTooManyActivePublicKeys
	}
	pkds := getPublicKeyDetails(rq)
	if err := k.storer.AddPublicKeys(pkds); err != nil {
		k.Logger.Error("storer add public keys error", zap.Error(err))
		return nil, ErrInternal
	}
	k.Logger.Info("added public keys", logAddPublicKeysRq(rq)...)
	return &api.AddPublicKeysResponse{}, nil
}

// GetPublicKeys returns the public keys of a given type for a given entity ID.
func (k *Key) GetPublicKeys(
	ctx context.Context, rq *api.GetPublicKeysRequest,
) (*api.GetPublicKeysResponse, error) {
	k.Logger.Debug("received get public keys request", logGetPublicKeysRq(rq)...)
	if err := api.ValidateGetPublicKeysRequest(rq); err != nil {
		k.Logger.Info("get public keys request invalid", zap.String(logErr, err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	pkds, err := k.storer.GetEntityPublicKeys(rq.EntityId, rq.KeyType)
	if err != nil {
		k.Logger.Error("storer get entity public keys error", zap.Error(err))
		return nil, ErrInternal
	}
	pks := make([][]byte, len(pkds))
	for i, pkd := range pkds {
		pks[i] = pkd.PublicKey
	}
	rp := &api.GetPublicKeysResponse{PublicKeys: pks}
	k.Logger.Info("got public keys", logGetPublicKeysRp(rq, rp)...)
	return rp, nil
}

// GetPublicKeyDetails gets the details (including their associated entity IDs) for a given set of
// public keys.
func (k *Key) GetPublicKeyDetails(
	ctx context.Context, rq *api.GetPublicKeyDetailsRequest,
) (*api.GetPublicKeyDetailsResponse, error) {
	k.Logger.Debug("received get public key details request",
		zap.Int(logNKeys, len(rq.PublicKeys)))
	if err := api.ValidateGetPublicKeyDetailsRequest(rq); err != nil {
		k.Logger.Info("get public key details request invalid",
			zap.String(logErr, err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	pkds, err := k.storer.GetPublicKeys(rq.PublicKeys)
	if err != nil && err == api.ErrNoSuchPublicKey {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		k.Logger.Error("storer get public keys error", zap.Error(err))
		return nil, ErrInternal
	}
	k.Logger.Info("got public key details", zap.Int(logNKeys, len(pkds)))
	return &api.GetPublicKeyDetailsResponse{
		PublicKeyDetails: pkds,
	}, nil
}

// SamplePublicKeys returns a sample of public keys of the given entity.
func (k *Key) SamplePublicKeys(
	ctx context.Context, rq *api.SamplePublicKeysRequest,
) (*api.SamplePublicKeysResponse, error) {
	k.Logger.Debug("received sample public keys request", logSamplePublicKeysRq(rq)...)
	if err := api.ValidateSamplePublicKeysRequest(rq); err != nil {
		k.Logger.Info("sample public keys request invalid", zap.String(logErr, err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	allPKDs, err := k.storer.GetEntityPublicKeys(rq.OfEntityId, api.KeyType_READER)
	if err != nil {
		k.Logger.Error("storer get entity public keys error", zap.Error(err))
		return nil, ErrInternal
	}
	orderKey := []byte(rq.RequesterEntityId)
	topOrdered := getOrderedLimit(allPKDs, orderKey, api.MaxSamplePublicKeysSize)
	rng := rand.New(rand.NewSource(int64(time.Now().Nanosecond()))) // good enough
	topSampled := sampleWithoutReplacement(topOrdered, rng, int(rq.NPublicKeys))
	rp := &api.SamplePublicKeysResponse{
		PublicKeyDetails: topSampled,
	}
	k.Logger.Info("sampled public keys", logSamplePublicKeysRp(rq, rp)...)
	return rp, nil
}

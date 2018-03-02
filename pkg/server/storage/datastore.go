package storage

import (
	"context"
	"encoding/hex"

	"cloud.google.com/go/datastore"
	api "github.com/elxirhealth/key/pkg/keyapi"
	bstorage "github.com/elxirhealth/service-base/pkg/server/storage"
	"go.uber.org/zap"
)

const (
	publicKeyKind = "public_key"
)

// PublicKeyDetail represents a public key and its publicKey, stored in DataStore.
type PublicKeyDetail struct {
	PublicKey string `datastore:"__key__"`
	EntityID  string `datastore:"entity_id"`
	KeyType   string `datastore:"key_type"`
	Disabled  bool   `datastore:"disabled"`
}

type datastoreStorer struct {
	params *Parameters
	client bstorage.DatastoreClient
	logger *zap.Logger
}

// NewDatastore creates a new Store backed by a GCP DataStore instance.
func NewDatastore(gcpProjectID string, params *Parameters, logger *zap.Logger) (Storer, error) {
	client, err := datastore.NewClient(context.Background(), gcpProjectID)
	if err != nil {
		return nil, err
	}
	return &datastoreStorer{
		params: params,
		client: &bstorage.DatastoreClientImpl{Inner: client},
		logger: logger,
	}, nil
}

func (s *datastoreStorer) AddPublicKeys(pkds []*api.PublicKeyDetail) error {
	if err := api.ValidatePublicKeyDetails(pkds); err != nil {
		return err
	}
	sKeys, sDetails := toStoredMulti(pkds)
	ctx, cancel := context.WithTimeout(context.Background(), s.params.AddQueryTimeout)
	defer cancel()
	if _, err := s.client.PutMulti(ctx, sKeys, sDetails); err != nil {
		return err
	}
	s.logger.Debug("added public keys to storage", zap.Int(logNPublicKeys, len(pkds)))
	return nil
}

func (s *datastoreStorer) GetPublicKeys(pks [][]byte) ([]*api.PublicKeyDetail, error) {
	if err := api.ValidatePublicKeys(pks); err != nil {
		return nil, err
	}
	spkds := make([]*PublicKeyDetail, len(pks))
	sKeys := toStoredKeys(pks)
	ctx, cancel := context.WithTimeout(context.Background(), s.params.GetQueryTimeout)
	defer cancel()
	if err := s.client.GetMulti(ctx, sKeys, spkds); err != nil {
		return nil, err
	}
	pkds, err := fromStoredMulti(spkds)
	if err != nil {
		return nil, err
	}
	s.logger.Debug("got public keys from storage", zap.Int(logNPublicKeys, len(pkds)))
	return pkds, nil
}

func toStoredKeys(pks [][]byte) []*datastore.Key {
	keys := make([]*datastore.Key, len(pks))
	for i, pk := range pks {
		keys[i] = datastore.NameKey(publicKeyKind, hex.EncodeToString(pk), nil)
	}
	return keys
}

func toStored(pkd *api.PublicKeyDetail) (*datastore.Key, *PublicKeyDetail) {
	pkHex := hex.EncodeToString(pkd.PublicKey)
	key := datastore.NameKey(publicKeyKind, pkHex, nil)
	return key, &PublicKeyDetail{
		PublicKey: pkHex,
		EntityID:  pkd.EntityId,
		KeyType:   pkd.KeyType.String(),
		Disabled:  false,
	}
}

func toStoredMulti(pkds []*api.PublicKeyDetail) ([]*datastore.Key, []*PublicKeyDetail) {
	keys := make([]*datastore.Key, len(pkds))
	spkds := make([]*PublicKeyDetail, len(pkds))
	for i, pkd := range pkds {
		keys[i], spkds[i] = toStored(pkd)
	}
	return keys, spkds
}

func fromStored(pkd *PublicKeyDetail) (*api.PublicKeyDetail, error) {
	pk, err := hex.DecodeString(pkd.PublicKey)
	if err != nil {
		return nil, err
	}
	return &api.PublicKeyDetail{
		PublicKey: pk,
		EntityId:  pkd.EntityID,
		KeyType:   api.KeyType(api.KeyType_value[pkd.KeyType]),
	}, nil
}

func fromStoredMulti(spkds []*PublicKeyDetail) ([]*api.PublicKeyDetail, error) {
	pkds := make([]*api.PublicKeyDetail, len(spkds))
	for i, spkd := range spkds {
		pkd, err := fromStored(spkd)
		if err != nil {
			return nil, err
		}
		pkds[i] = pkd
	}
	return pkds, nil
}

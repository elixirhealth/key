package datastore

import (
	"context"
	"encoding/hex"
	"time"

	"cloud.google.com/go/datastore"
	api "github.com/elixirhealth/key/pkg/keyapi"
	"github.com/elixirhealth/key/pkg/server/storage"
	bstorage "github.com/elixirhealth/service-base/pkg/server/storage"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
)

const (
	publicKeyKind = "public_key"

	secsPerDay = int64(3600 * 24 * 24)
)

// PublicKeyDetail represents a public key and its publicKey, stored in DataStore.
type PublicKeyDetail struct {
	PublicKey    *datastore.Key `datastore:"__key__"`
	EntityID     string         `datastore:"entity_id"`
	KeyType      string         `datastore:"key_type"`
	Disabled     bool           `datastore:"disabled"`
	ModifiedDate int32          `datastore:"modified_date"`
	ModifiedTime time.Time      `datastore:"modified_time,noindex"`
	AddedTime    time.Time      `datastore:"added_time,noindex"`
	DisabledTime time.Time      `datastore:"disabled_time,noindex"`
}

type storer struct {
	params *storage.Parameters
	client bstorage.DatastoreClient
	iter   bstorage.DatastoreIterator
	logger *zap.Logger
}

// New creates a new Storer backed by a GCP DataStore instance.
func New(
	gcpProjectID string, params *storage.Parameters, logger *zap.Logger,
) (storage.Storer, error) {
	client, err := datastore.NewClient(context.Background(), gcpProjectID)
	if err != nil {
		return nil, err
	}
	return &storer{
		params: params,
		client: &bstorage.DatastoreClientImpl{Inner: client},
		iter:   &bstorage.DatastoreIteratorImpl{},
		logger: logger,
	}, nil
}

func (s *storer) AddPublicKeys(pkds []*api.PublicKeyDetail) error {
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

func (s *storer) GetPublicKeys(pks [][]byte) ([]*api.PublicKeyDetail, error) {
	if err := api.ValidatePublicKeys(pks); err != nil {
		return nil, err
	}
	spkds := make([]*PublicKeyDetail, len(pks))
	sKeys := toStoredKeys(pks)
	ctx, cancel := context.WithTimeout(context.Background(), s.params.GetQueryTimeout)
	defer cancel()
	err := s.client.GetMulti(ctx, sKeys, spkds)
	if err != nil && firstMultiErrNotNil(err) == datastore.ErrNoSuchEntity {
		return nil, api.ErrNoSuchPublicKey
	} else if err != nil {
		return nil, err
	}
	pkds, err := fromStoredMulti(spkds)
	if err != nil {
		return nil, err
	}
	s.logger.Debug("got public keys from storage", zap.Int(logNPublicKeys, len(pkds)))
	return pkds, nil
}

func firstMultiErrNotNil(err error) error {
	switch et := err.(type) {
	case datastore.MultiError:
		for _, e := range et {
			if e != nil {
				return e
			}
		}
	}
	return err
}

func (s *storer) GetEntityPublicKeys(entityID string) ([]*api.PublicKeyDetail, error) {
	if entityID == "" {
		return nil, api.ErrEmptyEntityID
	}
	q := getEntityPublicKeysQuery(entityID, api.KeyType_READER).
		Limit(storage.MaxEntityKeyTypeKeys)
	ctx, cancel := context.WithTimeout(context.Background(), s.params.GetEntityQueryTimeout)
	defer cancel()
	iter := s.client.Run(ctx, q)
	s.iter.Init(iter)
	pkds := make([]*api.PublicKeyDetail, 0, storage.MaxEntityKeyTypeKeys)
	for {
		spkd := &PublicKeyDetail{}
		if _, err := s.iter.Next(spkd); err == iterator.Done {
			// no more results
			break
		} else if err != nil {
			return nil, err
		}
		pkd, err := fromStored(spkd)
		if err != nil {
			return nil, err
		}
		pkds = append(pkds, pkd)
	}
	s.logger.Debug("found public keys for entity", logGetEntityPubKeys(entityID, pkds)...)
	return pkds, nil
}

func (s *storer) CountEntityPublicKeys(entityID string, kt api.KeyType) (int, error) {
	n, err := s.client.Count(context.Background(), getEntityPublicKeysQuery(entityID, kt))
	if err != nil {
		return 0, err
	}
	s.logger.Debug("counted public keys for entity", logCountEntityPubKeys(entityID, kt)...)
	return n, nil
}

func (s *storer) Close() error {
	return nil
}

func getEntityPublicKeysQuery(entityID string, kt api.KeyType) *datastore.Query {
	return datastore.NewQuery(publicKeyKind).
		Filter("entity_id = ", entityID).
		Filter("key_type = ", kt.String()).
		Filter("disabled = ", false)
}

func toStoredKeys(pks [][]byte) []*datastore.Key {
	keys := make([]*datastore.Key, len(pks))
	for i, pk := range pks {
		keys[i] = datastore.NameKey(publicKeyKind, hex.EncodeToString(pk), nil)
	}
	return keys
}

func toStored(pkd *api.PublicKeyDetail, now time.Time) (*datastore.Key, *PublicKeyDetail) {
	pkHex := hex.EncodeToString(pkd.PublicKey)
	key := datastore.NameKey(publicKeyKind, pkHex, nil)
	return key, &PublicKeyDetail{
		PublicKey:    key,
		EntityID:     pkd.EntityId,
		KeyType:      pkd.KeyType.String(),
		Disabled:     false,
		AddedTime:    now,
		ModifiedTime: now,
		ModifiedDate: int32(now.Unix() / secsPerDay),
	}
}

func toStoredMulti(pkds []*api.PublicKeyDetail) ([]*datastore.Key, []*PublicKeyDetail) {
	keys := make([]*datastore.Key, len(pkds))
	spkds := make([]*PublicKeyDetail, len(pkds))
	now := time.Now()
	for i, pkd := range pkds {
		keys[i], spkds[i] = toStored(pkd, now)
	}
	return keys, spkds
}

func fromStored(spkd *PublicKeyDetail) (*api.PublicKeyDetail, error) {
	pk, err := hex.DecodeString(spkd.PublicKey.Name)
	if err != nil {
		return nil, err
	}
	return &api.PublicKeyDetail{
		PublicKey: pk,
		EntityId:  spkd.EntityID,
		KeyType:   api.KeyType(api.KeyType_value[spkd.KeyType]),
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

package storage

import (
	"context"
	"encoding/hex"
	"time"

	"cloud.google.com/go/datastore"
	api "github.com/elxirhealth/key/pkg/keyapi"
	bstorage "github.com/elxirhealth/service-base/pkg/server/storage"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
)

const (
	// MaxEntityKeyTypeKeys indicates the maximum number of public keys an entity can have for
	// a given key type.
	MaxEntityKeyTypeKeys = 256

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

type datastoreStorer struct {
	params *Parameters
	client bstorage.DatastoreClient
	iter   bstorage.DatastoreIterator
	logger *zap.Logger
}

// NewDatastore creates a new Storer backed by a GCP DataStore instance.
func NewDatastore(gcpProjectID string, params *Parameters, logger *zap.Logger) (Storer, error) {
	client, err := datastore.NewClient(context.Background(), gcpProjectID)
	if err != nil {
		return nil, err
	}
	return &datastoreStorer{
		params: params,
		client: &bstorage.DatastoreClientImpl{Inner: client},
		iter:   &bstorage.DatastoreIteratorImpl{},
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
	if err := s.client.GetMulti(ctx, sKeys, spkds); err == datastore.ErrNoSuchEntity {
		return nil, ErrNoSuchPublicKey
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

func (s *datastoreStorer) GetEntityPublicKeys(entityID string) ([]*api.PublicKeyDetail, error) {
	if entityID == "" {
		return nil, api.ErrEmptyEntityID
	}
	q := getEntityPublicKeysQuery(entityID, api.KeyType_READER).Limit(MaxEntityKeyTypeKeys)
	ctx, cancel := context.WithTimeout(context.Background(), s.params.GetEntityQueryTimeout)
	defer cancel()
	iter := s.client.Run(ctx, q)
	s.iter.Init(iter)
	pkds := make([]*api.PublicKeyDetail, 0, MaxEntityKeyTypeKeys)
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

func (s *datastoreStorer) CountEntityPublicKeys(entityID string, kt api.KeyType) (int, error) {
	n, err := s.client.Count(context.Background(), getEntityPublicKeysQuery(entityID, kt))
	if err != nil {
		return 0, err
	}
	s.logger.Debug("counted public keys for entity", logCountEntityPubKeys(entityID, kt)...)
	return n, nil
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

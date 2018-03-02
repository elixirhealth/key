package storage

import (
	"context"
	"math/rand"
	"testing"

	"cloud.google.com/go/datastore"
	api "github.com/elxirhealth/key/pkg/keyapi"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	errTest = errors.New("test error")
)

func TestDatastoreStorer_AddGetPublicKeys_ok(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	params := NewDefaultParameters()
	lg := zap.NewNop()
	s := &datastoreStorer{
		params: params,
		client: &fixedDatastoreClient{
			publicKey: make(map[string]*PublicKeyDetail),
		},
		logger: lg,
	}

	pkds1 := api.NewTestPublicKeyDetails(rng, 8)
	err := s.AddPublicKeys(pkds1)
	assert.Nil(t, err)

	pubKeys := make([][]byte, len(pkds1))
	for i, pkd := range pkds1 {
		pubKeys[i] = pkd.PublicKey
	}
	pkds2, err := s.GetPublicKeys(pubKeys)
	assert.Nil(t, err)
	assert.Equal(t, pkds1, pkds2)
}

func TestDatastoreStorer_AddPublicKeys_err(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	params := NewDefaultParameters()
	lg := zap.NewNop()
	s := &datastoreStorer{
		params: params,
		client: &fixedDatastoreClient{
			putMultiErr: errTest,
		},
		logger: lg,
	}
	pkds := api.NewTestPublicKeyDetails(rng, 8)

	// empty public key details
	err := s.AddPublicKeys(nil)
	assert.Equal(t, api.ErrEmptyPublicKeys, err)

	// datastore client PutMulti error
	err = s.AddPublicKeys(pkds)
	assert.Equal(t, errTest, err)
}

func TestDatastoreStorer_GetPublicKeys_err(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	params := NewDefaultParameters()
	lg := zap.NewNop()
	s := &datastoreStorer{
		params: params,
		client: &fixedDatastoreClient{
			getMultiErr: errTest,
		},
		logger: lg,
	}
	n := 8
	pubKeys := make([][]byte, n)
	for i, pkd := range api.NewTestPublicKeyDetails(rng, n) {
		pubKeys[i] = pkd.PublicKey
	}

	pkds, err := s.GetPublicKeys(nil)
	assert.Equal(t, api.ErrEmptyPublicKeys, err)
	assert.Nil(t, pkds)

	// datastore client GetMulti error
	pkds, err = s.GetPublicKeys(pubKeys)
	assert.Equal(t, errTest, err)
	assert.Nil(t, pkds)
}

func TestToFromStoredMulti(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	pkds1 := api.NewTestPublicKeyDetails(rng, 8)
	sKeys, spkds := toStoredMulti(pkds1)
	assert.Equal(t, len(pkds1), len(sKeys))
	assert.Equal(t, len(pkds1), len(spkds))
	for i, sKey := range sKeys {
		assert.Equal(t, sKey.Name, spkds[i].PublicKey)
	}

	pkds2, err := fromStoredMulti(spkds)
	assert.Nil(t, err)
	assert.Equal(t, pkds1, pkds2)
}

type fixedDatastoreClient struct {
	publicKey   map[string]*PublicKeyDetail
	putMultiErr error
	getMultiErr error
}

func (f *fixedDatastoreClient) PutMulti(
	ctx context.Context, keys []*datastore.Key, values interface{},
) ([]*datastore.Key, error) {
	if f.putMultiErr != nil {
		return nil, f.putMultiErr
	}
	for i, sKey := range keys {
		f.publicKey[sKey.Name] = values.([]*PublicKeyDetail)[i]
	}
	return keys, nil
}

func (f *fixedDatastoreClient) GetMulti(
	ctx context.Context, keys []*datastore.Key, dest interface{},
) error {
	if f.getMultiErr != nil {
		return f.getMultiErr
	}
	for i, sKey := range keys {
		dest.([]*PublicKeyDetail)[i] = f.publicKey[sKey.Name]
	}
	return nil
}

// the methods below aren't used in this test, so we leave them unimplemented

func (f *fixedDatastoreClient) Put(
	ctx context.Context, key *datastore.Key, value interface{},
) (*datastore.Key, error) {
	panic("implement me")
}

func (f *fixedDatastoreClient) Get(
	ctx context.Context, key *datastore.Key, dest interface{},
) error {
	panic("implement me")
}

func (f *fixedDatastoreClient) Delete(ctx context.Context, keys []*datastore.Key) error {
	panic("implement me")
}

func (f *fixedDatastoreClient) Count(ctx context.Context, q *datastore.Query) (int, error) {
	panic("implement me")
}

func (f *fixedDatastoreClient) Run(ctx context.Context, q *datastore.Query) *datastore.Iterator {
	panic("implement me")
}

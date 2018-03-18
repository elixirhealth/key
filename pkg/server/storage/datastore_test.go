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
	"google.golang.org/api/iterator"
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
		client: &fixedDatastoreClient{},
		logger: lg,
	}
	n := 8
	pubKeys := make([][]byte, n)
	for i, pkd := range api.NewTestPublicKeyDetails(rng, n) {
		pubKeys[i] = pkd.PublicKey
	}

	// bad request
	pkds, err := s.GetPublicKeys(nil)
	assert.Equal(t, api.ErrEmptyPublicKeys, err)
	assert.Nil(t, pkds)

	// missing key
	s.client = &fixedDatastoreClient{
		getMultiErr: datastore.MultiError{datastore.ErrNoSuchEntity},
	}
	pkds, err = s.GetPublicKeys(pubKeys)
	assert.Equal(t, api.ErrNoSuchPublicKey, err)
	assert.Nil(t, pkds)

	// other datastore client GetMulti error
	s.client = &fixedDatastoreClient{getMultiErr: errTest}
	pkds, err = s.GetPublicKeys(pubKeys)
	assert.Equal(t, errTest, err)
	assert.Nil(t, pkds)
}

func TestDatastoreStorer_GetEntityPublicKeys_ok(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	params := NewDefaultParameters()
	lg := zap.NewNop()
	pkds1 := api.NewTestPublicKeyDetails(rng, 8)
	keys, spkds := toStoredMulti(pkds1)
	s := &datastoreStorer{
		params: params,
		client: &fixedDatastoreClient{},
		iter: &fixedDatastoreIter{
			keys:   keys,
			values: spkds,
		},
		logger: lg,
	}

	pkds2, err := s.GetEntityPublicKeys("some entity ID")
	assert.Nil(t, err)
	assert.Equal(t, pkds1, pkds2)
}

func TestDatastoreStorer_GetEntityPublicKeys_err(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	params := NewDefaultParameters()
	lg := zap.NewNop()
	s := &datastoreStorer{
		params: params,
		client: &fixedDatastoreClient{},
		iter: &fixedDatastoreIter{
			err: errTest,
		},
		logger: lg,
	}

	// empty entity ID
	pkds, err := s.GetEntityPublicKeys("")
	assert.Equal(t, api.ErrEmptyEntityID, err)
	assert.Nil(t, pkds)

	// next error
	pkds, err = s.GetEntityPublicKeys("some entity ID")
	assert.Equal(t, errTest, err)
	assert.Nil(t, pkds)

	// bad stored value
	badKeys, badSpkds := toStoredMulti(api.NewTestPublicKeyDetails(rng, 1))
	badSpkds[0].PublicKey = datastore.NameKey(publicKeyKind, "*", nil)
	s = &datastoreStorer{
		params: params,
		client: &fixedDatastoreClient{},
		iter: &fixedDatastoreIter{
			keys:   badKeys,
			values: badSpkds,
		},
		logger: lg,
	}
	pkds, err = s.GetEntityPublicKeys("some entity ID")
	assert.NotNil(t, err)
	assert.Nil(t, pkds)
}

func TestDatastoreStorer_CountEntityPublicKeys(t *testing.T) {
	params := NewDefaultParameters()
	lg := zap.NewNop()
	count := 9
	s := &datastoreStorer{
		params: params,
		client: &fixedDatastoreClient{
			countValue: count,
		},
		logger: lg,
	}

	// ok
	val, err := s.CountEntityPublicKeys("some entity ID", api.KeyType_READER)
	assert.Nil(t, err)
	assert.Equal(t, count, val)

	// query err
	s = &datastoreStorer{
		params: params,
		client: &fixedDatastoreClient{
			countErr: errTest,
		},
		logger: lg,
	}
	val, err = s.CountEntityPublicKeys("some entity ID", api.KeyType_READER)
	assert.Equal(t, errTest, err)
	assert.Zero(t, val)
}

func TestToFromStoredMulti(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	pkds1 := api.NewTestPublicKeyDetails(rng, 8)
	sKeys, spkds := toStoredMulti(pkds1)
	assert.Equal(t, len(pkds1), len(sKeys))
	assert.Equal(t, len(pkds1), len(spkds))
	for i, sKey := range sKeys {
		assert.Equal(t, sKey.Name, spkds[i].PublicKey.Name)
		assert.NotZero(t, spkds[i].AddedTime)
		assert.NotZero(t, spkds[i].ModifiedTime)
		assert.Zero(t, spkds[i].DisabledTime)
	}

	pkds2, err := fromStoredMulti(spkds)
	assert.Nil(t, err)
	assert.Equal(t, pkds1, pkds2)
}

type fixedDatastoreClient struct {
	publicKey   map[string]*PublicKeyDetail
	putMultiErr error
	getMultiErr error
	countValue  int
	countErr    error
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
	return f.countValue, f.countErr
}

func (f *fixedDatastoreClient) Run(ctx context.Context, q *datastore.Query) *datastore.Iterator {
	return nil
}

type fixedDatastoreIter struct {
	err    error
	keys   []*datastore.Key
	values []*PublicKeyDetail
	offset int
}

func (f *fixedDatastoreIter) Init(iter *datastore.Iterator) {}

func (f *fixedDatastoreIter) Next(dst interface{}) (*datastore.Key, error) {
	if f.err != nil {
		return nil, f.err
	}
	defer func() { f.offset++ }()
	if f.offset == len(f.values) {
		return nil, iterator.Done
	}
	v := f.values[f.offset]
	dst.(*PublicKeyDetail).EntityID = v.EntityID
	dst.(*PublicKeyDetail).KeyType = v.KeyType
	dst.(*PublicKeyDetail).PublicKey = v.PublicKey
	dst.(*PublicKeyDetail).Disabled = v.Disabled
	return f.keys[f.offset], nil
}

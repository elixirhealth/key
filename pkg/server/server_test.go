package server

import (
	"context"
	"math/rand"
	"testing"

	api "github.com/elxirhealth/key/pkg/keyapi"
	"github.com/elxirhealth/key/pkg/server/storage"
	bserver "github.com/elxirhealth/service-base/pkg/server"
	bstorage "github.com/elxirhealth/service-base/pkg/server/storage"
	"github.com/elxirhealth/service-base/pkg/util"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	errTest = errors.New("some test error")
)

/* TODO (drausin) add back in once have in-memory storer
func TestNewKey_ok(t *testing.T) {
	config := NewDefaultConfig()
	c, err := newKey(config)
	assert.Nil(t, err)
	assert.Equal(t, config, c.config)
	assert.NotEmpty(t, c.storer)
}
*/

func TestNewKey_err(t *testing.T) {
	badConfigs := map[string]*Config{
		"empty ProjectID": NewDefaultConfig().WithStorage(
			&storage.Parameters{Type: bstorage.DataStore},
		),
	}
	for desc, badConfig := range badConfigs {
		c, err := newKey(badConfig)
		assert.NotNil(t, err, desc)
		assert.Nil(t, c)
	}
}

func TestKey_AddPublicKeys_ok(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	k := &Key{
		BaseServer: bserver.NewBaseServer(bserver.NewDefaultBaseConfig()),
		storer:     &fixedStorer{},
	}
	rq := &api.AddPublicKeysRequest{
		EntityId: "some entity ID",
		KeyType:  api.KeyType_READER,
		PublicKeys: [][]byte{
			util.RandBytes(rng, 33),
			util.RandBytes(rng, 33),
		},
	}
	rp, err := k.AddPublicKeys(context.Background(), rq)
	assert.Nil(t, err)
	assert.NotNil(t, rp)
}

func TestKey_AddPublicKeys_err(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	k := &Key{
		BaseServer: bserver.NewBaseServer(bserver.NewDefaultBaseConfig()),
		storer:     &fixedStorer{addErr: errTest},
	}

	// bad request
	rq := &api.AddPublicKeysRequest{}
	rp, err := k.AddPublicKeys(context.Background(), rq)
	assert.NotNil(t, err)
	assert.Nil(t, rp)

	// storer error
	rq = &api.AddPublicKeysRequest{
		EntityId: "some entity ID",
		KeyType:  api.KeyType_READER,
		PublicKeys: [][]byte{
			util.RandBytes(rng, 33),
			util.RandBytes(rng, 33),
		},
	}
	rp, err = k.AddPublicKeys(context.Background(), rq)
	assert.Equal(t, errTest, err)
	assert.Nil(t, rp)
}

func TestKey_GetPublicKeys_ok(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	pks := [][]byte{
		util.RandBytes(rng, 33),
		util.RandBytes(rng, 33),
	}
	k := &Key{
		BaseServer: bserver.NewBaseServer(bserver.NewDefaultBaseConfig()),
		storer: &fixedStorer{
			getPKDs: api.NewTestPublicKeyDetails(rng, len(pks)),
		},
	}
	rq := &api.GetPublicKeysRequest{PublicKeys: pks}
	rp, err := k.GetPublicKeys(context.Background(), rq)
	assert.Nil(t, err)
	assert.NotNil(t, rp)
	assert.Equal(t, len(pks), len(rq.PublicKeys))
}

func TestKey_GetPublicKeys_err(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	k := &Key{
		BaseServer: bserver.NewBaseServer(bserver.NewDefaultBaseConfig()),
		storer:     &fixedStorer{getErr: errTest},
	}

	// bad request
	rq := &api.GetPublicKeysRequest{}
	rp, err := k.GetPublicKeys(context.Background(), rq)
	assert.NotNil(t, err)
	assert.Nil(t, rp)

	// storer error
	rq = &api.GetPublicKeysRequest{
		PublicKeys: [][]byte{
			util.RandBytes(rng, 33),
			util.RandBytes(rng, 33),
		},
	}
	rp, err = k.GetPublicKeys(context.Background(), rq)
	assert.Equal(t, errTest, err)
	assert.Nil(t, rp)
}

type fixedStorer struct {
	addErr  error
	getPKDs []*api.PublicKeyDetail
	getErr  error
}

func (f *fixedStorer) GetEntityPublicKeysCount(entityID string, kt api.KeyType) (int, error) {
	panic("implement me")
}

func (f *fixedStorer) GetEntityPublicKeys(entityID string) ([]*api.PublicKeyDetail, error) {
	panic("implement me")
}

func (f *fixedStorer) AddPublicKeys(pkds []*api.PublicKeyDetail) error {
	return f.addErr
}

func (f *fixedStorer) GetPublicKeys(pks [][]byte) ([]*api.PublicKeyDetail, error) {
	return f.getPKDs, f.getErr
}

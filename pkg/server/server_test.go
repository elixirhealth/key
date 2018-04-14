package server

import (
	"context"
	"math/rand"
	"testing"

	api "github.com/elixirhealth/key/pkg/keyapi"
	"github.com/elixirhealth/key/pkg/server/storage"
	bserver "github.com/elixirhealth/service-base/pkg/server"
	bstorage "github.com/elixirhealth/service-base/pkg/server/storage"
	"github.com/elixirhealth/service-base/pkg/util"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errTest = errors.New("some test error")
)

func TestNewKey_ok(t *testing.T) {
	config := NewDefaultConfig()
	c, err := newKey(config)
	assert.Nil(t, err)
	assert.Equal(t, config, c.config)
	assert.NotEmpty(t, c.storer)
}

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
	baseServer := bserver.NewBaseServer(bserver.NewDefaultBaseConfig())
	okRq := &api.AddPublicKeysRequest{
		EntityId: "some entity ID",
		KeyType:  api.KeyType_READER,
		PublicKeys: [][]byte{
			util.RandBytes(rng, 33),
			util.RandBytes(rng, 33),
		},
	}
	cases := map[string]struct {
		k  *Key
		rq *api.AddPublicKeysRequest
	}{
		"bad request": {
			k: &Key{
				BaseServer: baseServer,
				storer:     &fixedStorer{},
			},
			rq: &api.AddPublicKeysRequest{},
		},
		"storer get count error": {
			k: &Key{
				BaseServer: baseServer,
				storer:     &fixedStorer{countEntityPKsErr: errTest},
			},
			rq: okRq,
		},
		"too many added": {
			k: &Key{
				BaseServer: baseServer,
				storer:     &fixedStorer{countEntityPKsValue: 255},
			},
			rq: okRq,
		},
		"storer add error": {
			k: &Key{
				BaseServer: baseServer,
				storer:     &fixedStorer{addErr: errTest},
			},
			rq: okRq,
		},
	}
	for desc, c := range cases {
		rp, err := c.k.AddPublicKeys(context.Background(), c.rq)
		assert.NotNil(t, err, desc)
		assert.Nil(t, rp, desc)
	}
}

func TestKey_GetPublicKeys_ok(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	n := 2
	k := &Key{
		BaseServer: bserver.NewBaseServer(bserver.NewDefaultBaseConfig()),
		storer: &fixedStorer{
			getEntityPKs: api.NewTestPublicKeyDetails(rng, n),
		},
	}
	rq := &api.GetPublicKeysRequest{EntityId: "some entity ID"}
	rp, err := k.GetPublicKeys(context.Background(), rq)
	assert.Nil(t, err)
	assert.NotNil(t, rp)
	assert.Equal(t, n, len(rp.PublicKeys))
}

func TestKey_GetPublicKeys_err(t *testing.T) {
	k := &Key{
		BaseServer: bserver.NewBaseServer(bserver.NewDefaultBaseConfig()),
		storer:     &fixedStorer{getEntityPKsErr: errTest},
	}

	// bad request
	rq := &api.GetPublicKeysRequest{}
	rp, err := k.GetPublicKeys(context.Background(), rq)
	assert.NotNil(t, err)
	assert.Nil(t, rp)

	// storer error
	rq = &api.GetPublicKeysRequest{EntityId: "some entity ID"}
	rp, err = k.GetPublicKeys(context.Background(), rq)
	assert.Equal(t, errTest, err)
	assert.Nil(t, rp)
}

func TestKey_GetPublicKeyDetails_ok(t *testing.T) {
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
	rq := &api.GetPublicKeyDetailsRequest{PublicKeys: pks}
	rp, err := k.GetPublicKeyDetails(context.Background(), rq)
	assert.Nil(t, err)
	assert.NotNil(t, rp)
	assert.Equal(t, len(pks), len(rp.PublicKeyDetails))
}

func TestKey_GetPublicKeyDetails_err(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	k := &Key{
		BaseServer: bserver.NewBaseServer(bserver.NewDefaultBaseConfig()),
		storer:     &fixedStorer{getErr: errTest},
	}

	// bad request
	rq := &api.GetPublicKeyDetailsRequest{}
	rp, err := k.GetPublicKeyDetails(context.Background(), rq)
	assert.NotNil(t, err)
	assert.Nil(t, rp)

	// storer error
	rq = &api.GetPublicKeyDetailsRequest{
		PublicKeys: [][]byte{
			util.RandBytes(rng, 33),
			util.RandBytes(rng, 33),
		},
	}
	rp, err = k.GetPublicKeyDetails(context.Background(), rq)
	assert.Equal(t, errTest, err)
	assert.Nil(t, rp)

	// no such pub key
	k = &Key{
		BaseServer: bserver.NewBaseServer(bserver.NewDefaultBaseConfig()),
		storer:     &fixedStorer{getErr: api.ErrNoSuchPublicKey},
	}
	rp, err = k.GetPublicKeyDetails(context.Background(), rq)
	assert.NotNil(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
	assert.Equal(t, api.ErrNoSuchPublicKey.Error(), status.Convert(err).Message())
	assert.Nil(t, rp)
}

func TestKey_SamplePublicKeys_ok(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	nEntityPKDs := 64
	ctx := context.Background()
	ofEntityID, rqEntityID := "some entity ID", "another entity ID"
	k := &Key{
		BaseServer: bserver.NewBaseServer(bserver.NewDefaultBaseConfig()),
		storer: &fixedStorer{
			getEntityPKs: api.NewTestPublicKeyDetails(rng, nEntityPKDs),
		},
	}

	rp1, err := k.SamplePublicKeys(ctx, &api.SamplePublicKeysRequest{
		OfEntityId:        ofEntityID,
		NPublicKeys:       2,
		RequesterEntityId: rqEntityID,
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(rp1.PublicKeyDetails))

	// check sample again yields diff result
	rp2, err := k.SamplePublicKeys(ctx, &api.SamplePublicKeysRequest{
		OfEntityId:        ofEntityID,
		NPublicKeys:       2,
		RequesterEntityId: rqEntityID,
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(rp2.PublicKeyDetails))
	assert.NotEqual(t, rp1.PublicKeyDetails, rp2.PublicKeyDetails)

	// check 2 samples of max public keys size yield same result
	rq := &api.SamplePublicKeysRequest{
		OfEntityId:        ofEntityID,
		NPublicKeys:       api.MaxSamplePublicKeysSize,
		RequesterEntityId: rqEntityID,
	}
	rp3, err := k.SamplePublicKeys(ctx, rq)
	assert.Nil(t, err)
	rp4, err := k.SamplePublicKeys(ctx, rq)
	assert.Nil(t, err)
	assert.Equal(t, rp3, rp4)

	// check another sample with diff requester has diff result
	rp5, err := k.SamplePublicKeys(ctx, &api.SamplePublicKeysRequest{
		OfEntityId:        ofEntityID,
		NPublicKeys:       api.MaxSamplePublicKeysSize,
		RequesterEntityId: "diff requester",
	})
	assert.Nil(t, err)
	assert.NotEqual(t, rp4, rp5)
}

func TestKey_SamplePublicKeys_err(t *testing.T) {
	k := &Key{
		BaseServer: bserver.NewBaseServer(bserver.NewDefaultBaseConfig()),
		storer:     &fixedStorer{getEntityPKsErr: errTest},
	}
	ofEntityID, rqEntityID := "some entity ID", "another entity ID"

	// bad request
	rq := &api.SamplePublicKeysRequest{}
	rp, err := k.SamplePublicKeys(context.Background(), rq)
	assert.NotNil(t, err)
	assert.Nil(t, rp)

	// storer error
	rq = &api.SamplePublicKeysRequest{
		OfEntityId:        ofEntityID,
		NPublicKeys:       api.MaxSamplePublicKeysSize,
		RequesterEntityId: rqEntityID,
	}
	rp, err = k.SamplePublicKeys(context.Background(), rq)
	assert.Equal(t, errTest, err)
	assert.Nil(t, rp)
}

type fixedStorer struct {
	addErr              error
	getPKDs             []*api.PublicKeyDetail
	getErr              error
	countEntityPKsValue int
	countEntityPKsErr   error
	getEntityPKs        []*api.PublicKeyDetail
	getEntityPKsErr     error
}

func (f *fixedStorer) CountEntityPublicKeys(entityID string, kt api.KeyType) (int, error) {
	return f.countEntityPKsValue, f.countEntityPKsErr
}

func (f *fixedStorer) GetEntityPublicKeys(
	entityID string, kt api.KeyType,
) ([]*api.PublicKeyDetail, error) {
	return f.getEntityPKs, f.getEntityPKsErr
}

func (f *fixedStorer) AddPublicKeys(pkds []*api.PublicKeyDetail) error {
	return f.addErr
}

func (f *fixedStorer) GetPublicKeys(pks [][]byte) ([]*api.PublicKeyDetail, error) {
	return f.getPKDs, f.getErr
}

func (f *fixedStorer) Close() error {
	return nil
}

package memory

import (
	"math/rand"
	"testing"

	api "github.com/elixirhealth/key/pkg/keyapi"
	"github.com/elixirhealth/key/pkg/server/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMemoryStorer_AddGetPublicKeys_ok(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	params := storage.NewDefaultParameters()
	lg := zap.NewNop()
	s := New(params, lg)

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

func TestMemoryStorer_AddPublicKeys_err(t *testing.T) {
	params := storage.NewDefaultParameters()
	lg := zap.NewNop()
	s := New(params, lg)

	// empty public key details
	err := s.AddPublicKeys(nil)
	assert.Equal(t, api.ErrEmptyPublicKeys, err)
}

func TestMemoryStorer_GetPublicKeys_err(t *testing.T) {
	params := storage.NewDefaultParameters()
	lg := zap.NewNop()
	s := New(params, lg)

	// bad request
	pkds, err := s.GetPublicKeys(nil)
	assert.Equal(t, api.ErrEmptyPublicKeys, err)
	assert.Nil(t, pkds)

	// missing key
	pkds, err = s.GetPublicKeys([][]byte{{1, 2, 3}})
	assert.Equal(t, api.ErrNoSuchPublicKey, err)
	assert.Nil(t, pkds)

	// missing key again, to check lock
	pkds, err = s.GetPublicKeys([][]byte{{1, 2, 3}})
	assert.Equal(t, api.ErrNoSuchPublicKey, err)
	assert.Nil(t, pkds)
}

func TestMemoryStorer_GetEntityPublicKeys_ok(t *testing.T) {
	params := storage.NewDefaultParameters()
	lg := zap.NewNop()
	s := New(params, lg)

	rng := rand.New(rand.NewSource(0))
	pkds1 := api.NewTestPublicKeyDetails(rng, 64)
	err := s.AddPublicKeys(pkds1)
	assert.Nil(t, err)

	pkds2, err := s.GetEntityPublicKeys(pkds1[0].EntityId, api.KeyType_READER)
	assert.Nil(t, err)
	expectedN := 0
	for _, pkd1 := range pkds1 {
		if pkd1.EntityId == pkds1[0].EntityId && pkd1.KeyType == api.KeyType_READER {
			expectedN++
		}
	}
	assert.Equal(t, expectedN, len(pkds2))
}

func TestMemoryStorer_GetEntityPublicKeys_err(t *testing.T) {
	params := storage.NewDefaultParameters()
	lg := zap.NewNop()
	s := New(params, lg)

	pkds, err := s.GetEntityPublicKeys("", api.KeyType_READER)
	assert.Equal(t, api.ErrEmptyEntityID, err)
	assert.Nil(t, pkds)
}

func TestMemoryStorer_CountEntityPublicKeys_ok(t *testing.T) {
	params := storage.NewDefaultParameters()
	lg := zap.NewNop()
	s := New(params, lg)

	rng := rand.New(rand.NewSource(0))
	pkds1 := api.NewTestPublicKeyDetails(rng, 64)
	err := s.AddPublicKeys(pkds1)
	assert.Nil(t, err)

	kt := api.KeyType_AUTHOR
	n, err := s.CountEntityPublicKeys(pkds1[0].EntityId, kt)
	assert.Nil(t, err)
	expectedN := 0
	for _, pkd1 := range pkds1 {
		if pkd1.EntityId == pkds1[0].EntityId && pkd1.KeyType == kt {
			expectedN++
		}
	}
	assert.Equal(t, expectedN, n)
}

func TestMemoryStorer_CountEntityPublicKeys_err(t *testing.T) {
	params := storage.NewDefaultParameters()
	lg := zap.NewNop()
	s := New(params, lg)

	n, err := s.CountEntityPublicKeys("", api.KeyType_AUTHOR)
	assert.Equal(t, api.ErrEmptyEntityID, err)
	assert.Zero(t, n)
}

package keyapi

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAddPublicKeysRequest(t *testing.T) {
	cases := map[string]struct {
		rq       *AddPublicKeysRequest
		expected error
	}{
		"ok": {
			rq: &AddPublicKeysRequest{
				EntityId:   "some entity ID",
				PublicKeys: [][]byte{{1, 2, 3}},
			},
			expected: nil,
		},
		"missing entity ID": {
			rq: &AddPublicKeysRequest{
				PublicKeys: [][]byte{{1, 2, 3}},
			},
			expected: ErrEmptyEntityID,
		},
		"missing public keys": {
			rq: &AddPublicKeysRequest{
				EntityId: "some entity ID",
			},
			expected: ErrEmptyPublicKeys,
		},
	}
	for _, c := range cases {
		err := ValidateAddPublicKeysRequest(c.rq)
		assert.Equal(t, c.expected, err)
	}
}

func TestValidateGetPublicKeysRequest(t *testing.T) {
	cases := map[string]struct {
		rq       *GetPublicKeysRequest
		expected error
	}{
		"ok": {
			rq: &GetPublicKeysRequest{
				PublicKeys: [][]byte{{1, 2, 3}},
			},
			expected: nil,
		},
		"missing public keys": {
			rq:       &GetPublicKeysRequest{},
			expected: ErrEmptyPublicKeys,
		},
	}
	for _, c := range cases {
		err := ValidateGetPublicKeysRequest(c.rq)
		assert.Equal(t, c.expected, err)
	}
}

func TestValidatePublicKeyDetails(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	okPKD := NewTestPublicKeyDetail(rng)
	cases := map[string]struct {
		pkds     []*PublicKeyDetail
		expected error
	}{
		"ok": {
			pkds:     []*PublicKeyDetail{okPKD},
			expected: nil,
		},
		"nil value": {
			pkds:     nil,
			expected: ErrEmptyPublicKeys,
		},
		"zero-len value": {
			pkds:     []*PublicKeyDetail{},
			expected: ErrEmptyPublicKeys,
		},
		"pkd missing required fields": {
			pkds:     []*PublicKeyDetail{{}},
			expected: ErrEmptyPublicKey,
		},
		"duplicate pkd": {
			pkds:     []*PublicKeyDetail{okPKD, okPKD},
			expected: ErrDupPublicKeys,
		},
	}
	for _, c := range cases {
		err := ValidatePublicKeyDetails(c.pkds)
		assert.Equal(t, c.expected, err)
	}
}

func TestValidatePublicKeyDetail(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	okPKD := NewTestPublicKeyDetail(rng)
	cases := map[string]struct {
		pkd      *PublicKeyDetail
		expected error
	}{
		"ok": {
			pkd:      okPKD,
			expected: nil,
		},
		"nil value": {
			pkd:      nil,
			expected: ErrEmptyPublicKeyDetail,
		},
		"empty public key": {
			pkd:      &PublicKeyDetail{},
			expected: ErrEmptyPublicKey,
		},
		"missing entity ID": {
			pkd: &PublicKeyDetail{
				PublicKey: []byte{1, 2, 3},
				EntityId:  "",
			},
			expected: ErrEmptyEntityID,
		},
	}
	for _, c := range cases {
		err := ValidatePublicKeyDetail(c.pkd)
		assert.Equal(t, c.expected, err)
	}
}

func TestValidatePublicKeys(t *testing.T) {
	cases := map[string]struct {
		pks      [][]byte
		expected error
	}{
		"ok": {
			pks:      [][]byte{{1, 2, 3}},
			expected: nil,
		},
		"nil value": {
			pks:      nil,
			expected: ErrEmptyPublicKeys,
		},
		"zero-len value": {
			pks:      [][]byte{},
			expected: ErrEmptyPublicKeys,
		},
		"empty pub key": {
			pks:      [][]byte{{}},
			expected: ErrEmptyPublicKey,
		},
		"duplicate pub keys": {
			pks:      [][]byte{{1, 2, 3}, {1, 2, 3}},
			expected: ErrDupPublicKeys,
		},
	}
	for _, c := range cases {
		err := ValidatePublicKeys(c.pks)
		assert.Equal(t, c.expected, err)
	}
}

func TestValidatePublicKey(t *testing.T) {
	cases := map[string]struct {
		pk       []byte
		expected error
	}{
		"ok": {
			pk:       []byte{1, 2, 3},
			expected: nil,
		},
		"nil value": {
			pk:       nil,
			expected: ErrEmptyPublicKey,
		},
		"zero-len value": {
			pk:       []byte{},
			expected: ErrEmptyPublicKey,
		},
	}
	for _, c := range cases {
		err := ValidatePublicKey(c.pk)
		assert.Equal(t, c.expected, err)
	}
}

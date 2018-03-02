package keyapi

import (
	"encoding/hex"
	"fmt"
	"math/rand"

	"github.com/elxirhealth/service-base/pkg/util"
	"github.com/pkg/errors"
)

var (
	// ErrEmptyPublicKeys indicates when a list of public keys is nil or zero length.
	ErrEmptyPublicKeys = errors.New("empty public keys list")

	// ErrDupPublicKeys indicates when a list of public keys or public key details has
	// duplicate public keys.
	ErrDupPublicKeys = errors.New("duplicate public keys in list")

	// ErrEmptyPublicKeyDetail indicates when a public key detail value is nil.
	ErrEmptyPublicKeyDetail = errors.New("empty public key detail value")

	// ErrEmptyPublicKey indicates when a public key is nil or zero length.
	ErrEmptyPublicKey = errors.New("empty public key field")

	// ErrEmptyEntityID indicates when the entity ID of a public key detail value is missing.
	ErrEmptyEntityID = errors.New("empty entity ID field")
)

// TODO add ValidateENDPOINTRequest method for each service ENDPOINT

// ValidatePublicKeyDetails checks that the list of public key details isn't empty, has no dups,
// and has valid public key detail elements.
func ValidatePublicKeyDetails(pkds []*PublicKeyDetail) error {
	if len(pkds) == 0 {
		return ErrEmptyPublicKeys
	}
	pks := map[string]struct{}{}
	for _, pkd := range pkds {
		if err := ValidatePublicKeyDetail(pkd); err != nil {
			return err
		}
		pkHex := hex.EncodeToString(pkd.PublicKey)
		if _, in := pks[pkHex]; in {
			return ErrDupPublicKeys
		}
		pks[pkHex] = struct{}{}
	}
	return nil
}

// ValidatePublicKeyDetail checks that a public key detail is not empty and has all fields
// populated.
func ValidatePublicKeyDetail(pkd *PublicKeyDetail) error {
	if pkd == nil {
		return ErrEmptyPublicKeyDetail
	}
	if err := ValidatePublicKey(pkd.PublicKey); err != nil {
		return err
	}
	if pkd.EntityId == "" {
		return ErrEmptyEntityID
	}
	return nil
}

// ValidatePublicKeys checks that a list of public keys is not empty, has no dups, and has
// non-empty elements.
func ValidatePublicKeys(pks [][]byte) error {
	if len(pks) == 0 {
		return ErrEmptyPublicKeys
	}
	pkSet := map[string]struct{}{}
	for _, pk := range pks {
		if err := ValidatePublicKey(pk); err != nil {
			return err
		}
		pkHex := hex.EncodeToString(pk)
		if _, in := pkSet[pkHex]; in {
			return ErrDupPublicKeys
		}
		pkSet[pkHex] = struct{}{}
	}
	return nil
}

// ValidatePublicKey checks that a public key is not nil or empty.
func ValidatePublicKey(pk []byte) error {
	if len(pk) == 0 {
		return ErrEmptyPublicKey
	}
	return nil
}

// NewTestPublicKeyDetail creates a random *PublicKeyDetail for use in testing.
func NewTestPublicKeyDetail(rng *rand.Rand) *PublicKeyDetail {
	return &PublicKeyDetail{
		PublicKey: util.RandBytes(rng, 33),
		EntityId:  fmt.Sprintf("EntityID-%d", rng.Intn(4)),
		KeyType:   KeyType(rng.Intn(2)),
	}
}

// NewTestPublicKeyDetails creates a list of random *PublicKeyDetails for use in testing.
func NewTestPublicKeyDetails(rng *rand.Rand, n int) []*PublicKeyDetail {
	pkds := make([]*PublicKeyDetail, n)
	for i := range pkds {
		pkds[i] = NewTestPublicKeyDetail(rng)
	}
	return pkds
}

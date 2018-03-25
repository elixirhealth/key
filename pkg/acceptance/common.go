package acceptance

import (
	"fmt"
	"math/rand"

	"github.com/elixirhealth/service-base/pkg/util"
)

// CreateTestEntityKeys creates a new entity (from index i) and some random author and reader
// public keys, suitable only for testing.
func CreateTestEntityKeys(rng *rand.Rand, i, nKeyTypeKeys uint) (string, [][]byte, [][]byte) {
	authorKeys := make([][]byte, nKeyTypeKeys)
	readerKeys := make([][]byte, nKeyTypeKeys)
	for i := range authorKeys {
		authorKeys[i] = util.RandBytes(rng, 33)
		readerKeys[i] = util.RandBytes(rng, 33)
	}
	return GetTestEntityID(i), authorKeys, readerKeys
}

// GetTestEntityID returns the ID for the i'th test entity.
func GetTestEntityID(i uint) string {
	return fmt.Sprintf("Entity-%d", i)
}

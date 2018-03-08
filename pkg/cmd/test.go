package cmd

import (
	"context"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/drausin/libri/libri/common/logging"
	"github.com/drausin/libri/libri/common/parse"
	"github.com/elxirhealth/key/pkg/acceptance"
	api "github.com/elxirhealth/key/pkg/keyapi"
	bcmd "github.com/elxirhealth/service-base/pkg/cmd"
	bserver "github.com/elxirhealth/service-base/pkg/server"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	timeoutFlag = "timeout"

	logEntityID          = "entity_id"
	logKeyType           = "key_type"
	logNKeys             = "n_keys"
	logExpected          = "expected"
	logActual            = "actual"
	logAuthorKeyShortHex = "author_key_short_hex"
	logReaderKeyShortHex = "reader_key_short_hex"
)

func testIO() error {
	rng := rand.New(rand.NewSource(0))
	logger := logging.NewDevLogger(logging.GetLogLevel(viper.GetString(logLevelFlag)))
	timeout := time.Duration(viper.GetInt(timeoutFlag) * 1e9)
	nEntities := uint(4)
	nKeyTypeKeys := uint(64)
	nGets := uint(16)

	clients, err := getClients()
	if err != nil {
		return err
	}

	entityAuthorKeys := make(map[string][][]byte)
	entityReaderKeys := make(map[string][][]byte)
	authorKeyEntities := make(map[string]string)
	readerKeyEntities := make(map[string]string)

	// create & add keys for each entity
	for c := uint(0); c < nEntities; c++ {
		entityID, authorKeys, readerKeys :=
			acceptance.CreateTestEntityKeys(rng, c, nKeyTypeKeys)
		entityAuthorKeys[entityID] = authorKeys
		entityReaderKeys[entityID] = readerKeys
		for i := range authorKeys {
			authorKeyEntities[hex.EncodeToString(authorKeys[i])] = entityID
			readerKeyEntities[hex.EncodeToString(readerKeys[i])] = entityID
		}

		rq := &api.AddPublicKeysRequest{
			EntityId:   entityID,
			KeyType:    api.KeyType_AUTHOR,
			PublicKeys: authorKeys,
		}
		client := clients[rng.Int31n(int32(len(clients)))]
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		_, err := client.AddPublicKeys(ctx, rq)
		cancel()
		if logAddPublicKeysRq(logger, rq, err) != nil {
			return err
		}

		rq = &api.AddPublicKeysRequest{
			EntityId:   entityID,
			KeyType:    api.KeyType_READER,
			PublicKeys: readerKeys,
		}
		client = clients[rng.Int31n(int32(len(clients)))]
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		_, err = client.AddPublicKeys(ctx, rq)
		cancel()
		if err2 := logAddPublicKeysRq(logger, rq, err); err2 != nil {
			return err
		}
	}

	// get keys
	for c := uint(0); c < nGets; c++ {
		entityID := acceptance.GetTestEntityID(c % 4)
		// get one random author key, and one random reader key
		authorKey := entityAuthorKeys[entityID][rng.Intn(len(entityAuthorKeys))]
		readerKey := entityReaderKeys[entityID][rng.Intn(len(entityReaderKeys))]
		authorEntityID := authorKeyEntities[hex.EncodeToString(authorKey)]
		readerEntityID := readerKeyEntities[hex.EncodeToString(readerKey)]

		rq := &api.GetPublicKeysRequest{
			PublicKeys: [][]byte{authorKey, readerKey},
		}
		client := clients[rng.Int31n(int32(len(clients)))]
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		rp, err := client.GetPublicKeys(ctx, rq)
		cancel()
		err2 := logGetPublicKeysRp(logger, authorEntityID, readerEntityID, authorKey,
			readerKey, rp, err)
		if err2 != nil {
			return err
		}

	}

	// sample key
	for c := uint(0); c < nEntities; c++ {
		entityID := acceptance.GetTestEntityID(c)
		rq := &api.SamplePublicKeysRequest{
			OfEntityId:        entityID,
			RequesterEntityId: "some requester",
			NPublicKeys:       1,
		}
		client := clients[rng.Int31n(int32(len(clients)))]
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		rp, err := client.SamplePublicKeys(ctx, rq)
		cancel()
		err2 := logSamplePublicKeysRp(logger, entityID, readerKeyEntities, rp, err)
		if err2 != nil {
			return err
		}
	}

	return nil
}

func logAddPublicKeysRq(logger *zap.Logger, rq *api.AddPublicKeysRequest, err error) error {
	if err != nil {
		logger.Error("adding public keys failed", zap.Error(err))
		return err
	}
	logger.Info("added public keys",
		zap.String(logEntityID, rq.EntityId),
		zap.Stringer(logKeyType, rq.KeyType),
		zap.Int(logNKeys, len(rq.PublicKeys)),
	)
	return nil
}

func logGetPublicKeysRp(
	logger *zap.Logger,
	authorEntityID, readerEntityID string,
	authorKey, readerKey []byte,
	rp *api.GetPublicKeysResponse, err error,
) error {
	if err != nil {
		logger.Error("get public keys failed", zap.Error(err))
		return err
	}
	if authorEntityID != rp.PublicKeyDetails[0].EntityId {
		logger.Error("unexpected entity ID for gotten author key",
			zap.String(logExpected, authorEntityID),
			zap.String(logActual, rp.PublicKeyDetails[0].EntityId),
		)
		return err
	}
	if readerEntityID != rp.PublicKeyDetails[1].EntityId {
		logger.Error("unexpected entity ID for gotten reader key",
			zap.String(logExpected, authorEntityID),
			zap.String(logActual, rp.PublicKeyDetails[0].EntityId),
		)
		return err
	}
	logger.Info("got public key details",
		zap.String(logAuthorKeyShortHex, hex.EncodeToString(authorKey[:8])),
		zap.String(logReaderKeyShortHex, hex.EncodeToString(readerKey[:8])),
	)
	return nil
}

func logSamplePublicKeysRp(
	logger *zap.Logger,
	entityID string,
	readerKeyEntities map[string]string,
	rp *api.SamplePublicKeysResponse,
	err error,
) error {
	if err != nil {
		logger.Error("sample public keys failed", zap.String(logEntityID, entityID))
		return err
	}
	pkHex := hex.EncodeToString(rp.PublicKeyDetails[0].PublicKey)
	if entityID != readerKeyEntities[pkHex] {
		logger.Error("unexpected entityID for sampled key",
			zap.String(logExpected, entityID),
			zap.String(logActual, readerKeyEntities[pkHex]),
		)
		return errors.New("unexpected entityID for sampled key")
	}
	pkShortHex := hex.EncodeToString(rp.PublicKeyDetails[0].PublicKey[:8])
	logger.Info("sampled public key",
		zap.String(logEntityID, entityID),
		zap.String(logReaderKeyShortHex, pkShortHex),
	)
	return nil
}

func getClients() ([]api.KeyClient, error) {
	addrs, err := parse.Addrs(viper.GetStringSlice(bcmd.AddressesFlag))
	if err != nil {
		return nil, err
	}
	dialer := bserver.NewInsecureDialer()
	clients := make([]api.KeyClient, len(addrs))
	for i, addr := range addrs {
		conn, err2 := dialer.Dial(addr.String())
		if err != nil {
			return nil, err2
		}
		clients[i] = api.NewKeyClient(conn)
	}
	return clients, nil
}

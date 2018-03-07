// +build acceptance

package acceptance

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"

	"github.com/drausin/libri/libri/common/errors"
	api "github.com/elxirhealth/key/pkg/keyapi"
	"github.com/elxirhealth/key/pkg/server"
	"github.com/elxirhealth/key/pkg/server/storage"
	bstorage "github.com/elxirhealth/service-base/pkg/server/storage"
	"github.com/elxirhealth/service-base/pkg/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

type parameters struct {
	nKeys        uint
	gcpProjectID string
	logLevel     zapcore.Level

	nEntities             uint
	nKeysPerEntityKeyType uint
	nGets                 uint
	timeout               time.Duration
}

type state struct {
	keys          []*server.Key
	keyClients    []api.KeyClient
	dataDir       string
	datastoreProc *os.Process
	rng           *rand.Rand

	entityAuthorKeys  map[string][][]byte
	entityReaderKeys  map[string][][]byte
	authorKeyEntities map[string]string
	readerKeyEntities map[string]string
}

func (st *state) randClient() api.KeyClient {
	return st.keyClients[st.rng.Int31n(int32(len(st.keyClients)))]
}

func TestAcceptance(t *testing.T) {
	params := &parameters{
		nKeys:        3,
		gcpProjectID: "dummy-acceptance-id",
		logLevel:     zapcore.InfoLevel,

		nEntities:             4,
		nKeysPerEntityKeyType: 64,
		nGets:   32,
		timeout: 1 * time.Second,
	}
	st := setUp(params)

	testAdd(t, params, st)

	testGet(t, params, st)

	testSample(t, params, st)

	tearDown(st)
}

func testAdd(t *testing.T, params *parameters, st *state) {
	for c := uint(0); c < params.nEntities; c++ {
		entityID := fmt.Sprintf("Entity-%d", c)
		authorKeys := make([][]byte, params.nKeysPerEntityKeyType)
		readerKeys := make([][]byte, params.nKeysPerEntityKeyType)
		st.entityAuthorKeys[entityID] = authorKeys
		st.entityReaderKeys[entityID] = readerKeys
		for i := range authorKeys {
			authorKeys[i] = util.RandBytes(st.rng, 33)
			readerKeys[i] = util.RandBytes(st.rng, 33)
			st.authorKeyEntities[hex.EncodeToString(authorKeys[i])] = entityID
			st.readerKeyEntities[hex.EncodeToString(readerKeys[i])] = entityID
		}
		st.entityAuthorKeys[entityID] = authorKeys
		st.entityReaderKeys[entityID] = readerKeys

		rq := &api.AddPublicKeysRequest{
			EntityId:   entityID,
			KeyType:    api.KeyType_AUTHOR,
			PublicKeys: authorKeys,
		}
		ctx, cancel := context.WithTimeout(context.Background(), params.timeout)
		_, err := st.randClient().AddPublicKeys(ctx, rq)
		cancel()
		assert.Nil(t, err)

		rq = &api.AddPublicKeysRequest{
			EntityId:   entityID,
			KeyType:    api.KeyType_READER,
			PublicKeys: readerKeys,
		}
		ctx, cancel = context.WithTimeout(context.Background(), params.timeout)
		_, err = st.randClient().AddPublicKeys(ctx, rq)
		cancel()
		assert.Nil(t, err)
	}
}

func testGet(t *testing.T, params *parameters, st *state) {
	for c := uint(0); c < params.nGets; c++ {
		entityID := fmt.Sprintf("Entity-%d", c%4)
		// get one random author key, and one random reader key
		authorKey := st.entityAuthorKeys[entityID][st.rng.Intn(len(st.entityAuthorKeys))]
		readerKey := st.entityReaderKeys[entityID][st.rng.Intn(len(st.entityReaderKeys))]
		authorEntityID := st.authorKeyEntities[hex.EncodeToString(authorKey)]
		readerEntityID := st.readerKeyEntities[hex.EncodeToString(readerKey)]

		rq := &api.GetPublicKeysRequest{
			PublicKeys: [][]byte{authorKey, readerKey},
		}
		ctx, cancel := context.WithTimeout(context.Background(), params.timeout)
		rp, err := st.randClient().GetPublicKeys(ctx, rq)
		cancel()
		assert.Nil(t, err)
		assert.Equal(t, authorEntityID, rp.PublicKeyDetails[0].EntityId)
		assert.Equal(t, readerEntityID, rp.PublicKeyDetails[1].EntityId)
	}
}

func testSample(t *testing.T, params *parameters, st *state) {
	for c := uint(0); c < params.nEntities; c++ {
		entityID := fmt.Sprintf("Entity-%d", c)
		rq := &api.SamplePublicKeysRequest{
			OfEntityId:        entityID,
			RequesterEntityId: "some requester",
			NPublicKeys:       1,
		}
		ctx, cancel := context.WithTimeout(context.Background(), params.timeout)
		rp, err := st.randClient().SamplePublicKeys(ctx, rq)
		cancel()
		assert.Nil(t, err)
		assert.Equal(t, 1, len(rp.PublicKeyDetails))
		pkHex := hex.EncodeToString(rp.PublicKeyDetails[0].PublicKey)
		assert.Equal(t, entityID, st.readerKeyEntities[pkHex])
	}
}

func setUp(params *parameters) *state {
	rng := rand.New(rand.NewSource(0))

	dataDir, err := ioutil.TempDir("", "key-datastore-test")
	errors.MaybePanic(err)
	datastoreProc := bstorage.StartDatastoreEmulator(dataDir)

	time.Sleep(5 * time.Second)
	st := &state{
		rng:           rng,
		dataDir:       dataDir,
		datastoreProc: datastoreProc,

		entityAuthorKeys:  make(map[string][][]byte),
		entityReaderKeys:  make(map[string][][]byte),
		authorKeyEntities: make(map[string]string),
		readerKeyEntities: make(map[string]string),
	}
	createAndStartKeys(params, st)
	return st
}

func createAndStartKeys(params *parameters, st *state) {
	configs, addrs := newKeyConfigs(params)
	keys := make([]*server.Key, params.nKeys)
	keyClients := make([]api.KeyClient, params.nKeys)
	up := make(chan *server.Key, 1)

	for i := uint(0); i < params.nKeys; i++ {
		go func() {
			err := server.Start(configs[i], up)
			errors.MaybePanic(err)
		}()

		// wait for server to come up
		keys[i] = <-up

		// set up client to it
		conn, err := grpc.Dial(addrs[i].String(), grpc.WithInsecure())
		errors.MaybePanic(err)
		keyClients[i] = api.NewKeyClient(conn)
	}

	st.keys = keys
	st.keyClients = keyClients
}

func newKeyConfigs(params *parameters) ([]*server.Config, []*net.TCPAddr) {
	startPort := uint(10100)
	configs := make([]*server.Config, params.nKeys)
	addrs := make([]*net.TCPAddr, params.nKeys)

	// set eviction params to ensure that evictions actually happen during test
	storageParams := storage.NewDefaultParameters()
	storageParams.Type = bstorage.DataStore

	for i := uint(0); i < params.nKeys; i++ {
		serverPort, metricsPort := startPort+i*10, startPort+i*10+1
		configs[i] = server.NewDefaultConfig().
			WithStorage(storageParams).
			WithGCPProjectID(params.gcpProjectID)
		configs[i].WithServerPort(uint(serverPort)).
			WithMetricsPort(uint(metricsPort)).
			WithLogLevel(params.logLevel)
		addrs[i] = &net.TCPAddr{IP: net.ParseIP("localhost"), Port: int(serverPort)}
	}
	return configs, addrs
}

func tearDown(st *state) {
	for _, c := range st.keys {
		c.StopServer()
	}
	bstorage.StopDatastoreEmulator(st.datastoreProc)
	err := os.RemoveAll(st.dataDir)
	errors.MaybePanic(err)
}

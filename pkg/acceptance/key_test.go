// +build acceptance

package acceptance

import (
	"context"
	"encoding/hex"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/drausin/libri/libri/common/errors"
	"github.com/drausin/libri/libri/common/logging"
	api "github.com/elixirhealth/key/pkg/keyapi"
	"github.com/elixirhealth/key/pkg/server"
	"github.com/elixirhealth/key/pkg/server/storage"
	"github.com/elixirhealth/key/pkg/server/storage/postgres/migrations"
	bstorage "github.com/elixirhealth/service-base/pkg/server/storage"
	"github.com/mattes/migrate/source/go-bindata"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

type parameters struct {
	nKeys        uint
	gcpProjectID string
	logLevel     zapcore.Level

	nEntities    uint
	nKeyTypeKeys uint
	nGets        uint
	timeout      time.Duration
}

type state struct {
	keys             []*server.Key
	keyClients       []api.KeyClient
	rng              *rand.Rand
	dbURL            string
	tearDownPostgres func() error

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

		nEntities:    4,
		nKeyTypeKeys: 64,
		nGets:        32,
		timeout:      1 * time.Second,
	}
	st := setUp(t, params)

	testAdd(t, params, st)

	testGet(t, params, st)

	testGetDetails(t, params, st)

	testSample(t, params, st)

	tearDown(t, st)
}

func testAdd(t *testing.T, params *parameters, st *state) {
	for c := uint(0); c < params.nEntities; c++ {
		entityID, authorKeys, readerKeys := CreateTestEntityKeys(st.rng, c, params.nKeyTypeKeys)
		st.entityAuthorKeys[entityID] = authorKeys
		st.entityReaderKeys[entityID] = readerKeys
		for i := range authorKeys {
			st.authorKeyEntities[hex.EncodeToString(authorKeys[i])] = entityID
			st.readerKeyEntities[hex.EncodeToString(readerKeys[i])] = entityID
		}

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
	for c := uint(0); c < params.nEntities; c++ {
		entityID := GetTestEntityID(c % 4)
		rq := &api.GetPublicKeysRequest{
			EntityId: entityID,
			KeyType:  api.KeyType_READER,
		}
		ctx, cancel := context.WithTimeout(context.Background(), params.timeout)
		rp, err := st.randClient().GetPublicKeys(ctx, rq)
		cancel()
		assert.Nil(t, err)
		assert.Equal(t, len(st.entityReaderKeys[entityID]), len(rp.PublicKeys))
		assert.Equal(t, getPKSet(st.entityReaderKeys[entityID]), getPKSet(rp.PublicKeys))
	}
}

func getPKSet(pks [][]byte) map[string]struct{} {
	pkSet := make(map[string]struct{})
	for _, pk := range pks {
		pkSet[hex.EncodeToString(pk)] = struct{}{}
	}
	return pkSet
}

func testGetDetails(t *testing.T, params *parameters, st *state) {
	for c := uint(0); c < params.nGets; c++ {
		entityID := GetTestEntityID(c % 4)
		// get one random author key, and one random reader key
		authorKey := st.entityAuthorKeys[entityID][st.rng.Intn(len(st.entityAuthorKeys))]
		readerKey := st.entityReaderKeys[entityID][st.rng.Intn(len(st.entityReaderKeys))]
		authorEntityID := st.authorKeyEntities[hex.EncodeToString(authorKey)]
		readerEntityID := st.readerKeyEntities[hex.EncodeToString(readerKey)]

		rq := &api.GetPublicKeyDetailsRequest{
			PublicKeys: [][]byte{authorKey, readerKey},
		}
		ctx, cancel := context.WithTimeout(context.Background(), params.timeout)
		rp, err := st.randClient().GetPublicKeyDetails(ctx, rq)
		cancel()
		assert.Nil(t, err)
		assert.Equal(t, authorEntityID, rp.PublicKeyDetails[0].EntityId)
		assert.Equal(t, readerEntityID, rp.PublicKeyDetails[1].EntityId)
	}
}

func testSample(t *testing.T, params *parameters, st *state) {
	for c := uint(0); c < params.nEntities; c++ {
		entityID := GetTestEntityID(c)
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
		if len(rp.PublicKeyDetails) == 1 {
			pkHex := hex.EncodeToString(rp.PublicKeyDetails[0].PublicKey)
			assert.Equal(t, entityID, st.readerKeyEntities[pkHex])
		}
	}
}

func setUp(t *testing.T, params *parameters) *state {
	rng := rand.New(rand.NewSource(0))
	dbURL, cleanup, err := bstorage.StartTestPostgres()
	if err != nil {
		t.Fatal(err)
	}

	st := &state{
		rng:              rng,
		dbURL:            dbURL,
		tearDownPostgres: cleanup,

		entityAuthorKeys:  make(map[string][][]byte),
		entityReaderKeys:  make(map[string][][]byte),
		authorKeyEntities: make(map[string]string),
		readerKeyEntities: make(map[string]string),
	}
	createAndStartKeys(params, st)
	return st
}

func createAndStartKeys(params *parameters, st *state) {
	configs, addrs := newKeyConfigs(params, st)
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

func newKeyConfigs(params *parameters, st *state) ([]*server.Config, []*net.TCPAddr) {
	startPort := uint(10100)
	configs := make([]*server.Config, params.nKeys)
	addrs := make([]*net.TCPAddr, params.nKeys)

	// set eviction params to ensure that evictions actually happen during test
	storageParams := storage.NewDefaultParameters()
	storageParams.Type = bstorage.Postgres

	for i := uint(0); i < params.nKeys; i++ {
		serverPort, metricsPort := startPort+i*10, startPort+i*10+1
		configs[i] = server.NewDefaultConfig().
			WithStorage(storageParams).
			WithDBUrl(st.dbURL)
		configs[i].WithServerPort(uint(serverPort)).
			WithMetricsPort(uint(metricsPort)).
			WithLogLevel(params.logLevel)
		addrs[i] = &net.TCPAddr{IP: net.ParseIP("localhost"), Port: int(serverPort)}
	}
	return configs, addrs
}

func tearDown(t *testing.T, st *state) {
	for _, c := range st.keys {
		c.StopServer()
	}
	logger := &bstorage.ZapLogger{Logger: logging.NewDevInfoLogger()}
	m := bstorage.NewBindataMigrator(
		st.dbURL,
		bindata.Resource(migrations.AssetNames(), migrations.Asset),
		logger,
	)
	err := m.Down()
	assert.Nil(t, err)

	err = st.tearDownPostgres()
	assert.Nil(t, err)
}

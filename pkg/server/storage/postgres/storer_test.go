package postgres

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"testing"

	sq "github.com/Masterminds/squirrel"
	errors2 "github.com/drausin/libri/libri/common/errors"
	"github.com/drausin/libri/libri/common/logging"
	api "github.com/elixirhealth/key/pkg/keyapi"
	"github.com/elixirhealth/key/pkg/server/storage"
	"github.com/elixirhealth/key/pkg/server/storage/postgres/migrations"
	bstorage "github.com/elixirhealth/service-base/pkg/server/storage"
	"github.com/mattes/migrate/source/go-bindata"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	errTest = errors.New("test error")
)

func setUpPostgresTest() (string, func() error) {
	dbURL, cleanup, err := bstorage.StartTestPostgres()
	errors2.MaybePanic(err)
	as := bindata.Resource(migrations.AssetNames(), migrations.Asset)
	logger := &bstorage.LogLogger{}
	m := bstorage.NewBindataMigrator(dbURL, as, logger)
	errors2.MaybePanic(m.Up())
	return dbURL, func() error {
		if err := m.Down(); err != nil {
			return err
		}
		return cleanup()
	}
}

func TestStorer_AddGetPublicKeys_ok(t *testing.T) {
	dbURL, tearDown := setUpPostgresTest()
	defer func() {
		err := tearDown()
		assert.Nil(t, err)
	}()

	rng := rand.New(rand.NewSource(0))
	pkds1 := api.NewTestPublicKeyDetails(rng, 8)
	params := storage.NewDefaultParameters()
	params.Type = bstorage.Postgres
	lg := logging.NewDevLogger(zap.DebugLevel)
	s, err := New(dbURL, params, lg)
	assert.Nil(t, err)

	err = s.AddPublicKeys(pkds1)
	assert.Nil(t, err)

	pubKeys := make([][]byte, len(pkds1))
	for i, pkd := range pkds1 {
		pubKeys[i] = pkd.PublicKey
	}
	pkds2, err := s.GetPublicKeys(pubKeys)
	assert.Nil(t, err)
	assert.Equal(t, len(pkds1), len(pkds2))
}

func TestStorer_AddPublicKeys_err(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	params := storage.NewDefaultParameters()
	params.Type = bstorage.Postgres
	lg := logging.NewDevLogger(zap.DebugLevel)

	cases := map[string]struct {
		s        *storer
		pkds     []*api.PublicKeyDetail
		expected error
	}{
		"bad PKDs": {
			s:        &storer{params: params},
			pkds:     []*api.PublicKeyDetail{},
			expected: api.ErrEmptyPublicKeys,
		},
		"batch too large": {
			s:        &storer{params: params},
			pkds:     api.NewTestPublicKeyDetails(rng, 128),
			expected: storage.ErrMaxBatchSizeExceeded,
		},
		"insert err": {
			s: &storer{
				params: params,
				logger: lg,
				qr: &fixedQuerier{
					insertErr: errTest,
				},
			},
			pkds:     api.NewTestPublicKeyDetails(rng, 8),
			expected: errTest,
		},
	}
	for desc, c := range cases {
		err := c.s.AddPublicKeys(c.pkds)
		assert.Equal(t, c.expected, err, desc)
	}
}

func TestStorer_GetPublicKeys_err(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	params := storage.NewDefaultParameters()
	params.Type = bstorage.Postgres
	lg := logging.NewDevLogger(zap.DebugLevel)
	n := 128
	pubKeys := make([][]byte, n)
	for i, pkd := range api.NewTestPublicKeyDetails(rng, n) {
		pubKeys[i] = pkd.PublicKey
	}

	cases := map[string]struct {
		s        *storer
		pks      [][]byte
		expected error
	}{
		"bad PKDs": {
			s:        &storer{params: params},
			pks:      [][]byte{},
			expected: api.ErrEmptyPublicKeys,
		},
		"batch too large": {
			s:        &storer{params: params},
			pks:      pubKeys,
			expected: storage.ErrMaxBatchSizeExceeded,
		},
		"select err": {
			s: &storer{
				params: params,
				logger: lg,
				qr: &fixedQuerier{
					selectErr: errTest,
				},
			},
			pks:      pubKeys[:8],
			expected: errTest,
		},
		"rows scan err": {
			s: &storer{
				params: params,
				logger: lg,
				qr: &fixedQuerier{
					selectResult: &fixedRowScanner{
						next:    true,
						scanErr: errTest,
					},
				},
			},
			pks:      pubKeys[:8],
			expected: errTest,
		},
		"rows err err": {
			s: &storer{
				params: params,
				logger: lg,
				qr: &fixedQuerier{
					selectResult: &fixedRowScanner{
						errErr: errTest,
					},
				},
			},
			pks:      pubKeys[:8],
			expected: errTest,
		},
	}
	for desc, c := range cases {
		pkds, err := c.s.GetPublicKeys(c.pks)
		assert.Equal(t, c.expected, err, desc)
		assert.Nil(t, pkds)
	}
}

func TestStorer_GetCountEntityPublicKeys_ok(t *testing.T) {
	dbURL, tearDown := setUpPostgresTest()
	defer func() {
		err := tearDown()
		assert.Nil(t, err)
	}()

	rng := rand.New(rand.NewSource(0))
	params := storage.NewDefaultParameters()
	params.Type = bstorage.Postgres
	lg := logging.NewDevLogger(zap.DebugLevel)
	pkds1 := api.NewTestPublicKeyDetails(rng, 64)

	s, err := New(dbURL, params, lg)
	assert.Nil(t, err)

	err = s.AddPublicKeys(pkds1)
	assert.Nil(t, err)

	entityID := pkds1[0].EntityId
	pkds2, err := s.GetEntityPublicKeys(entityID, api.KeyType_READER)
	assert.Nil(t, err)
	assert.True(t, len(pkds2) > 1)
	for _, pkd := range pkds2 {
		assert.Equal(t, entityID, pkd.EntityId)
		assert.Equal(t, api.KeyType_READER, pkd.KeyType)
		assert.NotEmpty(t, pkd.PublicKey)
	}

	n, err := s.CountEntityPublicKeys(entityID, api.KeyType_READER)
	assert.Nil(t, err)
	assert.Equal(t, len(pkds2), n)
}

func TestStorer_GetEntityPublicKeys_err(t *testing.T) {
	params := storage.NewDefaultParameters()
	params.Type = bstorage.Postgres
	lg := logging.NewDevLogger(zap.DebugLevel)
	entityID := "some entity ID"

	cases := map[string]struct {
		s        *storer
		entityID string
		expected error
	}{
		"bad entityID": {
			s:        &storer{params: params},
			entityID: "",
			expected: api.ErrEmptyEntityID,
		},
		"select err": {
			s: &storer{
				params: params,
				logger: lg,
				qr: &fixedQuerier{
					selectErr: errTest,
				},
			},
			entityID: entityID,
			expected: errTest,
		},
		"rows scan err": {
			s: &storer{
				params: params,
				logger: lg,
				qr: &fixedQuerier{
					selectResult: &fixedRowScanner{
						next:    true,
						scanErr: errTest,
					},
				},
			},
			entityID: entityID,
			expected: errTest,
		},
		"rows err err": {
			s: &storer{
				params: params,
				logger: lg,
				qr: &fixedQuerier{
					selectResult: &fixedRowScanner{
						errErr: errTest,
					},
				},
			},
			entityID: entityID,
			expected: errTest,
		},
	}
	for desc, c := range cases {
		pkds, err := c.s.GetEntityPublicKeys(c.entityID, api.KeyType_READER)
		assert.Equal(t, c.expected, err, desc)
		assert.Nil(t, pkds)
	}
}

func TestStorer_CountEntityPublicKeys_err(t *testing.T) {
	params := storage.NewDefaultParameters()
	params.Type = bstorage.Postgres
	lg := logging.NewDevLogger(zap.DebugLevel)
	entityID := "some entity ID"

	cases := map[string]struct {
		s        *storer
		entityID string
		expected error
	}{
		"bad entityID": {
			s:        &storer{params: params},
			entityID: "",
			expected: api.ErrEmptyEntityID,
		},
		"select row err": {
			s: &storer{
				params: params,
				logger: lg,
				qr: &fixedQuerier{
					selectRowResult: &fixedRowScanner{
						scanErr: errTest,
					},
				},
			},
			entityID: entityID,
			expected: errTest,
		},
	}
	for desc, c := range cases {
		pkds, err := c.s.CountEntityPublicKeys(c.entityID, api.KeyType_READER)
		assert.Equal(t, c.expected, err, desc)
		assert.Zero(t, pkds)
	}
}

type fixedQuerier struct {
	selectResult    bstorage.QueryRows
	selectErr       error
	selectRowResult sq.RowScanner
	insertResult    sql.Result
	insertErr       error
}

func (f *fixedQuerier) SelectQueryContext(
	ctx context.Context, b sq.SelectBuilder,
) (bstorage.QueryRows, error) {
	return f.selectResult, f.selectErr
}

func (f *fixedQuerier) SelectQueryRowContext(
	ctx context.Context, b sq.SelectBuilder,
) sq.RowScanner {
	return f.selectRowResult
}

func (f *fixedQuerier) InsertExecContext(
	ctx context.Context, b sq.InsertBuilder,
) (sql.Result, error) {
	return f.insertResult, f.insertErr
}

func (f *fixedQuerier) UpdateExecContext(
	ctx context.Context, b sq.UpdateBuilder,
) (sql.Result, error) {
	panic("implement me")
}

func (f *fixedQuerier) DeleteExecContext(
	ctx context.Context, b sq.DeleteBuilder,
) (sql.Result, error) {
	panic("implement me")
}

type fixedRowScanner struct {
	next    bool
	scanErr error
	errErr  error
}

func (f *fixedRowScanner) Next() bool {
	return f.next
}

func (f *fixedRowScanner) Close() error {
	panic("implement me")
}

func (f *fixedRowScanner) Err() error {
	return f.errErr
}

func (f *fixedRowScanner) Scan(...interface{}) error {
	return f.scanErr
}

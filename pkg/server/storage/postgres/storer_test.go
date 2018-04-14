package postgres

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"testing"

	"github.com/Masterminds/squirrel"
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
					selectRows: &fixedRowScanner{
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
					selectRows: &fixedRowScanner{
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

func TestStorer_GetEntityPublicKeys(t *testing.T) {

}

type fixedQuerier struct {
	selectRows   bstorage.QueryRows
	selectErr    error
	insertResult sql.Result
	insertErr    error
}

func (f *fixedQuerier) SelectQueryContext(ctx context.Context, b squirrel.SelectBuilder) (bstorage.QueryRows, error) {
	return f.selectRows, f.selectErr
}

func (f *fixedQuerier) SelectQueryRowContext(ctx context.Context, b squirrel.SelectBuilder) squirrel.RowScanner {
	panic("implement me")
}

func (f *fixedQuerier) InsertExecContext(ctx context.Context, b squirrel.InsertBuilder) (sql.Result, error) {
	return f.insertResult, f.insertErr
}

func (f *fixedQuerier) UpdateExecContext(ctx context.Context, b squirrel.UpdateBuilder) (sql.Result, error) {
	panic("implement me")
}

func (f *fixedQuerier) DeleteExecContext(ctx context.Context, b squirrel.DeleteBuilder) (sql.Result, error) {
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

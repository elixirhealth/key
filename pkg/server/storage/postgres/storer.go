package postgres

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"

	sq "github.com/Masterminds/squirrel"
	errors2 "github.com/drausin/libri/libri/common/errors"
	api "github.com/elixirhealth/key/pkg/keyapi"
	"github.com/elixirhealth/key/pkg/server/storage"
	bstorage "github.com/elixirhealth/service-base/pkg/server/storage"
	"go.uber.org/zap"
)

const (
	keySchema            = "key"
	publicKeyDetailTable = "public_key_detail"

	publicKeyCol = "public_key"
	keyTypeCol   = "key_type"
	entityIDCol  = "entity_id"

	count = "COUNT(*)"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	fqPublicKeyDetailTable = keySchema + "." + publicKeyDetailTable

	errEmptyDBUrl            = errors.New("empty DB URL")
	errUnexpectedStorageType = errors.New("unexpected storage type")
)

type storer struct {
	params  *storage.Parameters
	db      *sql.DB
	dbCache sq.DBProxyContext
	qr      bstorage.Querier
	logger  *zap.Logger
}

// New creates a new storage.Storer backed by a Postgres DB at the given dbURL.
func New(dbURL string, params *storage.Parameters, logger *zap.Logger) (storage.Storer, error) {
	if dbURL == "" {
		return nil, errEmptyDBUrl
	}
	if params.Type != bstorage.Postgres {
		return nil, errUnexpectedStorageType
	}
	db, err := sql.Open("postgres", dbURL)
	errors2.MaybePanic(err)
	return &storer{
		params:  params,
		db:      db,
		dbCache: sq.NewStmtCacher(db),
		qr:      bstorage.NewQuerier(),
		logger:  logger,
	}, nil
}

func (s *storer) AddPublicKeys(pkds []*api.PublicKeyDetail) error {
	if err := api.ValidatePublicKeyDetails(pkds); err != nil {
		return err
	}
	if len(pkds) > int(s.params.MaxBatchSize) {
		return storage.ErrMaxBatchSizeExceeded
	}
	q := psql.RunWith(s.db).
		Insert(fqPublicKeyDetailTable).
		Columns(pkdSQLCols...)
	for _, pkd := range pkds {
		q = q.Values(getPKDSQLValues(pkd)...)
	}
	s.logger.Debug("adding public keys to storage", logAddingPublicKeys(q, pkds)...)
	ctx, cancel := context.WithTimeout(context.Background(), s.params.AddQueryTimeout)
	defer cancel()
	_, err := s.qr.InsertExecContext(ctx, q)
	if err != nil {
		return err
	}
	s.logger.Debug("added public keys to storage", logAddedPublicKeys(pkds)...)
	return nil
}

func (s *storer) GetPublicKeys(pks [][]byte) ([]*api.PublicKeyDetail, error) {
	if err := api.ValidatePublicKeys(pks); err != nil {
		return nil, err
	}
	if len(pks) > int(s.params.MaxBatchSize) {
		return nil, storage.ErrMaxBatchSizeExceeded
	}
	cols, _, _ := prepPKDScan()
	q := psql.RunWith(s.dbCache).
		Select(cols...).
		From(fqPublicKeyDetailTable).
		Where(sq.Eq{publicKeyCol: pks})
	s.logger.Debug("getting public keys from storage", logGettingPublicKeys(q, pks)...)
	pkds, err := s.getPKDsFromQuery(q, len(pks))
	if err != nil {
		return nil, err
	}
	s.logger.Debug("got public keys from storage", zap.Int(logNPublicKeys, len(pkds)))
	return orderPKDs(pkds, pks), nil
}

func (s *storer) GetEntityPublicKeys(
	entityID string, kt api.KeyType,
) ([]*api.PublicKeyDetail, error) {
	if entityID == "" {
		return nil, api.ErrEmptyEntityID
	}
	cols, _, _ := prepPKDScan()
	q := psql.RunWith(s.dbCache).
		Select(cols...).
		From(fqPublicKeyDetailTable).
		Where(sq.Eq{entityIDCol: entityID, keyTypeCol: kt.String()})
	s.logger.Debug("getting entity public keys from storage",
		logGettingEntityPubKeys(q, entityID)...)
	pkds, err := s.getPKDsFromQuery(q, storage.MaxEntityKeyTypeKeys)
	if err != nil {
		return nil, err
	}
	s.logger.Debug("got entity public keys from storage",
		logGotEntityPubKeys(entityID, pkds)...)
	return pkds, nil
}

func (s *storer) CountEntityPublicKeys(entityID string, kt api.KeyType) (int, error) {
	if entityID == "" {
		return 0, api.ErrEmptyEntityID
	}
	q := psql.RunWith(s.dbCache).
		Select(count).
		From(fqPublicKeyDetailTable).
		Where(sq.Eq{
			entityIDCol: entityID,
			keyTypeCol:  kt.String(),
		})
	s.logger.Debug("counting public keys for entity",
		logCountingEntityPubKeys(q, entityID, kt)...)
	ctx, cancel := context.WithTimeout(context.Background(), s.params.GetEntityQueryTimeout)
	defer cancel()
	row := s.qr.SelectQueryRowContext(ctx, q)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	s.logger.Debug("counted public keys for entity",
		logCountEntityPubKeys(entityID, kt, count)...)
	return count, nil
}

func (s *storer) getPKDsFromQuery(q sq.SelectBuilder, size int) ([]*api.PublicKeyDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.params.GetQueryTimeout)
	defer cancel()
	rows, err := s.qr.SelectQueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	pkds := make([]*api.PublicKeyDetail, size)
	i := 0
	for rows.Next() {
		_, dest, create := prepPKDScan()
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		pkds[i] = create()
		i++
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return pkds[:i], nil
}

func (s *storer) Close() error {
	return s.db.Close()
}

func orderPKDs(pkds []*api.PublicKeyDetail, byPKs [][]byte) []*api.PublicKeyDetail {
	pkdsMap := make(map[string]*api.PublicKeyDetail)
	for _, pkd := range pkds {
		pkHex := hex.EncodeToString(pkd.PublicKey)
		pkdsMap[pkHex] = pkd
	}
	ordered := make([]*api.PublicKeyDetail, 0, len(pkds))
	for _, byPK := range byPKs {
		pkHex := hex.EncodeToString(byPK)
		if pkd, in := pkdsMap[pkHex]; in {
			ordered = append(ordered, pkd)
		}
	}
	return ordered
}

var pkdSQLCols = []string{
	publicKeyCol,
	keyTypeCol,
	entityIDCol,
}

func getPKDSQLValues(pkd *api.PublicKeyDetail) []interface{} {
	return []interface{}{
		pkd.PublicKey,
		pkd.KeyType.String(),
		pkd.EntityId,
	}
}

func prepPKDScan() ([]string, []interface{}, func() *api.PublicKeyDetail) {
	pkd := &api.PublicKeyDetail{}
	keyTypeStr := pkd.KeyType.String()
	cols, dests := bstorage.SplitColDests(0, []*bstorage.ColDest{
		{publicKeyCol, &pkd.PublicKey},
		{keyTypeCol, &keyTypeStr},
		{entityIDCol, &pkd.EntityId},
	})
	return cols, dests, func() *api.PublicKeyDetail {
		pkd.PublicKey = *dests[0].(*[]byte)
		pkd.KeyType = api.KeyType(api.KeyType_value[*dests[1].(*string)])
		pkd.EntityId = *dests[2].(*string)
		return pkd
	}
}

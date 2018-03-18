package storage

import (
	"encoding/hex"
	"sync"

	api "github.com/elxirhealth/key/pkg/keyapi"
	"go.uber.org/zap"
)

type memoryStorer struct {
	pkds   map[string]*api.PublicKeyDetail
	mu     sync.Mutex
	params *Parameters
	logger *zap.Logger
}

// NewMemoryStorer creates a new Storer backed by an in-memory map.
func NewMemoryStorer(params *Parameters, logger *zap.Logger) Storer {
	return &memoryStorer{
		pkds:   make(map[string]*api.PublicKeyDetail),
		params: params,
		logger: logger,
	}
}

func (s *memoryStorer) AddPublicKeys(pkds []*api.PublicKeyDetail) error {
	if err := api.ValidatePublicKeyDetails(pkds); err != nil {
		return err
	}
	for _, pkd := range pkds {
		pkHex := hex.EncodeToString(pkd.PublicKey)
		s.mu.Lock()
		s.pkds[pkHex] = pkd
		s.mu.Unlock()
	}
	s.logger.Debug("added public keys to storage", zap.Int(logNPublicKeys, len(pkds)))
	return nil
}

func (s *memoryStorer) GetPublicKeys(pks [][]byte) ([]*api.PublicKeyDetail, error) {
	if err := api.ValidatePublicKeys(pks); err != nil {
		return nil, err
	}
	pkds := make([]*api.PublicKeyDetail, 0, len(pks))
	for _, pk := range pks {
		pkHex := hex.EncodeToString(pk)
		s.mu.Lock()
		pkd, in := s.pkds[pkHex]
		if !in {
			s.mu.Unlock()
			return nil, api.ErrNoSuchPublicKey
		}
		pkds = append(pkds, pkd)
		s.mu.Unlock()
	}
	s.logger.Debug("got public keys from storage", zap.Int(logNPublicKeys, len(pkds)))
	return pkds, nil
}

func (s *memoryStorer) GetEntityPublicKeys(entityID string) ([]*api.PublicKeyDetail, error) {
	if entityID == "" {
		return nil, api.ErrEmptyEntityID
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	pkds := make([]*api.PublicKeyDetail, 0, MaxEntityKeyTypeKeys)
	for _, pkd := range s.pkds {
		if pkd.EntityId == entityID && pkd.KeyType == api.KeyType_READER {
			pkds = append(pkds, pkd)
		}
	}
	s.logger.Debug("found public keys for entity", logGetEntityPubKeys(entityID, pkds)...)
	return pkds, nil
}

func (s *memoryStorer) CountEntityPublicKeys(entityID string, kt api.KeyType) (int, error) {
	if entityID == "" {
		return 0, api.ErrEmptyEntityID
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	c := 0
	for _, pkd := range s.pkds {
		if pkd.EntityId == entityID && pkd.KeyType == kt {
			c++
		}
	}
	s.logger.Debug("counted public keys for entity", logCountEntityPubKeys(entityID, kt)...)
	return c, nil
}

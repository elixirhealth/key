package datastore

import (
	api "github.com/elixirhealth/key/pkg/keyapi"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logNPublicKeys = "n_public_keys"
	logEntityID    = "entity_id"
	logKeyType     = "key_type"
)

func logGetEntityPubKeys(entityID string, pkds []*api.PublicKeyDetail) []zapcore.Field {
	return []zapcore.Field{
		zap.String(logEntityID, entityID),
		zap.Int(logNPublicKeys, len(pkds)),
	}
}
func logCountEntityPubKeys(entityID string, kt api.KeyType) []zapcore.Field {
	return []zapcore.Field{
		zap.String(logEntityID, entityID),
		zap.Stringer(logKeyType, kt),
	}
}

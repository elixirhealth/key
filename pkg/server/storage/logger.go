package storage

import (
	api "github.com/elxirhealth/key/pkg/keyapi"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logType            = "type"
	logNPublicKeys     = "n_public_keys"
	logAddQueryTimeout = "add_query_timeout"
	logGetQueryTimeout = "get_query_timeout"
	logEntityID        = "entity_id"
	logKeyType         = "key_type"
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

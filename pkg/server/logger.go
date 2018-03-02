package server

import (
	"github.com/elxirhealth/key/pkg/keyapi"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logStorage  = "storage"
	logEntityID = "entity_id"
	logKeyType  = "key_type"
	logNKeys    = "n_keys"
)

func logAddPublicKeysRq(rq *keyapi.AddPublicKeysRequest) []zapcore.Field {
	return []zapcore.Field{
		zap.String(logEntityID, rq.EntityId),
		zap.Stringer(logKeyType, rq.KeyType),
		zap.Int(logNKeys, len(rq.PublicKeys)),
	}
}

package server

import (
	api "github.com/elxirhealth/key/pkg/keyapi"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logStorage            = "storage"
	logEntityID           = "entity_id"
	logKeyType            = "key_type"
	logNKeys              = "n_keys"
	logOfEntityID         = "of_entity_id"
	logRequersterEntityID = "requester_entity_id"
	logNPublicKeys        = "n_public_keys"
)

func logAddPublicKeysRq(rq *api.AddPublicKeysRequest) []zapcore.Field {
	return []zapcore.Field{
		zap.String(logEntityID, rq.EntityId),
		zap.Stringer(logKeyType, rq.KeyType),
		zap.Int(logNKeys, len(rq.PublicKeys)),
	}
}

func logSamplePublicKeysRq(rq *api.SamplePublicKeysRequest) []zapcore.Field {
	return []zapcore.Field{
		zap.String(logOfEntityID, rq.OfEntityId),
		zap.String(logRequersterEntityID, rq.RequesterEntityId),
		zap.Uint32(logNPublicKeys, rq.NPublicKeys),
	}
}

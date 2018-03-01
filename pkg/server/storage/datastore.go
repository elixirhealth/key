package storage

import (
	api "github.com/elxirhealth/key/pkg/keyapi"
	bstorage "github.com/elxirhealth/service-base/pkg/server/storage"
)

type PublicKeyDetails struct {
	PublicKey string `datastore:"__key__"`
	EntityID  string `datastore:"entity_id"`
	KeyType   string `datastore:"key_type"`
	Disabled  bool   `datastore:"disabled"`
}

type datastoreStorer struct {
	client bstorage.DatastoreClient
}

func (s *datastoreStorer) AddPublicKeys(details []*api.PublicKeyDetails) error {
	// validate details (unique pub keys, each detail valid)
	// convert api.PKD to PKD
	// client.PutMulti
	panic("not implemented")
}

func (s *datastoreStorer) GetPublicKeys(pubKeys [][]byte) ([]*api.PublicKeyDetails, error) {
	// validate pubKeys
	// convert pubKeys to hex strings
	// client.GetMulti
	panic("not implemented")
}

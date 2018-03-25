package client

import (
	api "github.com/elixirhealth/key/pkg/keyapi"
	"google.golang.org/grpc"
)

// NewInsecure returns a new KeyClient without any TLS on the connection.
func NewInsecure(address string) (api.KeyClient, error) {
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return api.NewKeyClient(cc), nil
}

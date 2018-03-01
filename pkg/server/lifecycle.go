package server

import (
	api "github.com/elxirhealth/key/pkg/keyapi"
	"google.golang.org/grpc"
)

// Start starts the server and eviction routines.
func Start(config *Config, up chan *Key) error {
	c, err := newKey(config)
	if err != nil {
		return err
	}

	// start Key aux routines
	// TODO add go x.auxRoutine() or delete comment

	registerServer := func(s *grpc.Server) { api.RegisterKeyServer(s, c) }
	return c.Serve(registerServer, func() { up <- c })
}

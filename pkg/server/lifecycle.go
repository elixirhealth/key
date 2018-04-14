package server

import (
	"github.com/drausin/libri/libri/common/errors"
	api "github.com/elixirhealth/key/pkg/keyapi"
	"github.com/elixirhealth/key/pkg/server/storage/postgres/migrations"
	bstorage "github.com/elixirhealth/service-base/pkg/server/storage"
	"github.com/mattes/migrate/source/go-bindata"
	"google.golang.org/grpc"
)

// Start starts the server and eviction routines.
func Start(config *Config, up chan *Key) error {
	c, err := newKey(config)
	if err != nil {
		return err
	}

	if err := c.maybeMigrateDB(); err != nil {
		return err
	}

	registerServer := func(s *grpc.Server) { api.RegisterKeyServer(s, c) }
	return c.Serve(registerServer, func() { up <- c })
}

// StopServer handles cleanup involved in closing down the server.
func (k *Key) StopServer() {
	k.BaseServer.StopServer()
	err := k.storer.Close()
	errors.MaybePanic(err)
}

func (k *Key) maybeMigrateDB() error {
	if k.config.Storage.Type != bstorage.Postgres {
		return nil
	}

	m := bstorage.NewBindataMigrator(
		k.config.DBUrl,
		bindata.Resource(migrations.AssetNames(), migrations.Asset),
		&bstorage.ZapLogger{Logger: k.Logger},
	)
	return m.Up()
}

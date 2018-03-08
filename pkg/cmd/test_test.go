package cmd

import (
	"fmt"
	"sync"
	"testing"

	"github.com/elxirhealth/key/pkg/server"
	bcmd "github.com/elxirhealth/service-base/pkg/cmd"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestTestIO(t *testing.T) {
	// start in-memory key
	config := server.NewDefaultConfig()
	config.LogLevel = zapcore.DebugLevel
	config.ServerPort = 10200
	config.MetricsPort = 10201

	up := make(chan *server.Key, 1)
	wg1 := new(sync.WaitGroup)
	wg1.Add(1)
	go func(wg2 *sync.WaitGroup) {
		defer wg2.Done()
		err := server.Start(config, up)
		assert.Nil(t, err)
	}(wg1)

	x := <-up
	viper.Set(bcmd.AddressesFlag, fmt.Sprintf("localhost:%d", config.ServerPort))

	err := testIO()
	assert.Nil(t, err)

	x.StopServer()
	wg1.Wait()
}

package cmd

import (
	"log"

	"github.com/drausin/libri/libri/common/parse"
	"github.com/elxirhealth/key/pkg/keyapi"
	"github.com/elxirhealth/service-base/pkg/cmd"
	"github.com/elxirhealth/service-base/pkg/server"
	"github.com/spf13/viper"
)

func testIO() error {
	//rng := rand.New(rand.NewSource(0))
	//logger := lserver.NewDevLogger(lserver.GetLogLevel(viper.GetString(logLevelFlag)))
	//timeout := time.Duration(viper.GetInt(timeoutFlag) * 1e9)
	// TODO get other I/O params here

	//clients, err := getClients()
	_, err := getClients()
	if err != nil {
		return err
	}

	// TODO add I/O logic here
	log.Println("here to fool linter, delete this line when this function is fleshed out")

	return nil
}

func getClients() ([]keyapi.KeyClient, error) {
	addrs, err := parse.Addrs(viper.GetStringSlice(cmd.AddressesFlag))
	if err != nil {
		return nil, err
	}
	dialer := server.NewInsecureDialer()
	clients := make([]keyapi.KeyClient, len(addrs))
	for i, addr := range addrs {
		conn, err2 := dialer.Dial(addr.String())
		if err != nil {
			return nil, err2
		}
		clients[i] = keyapi.NewKeyClient(conn)
	}
	return clients, nil
}

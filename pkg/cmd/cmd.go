package cmd

import (
	"errors"
	"log"

	errors2 "github.com/drausin/libri/libri/common/errors"
	"github.com/drausin/libri/libri/common/logging"
	"github.com/elxirhealth/key/pkg/server"
	"github.com/elxirhealth/key/version"
	"github.com/elxirhealth/service-base/pkg/cmd"
	bserver "github.com/elxirhealth/service-base/pkg/server"
	bstorage "github.com/elxirhealth/service-base/pkg/server/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	serviceNameLower     = "key"
	serviceNameCamel     = "Key"
	envVarPrefix         = "KEY"
	logLevelFlag         = "logLevel"
	storageMemoryFlag    = "storageMemory"
	storageDataStoreFlag = "storageDataStore"
	gcpProjectIDFlag     = "gcpProjectID"
)

var (
	errMultipleStorageTypes = errors.New("multiple storage types specified")
	errNoStorageType        = errors.New("no storage type specified")

	rootCmd = &cobra.Command{
		Short: "operate a Key server",
	}
)

func init() {
	rootCmd.PersistentFlags().String(logLevelFlag, bserver.DefaultLogLevel.String(),
		"log level")

	cmd.Start(serviceNameLower, serviceNameCamel, rootCmd, version.Current, start,
		func(flags *pflag.FlagSet) {
			flags.Bool(storageMemoryFlag, true, "use in-memory storage")
			flags.Bool(storageDataStoreFlag, false, "use GCP DataStore storage")
			flags.String(gcpProjectIDFlag, "", "GCP project ID")
		})

	testCmd := cmd.Test(serviceNameLower, rootCmd)
	cmd.TestHealth(serviceNameLower, testCmd)
	cmd.TestIO(serviceNameLower, testCmd, testIO, func(flags *pflag.FlagSet) {
		// add additional test flags here if needed
	})

	cmd.Version(serviceNameLower, rootCmd, version.Current)

	// bind viper flags
	viper.SetEnvPrefix(envVarPrefix) // look for env vars with prefix
	viper.AutomaticEnv()             // read in environment variables that match
	errors2.MaybePanic(viper.BindPFlags(rootCmd.Flags()))
}

// Execute runs the root key command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func start() error {
	config, err := getKeyConfig()
	if err != nil {
		return err
	}
	return server.Start(config, make(chan *server.Key, 1))
}

func getKeyConfig() (*server.Config, error) {
	c := server.NewDefaultConfig()
	c.WithServerPort(uint(viper.GetInt(cmd.ServerPortFlag))).
		WithMetricsPort(uint(viper.GetInt(cmd.MetricsPortFlag))).
		WithProfilerPort(uint(viper.GetInt(cmd.ProfilerPortFlag))).
		WithLogLevel(logging.GetLogLevel(viper.GetString(logLevelFlag))).
		WithProfile(viper.GetBool(cmd.ProfileFlag))
	st, err := getStorageType()
	if err != nil {
		return nil, err
	}
	c.Storage.Type = st
	c.GCPProjectID = viper.GetString(gcpProjectIDFlag)
	return c, nil
}

func getStorageType() (bstorage.Type, error) {
	if viper.GetBool(storageMemoryFlag) && viper.GetBool(storageDataStoreFlag) {
		return bstorage.Unspecified, errMultipleStorageTypes
	}
	if viper.GetBool(storageMemoryFlag) {
		return bstorage.Memory, nil
	}
	if viper.GetBool(storageDataStoreFlag) {
		return bstorage.DataStore, nil
	}
	return bstorage.Unspecified, errNoStorageType
}

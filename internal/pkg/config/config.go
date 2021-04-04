package config

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

const (
	configPrefixService = "SVCUSER"
	configPrefixDB      = "DB"

	DefListenAddress  = ":8080"
	DefLogLevel       = log.WarnLevel
	DefHTTPTimeoutSec = 60
)

// AppConfig is a structure for the Application related ENV data.
type AppConfig struct {
	ServiceAddress string
	Loglevel       log.Level
	HTTPTimeoutSec int
}

// DBConfig is a structure for DB related ENV data.
type DBConfig struct {
	Address  string
	Password string
}

// Load binds required values from the OS environment to the expected structures.
func Load() (*AppConfig, *DBConfig, error) {
	appConf := AppConfig{
		ServiceAddress: DefListenAddress,
		Loglevel:       DefLogLevel,
		HTTPTimeoutSec: DefHTTPTimeoutSec,
	}

	err := envconfig.Process(configPrefixService, &appConf)
	if err != nil {
		return nil, nil, err
	}

	dbConf := DBConfig{}
	err = envconfig.Process(configPrefixDB, &dbConf)
	if err != nil {
		return &appConf, nil, err
	}

	return &appConf, &dbConf, nil
}

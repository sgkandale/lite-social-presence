package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configFile = "config.yaml"
)

func ParseConfig() *Config {
	log.Println("[INFO] parsing config file : ", configFile)

	var readConfig FlatConfig

	err := cleanenv.ReadConfig(configFile, &readConfig)
	if err != nil {
		log.Fatal("[ERROR] reading config.yaml file : ", err.Error())
	}

	log.Println("[INFO] config file parsed successfully")

	return &Config{
		Server: ServerConfig{
			Port:        readConfig.ServerPort,
			TLS:         readConfig.ServerTLS,
			CertPath:    readConfig.ServerCertPath,
			KeyPath:     readConfig.ServerKeyPath,
			ServiceName: readConfig.ServerServiceName,
		},
		Database: DatabaseConfig{
			Type:      readConfig.DatabaseType,
			UriString: readConfig.DatabaseUriString,
			Timeout:   readConfig.DatabaseTimeout,
		},
		Cache: CacheConfig{
			Type: readConfig.CacheType,
		},
	}
}

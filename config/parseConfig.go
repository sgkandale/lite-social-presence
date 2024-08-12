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

	var readConfig *Config

	err := cleanenv.ReadConfig(configFile, &readConfig)
	if err != nil {
		log.Fatal("[ERROR] reading config.yaml file")
	}

	log.Println("[INFO] config file parsed successfully")

	return readConfig
}

package config

import (
	"log"
)

func Verify(cfg *Config) {
	// Server checks
	if cfg.Server.Port <= 0 {
		log.Fatal("[ERROR] server_port is empty in config")
	}
	if cfg.Server.TLS {
		log.Print("[INFO] server_tls is true in config")
		if cfg.Server.CertPath == "" {
			log.Fatal("[ERROR] server_cert_path is empty in config")
		}
		if cfg.Server.KeyPath == "" {
			log.Fatal("[ERROR] server_key_path is empty in config")
		}
	}
	if cfg.Server.ServiceName == "" {
		log.Fatal("[ERROR] server_service_name is empty in config")
	}

	// database checks
	if cfg.Database.Type == "" {
		log.Fatal("[ERROR] database_type is empty in config")
	}
	if cfg.Database.UriString == "" {
		log.Fatal("[ERROR] database_uri_string is empty in config")
	}
	if cfg.Database.Timeout <= 0 {
		log.Fatal("[ERROR] database_timeout is empty in config")
	}

	// cache checks
	if cfg.Cache.Type == "" {
		log.Fatal("[ERROR] cache_type is empty in config")
	}
}

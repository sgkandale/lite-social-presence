package config

import (
	"log"
)

func Verify(cfg *Config) {
	// Server checks
	if cfg.Server.Port <= 0 {
		log.Fatal("[ERROR] server.port is empty in config")
	}
	if cfg.Server.TLS {
		log.Print("[INFO] server.tls is true in config")
		if cfg.Server.CertPath == "" {
			log.Fatal("[ERROR] server.cert_path is empty in config")
		}
		if cfg.Server.KeyPath == "" {
			log.Fatal("[ERROR] server.key_path is empty in config")
		}
	}
	if cfg.Server.ServiceName == "" {
		log.Fatal("[ERROR] server.service_name is empty in config")
	}
}

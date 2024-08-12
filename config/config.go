package config

type ServerConfig struct {
	Port        int    `yaml:"port" env:"port"`
	TLS         bool   `yaml:"tls" env:"tls"`
	CertPath    string `yaml:"cert_path" env:"cert_path"`
	KeyPath     string `yaml:"key_path" env:"key_path"`
	ServiceName string `yaml:"service_name" env:"service_name"`
}

type Config struct {
	Server ServerConfig `yaml:"server" env:"server"`
}

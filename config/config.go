package config

type ServerConfig struct {
	Port        int    `yaml:"port" env:"port"`
	TLS         bool   `yaml:"tls" env:"tls"`
	CertPath    string `yaml:"cert_path" env:"cert_path"`
	KeyPath     string `yaml:"key_path" env:"key_path"`
	ServiceName string `yaml:"service_name" env:"service_name"`
}

type DatabaseConfig struct {
	Type      string `yaml:"type" env:"type"`
	UriString string `yaml:"uri_string" env:"uri_string"`
	Timeout   int    `yaml:"timeout" env:"timeout"`
}

type CacheConfig struct {
	Type string `yaml:"type" env:"type"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server" env:"server"`
	Database DatabaseConfig `yaml:"database" env:"database"`
	Cache    CacheConfig    `yaml:"cache" env:"cache"`
}

type FlatConfig struct {
	ServerPort        int    `yaml:"server_port" env:"server_port"`
	ServerTLS         bool   `yaml:"server_tls" env:"server_tls"`
	ServerCertPath    string `yaml:"server_cert_path" env:"server_cert_path"`
	ServerKeyPath     string `yaml:"server_key_path" env:"server_key_path"`
	ServerServiceName string `yaml:"server_service_name" env:"server_service_name"`

	DatabaseType      string `yaml:"database_type" env:"database_type"`
	DatabaseUriString string `yaml:"database_uri_string" env:"database_uri_string"`
	DatabaseTimeout   int    `yaml:"database_timeout" env:"database_timeout"`

	CacheType string `yaml:"cache_type" env:"cache_type"`
}

package config

// Config - конфиг сервиса
type Config struct {
	HTTP HTTPConfig
}

// HTTPConfig - конфиг http сервера
type HTTPConfig struct {
	Host string `default:"localhost"`
	Port int    `default:"8080"`
}

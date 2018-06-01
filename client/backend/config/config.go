package config

// Config - конфиг сервиса
type Config struct {
	HTTP HTTPConfig
}

// HTTPConfig - конфиг http сервера
type HTTPConfig struct {
	Host string `default:"0.0.0.0"`
	Port int    `default:"8080"`
}

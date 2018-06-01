package config

// Config - определение структуры конфига
type Config struct {
	Log    LogConfig
	Client User
	Server Server
}

// LogConfig - конфигурация логгера
type LogConfig struct {
	Writer   string
	MinLevel string `env:"LOG_LEVEL"`
}

// User - конфигурация клиентских данных
type User struct {
	ID    int
	Token string
}

// Server - конфигурация сервера
type Server struct {
	Address string
}

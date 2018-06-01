package config

// Config - определение структуры конфига
type Config struct {
	Log LogConfig
	Db  DbConnection
}

// LogConfig конфигурация логгера
type LogConfig struct {
	Writer   string
	MinLevel string `env:"LOG_LEVEL"`
}

// DbConnection - конфиг коннекта к базе данных
type DbConnection struct {
	Type     string
	User     string
	Password string
	DbName   string
	Host     string
	Port     int
}

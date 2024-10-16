package config

// PgConfig Конфигурация для подключения к Postgres
type PgConfig struct {
	Dsn string `env:"PG_DSN"`
}

package config

import "os"

const (
	PG_DSN           = "PG_DSN"
	TELEGRAM_API_KEY = "TELEGRAM_API_KEY"
)

type Config struct {
	PostgresDSN    string
	TelegramAPIKey string
}

func GetConfig() *Config {
	pgDsn := os.Getenv(PG_DSN)
	tgAPIKey := os.Getenv(TELEGRAM_API_KEY)

	cfg := &Config{
		PostgresDSN:    pgDsn,
		TelegramAPIKey: tgAPIKey,
	}

	return cfg
}

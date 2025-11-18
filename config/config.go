package config

import (
	"fmt"
	"time"
)

type Config struct {
	Postgres *Postgres `env:",prefix=POSTGRES_"`
	Telegram *Telegram `env:",prefix=TELEGRAM_"`
}

type Telegram struct {
	Token string `env:"TOKEN" envDefault:""`
}

type Postgres struct {
	Env                   string        `env:"ENV" envDefault:"dev"`
	Host                  string        `env:"HOST" envDefault:"127.0.0.1"`
	Port                  int           `env:"PORT" envDefault:"5432"`
	User                  string        `env:"USER" envDefault:"postgres"`
	Password              string        `env:"PASSWORD" envDefault:"postgres"`
	Database              string        `env:"DATABASE" envDefault:""`
	SSLMode               string        `env:"SSL_MODE" envDefault:"disable"`
	MaxIdleConnections    int           `env:"MAX_IDLE_CONNECTIONS" envDefault:"25"`
	MaxOpenConnections    int           `env:"MAX_OPEN_CONNECTIONS" envDefault:"25"`
	ConnectionMaxLifetime time.Duration `env:"CONNECTION_MAX_LIFETIME" envDefault:"5m"`
}

func (p *Postgres) ConnectionURL() string {
	if p.User == "" {
		return fmt.Sprintf("host=%s port=%d dbname=%s sslmode=%s",
			p.Host, p.Port, p.Database, p.SSLMode)
	}
	if p.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s",
			p.Host, p.Port, p.User, p.Database, p.SSLMode)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.Database, p.SSLMode)
}

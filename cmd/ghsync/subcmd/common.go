package subcmd

import "fmt"

type PostgresOpt struct {
	DB       string `long:"postgres-db" env:"GHSYNC_POSTGRES_DB" description:"PostgreSQL DB" default:"ghsync"`
	User     string `long:"postgres-user" env:"GHSYNC_POSTGRES_USER" description:"PostgreSQL user" default:"superset"`
	Password string `long:"postgres-password" env:"GHSYNC_POSTGRES_PASSWORD" description:"PostgreSQL password" default:"superset"`
	Host     string `long:"postgres-host" env:"GHSYNC_POSTGRES_HOST" description:"PostgreSQL host" default:"localhost"`
	Port     int    `long:"postgres-port" env:"GHSYNC_POSTGRES_PORT" description:"PostgreSQL port" default:"5432"`
}

func (o PostgresOpt) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		o.User, o.Password, o.Host, o.Port, o.DB)
}

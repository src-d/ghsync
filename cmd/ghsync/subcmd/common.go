package subcmd

import (
	"fmt"

	"github.com/src-d/ghsync/models/migrations"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

const maxVersion uint = 1558054487

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

func newMigrate(url string) (*migrate.Migrate, error) {
	// wrap assets into Resource
	s := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})

	d, err := bindata.WithInstance(s)
	if err != nil {
		return nil, err
	}
	return migrate.NewWithSourceInstance("go-bindata", d, url)
}

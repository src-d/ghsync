package subcmd

import (
	"github.com/src-d/ghsync/models/migrations"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"

	gocli "gopkg.in/src-d/go-cli.v0"
	"gopkg.in/src-d/go-log.v1"
)

type MigrateCommand struct {
	gocli.PlainCommand `name:"migrate" short-description:"performs a DB migration up to the latest version" long-description:"Performs a DB migration up to the latest version"`
	Postgres           PostgresOpt `group:"PostgreSQL connection options"`
}

func (c *MigrateCommand) Execute(args []string) error {
	// wrap assets into Resource
	s := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})

	d, err := bindata.WithInstance(s)
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("go-bindata", d, c.Postgres.URL())
	if err != nil {
		return err
	}

	err = m.Up()
	switch err {
	case nil:
		log.Infof("The DB was upgraded")
	case migrate.ErrNoChange:
		log.Infof("The DB is up to date")
	default:
		return err
	}

	return nil
}

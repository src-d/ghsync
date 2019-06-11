package subcmd

import (
	"github.com/golang-migrate/migrate/v4"

	gocli "gopkg.in/src-d/go-cli.v0"
	"gopkg.in/src-d/go-log.v1"
)

type MigrateCommand struct {
	gocli.PlainCommand `name:"migrate" short-description:"performs a DB migration up to the latest version" long-description:"Performs a DB migration up to the latest version"`
	Postgres           PostgresOpt `group:"PostgreSQL connection options"`
}

func (c *MigrateCommand) Execute(args []string) error {
	m, err := newMigrate(c.Postgres.URL())
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

package subcmd

import (
	"github.com/src-d/ghsync/shallow"

	"gopkg.in/src-d/go-cli.v0"
)

type ShallowCommand struct {
	cli.Command `name:"shallow" short-description:"Shallow sync of GitHub data" long-description:"Shallow sync of GitHub data"`

	Token string `long:"token" env:"GHSYNC_TOKEN" description:"GitHub personal access token" required:"true"`
	Org   string `long:"org" env:"GHSYNC_ORG" description:"Name of the GitHub organization" required:"true"`

	Postgres PostgresOpt `group:"PostgreSQL connection options"`
}

func (c *ShallowCommand) Execute(args []string) error {
	db, err := c.Postgres.initDB()
	if err != nil {
		return err
	}
	defer db.Close()

	client, err := newClient(c.Token)
	if err != nil {
		return err
	}

	orgSyncer := shallow.NewOrganizationSyncer(db, client)
	return orgSyncer.Sync(c.Org)
}

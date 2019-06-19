package subcmd

import (
	"strings"

	"github.com/src-d/ghsync/shallow"

	"gopkg.in/src-d/go-cli.v0"
)

type ShallowCommand struct {
	cli.Command `name:"shallow" short-description:"Shallow sync of GitHub data" long-description:"Shallow sync of GitHub data"`

	Token string `long:"token" env:"GHSYNC_TOKEN" description:"GitHub personal access token" required:"true"`
	Orgs  string `long:"orgs" env:"GHSYNC_ORGS" description:"Comma-separated list of GitHub organization names" required:"true"`

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
	for _, o := range strings.Split(c.Orgs, ",") {
		err = orgSyncer.Sync(o)
		if err != nil {
			return err
		}
	}

	return nil
}

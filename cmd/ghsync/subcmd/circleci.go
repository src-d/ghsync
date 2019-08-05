package subcmd

import (
	client "github.com/jszwedko/go-circleci"
	"github.com/src-d/ghsync/circleci"
	"gopkg.in/src-d/go-cli.v0"
)

type CircleCICommand struct {
	cli.Command `name:"circleci" short-description:"Shallow sync of GitHub data" long-description:"Shallow sync of GitHub data"`

	Token string `long:"token" env:"GHSYNC_TOKEN" description:"GitHub personal access token" required:"true"`
	Orgs  string `long:"orgs" env:"GHSYNC_ORGS" description:"Comma-separated list of GitHub organization names" required:"true"`

	Postgres PostgresOpt `group:"PostgreSQL connection options"`
}

func (c *CircleCICommand) Execute(args []string) error {
	db, err := c.Postgres.initDB()
	if err != nil {
		return err
	}
	defer db.Close()

	client := &client.Client{} // Token not rquired to query info for public projects

	syncer := circleci.NewBuildSyncer(db, client)
	return syncer.Sync("pytorch", "pytorch")
}

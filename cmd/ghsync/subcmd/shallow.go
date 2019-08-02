package subcmd

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/src-d/ghsync/shallow"

	"gopkg.in/src-d/go-cli.v0"
	"gopkg.in/src-d/go-log.v1"
)

type ShallowCommand struct {
	cli.Command `name:"shallow" short-description:"Shallow sync of GitHub data" long-description:"Shallow sync of GitHub data"`

	Token string `long:"token" env:"GHSYNC_TOKEN" description:"GitHub personal access token" required:"true"`
	Orgs  string `long:"orgs" env:"GHSYNC_ORGS" description:"Comma-separated list of GitHub organization names" required:"true"`

	NoForks bool `long:"no-forks"  env:"GHSYNC_NO_FORKS" description:"github forked repositories will be skipped"`

	Postgres PostgresOpt `group:"PostgreSQL connection options"`
}

func (c *ShallowCommand) Execute(args []string) error {
	if c.Orgs == "" {
		log.Warningf("no organizations found, at least one " +
			"organization must be provided")

		return nil
	}

	db, err := c.Postgres.initDB()
	if err != nil {
		return err
	}
	defer db.Close()

	client, err := newClient(c.Token)
	if err != nil {
		return err
	}

	orgs := strings.Split(c.Orgs, ",")
	if err = c.initStatus(db, statusTableName, orgs); err != nil {
		return err
	}

	orgSyncer := shallow.NewOrganizationSyncer(db, client, statusTableName, c.NoForks)
	for _, o := range orgs {
		err = orgSyncer.Sync(o)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ShallowCommand) initStatus(db *sql.DB, tableName string, orgs []string) error {
	log.Debugf("initializing status table for orgs: %v", orgs)
	var b strings.Builder

	for _, o := range orgs[:len(orgs)-1] {
		b.WriteString(fmt.Sprintf("('%s', 'repository'),", o))
		b.WriteString(fmt.Sprintf("('%s', 'user'),", o))
	}
	b.WriteString(fmt.Sprintf("('%s', 'repository'),", orgs[len(orgs)-1]))
	b.WriteString(fmt.Sprintf("('%s', 'user')", orgs[len(orgs)-1]))

	stm := fmt.Sprintf("INSERT INTO %s (org, entity) VALUES %s ON CONFLICT (org, entity) DO UPDATE SET failed=0, done=0, total=NULL;", tableName, b.String())
	log.Debugf("running statement: %s", stm)
	_, err := db.Exec(stm)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf(
			"an error occured while initializing %s table: %v", tableName, err))
	}

	return nil
}

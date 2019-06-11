package subcmd

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/src-d/ghsync"
	"github.com/src-d/ghsync/utils"

	"github.com/golang-migrate/migrate/v4"
	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-cli.v0"
	"gopkg.in/src-d/go-log.v1"
	"gopkg.in/src-d/go-queue.v1"
	_ "gopkg.in/src-d/go-queue.v1/amqp"
	_ "gopkg.in/src-d/go-queue.v1/memory"
)

type SyncCommand struct {
	cli.Command `name:"sync"`

	Token string `long:"token" env:"GHSYNC_TOKEN" description:"GitHub personal access token" required:"true"`
	Org   string `long:"org" env:"GHSYNC_ORG" description:"Name of the GitHub organization" required:"true"`

	QueueOpt struct {
		Queue  string `long:"queue" env:"GHSYNC_QUEUE" description:"queue name. If it's not set the organization name will be used"`
		Broker string `long:"broker" env:"GHSYNC_BROKER" default:"amqp://localhost:5672" description:"broker service URI"`
	} `group:"go-queue connection options"`

	Postgres PostgresOpt `group:"PostgreSQL connection options"`
}

func (c *SyncCommand) Execute(args []string) error {
	db, err := c.initDB()
	if err != nil {
		return err
	}
	defer db.Close()

	http := oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Token},
	))

	t := httpcache.NewTransport(diskcache.New("cache/" + c.Org))
	t.Transport = &RemoveHeaderTransport{utils.NewRateLimitTransport(http.Transport)}
	http.Transport = t

	client := github.NewClient(http)

	broker, err := queue.NewBroker(c.QueueOpt.Broker)
	if err != nil {
		return err
	}

	qName := c.QueueOpt.Queue
	if qName == "" {
		qName = c.Org
	}
	queue, err := broker.Queue(qName)
	if err != nil {
		return err
	}

	syncer := ghsync.NewSyncer(db, client, queue)

	go func() {
		err := syncer.DoOrganization(c.Org)
		if err != nil {
			log.Errorf(err, "syncer.DoOrganization finished with error")
		}
	}()

	return syncer.Wait()
}

func (c *SyncCommand) initDB() (db *sql.DB, err error) {
	db, err = sql.Open("postgres", c.Postgres.URL())
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			db.Close()
			db = nil
		}
	}()

	if err = db.Ping(); err != nil {
		return db, err
	}

	m, err := newMigrate(c.Postgres.URL())
	if err != nil {
		return db, err
	}

	dbVersion, _, err := m.Version()

	if err != nil && err != migrate.ErrNilVersion {
		return db, err
	}

	if dbVersion != maxVersion {
		return db, fmt.Errorf(
			"database version mismatch. Current version is %v, but this binary needs version %v. "+
				"Use the 'migrate' subcommand to upgrade your database", dbVersion, maxVersion)
	}

	log.With(log.Fields{"db-version": dbVersion}).Debugf("the DB version is up to date")
	log.Infof("connection with the DB established")
	return db, nil
}

type RemoveHeaderTransport struct {
	T http.RoundTripper
}

func (t *RemoveHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Del("X-Ratelimit-Limit")
	req.Header.Del("X-Ratelimit-Remaining")
	req.Header.Del("X-Ratelimit-Reset")
	return t.T.RoundTrip(req)
}

package subcmd

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/src-d/ghsync/models/migrations"
	"github.com/src-d/ghsync/utils"
	"gopkg.in/src-d/go-log.v1"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"golang.org/x/oauth2"
)

const maxVersion uint = 1560510971
const statusTableName = "status"

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

func (o PostgresOpt) initDB() (db *sql.DB, err error) {
	db, err = sql.Open("postgres", o.URL())
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

	return db, nil

	m, err := newMigrate(o.URL())
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
	if err = o.createStatusTable(); err != nil {
		return db, err
	}

	return db, nil
}

func (o PostgresOpt) createStatusTable() error {
	log.Debugf(fmt.Sprintf("creating status table '%s'", statusTableName))

	db, err := sql.Open("postgres", o.URL())
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			db.Close()
		}
	}()

	stm := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
    id serial PRIMARY KEY,
    org VARCHAR (50) NOT NULL,
    entity VARCHAR (20) NOT NULL,
    done INTEGER NOT NULL DEFAULT 0,
    failed INTEGER NOT NULL DEFAULT 0,
    total INTEGER DEFAULT NULL,
    UNIQUE (org, entity)
);`, statusTableName)
	log.Debugf("running statement: %s", stm)
	_, err = db.Exec(stm)
	if err != nil {
		return fmt.Errorf("an error occured while ensureing the status table: %v", err)
	}

	log.Infof("status table '%s' created", statusTableName)

	return nil
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

func newClient(token string) (*github.Client, error) {
	http := oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	dirPath := filepath.Join(os.TempDir(), "ghsync")
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error while creating directory %s: %v", dirPath, err)
	}

	t := httpcache.NewTransport(diskcache.New(dirPath))
	t.Transport = &RemoveHeaderTransport{utils.NewRateLimitTransport(http.Transport)}
	http.Transport = &RetryTransport{T: t}

	return github.NewClient(http), nil
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

type RetryTransport struct {
	T http.RoundTripper
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var r *http.Response
	var err error
	utils.Retry(func() error {
		r, err = t.T.RoundTrip(req)
		return err
	})

	return r, err
}

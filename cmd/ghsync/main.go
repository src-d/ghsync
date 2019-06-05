package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"

	"github.com/src-d/ghsync"
	"github.com/src-d/ghsync/utils"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-cli.v0"
	"gopkg.in/src-d/go-queue.v1"
	_ "gopkg.in/src-d/go-queue.v1/amqp"
	_ "gopkg.in/src-d/go-queue.v1/memory"
)

// rewritten during the CI build step
var (
	version = "master"
	build   = "dev"
)

var app = cli.New("ghsync", version, build, "GitHub metadata sync")

func main() {
	app.AddCommand(&syncCommand{})

	app.RunMain()
}

type syncCommand struct {
	cli.Command `name:"sync"`
}

func (c *syncCommand) Execute(args []string) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "superset", "superset", "ghsync")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	org := os.Getenv("GITHUB_ORG")

	http := oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	))

	t := httpcache.NewTransport(diskcache.New("cache/" + org))
	t.Transport = &RemoveHeaderTransport{utils.NewRateLimitTransport(http.Transport)}
	http.Transport = t

	client := github.NewClient(http)

	broker, err := queue.NewBroker("amqp://localhost:5672")
	if err != nil {
		return err
	}

	queue, err := broker.Queue(org)
	if err != nil {
		return err
	}

	syncer := ghsync.NewSyncer(db, client, queue)
	go syncer.DoOrganization(org)
	fmt.Println(syncer.Wait())

	return nil
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

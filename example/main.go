package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"

	"github.com/src-d/ghsync"
	"github.com/src-d/ghsync/utils"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-queue.v1"
	_ "gopkg.in/src-d/go-queue.v1/amqp"
	_ "gopkg.in/src-d/go-queue.v1/memory"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "superset", "superset", "ghsync")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http := oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "4b56d692e95dd8bae64e9d1092cd955d0db49fc3"},
	))

	t := httpcache.NewTransport(diskcache.New("cache/src-d"))
	t.Transport = utils.NewRateLimitTransport(http.Transport)
	http.Transport = t

	client := github.NewClient(http)

	broker, _ := queue.NewBroker("amqp://localhost:5672")
	queue, _ := broker.Queue("src-d")

	syncer := ghsync.NewSyncer(db, client, queue)
	go syncer.DoOrganization("src-d")
	fmt.Println(syncer.Wait())
}

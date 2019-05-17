package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/src-d/ghsync"
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

	client := github.NewClient(oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "7317b6ae552baec41d1304d8ac13b58f92dce868"},
	)))

	broker, _ := queue.NewBroker("amqp://localhost:5672")
	queue, _ := broker.Queue("src-d")

	syncer := ghsync.NewSyncer(db, client, queue)
	go syncer.DoOrganization("src-d")
	fmt.Println(syncer.Wait())
}

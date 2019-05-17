package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/src-d/ghsync"
	"golang.org/x/oauth2"
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

	syncUser := ghsync.NewUserSyncer(db, client)
	fmt.Println(syncUser.Sync("mcuadros"))
	//syncUser.QueueOrganization("src-d")

	syncOrg := ghsync.NewOrganizationSyncer(db, client)
	fmt.Println(syncOrg.Sync("src-d"))

	syncRepository := ghsync.NewRepositorySyncer(db, client)
	fmt.Println(syncRepository.Sync("src-d", "go-git"))
	//syncRepository.QueueOrganization("mcuadros")

	syncIssue := ghsync.NewIssueSyncer(db, client)
	fmt.Println(syncIssue.Sync("src-d", "go-git", 785))
	//syncIssue.QueueRepository("src-d", "go-git")

	syncPullRequest := ghsync.NewPullRequestSyncer(db, client)
	fmt.Println(syncPullRequest.Sync("src-d", "go-git", 1045))
	//syncPullRequest.QueueRepository("src-d", "go-git")

}

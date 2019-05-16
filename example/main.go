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
		&oauth2.Token{AccessToken: "9f2deb4e3649116d926a292406667fc87964cd18"},
	)))

	syncUser := ghsync.NewUserSyncer(db, client)
	fmt.Println(syncUser.Sync("mcuadros"))

	syncOrg := ghsync.NewOrganizationSyncer(db, client)
	fmt.Println(syncOrg.Sync("src-d"))

	syncRepository := ghsync.NewRepositorySyncer(db, client)
	fmt.Println(syncRepository.Sync("src-d", "go-git"))

	syncIssue := ghsync.NewIssueSyncer(db, client)
	fmt.Println(syncIssue.Sync("src-d", "go-git", 785))

	syncPullRequest := ghsync.NewPullRequestSyncer(db, client)
	fmt.Println(syncPullRequest.Sync("src-d", "go-git", 1045))

}

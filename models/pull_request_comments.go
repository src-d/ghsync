package models

import (
	"github.com/google/go-github/github"
	"github.com/src-d/ghsync/utils"
	"gopkg.in/src-d/go-kallax.v1"
)

type PullRequestComment struct {
	kallax.Model `table:"pull_request_comments" pk:"kallax_id" ignored:"User,URL,PullRequestURL"`
	github.PullRequestComment

	// int64 replacement for PullRequestComment.ID *int64, to be used as primary key
	KallaxID int64 `kallax:"kallax_id"`

	UserID            int64  `kallax:"user_id"`
	UserLogin         string `kallax:"user_login"`
	PullRequestNumber int    `kallax:"pull_request_number"`
	RepositoryOwner   string `kallax:"repository_owner"`
	RepositoryName    string `kallax:"repository_name"`
}

func (i *PullRequestComment) BeforeSave() error {
	i.KallaxID = i.PullRequestComment.GetID()

	var err error
	i.RepositoryOwner, i.RepositoryName, i.PullRequestNumber, err = utils.ParsePullRequestURL(i.GetPullRequestURL())
	if err != nil {
		return err
	}

	if i.User != nil {
		i.UserID = i.User.GetID()
		i.UserLogin = i.User.GetLogin()
	}

	return nil
}

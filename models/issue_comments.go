package models

import (
	"github.com/google/go-github/github"
	"github.com/src-d/ghsync/utils"
	"gopkg.in/src-d/go-kallax.v1"
)

type IssueComment struct {
	kallax.Model `table:"issue_comments" pk:"kallax_id" ignored:"User,URL,IssueURL"`
	github.IssueComment

	UserID          int64  `kallax:"user_id"`
	UserLogin       string `kallax:"user_login"`
	IssueNumber     int    `kallax:"issue_number"`
	RepositoryOwner string `kallax:"repository_owner"`
	RepositoryName  string `kallax:"repository_name"`

	// int64 replacement for IssueComment.ID *int64, to be used as primary key
	KallaxID int64 `kallax:"kallax_id"`
}

func (i *IssueComment) BeforeSave() error {
	var err error
	i.RepositoryOwner, i.RepositoryName, i.IssueNumber, err = utils.ParseIssueURL(i.GetIssueURL())
	if err != nil {
		return err
	}

	if i.User != nil {
		i.UserID = i.User.GetID()
		i.UserLogin = i.User.GetLogin()
	}

	i.KallaxID = i.IssueComment.GetID()

	return nil
}

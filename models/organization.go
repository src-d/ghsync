package models

import (
	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
)

type Organization struct {
	kallax.Model `table:"organizations" pk:"kallax_id" ignored:"Plan,URL,EventsURL,HooksURL,IssuesURL,MembersURL,PublicMembersURL,ReposURL,DefaultRepoPermission,DefaultRepoSettings,MembersCanCreateRepos"`
	github.Organization

	// int64 replacement for Organization.ID *int64, to be used as primary key
	KallaxID int64 `kallax:"kallax_id"`
}

func (o *Organization) BeforeSave() error {
	o.KallaxID = o.Organization.GetID()

	return nil
}

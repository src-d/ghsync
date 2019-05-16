package models

import (
	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
)

type Organization struct {
	kallax.Model `table:"organizations" pk:"id" ignored:"Plan,URL,EventsURL,HooksURL,IssuesURL,MembersURL,PublicMembersURL,ReposURL,DefaultRepoPermission,DefaultRepoSettings,MembersCanCreateRepos"`
	github.Organization
}

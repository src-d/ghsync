package models

import (
	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
)

type Repository struct {
	kallax.Model `table:"repositories" pk:"id" ignored:"Owner,Organization,TextMatches"`
	github.Repository

	OwnerID          int64  `kallax:"owner_id"`
	OwnerLogin       string `kallax:"owner_login"`
	OrganizationID   int64  `kallax:"organization_id"`
	OrganizationName string `kallax:"organization_name"`
}

func (r *Repository) BeforeSave() error {
	if r.Owner != nil {
		r.OwnerID = r.Owner.GetID()
		r.OwnerLogin = r.Owner.GetLogin()
	}

	if r.Organization != nil {
		r.OrganizationID = r.Organization.GetID()
		r.OrganizationName = r.Organization.GetName()
	}

	return nil
}

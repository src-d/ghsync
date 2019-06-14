package models

import (
	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
)

type User struct {
	kallax.Model `table:"users" pk:"kallax_id" ignored:"Plan,URL,EventsURL,FollowingURL,FollowersURL,GistsURL,OrganizationsURL,,ReceivedEventsURL,ReposURL,StarredURL,SubscriptionsURL,TextMatches,Permissions"`
	github.User

	// int64 replacement for User.ID *int64, to be used as primary key
	KallaxID int64 `kallax:"kallax_id"`
}

func (u *User) BeforeSave() error {
	u.KallaxID = u.User.GetID()

	return nil
}

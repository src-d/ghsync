package models

import (
	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
)

type User struct {
	kallax.Model `table:"users" pk:"id" ignored:"Plan,URL,EventsURL,FollowingURL,FollowersURL,GistsURL,OrganizationsURL,,ReceivedEventsURL,ReposURL,StarredURL,SubscriptionsURL,TextMatches,Permissions"`
	github.User
}

package models

import (
	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
)

type Repository struct {
	kallax.Model `table:"repositories" pk:"kallax_id" ignored:"Owner,Parent,Source,Organization,URL,ArchiveURL,AssigneesURL,BlobsURL,BranchesURL,CollaboratorsURL,CommentsURL,CommitsURL,CompareURL,ContentsURL,ContributorsURL,DeploymentsURL,DownloadsURL,EventsURL,ForksURL,GitCommitsURL,GitRefsURL,GitTagsURL,HooksURL,IssueCommentURL,IssueEventsURL,IssuesURL,KeysURL,LabelsURL,LanguagesURL,MergesURL,MilestonesURL,NotificationsURL,PullsURL,ReleasesURL,StargazersURL,StatusesURL,SubscribersURL,SubscriptionURL,TagsURL,TreesURL,TeamsURL,TextMatches"`
	github.Repository

	// int64 replacement for Repository.ID *int64, to be used as primary key
	KallaxID int64 `kallax:"kallax_id"`

	ParentRepository *RepositoryReference `kallax:"parent"`
	SourceRepository *RepositoryReference `kallax:"source"`

	OwnerID    int64  `kallax:"owner_id"`
	OwnerType  string `kallax:"owner_type"`
	OwnerLogin string `kallax:"owner_login"`

	OrganizationID   int64  `kallax:"organization_id"`
	OrganizationName string `kallax:"organization_name"`
}

func (r *Repository) BeforeSave() error {
	r.KallaxID = r.Repository.GetID()

	r.ParentRepository = NewRepositoryReference(r.Parent)
	r.SourceRepository = NewRepositoryReference(r.Parent)

	if r.Owner != nil {
		r.OwnerID = r.Owner.GetID()
		r.OwnerLogin = r.Owner.GetLogin()
		r.OwnerType = r.Owner.GetType()
	}

	if r.Organization != nil {
		r.OrganizationID = r.Organization.GetID()
		r.OrganizationName = r.Organization.GetName()
	}

	return nil
}

package models

import (
	"github.com/google/go-github/github"
	"github.com/src-d/ghsync/utils"
	"gopkg.in/src-d/go-kallax.v1"
)

type Issue struct {
	kallax.Model `table:"issues" pk:"id" ignored:"Labels,User,Assignee,ClosedBy,Repository,Milestone,PullRequestLinks,Assignees,URL,CommentsURL,EventsURL,LabelsURL,RepositoryURL,Milestone,PullRequestLinks,Reactions,ActiveLockReason,TextMatches"`
	github.Issue

	RepositoryOwner string `kallax:"repository_owner"`
	RepositoryName  string `kallax:"repository_name"`

	LabelList []string `kallax:"labels"`

	UserID        int64            `kallax:"user_id"`
	UserLogin     string           `kallax:"user_login"`
	AssigneeID    int64            `kallax:"assignee_id"`
	AssigneeLogin string           `kallax:"assignee_login"`
	AssigneesList []*UserReference `kallax:"assignees"`
	ClosedByID    int64            `kallax:"closed_by_id"`
	ClosedByLogin string           `kallax:"closed_by_login"`

	MilestoneID    int64  `kallax:"milestone_id"`
	MilestoneTitle string `kallax:"milestone_title"`
}

func (i *Issue) BeforeSave() error {
	var err error
	i.RepositoryOwner, i.RepositoryName, _, err = utils.ParseIssueURL(i.GetURL())
	if err != nil {
		return err
	}

	i.AssigneesList = NewUserRefernceList(i.Assignees)

	i.LabelList = make([]string, 0)
	for _, l := range i.Labels {
		i.LabelList = append(i.LabelList, l.GetName())
	}

	if i.User != nil {
		i.UserID = i.User.GetID()
		i.UserLogin = i.User.GetLogin()
	}

	if i.Assignee != nil {
		i.AssigneeID = i.Assignee.GetID()
		i.AssigneeLogin = i.Assignee.GetLogin()
	}

	if i.ClosedBy != nil {
		i.ClosedByID = i.ClosedBy.GetID()
		i.ClosedByLogin = i.ClosedBy.GetLogin()
	}

	if i.Milestone != nil {
		i.MilestoneID = i.Milestone.GetID()
		i.MilestoneTitle = i.Milestone.GetTitle()
	}

	return nil
}

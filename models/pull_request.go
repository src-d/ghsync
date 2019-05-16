package models

import (
	"github.com/google/go-github/github"
	"github.com/src-d/ghsync/utils"
	"gopkg.in/src-d/go-kallax.v1"
)

type PullRequest struct {
	kallax.Model `table:"pull_requests" pk:"id" ignored:"Labels,User,MergedBy,Assignee,Milestone,Assignees,Assignees  ,RequestedReviewers,Links,Head,Base,ActiveLockReason,RequestedTeams,URL,IssueURL,StatusesURL,DiffURL,PatchURL,CommitsURL,CommentsURL,ReviewCommentsURL,ReviewCommentURL"`
	github.PullRequest

	RepositoryOwner string `kallax:"repository_owner"`
	RepositoryName  string `kallax:"repository_name"`

	LabelList []string `kallax:"labels"`

	UserID                 int64            `kallax:"user_id"`
	UserLogin              string           `kallax:"user_login"`
	MergedByID             int64            `kallax:"merged_by_id"`
	MergedByLogin          string           `kallax:"merged_by_login"`
	AssigneeID             int64            `kallax:"assignee_id"`
	AssigneeLogin          string           `kallax:"assignee_login"`
	AssigneesList          []*UserReference `kallax:"assignees"`
	RequestedReviewersList []*UserReference `kallax:"requested_reviewers"`

	MilestoneID    int64  `kallax:"milestone_id"`
	MilestoneTitle string `kallax:"milestone_title"`

	HeadSHA             string `kallax:"head_sha"`
	HeadRef             string `kallax:"head_ref"`
	HeadLabel           string `kallax:"head_label"`
	HeadUser            string `kallax:"head_user"`
	HeadRepositoryOwner string `kallax:"head_repository_owner"`
	HeadRepositoryName  string `kallax:"head_repository_name"`

	BaseSHA             string `kallax:"base_sha"`
	BaseRef             string `kallax:"base_ref"`
	BaseLabel           string `kallax:"base_label"`
	BaseUser            string `kallax:"base_user"`
	BaseRepositoryOwner string `kallax:"base_repository_owner"`
	BaseRepositoryName  string `kallax:"base_repository_name"`
}

func (i *PullRequest) BeforeSave() error {
	var err error
	i.RepositoryOwner, i.RepositoryName, _, err = utils.ParsePullRequestURL(i.GetURL())
	if err != nil {
		return err
	}

	i.AssigneesList = NewUserRefernceList(i.Assignees)
	i.RequestedReviewersList = NewUserRefernceList(i.RequestedReviewers)

	for _, l := range i.Labels {
		i.LabelList = append(i.LabelList, l.GetName())
	}

	if i.User != nil {
		i.UserID = i.User.GetID()
		i.UserLogin = i.User.GetLogin()
	}

	if i.User != nil {
		i.MergedByID = i.MergedBy.GetID()
		i.MergedByLogin = i.MergedBy.GetLogin()
	}

	if i.Assignee != nil {
		i.AssigneeID = i.Assignee.GetID()
		i.AssigneeLogin = i.Assignee.GetLogin()
	}

	if i.Milestone != nil {
		i.MilestoneID = i.Milestone.GetID()
		i.MilestoneTitle = i.Milestone.GetTitle()
	}

	if i.Head != nil {
		i.HeadSHA = i.Head.GetSHA()
		i.HeadRef = i.Head.GetRef()
		i.HeadLabel = i.Head.GetLabel()

		if i.Head.User != nil {
			i.HeadUser = i.Head.User.GetLogin()
		}

		if i.Head.Repo != nil {
			i.HeadRepositoryOwner = i.Head.Repo.GetOwner().GetLogin()
			i.HeadRepositoryName = i.Head.Repo.GetName()
		}
	}

	if i.Base != nil {
		i.BaseSHA = i.Base.GetSHA()
		i.BaseRef = i.Base.GetRef()
		i.BaseLabel = i.Base.GetLabel()

		if i.Base.User != nil {
			i.BaseUser = i.Base.User.GetLogin()
		}

		if i.Base.Repo != nil {
			i.BaseRepositoryOwner = i.Base.Repo.GetOwner().GetLogin()
			i.BaseRepositoryName = i.Base.Repo.GetName()
		}
	}

	return nil
}

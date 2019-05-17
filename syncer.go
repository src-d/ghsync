package ghsync

import (
	"database/sql"
	"fmt"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-queue.v1"
)

type Syncer struct {
	c *github.Client
	q queue.Queue

	Organization       *OrganizationSyncer
	User               *UserSyncer
	Repository         *RepositorySyncer
	Issues             *IssueSyncer
	IssueComment       *IssueCommentsSyncer
	PullRequest        *PullRequestSyncer
	PullRequestComment *PullRequestCommentSyncer
	PullRequestReview  *PullRequestReviewSyncer
}

func NewSyncer(db *sql.DB, c *github.Client, q queue.Queue) *Syncer {
	return &Syncer{
		c: c,
		q: q,

		Organization:       NewOrganizationSyncer(db, c),
		User:               NewUserSyncer(db, c),
		Repository:         NewRepositorySyncer(db, c),
		Issues:             NewIssueSyncer(db, c),
		IssueComment:       NewIssueCommentsSyncer(db, c),
		PullRequest:        NewPullRequestSyncer(db, c),
		PullRequestComment: NewPullRequestCommentSyncer(db, c),
		PullRequestReview:  NewPullRequestReviewSyncer(db, c),
	}
}

func (s *Syncer) DoOrganization(org string) error {
	if err := s.Organization.Sync(org); err != nil {
		return err
	}

	if err := s.Repository.QueueOrganization(s.q, org); err != nil {
		return err
	}

	if err := s.User.QueueOrganization(s.q, org); err != nil {
		return err
	}

	return nil
}

func (s *Syncer) Wait() error {
	iter, err := s.q.Consume(1)
	if err != nil {
		return err
	}

	for {
		j, err := iter.Next()
		if err != nil {
			return err
		}

		var task *SyncTasks
		if err := j.Decode(&task); err != nil {
			return err
		}

		if err := s.handleSyncTasks(task); err != nil {
			return err
		}

		if err := j.Ack(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Syncer) handleSyncTasks(task *SyncTasks) error {
	fmt.Printf("Handling %q: %s\n", task.Type, task.Payload)

	payload := task.Payload.(map[interface{}]interface{})

	switch task.Type {
	case RepositorySyncTask:
		owner, name := payload["Owner"].(string), payload["Name"].(string)
		if err := s.Issues.QueueRepository(s.q, owner, name); err != nil {
			return err
		}

		if err := s.PullRequest.QueueRepository(s.q, owner, name); err != nil {
			return err
		}

		return s.Repository.Sync(owner, name)
	case UserSyncTask:
		login := payload["Login"].(string)
		return s.User.Sync(login)
	case IssueSyncTask:
		owner, name, number := payload["Owner"].(string), payload["Name"].(string), payload["Number"].(uint64)
		if err := s.IssueComment.QueueIssue(s.q, owner, name, int(number)); err != nil {
			return err
		}

		return s.Issues.Sync(owner, name, int(number))
	case IssueCommentSyncTask:
		owner, name, id := payload["Owner"].(string), payload["Name"].(string), payload["CommentID"].(uint64)
		return s.IssueComment.Sync(owner, name, int64(id))
	case PullRequestSyncTask:
		owner, name, number := payload["Owner"].(string), payload["Name"].(string), payload["Number"].(uint64)
		if err := s.PullRequestComment.QueuePullRequest(s.q, owner, name, int(number)); err != nil {
			return err
		}

		if err := s.PullRequestReview.QueuePullRequest(s.q, owner, name, int(number)); err != nil {
			return err
		}

		return s.PullRequest.Sync(owner, name, int(number))
	case PullRequestCommentSyncTask:
		owner, name, id := payload["Owner"].(string), payload["Name"].(string), payload["CommentID"].(uint64)
		return s.PullRequestComment.Sync(owner, name, int64(id))
	case PullRequestReviewSyncTask:
		owner, name := payload["Owner"].(string), payload["Name"].(string)
		number, id := payload["Number"].(uint), payload["CommentID"].(uint64)

		return s.PullRequestReview.Sync(owner, name, int(number), int64(id))
	}

	return fmt.Errorf("unexpected tasks: %s", task.Type)
}

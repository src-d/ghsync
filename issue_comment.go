package ghsync

import (
	"context"
	"database/sql"

	"github.com/src-d/ghsync/models"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
	log "gopkg.in/src-d/go-log.v1"
	"gopkg.in/src-d/go-queue.v1"
)

type IssueCommentsSyncer struct {
	s *models.IssueCommentStore
	c *github.Client
}

func NewIssueCommentsSyncer(db *sql.DB, c *github.Client) *IssueCommentsSyncer {
	return &IssueCommentsSyncer{
		s: models.NewIssueCommentStore(db),
		c: c,
	}
}

func (s *IssueCommentsSyncer) QueueIssue(q queue.Queue, owner, repo string, number int) error {
	opts := &github.IssueListCommentsOptions{}
	opts.ListOptions.PerPage = 10

	logger := log.New(log.Fields{
		"type":  IssueCommentSyncTask,
		"owner": owner, "repo": repo, "number": number,
	})

	for {
		comments, r, err := s.c.Issues.ListComments(context.TODO(), owner, repo, number, opts)
		if err != nil {
			return err
		}

		for _, c := range comments {
			j, err := NewIssueCommentSyncJob(owner, repo, c.GetID())
			if err != nil {
				return err
			}

			logger.Infof("queue request")
			if err := q.Publish(j); err != nil {
				logger.Errorf(err, "publishing job")
				return nil
			}
		}

		if r.NextPage == 0 {
			break
		}

		opts.Page = r.NextPage
	}

	return nil
}

func (s *IssueCommentsSyncer) Sync(owner string, repo string, commentID int64) error {
	comment, _, err := s.c.Issues.GetComment(context.TODO(), owner, repo, commentID)
	if err != nil {
		return err
	}

	record, err := s.s.FindOne(models.NewIssueCommentQuery().
		Where(kallax.Eq(models.Schema.IssueComment.ID, comment.GetID())),
	)

	if record == nil {
		record = models.NewIssueComment()
		record.IssueComment = *comment

		return s.s.Insert(record)
	}

	record.IssueComment = *comment
	_, err = s.s.Update(record)
	return err

}

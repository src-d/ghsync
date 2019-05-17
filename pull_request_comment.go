package ghsync

import (
	"context"
	"database/sql"

	"github.com/src-d/ghsync/models"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
	"gopkg.in/src-d/go-queue.v1"
)

type PullRequestCommentSyncer struct {
	s *models.PullRequestCommentStore
	c *github.Client
}

func NewPullRequestCommentSyncer(db *sql.DB, c *github.Client) *PullRequestCommentSyncer {
	return &PullRequestCommentSyncer{
		s: models.NewPullRequestCommentStore(db),
		c: c,
	}
}

func (s *PullRequestCommentSyncer) QueuePullRequest(q queue.Queue, owner, repo string, number int) error {
	opts := &github.PullRequestListCommentsOptions{}
	opts.ListOptions.PerPage = 10

	for {
		comments, r, err := s.c.PullRequests.ListComments(context.TODO(), owner, repo, number, opts)
		if err != nil {
			return err
		}

		for _, c := range comments {
			j, err := NewPullRequestCommentSyncJob(owner, repo, c.GetID())
			if err != nil {
				return err
			}

			if err := q.Publish(j); err != nil {
				return err
			}
		}

		if r.NextPage == 0 {
			break
		}

		opts.Page = r.NextPage
	}

	return nil
}

func (s *PullRequestCommentSyncer) Sync(owner string, repo string, commentID int64) error {
	comment, _, err := s.c.PullRequests.GetComment(context.TODO(), owner, repo, commentID)
	if err != nil {
		return err
	}

	record, err := s.s.FindOne(models.NewPullRequestCommentQuery().
		Where(kallax.And(
			kallax.Eq(models.Schema.PullRequestComment.ID, commentID),
		)),
	)
	if record == nil {
		record = models.NewPullRequestComment()
		record.PullRequestComment = *comment

		return s.s.Insert(record)
	}

	record.PullRequestComment = *comment
	_, err = s.s.Update(record)
	return err

}

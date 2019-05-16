package ghsync

import (
	"context"
	"database/sql"

	"github.com/src-d/ghsync/models"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
)

type PullRequestSyncer struct {
	s *models.PullRequestStore
	c *github.Client
}

func NewPullRequestSyncer(db *sql.DB, c *github.Client) *PullRequestSyncer {
	return &PullRequestSyncer{
		s: models.NewPullRequestStore(db),
		c: c,
	}
}

func (s *PullRequestSyncer) Sync(owner string, repo string, number int) error {
	pr, _, err := s.c.PullRequests.Get(context.TODO(), owner, repo, number)
	if err != nil {
		return err
	}

	record, err := s.s.FindOne(models.NewPullRequestQuery().
		Where(kallax.And(
			kallax.Eq(models.Schema.PullRequest.RepositoryOwner, owner),
			kallax.Eq(models.Schema.PullRequest.RepositoryName, repo),
			kallax.Eq(models.Schema.PullRequest.Number, number),
		)),
	)
	if record == nil {
		record = models.NewPullRequest()
		record.PullRequest = *pr

		return s.s.Insert(record)
	}

	record.PullRequest = *pr
	_, err = s.s.Update(record)
	return err

}

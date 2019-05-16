package ghsync

import (
	"context"
	"database/sql"

	"github.com/src-d/ghsync/models"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
)

type IssueSyncer struct {
	s *models.IssueStore
	c *github.Client
}

func NewIssueSyncer(db *sql.DB, c *github.Client) *IssueSyncer {
	return &IssueSyncer{
		s: models.NewIssueStore(db),
		c: c,
	}
}

func (s *IssueSyncer) Sync(owner string, repo string, number int) error {
	issue, _, err := s.c.Issues.Get(context.TODO(), owner, repo, number)
	if err != nil {
		return err
	}

	record, err := s.s.FindOne(models.NewIssueQuery().
		Where(kallax.And(
			kallax.Eq(models.Schema.Issue.RepositoryOwner, owner),
			kallax.Eq(models.Schema.Issue.RepositoryName, repo),
			kallax.Eq(models.Schema.Issue.Number, number),
		)),
	)

	if record == nil {
		record = models.NewIssue()
		record.Issue = *issue

		return s.s.Insert(record)
	}

	record.Issue = *issue
	_, err = s.s.Update(record)
	return err

}

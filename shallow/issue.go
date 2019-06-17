package shallow

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/src-d/ghsync/models"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
	"gopkg.in/src-d/go-log.v1"
)

type IssueSyncer struct {
	db     *sql.DB
	store  *models.IssueStore
	client *github.Client
}

func NewIssueSyncer(db *sql.DB, c *github.Client) *IssueSyncer {
	return &IssueSyncer{
		db:     db,
		store:  models.NewIssueStore(db),
		client: c,
	}
}

func (s *IssueSyncer) Sync(owner, repo string, logger log.Logger) error {
	opts := &github.IssueListByRepoOptions{}
	opts.ListOptions.PerPage = listOptionsPerPage
	opts.State = "all"

	logger.Infof("starting to retrieve issues")

	// TODO transaction for faster times

	// Get the list of all issues
	for {
		issues, r, err := s.client.Issues.ListByRepo(context.TODO(), owner, repo, opts)
		if err != nil {
			return err
		}

		for _, i := range issues {
			if i.IsPullRequest() {
				continue
			}

			logger := logger.With(log.Fields{"issue": i.GetNumber()})

			_, err := s.store.FindOne(models.NewIssueQuery().
				Where(kallax.And(
					kallax.Eq(models.Schema.Issue.RepositoryOwner, owner),
					kallax.Eq(models.Schema.Issue.RepositoryName, repo),
					kallax.Eq(models.Schema.Issue.Number, i.GetNumber()),
				)),
			)

			if err != nil && err != kallax.ErrNotFound {
				logger.Errorf(err, "failed to read the resource from the DB")
				return fmt.Errorf("failed to read the resource from the DB: %v", err)
			}

			if err == nil {
				logger.Infof("resource already exists, skipping")
				continue
			}

			record := models.NewIssue()
			record.Issue = *i

			err = s.store.Insert(record)
			if err != nil {
				logger.Errorf(err, "failed to write the resource into the DB")
				return fmt.Errorf("failed to write the resource into the DB: %v", err)
			}

			logger.Debugf("resource written in the DB")
		}

		if r.NextPage == 0 {
			break
		}

		opts.Page = r.NextPage
	}

	logger.Infof("finished to retrieve issues")

	return nil
}

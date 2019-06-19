package shallow

import (
	"context"
	"fmt"

	"github.com/src-d/ghsync/models"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-log.v1"
)

type IssueSyncer struct {
	store  *models.IssueStore
	client *github.Client
}

func NewIssueSyncer(s *models.IssueStore, c *github.Client) *IssueSyncer {
	return &IssueSyncer{
		store:  s,
		client: c,
	}
}

func (s *IssueSyncer) Sync(owner, repo string, logger log.Logger) error {
	return s.store.Transaction(func(store *models.IssueStore) error {
		return s.doIssues(store, owner, repo, logger)
	})
}

func (s *IssueSyncer) doIssues(store *models.IssueStore, owner, repo string, logger log.Logger) error {
	opts := &github.IssueListByRepoOptions{}
	opts.ListOptions.PerPage = listOptionsPerPage
	opts.State = "all"

	logger.Infof("starting to retrieve issues")

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

			record := models.NewIssue()
			record.Issue = *i

			err = store.Insert(record)
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

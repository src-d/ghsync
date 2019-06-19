package shallow

import (
	"context"
	"fmt"

	"github.com/src-d/ghsync/models"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-log.v1"
)

type PullRequestSyncer struct {
	store  *models.PullRequestStore
	client *github.Client
}

func NewPullRequestSyncer(s *models.PullRequestStore, c *github.Client) *PullRequestSyncer {
	return &PullRequestSyncer{
		store:  s,
		client: c,
	}
}

func (s *PullRequestSyncer) Sync(owner, repo string, logger log.Logger) error {
	return s.store.Transaction(func(store *models.PullRequestStore) error {
		return s.doPRs(store, owner, repo, logger)
	})
}

func (s *PullRequestSyncer) doPRs(store *models.PullRequestStore, owner, repo string, logger log.Logger) error {
	opts := &github.PullRequestListOptions{}
	opts.ListOptions.PerPage = listOptionsPerPage
	opts.State = "all"

	logger.Infof("starting to retrieve PRs")

	// Get the list of all PRs
	for {
		prs, r, err := s.client.PullRequests.List(context.TODO(), owner, repo, opts)
		if err != nil {
			return err
		}

		for _, pr := range prs {
			logger := logger.With(log.Fields{"pr": pr.GetNumber()})

			record := models.NewPullRequest()
			record.PullRequest = *pr

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

	logger.Infof("finished to retrieve PRs")

	return nil
}

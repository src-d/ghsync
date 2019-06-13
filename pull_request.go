package ghsync

import (
	"context"
	"database/sql"

	"github.com/src-d/ghsync/models"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
	"gopkg.in/src-d/go-log.v1"
	"gopkg.in/src-d/go-queue.v1"
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

func (s *PullRequestSyncer) QueueRepository(q queue.Queue, owner, repo string) error {
	opts := &github.PullRequestListOptions{}
	opts.ListOptions.PerPage = listOptionsPerPage
	opts.State = "all"

	logger := log.New(log.Fields{"type": PullRequestSyncTask, "owner": owner, "repo": repo})
	logger.Infof("starting to publish queue jobs")

	for {
		requests, r, err := s.c.PullRequests.List(context.TODO(), owner, repo, opts)
		if err != nil {
			return err
		}

		for _, r := range requests {
			j, err := NewPullRequestSyncJob(owner, repo, r.GetNumber())
			if err != nil {
				return err
			}

			l := logger.With(log.Fields{"pull-request": r.GetNumber()})
			l.Debugf("queue request")
			if err := q.Publish(j); err != nil {
				l.Errorf(err, "publishing job")
				return nil
			}
		}

		if r.NextPage == 0 {
			break
		}

		opts.Page = r.NextPage
	}

	logger.Infof("finished to publish queue jobs")

	return nil
}

func (s *PullRequestSyncer) Sync(owner string, repo string, number int) error {
	pr, _, err := s.c.PullRequests.Get(context.TODO(), owner, repo, number)
	if err != nil {
		return err
	}

	record, err := s.s.FindOne(models.NewPullRequestQuery().
		Where(kallax.And(
			kallax.Eq(models.Schema.PullRequest.ID, pr.GetID()),
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

package ghsync

import (
	"context"
	"database/sql"

	"github.com/src-d/ghsync/models"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
)

type PullRequestReviewSyncer struct {
	s *models.PullRequestReviewStore
	c *github.Client
}

func NewPullRequestReviewSyncer(db *sql.DB, c *github.Client) *PullRequestReviewSyncer {
	return &PullRequestReviewSyncer{
		s: models.NewPullRequestReviewStore(db),
		c: c,
	}
}
func (s *PullRequestReviewSyncer) SyncRepository(owner, repo string) error {
	return s.SyncPullRequest(owner, repo, 0)
}

func (s *PullRequestReviewSyncer) SyncPullRequest(owner, repo string, number int) error {
	opts := &github.ListOptions{}
	opts.PerPage = 10

	for {
		reviews, r, err := s.c.PullRequests.ListReviews(context.TODO(), owner, repo, number, opts)
		if err != nil {
			return err
		}

		for _, r := range reviews {
			if err := s.doSync(r); err != nil {
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

func (s *PullRequestReviewSyncer) Sync(owner string, repo string, number int, reviewID int64) error {
	review, _, err := s.c.PullRequests.GetReview(context.TODO(), owner, repo, number, reviewID)
	if err != nil {
		return err
	}

	return s.doSync(review)
}

func (s *PullRequestReviewSyncer) doSync(review *github.PullRequestReview) error {
	record, err := s.s.FindOne(models.NewPullRequestReviewQuery().
		Where(kallax.And(
			kallax.Eq(models.Schema.PullRequestReview.ID, review.GetID()),
		)),
	)
	if record == nil {
		record = models.NewPullRequestReview()
		record.PullRequestReview = *review

		return s.s.Insert(record)
	}

	record.PullRequestReview = *review
	_, err = s.s.Update(record)
	return err

}

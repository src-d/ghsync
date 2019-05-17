package ghsync

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/src-d/ghsync/models"
	"github.com/src-d/ghsync/utils"

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

func (s *PullRequestReviewSyncer) QueuePullRequest(owner, repo string, number int) error {
	opts := &github.ListOptions{}
	opts.PerPage = 10

	for {
		reviews, r, err := s.c.PullRequests.ListReviews(context.TODO(), owner, repo, number, opts)
		if err != nil {
			return err
		}

		for _, r := range reviews {
			_, _, number, _ := utils.ParsePullRequestURL(r.GetPullRequestURL())
			fmt.Println(s.Sync(owner, repo, number, r.GetID()))
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

	record, err := s.s.FindOne(models.NewPullRequestReviewQuery().
		Where(kallax.And(
			kallax.Eq(models.Schema.PullRequestReview.ID, reviewID),
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

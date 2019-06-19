package deep

import (
	"context"
	"database/sql"

	"github.com/src-d/ghsync/models"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
	log "gopkg.in/src-d/go-log.v1"
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

func (s *IssueCommentsSyncer) SyncRepository(owner, repo string) error {
	return s.SyncIssue(owner, repo, 0)
}

func (s *IssueCommentsSyncer) SyncIssue(owner, repo string, number int) error {
	opts := &github.IssueListCommentsOptions{}
	opts.ListOptions.PerPage = listOptionsPerPage

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
			if err := s.doSync(c); err != nil {
				logger.Errorf(err, "issue sync error")
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

	return s.doSync(comment)
}

func (s *IssueCommentsSyncer) doSync(comment *github.IssueComment) error {
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

package ghsync

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/src-d/ghsync/models"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
)

type RepositorySyncer struct {
	s *models.RepositoryStore
	c *github.Client
}

func NewRepositorySyncer(db *sql.DB, c *github.Client) *RepositorySyncer {
	return &RepositorySyncer{
		s: models.NewRepositoryStore(db),
		c: c,
	}
}

func (s *RepositorySyncer) QueueOrganization(owner string) error {
	opts := &github.RepositoryListOptions{}
	opts.ListOptions.PerPage = 10

	for {
		repositories, r, err := s.c.Repositories.List(context.TODO(), owner, opts)
		if err != nil {
			return err
		}

		for _, r := range repositories {
			fmt.Println(s.Sync(owner, r.GetName()))
		}

		if r.NextPage == 0 {
			break
		}

		opts.Page = r.NextPage
	}

	return nil
}

func (s *RepositorySyncer) Sync(owner, name string) error {
	repository, _, err := s.c.Repositories.Get(context.TODO(), owner, name)
	if err != nil {
		return err
	}

	record, err := s.s.FindOne(models.NewRepositoryQuery().
		Where(kallax.Eq(models.Schema.Repository.ID, repository.GetID())),
	)

	if record == nil {
		record = models.NewRepository()
		record.Repository = *repository

		return s.s.Insert(record)
	}

	record.Repository = *repository
	_, err = s.s.Update(record)
	return err

}

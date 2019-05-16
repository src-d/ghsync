package ghsync

import (
	"context"
	"database/sql"

	"github.com/src-d/ghsync/models"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
)

type OrganizationSyncer struct {
	s *models.OrganizationStore
	c *github.Client
}

func NewOrganizationSyncer(db *sql.DB, c *github.Client) *OrganizationSyncer {
	return &OrganizationSyncer{
		s: models.NewOrganizationStore(db),
		c: c,
	}
}

func (s *OrganizationSyncer) Sync(login string) error {
	org, _, err := s.c.Organizations.Get(context.TODO(), login)
	if err != nil {
		return err
	}

	record, err := s.s.FindOne(models.NewOrganizationQuery().
		Where(kallax.Eq(models.Schema.Organization.Login, login)),
	)

	if record == nil {
		record = models.NewOrganization()
		record.Organization = *org

		return s.s.Insert(record)
	}

	record.Organization = *org
	_, err = s.s.Update(record)
	return err

}

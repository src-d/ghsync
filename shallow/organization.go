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

type OrganizationSyncer struct {
	db              *sql.DB
	store           *models.OrganizationStore
	client          *github.Client
	statusTableName string
}

func NewOrganizationSyncer(db *sql.DB, c *github.Client, statusTableName string) *OrganizationSyncer {
	return &OrganizationSyncer{
		db:              db,
		store:           models.NewOrganizationStore(db),
		client:          c,
		statusTableName: statusTableName,
	}
}

func (s *OrganizationSyncer) Sync(login string) error {
	logger := log.With(log.Fields{"organization": login})

	_, err := s.store.FindOne(models.NewOrganizationQuery().
		Where(kallax.Eq(models.Schema.Organization.Login, login)),
	)

	if err != nil && err != kallax.ErrNotFound {
		logger.Errorf(err, "failed to read the resource from the DB")
		return fmt.Errorf("failed to read the resource from the DB: %v", err)
	}

	if err == nil {
		logger.Infof("resource already exists, skipping")
		stm := fmt.Sprintf("UPDATE %s SET total=0 WHERE org='%s'", s.statusTableName, login)
		_, err = s.db.Exec(stm)
		if err != nil {
			return fmt.Errorf("unable to update status for org %s: %v", login, err)
		}

		return nil
	}

	org, _, err := s.client.Organizations.Get(context.TODO(), login)
	if err != nil {
		return err
	}

	repoSyncer := NewRepositorySyncer(s.db, s.client, s.statusTableName)
	err = repoSyncer.Sync(login, logger)
	if err != nil {
		return err
	}

	userSyncer := NewUserSyncer(s.db, s.client, s.statusTableName)
	err = userSyncer.Sync(login, logger)
	if err != nil {
		return err
	}

	record := models.NewOrganization()
	record.Organization = *org

	logger.Debugf("inserting resource")

	err = s.store.Insert(record)
	if err != nil {
		logger.Errorf(err, "failed to write the resource into the DB")
		return fmt.Errorf("failed to write the resource into the DB: %v", err)
	}

	logger.Debugf("resource written in the DB")

	return nil
}

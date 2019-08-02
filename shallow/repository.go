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

type RepositorySyncer struct {
	db              *sql.DB
	store           *models.RepositoryStore
	client          *github.Client
	statusTableName string
	skipForks       bool
}

func NewRepositorySyncer(db *sql.DB, c *github.Client, statusTableName string, skipForks bool) *RepositorySyncer {
	return &RepositorySyncer{
		db:              db,
		store:           models.NewRepositoryStore(db),
		client:          c,
		statusTableName: statusTableName,
		skipForks:       skipForks,
	}
}

func (s *RepositorySyncer) Sync(owner string, logger log.Logger) error {
	opts := &github.RepositoryListByOrgOptions{}
	opts.ListOptions.PerPage = listOptionsPerPage

	logger.Infof("starting to retrieve repositories")

	repos := make([]*github.Repository, 0)

	// Get the list of all repositories
	for {
		repositories, r, err := s.client.Repositories.ListByOrg(context.TODO(), owner, opts)
		if err != nil {
			return err
		}

		for _, r := range repositories {
			if s.skipForks && r.GetFork() {
				continue
			}
			repos = append(repos, r)
		}

		if r.NextPage == 0 {
			break
		}

		opts.Page = r.NextPage
	}

	stm := fmt.Sprintf("UPDATE %s SET total=%d WHERE org='%s' AND entity='repository'",
		s.statusTableName, len(repos), owner)
	log.Debugf("running statement: %s", stm)
	if _, err := s.db.Exec(stm); err != nil {
		return fmt.Errorf("an error occured while updating %s table: %v",
			s.statusTableName, err)
	}

	// Process each one of them
	for _, repository := range repos {
		err := s.doRepo(repository, logger)
		if err != nil {
			stm := fmt.Sprintf("UPDATE %s SET failed=failed + 1 WHERE org='%s' AND entity='repository'",
				s.statusTableName, owner)
			if err = s.updateStatus(stm); err != nil {
				return err
			}

			return err
		}

		stm := fmt.Sprintf("UPDATE %s SET done=done + 1 WHERE org='%s' AND entity='repository'",
			s.statusTableName, owner)
		if err = s.updateStatus(stm); err != nil {
			return err
		}
	}

	logger.Infof("finished to retrieve repositories")

	return nil
}

func (s *RepositorySyncer) updateStatus(stm string) error {
	log.Debugf("running statement: %s", stm)
	if _, err := s.db.Exec(stm); err != nil {
		return fmt.Errorf("an error occured while updating %s table: %v",
			s.statusTableName, err)
	}

	return nil
}

func (s *RepositorySyncer) doRepo(repository *github.Repository, parentLogger log.Logger) error {
	logger := parentLogger.With(log.Fields{"repository": repository.GetName()})

	_, err := s.store.FindOne(models.NewRepositoryQuery().
		Where(kallax.Eq(models.Schema.Repository.ID, repository.GetID())),
	)

	if err != nil && err != kallax.ErrNotFound {
		logger.Errorf(err, "failed to read the resource from the DB")
		return fmt.Errorf("failed to read the resource from the DB: %v", err)
	}

	if err == nil {
		logger.Infof("resource already exists, skipping")
		return nil
	}

	prSyncer := NewPullRequestSyncer(s.db, s.client)
	err = prSyncer.Sync(repository.GetOwner().GetLogin(), repository.GetName(), logger)
	if err != nil {
		return err
	}

	issueSyncer := NewIssueSyncer(s.db, s.client)
	err = issueSyncer.Sync(repository.GetOwner().GetLogin(), repository.GetName(), logger)
	if err != nil {
		return err
	}

	record := models.NewRepository()
	record.Repository = *repository

	err = s.store.Insert(record)
	if err != nil {
		logger.Errorf(err, "failed to write the resource into the DB")
		return fmt.Errorf("failed to write the resource into the DB: %v", err)
	}

	logger.Debugf("resource written in the DB")

	return nil
}

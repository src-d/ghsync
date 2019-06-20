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

type UserSyncer struct {
	db              *sql.DB
	client          *github.Client
	statusTableName string
}

func NewUserSyncer(db *sql.DB, c *github.Client, statusTableName string) *UserSyncer {
	return &UserSyncer{
		db:              db,
		client:          c,
		statusTableName: statusTableName,
	}
}

func (s *UserSyncer) Sync(org string, logger log.Logger) error {
	store := models.NewUserStore(s.db)
	return store.Transaction(func(store *models.UserStore) error {
		return s.doUsers(store, org, logger)
	})
}

func (s *UserSyncer) doUsers(store *models.UserStore, org string, logger log.Logger) error {
	opts := &github.ListMembersOptions{}
	opts.ListOptions.PerPage = listOptionsPerPage

	logger.Infof("starting to retrieve users")

	allUsers := make([]*github.User, 0)

	// Get the list of all users
	for {
		users, r, err := s.client.Organizations.ListMembers(context.TODO(), org, opts)
		if err != nil {
			return err
		}

		for _, u := range users {
			allUsers = append(allUsers, u)
		}

		if r.NextPage == 0 {
			break
		}

		opts.Page = r.NextPage
	}

	stm := fmt.Sprintf("UPDATE %s SET total=%d WHERE org='%s' AND part='user'",
		s.statusTableName, len(allUsers), org)
	log.Debugf("running statement: %s", stm)
	if _, err := s.db.Exec(stm); err != nil {
		return fmt.Errorf("an error occured while updating %s table: %v",
			s.statusTableName, err)
	}

	for _, user := range allUsers {
		logger := logger.With(log.Fields{"user": user.GetLogin()})

		_, err := store.FindOne(models.NewUserQuery().
			Where(kallax.And(
				kallax.Eq(models.Schema.User.ID, user.GetID()),
			)),
		)
		if err != nil && err != kallax.ErrNotFound {
			logger.With(log.Fields{"user": user.GetLogin()}).Errorf(err, "failed to read the resource from the DB")
			return fmt.Errorf("failed to read the resource from the DB: %v", err)
		}

		if err == nil {
			logger.With(log.Fields{"user": user.GetLogin()}).Infof("resource already exists, skipping")
			continue
		}

		record := models.NewUser()
		record.User = *user

		err = store.Insert(record)
		if err != nil {
			logger.Errorf(err, "failed to write the resource into the DB")
			return fmt.Errorf("failed to write the resource into the DB: %v", err)
		}

		logger.Debugf("resource written in the DB")

		stm := fmt.Sprintf("UPDATE %s SET done=done + 1 WHERE org='%s' AND part='user'",
			s.statusTableName, org)
		log.Debugf("running statement: %s", stm)
		if _, err := s.db.Exec(stm); err != nil {
			return fmt.Errorf("an error occured while updating %s table: %v",
				s.statusTableName, err)
		}
	}

	logger.Infof("finished to retrieve users")

	return nil
}

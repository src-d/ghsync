package ghsync

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/src-d/ghsync/models"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
)

type UserSyncer struct {
	s *models.UserStore
	c *github.Client
}

func NewUserSyncer(db *sql.DB, c *github.Client) *UserSyncer {
	return &UserSyncer{
		s: models.NewUserStore(db),
		c: c,
	}
}

func (s *UserSyncer) QueueOrganization(org string) error {
	opts := &github.ListMembersOptions{}
	opts.ListOptions.PerPage = 10

	for {
		users, r, err := s.c.Organizations.ListMembers(context.TODO(), org, opts)
		if err != nil {
			return err
		}

		for _, u := range users {
			fmt.Println(s.Sync(u.GetLogin()))
		}

		if r.NextPage == 0 {
			break
		}

		opts.Page = r.NextPage
	}

	return nil
}

func (s *UserSyncer) Sync(login string) error {
	user, _, err := s.c.Users.Get(context.TODO(), login)
	if err != nil {
		return err
	}

	record, err := s.s.FindOne(models.NewUserQuery().
		Where(kallax.And(
			kallax.Eq(models.Schema.User.ID, user.GetID()),
		)),
	)

	if record == nil {
		record = models.NewUser()
		record.User = *user

		return s.s.Insert(record)
	}

	record.User = *user
	_, err = s.s.Update(record)
	return err

}

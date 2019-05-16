package models

import "github.com/google/go-github/github"

type UserReference struct {
	ID    int64
	Login string
}

func NewUserRefernceList(users []*github.User) []*UserReference {
	list := make([]*UserReference, len(users))
	for i, u := range users {
		list[i] = &UserReference{
			ID:    u.GetID(),
			Login: u.GetLogin(),
		}
	}

	return list
}

package models

import "github.com/google/go-github/github"

type UserReference struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
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

type RepositoryReference struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Fork       bool   `json:"fork"`
	Size       int    `json:"size"`
	OwnerLogin string `json:"owner_login"`
	OwnerType  string `json:"owner_type"`
	OwnerID    int64  `json:"owner_id"`
}

func NewRepositoryReference(r *github.Repository) *RepositoryReference {
	if r == nil {
		return nil
	}

	rr := &RepositoryReference{}
	rr.ID = r.GetID()
	rr.Name = r.GetName()
	rr.Fork = r.GetFork()
	rr.Size = r.GetSize()

	if r.GetOwner() != nil {
		rr.OwnerLogin = r.Owner.GetLogin()
		rr.OwnerType = r.Owner.GetType()
		rr.OwnerID = r.Owner.GetID()
	}
	return rr
}

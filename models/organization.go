package models

import (
	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-kallax.v1"
)

type Organization struct {
	kallax.Model `table:"organizations" pk:"id"`
	github.Organization
}

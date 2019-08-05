package circleci

import (
	"github.com/jszwedko/go-circleci"
	"gopkg.in/src-d/go-kallax.v1"
)

type Step struct {
	kallax.Model `table:"steps" pk:"id,autoincr" ignored:""`
	circleci.Action

	ID       int64
	BuildNum int
	Username string
	Reponame string
}

func (b *Step) BeforeSave() error {
	if len(b.Messages) == 0 {
		b.Messages = []string{}
	}

	return nil
}

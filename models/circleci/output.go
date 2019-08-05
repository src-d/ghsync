package circleci

import (
	"github.com/jszwedko/go-circleci"
	"gopkg.in/src-d/go-kallax.v1"
)

type Output struct {
	kallax.Model `table:"outputs" pk:"id,autoincr" ignored:""`
	circleci.Output

	ID       int64
	URL      string
	BuildNum int
	Username string
	Reponame string
}

func (o *Output) BeforeSave() error {
	return nil
}

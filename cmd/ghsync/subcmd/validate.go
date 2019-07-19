package subcmd

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/src-d/go-cli.v0"
)

type ValidateCommand struct {
	cli.Command `name:"validate" short-description:"Validate token and list of organizations" long-description:"Validate token and list of organizations.\nReturns 0 exit code in case of success and 1 exit code on error."`

	GithubOptions
}

func (c *ValidateCommand) Execute(args []string) error {
	client, err := newClient(c.Token)
	if err != nil {
		return err
	}

	_, resp, err := client.Users.Get(context.TODO(), "")
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("github token is not valid")
	}

	for _, org := range strings.Split(c.Orgs, ",") {
		_, resp, err := client.Organizations.Get(context.TODO(), org)
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("organization '%s' is not found", org)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

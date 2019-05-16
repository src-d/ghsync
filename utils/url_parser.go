package utils

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func ParsePullRequestURL(rawurl string) (owner, repo string, number int, err error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return
	}

	parts := strings.Split(u.RequestURI(), "/")

	switch u.Host {
	case "api.github.com":
		// https://api.github.com/repos/octocat/Hello-World/pulls/1347
		owner = parts[2]
		repo = parts[3]
		number, _ = strconv.Atoi(parts[5])

	case "github.com":
		// https://github.com/src-d/go-kallax/pull/309
		owner = parts[1]
		repo = parts[2]
		number, _ = strconv.Atoi(parts[4])

	default:
		err = fmt.Errorf("unsupported url: %s", rawurl)
	}

	return
}

func ParseIssueURL(rawurl string) (owner, repo string, number int, err error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return
	}

	parts := strings.Split(u.RequestURI(), "/")

	switch u.Host {
	case "api.github.com":
		// "https://api.github.com/repos/octocat/Hello-World/issues/1347",
		owner = parts[2]
		repo = parts[3]
		number, _ = strconv.Atoi(parts[5])

	case "github.com":
		// https://github.com/cncf/devstats-example/issues/2
		owner = parts[1]
		repo = parts[2]
		number, _ = strconv.Atoi(parts[4])

	default:
		err = fmt.Errorf("unsupported url: %s", rawurl)
	}

	return
}

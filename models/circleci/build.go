package circleci

import (
	"github.com/jszwedko/go-circleci"
	"gopkg.in/src-d/go-kallax.v1"
)

// context_ids
// all_Commit_details_Truncated
// owner
// has_Artifacts
// vcs_type
// no_depedency_cache
// ssh_disabled
// canceler

type Build struct {
	kallax.Model `table:"builds" pk:"id,autoincr" ignored:"User,Steps,FeaturedFlags,SSHEnabled,Workflows,AllCommitDetails,SSHUsers,Picard,Messages,Node,Previous,PreviousSuccessfulBuild,PullRequests"`
	circleci.Build

	ID int64

	PreviousBuildNum           int
	PreviousSuccessfulBuildNum int
	PullRequestURLs            []string
}

func (b *Build) BeforeSave() error {
	if b.Previous != nil {
		b.PreviousBuildNum = b.Previous.BuildNum
	}
	if b.PreviousSuccessfulBuild != nil {
		b.PreviousSuccessfulBuildNum = b.PreviousSuccessfulBuild.BuildNum
	}

	for _, pr := range b.PullRequests {
		b.PullRequestURLs = append(b.PullRequestURLs, pr.URL)
	}

	if b.Retries == nil {
		b.Retries = []int{}
	}

	if b.PullRequestURLs == nil {
		b.PullRequestURLs = []string{}
	}
	return nil
}

package ghsync

import (
	"strings"

	"gopkg.in/src-d/go-log.v1"
	"gopkg.in/src-d/go-queue.v1"
)

type SyncTaskType string

const (
	RepositorySyncTask         SyncTaskType = "repository"
	UserSyncTask               SyncTaskType = "user"
	IssueSyncTask              SyncTaskType = "issue"
	IssueCommentSyncTask       SyncTaskType = "issue-comment"
	PullRequestSyncTask        SyncTaskType = "pull-request"
	PullRequestCommentSyncTask SyncTaskType = "pull-request-comment"
	PullRequestReviewSyncTask  SyncTaskType = "pull-request-review"
)

type SyncTasks struct {
	Type    SyncTaskType
	Payload interface{}
}

func newSyncTasks(t SyncTaskType, payload interface{}) (*queue.Job, error) {
	j, err := queue.NewJob()
	if err != nil {
		return nil, err
	}

	err = j.Encode(&SyncTasks{
		Type:    t,
		Payload: payload,
	})

	if err != nil {
		return nil, err
	}

	return j, nil
}

type RepositorySyncPayload struct {
	Owner string
	Name  string
}

func NewRepositorySyncJob(owner, name string) (*queue.Job, error) {
	return newSyncTasks(RepositorySyncTask, RepositorySyncPayload{owner, name})
}

type UserSyncPayload struct {
	Login string
}

func NewUserSyncJob(login string) (*queue.Job, error) {
	return newSyncTasks(UserSyncTask, UserSyncPayload{login})
}

type IssueSyncPayload struct {
	Owner  string
	Name   string
	Number uint64
}

func NewIssueSyncJob(owner, name string, number int) (*queue.Job, error) {
	return newSyncTasks(IssueSyncTask, IssueSyncPayload{owner, name, uint64(number)})
}

func NewPullRequestSyncJob(owner, name string, number int) (*queue.Job, error) {
	return newSyncTasks(PullRequestSyncTask, IssueSyncPayload{owner, name, uint64(number)})
}

type IssueCommentSyncPayload struct {
	Owner     string
	Name      string
	CommentID uint64
}

func NewIssueCommentSyncJob(owner, name string, id int64) (*queue.Job, error) {
	return newSyncTasks(IssueCommentSyncTask, IssueCommentSyncPayload{owner, name, uint64(id)})
}

func NewPullRequestCommentSyncJob(owner, name string, id int64) (*queue.Job, error) {
	return newSyncTasks(PullRequestCommentSyncTask, IssueCommentSyncPayload{owner, name, uint64(id)})
}

type PullRequestReviewSyncPayload struct {
	Owner    string
	Name     string
	Number   uint64
	ReviewID uint64
}

func NewPullRequestReviewSyncJob(owner, name string, number int, id int64) (*queue.Job, error) {
	return newSyncTasks(PullRequestReviewSyncTask,
		PullRequestReviewSyncPayload{owner, name, uint64(number), uint64(id)})
}

func logFieldsFromPayload(payload map[interface{}]interface{}) log.Fields {
	fields := make(log.Fields, len(payload))
	for k, v := range payload {
		fields[strings.ToLower(k.(string))] = v
	}

	return fields
}

package circleci

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/jszwedko/go-circleci"
	"gopkg.in/src-d/go-kallax.v1"

	models "github.com/src-d/ghsync/models/circleci"
)

// CREATE INDEX build_search ON builds (build_num, reponame, username);
// CREATE INDEX step_search ON steps (build_num, reponame, username, name, index);
// CREATE INDEX output_search ON outputs (build_num, reponame, username, url);

type BuildSyncer struct {
	builds  *models.BuildStore
	steps   *models.StepStore
	outputs *models.OutputStore
	c       *circleci.Client
}

func NewBuildSyncer(db *sql.DB, c *circleci.Client) *BuildSyncer {
	return &BuildSyncer{
		builds:  models.NewBuildStore(db),
		steps:   models.NewStepStore(db),
		outputs: models.NewOutputStore(db),
		c:       c,
	}
}

func (s *BuildSyncer) Sync(account, repository string) error {
	builds, err := s.c.ListRecentBuildsForProject(account, repository, "", "", 100, 0)
	if err != nil {
		return err
	}

	fmt.Println(account, repository)
	for _, build := range builds {
		if err := s.SyncOne(account, repository, build.BuildNum); err != nil {
			return err
		}

		continue

		if err := s.doSyncBuild(build); err != nil {
			return err
		}
	}

	return err
}

func (s *BuildSyncer) SyncOne(account, repository string, buildNum int) error {
	build, err := s.c.GetBuild(account, repository, buildNum)
	if err != nil {
		return err
	}

	fmt.Println(build.Steps)
	return s.doSyncBuild(build)
}

func (s *BuildSyncer) doSyncBuild(build *circleci.Build) error {
	if err := s.doSyncSteps(build); err != nil {
		return err
	}

	record, err := s.builds.FindOne(models.NewBuildQuery().
		Where(kallax.And(
			kallax.Eq(models.Schema.Build.BuildNum, build.BuildNum),
			kallax.Eq(models.Schema.Build.Reponame, build.Reponame),
			kallax.Eq(models.Schema.Build.Username, build.Username),
		)),
	)

	if record == nil {
		record = models.NewBuild()
		record.Build = *build

		return s.builds.Insert(record)
	}

	record.Build = *build
	_, err = s.builds.Update(record)
	return err
}

func (s *BuildSyncer) doSyncSteps(build *circleci.Build) error {
	for _, step := range build.Steps {
		for _, action := range step.Actions {
			if err := s.doSyncAction(build, action); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *BuildSyncer) doSyncAction(build *circleci.Build, action *circleci.Action) error {
	if err := s.doSyncActionOutput(build, action); err != nil {
		return err
	}

	record, err := s.steps.FindOne(models.NewStepQuery().
		Where(kallax.And(
			kallax.Eq(models.Schema.Step.BuildNum, build.BuildNum),
			kallax.Eq(models.Schema.Step.Reponame, build.Reponame),
			kallax.Eq(models.Schema.Step.Username, build.Username),
			kallax.Eq(models.Schema.Step.Name, action.Name),
			kallax.Eq(models.Schema.Step.Index, action.Index),
		)),
	)

	if record == nil {
		record = models.NewStep()
		record.BuildNum = build.BuildNum
		record.Username = build.Username
		record.Reponame = build.Reponame
		record.Action = *action

		return s.steps.Insert(record)
	}

	record.Action = *action
	_, err = s.steps.Update(record)
	return err

}

func (s *BuildSyncer) doSyncActionOutput(build *circleci.Build, action *circleci.Action) error {
	url, _ := url.Parse(action.OutputURL)
	url.RawQuery = ""

	records, err := s.outputs.Count(models.NewOutputQuery().
		Where(kallax.And(
			kallax.Eq(models.Schema.Output.BuildNum, build.BuildNum),
			kallax.Eq(models.Schema.Output.Reponame, build.Reponame),
			kallax.Eq(models.Schema.Output.Username, build.Username),
			kallax.Eq(models.Schema.Output.URL, url.String()),
		)),
	)

	if err != nil || records != 0 {
		return err
	}

	outputs, err := s.c.GetActionOutputs(action)
	if err != nil {
		return err
	}

	for _, output := range outputs {
		if err := s.doInsertOutput(url.String(), build, output); err != nil {
			return err
		}
	}

	return err
}

func (s *BuildSyncer) doInsertOutput(url string, build *circleci.Build, output *circleci.Output) error {
	record := models.NewOutput()
	record.BuildNum = build.BuildNum
	record.Username = build.Username
	record.Reponame = build.Reponame
	record.URL = url
	record.Output = *output

	return s.outputs.Insert(record)
}

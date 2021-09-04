package github_test

import (
	"github.com/dev-this/terraform-gha-controller/internal/github"
	"github.com/dev-this/terraform-gha-controller/internal/runner"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckRunSession_SetStatus(t *testing.T) {
	session := github.CheckRunSession{}
	want := github.CheckRunSession{Status: "in_progress"}

	session.SetStatus("in_progress")

	assert.EqualValues(t, want, session)
}

func TestCheckRunSession_SetCompletion(t *testing.T) {
	session := github.CheckRunSession{}
	want := github.CheckRunSession{Status: "completed", Conclusion: "success"}

	session.SetCompletion("completed", "success")

	assert.EqualValues(t, want, session)
}

func TestCheckRunSession_SetTempRepoDir(t *testing.T) {
	session := github.CheckRunSession{}
	want := github.CheckRunSession{TempRepoDir: "/tmp/my-cool-repo"}

	session.SetTempRepoDir("/tmp/my-cool-repo")

	assert.EqualValues(t, want, session)
}

func TestCheckRunSession_AddLogSection(t *testing.T) {
	session := github.CheckRunSession{}
	execution := &runner.Execution{}
	want := github.CheckRunSession{Sections: []*runner.Execution{
		execution,
	}}

	session.AddLogSection(execution)

	assert.EqualValues(t, want, session)
}

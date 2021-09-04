package github

import (
	"fmt"
	"github.com/dev-this/terraform-gha-controller/internal/runner"
)

// CheckRunSession is a weird amalgamation of states from the GitHub repository, to information
// received from GitHub's Checks API, and everything in-between.
type CheckRunSession struct {
	Name     string
	Sha      string
	Title    string
	Sections []*runner.Execution
	// external identifier GitHub stores against the check run
	ExtID      string
	Repo       string
	Owner      string
	CheckID    int64
	Status     string
	Conclusion string

	TempRepoDir string

	State struct {
		TfInit bool
		TfPlan bool
	}
}

func (c *CheckRunSession) AddLogSection(execution *runner.Execution) {
	c.Sections = append(c.Sections, execution)
}

func (c *CheckRunSession) GetRepositoryURI() string {
	return fmt.Sprintf("https://github.com/%s/%s.git", c.Owner, c.Repo)
}

func (c *CheckRunSession) SetStatus(status string) {
	c.Status = status
}

func (c *CheckRunSession) SetCompletion(status, conclusion string) {
	c.Status = status
	c.Conclusion = conclusion
}

func (c *CheckRunSession) SetTempRepoDir(tempRepoDir string) {
	c.TempRepoDir = tempRepoDir
}

func NewCheckRunSession(name, title, sha, owner, repo, tempRepoDir string) *CheckRunSession {
	return &CheckRunSession{
		Name:        name,
		Title:       title,
		Sha:         sha,
		Repo:        repo,
		Owner:       owner,
		Sections:    []*runner.Execution{},
		TempRepoDir: tempRepoDir,
		State: struct {
			TfInit bool
			TfPlan bool
		}{},
	}
}

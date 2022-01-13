package webhook

import (
	"errors"
	"fmt"
	"github.com/google/go-github/v38/github"
	"log"
)

const (
	successConclusion = "success"
)

var (
	ErrNotDefaultBranch = errors.New("event was not emitted from default branch")
	ErrConclusion       = errors.New("conclusion was not satisfied")
)

type Validator struct {
}

func (v *Validator) ValidatePushEvent(pushEvent github.PushEvent) error {
	defaultBranch := fmt.Sprintf("refs/heads/%s", pushEvent.GetRepo().GetDefaultBranch())

	// ensure pull request was merged to default (primary)  branch, was closed + merged (ie. not just closed)
	// if !pullRequestEvent.GetPullRequest().GetMerged() || pullRequestEvent.GetAction() != "closed" || !isDefaultBranch {
	if pushEvent.GetRef() != defaultBranch {
		log.Println("Skipping webhook")

		return fmt.Errorf(
			"event ref branch %s does not match %s: %w",
			pushEvent.GetRef(),
			defaultBranch,
			ErrNotDefaultBranch,
		)
	}

	return nil
}

func (v *Validator) ValidateWorkflowRunEvent(workflowRunEvent github.WorkflowRunEvent) error {
	eventConclusion := workflowRunEvent.GetWorkflowRun().GetConclusion()

	if eventConclusion != successConclusion {
		return fmt.Errorf(
			"conclusion %s was not successful: %w",
			eventConclusion,
			ErrConclusion,
		)
	}

	return nil
}

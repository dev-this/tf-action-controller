package webhook_test

import (
	"github.com/dev-this/terraform-gha-controller/internal/webhook"
	"github.com/google/go-github/v38/github"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator_ValidatePushEvent(t *testing.T) {
	cases := []struct {
		name  string
		want  error
		event github.PushEvent
	}{
		{
			name: "SucceedsWhenEventRefMatchesDefaultBranch",
			event: github.PushEvent{
				Ref: newString("refs/heads/main"),
				Repo: &github.PushEventRepository{
					DefaultBranch: newString("main"),
				},
			},
		},
		{
			name: "FailsWhenEventRefIsNotDefaultBranch",
			want: webhook.ErrNotDefaultBranch,
			event: github.PushEvent{
				Ref: newString("refs/heads/not-the-default"),
				Repo: &github.PushEventRepository{
					DefaultBranch: newString("main"),
				},
			},
		},
	}

	for _, test := range cases {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			validator := webhook.Validator{}
			got := validator.ValidatePushEvent(tt.event)

			if tt.want != nil {
				assert.Error(t, tt.want, got)
			}
			assert.ErrorIs(t, got, tt.want)
		})
	}
}

func TestValidator_ValidateWorkflowRunEvent(t *testing.T) {
	cases := []struct {
		name  string
		want  error
		event github.WorkflowRunEvent
	}{
		{
			name: "SucceedsWhenConclusionSuccessful",
			event: github.WorkflowRunEvent{
				WorkflowRun: &github.WorkflowRun{
					Conclusion: newString("success"),
				},
			},
		},
		{
			name: "FailsWhenConclusionFailed",
			want: webhook.ErrConclusion,
			event: github.WorkflowRunEvent{
				WorkflowRun: &github.WorkflowRun{
					Conclusion: newString("failed"),
				},
			},
		},
	}

	for _, test := range cases {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			validator := webhook.Validator{}
			got := validator.ValidateWorkflowRunEvent(tt.event)

			if tt.want != nil {
				assert.Error(t, tt.want, got)
			}
			assert.ErrorIs(t, got, tt.want)
		})
	}
}

// newString is dirty helper to get a string pointer.
func newString(v string) *string {
	return &v
}

package github

import (
	"encoding/json"
	"github.com/google/go-github/v38/github"
	"log"
)

type WorkflowRunEvent struct {
	*github.WorkflowRunEvent
}

// ParseWorkflowRunEvent needs generics
func ParseWorkflowRunEvent(payload []byte) (*github.WorkflowRunEvent, error) {
	evt := &github.WorkflowRunEvent{}
	if err := json.Unmarshal(payload, evt); err != nil {
		log.Println("Invalid JSON?", err)

		return evt, err
	}

	return evt, nil
}

// ParsePullRequestEvent needs generics
func ParsePullRequestEvent(payload []byte) (*github.PullRequestEvent, error) {
	evt := &github.PullRequestEvent{}

	if err := json.Unmarshal(payload, evt); err != nil {
		log.Println("Invalid JSON?", err)

		return evt, err
	}

	return evt, nil
}

// ParsePushEvent needs generics
func ParsePushEvent(payload []byte) (*github.PushEvent, error) {
	evt := &github.PushEvent{}

	if err := json.Unmarshal(payload, evt); err != nil {
		log.Println("Invalid JSON?", err)

		return evt, err
	}

	return evt, nil
}

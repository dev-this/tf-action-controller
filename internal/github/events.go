package github

import (
	"encoding/json"
	"github.com/google/go-github/v38/github"
	"log"
)

// ParseWorkflowRunEvent needs generics
func ParseWorkflowRunEvent(payload []byte) (*github.WorkflowRunEvent, error) {
	evt := &github.WorkflowRunEvent{}

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

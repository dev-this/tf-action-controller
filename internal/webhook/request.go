package webhook

import (
	"fmt"
	"gopkg.in/rjz/githubhook.v0"
	"net/http"
)

var (
//ErrParsingRequest = errors.New("failed to parse request payload into GitHub event")
)

func ParseRequest(r *http.Request) (*githubhook.Hook, error) {
	hook, err := githubhook.New(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request payload into GitHub event: %w", err)
	}

	return hook, nil
}

package github

import (
	"context"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/dev-this/terraform-gha-controller/internal/runner"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	githubHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v38/github"
	"log"
	"net/http"
	"time"
)

type Client struct {
	githubClient     *github.Client
	appTokenFetcher  func() string
	checkRunSessions []*CheckRunSession
}

func (gh *Client) GetAppToken() string {
	return gh.appTokenFetcher()
}

func (gh *Client) Clone(gitRepository, clonePath, branch string) error {
	// Clones the repository into the given dir, just as a normal git clone does
	_, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		URL: gitRepository,
		Auth: &githubHttp.BasicAuth{
			Username: "someone", // username doesn't matter for installation generated tokens
			Password: gh.appTokenFetcher(),
		},
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Depth:         1,
	})

	if err != nil {
		return err
	}

	return nil
}

// CreateSession 'check run'.
func (gh *Client) CreateSession(ctx context.Context, session *CheckRunSession) (*CheckRunSession, error) {
	session.ExtID = "abvc123" // RANDOMLY GENERATE ME
	session.Status = "queued"

	startedAt := github.Timestamp{Time: time.Now().UTC()}

	runOptions := github.CreateCheckRunOptions{
		Name:       session.Name,
		HeadSHA:    session.Sha,
		ExternalID: &session.ExtID,
		Status:     &session.Status,
		StartedAt:  &startedAt,
	}

	check, _, err := gh.githubClient.Checks.CreateCheckRun(ctx, session.Owner, session.Repo, runOptions)
	if err != nil {
		return nil, err
	}

	session.CheckID = check.GetID()

	checkRunSession := &CheckRunSession{}
	gh.checkRunSessions = append(gh.checkRunSessions, checkRunSession)

	return checkRunSession, nil
}

func (gh *Client) FlushSession(ctx context.Context, session *CheckRunSession) error {
	// TODO check if this session is tracked by this Client

	details := runner.FormatSections(session.Sections...)
	summary := ""

	runOptions := github.UpdateCheckRunOptions{
		Name:   session.Name,
		Status: &session.Status,
		Output: &github.CheckRunOutput{
			Title:   &session.Title,
			Summary: &summary,
			Text:    &details,
		},
	}

	if session.Conclusion != "" {
		runOptions.Conclusion = &session.Conclusion
	}

	_, _, err := gh.githubClient.Checks.UpdateCheckRun(ctx, session.Owner, session.Repo, session.CheckID, runOptions)

	return err
}

func NewClient(applicationID, installationID int64, privateKey string) *Client {
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, applicationID, installationID, privateKey)
	if err != nil {
		log.Printf("init: failed to retrieve private key at path: %s\n", privateKey)
		log.Fatal(err)

		return nil
	}

	// context for appTokenFetcher
	ctx := context.Background()

	// Use installation transport with client.
	return &Client{
		githubClient: github.NewClient(&http.Client{Transport: itr}),
		appTokenFetcher: func() string {
			token, _ := itr.Token(ctx)

			return token
		},
		checkRunSessions: []*CheckRunSession{},
	}
}

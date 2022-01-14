package webhook

import (
	"context"
	"github.com/dev-this/terraform-gha-controller/internal/github"
	"github.com/dev-this/terraform-gha-controller/internal/runner"
	"io/ioutil"
	"log"
)

func Handler(ctx context.Context, client *github.Client, event Event, options ...HandlerOption) error {
	if event.owner != *GhOwner {
		log.Printf("%s does not match configured GH_OWNER %s - not processing event\n", event.owner, *GhOwner)

		return nil
	}

	// Tempdir to clone the repository
	dir, err := ioutil.TempDir("", "clone-example")
	if err != nil {
		log.Fatal(err)
	}

	tf, err := runner.NewRunner(*tfPath, dir)
	if err != nil {
		log.Fatal(err)
	}
	var name string
	var session *github.CheckRunSession
	var lastStep func() error

	for _, option := range options {
		switch option {
		case WithApply:
			name = "Apply"
			lastStep = func() error {
				section, invoker := tf.Apply(ctx)
				session.AddLogSection(section)

				// run it - flushing session as it progresses, return error if there is one
				return invoker(
					func() {
						_ = client.FlushSession(ctx, session)
					},
				)
			}

		case WithPlan:
			name = "Execution plan"
			lastStep = func() error {
				section, invoker := tf.Plan(ctx)
				session.AddLogSection(section)

				// run it - flushing session as it progresses, return error if there is one
				return invoker(
					func() {
						_ = client.FlushSession(ctx, session)
					},
				)
			}

		default:
			log.Println("failed to recognise webhook type")

			return nil
		}
	}

	session = github.NewCheckRunSession(name, name, event.sha, event.owner, event.repo, dir)

	steps := baseSteps(ctx, client, dir, session, event, tf)
	steps = append(steps, lastStep)

	for i, step := range steps {
		log.Printf("starting step %d\n", i)

		if err = step(); err != nil {
			log.Println("step FAILED", err)
			session.SetCompletion("completed", "failure")
			_ = client.FlushSession(ctx, session)

			return err
		}
	}

	session.SetCompletion("completed", "success")
	if err := client.FlushSession(ctx, session); err != nil {
		log.Println("COMPLETION FAILED", err)

		return err
	}

	return nil
}

func baseSteps(ctx context.Context, client *github.Client, dir string, session *github.CheckRunSession, event Event, tf *runner.TfRunner) []func() error {
	return []func() error{
		// STEP 1: Create Check Run on GitHub API
		func() error {
			if _, err := client.CreateSession(ctx, session); err != nil {
				return err
			}

			return nil
		},

		// STEP 2: Clone repository using App token
		func() error {
			// this should clone based on ref, not branch
			if err := client.Clone(session.GetRepositoryURI(), dir, event.ref); err != nil {
				return err
			}

			session.SetStatus("in_progress")
			return client.FlushSession(ctx, session)
		},

		// STEP 3: Invoke terraform init
		func() error {
			section, invoker := tf.Init(ctx)
			session.AddLogSection(section)

			// run it - flushing session as it progresses, return error if there is one
			return invoker(
				func() {
					_ = client.FlushSession(ctx, session)
				},
			)
		},
	}
}

package main

import (
	"fmt"
	"github.com/dev-this/terraform-gha-controller/internal/github"
	"github.com/dev-this/terraform-gha-controller/internal/webhook"
	"log"
	"net/http"
)

const (
	DefaultPort = "8080"
)

var (
	requiredEnvKeys = []string{"APP_ID", "INSTALLATION_ID", "PRIVATE_KEY", "GH_OWNER"}
)

func main() {
	// Ensure required environment variables exist.
	checkEnvKeys()

	runtimeParams := ParseRuntimeParameters()

	// Pre-checks TODO
	//  [ ] GitHub API Authentication
	//  [ ] GitHub API Permissions
	//  [ ] Filesystem writeable /tmp
	//  [ ] Terraform binary... executable?

	// Prepare GitHub client...
	ghClient := github.NewClient(
		runtimeParams.githubAppID,
		runtimeParams.githubAppInstallationID,
		runtimeParams.githubAppPrivateKey,
	)
	validator := webhook.Validator{}

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got connection: %s", r.Proto)

		hook, err := webhook.ParseRequest(r)
		if err != nil {
			log.Println(err)
			return
		}

		var (
			event   webhook.Event
			options []webhook.HandlerOption
		)

		// get event type from header (eg. X-GitHub-Event: check_suite)
		switch hook.Event {
		case "push":
			pushEvent, err := github.ParsePushEvent(hook.Payload)
			if err != nil {
				log.Println("Invalid JSON?", err)

				return
			}

			if err = validator.ValidatePushEvent(*pushEvent); err != nil {
				log.Println("event did not meet criteria to run", err)

				return
			}

			event = webhook.NewEvent(
				pushEvent.GetHeadCommit().GetID(),
				pushEvent.GetRepo().GetMasterBranch(),
				pushEvent.GetRepo().GetOwner().GetLogin(),
				pushEvent.GetRepo().GetName(),
			)
			options = append(options, webhook.WithApply)

		case "workflow_run":
			log.Println("Found workflow_run")

			workflowRun, err := github.ParseWorkflowRunEvent(hook.Payload)
			if err != nil {
				log.Println("Invalid JSON?", err)

				return
			}

			if err = validator.ValidateWorkflowRunEvent(*workflowRun); err != nil {
				log.Println("event did not meet criteria to run", err)

				return
			}

			event = webhook.NewEvent(
				workflowRun.GetWorkflowRun().GetHeadSHA(),
				workflowRun.GetWorkflowRun().GetHeadBranch(),
				workflowRun.GetWorkflowRun().GetHeadRepository().GetOwner().GetLogin(),
				workflowRun.GetRepo().GetName(),
			)
			options = append(options, webhook.WithPlan)

		default:
			w.WriteHeader(400)
			log.Printf("Unmatched event: %s", hook.Event)
			return
		}

		if err := webhook.Handler(r.Context(), ghClient, event, options...); err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Failed"))

			return
		}

		w.Write([]byte("Thanks!"))
	})

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", runtimeParams.servicePort),
	}

	log.Printf("Serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

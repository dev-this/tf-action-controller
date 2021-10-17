package main

import (
	"fmt"
	"github.com/dev-this/terraform-gha-controller/internal/github"
	"github.com/dev-this/terraform-gha-controller/internal/webhook"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	DefaultPort = "8080"
)

var (
	build          = "development"
	port           = os.Getenv("PORT")
	_              = os.Getenv("GH_SECRET")
	appID          = os.Getenv("APP_ID")
	installationID = os.Getenv("INSTALLATION_ID")
	privateKey     = os.Getenv("PRIVATE_KEY")

	requiredEnvKeys = []string{"APP_ID", "INSTALLATION_ID", "PRIVATE_KEY", "GH_OWNER"}
)

func main() {
	if port == "" {
		port = DefaultPort
	}

	// Ensure required environment variables exist.
	checkEnvKeys()

	// Pre-checks TODO
	//  [ ] GitHub API Authentication
	//  [ ] GitHub API Permissions
	//  [ ] Filesystem writeable /tmp
	//  [ ] Terraform binary... executable?

	applicationID, _ := strconv.ParseInt(appID, 10, 64)
	installationID, _ := strconv.ParseInt(installationID, 10, 64)

	// Prepare GitHub client...
	ghClient := github.NewClient(applicationID, installationID, privateKey)
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
				workflowRun.GetRepo().GetOwner().GetLogin(),
				workflowRun.GetRepo().GetName(),
			)
			options = append(options, webhook.WithPlan)

		default:
			w.WriteHeader(400)
			log.Printf("Unmatched event: %s", hook.Event)
			return
		}

		webhook.Handler(r.Context(), ghClient, event, options...)

		w.Write([]byte("Thanks!"))
	})

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
	}

	log.Printf("Serving on http://0.0.0.0:%s", port)
	log.Fatal(srv.ListenAndServe())
	// log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}

func checkEnvKeys() {
	missingEnvKeys := []string{}

	for _, requiredEnvKey := range requiredEnvKeys {
		if _, ok := os.LookupEnv(requiredEnvKey); !ok {
			missingEnvKeys = append(missingEnvKeys, requiredEnvKey)
		}
	}

	if len(missingEnvKeys) > 0 {
		log.Fatalf("missing defined env vars: %s", missingEnvKeys)
	}
}

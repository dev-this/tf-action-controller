package main

import (
	"fmt"
	github3 "github.com/dev-this/terraform-gha-controller/internal/github"
	"github.com/dev-this/terraform-gha-controller/internal/webhook"
	"gopkg.in/rjz/githubhook.v0"
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
)

func main() {
	if port == "" {
		port = DefaultPort
	}

	// Pre-checks TODO
	//  [ ] GitHub API Authentication
	//  [ ] GitHub API Permissions
	//  [ ] Filesystem writeable /tmp
	//  [ ] Terraform binary... executable?

	applicationID, _ := strconv.ParseInt(appID, 10, 64)
	installationID, _ := strconv.ParseInt(installationID, 10, 64)

	// Prepare GitHub client...
	ghClient := github3.NewClient(applicationID, installationID, privateKey)

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		// Log the request protocol3
		log.Printf("Got connection: %s", r.Proto)

		hook, err := githubhook.New(r)
		if err != nil {
			log.Println("itr.Token FAILED", err)

			return
		}

		var event webhook.Event
		var options []webhook.HandlerOption

		// check envs
		if webhook.GhOwner == nil {
			log.Fatal("Env var GH_OWNER has not been configured")
		}

		switch hook.Event {
		case "push":
			pushEvent, err := github3.ParsePushEvent(hook.Payload)
			if err != nil {
				log.Println("Invalid JSON?", err)

				return
			}

			isDefaultBranch := pushEvent.GetRef() == fmt.Sprintf("refs/heads/%s", pushEvent.GetRepo().GetDefaultBranch())
			head := pushEvent.GetHeadCommit()

			// ensure pull request was merged to default (primary)  branch, was closed + merged (ie. not just closed)
			//if !pullRequestEvent.GetPullRequest().GetMerged() || pullRequestEvent.GetAction() != "closed" || !isDefaultBranch {
			if !isDefaultBranch {
				log.Println("Skipping webhook")

				return
			}

			event = webhook.NewEvent(
				head.GetID(),
				pushEvent.GetRepo().GetMasterBranch(),
				pushEvent.GetPusher().GetLogin(),
				pushEvent.GetRepo().GetName(),
			)
			options = append(options, webhook.WithApply)

		case "workflow_run":
			log.Println("Found workflow_run")

			// get event type from header (eg. X-GitHub-Event: check_suite)
			workflowRun, err := github3.ParseWorkflowRunEvent(hook.Payload)
			if err != nil {
				log.Println("Invalid JSON?", err)

				return
			}

			if workflowRun.GetWorkflowRun().GetConclusion() != "success" {
				log.Println("Skipping cause conclusion was not success")

				// we only want to trigger on success
				return
			}

			event = webhook.NewEvent(
				workflowRun.GetWorkflowRun().GetHeadSHA(),
				workflowRun.GetWorkflowRun().GetHeadBranch(),
				workflowRun.GetOrg().GetLogin(),
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

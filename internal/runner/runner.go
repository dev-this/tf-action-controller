package runner

import (
	"context"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	"time"
)

type TfRunner struct {
	runner *tfexec.Terraform
}

type StreamingOutputCallback func(execution *Execution)

func (r *TfRunner) Apply(ctx context.Context) (*Execution, func(func()) error) {
	return r.newExecution(func(r *TfRunner) error {
		err := r.runner.Apply(ctx) // , tfexec.Parallelism(4))

		return err
	})
}

func (r *TfRunner) Init(ctx context.Context) (*Execution, func(func()) error) {
	return r.newExecution(func(r *TfRunner) error {
		return r.runner.Init(ctx, tfexec.Upgrade(false))
	})
}

func (r *TfRunner) Plan(ctx context.Context) (*Execution, func(func()) error) {
	return r.newExecution(func(r *TfRunner) error {
		_, err := r.runner.Plan(ctx)

		return err
	})
}

func (r *TfRunner) setupWriters(cb StdoutCallback) {
	writer := NewWriter(cb)
	r.runner.SetStdout(writer)
	r.runner.SetStderr(writer)
}
func (r *TfRunner) newExecution(action func(r *TfRunner) error) (*Execution, func(func()) error) {
	section := &Execution{
		Details:      "",
		ErrorDetails: "",
		Completed:    false,
		Successful:   false,
		StartedAt:    time.Now(),
		SecondsRan:   0,
	}

	return section, func(cb func()) error {
		wrappedCallback := func(chunk string) {
			section.AppendToDetails(chunk)
			cb()
		}

		r.setupWriters(wrappedCallback)

		err := action(r)

		section.MarkCompleted(err == nil)

		return err
	}
}

func NewRunner(binPath, workingDir string) (*TfRunner, error) {
	_, err := tfinstall.Find(context.Background(), tfinstall.ExactPath(binPath))

	if err != nil {
		return nil, err
	}

	tf, err := tfexec.NewTerraform(workingDir, binPath)

	if err != nil {
		return nil, err
	}

	return &TfRunner{
		runner: tf,
	}, nil
}

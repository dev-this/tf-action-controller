package runner

import (
	"context"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
)

type TfRunner struct {
	binPath string
}

func (r *TfRunner) Plan(ctx context.Context, cwd string, cb StdoutCallback) error {
	tf, err := tfexec.NewTerraform(cwd, r.binPath)

	if err != nil {
		return err
	}

	writer := NewWriter(cb)
	tf.SetStdout(writer)

	_, err = tf.Plan(ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewRunner(path string) *TfRunner {
	_, err := tfinstall.Find(context.Background(), tfinstall.ExactPath(path))

	if err != nil {
		panic(err)
	}

	return &TfRunner{
		binPath: path,
	}
}

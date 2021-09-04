package runner_test

import (
	"github.com/dev-this/terraform-gha-controller/internal/runner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecution_AppendToDetails(t *testing.T) {
	tests := []struct {
		name    string
		details []string
		want    runner.Execution
	}{
		{
			name:    "Single string",
			details: []string{"some new deets"},
			want:    runner.Execution{Details: "some new deets"},
		},
		{
			name:    "MultiStringsWithNewline",
			details: []string{"firstly we did X/Y/Z\n", "then some of this..."},
			want:    runner.Execution{Details: "firstly we did X/Y/Z\nthen some of this..."},
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			execution := runner.Execution{}

			for _, detail := range tt.details {
				execution.AppendToDetails(detail)
			}

			assert.EqualValues(t, tt.want, execution)
		})
	}
}

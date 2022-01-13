package runner_test

import (
	"github.com/dev-this/terraform-gha-controller/internal/runner"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

// TestOutputWriter_WriteInvokesCallback asserts callback passed to the Writer gets invoked, and in the right order.
func TestOutputWriter_WriteInvokesCallback(t *testing.T) {
	wantBufferedContents := []string{
		"This is the first line",
		"Hopefully this is the subsequent line.",
	}
	gotBufferedContents := []string{}
	// callback will append to gotBufferedContents
	callback := func(chunk string) {
		gotBufferedContents = append(gotBufferedContents, chunk)
	}
	writer := runner.NewWriter(callback)

	_, _ = writer.Write([]byte("This is the first line"))
	_, _ = writer.Write([]byte("Hopefully this is the subsequent line."))

	assert.Equal(t, wantBufferedContents, gotBufferedContents)
}

// TestOutputWriter_Write will assert the returned value make sense.
func TestOutputWriter_Write(t *testing.T) {
	wantWrittenBytes := 51 // 51 bytes = size of ascii characters in contents string.
	contents := []byte("This is a really cool test, don't try this at home.")
	callback := func(chunk string) {}
	writer := runner.NewWriter(callback)

	gotWrittenBytes, err := writer.Write(contents)

	assert.NoError(t, err)
	assert.Equal(t, wantWrittenBytes, gotWrittenBytes)
}

// TestNewWriter assets (indirectly)
func TestNewWriterImplementsWriterInterface(t *testing.T) {
	callback := func(chunk string) {}
	writer := runner.NewWriter(callback)

	assert.Implements(t, (*io.Writer)(nil), writer)
}

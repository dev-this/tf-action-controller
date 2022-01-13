package runner

import (
	"bytes"
	"io"
)

// Ensure interface is implemented
var _ io.Writer = (*OutputWriter)(nil)

// IOCallback is the signature for compatible callbacks with OutputWriter
type IOCallback func(chunk string)

// OutputWriter is a dumb IO writer that will parse bytes to a string, passing it to a callback.
type OutputWriter struct {
	callback IOCallback
}

// Write is method to make IO Writer interface happy.
func (w *OutputWriter) Write(p []byte) (n int, err error) {
	buf := bytes.NewBuffer(p)

	w.callback(buf.String())

	return len(p), nil
}

func NewWriter(callback IOCallback) *OutputWriter {
	return &OutputWriter{callback}
}

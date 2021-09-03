package runner

import (
	"bytes"
)

type Writer struct {
	callback StdoutCallback
}

func (w *Writer) Write(p []byte) (n int, err error) {
	buf := bytes.NewBuffer(p)

	w.callback(buf.String())

	return len(p), nil
}

func NewWriter(callback StdoutCallback) *Writer {
	return &Writer{callback}
}

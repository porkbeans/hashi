package ioutils

import (
	"io"
	"os"
)

// Remove removes specified file with suppressing error.
func Remove(filename string) {
	if len(filename) > 0 {
		_ = os.Remove(filename)
	}
}

// Close closes specified file with suppressing error.
func Close(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
}

// ErrorWriter helps error handling while writing.
type ErrorWriter struct {
	writer io.Writer
	err    error
}

// NewErrorWriter creates a ErrorWriter instance.
func NewErrorWriter(w io.Writer, err error) ErrorWriter {
	return ErrorWriter{
		writer: w,
		err:    err,
	}
}

// Err returns error if any error occurred in Write method.
func (w ErrorWriter) Err() error {
	return w.err
}

// Write writes only if no error occurred.
func (w ErrorWriter) Write(p []byte) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}

	if n, err = w.writer.Write(p); err != nil {
		w.err = err
	}
	return
}

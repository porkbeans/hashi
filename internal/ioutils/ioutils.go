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

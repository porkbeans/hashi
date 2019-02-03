package ioutils

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"
)

func TestErrWriter1(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("failed to create a filename")
	}
	defer Remove(file.Name())
	defer Close(file)

	w := NewErrorWriter(file, err)
	_, err = w.Write([]byte("content"))
	if err != nil {
		t.Error(err)
	}
}

func TestErrWriter2(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	dummyError := errors.New("dummy error")
	w := NewErrorWriter(buf, dummyError)

	_, err := w.Write([]byte("content"))
	if err.Error() != dummyError.Error() {
		t.Errorf("expected dummy error but got '%s'", err.Error())
	}

	if w.Err().Error() != dummyError.Error() {
		t.Errorf("expected dummy error but got '%s'", err.Error())
	}
}

type FailWriter struct {
	Err error
}

func (w FailWriter) Write(p []byte) (n int, err error) {
	return 0, w.Err
}

func TestErrWriter3(t *testing.T) {
	dummyError := errors.New("dummy error")
	w := NewErrorWriter(FailWriter{Err: dummyError}, nil)
	_, err := w.Write([]byte("content"))
	if err.Error() != dummyError.Error() {
		t.Errorf("expected dummy error but got '%s'", err.Error())
	}
}

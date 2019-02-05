package ioutils

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestRemoveClose(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	Close(file)
	Remove(file.Name())

	if err = file.Close(); err == nil {
		t.Errorf("file should be already closed")
	}

	if _, err = os.Stat(file.Name()); os.IsExist(err) {
		t.Errorf("file should be removed")
	}
}

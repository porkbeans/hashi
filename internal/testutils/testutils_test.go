package testutils

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/porkbeans/hashi/internal/ioutils"
)

func TestGenerateInvalidURL(t *testing.T) {
	if _, err := url.Parse(GenerateInvalidURL()); err == nil {
		t.Fatalf("must fail to parse")
	}
}

func TestTestServerHandler_ServeHTTP(t *testing.T) {
	server := httptest.NewServer(
		TestServerHandler{
			StatusCode: 200,
			Content:    "Hello",
		},
	)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error should not happen")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("status code must be 200")
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error should not happen")
	}

	if string(content) != "Hello" {
		t.Fatalf("content must be 'Hello'")
	}
}

func TestFailBodyHttpClient_Get(t *testing.T) {
	dummyError := errors.New("dummy error")
	c := FailBodyHTTPClient{Err: dummyError}
	r, err := c.Get("")
	if err != nil {
		t.Fatalf("failed to get dummy response")
	}
	defer ioutils.Close(r.Body)

	n, err := r.Body.Read([]byte{})
	if n != 0 || err != dummyError {
		t.Errorf("expected %s but got %s", dummyError, err)
	}
}

func TestTouchTempFile(t *testing.T) {
	filename, err := TouchTempFile()
	if err != nil {
		t.Fatalf("error should not happen")
	}
	defer os.Remove(filename)

	_, err = os.Stat(filename)
	if err != nil {
		t.Fatalf("file must exist")
	}
}

func TestCreateTempZip(t *testing.T) {
	filename, err := CreateTempZip("testexe", "Hello")
	if err != nil {
		t.Fatalf("error should not happen")
	}
	defer os.Remove(filename)

	zipReader, err := zip.OpenReader(filename)
	if err != nil {
		t.Fatalf("error should not happen")
	}
	defer zipReader.Close()

	if len(zipReader.File) != 1 {
		t.Fatalf("zip must contain only 1 file")
	}

	if zipReader.File[0].Name != "testexe" {
		t.Fatalf("error should not happen")
	}

	reader, err := zipReader.File[0].Open()
	if err != nil {
		t.Fatalf("error should not happen")
	}
	defer reader.Close()

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("error should not happen")
	}

	if string(content) != "Hello" {
		t.Fatalf("content must be 'Hello'")
	}
}

func TestCreateTempZip2(t *testing.T) {
	filename, err := CreateTempZip(strings.Repeat("a", math.MaxUint16+1), "test")
	if err == nil {
		defer os.Remove(filename)
		t.Error(err)
	}
}

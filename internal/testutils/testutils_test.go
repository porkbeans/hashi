package testutils

import (
	"archive/zip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
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

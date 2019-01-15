package testutils

import (
	"archive/zip"
	"crypto/rand"
	"io/ioutil"
	"net/http"
	"net/url"
)

func generateRandomBytes() string {
	buf := [32]byte{}
	rand.Read(buf[:])
	return string(buf[:])
}

// GenerateInvalidURL provides random invalid URL
func GenerateInvalidURL() string {
	var rawURL string
	var err error

	for err == nil {
		rawURL = generateRandomBytes()
		_, err = url.Parse(rawURL)
	}

	return rawURL
}

// TestServerHandler is a handler that returns specified contents for tests
type TestServerHandler struct {
	StatusCode int
	Content    string
}

// ServeHTTP returns specified contents in TestServerHandler
func (h TestServerHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(h.StatusCode)
	writer.Header().Add("Content-Type", "text/plain")
	writer.Write([]byte(h.Content))
}

// TouchTempFile creates a temporary file
func TouchTempFile() (string, error) {
	file, err := ioutil.TempFile("", "hashi-test-")
	defer file.Close()
	return file.Name(), err
}

// CreateTempZip creates a zip file that contains a file
func CreateTempZip(filenameInZip, content string) (string, error) {
	zipFile, err := ioutil.TempFile("", "hashi-test-")
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	fileWriter, err := zipWriter.Create(filenameInZip)
	if err != nil {
		return "", err
	}

	_, err = fileWriter.Write([]byte(content))
	if err != nil {
		return "", err
	}

	return zipFile.Name(), nil
}

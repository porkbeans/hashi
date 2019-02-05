package testutils

import (
	"archive/zip"
	"crypto/rand"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/porkbeans/hashi/internal/ioutils"
)

// GenerateInvalidURL provides random invalid URL
func GenerateInvalidURL() string {
	var rawURL string

	h := "ghijklmnopqrstuvwxyz"
	rands := [8]byte{}
	_, _ = rand.Read(rands[:])
	for i := 0; i < 4; i++ {
		h1 := string(h[int(rands[i*2])%len(h)])
		h2 := string(h[int(rands[i*2+1])%len(h)])
		rawURL += "%" + h1 + h2
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
	_, _ = writer.Write([]byte(h.Content))
}

type failReadCloser struct {
	Err error
}

func (w failReadCloser) Read(p []byte) (n int, err error) {
	return 0, w.Err
}

func (w failReadCloser) Close() (err error) {
	return nil
}

type failWriter struct {
	Err error
}

func (w failWriter) Write(p []byte) (n int, err error) {
	return 0, w.Err
}

// FailBodyHTTPClient helps simulating failure on reading response body.
type FailBodyHTTPClient struct {
	Err error
}

// Get returns response which returns error while reading.
func (c *FailBodyHTTPClient) Get(url string) (resp *http.Response, err error) {
	resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       failReadCloser{Err: c.Err},
	}
	return
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

// TouchTempFile creates a temporary file
func TouchTempFile() (string, error) {
	file, err := ioutil.TempFile("", "hashi-test-")
	defer ioutils.Close(file)
	return file.Name(), err
}

// CreateTempZip creates a zip file that contains a file
func CreateTempZip(filenameInZip, content string) (string, error) {
	tempFile, err := ioutil.TempFile("", "hashi-test-")
	zipFile := NewErrorWriter(tempFile, err)
	defer ioutils.Close(tempFile)

	zipWriter := zip.NewWriter(zipFile)
	defer ioutils.Close(zipWriter)

	fileWriter := NewErrorWriter(zipWriter.Create(filenameInZip))
	_, _ = fileWriter.Write([]byte(content))

	if fileWriter.Err() != nil {
		return "", fileWriter.Err()
	}

	return tempFile.Name(), nil
}

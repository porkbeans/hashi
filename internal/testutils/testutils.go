package testutils

import (
	"archive/zip"
	"crypto/rand"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GenerateRandomBytes() string {
	buf := [32]byte{}
	rand.Read(buf[:])
	return string(buf[:])
}

func GenerateInvalidURL() string {
	var rawURL string
	var err error

	for err == nil {
		rawURL = GenerateRandomBytes()
		_, err = url.Parse(rawURL)
	}

	return rawURL
}

type TestServerHandler struct {
	StatusCode int
	Content    string
}

func (h TestServerHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(h.StatusCode)
	writer.Header().Add("Content-Type", "text/plain")
	writer.Write([]byte(h.Content))
}

func TouchTempFile() (string, error) {
	file, err := ioutil.TempFile("", "hashi-test-")
	defer file.Close()
	return file.Name(), err
}

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

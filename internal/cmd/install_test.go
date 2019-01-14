package cmd

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"github.com/porkbeans/hashi/internal/testutils"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestDownloadToTempFile(t *testing.T) {
	content := "Hello"

	server := http.Server{
		Addr: "localhost:8989",
		Handler: testutils.TestServerHandler{
			StatusCode: 200,
			Content:    content,
		},
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Log(err)
		}
	}()

	tempFileName, checksum, err := downloadToTempFile("http://localhost:8989/", os.Stderr)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tempFileName)

	expected := sha256.Sum256([]byte(content))
	if !bytes.Equal(expected[:], checksum[:]) {
		t.Errorf("checksum failed")
	}

	server.Shutdown(context.Background())
}

func TestDownloadToTempFileFail(t *testing.T) {
	tempFileName, _, err := downloadToTempFile(testutils.GenerateInvalidURL(), os.Stderr)
	if err == nil {
		defer os.Remove(tempFileName)
		t.Error("error must happen")
	}
}

func TestGetChecksum(t *testing.T) {
	_, err := getChecksum("consul", "1.4.0", "linux", "amd64")
	if err != nil {
		t.Errorf("error should not happen")
	}

	_, err = getChecksum("unknown", "1.4.0", "linux", "amd64")
	if err == nil {
		t.Errorf("error must happen")
	}

	_, err = getChecksum("consul", "1.4.0", "unknown", "unknown")
	if err == nil {
		t.Errorf("error must happen")
	}
}

func TestOpenFileInZip(t *testing.T) {
	tempFileName, err := testutils.CreateTempZip("testexe", "hello")
	if err != nil {
		t.Errorf("error should not happen")
	}
	defer os.Remove(tempFileName)

	zipReader, err := zip.OpenReader(tempFileName)
	if err != nil {
		t.Errorf("error should not happen")
	}
	defer zipReader.Close()

	fileReader, _, err := openFileInZip(zipReader, "testexe")
	if err != nil {
		t.Errorf("error should not happen")
	}
	defer fileReader.Close()

	content, err := ioutil.ReadAll(fileReader)
	if err != nil {
		t.Errorf("error should not happen")
	}

	if string(content) != "hello" {
		t.Errorf("content must be 'Hello'")
	}
}

func TestOpenFileInZipNotFound(t *testing.T) {
	tempFileName, err := testutils.CreateTempZip("testexe", "hello")
	if err != nil {
		t.Errorf("error should not happen")
	}
	defer os.Remove(tempFileName)

	zipReader, err := zip.OpenReader(tempFileName)
	if err != nil {
		t.Errorf("error should not happen")
	}
	defer zipReader.Close()

	fileReader, _, err := openFileInZip(zipReader, "unknown")
	if err == nil {
		defer fileReader.Close()
		t.Errorf("error should not happen")
	}
}

func TestExtractBinaryInZip(t *testing.T) {
	tempZipName, err := testutils.CreateTempZip("testexe", "hello")
	if err != nil {
		t.Errorf("error should not happen")
	}
	defer os.Remove(tempZipName)

	tempBinName, err := testutils.TouchTempFile()
	if err != nil {
		t.Errorf("error should not happen")
	}
	defer os.Remove(tempBinName)

	err = extractBinaryInZip(tempBinName, tempZipName, "testexe", os.Stderr)
	if err != nil {
		t.Errorf("error should not happen")
	}
}

func TestExtractBinaryInZipNotFound1(t *testing.T) {
	tempZipName, err := testutils.TouchTempFile()
	if err != nil {
		t.Errorf("error should not happen")
	}
	os.Remove(tempZipName)

	tempBinName, err := testutils.TouchTempFile()
	if err != nil {
		t.Errorf("error should not happen")
	}
	defer os.Remove(tempBinName)

	err = extractBinaryInZip(tempBinName, tempZipName, "testexe", os.Stderr)
	if err == nil {
		t.Errorf("error must happen")
	}
}

func TestExtractBinaryInZipNotFound2(t *testing.T) {
	tempZipName, err := testutils.CreateTempZip("testexe", "hello")
	if err != nil {
		t.Errorf("error should not happen")
	}
	defer os.Remove(tempZipName)

	tempBinName, err := testutils.TouchTempFile()
	if err != nil {
		t.Errorf("error should not happen")
	}
	defer os.Remove(tempBinName)

	err = extractBinaryInZip(tempBinName, tempZipName, "unknown", os.Stderr)
	if err == nil {
		t.Errorf("error must happen")
	}
}

func TestInstallCmd(t *testing.T) {
	tempBinName, err := testutils.TouchTempFile()
	if err != nil {
		t.Errorf("error should not happen")
	}
	defer os.Remove(tempBinName)

	err = installCmd.RunE(installCmd, []string{"consul", "1.4.0", tempBinName})
	if err != nil {
		t.Errorf("error should not happen")
	}
}

func TestInstallCmdNotFound(t *testing.T) {
	tempBinName, err := testutils.TouchTempFile()
	if err != nil {
		t.Errorf("error should not happen")
	}
	defer os.Remove(tempBinName)

	err = installCmd.RunE(installCmd, []string{"unknown", "1.4.0", tempBinName})
	if err == nil {
		t.Errorf("error must happen")
	}
}
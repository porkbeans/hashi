package cmd

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/porkbeans/hashi/internal/ioutils"
	"github.com/porkbeans/hashi/internal/testutils"
)

func TestDownloadToTempFile(t *testing.T) {
	content := "Hello"

	server := httptest.NewServer(
		testutils.TestServerHandler{
			StatusCode: 200,
			Content:    content,
		},
	)
	defer server.Close()

	tempFileName, checksum, err := downloadToTempFile(server.URL, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}
	defer ioutils.Remove(tempFileName)

	expected := sha256.Sum256([]byte(content))
	if !bytes.Equal(expected[:], checksum[:]) {
		t.Errorf("checksum failed")
	}
}

func TestDownloadToTempFileFail(t *testing.T) {
	tempFileName, _, err := downloadToTempFile(testutils.GenerateInvalidURL(), os.Stderr)
	if err == nil {
		defer ioutils.Remove(tempFileName)
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
		t.Fatalf("error should not happen")
	}
	defer ioutils.Remove(tempFileName)

	zipReader, err := zip.OpenReader(tempFileName)
	if err != nil {
		t.Fatalf("error should not happen")
	}
	defer ioutils.Close(zipReader)

	fileReader, _, err := openFileInZip(zipReader, "testexe")
	if err != nil {
		t.Fatalf("error should not happen")
	}
	defer ioutils.Close(fileReader)

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
		t.Fatalf("error should not happen")
	}
	defer ioutils.Remove(tempFileName)

	zipReader, err := zip.OpenReader(tempFileName)
	if err != nil {
		t.Fatalf("error should not happen")
	}
	defer ioutils.Close(zipReader)

	fileReader, _, err := openFileInZip(zipReader, "unknown")
	if err == nil {
		defer ioutils.Close(fileReader)
		t.Errorf("error should not happen")
	}
}

func TestExtractBinaryInZip(t *testing.T) {
	tempZipName, err := testutils.CreateTempZip("testexe", "hello")
	if err != nil {
		t.Fatalf("error should not happen")
	}
	defer ioutils.Remove(tempZipName)

	tempBinName, err := testutils.TouchTempFile()
	if err != nil {
		t.Fatalf("error should not happen")
	}
	defer ioutils.Remove(tempBinName)

	err = extractBinaryInZip(tempBinName, tempZipName, "testexe", os.Stderr)
	if err != nil {
		t.Errorf("error should not happen")
	}
}

func TestExtractBinaryInZipNotFound1(t *testing.T) {
	tempZipName, err := testutils.TouchTempFile()
	if err != nil {
		t.Fatalf("error should not happen")
	}
	ioutils.Remove(tempZipName)

	tempBinName, err := testutils.TouchTempFile()
	if err != nil {
		t.Fatalf("error should not happen")
	}
	defer ioutils.Remove(tempBinName)

	err = extractBinaryInZip(tempBinName, tempZipName, "testexe", os.Stderr)
	if err == nil {
		t.Errorf("error must happen")
	}
}

func TestExtractBinaryInZipNotFound2(t *testing.T) {
	tempZipName, err := testutils.CreateTempZip("testexe", "hello")
	if err != nil {
		t.Fatalf("error should not happen")
	}
	defer ioutils.Remove(tempZipName)

	tempBinName, err := testutils.TouchTempFile()
	if err != nil {
		t.Fatalf("error should not happen")
	}
	defer ioutils.Remove(tempBinName)

	err = extractBinaryInZip(tempBinName, tempZipName, "unknown", os.Stderr)
	if err == nil {
		t.Errorf("error must happen")
	}
}

func TestInstallCmd(t *testing.T) {
	tempBinName, err := testutils.TouchTempFile()
	if err != nil {
		t.Fatalf("error should not happen")
	}
	defer ioutils.Remove(tempBinName)

	err = installCmd.RunE(installCmd, []string{"consul", "1.4.0", tempBinName})
	if err != nil {
		t.Errorf("error should not happen")
	}
}

func TestInstallCmdNotFound(t *testing.T) {
	tempBinName, err := testutils.TouchTempFile()
	if err != nil {
		t.Fatalf("error should not happen")
	}
	defer ioutils.Remove(tempBinName)

	err = installCmd.RunE(installCmd, []string{"unknown", "1.4.0", tempBinName})
	if err == nil {
		t.Errorf("error must happen")
	}
}

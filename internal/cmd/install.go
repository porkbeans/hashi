package cmd

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/porkbeans/hashi/internal/httputils"

	"github.com/mitchellh/ioprogress"
	"github.com/porkbeans/hashi/pkg/parseutils"
	"github.com/porkbeans/hashi/pkg/urlutils"
	"github.com/spf13/cobra"
)

var (
	targetGOOS   string
	targetGOARCH string
)

func progressReader(reader io.Reader, size int64, printer io.Writer, prefix string) *ioprogress.Reader {
	return &ioprogress.Reader{
		Reader:       reader,
		Size:         size,
		DrawInterval: 100 * time.Millisecond,
		DrawFunc: ioprogress.DrawTerminalf(printer, func(progress int64, total int64) string {
			return fmt.Sprintf("%s %s\r", prefix, ioprogress.DrawTextFormatBytes(progress, total))
		}),
	}
}

func downloadToTempFile(url string, printer io.Writer) (string, [32]byte, error) {
	checksum := [32]byte{}

	resp, err := httputils.Get(nil, url)
	if err != nil {
		return "", checksum, err
	}
	defer resp.Body.Close()

	tempFile, err := ioutil.TempFile("", "hashi-")
	if err != nil {
		return "", checksum, err
	}
	defer tempFile.Close()

	hash := sha256.New()
	tee := io.TeeReader(resp.Body, hash)

	_, err = io.Copy(tempFile, progressReader(tee, resp.ContentLength, printer, "Downloading..."))
	if err != nil {
		defer os.Remove(tempFile.Name())
		return "", checksum, err
	}

	copy(checksum[:], hash.Sum(nil)[0:32])

	return tempFile.Name(), checksum, nil
}

func getChecksum(product, version, goos, goarch string) ([32]byte, error) {
	url := urlutils.ProductZipChecksumURL(product, version)

	resp, err := httputils.Get(nil, url)
	if err != nil {
		return [32]byte{}, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return [32]byte{}, err
	}

	checksums := parseutils.ParseChecksumList(string(content))
	for _, checksum := range checksums {
		if checksum.Name == product && checksum.Version == version && checksum.Os == goos && checksum.Arch == goarch {
			return checksum.Checksum, nil
		}
	}

	return [32]byte{}, errors.New("checksum not found")
}

func openFileInZip(zipReader *zip.ReadCloser, filename string) (io.ReadCloser, *zip.File, error) {
	for _, file := range zipReader.File {
		if file.Name == filename {
			reader, err := file.Open()
			return reader, file, err
		}
	}

	return nil, nil, fmt.Errorf("%s not found in zip", filename)
}

func extractBinaryInZip(dst string, src string, filenameInZip string, printer io.Writer) error {
	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	rawReader, file, err := openFileInZip(zipReader, filenameInZip)
	if err != nil {
		return err
	}
	defer rawReader.Close()

	fileReader := progressReader(rawReader, int64(file.UncompressedSize64), printer, "Extracting...")

	fileWriter, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0755))
	if err != nil {
		return err
	}
	defer fileWriter.Close()

	_, err = io.Copy(fileWriter, fileReader)
	return err
}

var installCmd = &cobra.Command{
	Use:   "install <name> <version> <path>",
	Short: "Install HashiCorp tools.",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		product := args[0]
		version := args[1]
		goos := targetGOOS
		goarch := targetGOARCH
		installPath := args[2]

		zipURL := urlutils.ProductZipURL(product, version, goos, goarch)
		cmd.Printf("Retrieve %s\n", zipURL)

		tempFileName, actualChecksum, err := downloadToTempFile(zipURL, cmd.OutOrStderr())
		if err != nil {
			return err
		}
		defer os.Remove(tempFileName)

		expectedChecksum, err := getChecksum(product, version, goos, goarch)
		if err != nil {
			return err
		}

		if !bytes.Equal(expectedChecksum[:], actualChecksum[:]) {
			return errors.New("checksum failed")
		}
		cmd.Println("Checksum Passed")

		if err := extractBinaryInZip(installPath, tempFileName, product, cmd.OutOrStderr()); err != nil {
			return err
		}

		cmd.Printf("Installed %s successfully to %s\n", product, installPath)
		return nil
	},
}

func init() {
	installCmd.Flags().StringVarP(&targetGOOS, "os", "o", runtime.GOOS, "operating system")
	installCmd.Flags().StringVarP(&targetGOARCH, "arch", "a", runtime.GOARCH, "architecture")
}

package cmd

import (
	"bytes"
	"context"
	"github.com/porkbeans/hashi/internal/testutils"
	"github.com/porkbeans/hashi/pkg/urlutils"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

func TestGetList(t *testing.T) {
	_, err := getList(urlutils.HashicorpProductList)
	if err != nil {
		t.Errorf("error should not happen")
	}
}

func TestGetListInvalidURL(t *testing.T) {
	_, err := getList(testutils.GenerateInvalidURL())
	if err == nil {
		t.Errorf("error must happen")
	}
}

func TestGetListServerError(t *testing.T) {
	server := http.Server{
		Addr: "localhost:8989",
		Handler: testutils.TestServerHandler{
			StatusCode: 500,
			Content:    "Failed",
		},
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Log(err)
		}
	}()

	_, err := getList("http://localhost:8989/")
	if err == nil {
		t.Errorf("error must happen")
	}

	server.Shutdown(context.Background())
}

func TestShowList(t *testing.T) {
	buf := bytes.Buffer{}
	listCmd.SetOutput(&buf)
	listCmd.RunE(listCmd, nil)

	productNames := map[string]bool{}
	for _, productName := range strings.Split(buf.String(), "\n") {
		if len(productName) > 0 {
			productNames[productName] = true
		}
	}

	expectedProductNames := []string{
		"consul",
		"nomad",
		"terraform",
		"packer",
		"vagrant",
		"vault",
	}
	for _, expectedProductName := range expectedProductNames {
		if !productNames[expectedProductName] {
			t.Errorf("product list must contain %s", expectedProductName)
		}
	}
}

func TestShowVersionList(t *testing.T) {
	buf := bytes.Buffer{}
	listCmd.SetOutput(&buf)
	listCmd.RunE(listCmd, []string{"consul"})

	pattern := regexp.MustCompile(`\d+\.\d+\.\d+`)
	for _, version := range strings.Split(buf.String(), "\n") {
		if len(version) > 0 && !pattern.MatchString(version) {
			t.Errorf("%s doesn't match version format", version)
		}
	}
}

func TestShowZipList(t *testing.T) {
	buf := bytes.Buffer{}
	listCmd.SetOutput(&buf)
	listCmd.RunE(listCmd, []string{"consul", "1.4.0"})

	pattern := regexp.MustCompile(`(\w+) (\w+) \*?`)
	for _, version := range strings.Split(buf.String(), "\n") {
		if len(version) > 0 && !pattern.MatchString(version) {
			t.Errorf("%s doesn't match version format", version)
		}
	}
}

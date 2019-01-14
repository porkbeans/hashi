package ioutils

import (
	"context"
	"github.com/porkbeans/hashi/internal/testutils"
	"github.com/porkbeans/hashi/pkg/urlutils"
	"net/http"
	"testing"
)

func TestGetOK(t *testing.T) {
	_, err := Get(urlutils.HashicorpProductList)
	if err != nil {
		t.Errorf("error should not happen")
	}
}

func TestGetNil(t *testing.T) {
	_, err := Get("")
	if err == nil {
		t.Errorf("error must happen")
	}
}

func TestGetNotFound(t *testing.T) {
	_, err := Get(urlutils.HashicorpProductList + "nonexistence")
	if err == nil {
		t.Errorf("error must happen")
	}
}

func TestGetOther(t *testing.T) {
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

	_, err := Get("http://localhost:8989/")
	if err == nil {
		t.Errorf("error must happen")
	}

	server.Shutdown(context.Background())
}

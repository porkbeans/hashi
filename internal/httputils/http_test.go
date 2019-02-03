package httputils

import (
	"net/http/httptest"
	"testing"

	"github.com/porkbeans/hashi/internal/testutils"
	"github.com/porkbeans/hashi/pkg/urlutils"
)

func TestGetOK(t *testing.T) {
	_, err := Get(nil, urlutils.HashicorpProductList)
	if err != nil {
		t.Errorf("error should not happen")
	}
}

func TestGetNil(t *testing.T) {
	_, err := Get(nil, "")
	if err == nil {
		t.Errorf("error must happen")
	}
}

func TestGetNotFound(t *testing.T) {
	_, err := Get(nil, urlutils.HashicorpProductList+"nonexistence")
	if err == nil {
		t.Errorf("error must happen")
	}
}

func TestGetOther(t *testing.T) {
	server := httptest.NewServer(
		testutils.TestServerHandler{
			StatusCode: 500,
			Content:    "Failed",
		},
	)
	defer server.Close()

	client := server.Client()
	_, err := Get(client, server.URL)
	if err == nil {
		t.Errorf("error must happen")
	}
}

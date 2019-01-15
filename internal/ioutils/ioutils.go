package ioutils

import (
	"fmt"
	"net/http"
)

// Get retrieves resources from specified URL. returns error if status code is not 200.
func Get(client *http.Client, url string) (*http.Response, error) {
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 200:
		break
	case 403:
		return nil, fmt.Errorf("%s not found", url)
	default:
		return nil, fmt.Errorf("failed to get %s (status: %d)", url, resp.StatusCode)
	}

	return resp, nil
}

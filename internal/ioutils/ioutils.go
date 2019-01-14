package ioutils

import (
	"fmt"
	"net/http"
)

func Get(url string) (*http.Response, error) {
	resp, err := http.Get(url)
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

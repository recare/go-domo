package domo

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type payload struct {
	contentType string
	data        io.Reader
}

func (p *payload) getData() io.Reader {
	if p == nil {
		return nil
	}
	return p.data
}

func url(parts ...string) string {
	return strings.Join(parts, "/")
}

func (c HttpClient) do(method, path string, payload *payload) ([]byte, error) {
	token, err := c.auth("data")
	if err != nil {
		return nil, fmt.Errorf("error fetching oauth access_token - %w", err)
	}

	req, err := http.NewRequest(method, url(baseUrl, version, path), payload.getData())
	if err != nil {
		return nil, fmt.Errorf("error building request - %w", err)
	}

	req.Header.Add("Authorization", "bearer "+token.AccessToken)
	if payload != nil {
		req.Header.Add("Content-Type", payload.contentType)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing %s %s - %w", method, path, err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body of pdf generation request - %w", err)
	}

	if !isSuccess(res.StatusCode) {
		return nil, fmt.Errorf("expected status 2XX performing %s %s, got %d %s - %s", method, path, res.StatusCode, res.Status, string(body))
	}

	return body, nil
}

func isSuccess(status int) bool {
	return status >= 200 && status <= 299
}

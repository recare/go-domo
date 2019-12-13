package domo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strings"
	"time"
)

type token struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	ExpiresAt   time.Time
}

func (c *HttpClient) auth(grants ...string) (*token, error) {
	if c.token != nil && c.token.ExpiresAt.Before(time.Now()) {
		return c.token, nil
	}

	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s?grant_type=client_credentials&scope=%s",
			url(baseUrl, pathAuth),
			neturl.QueryEscape(strings.Join(grants, " "))), nil)
	if err != nil {
		return nil, fmt.Errorf("error building request - %w", err)
	}
	req.SetBasicAuth(c.clientID, c.secret)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting oauth access_token - %w", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body of pdf generation request - %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status 2XX getting oauth access_token, got %d %s - %s", res.StatusCode, res.Status, string(body))
	}

	var token *token
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("error deserializing access_token - %w", err)
	}
	token.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn-1) * time.Second)
	c.token = token

	return token, nil
}

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cli/go-gh"
)

type customClient struct {
	*http.Client
	token string
}

func (c *customClient) Request(method string, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, "https://api.github.com/"+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")

	return c.Do(req)
}

type GitHubClient struct {
	restClient interface {
		Request(method string, path string, body io.Reader) (*http.Response, error)
	}
}

func NewGitHubClient(token string) *GitHubClient {
	var restClient interface {
		Request(method string, path string, body io.Reader) (*http.Response, error)
	}

	if token == "" {
		ghClient, err := gh.RESTClient(nil)
		if err != nil {
			return nil
		}
		restClient = ghClient
	} else {
		restClient = &customClient{
			Client: http.DefaultClient,
			token:  token,
		}
	}

	return &GitHubClient{
		restClient: restClient,
	}
}

func (gh *GitHubClient) Get(ctx context.Context, uri string, result interface{}) error {
	if gh == nil || gh.restClient == nil {
		return fmt.Errorf("client not properly initialized")
	}

	resp, err := gh.restClient.Request("GET", uri, nil)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	return nil
}

func (gh *GitHubClient) Post(ctx context.Context, uri string, result interface{}) error {
	if gh == nil || gh.restClient == nil {
		return fmt.Errorf("client not properly initialized")
	}

	resp, err := gh.restClient.Request("POST", uri, nil)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	return nil
}

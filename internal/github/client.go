package github

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GitHubClient struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

func NewGitHubClient(token string) *GitHubClient {
	if token == "" {
		log.Fatalf("token is not set (use GH_TOKEN env)")
	}
	return &GitHubClient{
		BaseURL: "https://api.github.com/",
		Token:   token,
		Client:  &http.Client{},
	}
}

func (gh *GitHubClient) Get(ctx context.Context, uri string, result interface{}) error {
	apiURL := gh.BaseURL + uri
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+gh.Token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := gh.Client.Do(req)
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
	apiURL := gh.BaseURL + uri
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+gh.Token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := gh.Client.Do(req)
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

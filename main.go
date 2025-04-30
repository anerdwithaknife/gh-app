package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type AppDetails struct {
	AppId    int    `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	ClientID string `json:"client_id"`
	// Add other fields as needed
}

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

func main() {
	ctx := context.Background()
	token := os.Getenv("GH_TOKEN")
	client := NewGitHubClient(token)

	var appDetails AppDetails
	if err := client.Get(ctx, "apps/four-wards", &appDetails); err != nil {
		log.Fatalf("Error fetching app details: %v", err)
	}

	log.Printf("App ID: %d", appDetails.AppId)
	log.Printf("App Name: %s", appDetails.Name)
	log.Printf("App Slug: %s", appDetails.Slug)
	log.Printf("Client ID: %s", appDetails.ClientID)
}

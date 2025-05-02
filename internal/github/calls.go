package github

import (
	"context"
	"fmt"
	"os"
)

func GetAppDetails(slug string) (*AppDetails, error) {
	ctx := context.Background()
	token := os.Getenv("GH_TOKEN")
	client := NewGitHubClient(token)

	var appDetails AppDetails
	if err := client.Get(ctx, fmt.Sprintf("apps/%s", slug), &appDetails); err != nil {
		return nil, err
	}

	return &appDetails, nil
}

func GetPrivateKey(privateKeyFile string) (string, error) {
	if privateKeyFile == "" {
		return "", fmt.Errorf("path to private key file must not be empty")
	}
	privateKey, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return "", fmt.Errorf("error reading private key file: %w", err)
	}
	return string(privateKey), nil
}

func GetAppInstallations(jwtToken string, appId int) ([]AppInstallation, error) {
	ctx := context.Background()
	client := NewGitHubClient(jwtToken)

	var appInstallations []AppInstallation
	if err := client.Get(ctx, "app/installations", &appInstallations); err != nil {
		return nil, fmt.Errorf("error fetching app installations: %w", err)
	}

	return appInstallations, nil
}

func GenerateAccessToken(jwtToken string, appId int, installationId string) (string, error) {
	ctx := context.Background()
	client := NewGitHubClient(jwtToken)

	var appToken AppToken
	if err := client.Post(ctx, fmt.Sprintf("app/installations/%s/access_tokens", installationId), &appToken); err != nil {
		return "", fmt.Errorf("error generating app token: %w", err)
	}

	return appToken.Token, nil
}

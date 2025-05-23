package github

import (
	"context"
	"fmt"
	"log"
	"os"
)

func TestGithubApi() {
	ctx := context.Background()
	token := os.Getenv("GH_TOKEN")
	client := NewGitHubClient(token)

	var appDetails AppDetails
	if err := client.Get(ctx, "apps/four-wards-public", &appDetails); err != nil {
		log.Fatalf("Error fetching app details: %v", err)
	}

	log.Printf("App Id: %d", appDetails.AppId)
	log.Printf("App Name: %s", appDetails.Name)
	log.Printf("App Slug: %s", appDetails.Slug)
	log.Printf("Client Id: %s", appDetails.ClientId)

	privateKey, err := GetPrivateKey(os.Getenv("GH_APP_PRIVATE_KEY_FILE"))
	if err != nil {
		log.Fatalf("Get private key file: %v", err)
	}

	jwtToken, err := GenerateGithubAppJWT(appDetails.AppId, privateKey)
	if err != nil {
		log.Fatalf("Error generating JWT: %v", err)
	}
	log.Printf("Generated JWT: %s", jwtToken)

	appClient := NewGitHubClient(jwtToken)
	var appInstallations []AppInstallation
	if err := appClient.Get(ctx, "app/installations", &appInstallations); err != nil {
		log.Fatalf("Error fetching app details: %v", err)
	}

	log.Printf("App Installations: %+v", appInstallations)

	if len(appInstallations) == 0 {
		log.Fatalf("No installations found for the app")
	}

	installationId := appInstallations[0].Id
	log.Printf("Installation Id: %d", installationId)
	var appToken AppToken
	if err := appClient.Post(ctx, fmt.Sprintf("app/installations/%d/access_tokens", installationId), &appToken); err != nil {
		log.Fatalf("Error generating app token: %v", err)
	}
	log.Printf("Generated App Token: %s", appToken.Token)
	log.Printf("App Token Expires At: %s", appToken.ExpiresAt)
	log.Printf("App Token Permissions: %+v", appToken.Permissions)
}

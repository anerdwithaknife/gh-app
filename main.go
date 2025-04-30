package main

import (
	"context"
	"log"
	"net/url"
	"os"

	"github.com/google/go-github/v71/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()

	token := os.Getenv("GH_TOKEN")
	if token == "" {
		log.Fatalf("GH_TOKEN environment variable is not set")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	var err error
	client.BaseURL, err = url.Parse("http://localhost:8080/")
	if err != nil {
		log.Fatalf("Error parsing base URL: %v", err)
	}

	data, _, err := client.Apps.Get(ctx, "four-wards")
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	log.Printf("%v", data)
}

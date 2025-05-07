package cmd

import (
	"fmt"
	"log"

	"github.com/cursethevulgar/gh-app/internal/github"
	"github.com/cursethevulgar/gh-app/internal/store"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generates an access token for a given app installation",
	Long: `Generates an access token for a given app installation.

Requires an app slug and an installation id, use gh app installations 
to see available installation ids for a given app slug.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := store.NewDefaultStore(false)
		if err != nil {
			log.Fatalf("Error loading store: %v", err)
			return
		}

		slug, _ := cmd.Flags().GetString("slug")
		if slug == "" {
			log.Fatal("Slug must not be empty")
		}

		installationID, _ := cmd.Flags().GetString("installation-id")
		if installationID == "" {
			log.Fatal("Installation ID must not be empty")
		}

		app, err := db.GetAppBySlug(slug)
		if err != nil {
			log.Fatalf("Error getting app details: %v", err)
		}

		jwtToken, err := github.GenerateGithubAppJWT(app.AppID, app.PrivateKey)
		if err != nil {
			log.Fatalf("Error generating JWT: %v", err)
		}

		accessToken, err := github.GenerateAccessToken(jwtToken, app.AppID, installationID)
		if err != nil {
			log.Fatalf("Error generating access token: %v", err)
		}

		fmt.Println(accessToken)
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)

	tokenCmd.Flags().StringP("slug", "s", "", "The slug of the app")
	tokenCmd.Flags().StringP("installation-id", "i", "", "The associated installation id")

	tokenCmd.MarkFlagRequired("slug")
	tokenCmd.MarkFlagRequired("installation-id")
}

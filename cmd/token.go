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
		orgName, _ := cmd.Flags().GetString("org-name")

		if (installationID == "" && orgName == "") || (installationID != "" && orgName != "") {
			log.Fatal("You must provide either --installation-id (-i) or --org-name (-o)")
		}

		app, err := db.GetAppBySlug(slug)
		if err != nil {
			log.Fatalf("Error: Problem reading yaml store: %v", err)
		}

		if app == nil {
			log.Fatal("Error: App was not found in yaml store")
		}

		jwtToken, err := github.GenerateGithubAppJWT(app.AppID, app.PrivateKey)
		if err != nil {
			log.Fatalf("Error: Problem generating JWT: %v", err)
		}

		if installationID == "" {
			installations, err := github.GetAppInstallations(jwtToken, app.AppID)
			if err != nil {
				log.Fatalf("Error: Problem fetching installations: %v", err)
			}
			for _, installation := range installations {
				if installation.Account.Login == orgName {
					installationID = fmt.Sprintf("%d", installation.Id)
					break
				}
			}
			if installationID == "" {
				log.Fatalf("Error: No installation found for org: %s", orgName)
			}
		}

		accessToken, err := github.GenerateAccessToken(jwtToken, app.AppID, installationID)
		if err != nil {
			log.Fatalf("Error: Problem generating access token: %v", err)
		}

		fmt.Println(accessToken)
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)

	tokenCmd.Flags().StringP("slug", "s", "", "The slug of the app")
	tokenCmd.Flags().StringP("installation-id", "i", "", "The associated installation id")
	tokenCmd.Flags().StringP("org-name", "o", "", "The org to retrieve the installation id from")

	tokenCmd.MarkFlagRequired("slug")
}

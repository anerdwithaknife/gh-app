package cmd

import (
	"log"

	"github.com/cursethevulgar/gh-app/internal/github"
	"github.com/cursethevulgar/gh-app/internal/store"

	"github.com/spf13/cobra"
)

var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "Generate JWT for GitHub App",
	Long: `Generates a signed JWT token for the specified saved app.

The token can be used for calling the GitHub API /app endpoints.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := store.NewDefaultStore()
		if err != nil {
			log.Println("Error loading store:", err)
			return
		}
		slug, _ := cmd.Flags().GetString("slug")
		if slug == "" {
			log.Println("Slug must not be empty")
			return
		}

		app, err := db.GetAppBySlug(slug)
		if err != nil {
			log.Println("Error getting app details:", err)
			return
		}

		jwtToken, err := github.GenerateGithubAppJWT(app.AppID, app.PrivateKey)
		if err != nil {
			log.Fatalf("Error generating JWT: %v", err)
		}
		cmd.Println(jwtToken)
	},
}

func init() {
	rootCmd.AddCommand(jwtCmd)

	jwtCmd.Flags().StringP("slug", "s", "", "The slug of the app to show save for")
	jwtCmd.MarkFlagRequired("slug")
}

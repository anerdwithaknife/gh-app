package cmd

import (
	"log"

	"github.com/anerdwithaknife/gh-app/internal/github"
	"github.com/anerdwithaknife/gh-app/internal/store"

	"github.com/spf13/cobra"
)

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save a GitHub app to the local database",
	Long: `Saves a GitHub app with the specified slug to the local database.

If no app id is specified, the app details are fetched from GitHub API using GH_TOKEN.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := store.NewDefaultStore(true)
		if err != nil {
			log.Println("Error loading store:", err)
			return
		}

		appDetails := &github.AppDetails{}

		slug, _ := cmd.Flags().GetString("slug")

		appId, _ := cmd.Flags().GetInt("app-id")
		if appId != 0 {
			appDetails.AppId = appId
			appDetails.Slug = slug
			appDetails.Name = slug
		} else {
			cmd.Printf("Fetching details for app with slug: %s\n", slug)

			var err error
			appDetails, err = github.GetAppDetails(slug)
			if err != nil {
				cmd.Println("Error getting app details:", err)
				return
			}
		}

		privateKeyPath, _ := cmd.Flags().GetString("private-key")
		privateKey, err := github.GetPrivateKey(privateKeyPath)
		if err != nil {
			cmd.Println("Private key error:", err)
			return
		}

		clientSecret, _ := cmd.Flags().GetString("client-secret")

		app := store.App{
			Slug:         slug,
			Name:         appDetails.Name,
			AppID:        appDetails.AppId,
			ClientID:     appDetails.ClientId,
			ClientSecret: clientSecret,
			PrivateKey:   privateKey,
		}

		if err := db.SaveApp(&app); err != nil {
			cmd.Println("Error storing app:", err)
			return
		}

		cmd.Printf("App %s saved to local yaml store\n", slug)
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)

	saveCmd.Flags().StringP("slug", "s", "", "The slug of the app to show save for")
	saveCmd.Flags().StringP("private-key", "p", "", "Path to private key (*.pem) of the app")
	saveCmd.Flags().StringP("client-secret", "e", "", "Client secret of the app")
	saveCmd.Flags().IntP("app-id", "a", 0, "When specifying the app id, no details are fetched from GitHub")

	saveCmd.MarkFlagRequired("slug")
}

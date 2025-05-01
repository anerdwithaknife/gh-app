package cmd

import (
	"fmt"

	"gh-app/internal/github"
	"gh-app/internal/store"

	"github.com/spf13/cobra"
)

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Show all save for a given app",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		db := store.Store{}
		if err := db.Init(); err != nil {
			fmt.Println("Error initializing store:", err)
			return
		}

		slug, _ := cmd.Flags().GetString("slug")
		if slug == "" {
			fmt.Println("Slug must not be empty")
			return
		}

		fmt.Printf("Fetching details for app with slug: %s\n", slug)

		appDetails, err := github.GetAppDetails(slug)
		if err != nil {
			fmt.Println("Error getting app details:", err)
			return
		}

		privateKeyPath, _ := cmd.Flags().GetString("private-key")
		privateKey, err := github.GetPrivateKey(privateKeyPath)
		if err != nil {
			fmt.Println("Private key error:", err)
			return
		}

		app := store.App{
			Slug:       slug,
			Name:       appDetails.Name,
			AppID:      appDetails.AppId,
			ClientID:   appDetails.ClientId,
			PrivateKey: privateKey,
		}

		if err := db.SaveApp(&app); err != nil {
			fmt.Println("Error storing app:", err)
			return
		}

		fmt.Printf("App %s saved to local database\n", slug)
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)

	saveCmd.Flags().StringP("slug", "s", "", "The slug of the app to show save for")
	saveCmd.Flags().StringP("private-key", "p", "", "Path to private key (*.pem) of the app")

	saveCmd.MarkFlagRequired("slug")
}

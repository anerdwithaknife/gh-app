package cmd

import (
	"fmt"

	"gh-app/internal/lab"
	"gh-app/internal/store"

	"github.com/spf13/cobra"
)

var installationsCmd = &cobra.Command{
	Use:   "installations",
	Short: "Show all installations for a given app",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		slug, err := cmd.Flags().GetString("slug")
		if err != nil {
			fmt.Println("Error getting slug flag:", err)
			return
		}
		if slug == "" {
			fmt.Println("Slug is required")
			return
		}
		fmt.Printf("Showing installations for app with slug: %s\n", slug)

		lab.GetAppInstallations(slug)

		app := store.App{
			Slug: slug,
		}
	},
}

func init() {
	rootCmd.AddCommand(installationsCmd)

	installationsCmd.Flags().StringP("slug", "s", "", "The slug of the app to show installations for")
	installationsCmd.MarkFlagRequired("slug")
}

package cmd

import (
	"github.com/spf13/cobra"

	"gh-app/internal/store"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved apps",
	Long:  `Displays a list of all GitHub apps saved to the local data store.`,
	Run: func(cmd *cobra.Command, args []string) {
		db := store.Store{}
		if err := db.Init(); err != nil {
			cmd.Println("Error initializing store:", err)
			return
		}

		apps, err := db.GetAllApps()
		if err != nil {
			cmd.Println("Error getting apps:", err)
			return
		}

		for _, app := range apps {
			cmd.Printf("%d\t%s\n", app.AppID, app.Slug)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

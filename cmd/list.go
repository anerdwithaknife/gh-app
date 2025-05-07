package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/cursethevulgar/gh-app/internal/store"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all saved apps",
	Long:    `Displays a list of all GitHub apps saved to the local data store.`,
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		db, err := store.NewDefaultStore(false)
		if err != nil {
			log.Println("Error loading store:", err)
			return
		}

		apps, err := db.GetAllApps()
		if err != nil {
			cmd.Println("Error getting apps:", err)
			return
		}

		drawAppTable(apps)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func drawAppTable(apps []*store.App) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("App ID", "Slug", "Client ID")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, app := range apps {
		tbl.AddRow(app.AppID, app.Slug, app.ClientID)
	}

	tbl.Print()
}

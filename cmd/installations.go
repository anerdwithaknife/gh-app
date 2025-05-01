package cmd

import (
	"gh-app/internal/github"
	"gh-app/internal/store"
	"log"

	"github.com/fatih/color"
	"github.com/rodaine/table"
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
		db := store.Store{}
		if err := db.Init(); err != nil {
			log.Println("Error initializing store:", err)
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

		installations, err := github.GetAppInstallations(jwtToken, app.AppID)
		if err != nil {
			log.Fatalf("Error getting installations: %v", err)
		}

		drawAppInstallationsTable(installations)
	},
}

func init() {
	rootCmd.AddCommand(installationsCmd)

	installationsCmd.Flags().StringP("slug", "s", "", "The slug of the app to show installations for")
	installationsCmd.MarkFlagRequired("slug")
}

func drawAppInstallationsTable(installations []github.AppInstallation) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Installation ID", "Account Login", "Target Type")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, installation := range installations {
		tbl.AddRow(installation.Id, installation.Account.Login, installation.TargetType)
	}

	tbl.Print()
}

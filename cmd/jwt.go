package cmd

import (
	"gh-app/internal/lab"

	"github.com/spf13/cobra"
)

var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "Generate JWT for GitHub App",
	Long: `Generates a signed JWT token using APP_ID and APP_PRIVATE_KEY from environment.

The token can be used for calling the GitHub API /app endpoints.`,
	Run: func(cmd *cobra.Command, args []string) {
		lab.TestGithubApi()
	},
}

func init() {
	rootCmd.AddCommand(jwtCmd)
}

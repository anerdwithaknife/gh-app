package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gh-app",
	Short: "Manage GitHub apps",
	Long:  `Manage GitHub apps from the command line for convenient access.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	commonTemplate := `
{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}
{{end}}
USAGE
  {{.UseLine}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}

{{if .HasAvailableSubCommands}}COMMANDS
{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}  {{rpad .Name .NamePadding }} {{.Short}}
{{end}}{{end}}
{{end}}{{if .HasAvailableLocalFlags}}FLAGS
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}
{{end}}{{if .HasAvailableInheritedFlags}}GLOBAL FLAGS
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}
{{end}}
`
	rootCmd.SetUsageTemplate(commonTemplate)
	rootCmd.SetHelpTemplate(commonTemplate)
}

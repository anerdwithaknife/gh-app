/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// jwtCmd represents the jwt command
var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "Generate JWT for GitHub App",
	Long: `Generates a signed JWT token using APP_ID and APP_PRIVATE_KEY from environment.

The token can be used for calling the GitHub API /app endpoints.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("jwt called")
	},
}

func init() {
	rootCmd.AddCommand(jwtCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jwtCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jwtCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

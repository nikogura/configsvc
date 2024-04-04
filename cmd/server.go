/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/nikogura/configsvc/pkg/configsvc"
	"github.com/spf13/cobra"
	"log"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the HTTP service.",
	Long: `
Run the HTTP service.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := configsvc.Server(address, path)
		if err != nil {
			log.Fatalf("Error running server: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

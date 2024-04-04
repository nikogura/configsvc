package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var address string
var path string
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "configsvc",
	Short: "Serve JSON files over a network",
	Long:  "Serve JSON files over a network",
	//Run: func(cmd *cobra.Command, args []string) {
	//},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&address, "address", "a", "0.0.0.0:8888", "address on which to run (IP and Port)")
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "/opt/data", "path to look for data to serve")
}

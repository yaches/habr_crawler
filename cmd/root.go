package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(crawlCommand)
	rootCmd.AddCommand(infoCommand)
	rootCmd.AddCommand(histCommand)
	rootCmd.AddCommand(hubsCommand)
	rootCmd.AddCommand(tagsCommand)
}

var rootCmd = &cobra.Command{}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

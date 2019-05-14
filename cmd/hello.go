package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var helloCommand = &cobra.Command{
	Use: "hello",
	Run: func(cmd *cobra.Command, argv []string) {
		log.Printf("hello")
	},
}

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yaches/habr_crawler/state"
	"github.com/yaches/habr_crawler/tasks"
)

var helloCommand = &cobra.Command{
	Use: "hello",
	Run: func(cmd *cobra.Command, argv []string) {
		s := state.NewStorageRedis()
		defer s.Close()
		s.Add(tasks.Task{tasks.UserTask, "azaza", 0, 19})
	},
}

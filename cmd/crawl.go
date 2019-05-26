package cmd

import (
	"github.com/yaches/habr_crawler/content"
	"github.com/yaches/habr_crawler/crawler"
	"github.com/yaches/habr_crawler/state"
	"github.com/yaches/habr_crawler/tasks"

	"github.com/spf13/cobra"
)

var crawlCommand = &cobra.Command{
	Use: "crawl",
	Run: func(cmd *cobra.Command, argv []string) {
		cntStorage := content.NewStorageNative()
		stateStorage := state.NewStorageNative()
		queue := tasks.NewTaskManagerChan()
		worker := crawler.NewWorker(cntStorage, stateStorage, queue)
		queue.Push([]tasks.Task{tasks.Task{Type: tasks.PostTask, Body: "197598"}})
		worker.Work(5)
	},
}

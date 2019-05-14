package cmd

import (
	"github.com/yaches/habr_crawler/contentstore"
	"github.com/yaches/habr_crawler/crawler"
	"github.com/yaches/habr_crawler/statestore"
	"github.com/yaches/habr_crawler/tasks"

	"github.com/spf13/cobra"
)

var crawlCommand = &cobra.Command{
	Use: "crawl",
	Run: func(cmd *cobra.Command, argv []string) {
		cntStorage := contentstore.NewStorageNative()
		stateStorage := statestore.NewStorageNative()
		queue := tasks.NewTaskManagerChan()
		queue.Push([]tasks.Task{tasks.Task{Type: tasks.PostTask, Body: "451812"}})
		crawler.Work(&cntStorage, &stateStorage, &queue)
	},
}

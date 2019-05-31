package cmd

import (
	"github.com/go-redis/redis"
	"github.com/yaches/habr_crawler/content"
	"github.com/yaches/habr_crawler/crawler"
	"github.com/yaches/habr_crawler/state"
	"github.com/yaches/habr_crawler/tasks"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

var (
	threads int
	maxDeep int
)

func init() {
	crawlCommand.Flags().IntVarP(&threads, "threads", "t", 1, "Threads")
	crawlCommand.Flags().IntVarP(&maxDeep, "deep", "d", 1, "Max deep")
}

var crawlCommand = &cobra.Command{
	Use: "crawl",
	Run: func(cmd *cobra.Command, argv []string) {
		db := redis.NewClient(&redis.Options{})

		cntStorage, err := content.NewStorageES()
		if err != nil {
			zap.L().Fatal(err.Error())
		}
		stateStorage := state.NewStorageRedis(db)
		queue := tasks.NewManagerRedis(db)
		worker := crawler.NewWorker(cntStorage, stateStorage, queue)
		// queue.Channel() <- tasks.Task{Type: tasks.PostTask, Body: "453626"}
		queue.Push(tasks.Task{Type: tasks.PostTask, Body: "453596"})

		worker.Start(threads, maxDeep)

		// users, _ := cntStorage.GetAllUsers()
		// fmt.Println(users)
	},
}

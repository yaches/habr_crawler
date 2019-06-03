package cmd

import (
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/yaches/habr_crawler/content"
	"github.com/yaches/habr_crawler/crawler"
	"github.com/yaches/habr_crawler/tasks"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"

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

		dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
		if err != nil {
			zap.L().Fatal("Can't connect to the proxy", zap.Error(err))
		}
		httpTransport := &http.Transport{Dial: dialer.Dial}
		httpClient := &http.Client{Transport: httpTransport}
		// httpClient := &http.Client{}

		db := redis.NewClient(&redis.Options{})

		cntStorage, err := content.NewStorageES()
		if err != nil {
			zap.L().Fatal(err.Error())
		}
		queue := tasks.NewManagerRedis(db)
		worker := crawler.NewWorker(cntStorage, queue, httpClient)
		// err = queue.Push(tasks.Task{Type: tasks.PostTask, Body: "460000"})
		err = queue.Fill()
		if err != nil {
			zap.L().Fatal(err.Error())
		}
		zap.L().Info("Got new tasks", zap.Int("count", len(queue.Channel())))
		// postsIterate(queue)

		worker.Start(threads, maxDeep)
	},
}

func postsIterate(queue tasks.Manager) {
	zap.L().Info("Posts iterating started")
	for i := 0; i < 460000; i++ {
		task := tasks.Task{
			Type: tasks.PostTask,
			Body: strconv.Itoa(i),
		}
		zap.L().Info("Pushing", zap.Int("i", i))
		err := queue.Push(task)
		if err != nil {
			zap.L().Warn("", zap.Error(err))
		}
	}
	zap.L().Info("Posts iterating stopped")
}

package cmd

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/cobra"
	"github.com/yaches/habr_crawler/tasks"
	"go.uber.org/zap"
)

var helloCommand = &cobra.Command{
	Use: "hello",
	Run: func(cmd *cobra.Command, argv []string) {
		db := redis.NewClient(&redis.Options{})

		queue := tasks.NewManagerRedis(db)
		err := queue.Fill()
		if err != nil {
			zap.L().Debug("", zap.Error(err))
		}

		b := true
		for b {
			select {
			case t := <-queue.Channel():
				zap.L().Debug(fmt.Sprint(t))
			default:
				b = false
			}
		}
	},
}

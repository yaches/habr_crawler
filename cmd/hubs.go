package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yaches/habr_crawler/content"
	"github.com/yaches/habr_crawler/printer"
	"go.uber.org/zap"
)

var (
	usernameHubs string
)

func init() {
	hubsCommand.Flags().StringVarP(&usernameHubs, "username", "u", "", "Username")
}

var hubsCommand = &cobra.Command{
	Use: "hubs",
	Run: func(cmd *cobra.Command, argv []string) {
		cntStorage, err := content.NewStorageES()
		if err != nil {
			zap.L().Fatal("", zap.Error(err))
		}

		var hist map[string]int
		if usernameHubs == "" {
			hist, err = cntStorage.GetHubCommonHist("Hubs")
		} else {
			hist, err = cntStorage.GetHubUserHist(usernameHubs, "Hubs")
		}
		if err != nil {
			zap.L().Fatal("", zap.Error(err))
		}
		printer.PrintHubHist(hist)
	},
}

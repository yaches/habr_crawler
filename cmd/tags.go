package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yaches/habr_crawler/content"
	"github.com/yaches/habr_crawler/printer"
	"go.uber.org/zap"
)

var (
	usernameTags string
)

func init() {
	tagsCommand.Flags().StringVarP(&usernameTags, "username", "u", "", "Username")
}

var tagsCommand = &cobra.Command{
	Use: "tags",
	Run: func(cmd *cobra.Command, argv []string) {
		cntStorage, err := content.NewStorageES()
		if err != nil {
			zap.L().Fatal("", zap.Error(err))
		}

		var hist map[string]int
		if usernameTags == "" {
			hist, err = cntStorage.GetHubCommonHist("Tags")
		} else {
			hist, err = cntStorage.GetHubUserHist(usernameTags, "Tags")
		}
		if err != nil {
			zap.L().Fatal("", zap.Error(err))
		}
		printer.PrintHubHist(hist)
	},
}

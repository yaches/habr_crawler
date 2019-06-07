package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yaches/habr_crawler/content"
	"github.com/yaches/habr_crawler/printer"
	"go.uber.org/zap"
)

var (
	user  string
	index string
	gran  string
)

func init() {
	histCommand.Flags().StringVarP(&user, "user", "u", "", "Habr user")
	histCommand.Flags().StringVarP(&index, "index", "i", "posts", "Data index")
	histCommand.Flags().StringVarP(&gran, "gran", "g", "hour", "Hist granularity")
}

var histCommand = &cobra.Command{
	Use: "hist",
	Run: func(cmd *cobra.Command, argv []string) {
		cntStorage, err := content.NewStorageES()
		if err != nil {
			zap.L().Fatal("", zap.Error(err))
		}
		var m map[int]int
		if user == "" {
			m, err = cntStorage.GetCommonHist(index, gran)
		} else {
			m, err = cntStorage.GetTermFilteredHist(index, gran, "Author", user)
		}
		if err != nil {
			zap.L().Fatal("", zap.Error(err))
		}
		printer.PrintHist(index, gran, m)
	},
}

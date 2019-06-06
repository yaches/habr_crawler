package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yaches/habr_crawler/content"
	"github.com/yaches/habr_crawler/printer"
	"go.uber.org/zap"
)

var (
	user string
)

func init() {
	infoCommand.Flags().StringVarP(&user, "user", "u", "", "Habr user")
}

var infoCommand = &cobra.Command{
	Use: "info",
	Run: func(cmd *cobra.Command, argv []string) {
		cntStorage, err := content.NewStorageES()
		if err != nil {
			zap.L().Fatal("", zap.Error(err))
		}
		i, err := cntStorage.GetCommonInfo()
		if err != nil {
			zap.L().Fatal("", zap.Error(err))
		}

		printer.PrintCommons(i)
	},
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yaches/habr_crawler/content"
	"github.com/yaches/habr_crawler/printer"
	"go.uber.org/zap"
)

var (
	infoIndex string
	postID    string
	username  string
)

func init() {
	infoCommand.Flags().StringVarP(&infoIndex, "index", "i", "posts", "Data index")
	infoCommand.Flags().StringVarP(&postID, "post", "p", "1", "Post ID")
	infoCommand.Flags().StringVarP(&username, "username", "u", "bobuk", "Username")
}

var infoCommand = &cobra.Command{
	Use: "info",
	Run: func(cmd *cobra.Command, argv []string) {
		cntStorage, err := content.NewStorageES()
		if err != nil {
			zap.L().Fatal("", zap.Error(err))
		}

		if infoIndex == "users" {
			u, err := cntStorage.GetUser(username)
			if err != nil {
				zap.L().Fatal("", zap.Error(err))
			}
			printer.PrintUser(u)
		}
		if infoIndex == "posts" {
			p, err := cntStorage.GetPost(postID)
			if err != nil {
				zap.L().Fatal("", zap.Error(err))
			}
			fmt.Printf("%+v\n", p)
		}
	},
}

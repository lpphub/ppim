package main

import (
	"github.com/spf13/cobra"
	"ppim/internal/logic"
)

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "main",
		Short: "beacon-biz Main Function",
		Run: func(cmd *cobra.Command, args []string) {
			// 默认web服务
			logic.HttpServe()
		},
	}

	// 添加子命令
	//rootCmd.AddCommand(testCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}

func main() {
	Execute()
}

package main

import (
	"github.com/spf13/cobra"
	"ppim/internal/comet"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "main",
		Short: "Main Function",
		Run: func(cmd *cobra.Command, args []string) {
			comet.Serve()
		},
	}

	// 添加子命令
	//rootCmd.AddCommand(testCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}

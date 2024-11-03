package main

import (
	"github.com/spf13/cobra"
	"ppim/internal/logic"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "main",
		Short: "Main Function",
		Run: func(cmd *cobra.Command, args []string) {
			logic.Serve()
		},
	}

	// 添加子命令
	//rootCmd.AddCommand(testCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}

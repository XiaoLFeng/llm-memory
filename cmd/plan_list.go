package cmd

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

// planListCmd 列出所有计划
// 呀~ 查看所有计划！✨
var planListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有计划",
	Long:  `列出系统中的所有计划~ 📋`,
	Run: func(cmd *cobra.Command, args []string) {
		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewPlanHandler(bs)
		if err := handler.List(bs.Context()); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	planCmd.AddCommand(planListCmd)
}

package todo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/internal/models/dto"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var (
	todoBatchCreateJSON     string
	todoBatchCreateJSONFile string
)

// todoBatchCreateCmd 批量创建待办
var todoBatchCreateCmd = &cobra.Command{
	Use:   "batch-create",
	Short: "批量创建待办事项",
	Long:  `批量创建多个待办事项~ ✨`,
	Example: `  # 使用 JSON 字符串批量创建
  llm-memory todo batch-create --json '[
    {"code":"todo-1","title":"任务1","priority":3},
    {"code":"todo-2","title":"任务2","description":"详细描述"}
  ]'

  # 使用 JSON 文件批量创建
  llm-memory todo batch-create --json-file ./todos.json

  # JSON 文件格式示例 (todos.json):
  [
    {
      "code": "todo-task-one",
      "title": "第一个任务",
      "description": "任务详细说明",
      "priority": 3
    },
    {
      "code": "todo-task-two",
      "title": "第二个任务",
      "priority": 2
    }
  ]`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析 JSON 输入
		items, err := parseCreateItems(todoBatchCreateJSON, todoBatchCreateJSONFile)
		if err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}

		if len(items) == 0 {
			cli.PrintError("未提供有效的待办事项")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewTodoHandler(bs)
		if err := handler.BatchCreate(bs.Context(), items); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	todoBatchCreateCmd.Flags().StringVar(&todoBatchCreateJSON, "json", "", "JSON格式的待办列表")
	todoBatchCreateCmd.Flags().StringVar(&todoBatchCreateJSONFile, "json-file", "", "包含待办列表的JSON文件路径")

	todoCmd.AddCommand(todoBatchCreateCmd)
}

// parseCreateItems 解析 JSON 输入为待办创建项列表
func parseCreateItems(jsonStr string, jsonFile string) ([]dto.ToDoCreateDTO, error) {
	var data []byte
	var err error

	// 优先使用文件
	if jsonFile != "" {
		data, err = os.ReadFile(jsonFile)
		if err != nil {
			return nil, fmt.Errorf("读取文件失败: %w", err)
		}
	} else if jsonStr != "" {
		data = []byte(jsonStr)
	} else {
		return nil, fmt.Errorf("必须提供 --json 或 --json-file 参数")
	}

	var items []dto.ToDoCreateDTO
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("JSON 格式错误: %w", err)
	}

	return items, nil
}

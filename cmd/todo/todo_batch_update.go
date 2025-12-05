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
	todoBatchUpdateJSON     string
	todoBatchUpdateJSONFile string
)

// todoBatchUpdateCmd 批量更新待办
var todoBatchUpdateCmd = &cobra.Command{
	Use:   "batch-update",
	Short: "批量更新待办事项",
	Long:  `批量更新多个待办事项的标题、描述、优先级或状态~ 📝`,
	Example: `  # 使用 JSON 字符串批量更新
  llm-memory todo batch-update --json '[
    {"code":"todo-1","title":"新标题","priority":4},
    {"code":"todo-2","status":2}
  ]'

  # 使用 JSON 文件批量更新
  llm-memory todo batch-update --json-file ./updates.json

  # JSON 文件格式示例 (updates.json):
  [
    {
      "code": "todo-task-one",
      "title": "更新后的标题",
      "priority": 4,
      "status": 1
    },
    {
      "code": "todo-task-two",
      "description": "更新后的描述"
    }
  ]

  # 状态值：0=待处理, 1=进行中, 2=已完成, 3=已取消
  # 优先级值：1=低, 2=中, 3=高, 4=紧急`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析 JSON 输入
		items, err := parseUpdateItems(todoBatchUpdateJSON, todoBatchUpdateJSONFile)
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
		if err := handler.BatchUpdate(bs.Context(), items); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	todoBatchUpdateCmd.Flags().StringVar(&todoBatchUpdateJSON, "json", "", "JSON格式的更新列表")
	todoBatchUpdateCmd.Flags().StringVar(&todoBatchUpdateJSONFile, "json-file", "", "包含更新列表的JSON文件路径")

	todoCmd.AddCommand(todoBatchUpdateCmd)
}

// parseUpdateItems 解析 JSON 输入为待办更新项列表
func parseUpdateItems(jsonStr string, jsonFile string) ([]dto.ToDoUpdateDTO, error) {
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

	var items []dto.ToDoUpdateDTO
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("JSON 格式错误: %w", err)
	}

	return items, nil
}

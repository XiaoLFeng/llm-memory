package main

import (
	"github.com/XiaoLFeng/llm-memory/cmd"

	// 导入子命令包，触发 init() 注册命令
	_ "github.com/XiaoLFeng/llm-memory/cmd/group"
	_ "github.com/XiaoLFeng/llm-memory/cmd/memory"
	_ "github.com/XiaoLFeng/llm-memory/cmd/plan"
	_ "github.com/XiaoLFeng/llm-memory/cmd/todo"
)

// main 是程序的入口点
// 嘿嘿~ 这里是整个 LLM Memory 应用的起点呢！(´∀｀)💖
func main() {
	cmd.Execute()
}

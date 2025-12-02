package cli

import (
	"github.com/XiaoLFeng/llm-memory/startup"
)

// CLI 是 CLI 模块的核心结构
// 嘿嘿~ 封装了 Bootstrap 和输出配置！(´∀｀)💖
type CLI struct {
	bs     *startup.Bootstrap
	output OutputFormat
}

// OutputFormat 输出格式类型
type OutputFormat string

const (
	OutputTable OutputFormat = "table" // 表格格式
	OutputJSON  OutputFormat = "json"  // JSON格式
	OutputPlain OutputFormat = "plain" // 纯文本格式
)

// New 创建新的 CLI 实例
// 呀~ 初始化 CLI 核心结构！✨
func New(bs *startup.Bootstrap) *CLI {
	return &CLI{
		bs:     bs,
		output: OutputTable, // 默认表格输出
	}
}

// SetOutputFormat 设置输出格式
func (c *CLI) SetOutputFormat(format OutputFormat) {
	c.output = format
}

// OutputFormat 获取当前输出格式
func (c *CLI) GetOutputFormat() OutputFormat {
	return c.output
}

// Bootstrap 获取 Bootstrap 实例
func (c *CLI) Bootstrap() *startup.Bootstrap {
	return c.bs
}

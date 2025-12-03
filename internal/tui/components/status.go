package components

import (
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
)

// KeyHint 快捷键提示
type KeyHint struct {
	Key  string
	Desc string
}

func RenderKeys(keys []KeyHint) []string {
	var out []string
	for _, k := range keys {
		out = append(out, theme.KeyStyle.Render(k.Key)+" "+theme.ValueStyle.Render(k.Desc))
	}
	return out
}

// JoinKeys 将 KeyHint 渲染为一行
func JoinKeys(keys []KeyHint) string {
	parts := RenderKeys(keys)
	return strings.Join(parts, "  ")
}

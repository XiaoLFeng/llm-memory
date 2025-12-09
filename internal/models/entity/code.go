package entity

import (
	"errors"
	"regexp"
)

// Code 格式正则表达式
// 规则: 开头必须是小写字母，中间可含小写字母、数字、连字符，结尾必须是字母或数字，最少3个字符
// 示例有效值: "abc", "my-task", "task-123", "a1b", "test-1a"
// 示例无效值: "ab", "1abc", "abc-", "ABC", "my_task"
var codeRegex = regexp.MustCompile(`^[a-z][a-z0-9\-]*[a-z0-9]$`)

// ValidateCode 验证 code 格式
// 返回 nil 表示验证通过，返回 error 表示格式不正确
func ValidateCode(code string) error {
	if len(code) < 3 {
		return errors.New("code 长度至少为 3 个字符")
	}
	if !codeRegex.MatchString(code) {
		return errors.New("code 格式错误: 开头必须是小写字母，中间可含字母/数字/连字符，结尾必须是字母或数字")
	}
	return nil
}

// IsValidCode 检查 code 是否有效（便捷方法）
func IsValidCode(code string) bool {
	return ValidateCode(code) == nil
}

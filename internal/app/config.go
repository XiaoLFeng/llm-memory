package app

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Config 应用配置结构体 ✨
// 存储应用的各项配置信息，包括数据库路径、主题和调试模式
type Config struct {
	DBPath string `json:"db_path"` // 数据库文件路径
	Theme  string `json:"theme"`   // 主题名称
	Debug  bool   `json:"debug"`   // 调试模式开关
}

// DefaultConfig 返回默认配置 🎮
// 默认配置包括：
// - DBPath: ~/.llm-memory/data.db
// - Theme: default
// - Debug: false
func DefaultConfig() *Config {
	configDir := GetConfigDir()
	return &Config{
		DBPath: filepath.Join(configDir, "data.db"),
		Theme:  "default",
		Debug:  false,
	}
}

// GetConfigDir 获取配置目录路径 📁
// 返回 ~/.llm-memory 目录的绝对路径
func GetConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// 如果无法获取用户主目录，使用当前目录
		return ".llm-memory"
	}
	return filepath.Join(homeDir, ".llm-memory")
}

// LoadConfig 从配置文件加载配置 💾
// 如果配置文件不存在，则创建默认配置文件并返回默认配置
// 如果配置文件存在但读取失败，返回错误
func LoadConfig() (*Config, error) {
	configDir := GetConfigDir()
	configPath := filepath.Join(configDir, "config.json")

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		// 配置文件不存在，创建默认配置
		defaultCfg := DefaultConfig()
		if err := SaveConfig(defaultCfg); err != nil {
			return nil, err
		}
		return defaultCfg, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// 解析 JSON 配置
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig 保存配置到文件 💖
// 将配置序列化为 JSON 并写入 ~/.llm-memory/config.json
// 如果配置目录不存在，会自动创建
func SaveConfig(config *Config) error {
	configDir := GetConfigDir()

	// 确保配置目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")

	// 序列化配置为 JSON（格式化输出，方便人类阅读）
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// 写入配置文件
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return err
	}

	return nil
}

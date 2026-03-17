package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 全局配置结构
type Config struct {
	// 文件路径配置
	TargetPath string `yaml:"targetPath"`
	BackupPath string `yaml:"backupPath"`

	// 运行模式配置
	Mode        string `yaml:"mode"`        // interactive 或 batch
	DryRun      bool   `yaml:"dryRun"`      // 仅模拟操作
	Interactive bool   `yaml:"interactive"` // 操作前确认
}

// LoadConfig 从 YAML/JSON 配置文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		// 如果文件不存在，返回默认配置
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 应用默认值
	cfg.ApplyDefaults()

	return cfg, nil
}

// ApplyDefaults 为空字段应用默认值
func (c *Config) ApplyDefaults() {
	if c.Mode == "" {
		c.Mode = "interactive"
	}
}

// Validate 验证配置是否有效
func (c *Config) Validate() error {
	if c.TargetPath == "" {
		return fmt.Errorf("目标路径 (targetPath) 不能为空")
	}
	if c.BackupPath == "" {
		return fmt.Errorf("备份路径 (backupPath) 不能为空")
	}

	// 检查目录是否存在
	if info, err := os.Stat(c.TargetPath); err != nil || !info.IsDir() {
		return fmt.Errorf("目标路径不存在或不是目录: %s", c.TargetPath)
	}

	// 检查备份目录是否存在，不存在则创建
	if _, err := os.Stat(c.BackupPath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(c.BackupPath, 0755); err != nil {
				return fmt.Errorf("创建备份目录失败: %w", err)
			}
		} else {
			return fmt.Errorf("检查备份路径失败: %w", err)
		}
	}

	return nil
}

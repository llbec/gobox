package config

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Mode:        "interactive",
		DryRun:      false,
		Interactive: false,
	}
}

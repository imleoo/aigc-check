// Package gemini 提供 Gemini API 集成功能
// 用于深度语义分析和智能建议生成
package gemini

import (
	"os"
	"time"
)

// Config Gemini API 配置
type Config struct {
	// 是否启用 Gemini 分析
	Enabled bool `yaml:"enabled"`

	// API Key（可通过环境变量 GEMINI_API_KEY 设置）
	APIKey string `yaml:"api_key"`

	// 模型名称
	Model string `yaml:"model"`

	// 温度参数（0-1，越低越确定性）
	Temperature float64 `yaml:"temperature"`

	// 最大输出 token 数
	MaxTokens int `yaml:"max_tokens"`

	// 请求超时
	Timeout time.Duration `yaml:"timeout"`

	// 重试配置
	Retry RetryConfig `yaml:"retry"`

	// 缓存配置
	Cache CacheConfig `yaml:"cache"`

	// API 端点（可选，默认使用官方端点）
	Endpoint string `yaml:"endpoint"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	// 最大重试次数
	MaxAttempts int `yaml:"max_attempts"`

	// 初始退避时间
	InitialBackoff time.Duration `yaml:"initial_backoff"`

	// 最大退避时间
	MaxBackoff time.Duration `yaml:"max_backoff"`

	// 退避乘数
	Multiplier float64 `yaml:"multiplier"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	// 是否启用缓存
	Enabled bool `yaml:"enabled"`

	// 缓存 TTL
	TTL time.Duration `yaml:"ttl"`

	// 最大缓存条目数
	MaxEntries int `yaml:"max_entries"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Enabled:     false,
		APIKey:      "",
		Model:       "gemini-pro",
		Temperature: 0.3,
		MaxTokens:   500,
		Timeout:     30 * time.Second,
		Endpoint:    "https://generativelanguage.googleapis.com/v1beta",
		Retry: RetryConfig{
			MaxAttempts:    3,
			InitialBackoff: 1 * time.Second,
			MaxBackoff:     30 * time.Second,
			Multiplier:     2.0,
		},
		Cache: CacheConfig{
			Enabled:    true,
			TTL:        1 * time.Hour,
			MaxEntries: 1000,
		},
	}
}

// LoadFromEnv 从环境变量加载配置
func (c *Config) LoadFromEnv() {
	if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
		c.APIKey = apiKey
	}

	if enabled := os.Getenv("GEMINI_ENABLED"); enabled == "true" {
		c.Enabled = true
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if !c.Enabled {
		return nil // 未启用时不需要验证
	}

	if c.APIKey == "" {
		return ErrMissingAPIKey
	}

	if c.Model == "" {
		c.Model = "gemini-pro"
	}

	if c.Temperature < 0 || c.Temperature > 1 {
		c.Temperature = 0.3
	}

	if c.MaxTokens <= 0 {
		c.MaxTokens = 500
	}

	if c.Timeout <= 0 {
		c.Timeout = 30 * time.Second
	}

	return nil
}

// IsEnabled 检查是否启用
func (c *Config) IsEnabled() bool {
	return c.Enabled && c.APIKey != ""
}

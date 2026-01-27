package config

import (
	"os"

	"github.com/leoobai/aigc-check/internal/models"
	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Thresholds Thresholds              `yaml:"thresholds"` // 规则阈值配置
	Scoring    ScoringConfig           `yaml:"scoring"`    // 评分配置
	Output     OutputConfig            `yaml:"output"`     // 输出配置
	Rules      map[string]RuleConfig   `yaml:"rules"`      // 规则配置
}

// ScoringConfig 评分配置
type ScoringConfig struct {
	Weights models.DimensionWeights `yaml:"weights"` // 维度权重
}

// OutputConfig 输出配置
type OutputConfig struct {
	DefaultFormat string `yaml:"default_format"` // 默认输出格式: text, json
	Language      string `yaml:"language"`       // 语言: zh, en
	Verbose       bool   `yaml:"verbose"`        // 详细输出
	ColorEnabled  bool   `yaml:"color_enabled"`  // 启用颜色输出
}

// RuleConfig 规则配置
type RuleConfig struct {
	Enabled   bool                   `yaml:"enabled"`   // 是否启用
	Threshold int                    `yaml:"threshold"` // 阈值
	Severity  models.Severity        `yaml:"severity"`  // 严重程度
	Options   map[string]interface{} `yaml:"options"`   // 其他选项
}

// DefaultConfig 默认配置
var DefaultConfig = Config{
	Thresholds: DefaultThresholds,
	Scoring: ScoringConfig{
		Weights: models.DefaultDimensionWeights,
	},
	Output: OutputConfig{
		DefaultFormat: "text",
		Language:      "zh",
		Verbose:       false,
		ColorEnabled:  true,
	},
	Rules: map[string]RuleConfig{
		string(models.RuleTypeHighFreqWords): {
			Enabled:   true,
			Threshold: 3,
			Severity:  models.SeverityHigh,
		},
		string(models.RuleTypeSentenceStarters): {
			Enabled:   true,
			Threshold: 5,
			Severity:  models.SeverityHigh,
		},
		string(models.RuleTypeFalseRange): {
			Enabled:   true,
			Threshold: 2,
			Severity:  models.SeverityMedium,
		},
		string(models.RuleTypeCitationAnomaly): {
			Enabled:   true,
			Threshold: 1,
			Severity:  models.SeverityCritical,
		},
		string(models.RuleTypeEmDash): {
			Enabled:   true,
			Threshold: 5,
			Severity:  models.SeverityMedium,
		},
		string(models.RuleTypeMarkdown): {
			Enabled:   true,
			Threshold: 3,
			Severity:  models.SeverityHigh,
		},
		string(models.RuleTypeEmoji): {
			Enabled:   true,
			Threshold: 5,
			Severity:  models.SeverityMedium,
		},
		string(models.RuleTypeKnowledgeCutoff): {
			Enabled:   true,
			Threshold: 1,
			Severity:  models.SeverityCritical,
		},
		string(models.RuleTypeCollaborative): {
			Enabled:   true,
			Threshold: 3,
			Severity:  models.SeverityHigh,
		},
		string(models.RuleTypePerfectionism): {
			Enabled:   true,
			Threshold: 5,
			Severity:  models.SeverityHigh,
		},
	},
}

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*Config, error) {
	// 如果文件不存在，返回默认配置
	if _, err := os.Stat(path); os.IsNotExist(err) {
		config := DefaultConfig
		return &config, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// 合并默认配置
	mergeWithDefaults(&config)

	return &config, nil
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// mergeWithDefaults 合并默认配置
func mergeWithDefaults(config *Config) {
	// 如果规则配置为空，使用默认配置
	if config.Rules == nil {
		config.Rules = DefaultConfig.Rules
	}

	// 补充缺失的规则配置
	for ruleType, defaultRule := range DefaultConfig.Rules {
		if _, exists := config.Rules[ruleType]; !exists {
			config.Rules[ruleType] = defaultRule
		}
	}

	// 如果输出配置为空，使用默认值
	if config.Output.DefaultFormat == "" {
		config.Output.DefaultFormat = DefaultConfig.Output.DefaultFormat
	}
	if config.Output.Language == "" {
		config.Output.Language = DefaultConfig.Output.Language
	}
}

// GetRuleConfig 获取规则配置
func (c *Config) GetRuleConfig(ruleType models.RuleType) RuleConfig {
	if rule, exists := c.Rules[string(ruleType)]; exists {
		return rule
	}
	// 返回默认配置
	return RuleConfig{
		Enabled:   true,
		Threshold: 1,
		Severity:  models.SeverityMedium,
	}
}

// IsRuleEnabled 检查规则是否启用
func (c *Config) IsRuleEnabled(ruleType models.RuleType) bool {
	rule := c.GetRuleConfig(ruleType)
	return rule.Enabled
}

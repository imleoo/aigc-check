package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leoobai/aigc-check/internal/models"
)

func TestLoadConfig_NonExistentFile(t *testing.T) {
	cfg, err := LoadConfig("/non/existent/path/config.yaml")
	if err != nil {
		t.Fatalf("LoadConfig() error = %v, want nil", err)
	}

	if cfg == nil {
		t.Fatal("LoadConfig() returned nil config")
	}

	// 验证返回默认配置
	if cfg.Output.DefaultFormat != "text" {
		t.Errorf("DefaultFormat = %s, want text", cfg.Output.DefaultFormat)
	}
	if cfg.Output.Language != "zh" {
		t.Errorf("Language = %s, want zh", cfg.Output.Language)
	}
}

func TestLoadConfig_ValidFile(t *testing.T) {
	// 创建临时配置文件
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `
output:
  default_format: "json"
  language: "en"
  verbose: true
  color_enabled: false
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Output.DefaultFormat != "json" {
		t.Errorf("DefaultFormat = %s, want json", cfg.Output.DefaultFormat)
	}
	if cfg.Output.Language != "en" {
		t.Errorf("Language = %s, want en", cfg.Output.Language)
	}
	if !cfg.Output.Verbose {
		t.Error("Verbose = false, want true")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid.yaml")

	// 写入无效的YAML
	if err := os.WriteFile(configPath, []byte("invalid: yaml: content:"), 0644); err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}

	_, err := LoadConfig(configPath)
	if err == nil {
		t.Error("LoadConfig() expected error for invalid YAML, got nil")
	}
}

func TestSaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "save_test.yaml")

	cfg := &Config{
		Output: OutputConfig{
			DefaultFormat: "json",
			Language:      "en",
			Verbose:       true,
		},
	}

	if err := SaveConfig(cfg, configPath); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("SaveConfig() did not create file")
	}

	// 重新加载验证
	loaded, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() after save error = %v", err)
	}

	if loaded.Output.DefaultFormat != "json" {
		t.Errorf("Loaded DefaultFormat = %s, want json", loaded.Output.DefaultFormat)
	}
}

func TestConfig_GetRuleConfig(t *testing.T) {
	cfg := &DefaultConfig

	tests := []struct {
		ruleType    models.RuleType
		wantEnabled bool
	}{
		{models.RuleTypeHighFreqWords, true},
		{models.RuleTypeSentenceStarters, true},
		{models.RuleTypeCitationAnomaly, true},
		{models.RuleTypeKnowledgeCutoff, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.ruleType), func(t *testing.T) {
			rule := cfg.GetRuleConfig(tt.ruleType)
			if rule.Enabled != tt.wantEnabled {
				t.Errorf("GetRuleConfig(%s).Enabled = %v, want %v",
					tt.ruleType, rule.Enabled, tt.wantEnabled)
			}
		})
	}
}

func TestConfig_GetRuleConfig_NotExists(t *testing.T) {
	cfg := &Config{
		Rules: map[string]RuleConfig{},
	}

	rule := cfg.GetRuleConfig(models.RuleTypeHighFreqWords)

	// 应返回默认配置
	if !rule.Enabled {
		t.Error("GetRuleConfig() for non-existent rule should return enabled=true")
	}
	if rule.Threshold != 1 {
		t.Errorf("GetRuleConfig() Threshold = %d, want 1", rule.Threshold)
	}
}

func TestConfig_IsRuleEnabled(t *testing.T) {
	cfg := &Config{
		Rules: map[string]RuleConfig{
			string(models.RuleTypeHighFreqWords): {
				Enabled: true,
			},
			string(models.RuleTypeSentenceStarters): {
				Enabled: false,
			},
		},
	}

	if !cfg.IsRuleEnabled(models.RuleTypeHighFreqWords) {
		t.Error("IsRuleEnabled(HighFreqWords) = false, want true")
	}

	if cfg.IsRuleEnabled(models.RuleTypeSentenceStarters) {
		t.Error("IsRuleEnabled(SentenceStarters) = true, want false")
	}
}

func TestMergeWithDefaults(t *testing.T) {
	cfg := &Config{
		Rules: nil,
		Output: OutputConfig{
			DefaultFormat: "",
			Language:      "",
		},
	}

	mergeWithDefaults(cfg)

	// 验证规则被合并
	if cfg.Rules == nil {
		t.Error("mergeWithDefaults() did not set Rules")
	}

	// 验证输出配置被合并
	if cfg.Output.DefaultFormat != "text" {
		t.Errorf("mergeWithDefaults() DefaultFormat = %s, want text", cfg.Output.DefaultFormat)
	}
	if cfg.Output.Language != "zh" {
		t.Errorf("mergeWithDefaults() Language = %s, want zh", cfg.Output.Language)
	}
}

func TestDefaultConfig_AllRulesPresent(t *testing.T) {
	expectedRules := []models.RuleType{
		models.RuleTypeHighFreqWords,
		models.RuleTypeSentenceStarters,
		models.RuleTypeFalseRange,
		models.RuleTypeCitationAnomaly,
		models.RuleTypeEmDash,
		models.RuleTypeMarkdown,
		models.RuleTypeEmoji,
		models.RuleTypeKnowledgeCutoff,
		models.RuleTypeCollaborative,
		models.RuleTypePerfectionism,
	}

	for _, ruleType := range expectedRules {
		if _, exists := DefaultConfig.Rules[string(ruleType)]; !exists {
			t.Errorf("DefaultConfig missing rule: %s", ruleType)
		}
	}
}

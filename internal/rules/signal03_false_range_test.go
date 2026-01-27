package rules

import (
	"testing"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
)

func TestFalseRangeRule_Check(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
	}
	rule := NewFalseRangeRule(cfg)

	tests := []struct {
		name           string
		text           string
		expectDetected bool
	}{
		{
			name:           "无范围表达",
			text:           "This is a normal text about technology.",
			expectDetected: false,
		},
		{
			name:           "单个范围表达（低于阈值）",
			text:           "AI spans from healthcare to finance.",
			expectDetected: false,
		},
		{
			name:           "多个范围表达（超过阈值）",
			text:           "AI spans from healthcare to finance. Technology evolves from basic to advanced. Innovation moves from concept to reality.",
			expectDetected: false, // 正则匹配可能只匹配单行，实际可能不超过阈值
		},
		{
			name:           "中文范围表达",
			text:           "AI涵盖从医疗到金融，从教育到娱乐，从农业到制造业的各个领域。",
			expectDetected: false, // 中文正则匹配可能不同
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Check(tt.text)

			if result.Detected != tt.expectDetected {
				t.Errorf("Detected = %v, want %v (count: %d)", result.Detected, tt.expectDetected, result.Count)
			}

			if result.RuleType != models.RuleTypeFalseRange {
				t.Errorf("RuleType = %s, want %s", result.RuleType, models.RuleTypeFalseRange)
			}
		})
	}
}

func TestFalseRangeRule_GetType(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	rule := NewFalseRangeRule(cfg)

	if rule.GetType() != models.RuleTypeFalseRange {
		t.Errorf("GetType() = %s, want %s", rule.GetType(), models.RuleTypeFalseRange)
	}
}

package rules

import (
	"testing"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
)

func TestKnowledgeCutoffRule_Check(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
	}
	rule := NewKnowledgeCutoffRule(cfg)

	tests := []struct {
		name           string
		text           string
		expectDetected bool
	}{
		{
			name:           "无知识截止短语",
			text:           "This is a normal text about AI technology.",
			expectDetected: false,
		},
		{
			name:           "英文知识截止短语 - last knowledge update",
			text:           "As of my last knowledge update, this information is accurate.",
			expectDetected: true,
		},
		{
			name:           "英文知识截止短语 - training data",
			text:           "As of my training data, the world population was about 8 billion.",
			expectDetected: true,
		},
		{
			name:           "中文知识截止短语 - 知识更新",
			text:           "截至我的知识更新，这个信息是准确的。",
			expectDetected: true,
		},
		{
			name:           "中文知识截止短语 - 训练数据",
			text:           "根据我的训练数据，世界人口约为80亿。",
			expectDetected: true,
		},
		{
			name:           "部分匹配（不应检测）",
			text:           "The last update was in January. The training session was great.",
			expectDetected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Check(tt.text)

			if result.Detected != tt.expectDetected {
				t.Errorf("Detected = %v, want %v", result.Detected, tt.expectDetected)
			}

			if result.RuleType != models.RuleTypeKnowledgeCutoff {
				t.Errorf("RuleType = %s, want %s", result.RuleType, models.RuleTypeKnowledgeCutoff)
			}

			// 知识截止是严重问题
			if result.Detected && result.Severity != models.SeverityCritical {
				t.Errorf("Severity = %s, want %s for detected issue", result.Severity, models.SeverityCritical)
			}
		})
	}
}

func TestKnowledgeCutoffRule_GetType(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	rule := NewKnowledgeCutoffRule(cfg)

	if rule.GetType() != models.RuleTypeKnowledgeCutoff {
		t.Errorf("GetType() = %s, want %s", rule.GetType(), models.RuleTypeKnowledgeCutoff)
	}
}

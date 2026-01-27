package rules

import (
	"testing"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
)

func TestCollaborativeRule_Check(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
	}
	rule := NewCollaborativeRule(cfg)

	tests := []struct {
		name           string
		text           string
		expectDetected bool
	}{
		{
			name:           "无协作式短语",
			text:           "This is a normal article about technology. It discusses various topics.",
			expectDetected: false,
		},
		{
			name:           "少量协作式短语（低于阈值）",
			text:           "I hope this helps you understand the concept better.",
			expectDetected: false,
		},
		{
			name:           "大量协作式短语（超过阈值）",
			text:           "I hope this helps! Let me know if you have questions. Feel free to reach out. I hope this helps clarify things.",
			expectDetected: true,
		},
		{
			name:           "中文协作式短语",
			text:           "希望这能帮到你！如果需要更多信息，随时告诉我。希望这能帮到你理解这个概念。",
			expectDetected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Check(tt.text)

			if result.Detected != tt.expectDetected {
				t.Errorf("Detected = %v, want %v (count: %d)", result.Detected, tt.expectDetected, result.Count)
			}

			if result.RuleType != models.RuleTypeCollaborative {
				t.Errorf("RuleType = %s, want %s", result.RuleType, models.RuleTypeCollaborative)
			}
		})
	}
}

func TestCollaborativeRule_GetType(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	rule := NewCollaborativeRule(cfg)

	if rule.GetType() != models.RuleTypeCollaborative {
		t.Errorf("GetType() = %s, want %s", rule.GetType(), models.RuleTypeCollaborative)
	}
}

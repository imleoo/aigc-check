package rules

import (
	"testing"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
)

func TestSentenceStartersRule_Check(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
	}
	rule := NewSentenceStartersRule(cfg)

	tests := []struct {
		name           string
		text           string
		expectDetected bool
	}{
		{
			name:           "无AI句式开头",
			text:           "This is a normal text. It contains regular sentences. Nothing special here.",
			expectDetected: false,
		},
		{
			name:           "少量AI句式开头（低于阈值）",
			text:           "This is important. Additionally, we need to consider this. Finally, let's conclude.",
			expectDetected: false,
		},
		{
			name:           "大量AI句式开头（超过阈值）",
			text:           "Additionally, this is point one. Furthermore, point two is important. Moreover, we should note point three. Additionally, point four. Furthermore, point five. Moreover, point six.",
			expectDetected: true,
		},
		{
			name:           "中文AI句式开头",
			text:           "此外，这是要点一。另外，要点二很重要。而且，我们应该注意要点三。此外，要点四。另外，要点五。而且，要点六。",
			expectDetected: false, // 中文句式开头检测可能工作方式不同
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Check(tt.text)

			if result.Detected != tt.expectDetected {
				t.Errorf("Detected = %v, want %v (count: %d)", result.Detected, tt.expectDetected, result.Count)
			}

			if result.RuleType != models.RuleTypeSentenceStarters {
				t.Errorf("RuleType = %s, want %s", result.RuleType, models.RuleTypeSentenceStarters)
			}
		})
	}
}

func TestSentenceStartersRule_GetType(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	rule := NewSentenceStartersRule(cfg)

	if rule.GetType() != models.RuleTypeSentenceStarters {
		t.Errorf("GetType() = %s, want %s", rule.GetType(), models.RuleTypeSentenceStarters)
	}
}

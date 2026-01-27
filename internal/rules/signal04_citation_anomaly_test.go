package rules

import (
	"testing"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
)

func TestCitationAnomalyRule_Check(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
	}
	rule := NewCitationAnomalyRule(cfg)

	tests := []struct {
		name           string
		text           string
		expectDetected bool
	}{
		{
			name:           "无引用异常",
			text:           "This is a normal text with a link: https://example.com",
			expectDetected: false,
		},
		{
			name:           "UTM参数 - ChatGPT",
			text:           "Check out this link: https://example.com?utm_source=chatgpt.com",
			expectDetected: true,
		},
		{
			name:           "UTM参数 - OpenAI",
			text:           "Visit: https://example.com?utm_source=openai",
			expectDetected: true,
		},
		{
			name:           "幽灵标记 - oaicite",
			text:           "According to the study contentReference[oaicite:0] this is true.",
			expectDetected: true,
		},
		{
			name:           "幽灵标记 - oai_citation",
			text:           "The research shows [oai_citation: 1] important findings.",
			expectDetected: true,
		},
		{
			name:           "幽灵标记 - turn0search",
			text:           "Based on turn0search results, we can see...",
			expectDetected: true,
		},
		{
			name:           "占位符日期 - 2025-XX-XX",
			text:           "The event is scheduled for 2025-XX-XX.",
			expectDetected: true,
		},
		{
			name:           "占位符日期 - YYYY-MM-DD",
			text:           "Please enter the date in YYYY-MM-DD format.",
			expectDetected: true,
		},
		{
			name:           "正常日期（非占位符）",
			text:           "The meeting is on 2025-01-15.",
			expectDetected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Check(tt.text)

			if result.Detected != tt.expectDetected {
				t.Errorf("Detected = %v, want %v", result.Detected, tt.expectDetected)
			}

			if result.RuleType != models.RuleTypeCitationAnomaly {
				t.Errorf("RuleType = %s, want %s", result.RuleType, models.RuleTypeCitationAnomaly)
			}

			// 引用异常是严重问题
			if result.Detected && result.Severity != models.SeverityCritical {
				t.Errorf("Severity = %s, want %s for detected anomaly", result.Severity, models.SeverityCritical)
			}
		})
	}
}

func TestCitationAnomalyRule_GetType(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	rule := NewCitationAnomalyRule(cfg)

	if rule.GetType() != models.RuleTypeCitationAnomaly {
		t.Errorf("GetType() = %s, want %s", rule.GetType(), models.RuleTypeCitationAnomaly)
	}
}

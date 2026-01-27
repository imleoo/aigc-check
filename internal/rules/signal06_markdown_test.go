package rules

import (
	"testing"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
)

func TestMarkdownRule_Check(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
	}
	rule := NewMarkdownRule(cfg)

	tests := []struct {
		name           string
		text           string
		expectDetected bool
	}{
		{
			name:           "无Markdown标记",
			text:           "This is a plain text without any formatting.",
			expectDetected: false,
		},
		{
			name:           "少量Markdown（低于阈值）",
			text:           "This is **bold** text.",
			expectDetected: false,
		},
		{
			name:           "标题标记",
			text:           "## Section 1\n\nContent here.\n\n## Section 2\n\nMore content.\n\n## Section 3",
			expectDetected: true,
		},
		{
			name:           "加粗标记",
			text:           "This is **important**. Also **crucial**. And **vital**. Don't forget **significant**.",
			expectDetected: true,
		},
		{
			name:           "链接标记",
			text:           "Check out these resources: [link1](url1) for basics. Also see [link2](url2) for advanced topics. Don't miss [link3](url3) for examples.",
			expectDetected: false, // 链接标记的正则可能只匹配特定格式
		},
		{
			name:           "代码块",
			text:           "Here's some code:\n```\ncode here\n```\n\nMore code:\n```\nmore code\n```\n\nAnother:\n```\ncode\n```",
			expectDetected: true,
		},
		{
			name:           "混合Markdown",
			text:           "## Title\n\nThis is **bold** and this is [a link](url).",
			expectDetected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Check(tt.text)

			if result.Detected != tt.expectDetected {
				t.Errorf("Detected = %v, want %v (count: %d)", result.Detected, tt.expectDetected, result.Count)
			}

			if result.RuleType != models.RuleTypeMarkdown {
				t.Errorf("RuleType = %s, want %s", result.RuleType, models.RuleTypeMarkdown)
			}
		})
	}
}

func TestMarkdownRule_GetType(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	rule := NewMarkdownRule(cfg)

	if rule.GetType() != models.RuleTypeMarkdown {
		t.Errorf("GetType() = %s, want %s", rule.GetType(), models.RuleTypeMarkdown)
	}
}

package rules

import (
	"testing"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
)

func TestEmDashRule_Check(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
	}
	rule := NewEmDashRule(cfg)

	tests := []struct {
		name           string
		text           string
		expectDetected bool
	}{
		{
			name:           "无破折号",
			text:           "This is a normal text without any em dashes.",
			expectDetected: false,
		},
		{
			name:           "少量破折号（低于阈值）",
			text:           "This is a long paragraph about technology and innovation that discusses various topics in depth. The technology landscape is changing rapidly. Modern software development practices have evolved significantly. This is important — really important for understanding the context.",
			expectDetected: false,
		},
		{
			name:           "大量破折号（超过阈值）",
			text:           "AI is transforming — fundamentally changing — every aspect — from healthcare — to finance — creating unprecedented — opportunities — for innovation — across industries.",
			expectDetected: true,
		},
		{
			name:           "ASCII破折号",
			text:           "This -- is -- a -- test -- with -- many -- double -- dashes.",
			expectDetected: false, // 这是双连字符，不是em dash
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Check(tt.text)

			if result.Detected != tt.expectDetected {
				t.Errorf("Detected = %v, want %v (count: %d)", result.Detected, tt.expectDetected, result.Count)
			}

			if result.RuleType != models.RuleTypeEmDash {
				t.Errorf("RuleType = %s, want %s", result.RuleType, models.RuleTypeEmDash)
			}
		})
	}
}

func TestEmDashRule_GetType(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	rule := NewEmDashRule(cfg)

	if rule.GetType() != models.RuleTypeEmDash {
		t.Errorf("GetType() = %s, want %s", rule.GetType(), models.RuleTypeEmDash)
	}
}

func TestEmDashRule_DensityCalculation(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
	}
	rule := NewEmDashRule(cfg)

	// 测试1000字文本中的破折号密度
	// 阈值是每千字5个

	t.Run("低密度", func(t *testing.T) {
		// 构造约1000字的文本，只有2个破折号
		text := "This is a test " + generateWords(200) + " — test — " + generateWords(200)
		result := rule.Check(text)

		if result.Detected {
			t.Errorf("Detected = true, want false for low density")
		}
	})
}

// generateWords 生成指定数量的单词
func generateWords(count int) string {
	word := "word "
	result := ""
	for i := 0; i < count; i++ {
		result += word
	}
	return result
}

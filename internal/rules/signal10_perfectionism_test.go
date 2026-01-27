package rules

import (
	"testing"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
)

func TestPerfectionismRule_Check(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
	}
	rule := NewPerfectionismRule(cfg)

	tests := []struct {
		name           string
		text           string
		expectDetected bool
		description    string
	}{
		{
			name:           "充分的个人化表达",
			text:           "I think this is interesting. I feel that maybe we could try a different approach. I'm not sure if this will work, but I believe it's worth trying.",
			expectDetected: false,
			description:    "包含第一人称、情感词和不确定性标记",
		},
		{
			name:           "缺乏个人化表达（极短）",
			text:           "Data processed.",
			expectDetected: true,
			description:    "极短文本缺乏个人化指标",
		},
		{
			name:           "中文个人化表达",
			text:           "我觉得这个方案可能有效。我希望这能帮助大家理解。也许我们可以尝试不同的方法。",
			expectDetected: false,
			description:    "包含中文第一人称、情感词和不确定性标记",
		},
		{
			name:           "中文缺乏个人化",
			text:           "该系统高效处理数据。",
			expectDetected: true,
			description:    "极短中文没有个人化表达",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Check(tt.text)

			if result.Detected != tt.expectDetected {
				t.Errorf("Detected = %v, want %v (%s) [count: %d]",
					result.Detected, tt.expectDetected, tt.description, result.Count)
			}

			if result.RuleType != models.RuleTypePerfectionism {
				t.Errorf("RuleType = %s, want %s", result.RuleType, models.RuleTypePerfectionism)
			}
		})
	}
}

func TestPerfectionismRule_Score(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	rule := NewPerfectionismRule(cfg)

	tests := []struct {
		name     string
		text     string
		minScore float64
		maxScore float64
	}{
		{
			name:     "丰富个人化应得高分",
			text:     "I think maybe this could work. I feel it might be a good solution. I'm not sure, but I believe we should try. Perhaps it will succeed.",
			minScore: 80,
			maxScore: 100,
		},
		{
			name:     "极短文本应得低分",
			text:     "Done.",
			minScore: 0,
			maxScore: 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Check(tt.text)

			if result.Score < tt.minScore || result.Score > tt.maxScore {
				t.Errorf("Score = %.1f, want between %.1f and %.1f", result.Score, tt.minScore, tt.maxScore)
			}
		})
	}
}

func TestPerfectionismRule_GetType(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	rule := NewPerfectionismRule(cfg)

	if rule.GetType() != models.RuleTypePerfectionism {
		t.Errorf("GetType() = %s, want %s", rule.GetType(), models.RuleTypePerfectionism)
	}
}

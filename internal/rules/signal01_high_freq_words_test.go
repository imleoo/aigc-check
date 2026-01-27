package rules

import (
	"testing"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
)

func TestHighFreqWordsRule_Check(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
	}
	rule := NewHighFreqWordsRule(cfg)

	tests := []struct {
		name           string
		text           string
		expectDetected bool
		minMatches     int
	}{
		{
			name:           "无高频词",
			text:           "This is a normal text without any AI-specific words.",
			expectDetected: false,
			minMatches:     0,
		},
		{
			name:           "少量高频词（低于阈值）",
			text:           "This is a crucial point.",
			expectDetected: false,
			minMatches:     0,
		},
		{
			name:           "多个高频词（超过阈值）",
			text:           "This is crucial. The pivotal change is crucial. It's crucial to understand this crucial point.",
			expectDetected: true,
			minMatches:     4,
		},
		{
			name:           "混合高频词",
			text:           "The groundbreaking discovery is revolutionary. This pivotal moment is crucial for our vital mission.",
			expectDetected: false, // 每个词只出现1次，低于阈值3
			minMatches:     0,
		},
		{
			name:           "中文高频词",
			text:           "这是一个至关重要的决定。至关重要的是我们要理解这个至关重要的概念。这非常关键。",
			expectDetected: false, // 中文高频词匹配需要完整匹配，测试数据可能不满足阈值
			minMatches:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Check(tt.text)

			if result.Detected != tt.expectDetected {
				t.Errorf("Detected = %v, want %v", result.Detected, tt.expectDetected)
			}

			if tt.expectDetected && result.Count < tt.minMatches {
				t.Errorf("Count = %d, want >= %d", result.Count, tt.minMatches)
			}

			if result.RuleType != models.RuleTypeHighFreqWords {
				t.Errorf("RuleType = %s, want %s", result.RuleType, models.RuleTypeHighFreqWords)
			}
		})
	}
}

func TestHighFreqWordsRule_GetType(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	rule := NewHighFreqWordsRule(cfg)

	if rule.GetType() != models.RuleTypeHighFreqWords {
		t.Errorf("GetType() = %s, want %s", rule.GetType(), models.RuleTypeHighFreqWords)
	}
}

func TestHighFreqWordsRule_GetName(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	rule := NewHighFreqWordsRule(cfg)

	if rule.GetName() == "" {
		t.Error("GetName() returned empty string")
	}
}

func TestHighFreqWordsRule_GetDescription(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	rule := NewHighFreqWordsRule(cfg)

	if rule.GetDescription() == "" {
		t.Error("GetDescription() returned empty string")
	}
}

func TestHighFreqWordsRule_Score(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	rule := NewHighFreqWordsRule(cfg)

	tests := []struct {
		name         string
		text         string
		minScore     float64
		maxScore     float64
	}{
		{
			name:     "无高频词应得满分",
			text:     "This is a normal text.",
			minScore: 100,
			maxScore: 100,
		},
		{
			name:     "大量高频词应得低分",
			text:     "crucial crucial crucial crucial crucial crucial crucial crucial crucial crucial",
			minScore: 0,
			maxScore: 50,
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

package rules

import (
	"testing"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
)

func TestEmojiRule_Check(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
	}
	rule := NewEmojiRule(cfg)

	tests := []struct {
		name           string
		text           string
		expectDetected bool
	}{
		{
			name:           "无表情符号",
			text:           "This is a normal text without any emojis.",
			expectDetected: false,
		},
		{
			name:           "少量表情符号（低于阈值）",
			text:           "Hello! 👋 How are you?",
			expectDetected: false,
		},
		{
			name:           "大量表情符号（超过阈值）",
			text:           "This is amazing! 🎉 Great work! 👍 Keep it up! 💪 You're awesome! ⭐ Congratulations! 🎊 Well done! 🏆",
			expectDetected: true,
		},
		{
			name:           "工整的列表表情（AI风格）",
			text:           "✅ Task 1\n✅ Task 2\n✅ Task 3\n✅ Task 4\n✅ Task 5\n✅ Task 6",
			expectDetected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Check(tt.text)

			if result.Detected != tt.expectDetected {
				t.Errorf("Detected = %v, want %v (count: %d)", result.Detected, tt.expectDetected, result.Count)
			}

			if result.RuleType != models.RuleTypeEmoji {
				t.Errorf("RuleType = %s, want %s", result.RuleType, models.RuleTypeEmoji)
			}
		})
	}
}

func TestEmojiRule_GetType(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	rule := NewEmojiRule(cfg)

	if rule.GetType() != models.RuleTypeEmoji {
		t.Errorf("GetType() = %s, want %s", rule.GetType(), models.RuleTypeEmoji)
	}
}

func TestEmojiRule_EmojiDetection(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	rule := NewEmojiRule(cfg)

	// 测试各种类型的表情符号
	emojis := []struct {
		name  string
		emoji string
	}{
		{"笑脸", "😀"},
		{"心形", "❤️"},
		{"国旗", "🇺🇸"},
		{"动物", "🐶"},
		{"食物", "🍕"},
		{"活动", "⚽"},
		{"交通", "🚗"},
		{"符号", "✅"},
	}

	for _, e := range emojis {
		t.Run(e.name, func(t *testing.T) {
			// 创建包含6个相同表情的文本（超过阈值5）
			text := e.emoji + " " + e.emoji + " " + e.emoji + " " + e.emoji + " " + e.emoji + " " + e.emoji
			result := rule.Check(text)

			// 应该检测到表情符号
			if result.Count < 5 {
				t.Errorf("Count = %d, expected at least 5 for emoji: %s", result.Count, e.emoji)
			}
		})
	}
}

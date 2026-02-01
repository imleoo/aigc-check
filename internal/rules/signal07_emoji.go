package rules

import (
	"fmt"
	"strings"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
	"github.com/leoobai/aigc-check/internal/text"
)

// EmojiRule Signal 7: 表情符号异常检测
type EmojiRule struct {
	config    *config.Config
	processor *text.TextProcessor
}

// NewEmojiRule 创建表情符号异常检测规则
func NewEmojiRule(cfg *config.Config) *EmojiRule {
	return &EmojiRule{
		config:    cfg,
		processor: text.NewTextProcessor(),
	}
}

// Check 执行规则检测
func (rule *EmojiRule) Check(text string) models.RuleResult {
	result := models.RuleResult{
		RuleType:    models.RuleTypeEmoji,
		RuleName:    rule.GetName(),
		Description: rule.GetDescription(),
		Detected:    false,
		Score:       100.0,
		Severity:    models.SeverityMedium,
		Matches:     []models.Match{},
		Count:       0,
		Threshold:   rule.config.Thresholds.EmojiAnomaly.Threshold,
	}

	// 查找所有表情符号
	runes := []rune(text)
	offset := 0
	line := 1
	column := 1

	for i, char := range runes {
		if isEmoji(char) {
			result.Count++

			// 获取表情符号字符串
			emojiStr := string(char)

			result.Matches = append(result.Matches, models.Match{
				Text: emojiStr,
				Position: models.Position{
					Line:   line,
					Column: column,
					Offset: offset,
					Length: len(emojiStr),
				},
				Context: rule.getContext(text, offset, len(emojiStr)),
				Reason:  "检测到表情符号",
			})
		}

		// 更新位置
		if char == '\n' {
			line++
			column = 1
		} else {
			column++
		}
		offset = i
	}

	// 检查是否超过阈值
	if result.Count >= result.Threshold {
		result.Detected = true

		// 计算评分
		deduction := float64(result.Count-result.Threshold) * 8.0
		result.Score = 100.0 - deduction
		if result.Score < 0 {
			result.Score = 0
		}

		result.Message = fmt.Sprintf("检测到 %d 个表情符号，超过阈值 %d", result.Count, result.Threshold)
	} else {
		result.Message = fmt.Sprintf("检测到 %d 个表情符号，正常", result.Count)
	}

	return result
}

// GetType 获取规则类型
func (rule *EmojiRule) GetType() models.RuleType {
	return models.RuleTypeEmoji
}

// GetName 获取规则名称
func (rule *EmojiRule) GetName() string {
	return "表情符号异常检测"
}

// GetDescription 获取规则描述
func (rule *EmojiRule) GetDescription() string {
	return "检测过度工整的表情符号使用，AI生成内容倾向于规律性使用emoji"
}

// isEmoji 判断是否为表情符号
func isEmoji(r rune) bool {
	// 表情符号的Unicode范围
	return (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map
		(r >= 0x1F1E0 && r <= 0x1F1FF) || // Regional Indicators
		(r >= 0x2600 && r <= 0x26FF) || // Misc symbols
		(r >= 0x2700 && r <= 0x27BF) || // Dingbats
		(r >= 0xFE00 && r <= 0xFE0F) || // Variation Selectors
		(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols and Pictographs
		(r >= 0x1FA70 && r <= 0x1FAFF) // Symbols and Pictographs Extended-A
}

// getContext 获取匹配项的上下文
func (rule *EmojiRule) getContext(text string, offset, length int) string {
	const contextSize = 50

	// 将字节偏移量转换为 rune 索引

	runes := []rune(text)
	runeOffset := len([]rune(text[:offset]))
	runeLength := len([]rune(text[offset : offset+length]))
	start := runeOffset - contextSize
	if start < 0 {
		start = 0
	}

	end := runeOffset + runeLength + contextSize
	if end > len(runes) {
		end = len(runes)
	}

	context := string(runes[start:end])
	return strings.TrimSpace(context)
}

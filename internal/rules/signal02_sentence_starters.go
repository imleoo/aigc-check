package rules

import (
	"fmt"
	"strings"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
	"github.com/leoobai/aigc-check/internal/text"
)

// SentenceStartersRule Signal 2: 句式开头检测
type SentenceStartersRule struct {
	config    *config.Config
	processor *text.TextProcessor
}

// NewSentenceStartersRule 创建句式开头检测规则
func NewSentenceStartersRule(cfg *config.Config) *SentenceStartersRule {
	return &SentenceStartersRule{
		config:    cfg,
		processor: text.NewTextProcessor(),
	}
}

// Check 执行规则检测
func (r *SentenceStartersRule) Check(text string) models.RuleResult {
	result := models.RuleResult{
		RuleType:    models.RuleTypeSentenceStarters,
		RuleName:    r.GetName(),
		Description: r.GetDescription(),
		Detected:    false,
		Score:       100.0,
		Severity:    models.SeverityHigh,
		Matches:     []models.Match{},
		Count:       0,
		Threshold:   r.config.Thresholds.SentenceStarters.Threshold,
	}

	// 处理文本
	processed := r.processor.Process(text)

	patterns := r.config.Thresholds.SentenceStarters.Patterns

	// 检查每个句子的开头
	for _, sentence := range processed.Sentences {
		sentenceText := strings.TrimSpace(sentence.Text)
		if sentenceText == "" {
			continue
		}

		// 检查是否以指定模式开头
		for _, pattern := range patterns {
			if r.startsWithPattern(sentenceText, pattern) {
				result.Count++

				// 添加匹配项
				result.Matches = append(result.Matches, models.Match{
					Text:     pattern,
					Position: sentence.Position,
					Context:  r.truncate(sentenceText, 100),
					Reason:   fmt.Sprintf("句子以 '%s' 开头", pattern),
				})
				break // 每个句子只匹配一次
			}
		}
	}

	// 检查是否超过阈值
	if result.Count >= result.Threshold {
		result.Detected = true

		// 计算评分 (出现越多，分数越低)
		deduction := float64(result.Count-result.Threshold) * 8.0
		result.Score = 100.0 - deduction
		if result.Score < 0 {
			result.Score = 0
		}

		result.Message = fmt.Sprintf("检测到 %d 个句子使用重复的开头模式，超过阈值 %d", result.Count, result.Threshold)
	} else {
		result.Message = "未检测到异常的句式开头模式"
	}

	return result
}

// GetType 获取规则类型
func (r *SentenceStartersRule) GetType() models.RuleType {
	return models.RuleTypeSentenceStarters
}

// GetName 获取规则名称
func (r *SentenceStartersRule) GetName() string {
	return "句式开头检测"
}

// GetDescription 获取规则描述
func (r *SentenceStartersRule) GetDescription() string {
	return "检测重复使用的句式开头，如 Additionally, Furthermore, Moreover 等"
}

// startsWithPattern 检查句子是否以指定模式开头
func (r *SentenceStartersRule) startsWithPattern(sentence, pattern string) bool {
	sentence = strings.TrimSpace(sentence)
	pattern = strings.TrimSpace(pattern)

	// 不区分大小写比较
	sentenceLower := strings.ToLower(sentence)
	patternLower := strings.ToLower(pattern)

	if !strings.HasPrefix(sentenceLower, patternLower) {
		return false
	}

	// 确保模式后面是空格、逗号或句号
	if len(sentence) > len(pattern) {
		nextChar := sentence[len(pattern)]
		return nextChar == ' ' || nextChar == ',' || nextChar == '.'
	}

	return true
}

// truncate 截断文本
func (r *SentenceStartersRule) truncate(text string, maxLen int) string {
	runes := []rune(text)
	if len(runes) <= maxLen {
		return text
	}
	return string(runes[:maxLen]) + "..."
}

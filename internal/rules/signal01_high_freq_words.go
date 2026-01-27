package rules

import (
	"fmt"
	"strings"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
	"github.com/leoobai/aigc-check/internal/text"
)

// HighFreqWordsRule Signal 1: 高频词汇检测
type HighFreqWordsRule struct {
	config    *config.Config
	processor *text.TextProcessor
}

// NewHighFreqWordsRule 创建高频词汇检测规则
func NewHighFreqWordsRule(cfg *config.Config) *HighFreqWordsRule {
	return &HighFreqWordsRule{
		config:    cfg,
		processor: text.NewTextProcessor(),
	}
}

// Check 执行规则检测
func (r *HighFreqWordsRule) Check(text string) models.RuleResult {
	result := models.RuleResult{
		RuleType:    models.RuleTypeHighFreqWords,
		RuleName:    r.GetName(),
		Description: r.GetDescription(),
		Detected:    false,
		Score:       100.0,
		Severity:    models.SeverityHigh,
		Matches:     []models.Match{},
		Count:       0,
		Threshold:   r.config.Thresholds.HighFrequencyWords.Threshold,
	}

	// 处理文本
	processed := r.processor.Process(text)

	// 统计关键词出现次数
	keywordCounts := make(map[string]int)
	keywordPositions := make(map[string][]models.Position)

	keywords := r.config.Thresholds.HighFrequencyWords.Keywords

	// 遍历所有词汇
	for _, token := range processed.Words {
		for _, keyword := range keywords {
			if strings.EqualFold(token.Lower, strings.ToLower(keyword)) {
				keywordCounts[keyword]++
				keywordPositions[keyword] = append(keywordPositions[keyword], models.Position{
					Line:   token.Position.Line,
					Column: token.Position.Column,
					Offset: token.Position.Offset,
					Length: token.Position.Length,
				})
			}
		}
	}

	// 检查是否超过阈值
	totalCount := 0
	for keyword, count := range keywordCounts {
		if count >= result.Threshold {
			result.Detected = true
			totalCount += count

			// 添加匹配项
			for _, pos := range keywordPositions[keyword] {
				// 获取上下文
				context := r.getContext(text, pos.Offset, pos.Length)

				result.Matches = append(result.Matches, models.Match{
					Text:     keyword,
					Position: pos,
					Context:  context,
					Reason:   fmt.Sprintf("关键词 '%s' 出现 %d 次，超过阈值 %d", keyword, count, result.Threshold),
				})
			}
		}
	}

	result.Count = totalCount

	// 计算评分 (出现越多，分数越低)
	if result.Detected {
		// 每超过阈值1次，扣10分，最低0分
		deduction := float64(totalCount-result.Threshold) * 10.0
		result.Score = 100.0 - deduction
		if result.Score < 0 {
			result.Score = 0
		}

		result.Message = fmt.Sprintf("检测到 %d 个高频AI词汇，总计出现 %d 次", len(keywordCounts), totalCount)
	} else {
		result.Message = "未检测到异常的高频词汇使用"
	}

	return result
}

// GetType 获取规则类型
func (r *HighFreqWordsRule) GetType() models.RuleType {
	return models.RuleTypeHighFreqWords
}

// GetName 获取规则名称
func (r *HighFreqWordsRule) GetName() string {
	return "高频词汇检测"
}

// GetDescription 获取规则描述
func (r *HighFreqWordsRule) GetDescription() string {
	return "检测AI常用的高频词汇，如 crucial, pivotal, vital, groundbreaking 等"
}

// getContext 获取匹配项的上下文
func (r *HighFreqWordsRule) getContext(text string, offset, length int) string {
	const contextSize = 50 // 前后各50个字符

	runes := []rune(text)
	start := offset - contextSize
	if start < 0 {
		start = 0
	}

	end := offset + length + contextSize
	if end > len(runes) {
		end = len(runes)
	}

	context := string(runes[start:end])
	return strings.TrimSpace(context)
}

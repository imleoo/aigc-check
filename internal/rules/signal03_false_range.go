package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
	"github.com/leoobai/aigc-check/internal/text"
)

// FalseRangeRule Signal 3: 虚假范围表达检测
type FalseRangeRule struct {
	config    *config.Config
	processor *text.TextProcessor
}

// NewFalseRangeRule 创建虚假范围表达检测规则
func NewFalseRangeRule(cfg *config.Config) *FalseRangeRule {
	return &FalseRangeRule{
		config:    cfg,
		processor: text.NewTextProcessor(),
	}
}

// Check 执行规则检测
func (r *FalseRangeRule) Check(text string) models.RuleResult {
	result := models.RuleResult{
		RuleType:    models.RuleTypeFalseRange,
		RuleName:    r.GetName(),
		Description: r.GetDescription(),
		Detected:    false,
		Score:       100.0,
		Severity:    models.SeverityMedium,
		Matches:     []models.Match{},
		Count:       0,
		Threshold:   r.config.Thresholds.FalseRange.Threshold,
	}

	patterns := r.config.Thresholds.FalseRange.Patterns

	// 编译正则表达式
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}

		// 查找所有匹配
		matches := re.FindAllStringIndex(text, -1)
		for _, match := range matches {
			start := match[0]
			end := match[1]
			matchText := text[start:end]

			// 获取位置信息
			line, column := r.processor.GetLineColumn(text, start)

			result.Count++
			result.Matches = append(result.Matches, models.Match{
				Text: matchText,
				Position: models.Position{
					Line:   line,
					Column: column,
					Offset: start,
					Length: end - start,
				},
				Context: r.getContext(text, start, end-start),
				Reason:  "检测到可能的虚假范围表达",
			})
		}
	}

	// 检查是否超过阈值
	if result.Count >= result.Threshold {
		result.Detected = true

		// 计算评分
		deduction := float64(result.Count-result.Threshold) * 15.0
		result.Score = 100.0 - deduction
		if result.Score < 0 {
			result.Score = 0
		}

		result.Message = fmt.Sprintf("检测到 %d 个虚假范围表达，超过阈值 %d", result.Count, result.Threshold)
	} else {
		result.Message = "未检测到虚假范围表达"
	}

	return result
}

// GetType 获取规则类型
func (r *FalseRangeRule) GetType() models.RuleType {
	return models.RuleTypeFalseRange
}

// GetName 获取规则名称
func (r *FalseRangeRule) GetName() string {
	return "虚假范围表达检测"
}

// GetDescription 获取规则描述
func (r *FalseRangeRule) GetDescription() string {
	return "检测不连续的'from X to Y'范围表达"
}

// getContext 获取匹配项的上下文
func (r *FalseRangeRule) getContext(text string, offset, length int) string {
	const contextSize = 50

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

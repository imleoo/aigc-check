package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
	"github.com/leoobai/aigc-check/internal/text"
)

// MarkdownRule Signal 6: Markdown残留检测
type MarkdownRule struct {
	config    *config.Config
	processor *text.TextProcessor
}

// NewMarkdownRule 创建Markdown残留检测规则
func NewMarkdownRule(cfg *config.Config) *MarkdownRule {
	return &MarkdownRule{
		config:    cfg,
		processor: text.NewTextProcessor(),
	}
}

// Check 执行规则检测
func (r *MarkdownRule) Check(text string) models.RuleResult {
	result := models.RuleResult{
		RuleType:    models.RuleTypeMarkdown,
		RuleName:    r.GetName(),
		Description: r.GetDescription(),
		Detected:    false,
		Score:       100.0,
		Severity:    models.SeverityHigh,
		Matches:     []models.Match{},
		Count:       0,
		Threshold:   r.config.Thresholds.MarkdownResidue.Threshold,
	}

	patterns := r.config.Thresholds.MarkdownResidue.Patterns

	// 检测各种Markdown模式
	for _, pattern := range patterns {
		r.detectPattern(text, pattern, &result)
	}

	// 检查是否超过阈值
	if result.Count >= result.Threshold {
		result.Detected = true

		// 计算评分
		deduction := float64(result.Count-result.Threshold) * 12.0
		result.Score = 100.0 - deduction
		if result.Score < 0 {
			result.Score = 0
		}

		result.Message = fmt.Sprintf("检测到 %d 个Markdown格式残留，超过阈值 %d", result.Count, result.Threshold)
	} else {
		result.Message = "未检测到Markdown格式残留"
	}

	return result
}

// GetType 获取规则类型
func (r *MarkdownRule) GetType() models.RuleType {
	return models.RuleTypeMarkdown
}

// GetName 获取规则名称
func (r *MarkdownRule) GetName() string {
	return "Markdown残留检测"
}

// GetDescription 获取规则描述
func (r *MarkdownRule) GetDescription() string {
	return "检测未清理的Markdown格式标记，如 ##, **, []()"
}

// detectPattern 检测模式
func (r *MarkdownRule) detectPattern(text, pattern string, result *models.RuleResult) {
	// 尝试作为正则表达式
	re, err := regexp.Compile(pattern)
	if err != nil {
		// 如果不是正则表达式，作为普通字符串搜索
		r.detectLiteral(text, pattern, result)
		return
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
			Reason:  fmt.Sprintf("检测到Markdown格式: %s", pattern),
		})
	}
}

// detectLiteral 检测字面字符串
func (r *MarkdownRule) detectLiteral(text, pattern string, result *models.RuleResult) {
	offset := 0
	for {
		index := strings.Index(text[offset:], pattern)
		if index == -1 {
			break
		}

		actualOffset := offset + index
		line, column := r.processor.GetLineColumn(text, actualOffset)

		result.Count++
		result.Matches = append(result.Matches, models.Match{
			Text: pattern,
			Position: models.Position{
				Line:   line,
				Column: column,
				Offset: actualOffset,
				Length: len(pattern),
			},
			Context: r.getContext(text, actualOffset, len(pattern)),
			Reason:  fmt.Sprintf("检测到Markdown格式: %s", pattern),
		})

		offset = actualOffset + len(pattern)
	}
}

// getContext 获取匹配项的上下文
func (r *MarkdownRule) getContext(text string, offset, length int) string {
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

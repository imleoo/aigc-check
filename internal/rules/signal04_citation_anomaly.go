package rules

import (
	"fmt"
	"strings"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
	"github.com/leoobai/aigc-check/internal/text"
)

// CitationAnomalyRule Signal 4: 引用异常检测
type CitationAnomalyRule struct {
	config    *config.Config
	processor *text.TextProcessor
}

// NewCitationAnomalyRule 创建引用异常检测规则
func NewCitationAnomalyRule(cfg *config.Config) *CitationAnomalyRule {
	return &CitationAnomalyRule{
		config:    cfg,
		processor: text.NewTextProcessor(),
	}
}

// Check 执行规则检测
func (r *CitationAnomalyRule) Check(text string) models.RuleResult {
	result := models.RuleResult{
		RuleType:    models.RuleTypeCitationAnomaly,
		RuleName:    r.GetName(),
		Description: r.GetDescription(),
		Detected:    false,
		Score:       100.0,
		Severity:    models.SeverityCritical,
		Matches:     []models.Match{},
		Count:       0,
		Threshold:   1, // 任何引用异常都是严重问题
	}

	// 检测UTM参数
	r.checkUTMPatterns(text, &result)

	// 检测幽灵标记
	r.checkGhostMarkers(text, &result)

	// 检测占位符日期
	r.checkPlaceholderDates(text, &result)

	// 如果检测到任何异常
	if result.Count > 0 {
		result.Detected = true
		result.Score = 0.0 // 引用异常是严重问题，直接0分

		result.Message = fmt.Sprintf("检测到 %d 个引用异常标记，这是AI生成内容的明确证据", result.Count)
	} else {
		result.Message = "未检测到引用异常"
	}

	return result
}

// GetType 获取规则类型
func (r *CitationAnomalyRule) GetType() models.RuleType {
	return models.RuleTypeCitationAnomaly
}

// GetName 获取规则名称
func (r *CitationAnomalyRule) GetName() string {
	return "引用异常检测"
}

// GetDescription 获取规则描述
func (r *CitationAnomalyRule) GetDescription() string {
	return "检测AI生成内容的引用异常，包括UTM参数、幽灵标记和占位符日期"
}

// checkUTMPatterns 检测UTM参数
func (r *CitationAnomalyRule) checkUTMPatterns(text string, result *models.RuleResult) {
	patterns := r.config.Thresholds.CitationAnomaly.UTMPatterns

	for _, pattern := range patterns {
		positions := r.processor.FindPattern(text, pattern, false)
		for _, pos := range positions {
			result.Count++
			result.Matches = append(result.Matches, models.Match{
				Text:     pattern,
				Position: pos,
				Context:  r.getContext(text, pos.Offset, pos.Length),
				Reason:   fmt.Sprintf("检测到AI生成的UTM参数: %s", pattern),
			})
		}
	}
}

// checkGhostMarkers 检测幽灵标记
func (r *CitationAnomalyRule) checkGhostMarkers(text string, result *models.RuleResult) {
	markers := r.config.Thresholds.CitationAnomaly.GhostMarkers

	for _, marker := range markers {
		positions := r.processor.FindPattern(text, marker, false)
		for _, pos := range positions {
			result.Count++
			result.Matches = append(result.Matches, models.Match{
				Text:     marker,
				Position: pos,
				Context:  r.getContext(text, pos.Offset, pos.Length),
				Reason:   fmt.Sprintf("检测到AI生成的幽灵标记: %s", marker),
			})
		}
	}
}

// checkPlaceholderDates 检测占位符日期
func (r *CitationAnomalyRule) checkPlaceholderDates(text string, result *models.RuleResult) {
	placeholders := r.config.Thresholds.CitationAnomaly.PlaceholderDates

	for _, placeholder := range placeholders {
		positions := r.processor.FindPattern(text, placeholder, false)
		for _, pos := range positions {
			result.Count++
			result.Matches = append(result.Matches, models.Match{
				Text:     placeholder,
				Position: pos,
				Context:  r.getContext(text, pos.Offset, pos.Length),
				Reason:   fmt.Sprintf("检测到占位符日期: %s", placeholder),
			})
		}
	}
}

// getContext 获取匹配项的上下文
func (r *CitationAnomalyRule) getContext(text string, offset, length int) string {
	const contextSize = 50 // 前后各50个字符

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

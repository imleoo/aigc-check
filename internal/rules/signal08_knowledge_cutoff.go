package rules

import (
	"fmt"
	"strings"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
	"github.com/leoobai/aigc-check/internal/text"
)

// KnowledgeCutoffRule Signal 8: 知识截止日期检测
type KnowledgeCutoffRule struct {
	config    *config.Config
	processor *text.TextProcessor
}

// NewKnowledgeCutoffRule 创建知识截止日期检测规则
func NewKnowledgeCutoffRule(cfg *config.Config) *KnowledgeCutoffRule {
	return &KnowledgeCutoffRule{
		config:    cfg,
		processor: text.NewTextProcessor(),
	}
}

// Check 执行规则检测
func (r *KnowledgeCutoffRule) Check(text string) models.RuleResult {
	result := models.RuleResult{
		RuleType:    models.RuleTypeKnowledgeCutoff,
		RuleName:    r.GetName(),
		Description: r.GetDescription(),
		Detected:    false,
		Score:       100.0,
		Severity:    models.SeverityCritical,
		Matches:     []models.Match{},
		Count:       0,
		Threshold:   1, // 任何知识截止短语都是严重问题
	}

	phrases := r.config.Thresholds.KnowledgeCutoff.Phrases

	// 检测每个短语
	for _, phrase := range phrases {
		positions := r.processor.FindPattern(text, phrase, false)
		for _, pos := range positions {
			result.Count++
			result.Matches = append(result.Matches, models.Match{
				Text:     phrase,
				Position: pos,
				Context:  r.getContext(text, pos.Offset, pos.Length),
				Reason:   fmt.Sprintf("检测到AI知识截止短语: %s", phrase),
			})
		}
	}

	// 如果检测到任何知识截止短语
	if result.Count > 0 {
		result.Detected = true
		result.Score = 0.0 // 知识截止短语是明确的AI标记，直接0分

		result.Message = fmt.Sprintf("检测到 %d 个知识截止短语，这是AI生成内容的明确证据", result.Count)
	} else {
		result.Message = "未检测到知识截止短语"
	}

	return result
}

// GetType 获取规则类型
func (r *KnowledgeCutoffRule) GetType() models.RuleType {
	return models.RuleTypeKnowledgeCutoff
}

// GetName 获取规则名称
func (r *KnowledgeCutoffRule) GetName() string {
	return "知识截止日期检测"
}

// GetDescription 获取规则描述
func (r *KnowledgeCutoffRule) GetDescription() string {
	return "检测AI模型的知识截止日期相关短语，如'截至我的知识更新'"
}

// getContext 获取匹配项的上下文
func (r *KnowledgeCutoffRule) getContext(text string, offset, length int) string {
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

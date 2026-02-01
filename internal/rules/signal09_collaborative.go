package rules

import (
	"fmt"
	"strings"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
	"github.com/leoobai/aigc-check/internal/text"
)

// CollaborativeRule Signal 9: 协作式语气检测
type CollaborativeRule struct {
	config    *config.Config
	processor *text.TextProcessor
}

// NewCollaborativeRule 创建协作式语气检测规则
func NewCollaborativeRule(cfg *config.Config) *CollaborativeRule {
	return &CollaborativeRule{
		config:    cfg,
		processor: text.NewTextProcessor(),
	}
}

// Check 执行规则检测
func (r *CollaborativeRule) Check(text string) models.RuleResult {
	result := models.RuleResult{
		RuleType:    models.RuleTypeCollaborative,
		RuleName:    r.GetName(),
		Description: r.GetDescription(),
		Detected:    false,
		Score:       100.0,
		Severity:    models.SeverityHigh,
		Matches:     []models.Match{},
		Count:       0,
		Threshold:   r.config.Thresholds.CollaborativeTone.Threshold,
	}

	phrases := r.config.Thresholds.CollaborativeTone.Phrases

	// 检测每个短语
	for _, phrase := range phrases {
		positions := r.processor.FindPattern(text, phrase, false)
		for _, pos := range positions {
			result.Count++
			result.Matches = append(result.Matches, models.Match{
				Text:     phrase,
				Position: pos,
				Context:  r.getContext(text, pos.Offset, pos.Length),
				Reason:   fmt.Sprintf("检测到协作式语气: %s", phrase),
			})
		}
	}

	// 检查是否超过阈值
	if result.Count >= result.Threshold {
		result.Detected = true

		// 计算评分
		deduction := float64(result.Count-result.Threshold) * 10.0
		result.Score = 100.0 - deduction
		if result.Score < 0 {
			result.Score = 0
		}

		result.Message = fmt.Sprintf("检测到 %d 个协作式语气短语，超过阈值 %d", result.Count, result.Threshold)
	} else {
		result.Message = "未检测到异常的协作式语气"
	}

	return result
}

// GetType 获取规则类型
func (r *CollaborativeRule) GetType() models.RuleType {
	return models.RuleTypeCollaborative
}

// GetName 获取规则名称
func (r *CollaborativeRule) GetName() string {
	return "协作式语气检测"
}

// GetDescription 获取规则描述
func (r *CollaborativeRule) GetDescription() string {
	return "检测AI助手常用的协作式语气，如'希望这能帮到你'、'随时告诉我'"
}

// getContext 获取匹配项的上下文
func (r *CollaborativeRule) getContext(text string, offset, length int) string {
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

package rules

import (
	"fmt"
	"strings"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
	"github.com/leoobai/aigc-check/internal/text"
)

// PerfectionismRule Signal 10: 完美主义陷阱检测
type PerfectionismRule struct {
	config    *config.Config
	processor *text.TextProcessor
}

// NewPerfectionismRule 创建完美主义陷阱检测规则
func NewPerfectionismRule(cfg *config.Config) *PerfectionismRule {
	return &PerfectionismRule{
		config:    cfg,
		processor: text.NewTextProcessor(),
	}
}

// Check 执行规则检测
func (r *PerfectionismRule) Check(text string) models.RuleResult {
	result := models.RuleResult{
		RuleType:    models.RuleTypePerfectionism,
		RuleName:    r.GetName(),
		Description: r.GetDescription(),
		Detected:    false,
		Score:       100.0,
		Severity:    models.SeverityHigh,
		Matches:     []models.Match{},
		Count:       0,
		Threshold:   r.config.Thresholds.Perfectionism.Threshold,
	}

	// 检测第一人称代词
	firstPersonCount := r.detectFirstPerson(text, &result)

	// 检测情感词汇
	emotionalCount := r.detectEmotionalWords(text, &result)

	// 检测不确定性标记
	uncertaintyCount := r.detectUncertainty(text, &result)

	// 计算总分
	totalIndicators := firstPersonCount + emotionalCount + uncertaintyCount

	// 如果缺乏个人化表达（指标太少），说明可能是AI生成
	if totalIndicators < result.Threshold {
		result.Detected = true

		// 计算评分（个人化表达越少，分数越低）
		deficit := float64(result.Threshold - totalIndicators)
		deduction := deficit * 10.0
		result.Score = 100.0 - deduction
		if result.Score < 0 {
			result.Score = 0
		}

		result.Message = fmt.Sprintf("缺乏个人化表达：仅检测到 %d 个个人化指标，低于阈值 %d", totalIndicators, result.Threshold)
	} else {
		result.Message = fmt.Sprintf("检测到 %d 个个人化表达指标，表现正常", totalIndicators)
	}

	return result
}

// GetType 获取规则类型
func (r *PerfectionismRule) GetType() models.RuleType {
	return models.RuleTypePerfectionism
}

// GetName 获取规则名称
func (r *PerfectionismRule) GetName() string {
	return "完美主义陷阱检测"
}

// GetDescription 获取规则描述
func (r *PerfectionismRule) GetDescription() string {
	return "检测文本中个人化表达的缺失，包括第一人称、情感词汇和不确定性标记"
}

// detectFirstPerson 检测第一人称代词
func (r *PerfectionismRule) detectFirstPerson(text string, result *models.RuleResult) int {
	pronouns := r.config.Thresholds.Perfectionism.FirstPersonPronouns
	count := 0

	for _, pronoun := range pronouns {
		positions := r.processor.FindPattern(text, pronoun, false)
		for _, pos := range positions {
			count++
			result.Count++
			result.Matches = append(result.Matches, models.Match{
				Text:     pronoun,
				Position: pos,
				Context:  r.getContext(text, pos.Offset, pos.Length),
				Reason:   "第一人称代词",
			})
		}
	}

	return count
}

// detectEmotionalWords 检测情感词汇
func (r *PerfectionismRule) detectEmotionalWords(text string, result *models.RuleResult) int {
	words := r.config.Thresholds.Perfectionism.EmotionalWords
	count := 0

	for _, word := range words {
		positions := r.processor.FindPattern(text, word, false)
		for _, pos := range positions {
			count++
			result.Count++
			result.Matches = append(result.Matches, models.Match{
				Text:     word,
				Position: pos,
				Context:  r.getContext(text, pos.Offset, pos.Length),
				Reason:   "情感词汇",
			})
		}
	}

	return count
}

// detectUncertainty 检测不确定性标记
func (r *PerfectionismRule) detectUncertainty(text string, result *models.RuleResult) int {
	markers := r.config.Thresholds.Perfectionism.UncertaintyMarkers
	count := 0

	for _, marker := range markers {
		positions := r.processor.FindPattern(text, marker, false)
		for _, pos := range positions {
			count++
			result.Count++
			result.Matches = append(result.Matches, models.Match{
				Text:     marker,
				Position: pos,
				Context:  r.getContext(text, pos.Offset, pos.Length),
				Reason:   "不确定性标记",
			})
		}
	}

	return count
}

// getContext 获取匹配项的上下文
func (r *PerfectionismRule) getContext(text string, offset, length int) string {
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

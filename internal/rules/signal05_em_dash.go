package rules

import (
	"fmt"
	"strings"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
	"github.com/leoobai/aigc-check/internal/text"
)

// EmDashRule Signal 5: 破折号密度检测
type EmDashRule struct {
	config    *config.Config
	processor *text.TextProcessor
}

// NewEmDashRule 创建破折号密度检测规则
func NewEmDashRule(cfg *config.Config) *EmDashRule {
	return &EmDashRule{
		config:    cfg,
		processor: text.NewTextProcessor(),
	}
}

// Check 执行规则检测
func (r *EmDashRule) Check(text string) models.RuleResult {
	result := models.RuleResult{
		RuleType:    models.RuleTypeEmDash,
		RuleName:    r.GetName(),
		Description: r.GetDescription(),
		Detected:    false,
		Score:       100.0,
		Severity:    models.SeverityMedium,
		Matches:     []models.Match{},
		Count:       0,
		Threshold:   int(r.config.Thresholds.EmDashDensity),
	}

	// 统计破折号（em dash: —）
	emDash := "—"
	count := strings.Count(text, emDash)

	// 计算文本长度（字符数）
	charCount := len([]rune(text))
	if charCount == 0 {
		result.Message = "文本为空"
		return result
	}

	// 计算密度（每千字）
	density := float64(count) * 1000.0 / float64(charCount)

	result.Count = count

	// 查找所有破折号位置
	offset := 0
	for {
		index := strings.Index(text[offset:], emDash)
		if index == -1 {
			break
		}

		actualOffset := offset + index
		line, column := r.processor.GetLineColumn(text, actualOffset)

		result.Matches = append(result.Matches, models.Match{
			Text: emDash,
			Position: models.Position{
				Line:   line,
				Column: column,
				Offset: actualOffset,
				Length: len(emDash),
			},
			Context: r.getContext(text, actualOffset, len(emDash)),
			Reason:  "破折号使用",
		})

		offset = actualOffset + len(emDash)
	}

	// 检查密度是否超过阈值
	if density >= r.config.Thresholds.EmDashDensity {
		result.Detected = true

		// 计算评分
		excessDensity := density - r.config.Thresholds.EmDashDensity
		deduction := excessDensity * 5.0
		result.Score = 100.0 - deduction
		if result.Score < 0 {
			result.Score = 0
		}

		result.Message = fmt.Sprintf("破折号密度 %.2f/千字，超过阈值 %.2f/千字", density, r.config.Thresholds.EmDashDensity)
	} else {
		result.Message = fmt.Sprintf("破折号密度 %.2f/千字，正常", density)
	}

	return result
}

// GetType 获取规则类型
func (r *EmDashRule) GetType() models.RuleType {
	return models.RuleTypeEmDash
}

// GetName 获取规则名称
func (r *EmDashRule) GetName() string {
	return "破折号密度检测"
}

// GetDescription 获取规则描述
func (r *EmDashRule) GetDescription() string {
	return "检测破折号（—）的使用密度，AI生成内容倾向于过度使用破折号"
}

// getContext 获取匹配项的上下文
func (r *EmDashRule) getContext(text string, offset, length int) string {
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

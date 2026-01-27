package reporter

import (
	"fmt"
	"strings"

	"github.com/leoobai/aigc-check/internal/models"
)

// TextReporter 文本报告生成器
type TextReporter struct {
	colorEnabled bool
}

// NewTextReporter 创建文本报告生成器
func NewTextReporter(colorEnabled bool) *TextReporter {
	return &TextReporter{
		colorEnabled: colorEnabled,
	}
}

// Generate 生成文本报告
func (r *TextReporter) Generate(result *models.DetectionResult) (string, error) {
	var sb strings.Builder

	// 标题
	sb.WriteString("╔═══════════════════════════════════════════════════════════════╗\n")
	sb.WriteString("║           AIGC-Check 检测报告                                  ║\n")
	sb.WriteString("╚═══════════════════════════════════════════════════════════════╝\n\n")

	// 总体评分
	r.writeOverallScore(&sb, result)

	// 风险等级
	r.writeRiskLevel(&sb, result)

	// 维度评分
	r.writeDimensionScores(&sb, result)

	// 检测到的问题
	r.writeDetectedIssues(&sb, result)

	// 改进建议
	r.writeSuggestions(&sb, result)

	// 处理时间
	sb.WriteString(fmt.Sprintf("\n处理时间: %v\n", result.ProcessTime))
	sb.WriteString(fmt.Sprintf("检测时间: %s\n", result.DetectedAt.Format("2006-01-02 15:04:05")))

	return sb.String(), nil
}

// Format 获取报告格式
func (r *TextReporter) Format() string {
	return "text"
}

// writeOverallScore 写入总体评分
func (r *TextReporter) writeOverallScore(sb *strings.Builder, result *models.DetectionResult) {
	sb.WriteString("【总体评分】\n")
	sb.WriteString(strings.Repeat("─", 60) + "\n")

	scoreBar := r.createScoreBar(result.Score.Total)
	scoreColor := r.getScoreColor(result.Score.Total)

	if r.colorEnabled {
		sb.WriteString(fmt.Sprintf("%s%.1f / 100%s\n", scoreColor, result.Score.Total, r.colorReset()))
	} else {
		sb.WriteString(fmt.Sprintf("%.1f / 100\n", result.Score.Total))
	}

	sb.WriteString(scoreBar + "\n\n")
}

// writeRiskLevel 写入风险等级
func (r *TextReporter) writeRiskLevel(sb *strings.Builder, result *models.DetectionResult) {
	sb.WriteString("【风险等级】\n")
	sb.WriteString(strings.Repeat("─", 60) + "\n")

	riskIcon := r.getRiskIcon(result.RiskLevel)
	riskColor := r.getRiskColor(result.RiskLevel)

	if r.colorEnabled {
		sb.WriteString(fmt.Sprintf("%s%s %s%s\n\n", riskColor, riskIcon, result.RiskLevel.Description(), r.colorReset()))
	} else {
		sb.WriteString(fmt.Sprintf("%s %s\n\n", riskIcon, result.RiskLevel.Description()))
	}
}

// writeDimensionScores 写入维度评分
func (r *TextReporter) writeDimensionScores(sb *strings.Builder, result *models.DetectionResult) {
	sb.WriteString("【维度评分】\n")
	sb.WriteString(strings.Repeat("─", 60) + "\n")

	dimensions := []struct {
		name  string
		score models.DimensionScore
	}{
		{"词汇多样性", result.Score.Dimensions.VocabularyDiversity},
		{"句式复杂度", result.Score.Dimensions.SentenceComplexity},
		{"个人化表达", result.Score.Dimensions.Personalization},
		{"逻辑连贯性", result.Score.Dimensions.LogicalCoherence},
		{"情感真实度", result.Score.Dimensions.EmotionalAuthenticity},
	}

	for _, dim := range dimensions {
		percentage := dim.score.Percentage
		bar := r.createPercentageBar(percentage)

		sb.WriteString(fmt.Sprintf("%-12s %.1f/%.0f (%.0f%%) [%s] %s\n",
			dim.name,
			dim.score.Score,
			dim.score.MaxScore,
			percentage,
			dim.score.Level,
			bar,
		))

		if len(dim.score.Issues) > 0 {
			for _, issue := range dim.score.Issues {
				sb.WriteString(fmt.Sprintf("  ⚠ %s\n", issue))
			}
		}
		sb.WriteString("\n")
	}
}

// writeDetectedIssues 写入检测到的问题
func (r *TextReporter) writeDetectedIssues(sb *strings.Builder, result *models.DetectionResult) {
	sb.WriteString("【检测到的问题】\n")
	sb.WriteString(strings.Repeat("─", 60) + "\n")

	detectedCount := 0
	for _, ruleResult := range result.RuleResults {
		if ruleResult.Detected {
			detectedCount++
		}
	}

	if detectedCount == 0 {
		sb.WriteString("✓ 未检测到明显的AI生成特征\n\n")
		return
	}

	sb.WriteString(fmt.Sprintf("检测到 %d 个问题：\n\n", detectedCount))

	for _, ruleResult := range result.RuleResults {
		if !ruleResult.Detected {
			continue
		}

		severityIcon := r.getSeverityIcon(ruleResult.Severity)
		sb.WriteString(fmt.Sprintf("%s [%s] %s\n", severityIcon, ruleResult.Severity, ruleResult.RuleName))
		sb.WriteString(fmt.Sprintf("   评分: %.1f/100\n", ruleResult.Score))
		sb.WriteString(fmt.Sprintf("   消息: %s\n", ruleResult.Message))
		sb.WriteString(fmt.Sprintf("   匹配数: %d (阈值: %d)\n", ruleResult.Count, ruleResult.Threshold))

		// 显示前3个匹配项
		if len(ruleResult.Matches) > 0 {
			sb.WriteString("   示例:\n")
			maxShow := 3
			if len(ruleResult.Matches) < maxShow {
				maxShow = len(ruleResult.Matches)
			}
			for i := 0; i < maxShow; i++ {
				match := ruleResult.Matches[i]
				sb.WriteString(fmt.Sprintf("     - 行%d: %s\n", match.Position.Line, match.Text))
			}
			if len(ruleResult.Matches) > maxShow {
				sb.WriteString(fmt.Sprintf("     ... 还有 %d 个匹配项\n", len(ruleResult.Matches)-maxShow))
			}
		}
		sb.WriteString("\n")
	}
}

// writeSuggestions 写入改进建议
func (r *TextReporter) writeSuggestions(sb *strings.Builder, result *models.DetectionResult) {
	if len(result.Suggestions) == 0 {
		return
	}

	sb.WriteString("【改进建议】\n")
	sb.WriteString(strings.Repeat("─", 60) + "\n")

	for i, suggestion := range result.Suggestions {
		priorityIcon := r.getPriorityIcon(suggestion.Priority)
		sb.WriteString(fmt.Sprintf("%d. %s [%s] %s\n",
			i+1,
			priorityIcon,
			models.GetCategoryName(suggestion.Category),
			suggestion.Title,
		))
		sb.WriteString(fmt.Sprintf("   %s\n", suggestion.Description))

		if len(suggestion.Examples) > 0 {
			sb.WriteString("   示例:\n")
			for _, example := range suggestion.Examples {
				sb.WriteString(fmt.Sprintf("     修改前: %s\n", example.Before))
				sb.WriteString(fmt.Sprintf("     修改后: %s\n", example.After))
				if example.Reason != "" {
					sb.WriteString(fmt.Sprintf("     原因: %s\n", example.Reason))
				}
			}
		}
		sb.WriteString("\n")
	}
}

// createScoreBar 创建评分条
func (r *TextReporter) createScoreBar(score float64) string {
	barLength := 50
	filled := int(score / 100.0 * float64(barLength))
	if filled > barLength {
		filled = barLength
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barLength-filled)
	return fmt.Sprintf("[%s]", bar)
}

// createPercentageBar 创建百分比条
func (r *TextReporter) createPercentageBar(percentage float64) string {
	barLength := 20
	filled := int(percentage / 100.0 * float64(barLength))
	if filled > barLength {
		filled = barLength
	}

	return strings.Repeat("█", filled) + strings.Repeat("░", barLength-filled)
}

// getScoreColor 获取评分颜色
func (r *TextReporter) getScoreColor(score float64) string {
	if !r.colorEnabled {
		return ""
	}

	switch {
	case score >= 76:
		return "\033[32m" // 绿色
	case score >= 61:
		return "\033[33m" // 黄色
	case score >= 41:
		return "\033[31m" // 红色
	default:
		return "\033[35m" // 紫色
	}
}

// getRiskColor 获取风险等级颜色
func (r *TextReporter) getRiskColor(level models.RiskLevel) string {
	if !r.colorEnabled {
		return ""
	}

	switch level {
	case models.RiskLevelLow:
		return "\033[32m" // 绿色
	case models.RiskLevelMedium:
		return "\033[33m" // 黄色
	case models.RiskLevelHigh:
		return "\033[31m" // 红色
	case models.RiskLevelVeryHigh:
		return "\033[35m" // 紫色
	default:
		return ""
	}
}

// colorReset 重置颜色
func (r *TextReporter) colorReset() string {
	if !r.colorEnabled {
		return ""
	}
	return "\033[0m"
}

// getRiskIcon 获取风险等级图标
func (r *TextReporter) getRiskIcon(level models.RiskLevel) string {
	switch level {
	case models.RiskLevelLow:
		return "✓"
	case models.RiskLevelMedium:
		return "⚠"
	case models.RiskLevelHigh:
		return "⚠⚠"
	case models.RiskLevelVeryHigh:
		return "⚠⚠⚠"
	default:
		return "?"
	}
}

// getSeverityIcon 获取严重程度图标
func (r *TextReporter) getSeverityIcon(severity models.Severity) string {
	switch severity {
	case models.SeverityCritical:
		return "🔴"
	case models.SeverityHigh:
		return "🟠"
	case models.SeverityMedium:
		return "🟡"
	case models.SeverityLow:
		return "🟢"
	default:
		return "⚪"
	}
}

// getPriorityIcon 获取优先级图标
func (r *TextReporter) getPriorityIcon(priority models.Priority) string {
	switch priority {
	case models.PriorityHigh:
		return "🔴"
	case models.PriorityMedium:
		return "🟡"
	case models.PriorityLow:
		return "🟢"
	default:
		return "⚪"
	}
}

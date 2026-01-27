package analyzer

import (
	"time"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/detector"
	"github.com/leoobai/aigc-check/internal/models"
	"github.com/leoobai/aigc-check/internal/rules"
	"github.com/leoobai/aigc-check/internal/scorer"
	"github.com/leoobai/aigc-check/internal/text"
)

// Analyzer 主分析器
type Analyzer struct {
	config     *config.Config
	ruleEngine *detector.RuleEngine
	scorer     *scorer.Calculator
	processor  *text.TextProcessor
}

// NewAnalyzer 创建分析器
func NewAnalyzer(cfg *config.Config) *Analyzer {
	// 创建规则引擎
	ruleEngine := detector.NewRuleEngine(cfg)

	// 注册所有规则
	ruleEngine.RegisterRule(rules.NewHighFreqWordsRule(cfg))
	ruleEngine.RegisterRule(rules.NewSentenceStartersRule(cfg))
	ruleEngine.RegisterRule(rules.NewFalseRangeRule(cfg))
	ruleEngine.RegisterRule(rules.NewCitationAnomalyRule(cfg))
	ruleEngine.RegisterRule(rules.NewEmDashRule(cfg))
	ruleEngine.RegisterRule(rules.NewMarkdownRule(cfg))
	ruleEngine.RegisterRule(rules.NewEmojiRule(cfg))
	ruleEngine.RegisterRule(rules.NewKnowledgeCutoffRule(cfg))
	ruleEngine.RegisterRule(rules.NewCollaborativeRule(cfg))
	ruleEngine.RegisterRule(rules.NewPerfectionismRule(cfg))

	return &Analyzer{
		config:     cfg,
		ruleEngine: ruleEngine,
		scorer:     scorer.NewCalculator(cfg),
		processor:  text.NewTextProcessor(),
	}
}

// Analyze 执行完整分析
func (a *Analyzer) Analyze(request models.DetectionRequest) (*models.DetectionResult, error) {
	startTime := time.Now()

	// 执行规则检测
	ruleResults := a.ruleEngine.Check(request.Text)

	// 计算评分
	score := a.scorer.Calculate(ruleResults)

	// 生成建议
	suggestions := a.generateSuggestions(ruleResults)

	// 构建结果
	result := &models.DetectionResult{
		RequestID:   generateRequestID(),
		Text:        request.Text,
		Score:       score,
		RuleResults: ruleResults,
		Suggestions: suggestions,
		RiskLevel:   models.GetRiskLevel(score.Total),
		ProcessTime: time.Since(startTime),
		DetectedAt:  time.Now(),
	}

	return result, nil
}

// generateSuggestions 生成改进建议
func (a *Analyzer) generateSuggestions(results []models.RuleResult) []models.Suggestion {
	var suggestions []models.Suggestion

	for _, result := range results {
		if !result.Detected {
			continue
		}

		// 根据规则类型生成建议
		switch result.RuleType {
		case models.RuleTypeHighFreqWords:
			suggestion := models.NewSuggestion(
				models.CategoryVocabulary,
				models.PriorityHigh,
				"减少AI常用高频词汇",
				"避免过度使用 crucial, pivotal, vital 等AI常用词汇，使用更自然、多样化的表达方式。",
				result.RuleType,
			)
			suggestion.AddExample(
				"This is a crucial step in the process.",
				"This is an important step in the process.",
				"使用更自然的词汇替代AI高频词",
			)
			suggestions = append(suggestions, suggestion)

		case models.RuleTypeSentenceStarters:
			suggestion := models.NewSuggestion(
				models.CategorySentence,
				models.PriorityHigh,
				"增加句式开头的多样性",
				"避免重复使用 Additionally, Furthermore, Moreover 等连接词开头，尝试使用更多样化的句式结构。",
				result.RuleType,
			)
			suggestion.AddExample(
				"Additionally, we need to consider the cost. Furthermore, the timeline is important.",
				"We also need to consider the cost. The timeline matters too.",
				"使用更自然的连接方式",
			)
			suggestions = append(suggestions, suggestion)

		case models.RuleTypeCitationAnomaly:
			suggestion := models.NewSuggestion(
				models.CategoryFormatting,
				models.PriorityHigh,
				"清理AI生成的引用标记",
				"移除所有AI生成的UTM参数、幽灵标记和占位符日期。",
				result.RuleType,
			)
			suggestions = append(suggestions, suggestion)

		case models.RuleTypeEmDash:
			suggestion := models.NewSuggestion(
				models.CategorySentence,
				models.PriorityMedium,
				"减少破折号的使用",
				"适度使用破折号（—），过度使用会显得不自然。考虑使用其他标点符号或句式结构。",
				result.RuleType,
			)
			suggestions = append(suggestions, suggestion)

		case models.RuleTypeMarkdown:
			suggestion := models.NewSuggestion(
				models.CategoryFormatting,
				models.PriorityHigh,
				"清理Markdown格式残留",
				"移除所有Markdown格式标记，如 ##, **, []() 等，确保文本格式干净。",
				result.RuleType,
			)
			suggestions = append(suggestions, suggestion)

		case models.RuleTypeKnowledgeCutoff:
			suggestion := models.NewSuggestion(
				models.CategoryAuthenticity,
				models.PriorityHigh,
				"移除AI知识截止短语",
				"删除'截至我的知识更新'等明显的AI特征短语。",
				result.RuleType,
			)
			suggestions = append(suggestions, suggestion)

		case models.RuleTypeCollaborative:
			suggestion := models.NewSuggestion(
				models.CategoryTone,
				models.PriorityHigh,
				"调整协作式语气",
				"移除'希望这能帮到你'等AI助手特有的协作式语气，使用更自然的表达方式。",
				result.RuleType,
			)
			suggestions = append(suggestions, suggestion)

		case models.RuleTypePerfectionism:
			suggestion := models.NewSuggestion(
				models.CategoryAuthenticity,
				models.PriorityHigh,
				"增加个人化表达",
				"适当使用第一人称、情感词汇和不确定性表达，使文本更具人类特征。",
				result.RuleType,
			)
			suggestion.AddExample(
				"The solution is optimal and will work perfectly.",
				"I think this solution should work well, though we might need to adjust it.",
				"增加个人观点和适度的不确定性",
			)
			suggestions = append(suggestions, suggestion)
		}
	}

	return suggestions
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return time.Now().Format("20060102150405")
}

package analyzer

import (
	"context"
	"fmt"
	"time"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/detector"
	"github.com/leoobai/aigc-check/internal/gemini"
	"github.com/leoobai/aigc-check/internal/models"
	"github.com/leoobai/aigc-check/internal/rules"
	"github.com/leoobai/aigc-check/internal/scorer"
	"github.com/leoobai/aigc-check/internal/statistics"
	"github.com/leoobai/aigc-check/internal/text"
)

// Analyzer 主分析器
type Analyzer struct {
	config           *config.Config
	ruleEngine       *detector.RuleEngine
	scorer           *scorer.Calculator
	processor        *text.TextProcessor
	statsAnalyzer    *statistics.Analyzer
	geminiAnalyzer   *gemini.Analyzer
	geminiSuggester  *gemini.Suggester
	multimodalConfig models.MultimodalConfig
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

	// 创建统计分析器
	statsAnalyzer := statistics.NewAnalyzer()

	// 创建 Gemini 分析器和建议器（如果启用）
	var geminiAnalyzer *gemini.Analyzer
	var geminiSuggester *gemini.Suggester
	if cfg.Gemini.Enabled {
		client, err := gemini.NewClient(cfg.Gemini)
		if err == nil {
			geminiAnalyzer = gemini.NewAnalyzer(client)
			geminiSuggester = gemini.NewSuggester(client)
		}
	}

	// 多模态配置
	multimodalConfig := models.DefaultMultimodalConfig
	multimodalConfig.EnableStatistics = true
	multimodalConfig.EnableSemantic = cfg.Gemini.Enabled

	return &Analyzer{
		config:           cfg,
		ruleEngine:       ruleEngine,
		scorer:           scorer.NewCalculator(cfg),
		processor:        text.NewTextProcessor(),
		statsAnalyzer:    statsAnalyzer,
		geminiAnalyzer:   geminiAnalyzer,
		geminiSuggester:  geminiSuggester,
		multimodalConfig: multimodalConfig,
	}
}

// Analyze 执行完整分析（支持多模态检测）
func (a *Analyzer) Analyze(request models.DetectionRequest) (*models.DetectionResult, error) {
	startTime := time.Now()
	ctx := context.Background()

	// 如果启用多模态检测，使用多层分析
	if a.multimodalConfig.Enabled {
		return a.analyzeMultimodal(ctx, request, startTime)
	}

	// 否则使用传统单层检测
	return a.analyzeSingleLayer(request, startTime)
}

// analyzeSingleLayer 单层检测（传统模式）
func (a *Analyzer) analyzeSingleLayer(request models.DetectionRequest, startTime time.Time) (*models.DetectionResult, error) {
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

// analyzeMultimodal 多模态检测（分层触发策略）
func (a *Analyzer) analyzeMultimodal(ctx context.Context, request models.DetectionRequest, startTime time.Time) (*models.DetectionResult, error) {
	// Layer 1: 规则检测
	ruleResults := a.ruleEngine.Check(request.Text)
	ruleScore := a.scorer.Calculate(ruleResults)
	ruleConfidence := a.calculateRuleConfidence(ruleResults, ruleScore)

	// 初始化多模态结果
	multimodal := &models.MultimodalResult{
		RuleLayerScore: ruleScore.Total,
		LayerWeights:   a.multimodalConfig.GetEffectiveWeights(),
		DetectionMode:  a.multimodalConfig.GetDetectionMode(),
		RuleLayerDetails: &models.RuleLayerDetails{
			DetectedRulesCount: countDetectedRules(ruleResults),
			TotalRulesCount:    len(ruleResults),
			RedFlagCount:       countRedFlags(ruleResults),
			MainIssues:         extractMainIssues(ruleResults),
		},
	}

	return a.continueMultimodalAnalysis(ctx, request, ruleResults, ruleScore, ruleConfidence, multimodal, startTime)
}

// continueMultimodalAnalysis 继续多模态分析
func (a *Analyzer) continueMultimodalAnalysis(
	ctx context.Context,
	request models.DetectionRequest,
	ruleResults []models.RuleResult,
	ruleScore models.Score,
	ruleConfidence float64,
	multimodal *models.MultimodalResult,
	startTime time.Time,
) (*models.DetectionResult, error) {
	thresholds := a.multimodalConfig.ConfidenceThresholds

	// 判断是否需要统计分析
	if a.multimodalConfig.EnableStatistics && models.NeedsStatisticsAnalysis(ruleConfidence, thresholds) {
		statsResult := a.statsAnalyzer.Analyze(request.Text)
		multimodal.StatisticsLayerScore = statsResult.HumanScore
		multimodal.StatisticsLayerDetails = &models.StatisticsLayerDetails{
			TypeTokenRatio:         statsResult.Vocabulary.TTR,
			VocabularyRichness:     statsResult.Vocabulary.Richness,
			SentenceLengthVariance: statsResult.Sentence.LengthStdDev,
			SentenceComplexity:     statsResult.Sentence.ComplexityScore,
			PerplexityScore:        statsResult.Perplexity.Score,
			AIProbability:          statsResult.AIProbability,
			Details:                statsResult.Details,
		}
	}

	return a.finalizeMultimodalResult(ctx, request, ruleResults, ruleScore, ruleConfidence, multimodal, startTime)
}

// finalizeMultimodalResult 完成多模态分析并生成最终结果
func (a *Analyzer) finalizeMultimodalResult(
	ctx context.Context,
	request models.DetectionRequest,
	ruleResults []models.RuleResult,
	ruleScore models.Score,
	ruleConfidence float64,
	multimodal *models.MultimodalResult,
	startTime time.Time,
) (*models.DetectionResult, error) {
	statsConfidence := ruleConfidence
	if multimodal.StatisticsLayerDetails != nil {
		statsConfidence = 1.0 - multimodal.StatisticsLayerDetails.AIProbability
	}

	thresholds := a.multimodalConfig.ConfidenceThresholds

	// 判断是否需要语义分析
	if a.multimodalConfig.EnableSemantic && a.geminiAnalyzer != nil &&
		models.NeedsSemanticAnalysis(ruleConfidence, statsConfidence, thresholds) {
		// 调用 Gemini 进行语义分析
		analysisResult, err := a.geminiAnalyzer.AnalyzeText(ctx, request.Text)
		if err == nil {
			// 将 AI 概率转换为人类分数 (100 - AI概率)
			humanScore := 100.0 - analysisResult.AIProbability
			multimodal.SemanticLayerScore = humanScore

			// 提取特征名称
			features := make([]string, len(analysisResult.Features))
			for i, f := range analysisResult.Features {
				features[i] = f.Name
			}

			multimodal.SemanticLayerDetails = &models.SemanticLayerDetails{
				CoherenceScore:       50.0, // 默认值，可以通过单独调用获取
				PersonalizationScore: 50.0, // 默认值
				AIPatternScore:       analysisResult.AIProbability,
				DetectedFeatures:     features,
				Explanation:          analysisResult.Explanation,
				FromCache:            false,
			}
		}
	}

	// 融合分数
	weights := multimodal.LayerWeights
	multimodal.FinalScore = models.FuseScores(
		multimodal.RuleLayerScore,
		multimodal.StatisticsLayerScore,
		multimodal.SemanticLayerScore,
		weights,
	)

	// 计算融合置信度
	semanticConf := 0.5
	if multimodal.SemanticLayerDetails != nil {
		semanticConf = 1.0 - (multimodal.SemanticLayerDetails.AIPatternScore / 100.0)
	}
	multimodal.Confidence = models.CalculateFusionConfidence(ruleConfidence, statsConfidence, semanticConf, weights)

	// 生成融合说明
	multimodal.FusionExplanation = a.generateFusionExplanation(multimodal)

	// 生成建议
	suggestions := a.generateSuggestions(ruleResults)

	// 如果启用了智能建议，添加 Gemini 建议
	if a.geminiSuggester != nil && multimodal.SemanticLayerDetails != nil {
		issues := extractIssuesFromResults(ruleResults)
		geminiSuggestions, err := a.geminiSuggester.GenerateSuggestions(ctx, request.Text, issues)
		if err == nil {
			suggestions = append(suggestions, convertGeminiSuggestions(geminiSuggestions)...)
		}
	}

	// 构建最终结果
	result := &models.DetectionResult{
		RequestID:   generateRequestID(),
		Text:        request.Text,
		Score:       models.Score{Total: multimodal.FinalScore, Dimensions: ruleScore.Dimensions},
		RuleResults: ruleResults,
		Suggestions: suggestions,
		RiskLevel:   models.GetRiskLevel(multimodal.FinalScore),
		ProcessTime: time.Since(startTime),
		DetectedAt:  time.Now(),
	}

	return result, nil
}

// calculateRuleConfidence 计算规则检测置信度
func (a *Analyzer) calculateRuleConfidence(results []models.RuleResult, score models.Score) float64 {
	detectedCount := countDetectedRules(results)
	redFlagCount := countRedFlags(results)

	// 基础置信度基于分数
	baseConfidence := score.Total / 100.0

	// 红旗规则增加置信度
	if redFlagCount > 0 {
		baseConfidence += float64(redFlagCount) * 0.1
	}

	// 检测到的规则数量影响置信度
	if detectedCount >= 5 {
		baseConfidence += 0.15
	} else if detectedCount >= 3 {
		baseConfidence += 0.10
	}

	// 确保在 0-1 范围内
	if baseConfidence > 1.0 {
		return 1.0
	}
	if baseConfidence < 0.0 {
		return 0.0
	}

	return baseConfidence
}

// countDetectedRules 统计检测到的规则数量
func countDetectedRules(results []models.RuleResult) int {
	count := 0
	for _, r := range results {
		if r.Detected {
			count++
		}
	}
	return count
}

// countRedFlags 统计红旗规则数量
func countRedFlags(results []models.RuleResult) int {
	count := 0
	redFlagTypes := map[models.RuleType]bool{
		models.RuleTypeCitationAnomaly:  true,
		models.RuleTypeKnowledgeCutoff:  true,
		models.RuleTypeMarkdown:         true,
		models.RuleTypeEmoji:            true,
	}

	for _, r := range results {
		if r.Detected && redFlagTypes[r.RuleType] {
			count++
		}
	}
	return count
}

// extractMainIssues 提取主要问题
func extractMainIssues(results []models.RuleResult) []string {
	var issues []string
	for _, r := range results {
		if r.Detected {
			issues = append(issues, string(r.RuleType))
		}
	}
	return issues
}

// generateFusionExplanation 生成融合说明
func (a *Analyzer) generateFusionExplanation(multimodal *models.MultimodalResult) string {
	mode := multimodal.DetectionMode

	switch mode {
	case models.DetectionModeRuleOnly:
		return "仅使用规则检测"
	case models.DetectionModeRuleStatistics:
		return fmt.Sprintf("规则检测(%.1f) + 统计分析(%.1f) 融合",
			multimodal.RuleLayerScore, multimodal.StatisticsLayerScore)
	case models.DetectionModeMultimodal:
		return fmt.Sprintf("规则检测(%.1f) + 统计分析(%.1f) + 语义分析(%.1f) 融合",
			multimodal.RuleLayerScore, multimodal.StatisticsLayerScore, multimodal.SemanticLayerScore)
	default:
		return "未知检测模式"
	}
}

// extractIssuesFromResults 从规则结果中提取问题描述
func extractIssuesFromResults(results []models.RuleResult) []string {
	var issues []string
	for _, r := range results {
		if r.Detected {
			issues = append(issues, r.Message)
		}
	}
	return issues
}

// convertGeminiSuggestions 转换 Gemini 建议为标准建议格式
func convertGeminiSuggestions(geminiSuggestions []gemini.Suggestion) []models.Suggestion {
	suggestions := make([]models.Suggestion, 0, len(geminiSuggestions))
	for _, gs := range geminiSuggestions {
		// 根据优先级映射
		priority := models.PriorityMedium
		if gs.Priority <= 2 {
			priority = models.PriorityHigh
		} else if gs.Priority >= 4 {
			priority = models.PriorityLow
		}

		suggestion := models.NewSuggestion(
			models.CategoryAuthenticity,
			priority,
			gs.Title,
			gs.Description,
			"", // 不关联特定规则
		)

		// 如果有原文和建议文本，添加为示例
		if gs.OriginalText != "" && gs.SuggestedText != "" {
			suggestion.AddExample(gs.OriginalText, gs.SuggestedText, gs.Reason)
		}

		suggestions = append(suggestions, suggestion)
	}
	return suggestions
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 使用时间戳（到微秒）+ 随机数确保唯一性
	now := time.Now()
	timestamp := now.Format("20060102150405")
	microsecond := now.Nanosecond() / 1000
	return fmt.Sprintf("%s%06d", timestamp, microsecond)
}

package scorer

import (
	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
)

// Calculator 评分计算器
type Calculator struct {
	config *config.Config
}

// NewCalculator 创建评分计算器
func NewCalculator(cfg *config.Config) *Calculator {
	return &Calculator{
		config: cfg,
	}
}

// Calculate 计算完整评分
func (c *Calculator) Calculate(results []models.RuleResult) models.Score {
	// 计算维度评分
	dimensions := c.CalculateDimensions(results)

	// 计算总分（包含红旗规则）
	total := c.CalculateTotalWithRedFlags(dimensions, results)

	// 构建详细评分分解
	breakdown := make(map[string]float64)
	for _, result := range results {
		breakdown[string(result.RuleType)] = result.Score
	}

	return models.Score{
		Total:      total,
		Dimensions: dimensions,
		Breakdown:  breakdown,
	}
}

// CalculateDimensions 计算5维度评分
func (c *Calculator) CalculateDimensions(results []models.RuleResult) models.DimensionScores {
	weights := c.config.Scoring.Weights

	// 词汇多样性 (20分) - 基于Signal 1, 2
	vocabularyScore := c.calculateVocabularyDiversity(results, weights.VocabularyDiversity)

	// 句式复杂度 (15分) - 基于Signal 2, 5
	sentenceScore := c.calculateSentenceComplexity(results, weights.SentenceComplexity)

	// 个人化表达 (25分) - 基于Signal 10
	personalizationScore := c.calculatePersonalization(results, weights.Personalization)

	// 逻辑连贯性 (20分) - 基于Signal 3
	coherenceScore := c.calculateLogicalCoherence(results, weights.LogicalCoherence)

	// 情感真实度 (20分) - 基于Signal 9, 10
	authenticityScore := c.calculateEmotionalAuthenticity(results, weights.EmotionalAuthenticity)

	return models.DimensionScores{
		VocabularyDiversity:   vocabularyScore,
		SentenceComplexity:    sentenceScore,
		Personalization:       personalizationScore,
		LogicalCoherence:      coherenceScore,
		EmotionalAuthenticity: authenticityScore,
	}
}

// CalculateTotal 计算总分
func (c *Calculator) CalculateTotal(dimensions models.DimensionScores) float64 {
	total := dimensions.VocabularyDiversity.Score +
		dimensions.SentenceComplexity.Score +
		dimensions.Personalization.Score +
		dimensions.LogicalCoherence.Score +
		dimensions.EmotionalAuthenticity.Score

	// 确保总分在0-100范围内
	if total < 0 {
		total = 0
	}
	if total > 100 {
		total = 100
	}

	return total
}

// CalculateTotalWithRedFlags 计算总分（考虑红旗规则）
func (c *Calculator) CalculateTotalWithRedFlags(dimensions models.DimensionScores, results []models.RuleResult) float64 {
	total := c.CalculateTotal(dimensions)

	// 检查关键"红旗"规则
	// Signal 4: 引用异常（任何检测到都是严重问题）
	citationResult := c.findRuleResult(results, models.RuleTypeCitationAnomaly)
	if citationResult != nil && citationResult.Detected {
		// 引用异常是明确的AI标记，直接扣50分
		total -= 50
	}

	// Signal 8: 知识截止日期（任何检测到都是严重问题）
	knowledgeCutoffResult := c.findRuleResult(results, models.RuleTypeKnowledgeCutoff)
	if knowledgeCutoffResult != nil && knowledgeCutoffResult.Detected {
		// 知识截止日期是明确的AI标记，直接扣50分
		total -= 50
	}

	// Signal 6: Markdown残留
	markdownResult := c.findRuleResult(results, models.RuleTypeMarkdown)
	if markdownResult != nil && markdownResult.Detected {
		deduction := (100.0 - markdownResult.Score) / 100.0 * 15
		total -= deduction
	}

	// 确保总分在0-100范围内
	if total < 0 {
		total = 0
	}
	if total > 100 {
		total = 100
	}

	return total
}

// calculateVocabularyDiversity 计算词汇多样性评分
func (c *Calculator) calculateVocabularyDiversity(results []models.RuleResult, maxScore float64) models.DimensionScore {
	var issues []string
	var totalScore float64 = maxScore

	// Signal 1: 高频词汇检测
	highFreqResult := c.findRuleResult(results, models.RuleTypeHighFreqWords)
	if highFreqResult != nil {
		if highFreqResult.Detected {
			// 根据规则评分计算扣分
			deduction := (100.0 - highFreqResult.Score) / 100.0 * maxScore * 0.6
			totalScore -= deduction
			issues = append(issues, "检测到过度使用AI常用高频词汇")
		}
	}

	// Signal 2: 句式开头检测
	sentenceStarterResult := c.findRuleResult(results, models.RuleTypeSentenceStarters)
	if sentenceStarterResult != nil {
		if sentenceStarterResult.Detected {
			deduction := (100.0 - sentenceStarterResult.Score) / 100.0 * maxScore * 0.4
			totalScore -= deduction
			issues = append(issues, "检测到重复的句式开头模式")
		}
	}

	if totalScore < 0 {
		totalScore = 0
	}

	description := "词汇使用的丰富程度和多样性"
	if len(issues) == 0 {
		description = "词汇多样性良好，未检测到明显的AI生成模式"
	}

	return models.NewDimensionScore(totalScore, maxScore, issues, description)
}

// calculateSentenceComplexity 计算句式复杂度评分
func (c *Calculator) calculateSentenceComplexity(results []models.RuleResult, maxScore float64) models.DimensionScore {
	var issues []string
	var totalScore float64 = maxScore

	// Signal 2: 句式开头检测
	sentenceStarterResult := c.findRuleResult(results, models.RuleTypeSentenceStarters)
	if sentenceStarterResult != nil {
		if sentenceStarterResult.Detected {
			deduction := (100.0 - sentenceStarterResult.Score) / 100.0 * maxScore * 0.5
			totalScore -= deduction
			issues = append(issues, "句式结构单一，缺乏变化")
		}
	}

	// Signal 5: 破折号密度
	emDashResult := c.findRuleResult(results, models.RuleTypeEmDash)
	if emDashResult != nil {
		if emDashResult.Detected {
			deduction := (100.0 - emDashResult.Score) / 100.0 * maxScore * 0.5
			totalScore -= deduction
			issues = append(issues, "破折号使用过度，影响句式自然性")
		}
	}

	if totalScore < 0 {
		totalScore = 0
	}

	description := "句式结构的多样性和复杂度"
	if len(issues) == 0 {
		description = "句式复杂度良好，表现自然"
	}

	return models.NewDimensionScore(totalScore, maxScore, issues, description)
}

// calculatePersonalization 计算个人化表达评分
func (c *Calculator) calculatePersonalization(results []models.RuleResult, maxScore float64) models.DimensionScore {
	var issues []string
	var totalScore float64 = maxScore

	// Signal 10: 完美主义陷阱（缺乏个人化表达）
	perfectionismResult := c.findRuleResult(results, models.RuleTypePerfectionism)
	if perfectionismResult != nil {
		if perfectionismResult.Detected {
			// 完美主义检测到意味着缺乏个人化表达
			deduction := (100.0 - perfectionismResult.Score) / 100.0 * maxScore
			totalScore -= deduction
			issues = append(issues, "缺乏第一人称、情感词汇和不确定性表达")
		}
	}

	if totalScore < 0 {
		totalScore = 0
	}

	description := "个人风格和主观表达的程度"
	if len(issues) == 0 {
		description = "个人化表达充分，具有人类写作特征"
	}

	return models.NewDimensionScore(totalScore, maxScore, issues, description)
}

// calculateLogicalCoherence 计算逻辑连贯性评分
func (c *Calculator) calculateLogicalCoherence(results []models.RuleResult, maxScore float64) models.DimensionScore {
	var issues []string
	var totalScore float64 = maxScore

	// Signal 3: 虚假范围表达
	falseRangeResult := c.findRuleResult(results, models.RuleTypeFalseRange)
	if falseRangeResult != nil {
		if falseRangeResult.Detected {
			deduction := (100.0 - falseRangeResult.Score) / 100.0 * maxScore
			totalScore -= deduction
			issues = append(issues, "存在逻辑不连贯的范围表达")
		}
	}

	if totalScore < 0 {
		totalScore = 0
	}

	description := "逻辑结构的自然性和连贯性"
	if len(issues) == 0 {
		description = "逻辑连贯性良好，表达自然流畅"
	}

	return models.NewDimensionScore(totalScore, maxScore, issues, description)
}

// calculateEmotionalAuthenticity 计算情感真实度评分
func (c *Calculator) calculateEmotionalAuthenticity(results []models.RuleResult, maxScore float64) models.DimensionScore {
	var issues []string
	var totalScore float64 = maxScore

	// Signal 9: 协作式语气
	collaborativeResult := c.findRuleResult(results, models.RuleTypeCollaborative)
	if collaborativeResult != nil {
		if collaborativeResult.Detected {
			deduction := (100.0 - collaborativeResult.Score) / 100.0 * maxScore * 0.5
			totalScore -= deduction
			issues = append(issues, "使用AI助手特有的协作式语气")
		}
	}

	// Signal 10: 完美主义陷阱
	perfectionismResult := c.findRuleResult(results, models.RuleTypePerfectionism)
	if perfectionismResult != nil {
		if perfectionismResult.Detected {
			deduction := (100.0 - perfectionismResult.Score) / 100.0 * maxScore * 0.5
			totalScore -= deduction
			issues = append(issues, "情感表达不足，过于客观完美")
		}
	}

	if totalScore < 0 {
		totalScore = 0
	}

	description := "情感表达的真实性和自然性"
	if len(issues) == 0 {
		description = "情感表达真实自然，具有人类特征"
	}

	return models.NewDimensionScore(totalScore, maxScore, issues, description)
}

// findRuleResult 查找指定类型的规则结果
func (c *Calculator) findRuleResult(results []models.RuleResult, ruleType models.RuleType) *models.RuleResult {
	for i := range results {
		if results[i].RuleType == ruleType {
			return &results[i]
		}
	}
	return nil
}

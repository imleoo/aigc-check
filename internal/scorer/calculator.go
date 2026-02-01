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

// RedFlagPenalty 红旗规则惩罚因子配置
type RedFlagPenalty struct {
	// 乘法因子：分数 * factor，factor < 1 表示扣分
	Factor float64
	// 最大扣分比例（基于当前分数）
	MaxDeductionRate float64
}

// 红旗规则惩罚配置
var redFlagPenalties = map[models.RuleType]RedFlagPenalty{
	// Signal 4: 引用异常 - 严重问题，扣分50%
	models.RuleTypeCitationAnomaly: {Factor: 0.50, MaxDeductionRate: 0.50},
	// Signal 8: 知识截止日期 - 严重问题，扣分50%
	models.RuleTypeKnowledgeCutoff: {Factor: 0.50, MaxDeductionRate: 0.50},
	// Signal 6: Markdown残留 - 中等问题，最大扣分20%
	models.RuleTypeMarkdown: {Factor: 0.80, MaxDeductionRate: 0.20},
	// Signal 7: 表情符号异常 - 轻微问题，最大扣分10%
	models.RuleTypeEmoji: {Factor: 0.90, MaxDeductionRate: 0.10},
}

// 最低保底分数：即使检测到严重问题，也保留最低分数以区分程度
const MinimumScore = 5.0

// CalculateTotalWithRedFlags 计算总分（考虑红旗规则）
// 使用乘法递减机制替代减法扣分，避免分数骤降到0
func (c *Calculator) CalculateTotalWithRedFlags(dimensions models.DimensionScores, results []models.RuleResult) float64 {
	total := c.CalculateTotal(dimensions)

	// 使用乘法因子扣分，避免分数变成负数
	// 例如：基础分100，两个红旗各扣50% → 100 * 0.5 * 0.5 = 25，而非100 - 50 - 50 = 0

	// Signal 4: 引用异常（任何检测到都是严重问题）
	citationResult := c.findRuleResult(results, models.RuleTypeCitationAnomaly)
	if citationResult != nil && citationResult.Detected {
		penalty := redFlagPenalties[models.RuleTypeCitationAnomaly]
		// 根据检测严重程度调整惩罚因子
		adjustedFactor := c.adjustFactorBySeverity(penalty.Factor, citationResult.Score)
		total *= adjustedFactor
	}

	// Signal 8: 知识截止日期（任何检测到都是严重问题）
	knowledgeCutoffResult := c.findRuleResult(results, models.RuleTypeKnowledgeCutoff)
	if knowledgeCutoffResult != nil && knowledgeCutoffResult.Detected {
		penalty := redFlagPenalties[models.RuleTypeKnowledgeCutoff]
		adjustedFactor := c.adjustFactorBySeverity(penalty.Factor, knowledgeCutoffResult.Score)
		total *= adjustedFactor
	}

	// Signal 6: Markdown残留（中等问题，按严重程度递减）
	markdownResult := c.findRuleResult(results, models.RuleTypeMarkdown)
	if markdownResult != nil && markdownResult.Detected {
		penalty := redFlagPenalties[models.RuleTypeMarkdown]
		// Markdown的惩罚程度基于检测到的问题严重程度
		severityRate := (100.0 - markdownResult.Score) / 100.0
		// 惩罚因子范围: [penalty.Factor, 1.0]，问题越严重，惩罚越重
		factor := 1.0 - (severityRate * penalty.MaxDeductionRate)
		total *= factor
	}

	// Signal 7: 表情符号异常（轻微问题）
	emojiResult := c.findRuleResult(results, models.RuleTypeEmoji)
	if emojiResult != nil && emojiResult.Detected {
		penalty := redFlagPenalties[models.RuleTypeEmoji]
		severityRate := (100.0 - emojiResult.Score) / 100.0
		factor := 1.0 - (severityRate * penalty.MaxDeductionRate)
		total *= factor
	}

	// 应用最低保底分数
	// 即使检测到所有红旗，也保留最低分数以区分严重程度
	if total < MinimumScore {
		total = MinimumScore
	}
	if total > 100 {
		total = 100
	}

	return total
}

// adjustFactorBySeverity 根据检测严重程度调整惩罚因子
// 检测分数越低，问题越严重，惩罚因子越小
func (c *Calculator) adjustFactorBySeverity(baseFactor float64, detectionScore float64) float64 {
	// detectionScore: 0-100，0表示严重问题，100表示轻微问题
	// 严重程度: 0-1，1表示最严重
	severity := (100.0 - detectionScore) / 100.0

	// 基于严重程度调整因子
	// 轻微问题：factor接近1（扣分少）
	// 严重问题：factor接近baseFactor（扣分多）
	// 公式: factor = 1 - (1 - baseFactor) * severity
	// 当severity=1（最严重）时，factor=baseFactor
	// 当severity=0（最轻微）时，factor=1
	factor := 1.0 - (1.0-baseFactor)*severity

	// 确保因子不低于baseFactor
	if factor < baseFactor {
		factor = baseFactor
	}

	return factor
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

package scorer

import "github.com/leoobai/aigc-check/internal/models"

// Scorer 评分器接口
type Scorer interface {
	// Calculate 计算评分
	Calculate(results []models.RuleResult) models.Score

	// CalculateDimensions 计算维度评分
	CalculateDimensions(results []models.RuleResult) models.DimensionScores

	// CalculateTotal 计算总分
	CalculateTotal(dimensions models.DimensionScores) float64
}

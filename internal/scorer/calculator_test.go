package scorer

import (
	"testing"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
)

func TestNewCalculator(t *testing.T) {
	cfg := &config.Config{
		Scoring: config.ScoringConfig{
			Weights: models.DefaultDimensionWeights,
		},
	}

	calc := NewCalculator(cfg)
	if calc == nil {
		t.Error("NewCalculator() returned nil")
	}
}

func TestCalculator_Calculate(t *testing.T) {
	cfg := &config.Config{
		Scoring: config.ScoringConfig{
			Weights: models.DefaultDimensionWeights,
		},
	}
	calc := NewCalculator(cfg)

	// 测试空结果
	t.Run("空规则结果", func(t *testing.T) {
		results := []models.RuleResult{}
		score := calc.Calculate(results)

		// 空结果应该得满分（100分）
		if score.Total < 0 || score.Total > 100 {
			t.Errorf("Total score = %.1f, want between 0 and 100", score.Total)
		}
	})

	// 测试正常文本（无问题）
	t.Run("正常文本无问题", func(t *testing.T) {
		results := []models.RuleResult{
			{
				RuleType: models.RuleTypeHighFreqWords,
				Detected: false,
				Score:    100.0,
			},
			{
				RuleType: models.RuleTypeSentenceStarters,
				Detected: false,
				Score:    100.0,
			},
		}
		score := calc.Calculate(results)

		if score.Total < 80 {
			t.Errorf("Total score = %.1f, expected high score for normal text", score.Total)
		}
	})

	// 测试有问题的文本
	t.Run("有问题的文本", func(t *testing.T) {
		results := []models.RuleResult{
			{
				RuleType: models.RuleTypeHighFreqWords,
				Detected: true,
				Score:    50.0,
			},
			{
				RuleType: models.RuleTypeSentenceStarters,
				Detected: true,
				Score:    60.0,
			},
			{
				RuleType: models.RuleTypePerfectionism,
				Detected: true,
				Score:    40.0,
			},
		}
		score := calc.Calculate(results)

		if score.Total > 80 {
			t.Errorf("Total score = %.1f, expected lower score for problematic text", score.Total)
		}
	})
}

func TestCalculator_CalculateTotalWithRedFlags(t *testing.T) {
	cfg := &config.Config{
		Scoring: config.ScoringConfig{
			Weights: models.DefaultDimensionWeights,
		},
	}
	calc := NewCalculator(cfg)

	// 测试知识截止日期红旗
	t.Run("知识截止日期红旗", func(t *testing.T) {
		dimensions := models.DimensionScores{
			VocabularyDiversity:   models.DimensionScore{Score: 20},
			SentenceComplexity:    models.DimensionScore{Score: 15},
			Personalization:       models.DimensionScore{Score: 25},
			LogicalCoherence:      models.DimensionScore{Score: 20},
			EmotionalAuthenticity: models.DimensionScore{Score: 20},
		}

		results := []models.RuleResult{
			{
				RuleType: models.RuleTypeKnowledgeCutoff,
				Detected: true,
				Score:    0.0,
			},
		}

		total := calc.CalculateTotalWithRedFlags(dimensions, results)

		// 知识截止日期应该导致50分扣除
		if total > 55 {
			t.Errorf("Total = %.1f, expected significant deduction for knowledge cutoff", total)
		}
	})

	// 测试引用异常红旗
	t.Run("引用异常红旗", func(t *testing.T) {
		dimensions := models.DimensionScores{
			VocabularyDiversity:   models.DimensionScore{Score: 20},
			SentenceComplexity:    models.DimensionScore{Score: 15},
			Personalization:       models.DimensionScore{Score: 25},
			LogicalCoherence:      models.DimensionScore{Score: 20},
			EmotionalAuthenticity: models.DimensionScore{Score: 20},
		}

		results := []models.RuleResult{
			{
				RuleType: models.RuleTypeCitationAnomaly,
				Detected: true,
				Score:    0.0,
			},
		}

		total := calc.CalculateTotalWithRedFlags(dimensions, results)

		// 引用异常应该导致50分扣除
		if total > 55 {
			t.Errorf("Total = %.1f, expected significant deduction for citation anomaly", total)
		}
	})

	// 测试Markdown残留
	t.Run("Markdown残留", func(t *testing.T) {
		dimensions := models.DimensionScores{
			VocabularyDiversity:   models.DimensionScore{Score: 20},
			SentenceComplexity:    models.DimensionScore{Score: 15},
			Personalization:       models.DimensionScore{Score: 25},
			LogicalCoherence:      models.DimensionScore{Score: 20},
			EmotionalAuthenticity: models.DimensionScore{Score: 20},
		}

		results := []models.RuleResult{
			{
				RuleType: models.RuleTypeMarkdown,
				Detected: true,
				Score:    50.0, // 50%的问题
			},
		}

		total := calc.CalculateTotalWithRedFlags(dimensions, results)

		// Markdown应该有一定扣分但不像红旗那么严重
		if total > 95 || total < 85 {
			t.Errorf("Total = %.1f, expected moderate deduction for markdown", total)
		}
	})

	// 测试双重红旗 - 使用乘法递减机制
	t.Run("双重红旗", func(t *testing.T) {
		dimensions := models.DimensionScores{
			VocabularyDiversity:   models.DimensionScore{Score: 20},
			SentenceComplexity:    models.DimensionScore{Score: 15},
			Personalization:       models.DimensionScore{Score: 25},
			LogicalCoherence:      models.DimensionScore{Score: 20},
			EmotionalAuthenticity: models.DimensionScore{Score: 20},
		}

		results := []models.RuleResult{
			{
				RuleType: models.RuleTypeKnowledgeCutoff,
				Detected: true,
				Score:    0.0, // 最严重
			},
			{
				RuleType: models.RuleTypeCitationAnomaly,
				Detected: true,
				Score:    0.0, // 最严重
			},
		}

		total := calc.CalculateTotalWithRedFlags(dimensions, results)

		// 使用乘法递减: 100 * 0.5 * 0.5 = 25
		// 而非之前的减法: 100 - 50 - 50 = 0
		// 期望分数在20-30之间（有最低保底分数5分）
		if total < 20 || total > 35 {
			t.Errorf("Total = %.1f, expected between 20-35 for double red flags with multiplicative penalty", total)
		}
	})

	// 测试三重红旗（包含Markdown）
	t.Run("三重红旗", func(t *testing.T) {
		dimensions := models.DimensionScores{
			VocabularyDiversity:   models.DimensionScore{Score: 20},
			SentenceComplexity:    models.DimensionScore{Score: 15},
			Personalization:       models.DimensionScore{Score: 25},
			LogicalCoherence:      models.DimensionScore{Score: 20},
			EmotionalAuthenticity: models.DimensionScore{Score: 20},
		}

		results := []models.RuleResult{
			{
				RuleType: models.RuleTypeKnowledgeCutoff,
				Detected: true,
				Score:    0.0,
			},
			{
				RuleType: models.RuleTypeCitationAnomaly,
				Detected: true,
				Score:    0.0,
			},
			{
				RuleType: models.RuleTypeMarkdown,
				Detected: true,
				Score:    20.0, // 严重的Markdown问题
			},
		}

		total := calc.CalculateTotalWithRedFlags(dimensions, results)

		// 100 * 0.5 * 0.5 * 0.84 ≈ 21，应该大于最低保底分数
		if total < MinimumScore {
			t.Errorf("Total = %.1f, should be at least minimum score %.1f", total, MinimumScore)
		}
		if total > 25 {
			t.Errorf("Total = %.1f, expected lower score for triple red flags", total)
		}
	})
}

func TestCalculator_CalculateTotal(t *testing.T) {
	cfg := &config.Config{
		Scoring: config.ScoringConfig{
			Weights: models.DefaultDimensionWeights,
		},
	}
	calc := NewCalculator(cfg)

	tests := []struct {
		name       string
		dimensions models.DimensionScores
		minTotal   float64
		maxTotal   float64
	}{
		{
			name: "满分",
			dimensions: models.DimensionScores{
				VocabularyDiversity:   models.DimensionScore{Score: 20},
				SentenceComplexity:    models.DimensionScore{Score: 15},
				Personalization:       models.DimensionScore{Score: 25},
				LogicalCoherence:      models.DimensionScore{Score: 20},
				EmotionalAuthenticity: models.DimensionScore{Score: 20},
			},
			minTotal: 100,
			maxTotal: 100,
		},
		{
			name: "半分",
			dimensions: models.DimensionScores{
				VocabularyDiversity:   models.DimensionScore{Score: 10},
				SentenceComplexity:    models.DimensionScore{Score: 7.5},
				Personalization:       models.DimensionScore{Score: 12.5},
				LogicalCoherence:      models.DimensionScore{Score: 10},
				EmotionalAuthenticity: models.DimensionScore{Score: 10},
			},
			minTotal: 50,
			maxTotal: 50,
		},
		{
			name: "零分",
			dimensions: models.DimensionScores{
				VocabularyDiversity:   models.DimensionScore{Score: 0},
				SentenceComplexity:    models.DimensionScore{Score: 0},
				Personalization:       models.DimensionScore{Score: 0},
				LogicalCoherence:      models.DimensionScore{Score: 0},
				EmotionalAuthenticity: models.DimensionScore{Score: 0},
			},
			minTotal: 0,
			maxTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total := calc.CalculateTotal(tt.dimensions)

			if total < tt.minTotal || total > tt.maxTotal {
				t.Errorf("Total = %.1f, want between %.1f and %.1f", total, tt.minTotal, tt.maxTotal)
			}
		})
	}
}

func TestCalculator_CalculateDimensions(t *testing.T) {
	cfg := &config.Config{
		Scoring: config.ScoringConfig{
			Weights: models.DefaultDimensionWeights,
		},
	}
	calc := NewCalculator(cfg)

	// 测试无问题的维度计算
	t.Run("无问题维度", func(t *testing.T) {
		results := []models.RuleResult{
			{RuleType: models.RuleTypeHighFreqWords, Detected: false, Score: 100},
			{RuleType: models.RuleTypeSentenceStarters, Detected: false, Score: 100},
			{RuleType: models.RuleTypeFalseRange, Detected: false, Score: 100},
			{RuleType: models.RuleTypePerfectionism, Detected: false, Score: 100},
			{RuleType: models.RuleTypeCollaborative, Detected: false, Score: 100},
		}

		dimensions := calc.CalculateDimensions(results)

		// 所有维度应该是满分或接近满分
		if dimensions.VocabularyDiversity.Score < 18 {
			t.Errorf("VocabularyDiversity = %.1f, expected high score", dimensions.VocabularyDiversity.Score)
		}
		if dimensions.Personalization.Score < 23 {
			t.Errorf("Personalization = %.1f, expected high score", dimensions.Personalization.Score)
		}
	})

	// 测试有问题的维度计算
	t.Run("有问题维度", func(t *testing.T) {
		results := []models.RuleResult{
			{RuleType: models.RuleTypeHighFreqWords, Detected: true, Score: 30},
			{RuleType: models.RuleTypeSentenceStarters, Detected: true, Score: 40},
			{RuleType: models.RuleTypePerfectionism, Detected: true, Score: 20},
			{RuleType: models.RuleTypeCollaborative, Detected: true, Score: 30},
		}

		dimensions := calc.CalculateDimensions(results)

		// 词汇多样性应该较低
		if dimensions.VocabularyDiversity.Score > 15 {
			t.Errorf("VocabularyDiversity = %.1f, expected lower score for detected issues", dimensions.VocabularyDiversity.Score)
		}

		// 个人化表达应该较低
		if dimensions.Personalization.Score > 15 {
			t.Errorf("Personalization = %.1f, expected lower score for detected issues", dimensions.Personalization.Score)
		}
	})
}

func TestCalculator_ScoreBounds(t *testing.T) {
	cfg := &config.Config{
		Scoring: config.ScoringConfig{
			Weights: models.DefaultDimensionWeights,
		},
	}
	calc := NewCalculator(cfg)

	// 测试分数边界
	t.Run("分数不应超过100", func(t *testing.T) {
		dimensions := models.DimensionScores{
			VocabularyDiversity:   models.DimensionScore{Score: 50},
			SentenceComplexity:    models.DimensionScore{Score: 50},
			Personalization:       models.DimensionScore{Score: 50},
			LogicalCoherence:      models.DimensionScore{Score: 50},
			EmotionalAuthenticity: models.DimensionScore{Score: 50},
		}

		total := calc.CalculateTotal(dimensions)

		if total > 100 {
			t.Errorf("Total = %.1f, should not exceed 100", total)
		}
	})

	t.Run("分数不应低于0", func(t *testing.T) {
		dimensions := models.DimensionScores{
			VocabularyDiversity:   models.DimensionScore{Score: -10},
			SentenceComplexity:    models.DimensionScore{Score: -10},
			Personalization:       models.DimensionScore{Score: -10},
			LogicalCoherence:      models.DimensionScore{Score: -10},
			EmotionalAuthenticity: models.DimensionScore{Score: -10},
		}

		total := calc.CalculateTotal(dimensions)

		if total < 0 {
			t.Errorf("Total = %.1f, should not be below 0", total)
		}
	})
}

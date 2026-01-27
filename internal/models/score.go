package models

// Score 评分结果
type Score struct {
	Total      float64           `json:"total"`       // 总分 (0-100)
	Dimensions DimensionScores   `json:"dimensions"`  // 维度评分
	Breakdown  map[string]float64 `json:"breakdown"`  // 详细评分分解
}

// DimensionScores 5维度评分
type DimensionScores struct {
	VocabularyDiversity    DimensionScore `json:"vocabulary_diversity"`     // 词汇多样性 (20分)
	SentenceComplexity     DimensionScore `json:"sentence_complexity"`      // 句式复杂度 (15分)
	Personalization        DimensionScore `json:"personalization"`          // 个人化表达 (25分)
	LogicalCoherence       DimensionScore `json:"logical_coherence"`        // 逻辑连贯性 (20分)
	EmotionalAuthenticity  DimensionScore `json:"emotional_authenticity"`   // 情感真实度 (20分)
}

// DimensionScore 单个维度评分
type DimensionScore struct {
	Score       float64  `json:"score"`        // 得分
	MaxScore    float64  `json:"max_score"`    // 满分
	Percentage  float64  `json:"percentage"`   // 百分比
	Level       string   `json:"level"`        // 等级 (优秀/良好/一般/较差)
	Description string   `json:"description"`  // 描述
	Issues      []string `json:"issues"`       // 发现的问题
}

// DimensionWeights 维度权重配置
type DimensionWeights struct {
	VocabularyDiversity   float64 `yaml:"vocabulary_diversity"`    // 词汇多样性权重
	SentenceComplexity    float64 `yaml:"sentence_complexity"`     // 句式复杂度权重
	Personalization       float64 `yaml:"personalization"`         // 个人化表达权重
	LogicalCoherence      float64 `yaml:"logical_coherence"`       // 逻辑连贯性权重
	EmotionalAuthenticity float64 `yaml:"emotional_authenticity"`  // 情感真实度权重
}

// DefaultDimensionWeights 默认维度权重
var DefaultDimensionWeights = DimensionWeights{
	VocabularyDiversity:   20.0,
	SentenceComplexity:    15.0,
	Personalization:       25.0,
	LogicalCoherence:      20.0,
	EmotionalAuthenticity: 20.0,
}

// GetLevel 根据百分比获取等级
func GetLevel(percentage float64) string {
	switch {
	case percentage >= 90:
		return "优秀"
	case percentage >= 75:
		return "良好"
	case percentage >= 60:
		return "一般"
	default:
		return "较差"
	}
}

// CalculatePercentage 计算百分比
func CalculatePercentage(score, maxScore float64) float64 {
	if maxScore == 0 {
		return 0
	}
	return (score / maxScore) * 100
}

// NewDimensionScore 创建维度评分
func NewDimensionScore(score, maxScore float64, issues []string, description string) DimensionScore {
	percentage := CalculatePercentage(score, maxScore)
	return DimensionScore{
		Score:       score,
		MaxScore:    maxScore,
		Percentage:  percentage,
		Level:       GetLevel(percentage),
		Description: description,
		Issues:      issues,
	}
}

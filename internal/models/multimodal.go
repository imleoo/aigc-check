package models

// MultimodalResult 多模态检测结果
type MultimodalResult struct {
	// Layer 1: 规则检测分数 (0-100)
	RuleLayerScore float64 `json:"rule_layer_score"`

	// Layer 2: 统计分析分数 (0-100)
	StatisticsLayerScore float64 `json:"statistics_layer_score"`

	// Layer 3: 语义分析分数 (0-100, 仅在启用时有效)
	SemanticLayerScore float64 `json:"semantic_layer_score"`

	// 最终融合分数 (0-100)
	FinalScore float64 `json:"final_score"`

	// 置信度 (0-1)
	Confidence float64 `json:"confidence"`

	// 各层权重
	LayerWeights LayerWeights `json:"layer_weights"`

	// 检测模式
	DetectionMode DetectionMode `json:"detection_mode"`

	// 各层详细结果
	RuleLayerDetails       *RuleLayerDetails       `json:"rule_layer_details,omitempty"`
	StatisticsLayerDetails *StatisticsLayerDetails `json:"statistics_layer_details,omitempty"`
	SemanticLayerDetails   *SemanticLayerDetails   `json:"semantic_layer_details,omitempty"`

	// 融合说明
	FusionExplanation string `json:"fusion_explanation"`
}

// LayerWeights 各层权重配置
type LayerWeights struct {
	// 规则层权重 (默认 0.4)
	RuleLayer float64 `json:"rule_layer" yaml:"rule_layer"`

	// 统计层权重 (默认 0.3)
	StatisticsLayer float64 `json:"statistics_layer" yaml:"statistics_layer"`

	// 语义层权重 (默认 0.3)
	SemanticLayer float64 `json:"semantic_layer" yaml:"semantic_layer"`
}

// DefaultLayerWeights 默认层权重
var DefaultLayerWeights = LayerWeights{
	RuleLayer:       0.4,
	StatisticsLayer: 0.3,
	SemanticLayer:   0.3,
}

// TwoLayerWeights 两层模式权重 (无语义分析)
var TwoLayerWeights = LayerWeights{
	RuleLayer:       0.55,
	StatisticsLayer: 0.45,
	SemanticLayer:   0.0,
}

// DetectionMode 检测模式
type DetectionMode string

const (
	// DetectionModeRuleOnly 仅规则检测
	DetectionModeRuleOnly DetectionMode = "rule_only"

	// DetectionModeRuleStatistics 规则+统计分析
	DetectionModeRuleStatistics DetectionMode = "rule_statistics"

	// DetectionModeMultimodal 完整多模态检测
	DetectionModeMultimodal DetectionMode = "multimodal"
)

// RuleLayerDetails 规则层详细结果
type RuleLayerDetails struct {
	// 检测到的规则数量
	DetectedRulesCount int `json:"detected_rules_count"`

	// 总规则数量
	TotalRulesCount int `json:"total_rules_count"`

	// 红旗规则数量
	RedFlagCount int `json:"red_flag_count"`

	// 主要问题
	MainIssues []string `json:"main_issues"`
}

// StatisticsLayerDetails 统计层详细结果
type StatisticsLayerDetails struct {
	// 词汇多样性指标
	TypeTokenRatio float64 `json:"type_token_ratio"`

	// 词汇丰富度
	VocabularyRichness float64 `json:"vocabulary_richness"`

	// 句子长度变化
	SentenceLengthVariance float64 `json:"sentence_length_variance"`

	// 句式复杂度
	SentenceComplexity float64 `json:"sentence_complexity"`

	// 困惑度估算
	PerplexityScore float64 `json:"perplexity_score"`

	// AI概率估算
	AIProbability float64 `json:"ai_probability"`

	// 分析详情
	Details []string `json:"details"`
}

// SemanticLayerDetails 语义层详细结果
type SemanticLayerDetails struct {
	// 逻辑连贯性分数
	CoherenceScore float64 `json:"coherence_score"`

	// 个人化程度分数
	PersonalizationScore float64 `json:"personalization_score"`

	// AI模式检测分数
	AIPatternScore float64 `json:"ai_pattern_score"`

	// 检测到的特征
	DetectedFeatures []string `json:"detected_features"`

	// 分析说明
	Explanation string `json:"explanation"`

	// 是否使用了缓存
	FromCache bool `json:"from_cache"`
}

// ConfidenceThresholds 置信度阈值
type ConfidenceThresholds struct {
	// 高置信度阈值 - 直接输出，不需要额外分析
	High float64 `json:"high" yaml:"high"`

	// 中置信度阈值 - 需要统计分析验证
	Medium float64 `json:"medium" yaml:"medium"`

	// 低置信度阈值 - 需要LLM深度分析
	Low float64 `json:"low" yaml:"low"`
}

// DefaultConfidenceThresholds 默认置信度阈值
var DefaultConfidenceThresholds = ConfidenceThresholds{
	High:   0.85,
	Medium: 0.60,
	Low:    0.40,
}

// MultimodalConfig 多模态配置
type MultimodalConfig struct {
	// 是否启用多模态检测
	Enabled bool `json:"enabled" yaml:"enabled"`

	// 是否启用统计分析层
	EnableStatistics bool `json:"enable_statistics" yaml:"enable_statistics"`

	// 是否启用语义分析层 (需要Gemini API)
	EnableSemantic bool `json:"enable_semantic" yaml:"enable_semantic"`

	// 层权重
	Weights LayerWeights `json:"weights" yaml:"weights"`

	// 置信度阈值
	ConfidenceThresholds ConfidenceThresholds `json:"confidence_thresholds" yaml:"confidence_thresholds"`

	// 分层触发策略
	TieredTrigger bool `json:"tiered_trigger" yaml:"tiered_trigger"`
}

// DefaultMultimodalConfig 默认多模态配置
var DefaultMultimodalConfig = MultimodalConfig{
	Enabled:              false,
	EnableStatistics:     true,
	EnableSemantic:       false,
	Weights:              DefaultLayerWeights,
	ConfidenceThresholds: DefaultConfidenceThresholds,
	TieredTrigger:        true,
}

// GetDetectionMode 根据配置获取检测模式
func (c *MultimodalConfig) GetDetectionMode() DetectionMode {
	if !c.Enabled {
		return DetectionModeRuleOnly
	}
	if c.EnableSemantic {
		return DetectionModeMultimodal
	}
	if c.EnableStatistics {
		return DetectionModeRuleStatistics
	}
	return DetectionModeRuleOnly
}

// GetEffectiveWeights 获取有效的层权重
func (c *MultimodalConfig) GetEffectiveWeights() LayerWeights {
	mode := c.GetDetectionMode()
	switch mode {
	case DetectionModeRuleOnly:
		return LayerWeights{RuleLayer: 1.0, StatisticsLayer: 0.0, SemanticLayer: 0.0}
	case DetectionModeRuleStatistics:
		return TwoLayerWeights
	case DetectionModeMultimodal:
		return c.Weights
	default:
		return LayerWeights{RuleLayer: 1.0, StatisticsLayer: 0.0, SemanticLayer: 0.0}
	}
}

// NeedsStatisticsAnalysis 判断是否需要统计分析
func NeedsStatisticsAnalysis(ruleConfidence float64, thresholds ConfidenceThresholds) bool {
	return ruleConfidence < thresholds.High
}

// NeedsSemanticAnalysis 判断是否需要语义分析
func NeedsSemanticAnalysis(ruleConfidence float64, statsConfidence float64, thresholds ConfidenceThresholds) bool {
	// 如果规则和统计分析的平均置信度低于中等阈值，则需要语义分析
	avgConfidence := (ruleConfidence + statsConfidence) / 2
	return avgConfidence < thresholds.Medium
}

// FuseScores 融合多层分数
func FuseScores(ruleScore, statsScore, semanticScore float64, weights LayerWeights) float64 {
	totalWeight := weights.RuleLayer + weights.StatisticsLayer + weights.SemanticLayer
	if totalWeight == 0 {
		return ruleScore
	}

	fusedScore := (ruleScore*weights.RuleLayer +
		statsScore*weights.StatisticsLayer +
		semanticScore*weights.SemanticLayer) / totalWeight

	// 确保分数在 0-100 范围内
	if fusedScore < 0 {
		return 0
	}
	if fusedScore > 100 {
		return 100
	}

	return fusedScore
}

// CalculateFusionConfidence 计算融合置信度
func CalculateFusionConfidence(ruleConf, statsConf, semanticConf float64, weights LayerWeights) float64 {
	totalWeight := weights.RuleLayer + weights.StatisticsLayer + weights.SemanticLayer
	if totalWeight == 0 {
		return ruleConf
	}

	fusedConf := (ruleConf*weights.RuleLayer +
		statsConf*weights.StatisticsLayer +
		semanticConf*weights.SemanticLayer) / totalWeight

	// 置信度在 0-1 范围内
	if fusedConf < 0 {
		return 0
	}
	if fusedConf > 1 {
		return 1
	}

	return fusedConf
}

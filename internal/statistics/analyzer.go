// Package statistics 提供文本统计分析功能
// 用于检测AI生成内容的统计特征
package statistics

import (
	"math"
)

// Analyzer 统计分析器
type Analyzer struct {
	vocabAnalyzer    *VocabularyAnalyzer
	sentenceAnalyzer *SentenceAnalyzer
	perplexityCalc   *PerplexityCalculator
}

// NewAnalyzer 创建统计分析器
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		vocabAnalyzer:    NewVocabularyAnalyzer(),
		sentenceAnalyzer: NewSentenceAnalyzer(),
		perplexityCalc:   NewPerplexityCalculator(),
	}
}

// StatisticsResult 统计分析结果
type StatisticsResult struct {
	// 词汇统计
	Vocabulary VocabularyStats `json:"vocabulary"`

	// 句子统计
	Sentence SentenceStats `json:"sentence"`

	// 困惑度
	Perplexity PerplexityResult `json:"perplexity"`

	// 综合评分 (0-100，越高越像人类写作)
	HumanScore float64 `json:"human_score"`

	// AI可能性 (0-1，越高越可能是AI)
	AIProbability float64 `json:"ai_probability"`

	// 置信度 (0-1)
	Confidence float64 `json:"confidence"`

	// 分析细节
	Details []string `json:"details"`
}

// Analyze 执行完整的统计分析
func (a *Analyzer) Analyze(text string) StatisticsResult {
	// 分析词汇
	vocabStats := a.vocabAnalyzer.Analyze(text)

	// 分析句子
	sentenceStats := a.sentenceAnalyzer.Analyze(text)

	// 计算困惑度
	perplexity := a.perplexityCalc.Calculate(text)

	// 计算综合评分
	humanScore, aiProb, confidence, details := a.calculateCompositeScore(vocabStats, sentenceStats, perplexity)

	return StatisticsResult{
		Vocabulary:    vocabStats,
		Sentence:      sentenceStats,
		Perplexity:    perplexity,
		HumanScore:    humanScore,
		AIProbability: aiProb,
		Confidence:    confidence,
		Details:       details,
	}
}

// calculateCompositeScore 计算综合评分
func (a *Analyzer) calculateCompositeScore(
	vocab VocabularyStats,
	sentence SentenceStats,
	perplexity PerplexityResult,
) (humanScore float64, aiProb float64, confidence float64, details []string) {

	var scores []float64
	var weights []float64

	// 词汇多样性评分
	// TTR > 0.5 通常表示人类写作
	// TTR < 0.3 可能表示AI生成
	vocabScore := calculateVocabScore(vocab.TTR)
	scores = append(scores, vocabScore)
	weights = append(weights, 0.3) // 权重30%

	if vocab.TTR < 0.35 {
		details = append(details, "词汇多样性较低 (TTR < 0.35)")
	} else if vocab.TTR > 0.55 {
		details = append(details, "词汇多样性良好 (TTR > 0.55)")
	}

	// 词汇丰富度评分
	richnessScore := calculateRichnessScore(vocab.Richness)
	scores = append(scores, richnessScore)
	weights = append(weights, 0.15) // 权重15%

	if vocab.Richness < 5 {
		details = append(details, "词汇丰富度较低")
	}

	// 句子变化评分
	// 方差越大，越像人类写作
	varianceScore := calculateVarianceScore(sentence.LengthStdDev)
	scores = append(scores, varianceScore)
	weights = append(weights, 0.2) // 权重20%

	if sentence.LengthStdDev < 5 {
		details = append(details, "句子长度变化较小，结构较为单一")
	} else if sentence.LengthStdDev > 15 {
		details = append(details, "句子长度变化丰富，结构多样")
	}

	// 句式复杂度评分
	complexityScore := calculateComplexityScore(sentence.ComplexityScore)
	scores = append(scores, complexityScore)
	weights = append(weights, 0.15) // 权重15%

	// 困惑度评分
	// 困惑度过低可能表示AI生成（过于"流畅"）
	perplexityScore := calculatePerplexityScore(perplexity.Score)
	scores = append(scores, perplexityScore)
	weights = append(weights, 0.2) // 权重20%

	if perplexity.Score < 30 {
		details = append(details, "文本困惑度较低，可能过于流畅")
	} else if perplexity.Score > 100 {
		details = append(details, "文本困惑度较高，表达自然")
	}

	// 计算加权平均分
	var totalWeight float64
	for i, score := range scores {
		humanScore += score * weights[i]
		totalWeight += weights[i]
	}
	humanScore /= totalWeight

	// 计算AI可能性
	aiProb = 1 - (humanScore / 100)

	// 计算置信度
	// 基于样本量和分数分布
	confidence = calculateConfidence(vocab.TotalWords, scores)

	// 四舍五入到两位小数
	humanScore = math.Round(humanScore*100) / 100
	aiProb = math.Round(aiProb*100) / 100
	confidence = math.Round(confidence*100) / 100

	return humanScore, aiProb, confidence, details
}

// calculateVocabScore 计算词汇多样性评分
func calculateVocabScore(ttr float64) float64 {
	// TTR范围通常在0.2-0.7之间
	// 0.2以下 -> 很可能AI (评分20)
	// 0.35 -> 中等 (评分50)
	// 0.5以上 -> 很可能人类 (评分80+)
	if ttr < 0.2 {
		return 20
	} else if ttr < 0.35 {
		return 20 + (ttr-0.2)/0.15*30 // 20-50
	} else if ttr < 0.5 {
		return 50 + (ttr-0.35)/0.15*30 // 50-80
	} else if ttr < 0.7 {
		return 80 + (ttr-0.5)/0.2*15 // 80-95
	}
	return 95
}

// calculateRichnessScore 计算词汇丰富度评分
func calculateRichnessScore(richness float64) float64 {
	// 词汇丰富度: Hapax Legomena 比例
	// 越高表示更多独特词汇
	if richness < 3 {
		return 30
	} else if richness < 5 {
		return 30 + (richness-3)/2*30 // 30-60
	} else if richness < 10 {
		return 60 + (richness-5)/5*25 // 60-85
	}
	return 85 + min((richness-10)/10*10, 10) // 85-95
}

// calculateVarianceScore 计算句子变化评分
func calculateVarianceScore(stdDev float64) float64 {
	// 标准差越大，句子长度变化越大
	if stdDev < 3 {
		return 30
	} else if stdDev < 8 {
		return 30 + (stdDev-3)/5*35 // 30-65
	} else if stdDev < 15 {
		return 65 + (stdDev-8)/7*25 // 65-90
	}
	return 90
}

// calculateComplexityScore 计算句式复杂度评分
func calculateComplexityScore(complexity float64) float64 {
	// 复杂度评分直接使用
	return complexity
}

// calculatePerplexityScore 计算困惑度评分
func calculatePerplexityScore(perplexity float64) float64 {
	// 困惑度范围通常在10-200之间
	// 过低（<30）可能是AI生成
	// 过高（>200）可能是混乱文本
	// 适中（50-150）通常是人类写作
	if perplexity < 20 {
		return 20
	} else if perplexity < 40 {
		return 20 + (perplexity-20)/20*30 // 20-50
	} else if perplexity < 100 {
		return 50 + (perplexity-40)/60*40 // 50-90
	} else if perplexity < 200 {
		return 90
	}
	return 80 // 过高的困惑度反而降低评分
}

// calculateConfidence 计算置信度
func calculateConfidence(wordCount int, scores []float64) float64 {
	// 基于样本量的置信度
	var sampleConfidence float64
	if wordCount < 100 {
		sampleConfidence = float64(wordCount) / 100 * 0.5
	} else if wordCount < 500 {
		sampleConfidence = 0.5 + float64(wordCount-100)/400*0.3
	} else {
		sampleConfidence = 0.8 + min(float64(wordCount-500)/1000*0.2, 0.2)
	}

	// 基于分数一致性的置信度
	if len(scores) == 0 {
		return sampleConfidence
	}

	var mean float64
	for _, s := range scores {
		mean += s
	}
	mean /= float64(len(scores))

	var variance float64
	for _, s := range scores {
		variance += (s - mean) * (s - mean)
	}
	variance /= float64(len(scores))

	// 方差越小，一致性越高，置信度越高
	consistencyConfidence := 1 - min(math.Sqrt(variance)/50, 0.5)

	// 综合置信度
	return sampleConfidence * 0.6 + consistencyConfidence * 0.4
}

// min returns the smaller of two float64 values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

package statistics

import (
	"math"
	"regexp"
	"strings"
)

// PerplexityCalculator 困惑度计算器
// 使用简化的n-gram模型估算文本困惑度
type PerplexityCalculator struct {
	// n-gram 大小
	ngramSize int
	// 平滑参数
	smoothingFactor float64
}

// PerplexityResult 困惑度计算结果
type PerplexityResult struct {
	// 困惑度分数
	// 范围通常在 10-200 之间
	// 较低（<30）可能表示AI生成（过于"流畅"）
	// 较高（>150）可能表示混乱或非结构化文本
	// 适中（50-100）通常是人类自然写作
	Score float64 `json:"score"`

	// 置信度 (0-1)
	Confidence float64 `json:"confidence"`

	// 分析的n-gram数量
	NGramCount int `json:"ngram_count"`

	// 唯一n-gram数量
	UniqueNGrams int `json:"unique_ngrams"`

	// 分析指标
	Indicators []string `json:"indicators"`
}

// NewPerplexityCalculator 创建困惑度计算器
func NewPerplexityCalculator() *PerplexityCalculator {
	return &PerplexityCalculator{
		ngramSize:       3, // 使用trigram
		smoothingFactor: 1.0,
	}
}

// Calculate 计算文本困惑度
func (p *PerplexityCalculator) Calculate(text string) PerplexityResult {
	// 预处理文本
	tokens := p.tokenize(text)

	if len(tokens) < p.ngramSize+1 {
		return PerplexityResult{
			Score:      50, // 文本太短，返回中等分数
			Confidence: 0.3,
			Indicators: []string{"文本过短，无法准确估算困惑度"},
		}
	}

	// 构建n-gram模型
	ngramCounts := p.buildNGramModel(tokens)
	nMinus1GramCounts := p.buildNGramModelWithSize(tokens, p.ngramSize-1)

	// 计算困惑度
	perplexity := p.calculatePerplexity(tokens, ngramCounts, nMinus1GramCounts)

	// 分析指标
	var indicators []string
	uniqueNGrams := len(ngramCounts)
	totalNGrams := len(tokens) - p.ngramSize + 1

	// 计算n-gram重复率
	repeatRate := 1 - float64(uniqueNGrams)/float64(totalNGrams)
	if repeatRate > 0.3 {
		indicators = append(indicators, "n-gram重复率较高，可能缺乏多样性")
	}

	// 判断困惑度水平
	if perplexity < 30 {
		indicators = append(indicators, "困惑度较低，文本可能过于规整/模板化")
	} else if perplexity > 150 {
		indicators = append(indicators, "困惑度较高，文本可能较为随机/非结构化")
	} else {
		indicators = append(indicators, "困惑度适中，符合自然写作特征")
	}

	// 计算置信度
	confidence := p.calculateConfidence(len(tokens), uniqueNGrams)

	return PerplexityResult{
		Score:        math.Round(perplexity*100) / 100,
		Confidence:   math.Round(confidence*100) / 100,
		NGramCount:   totalNGrams,
		UniqueNGrams: uniqueNGrams,
		Indicators:   indicators,
	}
}

// tokenize 分词
func (p *PerplexityCalculator) tokenize(text string) []string {
	// 转小写
	text = strings.ToLower(text)

	// 使用正则分词，保留标点
	tokenRegex := regexp.MustCompile(`[\p{L}\p{N}]+|[.!?。！？，,;；:：]`)
	tokens := tokenRegex.FindAllString(text, -1)

	// 添加句子边界标记
	var result []string
	result = append(result, "<s>") // 开始标记

	for _, token := range tokens {
		result = append(result, token)
		// 在句子结束符后添加边界标记
		if token == "." || token == "!" || token == "?" ||
			token == "。" || token == "！" || token == "？" {
			result = append(result, "</s>")
			result = append(result, "<s>")
		}
	}

	result = append(result, "</s>") // 结束标记

	return result
}

// buildNGramModel 构建n-gram模型
func (p *PerplexityCalculator) buildNGramModel(tokens []string) map[string]int {
	return p.buildNGramModelWithSize(tokens, p.ngramSize)
}

// buildNGramModelWithSize 构建指定大小的n-gram模型
func (p *PerplexityCalculator) buildNGramModelWithSize(tokens []string, size int) map[string]int {
	counts := make(map[string]int)

	for i := 0; i <= len(tokens)-size; i++ {
		ngram := strings.Join(tokens[i:i+size], " ")
		counts[ngram]++
	}

	return counts
}

// calculatePerplexity 计算困惑度
// 使用简化的估算方法：基于n-gram的条件概率
func (p *PerplexityCalculator) calculatePerplexity(
	tokens []string,
	ngramCounts map[string]int,
	nMinus1GramCounts map[string]int,
) float64 {
	if len(tokens) <= p.ngramSize {
		return 50 // 默认值
	}

	var logProbSum float64
	count := 0
	vocabSize := float64(len(nMinus1GramCounts)) // 词汇表大小

	for i := 0; i <= len(tokens)-p.ngramSize; i++ {
		ngram := strings.Join(tokens[i:i+p.ngramSize], " ")
		context := strings.Join(tokens[i:i+p.ngramSize-1], " ")

		ngramCount := float64(ngramCounts[ngram])
		contextCount := float64(nMinus1GramCounts[context])

		// 拉普拉斯平滑
		prob := (ngramCount + p.smoothingFactor) / (contextCount + p.smoothingFactor*vocabSize)

		if prob > 0 {
			logProbSum += math.Log2(prob)
			count++
		}
	}

	if count == 0 {
		return 50 // 默认值
	}

	// 计算困惑度: 2^(-1/N * sum(log2(P)))
	avgLogProb := logProbSum / float64(count)
	perplexity := math.Pow(2, -avgLogProb)

	// 限制范围，避免极端值
	if perplexity > 500 {
		perplexity = 500
	}
	if perplexity < 1 {
		perplexity = 1
	}

	return perplexity
}

// calculateConfidence 计算置信度
func (p *PerplexityCalculator) calculateConfidence(tokenCount int, uniqueNGrams int) float64 {
	// 基于样本量的置信度
	var sampleConfidence float64
	if tokenCount < 50 {
		sampleConfidence = float64(tokenCount) / 50 * 0.4
	} else if tokenCount < 200 {
		sampleConfidence = 0.4 + float64(tokenCount-50)/150*0.3
	} else if tokenCount < 500 {
		sampleConfidence = 0.7 + float64(tokenCount-200)/300*0.2
	} else {
		sampleConfidence = 0.9
	}

	// 基于n-gram多样性的置信度
	diversityRatio := float64(uniqueNGrams) / float64(tokenCount)
	diversityConfidence := min(diversityRatio*2, 1.0)

	// 综合置信度
	return sampleConfidence*0.7 + diversityConfidence*0.3
}

// CalculateEntropyRate 计算熵率（补充指标）
func (p *PerplexityCalculator) CalculateEntropyRate(text string) float64 {
	tokens := p.tokenize(text)

	if len(tokens) < 2 {
		return 0
	}

	// 计算一元词频
	unigramCounts := make(map[string]int)
	for _, token := range tokens {
		unigramCounts[token]++
	}

	// 计算熵
	totalTokens := float64(len(tokens))
	var entropy float64

	for _, count := range unigramCounts {
		prob := float64(count) / totalTokens
		if prob > 0 {
			entropy -= prob * math.Log2(prob)
		}
	}

	return math.Round(entropy*100) / 100
}

// AnalyzePredictability 分析文本可预测性
func (p *PerplexityCalculator) AnalyzePredictability(text string) map[string]float64 {
	tokens := p.tokenize(text)

	bigramCounts := p.buildNGramModelWithSize(tokens, 2)
	trigramCounts := p.buildNGramModelWithSize(tokens, 3)
	unigramCounts := p.buildNGramModelWithSize(tokens, 1)

	// 计算bigram熵
	bigramEntropy := p.calculateNGramEntropy(bigramCounts, unigramCounts)

	// 计算trigram熵
	trigramEntropy := p.calculateNGramEntropy(trigramCounts, bigramCounts)

	// 计算熵率下降（条件熵的递减表明可预测性增加）
	// AI生成的文本通常有更低的条件熵（更可预测）
	entropyReduction := bigramEntropy - trigramEntropy

	return map[string]float64{
		"unigram_count":    float64(len(unigramCounts)),
		"bigram_entropy":   math.Round(bigramEntropy*100) / 100,
		"trigram_entropy":  math.Round(trigramEntropy*100) / 100,
		"entropy_reduction": math.Round(entropyReduction*100) / 100,
	}
}

// calculateNGramEntropy 计算n-gram条件熵
func (p *PerplexityCalculator) calculateNGramEntropy(
	ngramCounts map[string]int,
	contextCounts map[string]int,
) float64 {
	if len(ngramCounts) == 0 {
		return 0
	}

	var totalEntropy float64
	var totalCount float64

	for ngram, count := range ngramCounts {
		// 获取上下文（去掉最后一个词）
		parts := strings.Split(ngram, " ")
		if len(parts) < 2 {
			continue
		}
		context := strings.Join(parts[:len(parts)-1], " ")

		contextCount := float64(contextCounts[context])
		if contextCount == 0 {
			continue
		}

		prob := float64(count) / contextCount
		if prob > 0 {
			totalEntropy -= float64(count) * math.Log2(prob)
			totalCount += float64(count)
		}
	}

	if totalCount == 0 {
		return 0
	}

	return totalEntropy / totalCount
}

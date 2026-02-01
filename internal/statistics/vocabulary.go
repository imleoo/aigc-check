package statistics

import (
	"regexp"
	"strings"
	"unicode"
)

// VocabularyAnalyzer 词汇分析器
type VocabularyAnalyzer struct {
	// 停用词列表（不参与多样性计算）
	stopWords map[string]bool
}

// VocabularyStats 词汇统计结果
type VocabularyStats struct {
	// 总词数
	TotalWords int `json:"total_words"`

	// 唯一词数
	UniqueWords int `json:"unique_words"`

	// Type-Token Ratio (TTR) = 唯一词/总词数
	// 范围: 0-1，越高表示词汇越多样
	TTR float64 `json:"ttr"`

	// 标准化TTR (考虑文本长度)
	// 使用移动平均TTR (MATTR)
	StandardizedTTR float64 `json:"standardized_ttr"`

	// Hapax Legomena（只出现一次的词）数量
	HapaxLegomena int `json:"hapax_legomena"`

	// Hapax Legomena 比例
	HapaxRatio float64 `json:"hapax_ratio"`

	// 词汇丰富度指数 (Yule's K)
	Richness float64 `json:"richness"`

	// 词频分布
	FrequencyDistribution map[int]int `json:"frequency_distribution"`

	// 高频词（出现>3次）
	HighFrequencyWords []WordFrequency `json:"high_frequency_words"`
}

// WordFrequency 词频信息
type WordFrequency struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

// NewVocabularyAnalyzer 创建词汇分析器
func NewVocabularyAnalyzer() *VocabularyAnalyzer {
	return &VocabularyAnalyzer{
		stopWords: getStopWords(),
	}
}

// Analyze 分析文本词汇特征
func (v *VocabularyAnalyzer) Analyze(text string) VocabularyStats {
	// 分词
	words := v.tokenize(text)

	if len(words) == 0 {
		return VocabularyStats{}
	}

	// 计算词频
	wordFreq := make(map[string]int)
	for _, word := range words {
		wordFreq[word]++
	}

	totalWords := len(words)
	uniqueWords := len(wordFreq)

	// 计算 TTR
	ttr := float64(uniqueWords) / float64(totalWords)

	// 计算标准化 TTR (MATTR with window size 50)
	standardizedTTR := v.calculateMATTR(words, 50)

	// 计算 Hapax Legomena
	hapaxLegomena := 0
	freqDist := make(map[int]int)
	for _, count := range wordFreq {
		freqDist[count]++
		if count == 1 {
			hapaxLegomena++
		}
	}

	hapaxRatio := float64(hapaxLegomena) / float64(uniqueWords)
	if uniqueWords == 0 {
		hapaxRatio = 0
	}

	// 计算词汇丰富度 (Yule's K)
	richness := v.calculateYulesK(wordFreq, totalWords)

	// 获取高频词
	highFreqWords := v.getHighFrequencyWords(wordFreq, 3)

	return VocabularyStats{
		TotalWords:            totalWords,
		UniqueWords:           uniqueWords,
		TTR:                   ttr,
		StandardizedTTR:       standardizedTTR,
		HapaxLegomena:         hapaxLegomena,
		HapaxRatio:            hapaxRatio,
		Richness:              richness,
		FrequencyDistribution: freqDist,
		HighFrequencyWords:    highFreqWords,
	}
}

// tokenize 分词处理
func (v *VocabularyAnalyzer) tokenize(text string) []string {
	// 转小写
	text = strings.ToLower(text)

	// 使用正则分词
	wordRegex := regexp.MustCompile(`[\p{L}\p{N}]+`)
	tokens := wordRegex.FindAllString(text, -1)

	// 过滤停用词和单字符
	var words []string
	for _, token := range tokens {
		if len(token) > 1 && !v.stopWords[token] {
			words = append(words, token)
		}
	}

	return words
}

// calculateMATTR 计算移动平均TTR
// 将文本分成窗口，计算每个窗口的TTR，然后取平均
func (v *VocabularyAnalyzer) calculateMATTR(words []string, windowSize int) float64 {
	if len(words) <= windowSize {
		return float64(len(unique(words))) / float64(len(words))
	}

	var totalTTR float64
	count := 0

	for i := 0; i <= len(words)-windowSize; i++ {
		window := words[i : i+windowSize]
		uniqueInWindow := len(unique(window))
		windowTTR := float64(uniqueInWindow) / float64(windowSize)
		totalTTR += windowTTR
		count++
	}

	return totalTTR / float64(count)
}

// calculateYulesK 计算 Yule's K 特征值
// K = 10000 * (M2 - M1) / (M1 * M1)
// 其中 M1 是总词数，M2 是 Σ(fi * fi)
func (v *VocabularyAnalyzer) calculateYulesK(wordFreq map[string]int, totalWords int) float64 {
	if totalWords == 0 {
		return 0
	}

	var m2 float64
	for _, freq := range wordFreq {
		m2 += float64(freq * freq)
	}

	m1 := float64(totalWords)
	if m1*m1-m2 == 0 {
		return 0
	}

	// Yule's K - 越小表示词汇越丰富
	// 为了与其他指标一致（越大越好），我们取倒数并缩放
	k := 10000 * (m2 - m1) / (m1 * m1)

	// 转换为丰富度分数（越大越丰富）
	// k 通常在 50-200 之间，我们转换到 0-20 的范围
	richness := 200 / (k + 10)

	return richness
}

// getHighFrequencyWords 获取高频词
func (v *VocabularyAnalyzer) getHighFrequencyWords(wordFreq map[string]int, minCount int) []WordFrequency {
	var result []WordFrequency

	for word, count := range wordFreq {
		if count >= minCount {
			result = append(result, WordFrequency{
				Word:  word,
				Count: count,
			})
		}
	}

	// 按频率排序
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Count > result[i].Count {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	// 限制返回数量
	if len(result) > 20 {
		result = result[:20]
	}

	return result
}

// unique 返回唯一元素
func unique(words []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, word := range words {
		if !seen[word] {
			seen[word] = true
			result = append(result, word)
		}
	}

	return result
}

// getStopWords 获取停用词列表
func getStopWords() map[string]bool {
	// 中英文停用词
	stopWordList := []string{
		// 英文停用词
		"a", "an", "the", "is", "are", "was", "were", "be", "been", "being",
		"have", "has", "had", "do", "does", "did", "will", "would", "could",
		"should", "may", "might", "must", "shall", "can", "need", "dare",
		"to", "of", "in", "for", "on", "with", "at", "by", "from", "up",
		"about", "into", "through", "during", "before", "after", "above",
		"below", "between", "under", "again", "further", "then", "once",
		"here", "there", "when", "where", "why", "how", "all", "each",
		"few", "more", "most", "other", "some", "such", "no", "nor", "not",
		"only", "own", "same", "so", "than", "too", "very", "just", "and",
		"but", "if", "or", "because", "as", "until", "while", "although",
		"though", "after", "before", "this", "that", "these", "those",
		"it", "its", "he", "him", "his", "she", "her", "hers", "they",
		"them", "their", "theirs", "we", "us", "our", "ours", "you", "your",
		"yours", "i", "me", "my", "mine", "what", "which", "who", "whom",

		// 中文停用词
		"的", "了", "是", "在", "我", "有", "和", "就", "不", "人",
		"都", "一", "一个", "上", "也", "很", "到", "说", "要", "去",
		"你", "会", "着", "没有", "看", "好", "自己", "这", "他", "她",
		"它", "们", "那", "还", "被", "把", "让", "给", "从", "向",
		"对", "与", "为", "以", "及", "等", "但", "而", "或", "且",
		"因为", "所以", "如果", "虽然", "但是", "然后", "于是",
	}

	stopWords := make(map[string]bool)
	for _, word := range stopWordList {
		stopWords[word] = true
	}

	return stopWords
}

// IsLetter 检查字符是否为字母
func IsLetter(r rune) bool {
	return unicode.IsLetter(r)
}

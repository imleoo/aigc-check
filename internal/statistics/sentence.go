package statistics

import (
	"math"
	"regexp"
	"strings"
	"unicode"
)

// SentenceAnalyzer 句子分析器
type SentenceAnalyzer struct {
	// 复杂度计算权重
	clauseWeight      float64
	conjunctionWeight float64
	commaWeight       float64
}

// SentenceStats 句子统计结果
type SentenceStats struct {
	// 句子总数
	TotalSentences int `json:"total_sentences"`

	// 平均句子长度（词数）
	AverageLength float64 `json:"average_length"`

	// 句子长度标准差
	LengthStdDev float64 `json:"length_std_dev"`

	// 句子长度方差
	LengthVariance float64 `json:"length_variance"`

	// 最短句子长度
	MinLength int `json:"min_length"`

	// 最长句子长度
	MaxLength int `json:"max_length"`

	// 句式复杂度评分 (0-100)
	ComplexityScore float64 `json:"complexity_score"`

	// 复杂句比例（包含从句、连词等）
	ComplexSentenceRatio float64 `json:"complex_sentence_ratio"`

	// 简单句比例
	SimpleSentenceRatio float64 `json:"simple_sentence_ratio"`

	// 问句比例
	QuestionRatio float64 `json:"question_ratio"`

	// 感叹句比例
	ExclamationRatio float64 `json:"exclamation_ratio"`

	// 句子长度分布
	LengthDistribution map[string]int `json:"length_distribution"`
}

// NewSentenceAnalyzer 创建句子分析器
func NewSentenceAnalyzer() *SentenceAnalyzer {
	return &SentenceAnalyzer{
		clauseWeight:      0.4,
		conjunctionWeight: 0.3,
		commaWeight:       0.3,
	}
}

// Analyze 分析文本句子特征
func (s *SentenceAnalyzer) Analyze(text string) SentenceStats {
	sentences := s.splitSentences(text)

	if len(sentences) == 0 {
		return SentenceStats{}
	}

	// 计算句子长度
	var lengths []int
	var totalLength int
	minLen, maxLen := math.MaxInt, 0
	questionCount := 0
	exclamationCount := 0
	complexCount := 0

	for _, sentence := range sentences {
		wordCount := s.countWords(sentence)
		lengths = append(lengths, wordCount)
		totalLength += wordCount

		if wordCount < minLen {
			minLen = wordCount
		}
		if wordCount > maxLen {
			maxLen = wordCount
		}

		// 统计句子类型
		trimmed := strings.TrimSpace(sentence)
		if strings.HasSuffix(trimmed, "?") || strings.HasSuffix(trimmed, "？") {
			questionCount++
		}
		if strings.HasSuffix(trimmed, "!") || strings.HasSuffix(trimmed, "！") {
			exclamationCount++
		}
		if s.isComplexSentence(sentence) {
			complexCount++
		}
	}

	totalSentences := len(sentences)
	avgLength := float64(totalLength) / float64(totalSentences)

	// 计算方差和标准差
	var variance float64
	for _, length := range lengths {
		diff := float64(length) - avgLength
		variance += diff * diff
	}
	variance /= float64(totalSentences)
	stdDev := math.Sqrt(variance)

	// 计算长度分布
	lengthDist := s.calculateLengthDistribution(lengths)

	// 计算复杂度评分
	complexityScore := s.calculateComplexityScore(sentences, avgLength, stdDev)

	return SentenceStats{
		TotalSentences:       totalSentences,
		AverageLength:        math.Round(avgLength*100) / 100,
		LengthStdDev:         math.Round(stdDev*100) / 100,
		LengthVariance:       math.Round(variance*100) / 100,
		MinLength:            minLen,
		MaxLength:            maxLen,
		ComplexityScore:      math.Round(complexityScore*100) / 100,
		ComplexSentenceRatio: math.Round(float64(complexCount)/float64(totalSentences)*100) / 100,
		SimpleSentenceRatio:  math.Round(float64(totalSentences-complexCount)/float64(totalSentences)*100) / 100,
		QuestionRatio:        math.Round(float64(questionCount)/float64(totalSentences)*100) / 100,
		ExclamationRatio:     math.Round(float64(exclamationCount)/float64(totalSentences)*100) / 100,
		LengthDistribution:   lengthDist,
	}
}

// splitSentences 分割句子
func (s *SentenceAnalyzer) splitSentences(text string) []string {
	// 中英文句子分割正则
	// 匹配句号、问号、感叹号（中英文）
	sentenceEndRegex := regexp.MustCompile(`[.。!！?？]+\s*`)

	// 先用换行符分割段落
	paragraphs := strings.Split(text, "\n")

	var sentences []string
	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		// 分割句子
		parts := sentenceEndRegex.Split(para, -1)
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if len(part) > 0 && s.countWords(part) > 1 {
				sentences = append(sentences, part)
			}
		}
	}

	return sentences
}

// countWords 统计句子中的词数
func (s *SentenceAnalyzer) countWords(sentence string) int {
	// 对于中文，按字符计数（每2-3个字符约等于1个词）
	// 对于英文，按空格分词

	// 检测是否主要是中文
	chineseCount := 0
	totalCount := 0
	for _, r := range sentence {
		if unicode.Is(unicode.Han, r) {
			chineseCount++
		}
		if unicode.IsLetter(r) {
			totalCount++
		}
	}

	if totalCount > 0 && float64(chineseCount)/float64(totalCount) > 0.5 {
		// 中文：每2个汉字约等于1个词
		return (chineseCount + 1) / 2
	}

	// 英文：按空格和标点分词
	wordRegex := regexp.MustCompile(`[\p{L}\p{N}]+`)
	words := wordRegex.FindAllString(sentence, -1)
	return len(words)
}

// isComplexSentence 判断是否为复杂句
func (s *SentenceAnalyzer) isComplexSentence(sentence string) bool {
	sentence = strings.ToLower(sentence)

	// 英文从句标记词
	subordinatingConj := []string{
		"although", "though", "because", "since", "while", "whereas",
		"if", "unless", "until", "when", "whenever", "where", "wherever",
		"after", "before", "as", "that", "which", "who", "whom", "whose",
	}

	// 中文从句标记词
	chineseMarkers := []string{
		"虽然", "尽管", "因为", "由于", "如果", "假如", "除非",
		"当", "在...时", "无论", "不管", "只要", "一旦",
	}

	// 检查英文从句标记
	for _, conj := range subordinatingConj {
		if strings.Contains(sentence, " "+conj+" ") {
			return true
		}
	}

	// 检查中文从句标记
	for _, marker := range chineseMarkers {
		if strings.Contains(sentence, marker) {
			return true
		}
	}

	// 检查逗号数量（多个逗号通常表示复杂句）
	commaCount := strings.Count(sentence, ",") + strings.Count(sentence, "，")
	if commaCount >= 2 {
		return true
	}

	return false
}

// calculateLengthDistribution 计算句子长度分布
func (s *SentenceAnalyzer) calculateLengthDistribution(lengths []int) map[string]int {
	dist := map[string]int{
		"very_short": 0, // <5 词
		"short":      0, // 5-10 词
		"medium":     0, // 11-20 词
		"long":       0, // 21-35 词
		"very_long":  0, // >35 词
	}

	for _, length := range lengths {
		switch {
		case length < 5:
			dist["very_short"]++
		case length < 11:
			dist["short"]++
		case length < 21:
			dist["medium"]++
		case length < 36:
			dist["long"]++
		default:
			dist["very_long"]++
		}
	}

	return dist
}

// calculateComplexityScore 计算句式复杂度评分
func (s *SentenceAnalyzer) calculateComplexityScore(sentences []string, avgLength float64, stdDev float64) float64 {
	if len(sentences) == 0 {
		return 0
	}

	var totalComplexity float64

	for _, sentence := range sentences {
		sentenceComplexity := 0.0

		// 基于句子长度的复杂度（10-25词最佳）
		wordCount := float64(s.countWords(sentence))
		if wordCount >= 10 && wordCount <= 25 {
			sentenceComplexity += 30
		} else if wordCount >= 5 && wordCount <= 35 {
			sentenceComplexity += 20
		} else {
			sentenceComplexity += 10
		}

		// 基于从句数量的复杂度
		if s.isComplexSentence(sentence) {
			sentenceComplexity += 25
		}

		// 基于逗号使用的复杂度（适当使用逗号）
		commaCount := strings.Count(sentence, ",") + strings.Count(sentence, "，")
		if commaCount >= 1 && commaCount <= 3 {
			sentenceComplexity += 15
		} else if commaCount > 3 {
			sentenceComplexity += 10
		}

		// 基于词汇多样性的复杂度（粗略估计）
		uniqueRatio := s.estimateUniqueWordRatio(sentence)
		sentenceComplexity += uniqueRatio * 30

		totalComplexity += sentenceComplexity
	}

	avgComplexity := totalComplexity / float64(len(sentences))

	// 加入句子长度变化的影响
	// 变化越大，表示写作风格越多样
	variationBonus := min(stdDev/avgLength*20, 15)
	avgComplexity += variationBonus

	// 确保分数在 0-100 范围内
	if avgComplexity > 100 {
		avgComplexity = 100
	}
	if avgComplexity < 0 {
		avgComplexity = 0
	}

	return avgComplexity
}

// estimateUniqueWordRatio 估算句子中的唯一词比例
func (s *SentenceAnalyzer) estimateUniqueWordRatio(sentence string) float64 {
	wordRegex := regexp.MustCompile(`[\p{L}\p{N}]+`)
	words := wordRegex.FindAllString(strings.ToLower(sentence), -1)

	if len(words) == 0 {
		return 0
	}

	seen := make(map[string]bool)
	for _, word := range words {
		seen[word] = true
	}

	return float64(len(seen)) / float64(len(words))
}

package text

import (
	"strings"
	"unicode"

	"github.com/leoobai/aigc-check/internal/models"
)

// ProcessedText 处理后的文本
type ProcessedText struct {
	Original   string     // 原始文本
	Normalized string     // 标准化文本
	Sentences  []Sentence // 句子列表
	Words      []Token    // 词汇列表
	Lines      []string   // 行列表
	WordCount  int        // 词数
	CharCount  int        // 字符数
}

// Sentence 句子
type Sentence struct {
	Text     string          // 句子文本
	Position models.Position // 位置信息
	Words    []Token         // 句子中的词汇
}

// Token 词汇标记
type Token struct {
	Text     string          // 词汇文本
	Position models.Position // 位置信息
	Lower    string          // 小写形式
}

// TextProcessor 文本处理器
type TextProcessor struct{}

// NewTextProcessor 创建文本处理器
func NewTextProcessor() *TextProcessor {
	return &TextProcessor{}
}

// Process 处理文本
func (p *TextProcessor) Process(text string) *ProcessedText {
	processed := &ProcessedText{
		Original: text,
	}

	// 标准化文本
	processed.Normalized = p.normalize(text)

	// 分割行
	processed.Lines = strings.Split(text, "\n")

	// 分割句子
	processed.Sentences = p.splitSentences(text)

	// 提取词汇
	processed.Words = p.extractWords(text)

	// 统计
	processed.WordCount = len(processed.Words)
	processed.CharCount = len([]rune(text))

	return processed
}

// normalize 标准化文本
func (p *TextProcessor) normalize(text string) string {
	// 统一空白符
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// 移除多余空白
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	return strings.Join(lines, "\n")
}

// splitSentences 分割句子
func (p *TextProcessor) splitSentences(text string) []Sentence {
	var sentences []Sentence
	var currentSentence strings.Builder
	var startOffset int
	var startLine, startColumn int = 1, 1

	runes := []rune(text)
	line, column := 1, 1

	for i, r := range runes {
		currentSentence.WriteRune(r)

		// 更新位置
		if r == '\n' {
			line++
			column = 1
		} else {
			column++
		}

		// 句子结束标记
		if p.isSentenceEnd(runes, i) {
			sentenceText := strings.TrimSpace(currentSentence.String())
			if sentenceText != "" {
				// 提取句子中的词汇
				words := p.extractWordsFromRange(text, startOffset, i+1)

				sentences = append(sentences, Sentence{
					Text: sentenceText,
					Position: models.Position{
						Line:   startLine,
						Column: startColumn,
						Offset: startOffset,
						Length: len(sentenceText),
					},
					Words: words,
				})
			}

			// 重置
			currentSentence.Reset()
			startOffset = i + 1
			startLine = line
			startColumn = column
		}
	}

	// 处理最后一个句子
	if currentSentence.Len() > 0 {
		sentenceText := strings.TrimSpace(currentSentence.String())
		if sentenceText != "" {
			words := p.extractWordsFromRange(text, startOffset, len(runes))
			sentences = append(sentences, Sentence{
				Text: sentenceText,
				Position: models.Position{
					Line:   startLine,
					Column: startColumn,
					Offset: startOffset,
					Length: len(sentenceText),
				},
				Words: words,
			})
		}
	}

	return sentences
}

// isSentenceEnd 判断是否为句子结束
func (p *TextProcessor) isSentenceEnd(runes []rune, i int) bool {
	if i >= len(runes) {
		return false
	}

	r := runes[i]

	// 句子结束标点
	if r == '.' || r == '!' || r == '?' || r == '。' || r == '！' || r == '？' {
		// 检查下一个字符
		if i+1 < len(runes) {
			next := runes[i+1]
			// 如果下一个是空白或换行，则为句子结束
			if unicode.IsSpace(next) || next == '\n' {
				return true
			}
		} else {
			// 文本结束
			return true
		}
	}

	// 换行也可能是句子结束
	if r == '\n' {
		return true
	}

	return false
}

// extractWords 提取所有词汇
func (p *TextProcessor) extractWords(text string) []Token {
	return p.extractWordsFromRange(text, 0, len([]rune(text)))
}

// extractWordsFromRange 从指定范围提取词汇
func (p *TextProcessor) extractWordsFromRange(text string, start, end int) []Token {
	var tokens []Token
	var currentWord strings.Builder
	var wordStart int
	var wordLine, wordColumn int = 1, 1

	runes := []rune(text)
	line, column := 1, 1
	offset := 0

	for i, r := range runes {
		// 跳过范围外的字符
		if i < start {
			if r == '\n' {
				line++
				column = 1
			} else {
				column++
			}
			offset++
			continue
		}
		if i >= end {
			break
		}

		// 判断是否为词汇字符
		if p.isWordChar(r) {
			if currentWord.Len() == 0 {
				wordStart = offset
				wordLine = line
				wordColumn = column
			}
			currentWord.WriteRune(r)
		} else {
			// 保存当前词汇
			if currentWord.Len() > 0 {
				word := currentWord.String()
				tokens = append(tokens, Token{
					Text:  word,
					Lower: strings.ToLower(word),
					Position: models.Position{
						Line:   wordLine,
						Column: wordColumn,
						Offset: wordStart,
						Length: len(word),
					},
				})
				currentWord.Reset()
			}
		}

		// 更新位置
		if r == '\n' {
			line++
			column = 1
		} else {
			column++
		}
		offset++
	}

	// 处理最后一个词汇
	if currentWord.Len() > 0 {
		word := currentWord.String()
		tokens = append(tokens, Token{
			Text:  word,
			Lower: strings.ToLower(word),
			Position: models.Position{
				Line:   wordLine,
				Column: wordColumn,
				Offset: wordStart,
				Length: len(word),
			},
		})
	}

	return tokens
}

// isWordChar 判断是否为词汇字符
func (p *TextProcessor) isWordChar(r rune) bool {
	// 字母、数字、中文字符、连字符
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '\''
}

// FindPattern 查找模式匹配
func (p *TextProcessor) FindPattern(text, pattern string, caseSensitive bool) []models.Position {
	var positions []models.Position

	searchText := text
	searchPattern := pattern

	if !caseSensitive {
		searchText = strings.ToLower(text)
		searchPattern = strings.ToLower(pattern)
	}

	offset := 0
	for {
		index := strings.Index(searchText[offset:], searchPattern)
		if index == -1 {
			break
		}

		actualOffset := offset + index
		line, column := p.GetLineColumn(text, actualOffset)

		positions = append(positions, models.Position{
			Line:   line,
			Column: column,
			Offset: actualOffset,
			Length: len(pattern),
		})

		offset = actualOffset + len(pattern)
	}

	return positions
}

// GetLineColumn 获取行列号
func (p *TextProcessor) GetLineColumn(text string, offset int) (int, int) {
	line := 1
	column := 1

	for i, r := range text {
		if i >= offset {
			break
		}
		if r == '\n' {
			line++
			column = 1
		} else {
			column++
		}
	}

	return line, column
}

// CountPattern 统计模式出现次数
func (p *TextProcessor) CountPattern(text, pattern string, caseSensitive bool) int {
	searchText := text
	searchPattern := pattern

	if !caseSensitive {
		searchText = strings.ToLower(text)
		searchPattern = strings.ToLower(pattern)
	}

	return strings.Count(searchText, searchPattern)
}

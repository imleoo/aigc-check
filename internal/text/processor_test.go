package text

import (
	"testing"
)

func TestNewTextProcessor(t *testing.T) {
	processor := NewTextProcessor()
	if processor == nil {
		t.Error("NewTextProcessor() returned nil")
	}
}

func TestProcess(t *testing.T) {
	processor := NewTextProcessor()

	tests := []struct {
		name          string
		input         string
		expectedWords int
		expectedLines int
	}{
		{
			name:          "简单英文文本",
			input:         "Hello world. This is a test.",
			expectedWords: 6, // "Hello", "world", "This", "is", "a", "test" (标点不计入)
			expectedLines: 1,
		},
		{
			name:          "多行文本",
			input:         "Line one.\nLine two.\nLine three.",
			expectedWords: 6,
			expectedLines: 3,
		},
		{
			name:          "空文本",
			input:         "",
			expectedWords: 0,
			expectedLines: 1,
		},
		{
			name:          "带标点的文本",
			input:         "Hello, world! How are you?",
			expectedWords: 5,
			expectedLines: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.Process(tt.input)

			if result == nil {
				t.Fatal("Process() returned nil")
			}

			if result.WordCount != tt.expectedWords {
				t.Errorf("WordCount = %d, want %d", result.WordCount, tt.expectedWords)
			}

			if len(result.Lines) != tt.expectedLines {
				t.Errorf("Lines count = %d, want %d", len(result.Lines), tt.expectedLines)
			}

			if result.Original != tt.input {
				t.Errorf("Original = %q, want %q", result.Original, tt.input)
			}
		})
	}
}

func TestSplitSentences(t *testing.T) {
	processor := NewTextProcessor()

	tests := []struct {
		name              string
		input             string
		expectedSentences int
	}{
		{
			name:              "单句",
			input:             "This is a sentence.",
			expectedSentences: 1,
		},
		{
			name:              "多句",
			input:             "First sentence. Second sentence. Third sentence.",
			expectedSentences: 3,
		},
		{
			name:              "带问号和感叹号",
			input:             "Is this a question? Yes! It is.",
			expectedSentences: 3,
		},
		{
			name:              "换行分隔",
			input:             "Line one\nLine two\nLine three",
			expectedSentences: 3,
		},
		{
			name:              "中文句子",
			input:             "这是第一句。\n这是第二句！\n这是第三句？",
			expectedSentences: 3, // 中文句子需要换行来分隔
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.Process(tt.input)

			if len(result.Sentences) != tt.expectedSentences {
				t.Errorf("Sentences count = %d, want %d", len(result.Sentences), tt.expectedSentences)
			}
		})
	}
}

func TestFindPattern(t *testing.T) {
	processor := NewTextProcessor()

	tests := []struct {
		name          string
		text          string
		pattern       string
		caseSensitive bool
		expectedCount int
	}{
		{
			name:          "精确匹配",
			text:          "The quick brown fox jumps over the lazy dog.",
			pattern:       "the",
			caseSensitive: true,
			expectedCount: 1,
		},
		{
			name:          "不区分大小写",
			text:          "The quick brown fox jumps over the lazy dog.",
			pattern:       "the",
			caseSensitive: false,
			expectedCount: 2,
		},
		{
			name:          "无匹配",
			text:          "Hello world",
			pattern:       "foo",
			caseSensitive: false,
			expectedCount: 0,
		},
		{
			name:          "多次匹配",
			text:          "crucial crucial crucial",
			pattern:       "crucial",
			caseSensitive: false,
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			positions := processor.FindPattern(tt.text, tt.pattern, tt.caseSensitive)

			if len(positions) != tt.expectedCount {
				t.Errorf("FindPattern() found %d matches, want %d", len(positions), tt.expectedCount)
			}
		})
	}
}

func TestCountPattern(t *testing.T) {
	processor := NewTextProcessor()

	tests := []struct {
		name          string
		text          string
		pattern       string
		caseSensitive bool
		expectedCount int
	}{
		{
			name:          "区分大小写",
			text:          "The The the the",
			pattern:       "The",
			caseSensitive: true,
			expectedCount: 2,
		},
		{
			name:          "不区分大小写",
			text:          "The The the the",
			pattern:       "the",
			caseSensitive: false,
			expectedCount: 4,
		},
		{
			name:          "无匹配",
			text:          "Hello world",
			pattern:       "foo",
			caseSensitive: false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := processor.CountPattern(tt.text, tt.pattern, tt.caseSensitive)

			if count != tt.expectedCount {
				t.Errorf("CountPattern() = %d, want %d", count, tt.expectedCount)
			}
		})
	}
}

func TestGetLineColumn(t *testing.T) {
	processor := NewTextProcessor()

	text := "Line 1\nLine 2\nLine 3"

	tests := []struct {
		offset       int
		expectedLine int
		expectedCol  int
	}{
		{0, 1, 1},   // 'L' of Line 1
		{5, 1, 6},   // '1' of Line 1
		{7, 2, 1},   // 'L' of Line 2
		{14, 3, 1},  // 'L' of Line 3
	}

	for _, tt := range tests {
		line, col := processor.GetLineColumn(text, tt.offset)

		if line != tt.expectedLine || col != tt.expectedCol {
			t.Errorf("GetLineColumn(%d) = (%d, %d), want (%d, %d)",
				tt.offset, line, col, tt.expectedLine, tt.expectedCol)
		}
	}
}

func TestExtractWords(t *testing.T) {
	processor := NewTextProcessor()

	tests := []struct {
		name          string
		input         string
		expectedWords []string
	}{
		{
			name:          "简单词汇",
			input:         "hello world",
			expectedWords: []string{"hello", "world"},
		},
		{
			name:          "带标点",
			input:         "Hello, world!",
			expectedWords: []string{"Hello", "world"},
		},
		{
			name:          "带连字符",
			input:         "self-driving car",
			expectedWords: []string{"self-driving", "car"},
		},
		{
			name:          "带撇号",
			input:         "I'm fine",
			expectedWords: []string{"I'm", "fine"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.Process(tt.input)

			if len(result.Words) != len(tt.expectedWords) {
				t.Errorf("Words count = %d, want %d", len(result.Words), len(tt.expectedWords))
				return
			}

			for i, word := range result.Words {
				if word.Text != tt.expectedWords[i] {
					t.Errorf("Word[%d] = %q, want %q", i, word.Text, tt.expectedWords[i])
				}
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	processor := NewTextProcessor()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "CRLF转LF",
			input:    "line1\r\nline2",
			expected: "line1\nline2",
		},
		{
			name:     "CR转LF",
			input:    "line1\rline2",
			expected: "line1\nline2",
		},
		{
			name:     "去除行首尾空白",
			input:    "  hello  \n  world  ",
			expected: "hello\nworld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.Process(tt.input)

			if result.Normalized != tt.expected {
				t.Errorf("Normalized = %q, want %q", result.Normalized, tt.expected)
			}
		})
	}
}

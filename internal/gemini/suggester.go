package gemini

import (
	"context"
	"fmt"
	"strings"
)

// Suggester 智能建议生成器
type Suggester struct {
	client *Client
}

// NewSuggester 创建建议生成器
func NewSuggester(client *Client) *Suggester {
	return &Suggester{
		client: client,
	}
}

// Suggestion 建议
type Suggestion struct {
	// 建议类型
	Type string `json:"type"`

	// 优先级 (1-5，1最高)
	Priority int `json:"priority"`

	// 标题
	Title string `json:"title"`

	// 详细描述
	Description string `json:"description"`

	// 原文片段
	OriginalText string `json:"original_text,omitempty"`

	// 建议修改后的文本
	SuggestedText string `json:"suggested_text,omitempty"`

	// 改进理由
	Reason string `json:"reason"`
}

// RewriteResult 重写结果
type RewriteResult struct {
	// 重写后的文本
	RewrittenText string `json:"rewritten_text"`

	// 主要修改点
	Changes []Change `json:"changes"`

	// 修改说明
	Explanation string `json:"explanation"`
}

// Change 修改点
type Change struct {
	Original string `json:"original"`
	Modified string `json:"modified"`
	Reason   string `json:"reason"`
}

// GenerateSuggestions 根据检测结果生成建议
func (s *Suggester) GenerateSuggestions(ctx context.Context, text string, issues []string) ([]Suggestion, error) {
	if len(issues) == 0 {
		return []Suggestion{}, nil
	}

	issuesList := strings.Join(issues, "\n- ")

	prompt := fmt.Sprintf(`你是一个写作顾问。根据以下检测到的问题，为文本提供具体的改进建议。

原文：
"""
%s
"""

检测到的问题：
- %s

请提供3-5条具体的改进建议，每条建议包括：
1. 问题所在的具体文本片段
2. 建议的修改方式
3. 修改后的示例
4. 为什么这样修改可以让文本更自然

请以JSON数组格式返回：
[
  {
    "type": "<问题类型>",
    "priority": <1-5>,
    "title": "<建议标题>",
    "description": "<详细描述>",
    "original_text": "<原文片段>",
    "suggested_text": "<修改后的文本>",
    "reason": "<改进理由>"
  }
]

请只返回JSON数组，不要有其他内容。`, text, issuesList)

	response, err := s.client.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate suggestions: %w", err)
	}

	var suggestions []Suggestion
	if err := parseJSONArrayResponse(response, &suggestions); err != nil {
		// 返回默认建议
		return s.getDefaultSuggestions(issues), nil
	}

	return suggestions, nil
}

// RewriteText 重写文本以降低AI痕迹
func (s *Suggester) RewriteText(ctx context.Context, text string, instructions string) (*RewriteResult, error) {
	if instructions == "" {
		instructions = "降低AI生成痕迹，使文本更加自然和人性化"
	}

	prompt := fmt.Sprintf(`你是一个文本改写专家。请根据以下要求改写文本：

要求：%s

原文：
"""
%s
"""

改写原则：
1. 保持原意不变
2. 添加适当的个人化表达
3. 使用更自然的词汇和句式
4. 避免过于完美的结构
5. 适当加入口语化表达

请以JSON格式返回：
{
  "rewritten_text": "<改写后的完整文本>",
  "changes": [
    {
      "original": "<原文片段>",
      "modified": "<修改后>",
      "reason": "<修改理由>"
    }
  ],
  "explanation": "<整体改写说明>"
}

请只返回JSON，不要有其他内容。`, instructions, text)

	response, err := s.client.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to rewrite text: %w", err)
	}

	result := &RewriteResult{}
	if err := parseJSONResponse(response, result); err != nil {
		return &RewriteResult{
			RewrittenText: text,
			Explanation:   "无法解析API响应，返回原文",
		}, nil
	}

	return result, nil
}

// ProvideAlternative 为特定片段提供替代表达
func (s *Suggester) ProvideAlternative(ctx context.Context, phrase string, context string) ([]string, error) {
	prompt := fmt.Sprintf(`请为以下短语/句子提供3-5个更自然、更人性化的替代表达。

短语："%s"
上下文：%s

要求：
1. 保持原意
2. 避免AI常用的表达方式
3. 使用更口语化或个性化的表达

请直接返回替代表达，每行一个，不要编号或其他格式。`, phrase, context)

	response, err := s.client.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to provide alternatives: %w", err)
	}

	// 按行分割
	lines := strings.Split(strings.TrimSpace(response), "\n")
	var alternatives []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 移除可能的编号
		line = strings.TrimLeft(line, "0123456789.-) ")
		if line != "" {
			alternatives = append(alternatives, line)
		}
	}

	return alternatives, nil
}

// getDefaultSuggestions 返回默认建议
func (s *Suggester) getDefaultSuggestions(issues []string) []Suggestion {
	var suggestions []Suggestion

	for i, issue := range issues {
		if i >= 5 {
			break
		}
		suggestions = append(suggestions, Suggestion{
			Type:        "general",
			Priority:    i + 1,
			Title:       "改进建议",
			Description: issue,
			Reason:      "基于检测到的问题提供的通用建议",
		})
	}

	return suggestions
}

// parseJSONArrayResponse 解析JSON数组响应
func parseJSONArrayResponse(response string, result interface{}) error {
	response = strings.TrimSpace(response)

	// 如果响应被包裹在```json...```中，提取内容
	if strings.HasPrefix(response, "```") {
		lines := strings.Split(response, "\n")
		var jsonLines []string
		inJSON := false
		for _, line := range lines {
			if strings.HasPrefix(line, "```") {
				if inJSON {
					break
				}
				inJSON = true
				continue
			}
			if inJSON {
				jsonLines = append(jsonLines, line)
			}
		}
		response = strings.Join(jsonLines, "\n")
	}

	// 尝试找到JSON数组的开始和结束
	start := strings.Index(response, "[")
	end := strings.LastIndex(response, "]")
	if start >= 0 && end > start {
		response = response[start : end+1]
	}

	return parseJSONResponse(response, result)
}

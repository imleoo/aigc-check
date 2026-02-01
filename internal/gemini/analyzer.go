package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// Analyzer 语义分析器
type Analyzer struct {
	client *Client
}

// NewAnalyzer 创建语义分析器
func NewAnalyzer(client *Client) *Analyzer {
	return &Analyzer{
		client: client,
	}
}

// AnalysisResult 分析结果
type AnalysisResult struct {
	// AI生成可能性 (0-100)
	AIProbability float64 `json:"ai_probability"`

	// 置信度 (0-1)
	Confidence float64 `json:"confidence"`

	// 检测到的特征
	Features []DetectedFeature `json:"features"`

	// 详细解释
	Explanation string `json:"explanation"`

	// 建议
	Suggestions []string `json:"suggestions"`
}

// DetectedFeature 检测到的特征
type DetectedFeature struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Severity    string  `json:"severity"` // low, medium, high
	Score       float64 `json:"score"`
}

// CoherenceResult 逻辑连贯性分析结果
type CoherenceResult struct {
	// 连贯性评分 (0-100)
	Score float64 `json:"score"`

	// 问题列表
	Issues []CoherenceIssue `json:"issues"`

	// 整体评价
	Assessment string `json:"assessment"`
}

// CoherenceIssue 连贯性问题
type CoherenceIssue struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Suggestion  string `json:"suggestion"`
}

// StyleResult 风格分析结果
type StyleResult struct {
	// 个人化程度 (0-100)
	PersonalizationScore float64 `json:"personalization_score"`

	// 检测到的风格特征
	StyleFeatures []string `json:"style_features"`

	// 缺失的人类写作特征
	MissingFeatures []string `json:"missing_features"`

	// 评估
	Assessment string `json:"assessment"`
}

// AnalyzeText 分析文本
func (a *Analyzer) AnalyzeText(ctx context.Context, text string) (*AnalysisResult, error) {
	prompt := fmt.Sprintf(`你是一个AI内容检测专家。请分析以下文本是否可能是AI生成的。

文本内容：
"""
%s
"""

请以JSON格式返回分析结果，包含以下字段：
{
  "ai_probability": <0-100的数字，表示AI生成的可能性>,
  "confidence": <0-1的数字，表示你对判断的置信度>,
  "features": [
    {
      "name": "<特征名称>",
      "description": "<特征描述>",
      "severity": "<low/medium/high>",
      "score": <0-100>
    }
  ],
  "explanation": "<详细解释为什么做出这个判断>",
  "suggestions": ["<改进建议1>", "<改进建议2>"]
}

请只返回JSON，不要有其他内容。`, text)

	response, err := a.client.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze text: %w", err)
	}

	// 解析响应
	result := &AnalysisResult{}
	if err := parseJSONResponse(response, result); err != nil {
		// 如果解析失败，返回默认结果
		return &AnalysisResult{
			AIProbability: 50,
			Confidence:    0.3,
			Explanation:   "无法解析API响应",
		}, nil
	}

	return result, nil
}

// AnalyzeLogicalCoherence 分析逻辑连贯性
func (a *Analyzer) AnalyzeLogicalCoherence(ctx context.Context, text string) (*CoherenceResult, error) {
	prompt := fmt.Sprintf(`你是一个文本分析专家。请分析以下文本的逻辑连贯性，特别关注：
1. 是否存在不自然的范围表达（如"从X到Y"但X和Y没有逻辑关联）
2. 是否存在论点之间的逻辑跳跃
3. 前后文是否一致
4. 是否存在AI生成常见的逻辑问题

文本内容：
"""
%s
"""

请以JSON格式返回分析结果：
{
  "score": <0-100的连贯性评分，越高越好>,
  "issues": [
    {
      "type": "<问题类型>",
      "description": "<问题描述>",
      "location": "<问题所在位置>",
      "suggestion": "<改进建议>"
    }
  ],
  "assessment": "<整体评价>"
}

请只返回JSON，不要有其他内容。`, text)

	response, err := a.client.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze coherence: %w", err)
	}

	result := &CoherenceResult{}
	if err := parseJSONResponse(response, result); err != nil {
		return &CoherenceResult{
			Score:      70,
			Assessment: "无法解析API响应",
		}, nil
	}

	return result, nil
}

// AnalyzePersonalStyle 分析个人风格
func (a *Analyzer) AnalyzePersonalStyle(ctx context.Context, text string) (*StyleResult, error) {
	prompt := fmt.Sprintf(`你是一个写作风格分析专家。请分析以下文本的个人化程度和写作风格特征。

人类写作通常具有以下特征：
- 使用第一人称表达观点
- 包含情感词汇和主观判断
- 使用不确定性表达（如"我认为"、"可能"）
- 具有个人经历的引用
- 口语化表达和语气词

AI生成的文本通常：
- 过于客观和正式
- 缺乏个人色彩
- 结构过于整齐
- 使用模板化的过渡词

文本内容：
"""
%s
"""

请以JSON格式返回分析结果：
{
  "personalization_score": <0-100，个人化程度，越高越像人类写作>,
  "style_features": ["<检测到的风格特征1>", "<风格特征2>"],
  "missing_features": ["<缺失的人类写作特征1>", "<缺失特征2>"],
  "assessment": "<整体评价>"
}

请只返回JSON，不要有其他内容。`, text)

	response, err := a.client.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze style: %w", err)
	}

	result := &StyleResult{}
	if err := parseJSONResponse(response, result); err != nil {
		return &StyleResult{
			PersonalizationScore: 50,
			Assessment:           "无法解析API响应",
		}, nil
	}

	return result, nil
}

// parseJSONResponse 解析JSON响应
func parseJSONResponse(response string, result interface{}) error {
	// 尝试提取JSON部分
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

	// 尝试找到JSON对象的开始和结束
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start >= 0 && end > start {
		response = response[start : end+1]
	}

	return json.Unmarshal([]byte(response), result)
}

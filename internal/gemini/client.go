package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client Gemini API 客户端
type Client struct {
	config     Config
	httpClient *http.Client
	cache      *Cache
}

// NewClient 创建 Gemini 客户端
func NewClient(cfg Config) (*Client, error) {
	// 从环境变量加载配置
	cfg.LoadFromEnv()

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// 创建 HTTP 客户端
	httpClient := &http.Client{
		Timeout: cfg.Timeout,
	}

	client := &Client{
		config:     cfg,
		httpClient: httpClient,
	}

	// 初始化缓存
	if cfg.Cache.Enabled {
		client.cache = NewCache(cfg.Cache)
	}

	return client, nil
}

// GenerateContentRequest Gemini API 请求
type GenerateContentRequest struct {
	Contents         []Content         `json:"contents"`
	GenerationConfig *GenerationConfig `json:"generationConfig,omitempty"`
}

// Content 内容
type Content struct {
	Parts []Part `json:"parts"`
	Role  string `json:"role,omitempty"`
}

// Part 内容部分
type Part struct {
	Text string `json:"text"`
}

// GenerationConfig 生成配置
type GenerationConfig struct {
	Temperature     float64 `json:"temperature,omitempty"`
	TopK            int     `json:"topK,omitempty"`
	TopP            float64 `json:"topP,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

// GenerateContentResponse Gemini API 响应
type GenerateContentResponse struct {
	Candidates []Candidate `json:"candidates"`
	Error      *APIError   `json:"error,omitempty"`
}

// Candidate 候选结果
type Candidate struct {
	Content       Content        `json:"content"`
	FinishReason  string         `json:"finishReason"`
	SafetyRatings []SafetyRating `json:"safetyRatings"`
}

// SafetyRating 安全评级
type SafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

// APIError API 错误
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// GenerateContent 生成内容
func (c *Client) GenerateContent(ctx context.Context, prompt string) (string, error) {
	if !c.config.IsEnabled() {
		return "", ErrNotEnabled
	}

	// 检查缓存
	if c.cache != nil {
		if cached, found := c.cache.Get(prompt); found {
			return cached, nil
		}
	}

	// 构建请求
	request := GenerateContentRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: &GenerationConfig{
			Temperature:     c.config.Temperature,
			MaxOutputTokens: c.config.MaxTokens,
		},
	}

	// 执行请求（带重试）
	var response GenerateContentResponse
	var lastErr error

	for attempt := 0; attempt < c.config.Retry.MaxAttempts; attempt++ {
		if attempt > 0 {
			// 计算退避时间
			backoff := c.calculateBackoff(attempt)
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(backoff):
			}
		}

		response, lastErr = c.doRequest(ctx, request)
		if lastErr == nil {
			break
		}

		// 检查是否应该重试
		if !c.shouldRetry(lastErr) {
			break
		}
	}

	if lastErr != nil {
		return "", lastErr
	}

	// 提取文本
	if len(response.Candidates) == 0 {
		return "", ErrNoResponse
	}

	var result string
	for _, part := range response.Candidates[0].Content.Parts {
		result += part.Text
	}

	// 存入缓存
	if c.cache != nil {
		c.cache.Set(prompt, result)
	}

	return result, nil
}

// doRequest 执行 HTTP 请求
func (c *Client) doRequest(ctx context.Context, request GenerateContentRequest) (GenerateContentResponse, error) {
	var response GenerateContentResponse

	// 序列化请求
	body, err := json.Marshal(request)
	if err != nil {
		return response, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 构建 URL
	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s",
		c.config.Endpoint, c.config.Model, c.config.APIKey)

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return response, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return response, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("failed to read response: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	if err := json.Unmarshal(respBody, &response); err != nil {
		return response, fmt.Errorf("failed to parse response: %w", err)
	}

	// 检查 API 错误
	if response.Error != nil {
		return response, fmt.Errorf("API error: %s", response.Error.Message)
	}

	return response, nil
}

// calculateBackoff 计算退避时间
func (c *Client) calculateBackoff(attempt int) time.Duration {
	backoff := c.config.Retry.InitialBackoff
	for i := 0; i < attempt; i++ {
		backoff = time.Duration(float64(backoff) * c.config.Retry.Multiplier)
		if backoff > c.config.Retry.MaxBackoff {
			backoff = c.config.Retry.MaxBackoff
			break
		}
	}
	return backoff
}

// shouldRetry 判断是否应该重试
func (c *Client) shouldRetry(err error) bool {
	// 可以根据错误类型决定是否重试
	// 例如：超时、5xx 错误应该重试
	// 4xx 错误通常不应该重试
	return true // 简单实现：总是重试
}

// Close 关闭客户端
func (c *Client) Close() error {
	if c.cache != nil {
		c.cache.Clear()
	}
	return nil
}

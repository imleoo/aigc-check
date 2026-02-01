package gemini

import "errors"

// 错误定义
var (
	// ErrMissingAPIKey API Key 未配置
	ErrMissingAPIKey = errors.New("gemini: API key is required")

	// ErrNotEnabled Gemini 未启用
	ErrNotEnabled = errors.New("gemini: not enabled")

	// ErrNoResponse 无响应
	ErrNoResponse = errors.New("gemini: no response from API")

	// ErrInvalidResponse 无效响应
	ErrInvalidResponse = errors.New("gemini: invalid response format")

	// ErrRateLimited 请求被限流
	ErrRateLimited = errors.New("gemini: rate limited")

	// ErrTimeout 请求超时
	ErrTimeout = errors.New("gemini: request timeout")

	// ErrContentBlocked 内容被安全过滤器阻止
	ErrContentBlocked = errors.New("gemini: content blocked by safety filters")
)

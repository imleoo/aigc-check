package models

import "time"

// DetectionRequest 表示检测请求
type DetectionRequest struct {
	Text     string            `json:"text"`      // 待检测文本
	Options  DetectionOptions  `json:"options"`   // 检测选项
	Metadata map[string]string `json:"metadata"`  // 元数据
}

// DetectionOptions 检测选项
type DetectionOptions struct {
	EnabledRules []string `json:"enabled_rules"` // 启用的规则列表，空表示全部启用
	Language     string   `json:"language"`      // 语言，默认"zh"
	OutputFormat string   `json:"output_format"` // 输出格式：text, json
}

// DetectionResult 表示检测结果
type DetectionResult struct {
	RequestID    string          `json:"request_id"`     // 请求ID
	Text         string          `json:"text"`           // 原始文本
	Score        Score           `json:"score"`          // 总体评分
	RuleResults  []RuleResult    `json:"rule_results"`   // 规则检测结果
	Suggestions  []Suggestion    `json:"suggestions"`    // 改进建议
	RiskLevel    RiskLevel       `json:"risk_level"`     // 风险等级
	ProcessTime  time.Duration   `json:"process_time"`   // 处理时间
	DetectedAt   time.Time       `json:"detected_at"`    // 检测时间
}

// RiskLevel 风险等级
type RiskLevel string

const (
	RiskLevelVeryHigh RiskLevel = "very_high" // 极高风险 (0-40)
	RiskLevelHigh     RiskLevel = "high"      // 高风险 (41-60)
	RiskLevelMedium   RiskLevel = "medium"    // 中等风险 (61-75)
	RiskLevelLow      RiskLevel = "low"       // 低风险 (76-100)
)

// GetRiskLevel 根据分数获取风险等级
func GetRiskLevel(score float64) RiskLevel {
	switch {
	case score <= 40:
		return RiskLevelVeryHigh
	case score <= 60:
		return RiskLevelHigh
	case score <= 75:
		return RiskLevelMedium
	default:
		return RiskLevelLow
	}
}

// GetRiskLevelDescription 获取风险等级描述
func (r RiskLevel) Description() string {
	switch r {
	case RiskLevelVeryHigh:
		return "极高风险 - 极可能为AI生成内容"
	case RiskLevelHigh:
		return "高风险 - 很可能为AI生成内容"
	case RiskLevelMedium:
		return "中等风险 - 可能包含AI生成内容"
	case RiskLevelLow:
		return "低风险 - 可能为人类编写"
	default:
		return "未知风险"
	}
}

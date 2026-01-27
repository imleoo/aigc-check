package models

// Suggestion 改进建议
type Suggestion struct {
	Category    SuggestionCategory `json:"category"`     // 建议类别
	Priority    Priority           `json:"priority"`     // 优先级
	Title       string             `json:"title"`        // 标题
	Description string             `json:"description"`  // 详细描述
	Examples    []Example          `json:"examples"`     // 示例
	RelatedRule RuleType           `json:"related_rule"` // 相关规则
}

// SuggestionCategory 建议类别
type SuggestionCategory string

const (
	CategoryVocabulary    SuggestionCategory = "vocabulary"     // 词汇
	CategorySentence      SuggestionCategory = "sentence"       // 句式
	CategoryTone          SuggestionCategory = "tone"           // 语气
	CategoryStructure     SuggestionCategory = "structure"      // 结构
	CategoryAuthenticity  SuggestionCategory = "authenticity"   // 真实性
	CategoryFormatting    SuggestionCategory = "formatting"     // 格式
)

// Priority 优先级
type Priority string

const (
	PriorityHigh   Priority = "high"    // 高优先级
	PriorityMedium Priority = "medium"  // 中优先级
	PriorityLow    Priority = "low"     // 低优先级
)

// Example 示例
type Example struct {
	Before string `json:"before"` // 修改前
	After  string `json:"after"`  // 修改后
	Reason string `json:"reason"` // 原因说明
}

// NewSuggestion 创建建议
func NewSuggestion(category SuggestionCategory, priority Priority, title, description string, relatedRule RuleType) Suggestion {
	return Suggestion{
		Category:    category,
		Priority:    priority,
		Title:       title,
		Description: description,
		Examples:    []Example{},
		RelatedRule: relatedRule,
	}
}

// AddExample 添加示例
func (s *Suggestion) AddExample(before, after, reason string) {
	s.Examples = append(s.Examples, Example{
		Before: before,
		After:  after,
		Reason: reason,
	})
}

// GetCategoryName 获取类别名称
func GetCategoryName(category SuggestionCategory) string {
	names := map[SuggestionCategory]string{
		CategoryVocabulary:   "词汇改进",
		CategorySentence:     "句式改进",
		CategoryTone:         "语气调整",
		CategoryStructure:    "结构优化",
		CategoryAuthenticity: "真实性提升",
		CategoryFormatting:   "格式清理",
	}
	if name, ok := names[category]; ok {
		return name
	}
	return string(category)
}

// GetPriorityName 获取优先级名称
func GetPriorityName(priority Priority) string {
	names := map[Priority]string{
		PriorityHigh:   "高",
		PriorityMedium: "中",
		PriorityLow:    "低",
	}
	if name, ok := names[priority]; ok {
		return name
	}
	return string(priority)
}

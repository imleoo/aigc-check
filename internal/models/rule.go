package models

// RuleType 规则类型
type RuleType string

const (
	// Signal 1: 高频词汇检测
	RuleTypeHighFreqWords RuleType = "high_frequency_words"

	// Signal 2: 句式开头检测
	RuleTypeSentenceStarters RuleType = "sentence_starters"

	// Signal 3: 虚假范围表达
	RuleTypeFalseRange RuleType = "false_range"

	// Signal 4: 引用异常检测
	RuleTypeCitationAnomaly RuleType = "citation_anomaly"

	// Signal 5: 破折号密度
	RuleTypeEmDash RuleType = "em_dash_density"

	// Signal 6: Markdown残留
	RuleTypeMarkdown RuleType = "markdown_residue"

	// Signal 7: 表情符号异常
	RuleTypeEmoji RuleType = "emoji_anomaly"

	// Signal 8: 知识截止日期
	RuleTypeKnowledgeCutoff RuleType = "knowledge_cutoff"

	// Signal 9: 协作式语气
	RuleTypeCollaborative RuleType = "collaborative_tone"

	// Signal 10: 完美主义陷阱
	RuleTypePerfectionism RuleType = "perfectionism"
)

// RuleResult 规则检测结果
type RuleResult struct {
	RuleType    RuleType      `json:"rule_type"`    // 规则类型
	RuleName    string        `json:"rule_name"`    // 规则名称
	Description string        `json:"description"`  // 规则描述
	Detected    bool          `json:"detected"`     // 是否检测到问题
	Score       float64       `json:"score"`        // 规则评分 (0-100)
	Severity    Severity      `json:"severity"`     // 严重程度
	Matches     []Match       `json:"matches"`      // 匹配项列表
	Count       int           `json:"count"`        // 匹配数量
	Threshold   int           `json:"threshold"`    // 阈值
	Message     string        `json:"message"`      // 结果消息
}

// Severity 严重程度
type Severity string

const (
	SeverityCritical Severity = "critical" // 严重
	SeverityHigh     Severity = "high"     // 高
	SeverityMedium   Severity = "medium"   // 中等
	SeverityLow      Severity = "low"      // 低
	SeverityInfo     Severity = "info"     // 信息
)

// Match 匹配项
type Match struct {
	Text     string   `json:"text"`      // 匹配的文本
	Position Position `json:"position"`  // 位置信息
	Context  string   `json:"context"`   // 上下文
	Reason   string   `json:"reason"`    // 匹配原因
}

// Position 位置信息
type Position struct {
	Line   int `json:"line"`   // 行号 (从1开始)
	Column int `json:"column"` // 列号 (从1开始)
	Offset int `json:"offset"` // 字符偏移量 (从0开始)
	Length int `json:"length"` // 匹配长度
}

// Rule 规则接口
type Rule interface {
	// Check 执行规则检测
	Check(text string) RuleResult

	// GetType 获取规则类型
	GetType() RuleType

	// GetName 获取规则名称
	GetName() string

	// GetDescription 获取规则描述
	GetDescription() string
}

// RuleConfig 规则配置
type RuleConfig struct {
	Enabled   bool                   `yaml:"enabled"`   // 是否启用
	Threshold int                    `yaml:"threshold"` // 阈值
	Severity  Severity               `yaml:"severity"`  // 严重程度
	Options   map[string]interface{} `yaml:"options"`   // 其他选项
}

// GetRuleTypeName 获取规则类型的显示名称
func GetRuleTypeName(ruleType RuleType) string {
	names := map[RuleType]string{
		RuleTypeHighFreqWords:     "高频词汇检测",
		RuleTypeSentenceStarters:  "句式开头检测",
		RuleTypeFalseRange:        "虚假范围表达",
		RuleTypeCitationAnomaly:   "引用异常检测",
		RuleTypeEmDash:            "破折号密度",
		RuleTypeMarkdown:          "Markdown残留",
		RuleTypeEmoji:             "表情符号异常",
		RuleTypeKnowledgeCutoff:   "知识截止日期",
		RuleTypeCollaborative:     "协作式语气",
		RuleTypePerfectionism:     "完美主义陷阱",
	}
	if name, ok := names[ruleType]; ok {
		return name
	}
	return string(ruleType)
}

// GetAllRuleTypes 获取所有规则类型
func GetAllRuleTypes() []RuleType {
	return []RuleType{
		RuleTypeHighFreqWords,
		RuleTypeSentenceStarters,
		RuleTypeFalseRange,
		RuleTypeCitationAnomaly,
		RuleTypeEmDash,
		RuleTypeMarkdown,
		RuleTypeEmoji,
		RuleTypeKnowledgeCutoff,
		RuleTypeCollaborative,
		RuleTypePerfectionism,
	}
}

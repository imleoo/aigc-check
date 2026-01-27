package config

// Thresholds 规则阈值配置
type Thresholds struct {
	HighFrequencyWords  HighFreqWordsThresholds  `yaml:"high_frequency_words"`  // Signal 1
	SentenceStarters    SentenceStartersThresholds `yaml:"sentence_starters"`   // Signal 2
	FalseRange          FalseRangeThresholds     `yaml:"false_range"`           // Signal 3
	CitationAnomaly     CitationAnomalyThresholds `yaml:"citation_anomaly"`     // Signal 4
	EmDashDensity       float64                  `yaml:"em_dash_density"`       // Signal 5: 每千字破折号数量
	MarkdownResidue     MarkdownThresholds       `yaml:"markdown_residue"`      // Signal 6
	EmojiAnomaly        EmojiThresholds          `yaml:"emoji_anomaly"`         // Signal 7
	KnowledgeCutoff     KnowledgeCutoffThresholds `yaml:"knowledge_cutoff"`     // Signal 8
	CollaborativeTone   CollaborativeThresholds  `yaml:"collaborative_tone"`    // Signal 9
	Perfectionism       PerfectionismThresholds  `yaml:"perfectionism"`         // Signal 10
}

// HighFreqWordsThresholds Signal 1 阈值
type HighFreqWordsThresholds struct {
	Keywords  []string `yaml:"keywords"`  // 关键词列表
	Threshold int      `yaml:"threshold"` // 出现次数阈值
}

// SentenceStartersThresholds Signal 2 阈值
type SentenceStartersThresholds struct {
	Patterns  []string `yaml:"patterns"`  // 句式开头模式
	Threshold int      `yaml:"threshold"` // 句子数量阈值
}

// FalseRangeThresholds Signal 3 阈值
type FalseRangeThresholds struct {
	Patterns  []string `yaml:"patterns"`  // 范围表达模式
	Threshold int      `yaml:"threshold"` // 检测阈值
}

// CitationAnomalyThresholds Signal 4 阈值
type CitationAnomalyThresholds struct {
	UTMPatterns       []string `yaml:"utm_patterns"`       // UTM参数模式
	GhostMarkers      []string `yaml:"ghost_markers"`      // 幽灵标记
	PlaceholderDates  []string `yaml:"placeholder_dates"`  // 占位符日期
}

// MarkdownThresholds Signal 6 阈值
type MarkdownThresholds struct {
	Patterns  []string `yaml:"patterns"`  // Markdown模式
	Threshold int      `yaml:"threshold"` // 检测阈值
}

// EmojiThresholds Signal 7 阈值
type EmojiThresholds struct {
	Threshold int `yaml:"threshold"` // 表情符号数量阈值
}

// KnowledgeCutoffThresholds Signal 8 阈值
type KnowledgeCutoffThresholds struct {
	Phrases []string `yaml:"phrases"` // 知识截止短语
}

// CollaborativeThresholds Signal 9 阈值
type CollaborativeThresholds struct {
	Phrases   []string `yaml:"phrases"`   // 协作式短语
	Threshold int      `yaml:"threshold"` // 检测阈值
}

// PerfectionismThresholds Signal 10 阈值
type PerfectionismThresholds struct {
	FirstPersonPronouns []string `yaml:"first_person_pronouns"` // 第一人称代词
	EmotionalWords      []string `yaml:"emotional_words"`       // 情感词汇
	UncertaintyMarkers  []string `yaml:"uncertainty_markers"`   // 不确定性标记
	Threshold           int      `yaml:"threshold"`             // 检测阈值
}

// DefaultThresholds 默认阈值配置
var DefaultThresholds = Thresholds{
	HighFrequencyWords: HighFreqWordsThresholds{
		Keywords: []string{
			"crucial", "pivotal", "vital", "groundbreaking",
			"revolutionary", "profound", "significant",
			"关键", "至关重要", "革命性", "突破性",
		},
		Threshold: 3,
	},
	SentenceStarters: SentenceStartersThresholds{
		Patterns: []string{
			"Additionally", "Furthermore", "Moreover",
			"此外", "另外", "而且",
		},
		Threshold: 5,
	},
	FalseRange: FalseRangeThresholds{
		Patterns: []string{
			"from .+ to .+",
			"从.+到.+",
		},
		Threshold: 2,
	},
	CitationAnomaly: CitationAnomalyThresholds{
		UTMPatterns: []string{
			"utm_source=chatgpt.com",
			"utm_source=openai",
		},
		GhostMarkers: []string{
			"contentReference[oaicite:",
			"[oai_citation:",
			"turn0search",
		},
		PlaceholderDates: []string{
			"2025-XX-XX",
			"YYYY-MM-DD",
			"20XX",
		},
	},
	EmDashDensity: 5.0, // 每千字5个破折号
	MarkdownResidue: MarkdownThresholds{
		Patterns: []string{
			"##",
			"**",
			"[.+]\\(.+\\)",
			"```",
		},
		Threshold: 3,
	},
	EmojiAnomaly: EmojiThresholds{
		Threshold: 5,
	},
	KnowledgeCutoff: KnowledgeCutoffThresholds{
		Phrases: []string{
			"As of my last knowledge update",
			"As of my training data",
			"截至我的知识更新",
			"根据我的训练数据",
		},
	},
	CollaborativeTone: CollaborativeThresholds{
		Phrases: []string{
			"I hope this helps",
			"Let me know if",
			"Feel free to",
			"希望这能帮到你",
			"如果需要",
			"随时告诉我",
		},
		Threshold: 3,
	},
	Perfectionism: PerfectionismThresholds{
		FirstPersonPronouns: []string{
			"I", "me", "my", "mine",
			"我", "我的",
		},
		EmotionalWords: []string{
			"feel", "think", "believe", "hope",
			"感觉", "认为", "相信", "希望",
		},
		UncertaintyMarkers: []string{
			"might", "maybe", "perhaps", "possibly",
			"可能", "也许", "或许",
		},
		Threshold: 5,
	},
}

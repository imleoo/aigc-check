package models

import (
	"testing"
)

func TestRuleType_String(t *testing.T) {
	tests := []struct {
		ruleType RuleType
		want     string
	}{
		{RuleTypeHighFreqWords, "high_frequency_words"},
		{RuleTypeSentenceStarters, "sentence_starters"},
		{RuleTypeFalseRange, "false_range"},
		{RuleTypeCitationAnomaly, "citation_anomaly"},
		{RuleTypeEmDash, "em_dash_density"},
		{RuleTypeMarkdown, "markdown_residue"},
		{RuleTypeEmoji, "emoji_anomaly"},
		{RuleTypeKnowledgeCutoff, "knowledge_cutoff"},
		{RuleTypeCollaborative, "collaborative_tone"},
		{RuleTypePerfectionism, "perfectionism"},
	}

	for _, tt := range tests {
		if string(tt.ruleType) != tt.want {
			t.Errorf("RuleType %s = %s, want %s", tt.ruleType, string(tt.ruleType), tt.want)
		}
	}
}

func TestSeverity_Constants(t *testing.T) {
	// 验证严重程度常量存在
	severities := []Severity{
		SeverityCritical,
		SeverityHigh,
		SeverityMedium,
		SeverityLow,
	}

	for _, s := range severities {
		if string(s) == "" {
			t.Errorf("Severity constant is empty")
		}
	}
}

func TestNewMatch(t *testing.T) {
	pos := Position{Line: 1, Column: 5, Offset: 10}
	match := Match{Text: "test text", Position: pos}

	if match.Text != "test text" {
		t.Errorf("Match.Text = %s, want 'test text'", match.Text)
	}
	if match.Position.Line != 1 {
		t.Errorf("Match.Position.Line = %d, want 1", match.Position.Line)
	}
	if match.Position.Column != 5 {
		t.Errorf("Match.Position.Column = %d, want 5", match.Position.Column)
	}
}

func TestRuleResult_Fields(t *testing.T) {
	result := RuleResult{
		RuleType:  RuleTypeHighFreqWords,
		RuleName:  "高频词汇检测",
		Detected:  true,
		Score:     75.5,
		Count:     3,
		Threshold: 3,
		Severity:  SeverityHigh,
		Message:   "检测到高频词汇",
		Matches:   []Match{{Text: "crucial", Position: Position{Line: 1}}},
	}

	if result.RuleType != RuleTypeHighFreqWords {
		t.Errorf("RuleResult.RuleType = %s, want high_freq_words", result.RuleType)
	}
	if !result.Detected {
		t.Error("RuleResult.Detected = false, want true")
	}
	if len(result.Matches) != 1 {
		t.Errorf("RuleResult.Matches length = %d, want 1", len(result.Matches))
	}
}

func TestPosition(t *testing.T) {
	pos := Position{
		Line:   10,
		Column: 20,
		Offset: 150,
	}

	if pos.Line != 10 {
		t.Errorf("Position.Line = %d, want 10", pos.Line)
	}
	if pos.Column != 20 {
		t.Errorf("Position.Column = %d, want 20", pos.Column)
	}
	if pos.Offset != 150 {
		t.Errorf("Position.Offset = %d, want 150", pos.Offset)
	}
}

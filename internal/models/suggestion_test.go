package models

import (
	"testing"
)

func TestSuggestionCategory_Constants(t *testing.T) {
	categories := []SuggestionCategory{
		CategoryVocabulary,
		CategorySentence,
		CategoryFormatting,
		CategoryTone,
		CategoryAuthenticity,
	}

	for _, c := range categories {
		if string(c) == "" {
			t.Error("SuggestionCategory constant is empty")
		}
	}
}

func TestPriority_Constants(t *testing.T) {
	priorities := []Priority{
		PriorityHigh,
		PriorityMedium,
		PriorityLow,
	}

	for _, p := range priorities {
		if string(p) == "" {
			t.Error("Priority constant is empty")
		}
	}
}

func TestNewSuggestion(t *testing.T) {
	suggestion := NewSuggestion(
		CategoryVocabulary,
		PriorityHigh,
		"Test Title",
		"Test Description",
		RuleTypeHighFreqWords,
	)

	if suggestion.Category != CategoryVocabulary {
		t.Errorf("Category = %s, want vocabulary", suggestion.Category)
	}
	if suggestion.Priority != PriorityHigh {
		t.Errorf("Priority = %s, want high", suggestion.Priority)
	}
	if suggestion.Title != "Test Title" {
		t.Errorf("Title = %s, want 'Test Title'", suggestion.Title)
	}
	if suggestion.Description != "Test Description" {
		t.Errorf("Description = %s, want 'Test Description'", suggestion.Description)
	}
	if suggestion.RelatedRule != RuleTypeHighFreqWords {
		t.Errorf("RelatedRule = %s, want high_freq_words", suggestion.RelatedRule)
	}
}

func TestSuggestion_AddExample(t *testing.T) {
	suggestion := NewSuggestion(
		CategoryVocabulary,
		PriorityHigh,
		"Title",
		"Description",
		RuleTypeHighFreqWords,
	)

	suggestion.AddExample("before text", "after text", "reason")

	if len(suggestion.Examples) != 1 {
		t.Fatalf("Examples length = %d, want 1", len(suggestion.Examples))
	}

	example := suggestion.Examples[0]
	if example.Before != "before text" {
		t.Errorf("Example.Before = %s, want 'before text'", example.Before)
	}
	if example.After != "after text" {
		t.Errorf("Example.After = %s, want 'after text'", example.After)
	}
	if example.Reason != "reason" {
		t.Errorf("Example.Reason = %s, want 'reason'", example.Reason)
	}
}

func TestGetCategoryName(t *testing.T) {
	tests := []struct {
		category SuggestionCategory
		want     string
	}{
		{CategoryVocabulary, "词汇改进"},
		{CategorySentence, "句式改进"},
		{CategoryFormatting, "格式清理"},
		{CategoryTone, "语气调整"},
		{CategoryAuthenticity, "真实性提升"},
	}

	for _, tt := range tests {
		got := GetCategoryName(tt.category)
		if got != tt.want {
			t.Errorf("GetCategoryName(%s) = %s, want %s", tt.category, got, tt.want)
		}
	}
}

func TestGetCategoryName_Unknown(t *testing.T) {
	got := GetCategoryName(SuggestionCategory("unknown"))
	if got != "unknown" {
		t.Errorf("GetCategoryName(unknown) = %s, want 'unknown'", got)
	}
}

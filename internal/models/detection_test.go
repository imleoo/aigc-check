package models

import (
	"testing"
	"time"
)

func TestDetectionRequest_Fields(t *testing.T) {
	request := DetectionRequest{
		Text: "Test text content",
		Options: DetectionOptions{
			EnabledRules: []string{"high_freq_words"},
			Language:     "zh",
			OutputFormat: "text",
		},
		Metadata: map[string]string{
			"source": "test",
		},
	}

	if request.Text != "Test text content" {
		t.Errorf("Text = %s, want 'Test text content'", request.Text)
	}
	if request.Options.Language != "zh" {
		t.Errorf("Options.Language = %s, want zh", request.Options.Language)
	}
	if len(request.Metadata) != 1 {
		t.Errorf("Metadata length = %d, want 1", len(request.Metadata))
	}
}

func TestDetectionResult_Fields(t *testing.T) {
	now := time.Now()
	result := DetectionResult{
		RequestID:   "test123",
		Text:        "Test text",
		Score:       Score{Total: 75.5},
		RuleResults: []RuleResult{{RuleType: RuleTypeHighFreqWords}},
		Suggestions: []Suggestion{{Title: "Test"}},
		RiskLevel:   RiskLevelMedium,
		ProcessTime: 100 * time.Millisecond,
		DetectedAt:  now,
	}

	if result.RequestID != "test123" {
		t.Errorf("RequestID = %s, want test123", result.RequestID)
	}
	if result.Score.Total != 75.5 {
		t.Errorf("Score.Total = %.1f, want 75.5", result.Score.Total)
	}
	if result.RiskLevel != RiskLevelMedium {
		t.Errorf("RiskLevel = %s, want medium", result.RiskLevel)
	}
	if result.ProcessTime != 100*time.Millisecond {
		t.Errorf("ProcessTime = %v, want 100ms", result.ProcessTime)
	}
}

func TestDetectionOptions_Defaults(t *testing.T) {
	options := DetectionOptions{}

	// 验证默认值为空
	if options.Language != "" {
		t.Errorf("Language default = %s, want empty", options.Language)
	}
	if len(options.EnabledRules) != 0 {
		t.Errorf("EnabledRules default length = %d, want 0", len(options.EnabledRules))
	}
}

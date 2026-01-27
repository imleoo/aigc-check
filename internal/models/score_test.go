package models

import (
	"testing"
)

func TestGetRiskLevel(t *testing.T) {
	tests := []struct {
		score float64
		want  RiskLevel
	}{
		{0, RiskLevelVeryHigh},
		{20, RiskLevelVeryHigh},
		{40, RiskLevelVeryHigh},
		{41, RiskLevelHigh},
		{50, RiskLevelHigh},
		{60, RiskLevelHigh},
		{61, RiskLevelMedium},
		{70, RiskLevelMedium},
		{75, RiskLevelMedium},
		{76, RiskLevelLow},
		{90, RiskLevelLow},
		{100, RiskLevelLow},
	}

	for _, tt := range tests {
		got := GetRiskLevel(tt.score)
		if got != tt.want {
			t.Errorf("GetRiskLevel(%.0f) = %s, want %s", tt.score, got, tt.want)
		}
	}
}

func TestRiskLevel_Description(t *testing.T) {
	tests := []struct {
		level       RiskLevel
		wantContain string
	}{
		{RiskLevelVeryHigh, "极高风险"},
		{RiskLevelHigh, "高风险"},
		{RiskLevelMedium, "中等风险"},
		{RiskLevelLow, "低风险"},
	}

	for _, tt := range tests {
		got := tt.level.Description()
		if got == "" {
			t.Errorf("RiskLevel(%s).Description() returned empty", tt.level)
		}
	}
}

func TestGetLevel(t *testing.T) {
	tests := []struct {
		percentage float64
		want       string
	}{
		{95, "优秀"},
		{90, "优秀"},
		{85, "良好"},
		{75, "良好"},
		{65, "一般"},
		{60, "一般"},
		{50, "较差"},
		{0, "较差"},
	}

	for _, tt := range tests {
		got := GetLevel(tt.percentage)
		if got != tt.want {
			t.Errorf("GetLevel(%.0f) = %s, want %s", tt.percentage, got, tt.want)
		}
	}
}

func TestCalculatePercentage(t *testing.T) {
	tests := []struct {
		score    float64
		maxScore float64
		want     float64
	}{
		{50, 100, 50},
		{20, 20, 100},
		{15, 30, 50},
		{0, 100, 0},
		{0, 0, 0}, // 除以零情况
	}

	for _, tt := range tests {
		got := CalculatePercentage(tt.score, tt.maxScore)
		if got != tt.want {
			t.Errorf("CalculatePercentage(%.0f, %.0f) = %.0f, want %.0f",
				tt.score, tt.maxScore, got, tt.want)
		}
	}
}

func TestNewDimensionScore(t *testing.T) {
	issues := []string{"Issue 1", "Issue 2"}
	ds := NewDimensionScore(15, 20, issues, "Test description")

	if ds.Score != 15 {
		t.Errorf("Score = %.0f, want 15", ds.Score)
	}
	if ds.MaxScore != 20 {
		t.Errorf("MaxScore = %.0f, want 20", ds.MaxScore)
	}
	if ds.Percentage != 75 {
		t.Errorf("Percentage = %.0f, want 75", ds.Percentage)
	}
	if ds.Level != "良好" {
		t.Errorf("Level = %s, want 良好", ds.Level)
	}
	if ds.Description != "Test description" {
		t.Errorf("Description = %s, want 'Test description'", ds.Description)
	}
	if len(ds.Issues) != 2 {
		t.Errorf("Issues length = %d, want 2", len(ds.Issues))
	}
}

func TestDefaultDimensionWeights(t *testing.T) {
	total := DefaultDimensionWeights.VocabularyDiversity +
		DefaultDimensionWeights.SentenceComplexity +
		DefaultDimensionWeights.Personalization +
		DefaultDimensionWeights.LogicalCoherence +
		DefaultDimensionWeights.EmotionalAuthenticity

	if total != 100 {
		t.Errorf("Total weights = %.0f, want 100", total)
	}
}

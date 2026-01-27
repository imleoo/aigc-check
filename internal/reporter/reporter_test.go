package reporter

import (
	"strings"
	"testing"
	"time"

	"github.com/leoobai/aigc-check/internal/models"
)

func TestNewTextReporter(t *testing.T) {
	reporter := NewTextReporter(true)
	if reporter == nil {
		t.Fatal("NewTextReporter() returned nil")
	}
	if !reporter.colorEnabled {
		t.Error("NewTextReporter(true) colorEnabled = false, want true")
	}

	reporter2 := NewTextReporter(false)
	if reporter2.colorEnabled {
		t.Error("NewTextReporter(false) colorEnabled = true, want false")
	}
}

func TestTextReporter_Format(t *testing.T) {
	reporter := NewTextReporter(false)
	if reporter.Format() != "text" {
		t.Errorf("Format() = %s, want text", reporter.Format())
	}
}

func TestTextReporter_Generate(t *testing.T) {
	reporter := NewTextReporter(false)

	result := &models.DetectionResult{
		RequestID: "test123",
		Text:      "Test text content",
		Score: models.Score{
			Total: 75.5,
			Dimensions: models.DimensionScores{
				VocabularyDiversity:   models.NewDimensionScore(18, 20, nil, "Good"),
				SentenceComplexity:    models.NewDimensionScore(12, 15, nil, "Good"),
				Personalization:       models.NewDimensionScore(20, 25, nil, "Good"),
				LogicalCoherence:      models.NewDimensionScore(16, 20, nil, "Good"),
				EmotionalAuthenticity: models.NewDimensionScore(16, 20, nil, "Good"),
			},
		},
		RuleResults: []models.RuleResult{
			{
				RuleType:  models.RuleTypeHighFreqWords,
				RuleName:  "高频词汇检测",
				Detected:  true,
				Score:     70,
				Count:     3,
				Threshold: 3,
				Severity:  models.SeverityHigh,
				Message:   "检测到高频AI词汇",
			},
		},
		Suggestions: []models.Suggestion{
			{
				Category:    models.CategoryVocabulary,
				Priority:    models.PriorityHigh,
				Title:       "减少高频词汇",
				Description: "使用更多样化的词汇",
			},
		},
		RiskLevel:   models.RiskLevelMedium,
		ProcessTime: 100 * time.Millisecond,
		DetectedAt:  time.Now(),
	}

	output, err := reporter.Generate(result)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// 验证输出包含关键内容
	checks := []string{
		"AIGC-Check",
		"总体评分",
		"风险等级",
		"维度评分",
		"检测到的问题",
		"改进建议",
		"处理时间",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("Generate() output missing: %s", check)
		}
	}
}

func TestTextReporter_Generate_NoIssues(t *testing.T) {
	reporter := NewTextReporter(false)

	result := &models.DetectionResult{
		Score: models.Score{
			Total: 95,
			Dimensions: models.DimensionScores{
				VocabularyDiversity:   models.NewDimensionScore(19, 20, nil, "Excellent"),
				SentenceComplexity:    models.NewDimensionScore(14, 15, nil, "Excellent"),
				Personalization:       models.NewDimensionScore(24, 25, nil, "Excellent"),
				LogicalCoherence:      models.NewDimensionScore(19, 20, nil, "Excellent"),
				EmotionalAuthenticity: models.NewDimensionScore(19, 20, nil, "Excellent"),
			},
		},
		RuleResults: []models.RuleResult{
			{Detected: false},
		},
		Suggestions: []models.Suggestion{},
		RiskLevel:   models.RiskLevelLow,
		ProcessTime: 50 * time.Millisecond,
		DetectedAt:  time.Now(),
	}

	output, err := reporter.Generate(result)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(output, "未检测到明显的AI生成特征") {
		t.Error("Generate() should indicate no issues detected")
	}
}

func TestTextReporter_CreateScoreBar(t *testing.T) {
	reporter := NewTextReporter(false)

	tests := []struct {
		score    float64
		wantLen  int
	}{
		{0, 52},   // [50个░]
		{50, 52},  // [25个█ + 25个░]
		{100, 52}, // [50个█]
	}

	for _, tt := range tests {
		bar := reporter.createScoreBar(tt.score)
		// 验证长度合理
		if len(bar) < 10 {
			t.Errorf("createScoreBar(%.0f) too short: %d", tt.score, len(bar))
		}
	}
}

func TestTextReporter_GetRiskIcon(t *testing.T) {
	reporter := NewTextReporter(false)

	tests := []struct {
		level    models.RiskLevel
		wantIcon string
	}{
		{models.RiskLevelLow, "✓"},
		{models.RiskLevelMedium, "⚠"},
		{models.RiskLevelHigh, "⚠⚠"},
		{models.RiskLevelVeryHigh, "⚠⚠⚠"},
	}

	for _, tt := range tests {
		icon := reporter.getRiskIcon(tt.level)
		if icon != tt.wantIcon {
			t.Errorf("getRiskIcon(%s) = %s, want %s", tt.level, icon, tt.wantIcon)
		}
	}
}

func TestTextReporter_GetSeverityIcon(t *testing.T) {
	reporter := NewTextReporter(false)

	tests := []struct {
		severity models.Severity
		wantIcon string
	}{
		{models.SeverityCritical, "🔴"},
		{models.SeverityHigh, "🟠"},
		{models.SeverityMedium, "🟡"},
		{models.SeverityLow, "🟢"},
	}

	for _, tt := range tests {
		icon := reporter.getSeverityIcon(tt.severity)
		if icon != tt.wantIcon {
			t.Errorf("getSeverityIcon(%s) = %s, want %s", tt.severity, icon, tt.wantIcon)
		}
	}
}

func TestTextReporter_GetPriorityIcon(t *testing.T) {
	reporter := NewTextReporter(false)

	tests := []struct {
		priority models.Priority
		wantIcon string
	}{
		{models.PriorityHigh, "🔴"},
		{models.PriorityMedium, "🟡"},
		{models.PriorityLow, "🟢"},
	}

	for _, tt := range tests {
		icon := reporter.getPriorityIcon(tt.priority)
		if icon != tt.wantIcon {
			t.Errorf("getPriorityIcon(%s) = %s, want %s", tt.priority, icon, tt.wantIcon)
		}
	}
}

func TestTextReporter_ColorEnabled(t *testing.T) {
	reporterWithColor := NewTextReporter(true)
	reporterNoColor := NewTextReporter(false)

	// 验证颜色功能
	colorScore := reporterWithColor.getScoreColor(90)
	noColorScore := reporterNoColor.getScoreColor(90)

	if colorScore == "" {
		t.Error("getScoreColor() with color enabled should return color code")
	}
	if noColorScore != "" {
		t.Error("getScoreColor() with color disabled should return empty string")
	}

	// 验证重置
	if reporterWithColor.colorReset() == "" {
		t.Error("colorReset() with color enabled should return reset code")
	}
	if reporterNoColor.colorReset() != "" {
		t.Error("colorReset() with color disabled should return empty string")
	}
}

func TestTextReporter_GetScoreColor(t *testing.T) {
	reporter := NewTextReporter(true)

	tests := []struct {
		score     float64
		wantColor string
	}{
		{90, "\033[32m"},  // 绿色
		{70, "\033[33m"},  // 黄色
		{50, "\033[31m"},  // 红色
		{30, "\033[35m"},  // 紫色
	}

	for _, tt := range tests {
		color := reporter.getScoreColor(tt.score)
		if color != tt.wantColor {
			t.Errorf("getScoreColor(%.0f) = %s, want %s", tt.score, color, tt.wantColor)
		}
	}
}

func TestTextReporter_GetRiskColor(t *testing.T) {
	reporter := NewTextReporter(true)

	tests := []struct {
		level     models.RiskLevel
		wantColor string
	}{
		{models.RiskLevelLow, "\033[32m"},
		{models.RiskLevelMedium, "\033[33m"},
		{models.RiskLevelHigh, "\033[31m"},
		{models.RiskLevelVeryHigh, "\033[35m"},
	}

	for _, tt := range tests {
		color := reporter.getRiskColor(tt.level)
		if color != tt.wantColor {
			t.Errorf("getRiskColor(%s) = %s, want %s", tt.level, color, tt.wantColor)
		}
	}
}

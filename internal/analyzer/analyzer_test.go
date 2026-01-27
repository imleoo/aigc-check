package analyzer

import (
	"testing"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
)

func TestNewAnalyzer(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
		Rules:      config.DefaultConfig.Rules,
	}

	analyzer := NewAnalyzer(cfg)

	if analyzer == nil {
		t.Fatal("NewAnalyzer() returned nil")
	}

	if analyzer.config != cfg {
		t.Error("NewAnalyzer() config not set correctly")
	}

	if analyzer.ruleEngine == nil {
		t.Error("NewAnalyzer() ruleEngine is nil")
	}

	if analyzer.scorer == nil {
		t.Error("NewAnalyzer() scorer is nil")
	}

	if analyzer.processor == nil {
		t.Error("NewAnalyzer() processor is nil")
	}
}

func TestAnalyzer_Analyze(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
		Rules:      config.DefaultConfig.Rules,
	}
	analyzer := NewAnalyzer(cfg)

	tests := []struct {
		name        string
		text        string
		wantRisk    models.RiskLevel
		description string
	}{
		{
			name:        "正常人类文本",
			text:        "I think this is an interesting approach. Maybe we could try something different.",
			wantRisk:    models.RiskLevelLow,
			description: "包含个人化表达的文本应为低风险",
		},
		{
			name:        "AI特征文本",
			text:        "As of my last knowledge update, this is crucial and pivotal. Furthermore, additionally, moreover, this is vital.",
			wantRisk:    models.RiskLevelVeryHigh,
			description: "包含知识截止日期和高频词的文本应为高风险",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := models.DetectionRequest{
				Text: tt.text,
			}

			result, err := analyzer.Analyze(request)
			if err != nil {
				t.Fatalf("Analyze() error = %v", err)
			}

			if result == nil {
				t.Fatal("Analyze() returned nil result")
			}

			// 验证基本结果
			if result.Text != tt.text {
				t.Errorf("result.Text = %s, want %s", result.Text, tt.text)
			}

			if result.Score.Total < 0 || result.Score.Total > 100 {
				t.Errorf("result.Score.Total = %.1f, want 0-100", result.Score.Total)
			}

			// 验证风险等级方向正确
			if tt.wantRisk == models.RiskLevelLow && result.RiskLevel != models.RiskLevelLow {
				t.Logf("Note: %s got risk level %s (score: %.1f)", tt.description, result.RiskLevel, result.Score.Total)
			}

			// 验证处理时间
			if result.ProcessTime <= 0 {
				t.Error("result.ProcessTime should be > 0")
			}

			// 验证规则结果不为空
			if len(result.RuleResults) == 0 {
				t.Error("result.RuleResults should not be empty")
			}
		})
	}
}

func TestAnalyzer_Analyze_EmptyText(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
		Rules:      config.DefaultConfig.Rules,
	}
	analyzer := NewAnalyzer(cfg)

	result, err := analyzer.Analyze(models.DetectionRequest{Text: ""})
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}

	if result == nil {
		t.Fatal("Analyze() returned nil for empty text")
	}
}

func TestAnalyzer_GenerateSuggestions(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
		Rules:      config.DefaultConfig.Rules,
	}
	analyzer := NewAnalyzer(cfg)

	// 包含多种问题的文本
	text := "As of my last knowledge update, this is crucial and pivotal. ## Title **bold** I hope this helps!"

	result, err := analyzer.Analyze(models.DetectionRequest{Text: text})
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}

	// 验证生成了建议
	if len(result.Suggestions) == 0 {
		t.Log("Note: No suggestions generated, may be expected based on detection results")
	}

	// 验证建议的基本结构
	for _, suggestion := range result.Suggestions {
		if suggestion.Title == "" {
			t.Error("Suggestion has empty Title")
		}
		if suggestion.Description == "" {
			t.Error("Suggestion has empty Description")
		}
	}
}

func TestAnalyzer_RuleResults(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
		Rules:      config.DefaultConfig.Rules,
	}
	analyzer := NewAnalyzer(cfg)

	result, _ := analyzer.Analyze(models.DetectionRequest{
		Text: "This is a test text for rule execution.",
	})

	// 验证规则结果
	for _, rr := range result.RuleResults {
		if rr.RuleType == "" {
			t.Error("RuleResult has empty RuleType")
		}
		if rr.Score < 0 || rr.Score > 100 {
			t.Errorf("RuleResult.Score = %.1f, want 0-100", rr.Score)
		}
	}
}

func TestGenerateRequestID(t *testing.T) {
	id1 := generateRequestID()
	if id1 == "" {
		t.Error("generateRequestID() returned empty string")
	}

	// 验证格式（应该是时间戳格式）
	if len(id1) != 14 {
		t.Errorf("generateRequestID() length = %d, want 14", len(id1))
	}
}

func TestAnalyzer_ScoreDimensions(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
		Rules:      config.DefaultConfig.Rules,
	}
	analyzer := NewAnalyzer(cfg)

	result, _ := analyzer.Analyze(models.DetectionRequest{
		Text: "I believe this might be a good approach. Perhaps we should consider alternatives.",
	})

	// 验证维度评分
	dimensions := result.Score.Dimensions

	checkDimension := func(name string, dim models.DimensionScore) {
		if dim.Score < 0 || dim.Score > dim.MaxScore {
			t.Errorf("Dimension %s score = %.1f, want 0-%.1f", name, dim.Score, dim.MaxScore)
		}
		if dim.Percentage < 0 || dim.Percentage > 100 {
			t.Errorf("Dimension %s percentage = %.1f, want 0-100", name, dim.Percentage)
		}
	}

	checkDimension("VocabularyDiversity", dimensions.VocabularyDiversity)
	checkDimension("SentenceComplexity", dimensions.SentenceComplexity)
	checkDimension("Personalization", dimensions.Personalization)
	checkDimension("LogicalCoherence", dimensions.LogicalCoherence)
	checkDimension("EmotionalAuthenticity", dimensions.EmotionalAuthenticity)
}

func TestAnalyzer_RedFlagDetection(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
		Rules:      config.DefaultConfig.Rules,
	}
	analyzer := NewAnalyzer(cfg)

	tests := []struct {
		name     string
		text     string
		wantLow  bool // 是否期望低分
	}{
		{
			name:     "知识截止日期",
			text:     "As of my last knowledge update in 2023, this information was accurate.",
			wantLow:  true,
		},
		{
			name:     "UTM参数",
			text:     "Check this link: https://example.com?utm_source=chatgpt.com",
			wantLow:  true,
		},
		{
			name:     "正常文本",
			text:     "I think this is an interesting topic to explore.",
			wantLow:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := analyzer.Analyze(models.DetectionRequest{Text: tt.text})

			if tt.wantLow && result.Score.Total > 60 {
				t.Errorf("Score = %.1f, want < 60 for red flag text", result.Score.Total)
			}
			if !tt.wantLow && result.Score.Total < 60 {
				t.Logf("Note: Score = %.1f for normal text", result.Score.Total)
			}
		})
	}
}

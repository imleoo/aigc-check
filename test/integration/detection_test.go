package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leoobai/aigc-check/internal/analyzer"
	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
)

func TestAnalyzer_Integration(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
		Scoring: config.ScoringConfig{
			Weights: models.DefaultDimensionWeights,
		},
	}

	analyzer := analyzer.NewAnalyzer(cfg)

	t.Run("分析AI生成文本", func(t *testing.T) {
		text := `This is a comprehensive guide to understanding artificial intelligence.

Additionally, it is crucial to understand the pivotal role that machine learning plays in modern technology. Furthermore, the groundbreaking advancements in deep learning have revolutionized many industries. Moreover, the profound impact of AI on society cannot be overstated.

This technology spans from healthcare to finance — transforming how we approach complex problems — enabling unprecedented efficiency — and creating new opportunities for innovation.

As of my last knowledge update, there are several key areas where AI has made significant progress. The field continues to evolve rapidly, and I hope this helps you understand the basics.

## Key Benefits

**Improved Efficiency**: AI systems can process vast amounts of data quickly.

**Better Decision Making**: Machine learning algorithms can identify patterns that humans might miss.

Let me know if you have any questions about this topic. Feel free to reach out if you need further clarification.

The solution is optimal and will work perfectly in most scenarios. This approach is comprehensive and covers all the essential aspects of the subject matter.`

		request := models.DetectionRequest{
			Text: text,
		}

		result, err := analyzer.Analyze(request)
		if err != nil {
			t.Fatalf("Analyze() error = %v", err)
		}

		// AI生成文本应该得到较低的分数
		if result.Score.Total > 60 {
			t.Errorf("AI text score = %.1f, expected < 60", result.Score.Total)
		}

		// 应该检测到高风险
		if result.RiskLevel != models.RiskLevelHigh && result.RiskLevel != models.RiskLevelVeryHigh {
			t.Errorf("RiskLevel = %s, expected High or VeryHigh", result.RiskLevel)
		}

		// 应该有一些检测结果
		detectedCount := 0
		for _, r := range result.RuleResults {
			if r.Detected {
				detectedCount++
			}
		}
		if detectedCount < 3 {
			t.Errorf("Detected rules = %d, expected >= 3", detectedCount)
		}

		// 应该生成建议
		if len(result.Suggestions) == 0 {
			t.Error("Expected suggestions, got none")
		}
	})

	t.Run("分析人类编写文本", func(t *testing.T) {
		text := `I've been thinking about artificial intelligence lately, and honestly, I'm not entirely sure what to make of it all.

Last week, I tried out one of those AI chatbots everyone's been talking about. It was pretty cool, I guess, though sometimes the responses felt a bit... robotic? Hard to explain, really.

My friend Sarah disagrees with me on this. She thinks AI is going to change everything in the next decade or so. Maybe she's right, but I've seen too many hyped technologies fizzle out to be completely convinced.

The thing that struck me most was how the AI handled ambiguity. When I asked vague questions, it would sometimes give these overly confident answers that turned out to be wrong. That made me a bit nervous, to be honest.

I don't know. I think AI has potential, but I also think we need to be careful about how much we rely on it. Just my two cents, anyway.

What do you all think? Am I being too skeptical here?`

		request := models.DetectionRequest{
			Text: text,
		}

		result, err := analyzer.Analyze(request)
		if err != nil {
			t.Fatalf("Analyze() error = %v", err)
		}

		// 人类文本应该得到较高的分数
		if result.Score.Total < 70 {
			t.Errorf("Human text score = %.1f, expected > 70", result.Score.Total)
		}

		// 应该是低风险或中等风险
		if result.RiskLevel == models.RiskLevelVeryHigh {
			t.Errorf("RiskLevel = %s, expected not VeryHigh for human text", result.RiskLevel)
		}
	})
}

func TestAnalyzer_WithTestData(t *testing.T) {
	// 获取测试数据目录
	testDataDir := filepath.Join("..", "..", "test", "testdata")

	// 检查目录是否存在
	if _, err := os.Stat(testDataDir); os.IsNotExist(err) {
		t.Skip("Test data directory not found, skipping file-based tests")
	}

	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
		Scoring: config.ScoringConfig{
			Weights: models.DefaultDimensionWeights,
		},
	}

	analyzer := analyzer.NewAnalyzer(cfg)

	t.Run("AI样本文件", func(t *testing.T) {
		samplePath := filepath.Join(testDataDir, "sample.txt")
		data, err := os.ReadFile(samplePath)
		if err != nil {
			t.Skip("sample.txt not found, skipping")
		}

		request := models.DetectionRequest{
			Text: string(data),
		}

		result, err := analyzer.Analyze(request)
		if err != nil {
			t.Fatalf("Analyze() error = %v", err)
		}

		// AI样本应该得到较低分数
		if result.Score.Total > 60 {
			t.Errorf("AI sample score = %.1f, expected < 60", result.Score.Total)
		}
	})

	t.Run("人类样本文件", func(t *testing.T) {
		samplePath := filepath.Join(testDataDir, "human_sample.txt")
		data, err := os.ReadFile(samplePath)
		if err != nil {
			t.Skip("human_sample.txt not found, skipping")
		}

		request := models.DetectionRequest{
			Text: string(data),
		}

		result, err := analyzer.Analyze(request)
		if err != nil {
			t.Fatalf("Analyze() error = %v", err)
		}

		// 人类样本应该得到较高分数
		if result.Score.Total < 70 {
			t.Errorf("Human sample score = %.1f, expected > 70", result.Score.Total)
		}
	})
}

func TestAnalyzer_EdgeCases(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
		Scoring: config.ScoringConfig{
			Weights: models.DefaultDimensionWeights,
		},
	}

	analyzer := analyzer.NewAnalyzer(cfg)

	t.Run("空文本", func(t *testing.T) {
		request := models.DetectionRequest{
			Text: "",
		}

		result, err := analyzer.Analyze(request)
		if err != nil {
			t.Fatalf("Analyze() error = %v", err)
		}

		// 空文本应该有结果
		if result == nil {
			t.Error("Expected result, got nil")
		}
	})

	t.Run("纯标点文本", func(t *testing.T) {
		request := models.DetectionRequest{
			Text: "... !!! ???",
		}

		result, err := analyzer.Analyze(request)
		if err != nil {
			t.Fatalf("Analyze() error = %v", err)
		}

		if result == nil {
			t.Error("Expected result, got nil")
		}
	})

	t.Run("纯数字文本", func(t *testing.T) {
		request := models.DetectionRequest{
			Text: "123 456 789 000",
		}

		result, err := analyzer.Analyze(request)
		if err != nil {
			t.Fatalf("Analyze() error = %v", err)
		}

		if result == nil {
			t.Error("Expected result, got nil")
		}
	})

	t.Run("超长文本", func(t *testing.T) {
		// 生成一个较长的文本
		longText := ""
		for i := 0; i < 100; i++ {
			longText += "This is a test sentence that is repeated many times. "
		}

		request := models.DetectionRequest{
			Text: longText,
		}

		result, err := analyzer.Analyze(request)
		if err != nil {
			t.Fatalf("Analyze() error = %v", err)
		}

		if result == nil {
			t.Error("Expected result, got nil")
		}

		// 分数应该在有效范围内
		if result.Score.Total < 0 || result.Score.Total > 100 {
			t.Errorf("Score = %.1f, expected between 0 and 100", result.Score.Total)
		}
	})
}

func TestAnalyzer_RiskLevelClassification(t *testing.T) {
	tests := []struct {
		score    float64
		expected models.RiskLevel
	}{
		{0, models.RiskLevelVeryHigh},
		{20, models.RiskLevelVeryHigh},
		{40, models.RiskLevelVeryHigh},
		{41, models.RiskLevelHigh},
		{50, models.RiskLevelHigh},
		{60, models.RiskLevelHigh},
		{61, models.RiskLevelMedium},
		{70, models.RiskLevelMedium},
		{75, models.RiskLevelMedium},
		{76, models.RiskLevelLow},
		{90, models.RiskLevelLow},
		{100, models.RiskLevelLow},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := models.GetRiskLevel(tt.score)
			if result != tt.expected {
				t.Errorf("GetRiskLevel(%.1f) = %s, want %s", tt.score, result, tt.expected)
			}
		})
	}
}

func TestAnalyzer_Performance(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
		Scoring: config.ScoringConfig{
			Weights: models.DefaultDimensionWeights,
		},
	}

	analyzer := analyzer.NewAnalyzer(cfg)

	// 生成约1000字的文本
	text := ""
	for i := 0; i < 200; i++ {
		text += "This is a sample sentence with some words. "
	}

	request := models.DetectionRequest{
		Text: text,
	}

	result, err := analyzer.Analyze(request)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}

	// 处理时间应该在3秒内
	if result.ProcessTime.Seconds() > 3 {
		t.Errorf("ProcessTime = %v, expected < 3s", result.ProcessTime)
	}
}

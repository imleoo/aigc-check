package detector

import (
	"testing"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
)

// MockRule 用于测试的模拟规则
type MockRule struct {
	ruleType  models.RuleType
	detected  bool
	score     float64
	count     int
}

func (m *MockRule) Check(text string) models.RuleResult {
	return models.RuleResult{
		RuleType: m.ruleType,
		Detected: m.detected,
		Score:    m.score,
		Count:    m.count,
	}
}

func (m *MockRule) GetType() models.RuleType {
	return m.ruleType
}

func (m *MockRule) GetName() string {
	return string(m.ruleType)
}

func (m *MockRule) GetDescription() string {
	return "Mock rule for testing"
}

func TestNewRuleEngine(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	engine := NewRuleEngine(cfg)

	if engine == nil {
		t.Fatal("NewRuleEngine() returned nil")
	}

	if engine.config != cfg {
		t.Error("NewRuleEngine() config not set correctly")
	}
}

func TestRuleEngine_RegisterRule(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	engine := NewRuleEngine(cfg)

	rule := &MockRule{
		ruleType: models.RuleTypeHighFreqWords,
		detected: true,
		score:    80.0,
	}

	engine.RegisterRule(rule)

	if engine.CountRules() != 1 {
		t.Errorf("CountRules() = %d, want 1", engine.CountRules())
	}
}

func TestRuleEngine_UnregisterRule(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	engine := NewRuleEngine(cfg)

	rule := &MockRule{
		ruleType: models.RuleTypeHighFreqWords,
	}

	engine.RegisterRule(rule)
	engine.UnregisterRule(models.RuleTypeHighFreqWords)

	if engine.CountRules() != 0 {
		t.Errorf("CountRules() after unregister = %d, want 0", engine.CountRules())
	}
}

func TestRuleEngine_GetRule(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	engine := NewRuleEngine(cfg)

	rule := &MockRule{
		ruleType: models.RuleTypeHighFreqWords,
	}

	engine.RegisterRule(rule)

	// 测试存在的规则
	got, exists := engine.GetRule(models.RuleTypeHighFreqWords)
	if !exists {
		t.Error("GetRule() exists = false, want true")
	}
	if got != rule {
		t.Error("GetRule() returned different rule")
	}

	// 测试不存在的规则
	_, exists = engine.GetRule(models.RuleTypeSentenceStarters)
	if exists {
		t.Error("GetRule() for non-existent rule exists = true, want false")
	}
}

func TestRuleEngine_GetAllRules(t *testing.T) {
	cfg := &config.Config{Thresholds: config.DefaultThresholds}
	engine := NewRuleEngine(cfg)

	rule1 := &MockRule{ruleType: models.RuleTypeHighFreqWords}
	rule2 := &MockRule{ruleType: models.RuleTypeSentenceStarters}

	engine.RegisterRule(rule1)
	engine.RegisterRule(rule2)

	rules := engine.GetAllRules()
	if len(rules) != 2 {
		t.Errorf("GetAllRules() returned %d rules, want 2", len(rules))
	}
}

func TestRuleEngine_Check(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
		Rules: map[string]config.RuleConfig{
			string(models.RuleTypeHighFreqWords): {Enabled: true},
			string(models.RuleTypeSentenceStarters): {Enabled: true},
		},
	}
	engine := NewRuleEngine(cfg)

	rule1 := &MockRule{
		ruleType: models.RuleTypeHighFreqWords,
		detected: true,
		score:    70.0,
		count:    5,
	}
	rule2 := &MockRule{
		ruleType: models.RuleTypeSentenceStarters,
		detected: false,
		score:    100.0,
		count:    0,
	}

	engine.RegisterRule(rule1)
	engine.RegisterRule(rule2)

	results := engine.Check("test text")

	if len(results) != 2 {
		t.Errorf("Check() returned %d results, want 2", len(results))
	}

	// 验证结果包含预期的规则类型
	hasHighFreq := false
	hasSentenceStarters := false
	for _, r := range results {
		if r.RuleType == models.RuleTypeHighFreqWords {
			hasHighFreq = true
			if !r.Detected {
				t.Error("HighFreqWords result Detected = false, want true")
			}
		}
		if r.RuleType == models.RuleTypeSentenceStarters {
			hasSentenceStarters = true
			if r.Detected {
				t.Error("SentenceStarters result Detected = true, want false")
			}
		}
	}

	if !hasHighFreq {
		t.Error("Check() missing HighFreqWords result")
	}
	if !hasSentenceStarters {
		t.Error("Check() missing SentenceStarters result")
	}
}

func TestRuleEngine_Check_DisabledRule(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
		Rules: map[string]config.RuleConfig{
			string(models.RuleTypeHighFreqWords): {Enabled: false},
		},
	}
	engine := NewRuleEngine(cfg)

	rule := &MockRule{
		ruleType: models.RuleTypeHighFreqWords,
		detected: true,
	}

	engine.RegisterRule(rule)

	results := engine.Check("test text")

	// 禁用的规则不应该执行
	if len(results) != 0 {
		t.Errorf("Check() with disabled rule returned %d results, want 0", len(results))
	}
}

func TestRuleEngine_CheckWithRules(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
		Rules: map[string]config.RuleConfig{
			string(models.RuleTypeHighFreqWords): {Enabled: true},
			string(models.RuleTypeSentenceStarters): {Enabled: true},
			string(models.RuleTypeEmDash): {Enabled: true},
		},
	}
	engine := NewRuleEngine(cfg)

	rule1 := &MockRule{ruleType: models.RuleTypeHighFreqWords, detected: true}
	rule2 := &MockRule{ruleType: models.RuleTypeSentenceStarters, detected: true}
	rule3 := &MockRule{ruleType: models.RuleTypeEmDash, detected: true}

	engine.RegisterRule(rule1)
	engine.RegisterRule(rule2)
	engine.RegisterRule(rule3)

	// 只执行指定的规则
	results := engine.CheckWithRules("test", []models.RuleType{
		models.RuleTypeHighFreqWords,
		models.RuleTypeSentenceStarters,
	})

	if len(results) != 2 {
		t.Errorf("CheckWithRules() returned %d results, want 2", len(results))
	}
}

func TestRuleEngine_GetEnabledRules(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
		Rules: map[string]config.RuleConfig{
			string(models.RuleTypeHighFreqWords): {Enabled: true},
			string(models.RuleTypeSentenceStarters): {Enabled: false},
		},
	}
	engine := NewRuleEngine(cfg)

	rule1 := &MockRule{ruleType: models.RuleTypeHighFreqWords}
	rule2 := &MockRule{ruleType: models.RuleTypeSentenceStarters}

	engine.RegisterRule(rule1)
	engine.RegisterRule(rule2)

	enabled := engine.GetEnabledRules()
	if len(enabled) != 1 {
		t.Errorf("GetEnabledRules() returned %d rules, want 1", len(enabled))
	}
}

func TestRuleEngine_CountEnabledRules(t *testing.T) {
	cfg := &config.Config{
		Thresholds: config.DefaultThresholds,
		Rules: map[string]config.RuleConfig{
			string(models.RuleTypeHighFreqWords): {Enabled: true},
			string(models.RuleTypeSentenceStarters): {Enabled: false},
			string(models.RuleTypeEmDash): {Enabled: true},
		},
	}
	engine := NewRuleEngine(cfg)

	engine.RegisterRule(&MockRule{ruleType: models.RuleTypeHighFreqWords})
	engine.RegisterRule(&MockRule{ruleType: models.RuleTypeSentenceStarters})
	engine.RegisterRule(&MockRule{ruleType: models.RuleTypeEmDash})

	count := engine.CountEnabledRules()
	if count != 2 {
		t.Errorf("CountEnabledRules() = %d, want 2", count)
	}
}

package detector

import (
	"sync"

	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
)

// RuleEngine 规则引擎
type RuleEngine struct {
	rules  map[models.RuleType]models.Rule
	config *config.Config
	mu     sync.RWMutex
}

// NewRuleEngine 创建规则引擎
func NewRuleEngine(cfg *config.Config) *RuleEngine {
	return &RuleEngine{
		rules:  make(map[models.RuleType]models.Rule),
		config: cfg,
	}
}

// RegisterRule 注册规则
func (e *RuleEngine) RegisterRule(rule models.Rule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules[rule.GetType()] = rule
}

// UnregisterRule 注销规则
func (e *RuleEngine) UnregisterRule(ruleType models.RuleType) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.rules, ruleType)
}

// GetRule 获取规则
func (e *RuleEngine) GetRule(ruleType models.RuleType) (models.Rule, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	rule, exists := e.rules[ruleType]
	return rule, exists
}

// GetAllRules 获取所有规则
func (e *RuleEngine) GetAllRules() []models.Rule {
	e.mu.RLock()
	defer e.mu.RUnlock()

	rules := make([]models.Rule, 0, len(e.rules))
	for _, rule := range e.rules {
		rules = append(rules, rule)
	}
	return rules
}

// Check 执行所有启用的规则检测
func (e *RuleEngine) Check(text string) []models.RuleResult {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var results []models.RuleResult

	// 并发执行规则检测
	resultChan := make(chan models.RuleResult, len(e.rules))
	var wg sync.WaitGroup

	for ruleType, rule := range e.rules {
		// 检查规则是否启用
		if !e.config.IsRuleEnabled(ruleType) {
			continue
		}

		wg.Add(1)
		go func(r models.Rule) {
			defer wg.Done()
			result := r.Check(text)
			resultChan <- result
		}(rule)
	}

	// 等待所有规则执行完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// CheckWithRules 执行指定规则检测
func (e *RuleEngine) CheckWithRules(text string, ruleTypes []models.RuleType) []models.RuleResult {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var results []models.RuleResult

	// 并发执行规则检测
	resultChan := make(chan models.RuleResult, len(ruleTypes))
	var wg sync.WaitGroup

	for _, ruleType := range ruleTypes {
		rule, exists := e.rules[ruleType]
		if !exists {
			continue
		}

		// 检查规则是否启用
		if !e.config.IsRuleEnabled(ruleType) {
			continue
		}

		wg.Add(1)
		go func(r models.Rule) {
			defer wg.Done()
			result := r.Check(text)
			resultChan <- result
		}(rule)
	}

	// 等待所有规则执行完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// GetEnabledRules 获取所有启用的规则
func (e *RuleEngine) GetEnabledRules() []models.Rule {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var enabledRules []models.Rule
	for ruleType, rule := range e.rules {
		if e.config.IsRuleEnabled(ruleType) {
			enabledRules = append(enabledRules, rule)
		}
	}
	return enabledRules
}

// CountRules 统计规则数量
func (e *RuleEngine) CountRules() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.rules)
}

// CountEnabledRules 统计启用的规则数量
func (e *RuleEngine) CountEnabledRules() int {
	e.mu.RLock()
	defer e.mu.RUnlock()

	count := 0
	for ruleType := range e.rules {
		if e.config.IsRuleEnabled(ruleType) {
			count++
		}
	}
	return count
}

package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/leoobai/aigc-check/internal/analyzer"
	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
	"github.com/leoobai/aigc-check/internal/repository"
)

// DetectionService 检测服务接口
type DetectionService interface {
	Detect(text string, options DetectionOptions) (*DetectionResult, error)
	GetResult(id string) (*DetectionResult, error)
}

// DetectionOptions 检测选项
type DetectionOptions struct {
	EnableMultimodal bool
	EnableStatistics bool
	EnableSemantic   bool
	Language         string
}

// DetectionResult 检测结果
type DetectionResult struct {
	ID               string                  `json:"id"`
	RequestID        string                  `json:"request_id"`
	Text             string                  `json:"text"`
	Score            *models.Score           `json:"score"`
	RiskLevel        string                  `json:"risk_level"`
	RuleResults      []*models.RuleResult    `json:"rule_results"`
	Suggestions      []*models.Suggestion    `json:"suggestions"`
	MultimodalResult *models.MultimodalResult `json:"multimodal,omitempty"`
	ProcessTime      string                  `json:"process_time"`
	DetectedAt       time.Time               `json:"detected_at"`
}

// detectionService 检测服务实现
type detectionService struct {
	analyzer   *analyzer.Analyzer
	repository repository.DetectionRepository
}

// NewDetectionService 创建检测服务
func NewDetectionService(cfg *config.Config, repo repository.DetectionRepository) DetectionService {
	return &detectionService{
		analyzer:   analyzer.NewAnalyzer(cfg),
		repository: repo,
	}
}

// Detect 执行文本检测
func (s *detectionService) Detect(text string, options DetectionOptions) (*DetectionResult, error) {
	// 构建检测请求
	request := models.DetectionRequest{
		Text: text,
		Options: models.DetectionOptions{
			Language: options.Language,
		},
	}

	// 执行分析
	result, err := s.analyzer.Analyze(request)
	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	// 生成 ID
	id := uuid.New().String()

	// 转换结果
	detectionResult := &DetectionResult{
		ID:          id,
		RequestID:   result.RequestID,
		Text:        result.Text,
		Score:       &result.Score,
		RiskLevel:   string(result.RiskLevel),
		ProcessTime: result.ProcessTime.String(),
		DetectedAt:  result.DetectedAt,
	}

	// 转换 RuleResults
	detectionResult.RuleResults = make([]*models.RuleResult, len(result.RuleResults))
	for i := range result.RuleResults {
		detectionResult.RuleResults[i] = &result.RuleResults[i]
	}

	// 转换 Suggestions
	detectionResult.Suggestions = make([]*models.Suggestion, len(result.Suggestions))
	for i := range result.Suggestions {
		detectionResult.Suggestions[i] = &result.Suggestions[i]
	}

	// 保存到数据库
	if err := s.saveToRepository(detectionResult); err != nil {
		return nil, fmt.Errorf("failed to save result: %w", err)
	}

	return detectionResult, nil
}

// saveToRepository 保存检测结果到数据库
func (s *detectionService) saveToRepository(result *DetectionResult) error {
	// 序列化 JSON 字段
	ruleResultsJSON, err := json.Marshal(result.RuleResults)
	if err != nil {
		return fmt.Errorf("failed to marshal rule results: %w", err)
	}

	suggestionsJSON, err := json.Marshal(result.Suggestions)
	if err != nil {
		return fmt.Errorf("failed to marshal suggestions: %w", err)
	}

	var multimodalJSON []byte
	if result.MultimodalResult != nil {
		multimodalJSON, err = json.Marshal(result.MultimodalResult)
		if err != nil {
			return fmt.Errorf("failed to marshal multimodal result: %w", err)
		}
	}

	// 创建文本预览（前100字）
	textPreview := result.Text
	if len(textPreview) > 100 {
		textPreview = textPreview[:100] + "..."
	}

	// 创建数据库记录
	record := &repository.DetectionRecord{
		ID:               result.ID,
		RequestID:        result.RequestID,
		Text:             result.Text,
		TextPreview:      textPreview,
		Score:            result.Score.Total,
		RiskLevel:        result.RiskLevel,
		RuleResults:      string(ruleResultsJSON),
		Suggestions:      string(suggestionsJSON),
		MultimodalResult: string(multimodalJSON),
		ProcessTime:      result.ProcessTime,
	}

	return s.repository.Create(record)
}

// GetResult 根据 ID 获取检测结果
func (s *detectionService) GetResult(id string) (*DetectionResult, error) {
	record, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 反序列化 JSON 字段
	var ruleResults []*models.RuleResult
	if err := json.Unmarshal([]byte(record.RuleResults), &ruleResults); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rule results: %w", err)
	}

	var suggestions []*models.Suggestion
	if err := json.Unmarshal([]byte(record.Suggestions), &suggestions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal suggestions: %w", err)
	}

	var multimodalResult *models.MultimodalResult
	if record.MultimodalResult != "" {
		if err := json.Unmarshal([]byte(record.MultimodalResult), &multimodalResult); err != nil {
			return nil, fmt.Errorf("failed to unmarshal multimodal result: %w", err)
		}
	}

	return &DetectionResult{
		ID:               record.ID,
		RequestID:        record.RequestID,
		Text:             record.Text,
		Score:            &models.Score{Total: record.Score},
		RiskLevel:        record.RiskLevel,
		RuleResults:      ruleResults,
		Suggestions:      suggestions,
		MultimodalResult: multimodalResult,
		ProcessTime:      record.ProcessTime,
		DetectedAt:       record.CreatedAt,
	}, nil
}

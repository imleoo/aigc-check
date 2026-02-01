package service

import (
	"encoding/json"
	"fmt"

	"github.com/leoobai/aigc-check/internal/models"
	"github.com/leoobai/aigc-check/internal/repository"
)

// HistoryService 历史服务接口
type HistoryService interface {
	List(page, pageSize int, sortBy, order string) (*HistoryListResult, error)
	GetByID(id string) (*DetectionResult, error)
	Delete(id string) error
	DeleteAll() error
}

// HistoryListResult 历史列表结果
type HistoryListResult struct {
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
	Items    []*HistoryListItem `json:"items"`
}

// HistoryListItem 历史列表项
type HistoryListItem struct {
	ID          string `json:"id"`
	RequestID   string `json:"request_id"`
	TextPreview string `json:"text_preview"`
	Score       float64 `json:"score"`
	RiskLevel   string `json:"risk_level"`
	CreatedAt   string `json:"created_at"`
}

// historyService 历史服务实现
type historyService struct {
	repository repository.DetectionRepository
}

// NewHistoryService 创建历史服务
func NewHistoryService(repo repository.DetectionRepository) HistoryService {
	return &historyService{
		repository: repo,
	}
}

// List 获取历史记录列表
func (s *historyService) List(page, pageSize int, sortBy, order string) (*HistoryListResult, error) {
	records, total, err := s.repository.List(page, pageSize, sortBy, order)
	if err != nil {
		return nil, fmt.Errorf("failed to list history: %w", err)
	}

	items := make([]*HistoryListItem, len(records))
	for i, record := range records {
		items[i] = &HistoryListItem{
			ID:          record.ID,
			RequestID:   record.RequestID,
			TextPreview: record.TextPreview,
			Score:       record.Score,
			RiskLevel:   record.RiskLevel,
			CreatedAt:   record.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &HistoryListResult{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    items,
	}, nil
}

// GetByID 根据 ID 获取历史记录
func (s *historyService) GetByID(id string) (*DetectionResult, error) {
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

	return &DetectionResult{
		ID:          record.ID,
		RequestID:   record.RequestID,
		Text:        record.Text,
		Score:       &models.Score{Total: record.Score},
		RiskLevel:   record.RiskLevel,
		RuleResults: ruleResults,
		Suggestions: suggestions,
		ProcessTime: record.ProcessTime,
		DetectedAt:  record.CreatedAt,
	}, nil
}

// Delete 删除历史记录
func (s *historyService) Delete(id string) error {
	return s.repository.Delete(id)
}

// DeleteAll 删除所有历史记录
func (s *historyService) DeleteAll() error {
	return s.repository.DeleteAll()
}

package repository

import (
	"fmt"

	"gorm.io/gorm"
)

// DetectionRepository 检测记录仓储接口
type DetectionRepository interface {
	Create(record *DetectionRecord) error
	GetByID(id string) (*DetectionRecord, error)
	GetByRequestID(requestID string) (*DetectionRecord, error)
	List(page, pageSize int, sortBy, order string) ([]*DetectionRecord, int64, error)
	Delete(id string) error
	DeleteAll() error
}

// detectionRepository 检测记录仓储实现
type detectionRepository struct {
	db *gorm.DB
}

// NewDetectionRepository 创建检测记录仓储
func NewDetectionRepository(db *gorm.DB) DetectionRepository {
	return &detectionRepository{db: db}
}

// Create 创建检测记录
func (r *detectionRepository) Create(record *DetectionRecord) error {
	if err := r.db.Create(record).Error; err != nil {
		return fmt.Errorf("failed to create detection record: %w", err)
	}
	return nil
}

// GetByID 根据 ID 获取检测记录
func (r *detectionRepository) GetByID(id string) (*DetectionRecord, error) {
	var record DetectionRecord
	if err := r.db.Where("id = ?", id).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("detection record not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get detection record: %w", err)
	}
	return &record, nil
}

// GetByRequestID 根据 RequestID 获取检测记录
func (r *detectionRepository) GetByRequestID(requestID string) (*DetectionRecord, error) {
	var record DetectionRecord
	if err := r.db.Where("request_id = ?", requestID).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("detection record not found: %s", requestID)
		}
		return nil, fmt.Errorf("failed to get detection record: %w", err)
	}
	return &record, nil
}

// List 获取检测记录列表（分页）
func (r *detectionRepository) List(page, pageSize int, sortBy, order string) ([]*DetectionRecord, int64, error) {
	var records []*DetectionRecord
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询
	query := r.db.Model(&DetectionRecord{})

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count records: %w", err)
	}

	// 排序和分页
	orderClause := fmt.Sprintf("%s %s", sortBy, order)
	if err := query.Order(orderClause).Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list records: %w", err)
	}

	return records, total, nil
}

// Delete 删除检测记录
func (r *detectionRepository) Delete(id string) error {
	result := r.db.Where("id = ?", id).Delete(&DetectionRecord{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete detection record: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("detection record not found: %s", id)
	}
	return nil
}

// DeleteAll 删除所有检测记录
func (r *detectionRepository) DeleteAll() error {
	if err := r.db.Exec("DELETE FROM detection_records").Error; err != nil {
		return fmt.Errorf("failed to delete all records: %w", err)
	}
	return nil
}

package repository

import (
	"time"

	"gorm.io/gorm"
)

// DetectionRecord 检测记录数据库模型
type DetectionRecord struct {
	ID               string    `gorm:"primaryKey;type:text"`
	RequestID        string    `gorm:"uniqueIndex;type:text;not null"`
	Text             string    `gorm:"type:text;not null"`
	TextPreview      string    `gorm:"type:text"`
	Score            float64   `gorm:"not null"`
	RiskLevel        string    `gorm:"type:text;not null;index"`
	RuleResults      string    `gorm:"type:text"` // JSON
	Suggestions      string    `gorm:"type:text"` // JSON
	MultimodalResult string    `gorm:"type:text"` // JSON
	ProcessTime      string    `gorm:"type:text"`
	CreatedAt        time.Time `gorm:"index"`
	UpdatedAt        time.Time
}

// TableName 指定表名
func (DetectionRecord) TableName() string {
	return "detection_records"
}

// BeforeCreate GORM 钩子：创建前设置时间
func (r *DetectionRecord) BeforeCreate(tx *gorm.DB) error {
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now()
	}
	if r.UpdatedAt.IsZero() {
		r.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate GORM 钩子：更新前设置时间
func (r *DetectionRecord) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedAt = time.Now()
	return nil
}

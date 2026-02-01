package migrations

import (
	"fmt"
	"log"

	"github.com/leoobai/aigc-check/internal/repository"
	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表
func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database migration...")

	// 自动迁移所有模型
	if err := db.AutoMigrate(
		&repository.DetectionRecord{},
	); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

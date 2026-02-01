package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Config 数据库配置
type Config struct {
	Type     string // sqlite, postgres, mysql
	DSN      string // 数据源名称
	LogLevel logger.LogLevel
}

// Initialize 初始化数据库连接
func Initialize(config Config) error {
	var dialector gorm.Dialector

	switch config.Type {
	case "sqlite":
		// 确保数据目录存在
		dbDir := filepath.Dir(config.DSN)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}
		dialector = sqlite.Open(config.DSN)
	default:
		return fmt.Errorf("unsupported database type: %s", config.Type)
	}

	// 配置 GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	}

	// 连接数据库
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	log.Printf("Database connected successfully (type: %s)", config.Type)

	return nil
}

// Close 关闭数据库连接
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	log.Println("Database connection closed")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

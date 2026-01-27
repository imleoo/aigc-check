package reporter

import "github.com/leoobai/aigc-check/internal/models"

// Reporter 报告生成器接口
type Reporter interface {
	// Generate 生成报告
	Generate(result *models.DetectionResult) (string, error)

	// Format 获取报告格式
	Format() string
}

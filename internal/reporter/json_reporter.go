package reporter

import (
	"encoding/json"

	"github.com/leoobai/aigc-check/internal/models"
)

// JSONReporter JSON报告生成器
type JSONReporter struct {
	pretty bool
}

// NewJSONReporter 创建JSON报告生成器
func NewJSONReporter(pretty bool) *JSONReporter {
	return &JSONReporter{
		pretty: pretty,
	}
}

// Generate 生成JSON报告
func (r *JSONReporter) Generate(result *models.DetectionResult) (string, error) {
	var data []byte
	var err error

	if r.pretty {
		data, err = json.MarshalIndent(result, "", "  ")
	} else {
		data, err = json.Marshal(result)
	}

	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Format 获取报告格式
func (r *JSONReporter) Format() string {
	return "json"
}

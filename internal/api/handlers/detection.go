package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leoobai/aigc-check/internal/service"
)

// DetectionHandler 检测处理器
type DetectionHandler struct {
	detectionService service.DetectionService
}

// NewDetectionHandler 创建检测处理器
func NewDetectionHandler(detectionService service.DetectionService) *DetectionHandler {
	return &DetectionHandler{
		detectionService: detectionService,
	}
}

// DetectRequest 检测请求
// @Description 检测请求参数
type DetectRequest struct {
	Text    string        `json:"text" binding:"required" example:"这是一段需要检测的文本"`
	Options DetectOptions `json:"options"`
}

// DetectOptions 检测选项
// @Description 检测选项配置
type DetectOptions struct {
	EnableMultimodal bool   `json:"enable_multimodal" example:"false"`
	EnableStatistics bool   `json:"enable_statistics" example:"false"`
	EnableSemantic   bool   `json:"enable_semantic" example:"false"`
	Language         string `json:"language" example:"zh"`
}

// Response 通用响应
// @Description API 通用响应格式
type Response struct {
	Code    int         `json:"code" example:"0"`
	Message string      `json:"message" example:"success"`
	Data    interface{} `json:"data,omitempty"`
}

// DetectionResultResponse 检测结果响应
// @Description 检测结果数据
type DetectionResultResponse struct {
	ID          string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	RequestID   string  `json:"request_id" example:"req-12345"`
	Text        string  `json:"text" example:"检测的原始文本"`
	Score       float64 `json:"score" example:"75.5"`
	RiskLevel   string  `json:"risk_level" example:"medium"`
	ProcessTime string  `json:"process_time" example:"150ms"`
	DetectedAt  string  `json:"detected_at" example:"2024-01-15T10:30:00Z"`
}

// Detect 执行检测
// @Summary      执行AI内容检测
// @Description  对输入文本进行AI生成内容检测，返回检测分数和风险等级
// @Tags         detection
// @Accept       json
// @Produce      json
// @Param        request body DetectRequest true "检测请求参数"
// @Success      200 {object} Response{data=DetectionResultResponse} "检测成功"
// @Failure      400 {object} Response "请求参数错误"
// @Failure      500 {object} Response "服务器内部错误"
// @Router       /api/v1/detect [post]
func (h *DetectionHandler) Detect(c *gin.Context) {
	var req DetectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// 转换选项
	options := service.DetectionOptions{
		EnableMultimodal: req.Options.EnableMultimodal,
		EnableStatistics: req.Options.EnableStatistics,
		EnableSemantic:   req.Options.EnableSemantic,
		Language:         req.Options.Language,
	}

	// 执行检测
	result, err := h.detectionService.Detect(req.Text, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Detection failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    result,
	})
}

// GetByID 根据 ID 获取检测结果
// @Summary      获取检测结果详情
// @Description  根据检测ID获取详细的检测结果
// @Tags         detection
// @Accept       json
// @Produce      json
// @Param        id path string true "检测结果ID"
// @Success      200 {object} Response{data=DetectionResultResponse} "获取成功"
// @Failure      400 {object} Response "请求参数错误"
// @Failure      404 {object} Response "结果不存在"
// @Router       /api/v1/detect/{id} [get]
func (h *DetectionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "ID is required",
		})
		return
	}

	result, err := h.detectionService.GetResult(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "Result not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    result,
	})
}

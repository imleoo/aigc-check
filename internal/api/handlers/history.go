package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leoobai/aigc-check/internal/service"
)

// HistoryHandler 历史记录处理器
type HistoryHandler struct {
	historyService service.HistoryService
}

// NewHistoryHandler 创建历史记录处理器
func NewHistoryHandler(historyService service.HistoryService) *HistoryHandler {
	return &HistoryHandler{
		historyService: historyService,
	}
}

// HistoryListResponse 历史列表响应
// @Description 历史记录列表
type HistoryListResponse struct {
	Total    int64               `json:"total" example:"100"`
	Page     int                 `json:"page" example:"1"`
	PageSize int                 `json:"page_size" example:"20"`
	Items    []HistoryItemResult `json:"items"`
}

// HistoryItemResult 历史记录项
// @Description 历史记录条目
type HistoryItemResult struct {
	ID          string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	RequestID   string  `json:"request_id" example:"req-12345"`
	TextPreview string  `json:"text_preview" example:"这是检测文本的预览..."`
	Score       float64 `json:"score" example:"75.5"`
	RiskLevel   string  `json:"risk_level" example:"medium"`
	CreatedAt   string  `json:"created_at" example:"2024-01-15 10:30:00"`
}

// allowedSortColumns 允许的排序字段白名单
var allowedSortColumns = map[string]bool{
	"created_at": true,
	"score":      true,
	"risk_level": true,
}

// allowedOrders 允许的排序方向
var allowedOrders = map[string]bool{
	"asc":  true,
	"desc": true,
}

// List 获取历史记录列表
// @Summary      获取历史记录列表
// @Description  分页获取检测历史记录列表
// @Tags         history
// @Accept       json
// @Produce      json
// @Param        page query int false "页码" default(1) minimum(1)
// @Param        page_size query int false "每页数量" default(20) minimum(1) maximum(100)
// @Param        sort query string false "排序字段" Enums(created_at,score,risk_level) default(created_at)
// @Param        order query string false "排序方向" Enums(asc,desc) default(desc)
// @Success      200 {object} Response{data=HistoryListResponse} "获取成功"
// @Failure      400 {object} Response "请求参数错误"
// @Failure      500 {object} Response "服务器内部错误"
// @Router       /api/v1/history [get]
func (h *HistoryHandler) List(c *gin.Context) {
	// 解析分页参数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 验证排序字段（防止 SQL 注入）
	sortBy := c.DefaultQuery("sort", "created_at")
	if !allowedSortColumns[sortBy] {
		sortBy = "created_at"
	}

	// 验证排序方向
	order := c.DefaultQuery("order", "desc")
	if !allowedOrders[order] {
		order = "desc"
	}

	// 获取历史记录
	result, err := h.historyService.List(page, pageSize, sortBy, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to get history: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    result,
	})
}

// GetByID 根据 ID 获取历史记录
// @Summary      获取历史记录详情
// @Description  根据ID获取单条历史记录的详细信息
// @Tags         history
// @Accept       json
// @Produce      json
// @Param        id path string true "历史记录ID"
// @Success      200 {object} Response{data=DetectionResultResponse} "获取成功"
// @Failure      400 {object} Response "请求参数错误"
// @Failure      404 {object} Response "记录不存在"
// @Router       /api/v1/history/{id} [get]
func (h *HistoryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "ID is required",
		})
		return
	}

	result, err := h.historyService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "History not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    result,
	})
}

// Delete 删除历史记录
// @Summary      删除历史记录
// @Description  根据ID删除单条历史记录
// @Tags         history
// @Accept       json
// @Produce      json
// @Param        id path string true "历史记录ID"
// @Success      200 {object} Response "删除成功"
// @Failure      400 {object} Response "请求参数错误"
// @Failure      500 {object} Response "服务器内部错误"
// @Router       /api/v1/history/{id} [delete]
func (h *HistoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "ID is required",
		})
		return
	}

	if err := h.historyService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to delete history: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
	})
}

// DeleteAll 删除所有历史记录
// @Summary      删除所有历史记录
// @Description  清空所有检测历史记录
// @Tags         history
// @Accept       json
// @Produce      json
// @Success      200 {object} Response "删除成功"
// @Failure      500 {object} Response "服务器内部错误"
// @Router       /api/v1/history [delete]
func (h *HistoryHandler) DeleteAll(c *gin.Context) {
	if err := h.historyService.DeleteAll(); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to delete all history: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
	})
}

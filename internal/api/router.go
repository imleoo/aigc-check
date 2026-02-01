package api

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/leoobai/aigc-check/internal/api/handlers"
	"github.com/leoobai/aigc-check/internal/api/middleware"

	_ "github.com/leoobai/aigc-check/docs"
)

// SetupRouter 设置路由
func SetupRouter(
	detectionHandler *handlers.DetectionHandler,
	historyHandler *handlers.HistoryHandler,
) *gin.Engine {
	router := gin.New()

	// 使用中间件
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	// Swagger API 文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 静态文件服务 (用于 Docker 部署)
	// 检查 ./static 目录是否存在，如果存在则提供静态文件服务
	if _, err := os.Stat("./static"); err == nil {
		// 服务静态资源目录
		router.Static("/assets", "./static/assets")

		// 处理根路径和其他前端路由
		router.NoRoute(func(c *gin.Context) {
			// 如果是 API 路径但未匹配，返回 404 JSON
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(404, gin.H{"error": "API not found"})
				return
			}

			// 如果是 Swagger 路径，不处理
			if strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
				return
			}

			// 检查请求的文件是否存在于 static 根目录（如 vite.svg, favicon.ico 等）
			filepath := "./static" + c.Request.URL.Path
			if _, err := os.Stat(filepath); err == nil {
				c.File(filepath)
				return
			}

			// 默认返回 index.html (SPA 支持)
			c.File("./static/index.html")
		})
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		// 检测相关 API
		v1.POST("/detect", detectionHandler.Detect)
		v1.GET("/detect/:id", detectionHandler.GetByID)

		// 历史记录相关 API
		v1.GET("/history", historyHandler.List)
		v1.GET("/history/:id", historyHandler.GetByID)
		v1.DELETE("/history/:id", historyHandler.Delete)
		v1.DELETE("/history", historyHandler.DeleteAll)
	}

	return router
}

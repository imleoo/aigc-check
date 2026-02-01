package middleware

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS 返回 CORS 中间件
func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()

	// 从环境变量读取允许的源
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins != "" {
		// 支持逗号分隔的多个源
		origins := strings.Split(allowedOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
		config.AllowOrigins = origins
		config.AllowCredentials = true
	} else {
		// 生产环境默认允许所有源（通过反向代理控制）
		config.AllowAllOrigins = true
		config.AllowCredentials = false
	}

	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}

	return cors.New(config)
}

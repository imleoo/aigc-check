package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 返回日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(startTime)

		// 记录日志
		log.Printf("[%s] %s %s - %d (%v)",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			duration,
		)
	}
}

package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rhoat/go-exercise/pkg/system"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logging(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Set("start_time", start)

		c.Next() // Process request
		startTime := c.GetTime("start_time")

		// Log after request is processed
		latency := time.Since(startTime)
		fields := []zapcore.Field{
			zap.Int("status", c.Writer.Status()),              // Status code after processing
			zap.String("method", c.Request.Method),            // HTTP method
			zap.String("ApplicationID", system.ApplicationID), // ApplicationID
			zap.String("BuildVersion", system.BuildVersion),   // BuildVersion
			zap.String("path", c.Request.URL.Path),            // Request path
			zap.String("query", c.Request.URL.RawQuery),       // Query string
			zap.String("ip", c.ClientIP()),                    // Client IP
			zap.String("user-agent", c.Request.UserAgent()),   // User-Agent header
			zap.Int64("usLatency", latency.Microseconds()),    // Latency
			zap.Int64("msLatency", latency.Milliseconds()),    // Latency
		}
		log.Info("Request completed", fields...)
	}
}

package middlewares

import (
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/pkg/constants"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"github.com/aditya109/amrutha_assignment/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func JSONLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.Request.URL.Path, "/healthz") && !strings.Contains(c.Request.URL.Path, "/health") {
			start := time.Now()
			c.Next()
			responseTime := time.Since(start)
			contextLogger := logger.GetContextLoggerForGin(
				logger.ServerHTTPLog{
					Time:          fmt.Sprintf("%v", time.Now().UTC()),
					Hostname:      c.Request.Host,
					RemoteAddress: c.ClientIP(),
					Method:        c.Request.Method,
					Path:          c.Request.URL.Path,
					ResponseTime:  responseTime.String(),
					StatusCode:    c.Writer.Status(),
					UserAgent:     c.Request.UserAgent(),
					AppVersion:    c.Keys["appVersion"],
					TraceId:       context.GetTraceId(c),
				},
			)
			if c.Writer.Status() >= 500 {
				contextLogger.Error(c.Errors.String())
			} else {
				contextLogger.Info("")
			}
		}
	}
}

func BindTraceIdToRequestHeaderMiddleware(operationName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var shouldManuallyGenerateTraceId = viper.GetString("SHOULD_MANUALLY_GENERATED_TRACE_ID")
		if _, exists := c.Keys[constants.ApplicationTraceKey]; !exists && shouldManuallyGenerateTraceId == "true" {
			UUID := uuid.New()
			c.Request.Header.Set(constants.ApplicationTraceKey, UUID.String())
		}
		c.Next()
	}
}

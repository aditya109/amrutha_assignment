package middlewares

import (
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/pkg/constants"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"github.com/aditya109/amrutha_assignment/pkg/logger"
	"github.com/aditya109/amrutha_assignment/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"net/http"
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

func MethodNotAllowedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, models.ErrorResponse(models.Error{
			Code:              "METHOD_NOT_ALLOWED",
			Message:           constants.METHOD_NOT_ALLOWED_MESSAGE,
			ResolutionMessage: "Please check your request method",
			Data:              nil,
		}))
	}
}
func ResourceNotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, models.ErrorResponse(models.Error{
			Code:              "RESOURCE_NOT_FOUND",
			Message:           constants.RESOURCE_NOT_FOUND_MESSAGE,
			ResolutionMessage: "Please check your resource address",
			Data:              nil,
		}))
	}
}

func BindTraceIdToRequestHeaderMiddleware(operationName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var shouldManuallyGenerateTraceId = viper.GetString("SHOULD_MANUALLY_GENERATED_TRACE_ID")
		if _, exists := c.Keys[constants.APPLICATION_TRACE_KEY]; !exists && shouldManuallyGenerateTraceId == "true" {
			UUID := uuid.New()
			c.Request.Header.Set(constants.APPLICATION_TRACE_KEY, UUID.String())
		}
		c.Next()
	}
}

package logger

import (
	"context"
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type LogLevel string

type ServerHTTPLog struct {
	Time          string `json:"time,omitempty"`
	Hostname      string `json:"hostname,omitempty"`
	RemoteAddress string `json:"remoteAddress,omitempty"`
	Method        string `json:"method,omitempty"`
	Path          string `json:"path,omitempty"`
	StatusCode    int    `json:"statusCode,omitempty"`
	ResponseTime  string `json:"responseTime,omitempty"`
	UserAgent     string `json:"userAgent,omitempty"`
	Platform      string `json:"platform,omitempty"`
	AppVersion    any    `json:"appVersion,omitempty"`
	TraceId       string `json:"traceId,omitempty"`
}

func GetInternalContextLogger(function interface{}, meta ...interface{}) *logrus.Entry {
	var traceId, spanId string
	if len(meta) >= 1 {
		traceId = meta[0].(string)
	}
	if len(meta) >= 2 {
		spanId = meta[1].(string)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: true,
	})

	logrus.SetOutput(os.Stdout)
	contextLogger := logrus.WithFields(logrus.Fields{
		"function_name": getFunctionName(function),
		"trace_id":      traceId,
		"span_id":       spanId,
	})
	return contextLogger
}

func GetSupportContextLogger(function interface{}) *logrus.Entry {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.SetOutput(os.Stdout)
	contextLogger := logrus.WithFields(logrus.Fields{
		"function_name": getFunctionName(function),
	})
	return contextLogger
}

func getFunctionName(i interface{}) string {
	if i != nil {
		return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	} else {
		return ""
	}
}

func GetContextLoggerForGin(data ServerHTTPLog) *logrus.Entry {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	return logrus.WithFields(logrus.Fields{
		"time":           data.Time,
		"hostname":       data.Hostname,
		"remote_address": data.RemoteAddress,
		"method":         data.Method,
		"path":           data.Path,
		"status_code":    data.StatusCode,
		"response_time":  data.ResponseTime,
		"user_agent":     data.UserAgent,
		"user":           data.UserId,
		"platform":       data.Platform,
		"app_version":    data.AppVersion,
		"trace_id":       data.TraceId,
		"span_id":        data.SpanId,
	})
}

// GormLogger struct
type GormLogger struct {
	L    *logrus.Logger
	Meta []interface{}
}

// Error implements logger.Interface.
func (g GormLogger) Error(c context.Context, msg string, data ...interface{}) {
	g.L.WithFields(logrus.Fields{
		"trace_id": g.Meta[0].(string),
	}).Warn(msg, data)
}

// Info implements logger.Interface.
func (g GormLogger) Info(c context.Context, msg string, data ...interface{}) {
	g.L.WithFields(logrus.Fields{
		"trace_id": g.Meta[0].(string),
	}).Info(msg, data)
}

// LogMode implements logger.Interface.
func (g *GormLogger) LogMode(l logger.LogLevel) logger.Interface {
	g.L.SetLevel(logLevelToLogrusLevel(l))
	return g
}

// Trace implements logger.Interface.
func (g GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	if err != nil {
		g.L.WithFields(logrus.Fields{
			"error": err,
			"rows":  rows,
		}).Error(sql)
	} else {
		g.L.WithFields(logrus.Fields{
			"rows": rows,
		}).Info(sql)
	}
}

// Warn implements logger.Interface.
func (g GormLogger) Warn(c context.Context, msg string, data ...interface{}) {
	g.L.WithFields(logrus.Fields{
		"trace_id": g.Meta[0].(string),
	}).Warn(msg, data)
}

func GetCustomGormLogger(meta ...interface{}) logger.Interface {
	gormLogger := GormLogger{
		L:    NewLogrusLogger(),
		Meta: meta,
	}
	gormLogger.LogMode(logger.Info)
	return &gormLogger
}

// logLevelToLogrusLevel converts GORM log level to Logrus log level
func logLevelToLogrusLevel(level logger.LogLevel) logrus.Level {
	switch level {
	case logger.Silent:
		return logrus.PanicLevel
	case logger.Error:
		return logrus.ErrorLevel
	case logger.Warn:
		return logrus.WarnLevel
	case logger.Info:
		return logrus.InfoLevel
	default:
		return logrus.DebugLevel
	}
}

func NewLogrusLogger() *logrus.Logger {
	return &logrus.Logger{
		Formatter: &logrus.TextFormatter{
			DisableTimestamp: false,
		},
		Out: os.Stdout,
	}
}

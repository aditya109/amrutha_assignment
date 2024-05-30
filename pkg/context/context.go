package context

import (
	"context"
	"github.com/aditya109/amrutha_assignment/billing/pkg/constants"
	"github.com/aditya109/amrutha_assignment/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Backdrop interface {
	// Error(statusCode int, err fallback.Error)
	Response(statusCode int, data interface{})
	GetLogger(function interface{}) *logrus.Entry
	GetMetaStore() *MetaStore
	GetAllHeaders() map[string]string
	SetCustomErrorMessage(msg string)
	SetCustomResolutionMessage(msg string)
	SetCustomTraceId(traceId string)
	GetDatabaseInstance() *gorm.DB
	SetDatabaseLogger()
	GetContext() context.Context
	SetDatabaseInstance(*gorm.DB)
}

func GetNewBackdrop(c context.Context) Backdrop {
	return &callContext{
		metaStore: &MetaStore{
			TraceId: bindTraceIdToContextAsValue(),
		},
		processContext: c,
	}
}

func (c callContext) SetCustomErrorMessage(msg string) {
	c.metaStore.CustomErrorMessage = msg
}

func (c callContext) SetCustomResolutionMessage(msg string) {
	c.metaStore.CustomResolutionMessage = msg
}

func (c callContext) GetAllHeaders() map[string]string {
	var fHeaders = make(map[string]string)
	fHeaders[constants.APPLICATION_TRACE_KEY] = c.metaStore.TraceId
	return fHeaders
}

func (c callContext) GetMetaStore() *MetaStore {
	return c.metaStore
}

func (c callContext) GetLogger(function interface{}) *logrus.Entry {
	return logger.GetInternalContextLogger(function, c.metaStore.TraceId)
}

func (c callContext) Response(statusCode int, data interface{}) {
	log := c.GetLogger(c.Response)
	log.WithFields(logrus.Fields{"response": data}).Info()
}

// GetDatabaseInstance implements Backdrop.
func (c *callContext) GetDatabaseInstance() *gorm.DB {
	return c.dbInstance
}
func (c callContext) GetContext() context.Context {
	return c.processContext
}
func (c callContext) GetMode() string {
	return c.mode
}
func (c callContext) SetCustomTraceId(traceId string) {
	c.metaStore.TraceId = traceId
	// c.SetDatabaseLogger()
}

// SetDatabaseInstance implements Backdrop.
func (c *callContext) SetDatabaseInstance(db *gorm.DB) {
	c.dbInstance = db
}
func (c *callContext) SetDatabaseLogger() {
	c.dbInstance.Logger = logger.GetCustomGormLogger(c.metaStore.TraceId)
}

// SetMode implements Backdrop.
func (c *callContext) SetMode(mode string) {
	c.mode = mode
}

func bindTraceIdToContextAsValue() string {
	return uuid.New().String()
}

func GetTraceId(c *gin.Context) string {
	if c != nil {
		return c.Request.Header.Get(constants.APPLICATION_TRACE_KEY)
	}
	return ""
}

package context

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type callContext struct {
	metaStore  *MetaStore
	ginContext *gin.Context
	dbInstance *gorm.DB
	mode       string
}

type MetaStore struct {
	TraceId                 string
	CustomErrorMessage      string
	CustomResolutionMessage string
	StatusCode              int
}

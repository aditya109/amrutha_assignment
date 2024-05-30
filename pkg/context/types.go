package context

import (
	"context"
	"gorm.io/gorm"
)

type callContext struct {
	metaStore      *MetaStore
	processContext context.Context
	dbInstance     *gorm.DB
	mode           string
}

type MetaStore struct {
	TraceId                 string
	CustomErrorMessage      string
	CustomResolutionMessage string
}

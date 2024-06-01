package database

import (
	"gorm.io/gorm"
)

type Config struct {
	DatabaseUrl  string
	DatabaseName string
	SchemaName   string
}

type Connector interface {
	Connect() error
	Close() error
	GetDbInstance() *gorm.DB
}

func GetDbInstance(connector Connector) (*gorm.DB, error) {
	if connector.GetDbInstance() == nil {
		if err := connector.Connect(); err != nil {
			return nil, err
		}
	}
	return connector.GetDbInstance(), nil
}

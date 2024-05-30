package database

import (
	"fmt"
	"github.com/aditya109/amrutha_assignment/pkg/logger"
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

func CloseDatabaseConnection(dbConnector Connector) error {
	log := logger.GetSupportContextLogger(CloseDatabaseConnection)
	if err := dbConnector.Close(); err != nil {
		return fmt.Errorf("error while closing connection from database, err: %v", err)
	}
	log.Println("database connection successfully closed")
	return nil
}

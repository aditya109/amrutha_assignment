package postgres

import (
	"fmt"
	"github.com/aditya109/amrutha_assignment/pkg/connectors/database"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	appLogger "github.com/aditya109/amrutha_assignment/pkg/logger"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type PostgresConstruct struct {
	Config database.Config
	dB     *gorm.DB
}

func (construct *PostgresConstruct) Connect() error {
	var db *gorm.DB
	var err error
	log := appLogger.GetSupportContextLogger(construct.Connect)

	db, err = gorm.Open(postgres.Open(construct.Config.DatabaseUrl), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: fmt.Sprintf("%s.", construct.Config.SchemaName),
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to connect to Database"))
		return err
	}

	construct.instantiateDbInstance(db)
	return nil
}

func (construct *PostgresConstruct) Close() error {
	log := appLogger.GetSupportContextLogger(construct.Connect)
	if connection, err := construct.dB.DB(); err != nil {
		log.Printf("error while getting connection from db instance, err: %v", err)
	} else {
		defer connection.Close()
	}
	return nil
}

func (construct *PostgresConstruct) instantiateDbInstance(db *gorm.DB) {
	construct.dB = db
}

func (construct *PostgresConstruct) GetDbInstance() *gorm.DB {
	return construct.dB
}

func ExecuteTransaction(b context.Backdrop, tx *gorm.DB) error {
	var txError error
	if txError = tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Error; err != nil {
			return fmt.Errorf("error while creating source event entry, err: %v", err)
		}
		return nil
	}); txError != nil {
		return txError
	}
	if b != nil {
		b.GetLogger(ExecuteTransaction).Printf("%d rows got affected", tx.RowsAffected)
	} else {
		appLogger.GetSupportContextLogger(ExecuteTransaction).Printf("%d rows got affected", tx.RowsAffected)
	}
	return nil
}

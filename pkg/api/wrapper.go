package api

import (
	"github.com/aditya109/amrutha_assignment/pkg/connectors/database"
	"github.com/aditya109/amrutha_assignment/pkg/connectors/database/postgres"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"github.com/aditya109/amrutha_assignment/pkg/helpers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
)

func WrapHighOrderControl(baseController func(b context.Backdrop)) func(c *gin.Context) {
	return func(c *gin.Context) {
		helpers.PrintLandingRequestCurl(c)
		var err error
		// establishing database connection
		var db *gorm.DB
		databaseUrl := viper.GetString("DATABASE_URL")
		databaseName := viper.GetString("DATABASE_NAME")
		var construct = postgres.Construct{
			Config: database.Config{
				DatabaseUrl:  viper.GetString("DATABASE_URL"),
				DatabaseName: viper.GetString("DATABASE_NAME"),
				SchemaName:   viper.GetString("DATABASE_SCHEMA"),
			},
		}
		if databaseUrl != "" && databaseName != "" {
			if db, err = database.GetDbInstance(&construct); err != nil {
				log.Fatal("error while getting database instance")
			}
		} else {
			log.Fatal("no database url or database name found")
		}
		baseController(context.GetNewBackdrop(c, db))
	}
}

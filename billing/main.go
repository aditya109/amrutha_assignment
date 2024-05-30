package main

import (
	"context"
	"github.com/aditya109/amrutha_assignment/billing/pkg/appinit"
	"github.com/aditya109/amrutha_assignment/billing/pkg/constants"
	"github.com/aditya109/amrutha_assignment/pkg/connectors/database"
	"github.com/aditya109/amrutha_assignment/pkg/connectors/database/postgres"
	"github.com/aditya109/amrutha_assignment/pkg/logger"
	"github.com/aditya109/amrutha_assignment/pkg/middlewares"
	"github.com/aditya109/amrutha_assignment/pkg/recovery"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	appinit.InitializeApplication(constants.SERVICE_IDENTIFIER)
}

func main() {
	log := logger.GetSupportContextLogger(main)
	var err error
	// establishing database connection
	var db *gorm.DB
	databaseUrl := viper.GetString("DATABASE_URL")
	databaseName := viper.GetString("DATABASE_NAME")
	var construct = postgres.PostgresConstruct{
		Config: database.Config{
			DatabaseUrl:  databaseUrl,
			DatabaseName: databaseName,
		},
	}
	if databaseUrl != "" && databaseName != "" {
		if db, err = database.GetDbInstance(&construct); err != nil {
			log.Fatal("error while getting database instance")
		}
	} else {
		log.Fatal("no database url or database name found")
	}
	// starting server
	srv, port := acquireHttpServer()
	PORT := viper.GetString("SERVER_PORT")

	engine := gin.New()
	engine.Use(middlewares.JSONLogMiddleware())
	engine.Use(gin.CustomRecovery(recovery.ApiPanicRecovery))
	engine.Use(gin.Recovery())
	engine.Use(middlewares.BindTraceIdToRequestHeaderMiddleware(constants.SERVICE_IDENTIFIER))
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[SERVER] stopped")

	var serverShutdownTimeout time.Duration
	if ENV, exists := os.LookupEnv("ENV"); exists {
		switch ENV {
		case "production":
			serverShutdownTimeout = 5
		case "staging":
			serverShutdownTimeout = 3
		default:
			serverShutdownTimeout = 0
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[SERVER] shutdown failed: %+v", err)
	}
	defer func() {
		err := database.CloseDatabaseConnection()
		if err != nil {
			log.Fatalf("[SERVER] failed to close database connection: %+v", err)
		}
	}()
	state := <-ctx.Done()
	log.Printf("timeout of %d seconds.: %v", state, serverShutdownTimeout)
	log.Print("[SERVER] graceful shutdown completed")
}

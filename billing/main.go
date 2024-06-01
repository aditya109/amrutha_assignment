package main

import (
	"context"
	"errors"
	"github.com/aditya109/amrutha_assignment/billing/api/v1/modules"
	"github.com/aditya109/amrutha_assignment/billing/pkg/appinit"
	"github.com/aditya109/amrutha_assignment/billing/pkg/constants"
	"github.com/aditya109/amrutha_assignment/pkg/api"
	"github.com/aditya109/amrutha_assignment/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func init() {
	appinit.InitializeApplication(constants.ServiceIdentifier)
}

func main() {
	log := logger.GetSupportContextLogger(main)

	// starting server
	srv, port := api.AcquireHttpServer(modules.GetModule(), constants.ServiceIdentifier)

	go func() {
		log.Printf("%s server running on port - %s", strings.ToLower(constants.ServiceIdentifier), port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server start failed: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[SERVER] stopped")

	var serverShutdownTimeout = 1 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[SERVER] shutdown failed: %+v", err)
	}
	//defer func() {
	//	err := database.CloseDatabaseConnection(&postgres.PostgresConstruct{
	//		Config: database.Config{
	//			DatabaseUrl:  viper.GetString("DATABASE_URL"),
	//			DatabaseName: viper.GetString("DATABASE_NAME"),
	//		},
	//	})
	//	if err != nil {
	//		log.Fatalf("[SERVER] failed to close database connection: %+v", err)
	//	}
	//}()
	_ = <-ctx.Done()
	log.Printf("timeout: %v", serverShutdownTimeout)
	log.Print("[SERVER] graceful shutdown completed")
}

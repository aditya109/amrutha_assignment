package appinit

import (
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/config"
	"github.com/aditya109/amrutha_assignment/pkg/logger"
)

func InitializeApplication(serviceTag string) {
	log := logger.GetSupportContextLogger(InitializeApplication)
	log.Infoln("initializing application:", fmt.Sprintf("%s", serviceTag))
	log.Infoln("trying to load application configuration")
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error in loading config: %v", err)
	}
	log.Infoln("application configuration loaded")
}

package loan_config_repository

import (
	"errors"
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/internal/models"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"gorm.io/gorm"
)

func FindOne(b context.Backdrop, loanConfig *models.LoanConfig) error {
	db := b.GetDatabaseInstance()
	if result := db.Where(models.LoanConfig{Id: loanConfig.Id}).First(&loanConfig); result.Error == nil {
		return nil
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error occured while looking for LoanConfig with id %v", loanConfig.Id)
	} else {
		return fmt.Errorf("loan config %v not found", loanConfig.Id)
	}
}

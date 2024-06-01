package loan_account_repository

import (
	"errors"
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/internal/models"
	"github.com/aditya109/amrutha_assignment/pkg/connectors/database/postgres"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"gorm.io/gorm"
)

func Update(b context.Backdrop, loanAccount *models.LoanAccount) error {
	db := b.GetDatabaseInstance()
	if err := postgres.ExecuteTransaction(b, db.Save(loanAccount)); err != nil {
		return fmt.Errorf("failed to update loan account: %w", err)
	}
	return nil
}

func FindOne(b context.Backdrop, loanAccount *models.LoanAccount) error {
	db := b.GetDatabaseInstance()
	if result := db.Where(models.LoanAccount{Id: loanAccount.Id}).First(&loanAccount); result.Error == nil {
		return nil
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error occured while looking for loanAccount with id %v", loanAccount.DisplayId)
	} else {
		return fmt.Errorf("loanAccount with id %v not found", loanAccount.DisplayId)
	}
}

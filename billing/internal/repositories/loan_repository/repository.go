package loan_repository

import (
	"errors"
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/internal/models"
	"github.com/aditya109/amrutha_assignment/pkg/connectors/database/postgres"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"gorm.io/gorm"
)

func FindOne(b context.Backdrop, loan *models.Loan) error {
	db := b.GetDatabaseInstance()

	if result := db.Where(models.Loan{Customer: models.Customer{DisplayId: loan.Customer.DisplayId}}).First(&loan); result.Error == nil {
		return nil
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error occured while looking for customer with id %v", loan.DisplayId)
	} else {
		return fmt.Errorf("loan with id %v not found", loan.Customer.DisplayId)
	}
}

func IfExists(b context.Backdrop, loan *models.Loan) (bool, error) {
	db := b.GetDatabaseInstance()

	if result := db.Where(models.Loan{Customer: models.Customer{DisplayId: loan.Customer.DisplayId}}).First(&loan); result.Error == nil {
		return true, nil
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("error occured while looking for active loan with customer with id %v", loan.Customer.DisplayId)
	} else {
		return false, nil
	}
}

func Update(b context.Backdrop, loan *models.Loan) error {
	db := b.GetDatabaseInstance()
	if err := postgres.ExecuteTransaction(b, db.Save(loan)); err != nil {
		return fmt.Errorf("failed to update loan: %w", err)
	}
	return nil
}

package loan_account_repository

import (
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/internal/models"
	"github.com/aditya109/amrutha_assignment/pkg/connectors/database/postgres"
	"github.com/aditya109/amrutha_assignment/pkg/context"
)

func Update(b context.Backdrop, loanAccount *models.LoanAccount) error {
	db := b.GetDatabaseInstance()
	if err := postgres.ExecuteTransaction(b, db.Save(loanAccount)); err != nil {
		return fmt.Errorf("failed to update loan account: %w", err)
	}
	return nil
}

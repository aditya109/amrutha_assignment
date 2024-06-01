package payment_repository

import (
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/internal/models"
	"github.com/aditya109/amrutha_assignment/pkg/connectors/database/postgres"
	"github.com/aditya109/amrutha_assignment/pkg/context"
)

func Update(b context.Backdrop, payment *models.Payment) error {
	db := b.GetDatabaseInstance()
	if err := postgres.ExecuteTransaction(b, db.Save(payment)); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}
	return nil
}

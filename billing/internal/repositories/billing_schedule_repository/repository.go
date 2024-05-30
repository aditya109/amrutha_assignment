package billing_schedule_repository

import (
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/internal/models"
	"github.com/aditya109/amrutha_assignment/pkg/connectors/database/postgres"
	"github.com/aditya109/amrutha_assignment/pkg/context"
)

func Update(b context.Backdrop, schedule *models.BillingSchedule) error {
	db := b.GetDatabaseInstance()
	if err := postgres.ExecuteTransaction(b, db.Save(schedule)); err != nil {
		return fmt.Errorf("failed to update billing schedule: %w", err)
	}
	return nil
}

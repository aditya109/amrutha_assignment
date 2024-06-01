package billing_schedule_repository

import (
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/internal/models"
	"github.com/aditya109/amrutha_assignment/pkg/connectors/database/postgres"
	"github.com/aditya109/amrutha_assignment/pkg/context"
)

func Save(b context.Backdrop, schedule *models.BillingSchedule) error {
	db := b.GetDatabaseInstance()
	if err := postgres.ExecuteTransaction(b, db.Save(schedule)); err != nil {
		return fmt.Errorf("failed to update billing schedule: %w", err)
	}
	return nil
}

func UpdateOnly(b context.Backdrop, schedule *models.BillingSchedule) error {
	db := b.GetDatabaseInstance()
	if err := postgres.ExecuteTransaction(b, db.Updates(schedule)); err != nil {
		return fmt.Errorf("failed to update billing schedule: %w", err)
	}
	return nil
}
func FindAllUnpaidSchedules(b context.Backdrop, customerDisplayIds []string, date string) ([]models.BillingScheduleCombinedInfo, error) {
	db := b.GetDatabaseInstance()
	var unpaidSchedules []models.BillingScheduleCombinedInfo

	var query = db.
		Table("billing.billing_schedules as bs").
		Joins("join billing.loan_accounts as la on la.id = bs.loan_account_id").
		Joins("join billing.loans as l on l.id = la.loan_id").
		Joins("join billing.customers as cs on cs.id = l.customer_id").
		Select(" bs.id, bs.created_at, bs.updated_at, bs.loan_account_id, bs.start_date, bs.end_date, bs.week_count, bs.installment_amount, cs.display_id as customer_display_id,l.display_id as loan_display_id, cs.is_active as is_customer_active, l.loan_state as loan_state, l.id as loan_id, cs.id as customer_id").
		Where("bs.is_payment_done", false).
		Where("l.loan_state", "ACTIVE").
		Where("cs.is_active", true)

	if len(customerDisplayIds) > 0 {
		query = query.Where("cs.display_id IN (?)", customerDisplayIds)
	}

	if len(date) > 0 {
		query = query.Where("bs.end_date < ?", date)
	}

	if err := query.Find(&unpaidSchedules).Error; err != nil {
	}

	return unpaidSchedules, nil
}

func FindAllDefaultedSchedules(b context.Backdrop, customerDisplayIds []string, date string) ([]models.BillingScheduleCombinedInfo, error) {
	db := b.GetDatabaseInstance()
	var defaultedSchedules []models.BillingScheduleCombinedInfo

	var query = db.
		Table("billing.billing_schedules as bs").
		Joins("join billing.loan_accounts as la on la.id = bs.loan_account_id").
		Joins("join billing.loans as l on l.id = la.loan_id").
		Joins("join billing.customers as cs on cs.id = l.customer_id").
		Select(" bs.id, bs.created_at, bs.updated_at, bs.loan_account_id, bs.start_date, bs.end_date, bs.week_count, bs.installment_amount, cs.display_id as customer_display_id,l.display_id as loan_display_id, cs.is_active as is_customer_active, l.loan_state as loan_state, l.id as loan_id, cs.id as customer_id").
		Where("bs.is_default", true).Or("bs.end_date < ?", date).
		Where("l.loan_state", "ACTIVE").
		Where("cs.is_active", true)

	if len(customerDisplayIds) > 0 {
		query = query.Where("cs.display_id IN (?)", customerDisplayIds)
	}

	if len(date) > 0 {
		query = query.Where("bs.end_date < ?", date)
	}

	if err := query.Find(&defaultedSchedules).Error; err != nil {
	}

	return defaultedSchedules, nil
}

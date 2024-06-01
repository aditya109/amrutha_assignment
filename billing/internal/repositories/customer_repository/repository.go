package customer_repository

import (
	"errors"
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/internal/models"
	"github.com/aditya109/amrutha_assignment/pkg/connectors/database/postgres"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"gorm.io/gorm"
	"net/http"
)

func UniqueSave(b context.Backdrop, customer *models.Customer) (models.Customer, error) {
	db := b.GetDatabaseInstance()
	var existingCustomer models.Customer

	if result := db.Where(models.Customer{DisplayId: customer.DisplayId}).First(&existingCustomer); result.Error == nil {
		b.SetCustomErrorMessage("customer already exists")
		b.SetStatusCodeForResponse(http.StatusBadRequest)
		return existingCustomer, fmt.Errorf("customer with displayId: %v already exists, cannot proceed with this data", customer.DisplayId)
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return existingCustomer, fmt.Errorf("error occured while looking for customer with id %v", customer.DisplayId)
	}

	if err := postgres.ExecuteTransaction(b, db.FirstOrCreate(&existingCustomer, customer)); err != nil {
		return models.Customer{}, fmt.Errorf("failed to save customer: %w", err)
	}
	return existingCustomer, nil
}

func FindOne(b context.Backdrop, customer *models.Customer) error {
	db := b.GetDatabaseInstance()
	if result := db.Where(models.Customer{DisplayId: customer.DisplayId}).First(&customer); result.Error == nil {
		return nil
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error occured while looking for customer with id %v", customer.DisplayId)
	} else {
		return fmt.Errorf("customer with id %v not found", customer.DisplayId)
	}
}

func FindAllActiveCustomers(b context.Backdrop, ids []string) ([]models.Customer, error) {
	db := b.GetDatabaseInstance()
	var customers []models.Customer
	db = db.Where("is_active = ?", true)
	if len(ids) > 0 {
		db = db.Where("id IN ?", ids)
	}

	if result := db.Find(&customers); result.Error == nil {
		return customers, nil
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error occured while looking for customer with active loans")
	} else {
		return nil, fmt.Errorf("customer with active loans not found")
	}
}

func Update(b context.Backdrop, customer *models.Customer) error {
	db := b.GetDatabaseInstance()
	if err := postgres.ExecuteTransaction(b, db.Save(customer)); err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}
	return nil
}

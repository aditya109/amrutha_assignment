package customer

import (
	"github.com/aditya109/amrutha_assignment/billing/internal/models"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/customer_repository"
	"github.com/aditya109/amrutha_assignment/billing/pkg/constants"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"github.com/aditya109/amrutha_assignment/pkg/helpers"
	"time"
)

type MakeNewCustomerInputConstruct struct {
	Name    string
	Address string
}

func (c MakeNewCustomerInputConstruct) MakeNewCustomer(b context.Backdrop) (*models.Customer, error) {
	var customer = &models.Customer{
		Name:      c.Name,
		Address:   c.Address,
		IsActive:  false,
		Type:      models.RegularCustomerState,
		CreatedAt: helpers.CreatePointerForValue(time.Now()),
		UpdatedAt: helpers.CreatePointerForValue(time.Now()),
	}
	var err error
	if customer.DisplayId, err = helpers.CreateUniqueDisplayId(models.Customer{
		Name:    customer.Name,
		Address: customer.Address,
	}, constants.CUSTOMER_PREFIX); err != nil {
		return nil, err
	}
	if *customer, err = customer_repository.UniqueSave(b, customer); err != nil {
		return nil, err
	} else {
		return &models.Customer{
			Name:      customer.Name,
			Address:   customer.Address,
			DisplayId: customer.DisplayId,
		}, nil
	}
}

package customer

import (
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/internal/models"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/billing_schedule_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/customer_repository"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"github.com/aditya109/amrutha_assignment/pkg/helpers"
)

type InputForCheckCustomerState struct {
	CustomerId string
}

type Output struct {
	Message *string `json:"message"`
}

func (i InputForCheckCustomerState) IsDelinquent(b context.Backdrop) (*Output, error) {
	var existingCustomer = &models.Customer{DisplayId: i.CustomerId}
	var output = Output{}
	var err error

	if err = customer_repository.FindOne(b, existingCustomer); err != nil {
		return nil, err
	}

	switch {
	case !existingCustomer.IsActive:
		switch existingCustomer.Type {
		case models.RegularCustomerState:
			output.Message = helpers.CreatePointerForValue(fmt.Sprintf("the customer with id %s is not active, current state = REGULAR", existingCustomer.DisplayId))
		case models.DelinquentCustomerState:
			output.Message = helpers.CreatePointerForValue(fmt.Sprintf("the customer with id %s is currently not in active state, recorded state = DELINQUENT", existingCustomer.DisplayId))
		default:
			return nil, fmt.Errorf("unknown customer state")
		}
	case existingCustomer.IsActive:
		switch existingCustomer.Type {
		case models.RegularCustomerState:
			defaultedSchedules, err := billing_schedule_repository.FindAllDefaultedSchedules(b, []string{i.CustomerId}, "")
			if err != nil {
				return nil, err
			}
			if len(defaultedSchedules) >= 2 {
				var customer = &models.Customer{
					ID:   existingCustomer.ID,
					Type: models.DelinquentCustomerState,
				}
				if err := customer_repository.Update(b, customer); err != nil {
					return nil, err
				}
				output.Message = helpers.CreatePointerForValue(fmt.Sprintf("the customer with id %s is currently in active state, recorded state = DELINQUENT", existingCustomer.DisplayId))
			}
			output.Message = helpers.CreatePointerForValue(fmt.Sprintf("the customer with id %s is active, current state = REGULAR", existingCustomer.DisplayId))
		case models.DelinquentCustomerState:
			output.Message = helpers.CreatePointerForValue(fmt.Sprintf("the customer with id %s is currently in active state, recorded state = DELINQUENT", existingCustomer.DisplayId))
		default:
			return nil, fmt.Errorf("unknown customer state")
		}
	default:
		return nil, fmt.Errorf("reached unreachable state")
	}
	return &output, nil
}

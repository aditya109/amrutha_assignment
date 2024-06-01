package loan_account

import (
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/loan_account_repository"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"github.com/aditya109/amrutha_assignment/pkg/helpers"
	"github.com/shopspring/decimal"
)

type InputForGetOutstandingAmount struct {
	CustomerId string
}

type Output struct {
	OutstandingAmount string `json:"outstandingAmount"`
	Message           string `json:"message"`
}

func (i InputForGetOutstandingAmount) GetOutstanding(b context.Backdrop) (*Output, error) {
	if loanAccount, err := loan_account_repository.FindOneWithCustomerId(b, i.CustomerId); err != nil {
		return nil, fmt.Errorf("error finding outstanding amount for customerId %d: %w", i.CustomerId, err)
	} else {
		outstandingAmount, err := decimal.NewFromString(loanAccount.OutstandingAmount)
		if err != nil {
			return nil, fmt.Errorf("cannot convert outstanding amount to decimal, err: %v", err)
		}
		return &Output{
			OutstandingAmount: helpers.FormatCurrency(outstandingAmount),
			Message:           fmt.Sprintf("loan account id: %s has outstanding of %s", loanAccount.DisplayId, helpers.FormatCurrency(outstandingAmount)),
		}, nil
	}
}

package interest_calculation_rules

import (
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/internal/models"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"github.com/shopspring/decimal"
)

type ResultantConstructForLoanAccount struct {
	PayablePrincipalAmount  decimal.Decimal
	RateOfInterest          decimal.Decimal
	TotalTimeInWeeks        decimal.Decimal
	AccruedInterest         decimal.Decimal
	TotalPayableAmount      decimal.Decimal
	WeeklyInstallmentAmount decimal.Decimal
	TotalPaidAmount         decimal.Decimal
	OutstandingAmount       decimal.Decimal
}

func CalculateInitialConstructForLoanAccount(b context.Backdrop, loanConfig *models.LoanConfig) (*ResultantConstructForLoanAccount, error) {
	var construct = &ResultantConstructForLoanAccount{}
	var err error
	switch loanConfig.TypeOfLoan {
	case models.FlatInterestType:

		construct.PayablePrincipalAmount, err = decimal.NewFromString(loanConfig.PrincipalAmount)
		if err != nil {
			return nil, fmt.Errorf("error while parsing principal amount: %v", err)
		}
		construct.RateOfInterest, err = decimal.NewFromString(loanConfig.RateOfInterest)
		if err != nil {
			return nil, fmt.Errorf("error while parsing principal amount: %v", err)
		}
		rateOfInterestWeekly := construct.RateOfInterest.Div(decimal.NewFromInt(52))
		construct.TotalTimeInWeeks = decimal.NewFromInt(int64(loanConfig.MaxSpan))
		construct.AccruedInterest = rateOfInterestWeekly.Mul(construct.PayablePrincipalAmount).Mul(construct.TotalTimeInWeeks).Round(2)
		construct.TotalPayableAmount = construct.AccruedInterest.Add(construct.PayablePrincipalAmount).Round(2)
		construct.WeeklyInstallmentAmount = construct.TotalPayableAmount.Div(decimal.NewFromInt(int64(loanConfig.MaxSpan))).Round(2)
		construct.TotalPaidAmount = decimal.Zero
		return construct, nil
	default:
		return nil, fmt.Errorf("loan type not supported")
	}
}

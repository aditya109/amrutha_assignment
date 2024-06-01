package loan

import (
	"fmt"
	interestcalculationrules "github.com/aditya109/amrutha_assignment/billing/internal/domain/rules/interest_calculation"
	"github.com/aditya109/amrutha_assignment/billing/internal/models"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/billing_schedule_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/customer_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/loan_account_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/loan_config_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/loan_repository"
	"github.com/aditya109/amrutha_assignment/billing/pkg/constants"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"github.com/aditya109/amrutha_assignment/pkg/helpers"
	"log"
	"net/http"
	"time"
)

type TransitionLoanConstruct struct {
	CustomerId      string
	ConfigurationId *int
}

type Output struct {
	Loan            *models.Loan            `json:"loan,omitempty"`
	LoanAccount     *models.LoanAccount     `json:"loanAccount,omitempty"`
	NearestSchedule *models.BillingSchedule `json:"nearestSchedule,omitempty"`
	Customer        *models.Customer        `json:"customer,omitempty"`
}

func (c TransitionLoanConstruct) TransitionLoan(b context.Backdrop) (*Output, error) {
	var loan = &models.Loan{
		Customer: &models.Customer{DisplayId: c.CustomerId},
	}
	var existingCustomer = &models.Customer{DisplayId: c.CustomerId}
	var doesLoanExists bool
	var err error

	if err = customer_repository.FindOne(b, existingCustomer); err != nil {
		return nil, err
	}
	loan.Customer = existingCustomer
	if doesLoanExists, err = loan_repository.IfExists(b, loan); err != nil {
		return nil, err
	}

	switch {
	case !existingCustomer.IsActive:
		switch existingCustomer.Type {
		case models.RegularCustomerState:
			// customer has no active, inactive loan, hence a new inactive loan is opened
			if doesLoanExists {
				switch loan.LoanState {
				case models.ActiveLoanType:
					b.SetStatusCodeForResponse(http.StatusConflict)
					return nil, fmt.Errorf("there is an active loan on the user, loan id: %v", loan.DisplayId)
				case models.InactiveLoanType:
					loan.LoanState = models.ActiveLoanType
					loan.Customer.IsActive = true
					if err := loan_repository.Update(b, loan); err != nil {
						return nil, fmt.Errorf("error while updating loan state: %v from %s to %s", err, loan.LoanState, models.ActiveLoanType)
					}
					if err := customer_repository.Update(b, loan.Customer); err != nil {
						return nil, fmt.Errorf("error while changing activity state of customer, err: %v", err)
					}
					if result, err := attachNewLoanAccount(b, loan); err != nil {
						loan.LoanState = models.InactiveLoanType
						if err := loan_repository.Update(b, loan); err != nil {
							return nil, fmt.Errorf("error while updating loan state: %v from %s to %s", err, loan.LoanState, models.InactiveLoanType)
						}
						return nil, fmt.Errorf("error while updating loan schedule: %v", err)
					} else {
						return result, nil
					}
				case models.PaidLoanType:
					if result, err := createNewLoan(b, loan, c.ConfigurationId); err != nil {
						return nil, fmt.Errorf("error while creating new loan: %v", err)
					} else {
						return result, nil
					}
				default:
					return nil, fmt.Errorf("there is an unknown loan state: %v", loan.LoanState)
				}
			} else {
				if result, err := createNewLoan(b, loan, c.ConfigurationId); err != nil {
					return nil, fmt.Errorf("error while creating new loan: %v", err)
				} else {
					return result, nil
				}
			}
		case models.DelinquentCustomerState:
			// this is no yet covered in the requirements, hence ignoring this
		default:
		}
	case existingCustomer.IsActive:
		switch existingCustomer.Type {
		case models.RegularCustomerState:
			if doesLoanExists {
				if loan.PaymentCompletionCount == loan.LoanConfig.MaxSpan {
					loan.LoanState = models.PaidLoanType
					loan.Customer.Type = models.RegularCustomerState
				} else if loan.PaymentCompletionCount < loan.LoanConfig.MaxSpan && loan.MissedPaymentCount > 0 {
					loan.Customer.Type = models.DelinquentCustomerState
				}

				if err := loan_repository.Update(b, loan); err != nil {
					return nil, fmt.Errorf("error while updating loan state: %v from %s to %s", err, loan.LoanState, models.PaidLoanType)
				}
			} else {
				log.Println("loan for this user does not already exists, not possible")
			}
			// can only be called if customer has paid off his debt in full
		case models.DelinquentCustomerState:
			// this is no yet covered in the requirements, hence ignoring this
		default:
		}

	default:
	}

	return nil, nil
}

func createNewLoan(b context.Backdrop, loan *models.Loan, configurationId *int) (*Output, error) {
	var loanConfig = &models.LoanConfig{}
	var err error
	if configurationId == nil {
		b.SetStatusCodeForResponse(http.StatusBadRequest)
		return nil, fmt.Errorf("configuration id is required")
	}
	loan.LoanState = models.InactiveLoanType
	if err = loan_config_repository.FindOne(b, loanConfig); err != nil {
		return nil, err
	}
	loan.LoanConfig = loanConfig

	if loan.DisplayId, err = helpers.CreateUniqueDisplayId(models.Loan{
		Customer:   &models.Customer{DisplayId: loan.DisplayId, ID: loan.Customer.ID},
		LoanConfig: &models.LoanConfig{Id: loan.LoanConfig.Id},
	}, constants.LOAN_PREFIX); err != nil {
		return nil, err
	}
	err = loan_repository.Update(b, loan)
	if err != nil {
		return nil, err
	}
	return &Output{
		Loan: &models.Loan{
			DisplayId:              loan.DisplayId,
			LoanState:              loan.LoanState,
			MissedPaymentCount:     loan.MissedPaymentCount,
			PaymentCompletionCount: loan.PaymentCompletionCount,
		},
		Customer: &models.Customer{
			Name:      loan.Customer.Name,
			Address:   loan.Customer.Address,
			DisplayId: loan.Customer.DisplayId,
			Type:      loan.Customer.Type,
			IsActive:  loan.Customer.IsActive,
		},
	}, nil
}

func attachNewLoanAccount(b context.Backdrop, loan *models.Loan) (*Output, error) {
	var loanConfig = loan.LoanConfig
	var construct *interestcalculationrules.ResultantConstructForLoanAccount
	var err error
	if construct, err = interestcalculationrules.CalculateInitialConstructForLoanAccount(b, loanConfig); err != nil {
		return nil, fmt.Errorf("error while calculating construct for loan account: %v", err)
	}
	var loanAccount = &models.LoanAccount{
		CreatedAt:              helpers.CreatePointerForValue(time.Now()),
		UpdatedAt:              helpers.CreatePointerForValue(time.Now()),
		LoanId:                 int(loan.Id),
		Loan:                   loan,
		PayablePrincipalAmount: construct.PayablePrincipalAmount.StringFixed(2),
		AccruedInterest:        construct.AccruedInterest.StringFixed(2),
		TotalPayableAmount:     construct.TotalPayableAmount.StringFixed(2),
		TotalPaidAmount:        construct.TotalPaidAmount.StringFixed(2),
		OutstandingAmount:      construct.TotalPayableAmount.Sub(construct.TotalPaidAmount).StringFixed(2),
		InstallmentAmount:      construct.WeeklyInstallmentAmount.StringFixed(2),
	}
	if loanAccount.DisplayId, err = helpers.CreateUniqueDisplayId(models.LoanAccount{
		LoanId: int(loan.Id),
		Loan: &models.Loan{
			Id:         loan.Id,
			DisplayId:  loan.DisplayId,
			Customer:   &models.Customer{DisplayId: loan.DisplayId, ID: loan.Customer.ID},
			LoanConfig: &models.LoanConfig{Id: loan.LoanConfig.Id},
		},
	}, constants.LOAN_ACCOUNT_PREFIX); err != nil {
		return nil, err
	}
	if err = loan_account_repository.Update(b, loanAccount); err != nil {
		return nil, fmt.Errorf("error while updating loan account: %v", err)
	}

	// create first schedule
	var nearestSchedule = &models.BillingSchedule{
		CreatedAt:         helpers.CreatePointerForValue(time.Now()),
		UpdatedAt:         helpers.CreatePointerForValue(time.Now()),
		LoanAccountId:     loanAccount.Id,
		StartDate:         helpers.CreatePointerForValue(time.Now()),
		EndDate:           helpers.CreatePointerForValue(time.Now().Add(7 * 24 * time.Hour)),
		WeekCount:         1,
		InstallmentAmount: construct.WeeklyInstallmentAmount.String(),
	}
	if err = billing_schedule_repository.Save(b, nearestSchedule); err != nil {
		return nil, fmt.Errorf("error while updating loan schedule: %v", err)
	}
	return &Output{

		LoanAccount: &models.LoanAccount{
			PayablePrincipalAmount: helpers.FormatCurrency(construct.PayablePrincipalAmount),
			AccruedInterest:        helpers.FormatCurrency(construct.AccruedInterest),
			TotalPayableAmount:     helpers.FormatCurrency(construct.TotalPayableAmount),
			TotalPaidAmount:        helpers.FormatCurrency(construct.TotalPaidAmount),
			OutstandingAmount:      helpers.FormatCurrency(construct.OutstandingAmount),
			InstallmentAmount:      helpers.FormatCurrency(construct.WeeklyInstallmentAmount),
			DisplayId:              loanAccount.DisplayId,
		},
		NearestSchedule: &models.BillingSchedule{
			StartDate:         nearestSchedule.StartDate,
			EndDate:           nearestSchedule.EndDate,
			WeekCount:         nearestSchedule.WeekCount,
			InstallmentAmount: nearestSchedule.InstallmentAmount,
		},
		Loan: &models.Loan{
			DisplayId:              loan.DisplayId,
			LoanState:              loan.LoanState,
			MissedPaymentCount:     loan.MissedPaymentCount,
			PaymentCompletionCount: loan.PaymentCompletionCount,
		},
		Customer: &models.Customer{
			Name:      loan.Customer.Name,
			Address:   loan.Customer.Address,
			DisplayId: loan.Customer.DisplayId,
			Type:      loan.Customer.Type,
			IsActive:  loan.Customer.IsActive,
		},
	}, nil
}

package loan

import (
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/internal/models"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/billing_schedule_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/customer_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/loan_account_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/loan_config_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/loan_repository"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"github.com/aditya109/amrutha_assignment/pkg/helpers"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

type TransistionLoanConstruct struct {
	CustomerId      string
	ConfigurationId *int
}

func (c TransistionLoanConstruct) TransistionLoan(b context.Backdrop) (*models.Loan, error) {
	var loan = &models.Loan{
		Customer: models.Customer{DisplayId: c.CustomerId},
	}
	var existingCustomer = &models.Customer{DisplayId: c.CustomerId}
	var doesLoanExists bool
	var err error
	g := new(errgroup.Group)
	g.Go(func() error {
		var e error
		defer func() {
			if doesLoanExists, err = loan_repository.IfExists(b, loan); err != nil {
				e = err
			}
		}()
		return e
	})
	g.Go(func() error {
		var e error
		defer func() {
			if err = customer_repository.FindOne(b, existingCustomer); err != nil {
				e = err
			}

		}()
		return e
	})

	if syncErr := g.Wait(); syncErr != nil {
		return nil, fmt.Errorf("error while making concurrent database calls: err : %v", syncErr)
	}
	loan.Customer = *existingCustomer
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
					if err := loan_repository.Update(b, loan); err != nil {
						return nil, fmt.Errorf("error while updating loan state: %v from %s to %s", err, loan.LoanState, models.ActiveLoanType)
					}
					if err := attachNewLoanAccount(b, loan); err != nil {
						return nil, fmt.Errorf("error while updating loan schedule: %v", err)
					}
				case models.PaidLoanType:
					if err := createNewLoan(b, loan, c.ConfigurationId); err != nil {
						return nil, fmt.Errorf("error while creating new loan: %v", err)
					}
				default:
					return nil, fmt.Errorf("there is an unknown loan state: %v", loan.LoanState)
				}
			} else {
				if err := createNewLoan(b, loan, c.ConfigurationId); err != nil {
					return nil, fmt.Errorf("error while creating new loan: %v", err)
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
					return nil, fmt.Errorf("error while updating loan state: %v from %s to %s", err, loan.LoanState, models.ActiveLoanType)
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

func createNewLoan(b context.Backdrop, loan *models.Loan, configurationId *int) error {
	var loanConfig = &models.LoanConfig{}
	var err error
	if configurationId == nil {
		b.SetStatusCodeForResponse(http.StatusBadRequest)
		return fmt.Errorf("configuration id is required")
	}
	loan.LoanState = models.InactiveLoanType
	if err = loan_config_repository.FindOne(b, loanConfig); err != nil {
		return err
	}
	loan.LoanConfig = *loanConfig

	if loan.DisplayId, err = helpers.CreateUniqueDisplayId(models.Loan{
		Customer:   models.Customer{DisplayId: loan.DisplayId, ID: loan.Customer.ID},
		LoanConfig: models.LoanConfig{Id: loan.LoanConfig.Id},
	}); err != nil {
		return err
	}

	return loan_repository.Update(b, loan)
}

func attachNewLoanAccount(b context.Backdrop, loan *models.Loan) error {
	var loanConfig = loan.LoanConfig
	payablePrincipalAmount, err := decimal.NewFromString(loanConfig.PrincipalAmount)
	if err != nil {
		return fmt.Errorf("error while parsing principal amount: %v", err)
	}
	rateOfInterest, err := decimal.NewFromString(loanConfig.RateOfInterest)
	if err != nil {
		return fmt.Errorf("error while parsing principal amount: %v", err)
	}
	var accruedInterest = rateOfInterest.Mul(payablePrincipalAmount).Mul(decimal.NewFromInt(int64(loanConfig.MaxSpan)).Div(decimal.NewFromInt(52))).Round(2)
	var totalPayableAmount = accruedInterest.Add(payablePrincipalAmount).Round(2)
	var weeklyInstallmentAmount = totalPayableAmount.Div(decimal.NewFromInt(int64(loanConfig.MaxSpan))).Round(2)
	var totalPaidAmount = decimal.Zero
	var loanAccount = &models.LoanAccount{
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
		LoanId:                 0,
		Loan:                   *loan,
		PayablePrincipalAmount: payablePrincipalAmount.String(),
		AccruedInterest:        accruedInterest.String(),
		TotalPayableAmount:     totalPayableAmount.String(),
		TotalPaidAmount:        totalPaidAmount.String(),
		OutstandingAmount:      totalPayableAmount.Sub(totalPaidAmount).String(),
		InstallmentAmount:      weeklyInstallmentAmount.String(),
	}
	if err = loan_account_repository.Update(b, loanAccount); err != nil {
		return fmt.Errorf("error while updating loan account: %v", err)
	}

	// create first schedule
	var schedule = &models.BillingSchedule{
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		LoanId:            0,
		Loan:              *loan,
		StartDate:         time.Now(),
		EndDate:           time.Now().Add(time.Duration(7 * 24 * time.Hour)),
		WeekCount:         1,
		InstallmentAmount: weeklyInstallmentAmount.String(),
	}
	if err = billing_schedule_repository.Update(b, schedule); err != nil {
		return fmt.Errorf("error while updating loan schedule: %v", err)
	}
	return nil
}

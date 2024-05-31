package loan

import (
	"fmt"
	interest_calculation_rules "github.com/aditya109/amrutha_assignment/billing/internal/domain/rules/interest_calculation"
	"github.com/aditya109/amrutha_assignment/billing/internal/models"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/billing_schedule_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/customer_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/loan_account_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/loan_config_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/loan_repository"
	"github.com/aditya109/amrutha_assignment/billing/pkg/constants"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"github.com/aditya109/amrutha_assignment/pkg/helpers"
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
		Customer: &models.Customer{DisplayId: c.CustomerId},
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
	loan.Customer = existingCustomer
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
						loan.LoanState = models.InactiveLoanType
						if err := loan_repository.Update(b, loan); err != nil {
							return nil, fmt.Errorf("error while updating loan state: %v from %s to %s", err, loan.LoanState, models.ActiveLoanType)
						}
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
	loan.LoanConfig = loanConfig

	if loan.DisplayId, err = helpers.CreateUniqueDisplayId(models.Loan{
		Customer:   &models.Customer{DisplayId: loan.DisplayId, ID: loan.Customer.ID},
		LoanConfig: &models.LoanConfig{Id: loan.LoanConfig.Id},
	}, constants.LOAN_PREFIX); err != nil {
		return err
	}
	return loan_repository.Update(b, loan)
}

func attachNewLoanAccount(b context.Backdrop, loan *models.Loan) error {
	var loanConfig = loan.LoanConfig
	var construct *interest_calculation_rules.ResultantConstructForLoanAccount
	var err error
	if construct, err = interest_calculation_rules.CalculateInitialConstructForLoanAccount(b, loanConfig); err != nil {
		return fmt.Errorf("error while calculating construct for loan account: %v", err)
	}
	var loanAccount = &models.LoanAccount{
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
		LoanId:                 0,
		Loan:                   *loan,
		PayablePrincipalAmount: construct.PayablePrincipalAmount.String(),
		AccruedInterest:        construct.AccruedInterest.String(),
		TotalPayableAmount:     construct.TotalPayableAmount.String(),
		TotalPaidAmount:        construct.TotalPaidAmount.String(),
		OutstandingAmount:      construct.TotalPayableAmount.Sub(construct.TotalPaidAmount).String(),
		InstallmentAmount:      construct.WeeklyInstallmentAmount.String(),
	}
	if loanAccount.DisplayId, err = helpers.CreateUniqueDisplayId(models.LoanAccount{
		LoanId: int(loan.Id),
		Loan: models.Loan{
			Id:         loan.Id,
			DisplayId:  loan.DisplayId,
			Customer:   &models.Customer{DisplayId: loan.DisplayId, ID: loan.Customer.ID},
			LoanConfig: &models.LoanConfig{Id: loan.LoanConfig.Id},
		},
	}, constants.LOAN_ACCOUNT_PREFIX); err != nil {
		return err
	}
	if err = loan_account_repository.Update(b, loanAccount); err != nil {
		return fmt.Errorf("error while updating loan account: %v", err)
	}

	// create first schedule
	var schedule = &models.BillingSchedule{
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		LoanAccountId:     loanAccount.Id,
		StartDate:         time.Now(),
		EndDate:           time.Now().Add(7 * 24 * time.Hour),
		WeekCount:         1,
		InstallmentAmount: construct.WeeklyInstallmentAmount.String(),
	}
	if err = billing_schedule_repository.Update(b, schedule); err != nil {
		return fmt.Errorf("error while updating loan schedule: %v", err)
	}
	return nil
}

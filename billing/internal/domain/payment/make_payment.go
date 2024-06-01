package payment

import (
	"fmt"
	domainmodels "github.com/aditya109/amrutha_assignment/billing/internal/domain/models"
	"github.com/aditya109/amrutha_assignment/billing/internal/models"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/billing_schedule_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/customer_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/loan_account_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/loan_repository"
	"github.com/aditya109/amrutha_assignment/billing/internal/repositories/payment_repository"
	"github.com/aditya109/amrutha_assignment/billing/pkg/constants"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"github.com/aditya109/amrutha_assignment/pkg/helpers"
	"github.com/shopspring/decimal"
	"time"
)

type MakePaymentInputConstruct struct {
	CustomerId             string
	Amount                 string
	DateOfTransaction      string
	TransactionReferenceId string
}

// MakePayment Due to lack of time, I wrote the entire logic as a part of API core, ideally it ought to be run with an asynchronous function commanded by Kafka messages published within this API.
func (c MakePaymentInputConstruct) MakePayment(b context.Backdrop) (*domainmodels.Output, error) {
	var payment *models.Payment
	var existingCustomer = &models.Customer{DisplayId: c.CustomerId}
	var err error

	if err = customer_repository.FindOne(b, existingCustomer); err != nil {
		return nil, err
	}
	unpaidSchedules, err := billing_schedule_repository.FindAllUnpaidSchedules(b, []string{c.CustomerId}, "")
	if err != nil {
		return nil, err
	}
	loanAccount := &models.LoanAccount{
		Id: unpaidSchedules[0].LoanAccountId,
	}

	if err := loan_account_repository.FindOne(b, loanAccount); err != nil {
		return nil, fmt.Errorf("loan_account_repository.FindOne: %v", err)
	}
	switch {
	case len(unpaidSchedules) == 0:
		return nil, fmt.Errorf("no schedule for upcoming payment exists, payment can not be accepted, please try again in some time")
	case len(unpaidSchedules) == 1:
		paidAmount, err := decimal.NewFromString(c.Amount)
		if err != nil {
			return nil, fmt.Errorf("cannot convert amount to decimal, err: %v", err)
		}
		installmentAmount, _ := decimal.NewFromString(unpaidSchedules[0].InstallmentAmount)

		if !paidAmount.Equal(installmentAmount) {
			b.SetCustomErrorMessage("paid amount and installment amount do not match, kindly pay the exact amount to proceed")
			return nil, fmt.Errorf("paid amount and installment amount do not match, unable to register payment")
		} else {
			payment = &models.Payment{
				CreatedAt:                    helpers.CreatePointerForValue(time.Now()),
				UpdatedAt:                    helpers.CreatePointerForValue(time.Now()),
				CustomerId:                   unpaidSchedules[0].CustomerId,
				LoanAccountId:                unpaidSchedules[0].LoanAccountId,
				PaidAmount:                   c.Amount,
				ClientTransactionReferenceId: c.TransactionReferenceId,
				DateOfTransaction:            helpers.GetDateAsTimeFromString(c.DateOfTransaction),
				IsAccepted:                   true,
				ScheduleId:                   unpaidSchedules[0].Id,
			}
			if payment.PaymentDisplayId, err = helpers.CreateUniqueDisplayId(models.Payment{
				CreatedAt:                    payment.CreatedAt,
				CustomerId:                   payment.CustomerId,
				LoanAccountId:                payment.LoanAccountId,
				ClientTransactionReferenceId: payment.ClientTransactionReferenceId,
				DateOfTransaction:            payment.DateOfTransaction,
			}, constants.PaymentPrefix); err != nil {
				return nil, err
			}

			if err := payment_repository.Update(b, payment); err != nil {
				return nil, err
			}
			payment.LoanAccount = loanAccount
			billingSchedule := &models.BillingSchedule{
				Id:            unpaidSchedules[0].Id,
				IsPaymentDone: true,
				LoanAccountId: payment.LoanAccountId,
			}
			totalPaidAmount, err := decimal.NewFromString(payment.LoanAccount.TotalPaidAmount)
			if err != nil {
				return nil, fmt.Errorf("cannot convert total paid amount to decimal, err: %v", err)
			}
			outstandingAmount, err := decimal.NewFromString(payment.LoanAccount.OutstandingAmount)
			if err != nil {
				return nil, fmt.Errorf("cannot convert outstanding amount to decimal, err: %v", err)
			}

			installmentAmount, err := decimal.NewFromString(payment.LoanAccount.InstallmentAmount)
			if err != nil {
				return nil, fmt.Errorf("cannot convert installment amount to decimal, err: %v", err)
			}

			loanAccount.TotalPaidAmount = totalPaidAmount.Add(paidAmount).Round(2).StringFixed(2)
			loanAccount.OutstandingAmount = outstandingAmount.Sub(paidAmount).Round(2).StringFixed(2)
			var shouldIncrementMissedPaymentCount bool
			if helpers.GetDateAsTimeFromString(c.DateOfTransaction).After(unpaidSchedules[0].EndDate) {
				billingSchedule.IsDefault = true
				shouldIncrementMissedPaymentCount = true
			}
			if err := billing_schedule_repository.UpdateOnly(b, billingSchedule); err != nil {
				return nil, err
			}

			if err := loan_account_repository.Update(b, loanAccount); err != nil {
				return nil, err
			}
			if err := loan_account_repository.FindOne(b, loanAccount); err != nil {
				return nil, err
			}
			var loan = &models.Loan{Id: unpaidSchedules[0].LoanId}
			if err := loan_repository.UpdateAfterPayment(b, loan, shouldIncrementMissedPaymentCount); err != nil {
				return nil, err
			}
			if err := loan_repository.FindOne(b, loan); err != nil {
				return nil, err
			}

			// create next schedule
			var nextBillingScheduleStartDate = unpaidSchedules[0].EndDate.Add(1 * time.Second)
			var nextBillingScheduleEndDate = nextBillingScheduleStartDate.Add(7 * 24 * time.Hour)
			var nextSchedule = &models.BillingSchedule{
				CreatedAt:         helpers.CreatePointerForValue(time.Now()),
				UpdatedAt:         helpers.CreatePointerForValue(time.Now()),
				LoanAccountId:     loanAccount.Id,
				StartDate:         helpers.CreatePointerForValue(nextBillingScheduleStartDate),
				EndDate:           helpers.CreatePointerForValue(nextBillingScheduleEndDate),
				WeekCount:         1,
				InstallmentAmount: unpaidSchedules[0].InstallmentAmount,
			}
			if nextBillingScheduleEndDate.Before(time.Now()) {
				nextSchedule.IsDefault = true
			}
			if err := billing_schedule_repository.Save(b, nextSchedule); err != nil {
				return nil, err
			}

			if shouldIncrementMissedPaymentCount && nextSchedule.IsDefault {
				var customer = &models.Customer{
					ID:   uint(unpaidSchedules[0].CustomerId),
					Type: models.DelinquentCustomerState,
				}
				if err := customer_repository.Update(b, customer); err != nil {
					return nil, err
				}
			}
			return &domainmodels.Output{
				LoanAccount: &models.LoanAccount{
					TotalPaidAmount:   helpers.FormatCurrency(totalPaidAmount),
					OutstandingAmount: helpers.FormatCurrency(outstandingAmount),
					DisplayId:         loanAccount.DisplayId,
				},
				NextBillingSchedule: &models.BillingSchedule{
					StartDate:         nextSchedule.StartDate,
					EndDate:           nextSchedule.EndDate,
					WeekCount:         nextSchedule.WeekCount,
					InstallmentAmount: helpers.FormatCurrency(installmentAmount),
				},
				Loan: &models.Loan{
					DisplayId:              loan.DisplayId,
					LoanState:              loan.LoanState,
					MissedPaymentCount:     loan.MissedPaymentCount,
					PaymentCompletionCount: loan.PaymentCompletionCount,
				},
				Customer: &models.Customer{
					Name:      existingCustomer.Name,
					Address:   existingCustomer.Address,
					IsActive:  existingCustomer.IsActive,
					Type:      existingCustomer.Type,
					DisplayId: existingCustomer.DisplayId,
				},
				Payment: &models.Payment{
					PaidAmount:                   helpers.FormatCurrency(paidAmount),
					PaymentDisplayId:             payment.PaymentDisplayId,
					ClientTransactionReferenceId: payment.ClientTransactionReferenceId,
					DateOfTransaction:            payment.DateOfTransaction,
				},
			}, nil
		}
	default:
		return nil, fmt.Errorf("unhandled payment operation")
	}
}

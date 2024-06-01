package modules

import (
	"github.com/aditya109/amrutha_assignment/billing/internal/domain/customer"
	"github.com/aditya109/amrutha_assignment/billing/internal/domain/loan"
	"github.com/aditya109/amrutha_assignment/billing/internal/domain/loan_account"
	"github.com/aditya109/amrutha_assignment/billing/internal/domain/payment"
	"github.com/aditya109/amrutha_assignment/pkg/context"
)

func createNewCustomerService(b context.Backdrop, body CustomerDto) (interface{}, error) {
	return customer.MakeNewCustomerInputConstruct{
		Name:    body.Name,
		Address: body.Address,
	}.MakeNewCustomer(b)
}

func transitionLoanService(b context.Backdrop, body TransitionLoanRequestDto) (interface{}, error) {
	return loan.TransitionLoanConstruct{
		CustomerId:      body.CustomerId,
		ConfigurationId: body.ConfigurationId,
	}.TransitionLoan(b)
}

func makePaymentService(b context.Backdrop, body MakePaymentRequestDto) (interface{}, error) {
	return payment.MakePaymentInputConstruct{
		CustomerId:             body.CustomerId,
		Amount:                 body.Amount,
		DateOfTransaction:      body.DateOfTransaction,
		TransactionReferenceId: body.TransactionReferenceId,
	}.MakePayment(b)
}

func getCustomerStateForDelinquencyService(b context.Backdrop, customerId string) (interface{}, error) {
	return customer.InputForCheckCustomerState{
		CustomerId: customerId,
	}.IsDelinquent(b)
}

func getOutstandingAmountForDelinquencyService(b context.Backdrop, customerId string) (interface{}, error) {
	return loan_account.InputForGetOutstandingAmount{
		CustomerId: customerId,
	}.GetOutstanding(b)
}

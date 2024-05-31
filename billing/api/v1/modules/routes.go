package modules

import (
	"github.com/aditya109/amrutha_assignment/billing/pkg/constants"
	"github.com/aditya109/amrutha_assignment/pkg/api"
	"net/http"
)

const (
	getOutstandingAmountForCustomerRoute = "/customers/:customerId/outstanding-amount"
	getCustomerStateForDeliquencyRoute   = "/customers/:customerId/upstate"
	postMakePaymentRoute                 = "/customers/:customerId/make-payment"
	postCreateNewCustomer                = "/customer/new"
	putActivateLoan                      = "/customers/:customerId/loan/transition-state"
)

func GetModule() api.ApiModule {
	return api.ApiModule{
		Module:     constants.SERVICE_IDENTIFIER,
		ApiVersion: "v1",
		Routes: []api.Route{
			{
				Path:       getOutstandingAmountForCustomerRoute,
				Method:     http.MethodGet,
				Controller: api.WrapHighOrderControl(GetOutstandingAmountForCustomerController),
			},
			{
				Path:       getCustomerStateForDeliquencyRoute,
				Method:     http.MethodGet,
				Controller: api.WrapHighOrderControl(GetCustomerStateForDeliquencyRoute),
			},
			{
				Path:       postMakePaymentRoute,
				Method:     http.MethodPost,
				Controller: api.WrapHighOrderControl(MakePaymentController),
			},
			{
				Path:       postCreateNewCustomer,
				Method:     http.MethodPost,
				Controller: api.WrapHighOrderControl(CreateNewCustomerController),
			},
			{
				Path:       putActivateLoan,
				Method:     http.MethodPut,
				Controller: api.WrapHighOrderControl(ActivateLoanController),
			},
		},
	}
}

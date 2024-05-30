package modules

import (
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/pkg/constants"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"github.com/aditya109/amrutha_assignment/pkg/models"
	"net/http"
)

func GetOutstandingAmountForCustomerController(b context.Backdrop) {
}

func GetCustomerStateForDeliquencyRoute(b context.Backdrop) {

}

func MakePaymentController(b context.Backdrop) {}

func CreateNewCustomerController(b context.Backdrop) {
	var body CustomerDto
	if err := b.ReadRequestPayload(&body); err != nil {
		b.Error(http.StatusBadRequest, models.Error{
			Code:              constants.VALIDATION_FAILED,
			Message:           constants.GENERIC_ERROR_MESSAGE,
			ResolutionMessage: fmt.Errorf("%s, err: %v", constants.REQUEST_BODY_VALIDATION_FAILED, err).Error(),
			Data:              nil,
		})
		return
	}
	if result, err := createNewCustomerService(b, body); err != nil {
		b.Error(http.StatusInternalServerError, models.Error{
			Code:              models.INTERNAL_SERVER_ERROR_MESSAGE,
			Message:           constants.GENERIC_ERROR_MESSAGE,
			ResolutionMessage: fmt.Errorf("err: %v", err).Error(),
			Data:              nil,
		})
		return
	} else {
		b.Response(http.StatusOK, result)
	}
}

func ActivateLoanController(b context.Backdrop) {
	var body ActivateLoanRequestDto
	body.CustomerId = b.GetContext().Param("customerId")

	if err := b.ReadRequestPayload(&body); err != nil {
		b.Error(http.StatusBadRequest, models.Error{
			Code:              constants.VALIDATION_FAILED,
			Message:           constants.GENERIC_ERROR_MESSAGE,
			ResolutionMessage: fmt.Errorf("%s, err: %v", constants.REQUEST_BODY_VALIDATION_FAILED, err).Error(),
			Data:              nil,
		})
		return
	}

	if result, err := activateLoanService(b, body); err != nil {
		b.Error(http.StatusInternalServerError, models.Error{
			Code:              models.INTERNAL_SERVER_ERROR_MESSAGE,
			Message:           constants.GENERIC_ERROR_MESSAGE,
			ResolutionMessage: fmt.Errorf("err: %v", err).Error(),
			Data:              nil,
		})
		return
	} else {
		b.Response(http.StatusOK, result)
	}
}

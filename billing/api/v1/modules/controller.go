package modules

import (
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/pkg/constants"
	"github.com/aditya109/amrutha_assignment/pkg/context"
	"github.com/aditya109/amrutha_assignment/pkg/models"
	"net/http"
)

func GetOutstandingAmountForCustomerController(b context.Backdrop) {
	customerId := b.GetContext().Param("customerId")
	if result, err := getOutstandingAmountForDelinquencyService(b, customerId); err != nil {
		b.Error(http.StatusInternalServerError, models.Error{
			Code:              models.InternalServerErrorMessage,
			Message:           constants.GenericErrorMessage,
			ResolutionMessage: fmt.Errorf("err: %v", err).Error(),
			Data:              nil,
		})
		return
	} else {
		b.Response(http.StatusOK, result)
	}
}

func GetCustomerStateForDelinquencyRoute(b context.Backdrop) {
	customerId := b.GetContext().Param("customerId")
	if result, err := getCustomerStateForDelinquencyService(b, customerId); err != nil {
		b.Error(http.StatusInternalServerError, models.Error{
			Code:              models.InternalServerErrorMessage,
			Message:           constants.GenericErrorMessage,
			ResolutionMessage: fmt.Errorf("err: %v", err).Error(),
			Data:              nil,
		})
		return
	} else {
		b.Response(http.StatusOK, result)
	}
}

func MakePaymentController(b context.Backdrop) {
	var body MakePaymentRequestDto
	body.CustomerId = b.GetContext().Param("customerId")
	if err := b.ReadRequestPayload(&body); err != nil {
		b.Error(http.StatusBadRequest, models.Error{
			Code:              constants.ValidationFailed,
			Message:           constants.GenericErrorMessage,
			ResolutionMessage: fmt.Errorf("%s, err: %v", constants.RequestBodyValidationFailed, err).Error(),
			Data:              nil,
		})
		return
	}
	if result, err := makePaymentService(b, body); err != nil {
		b.Error(http.StatusInternalServerError, models.Error{
			Code:              models.InternalServerErrorMessage,
			Message:           constants.GenericErrorMessage,
			ResolutionMessage: fmt.Errorf("err: %v", err).Error(),
			Data:              nil,
		})
		return
	} else {
		b.Response(http.StatusOK, result)
	}
}

func CreateNewCustomerController(b context.Backdrop) {
	var body CustomerDto
	if err := b.ReadRequestPayload(&body); err != nil {
		b.Error(http.StatusBadRequest, models.Error{
			Code:              constants.ValidationFailed,
			Message:           constants.GenericErrorMessage,
			ResolutionMessage: fmt.Errorf("%s, err: %v", constants.RequestBodyValidationFailed, err).Error(),
			Data:              nil,
		})
		return
	}
	if result, err := createNewCustomerService(b, body); err != nil {
		b.Error(http.StatusInternalServerError, models.Error{
			Code:              models.InternalServerErrorMessage,
			Message:           constants.GenericErrorMessage,
			ResolutionMessage: fmt.Errorf("err: %v", err).Error(),
			Data:              nil,
		})
		return
	} else {
		b.Response(http.StatusOK, result)
	}
}

func TransitionLoanController(b context.Backdrop) {
	var body TransitionLoanRequestDto
	body.CustomerId = b.GetContext().Param("customerId")

	if err := b.ReadRequestPayload(&body); err != nil {
		b.Error(http.StatusBadRequest, models.Error{
			Code:              constants.ValidationFailed,
			Message:           constants.GenericErrorMessage,
			ResolutionMessage: fmt.Errorf("%s, err: %v", constants.RequestBodyValidationFailed, err).Error(),
			Data:              nil,
		})
		return
	}

	if result, err := transitionLoanService(b, body); err != nil {
		b.Error(http.StatusInternalServerError, models.Error{
			Code:              models.InternalServerErrorMessage,
			Message:           constants.GenericErrorMessage,
			ResolutionMessage: fmt.Errorf("err: %v", err).Error(),
			Data:              nil,
		})
		return
	} else {
		b.Response(http.StatusOK, result)
	}
}

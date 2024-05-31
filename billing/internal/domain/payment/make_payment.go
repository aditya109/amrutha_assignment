package payment

import (
	"github.com/aditya109/amrutha_assignment/billing/internal/models"
	"github.com/aditya109/amrutha_assignment/pkg/context"
)

type MakePaymentInputConstruct struct {
}

func (c MakePaymentInputConstruct) MakePayment(b context.Backdrop) (*models.Payment, error) {
	var payment *models.Payment
	return payment, nil
}

package domainmodels

import "github.com/aditya109/amrutha_assignment/billing/internal/models"

type Output struct {
	Loan                *models.Loan            `json:"loan,omitempty"`
	LoanAccount         *models.LoanAccount     `json:"loanAccount,omitempty"`
	NearestSchedule     *models.BillingSchedule `json:"nearestSchedule,omitempty"`
	Customer            *models.Customer        `json:"customer,omitempty"`
	Payment             *models.Payment         `json:"payment,omitempty"`
	NextBillingSchedule *models.BillingSchedule `json:"nextBillingSchedule,omitempty"`
}

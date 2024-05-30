package models

import "time"

type Loan struct {
	Id                     uint64     `json:"id" gorm:"column:id;primaryKey"`
	CreatedAt              time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt              time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	CustomerId             uint64     `gorm:"column:customer_id" json:"customerId"`
	LoanConfigId           int        `gorm:"column:loan_config_id" json:"loanConfigId"`
	LoanConfig             LoanConfig `gorm:"foreignKey:LoanConfigId"`
	Customer               Customer   `gorm:"foreignKey:CustomerId"`
	PaymentCompletionCount int        `json:"payment_completion_count" gorm:"column:payment_completion_count"`
	MissedPaymentCount     int        `json:"missed_payment_count" gorm:"column:missed_payment_count"`
	LoanState              LoanState  `json:"loan_state" gorm:"column:loan_state"`
	DisplayId              string     `json:"displayId" gorm:"column:display_id"`
}

type LoanState string

const (
	ActiveLoanType   LoanState = "ACTIVE"
	PaidLoanType     LoanState = "PAID"
	InactiveLoanType LoanState = "INACTIVE"
)

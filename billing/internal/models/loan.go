package models

import "time"

type Loan struct {
	Id                     uint64      `json:"-" gorm:"column:id;primaryKey"`
	CreatedAt              *time.Time  `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt              *time.Time  `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	CustomerId             uint64      `gorm:"column:customer_id" json:"customerId,omitempty"`
	LoanConfigId           int         `gorm:"column:loan_config_id" json:"loanConfigId,omitempty"`
	LoanConfig             *LoanConfig `gorm:"foreignKey:LoanConfigId" json:"-"`
	Customer               *Customer   `gorm:"foreignKey:CustomerId" json:"-"`
	PaymentCompletionCount int         `json:"paymentCompletionCount" gorm:"column:payment_completion_count"`
	MissedPaymentCount     int         `json:"missedPaymentCount" gorm:"column:missed_payment_count"`
	LoanState              LoanState   `json:"loanState" gorm:"column:loan_state"`
	DisplayId              string      `json:"displayId" gorm:"column:display_id"`
}

type LoanState string

const (
	ActiveLoanType   LoanState = "ACTIVE"
	PaidLoanType     LoanState = "PAID"
	InactiveLoanType LoanState = "INACTIVE"
)

package models

import "time"

type Loan struct {
	Id                     uint64      `gorm:"column:id;primaryKey" json:"-"`
	CreatedAt              *time.Time  `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt              *time.Time  `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	CustomerId             uint64      `gorm:"column:customer_id" json:"customerId,omitempty"`
	LoanConfigId           int         `gorm:"column:loan_config_id" json:"loanConfigId,omitempty"`
	LoanConfig             *LoanConfig `gorm:"foreignKey:LoanConfigId" json:"-"`
	Customer               *Customer   `gorm:"foreignKey:CustomerId" json:"-"`
	PaymentCompletionCount int         `gorm:"column:payment_completion_count" json:"paymentCompletionCount"`
	MissedPaymentCount     int         `gorm:"column:missed_payment_count" json:"missedPaymentCount"`
	LoanState              LoanState   `gorm:"column:loan_state" json:"loanState"`
	DisplayId              string      `gorm:"column:display_id" json:"displayId"`
}

type LoanState string

const (
	ActiveLoanType   LoanState = "ACTIVE"
	PaidLoanType     LoanState = "PAID"
	InactiveLoanType LoanState = "INACTIVE"
)

package models

import (
	"time"
)

type BillingSchedule struct {
	Id                uint64       `json:"-" gorm:"column:id;primaryKey"`
	CreatedAt         *time.Time   `gorm:"column:created_at" json:"-"`
	UpdatedAt         *time.Time   `gorm:"column:updated_at" json:"-"`
	LoanAccountId     uint64       `gorm:"column:loan_account_id" json:"-"`
	LoanAccount       *LoanAccount `gorm:"foreignKey:LoanAccountId" json:"-"`
	StartDate         *time.Time   `gorm:"column:start_date" json:"startDate"`
	EndDate           *time.Time   `gorm:"column:end_date" json:"endDate"`
	WeekCount         int          `gorm:"column:week_count" json:"weekCount"`
	InstallmentAmount string       `gorm:"column:installment_amount" json:"installmentAmount"`
	IsDefault         bool         `gorm:"column:is_default" json:"isDefault"`
	IsPaymentDone     bool         `gorm:"column:is_payment_done" json:"isPaymentDone"`
}

type BillingScheduleCombinedInfo struct {
	Id                uint64    `json:"id" gorm:"column:id;primaryKey"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt         time.Time `gorm:"column:updated_at" json:"updatedAt"`
	LoanAccountId     uint64    `gorm:"column:loan_account_id" json:"loanAccountId"`
	StartDate         time.Time `gorm:"column:start_date" json:"startDate"`
	EndDate           time.Time `gorm:"column:end_date" json:"endDate"`
	WeekCount         int       `gorm:"column:week_count" json:"weekCount"`
	InstallmentAmount string    `gorm:"column:installment_amount" json:"installmentAmount"`
	CustomerDisplayId string    `gorm:"column:customer_display_id" json:"customerDisplayId"`
	LoanDisplayId     string    `gorm:"column:loan_display_id" json:"loanDisplayId"`
	IsCustomerActive  bool      `gorm:"column:is_customer_active" json:"isCustomerActive"`
	LoanState         LoanState `gorm:"column:loan_state" json:"loanState"`
	LoanId            uint64    `gorm:"column:loan_id" json:"loanId"`
	CustomerId        uint64    `gorm:"column:customer_id" json:"customerId"`
}

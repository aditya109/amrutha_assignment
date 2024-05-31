package models

import "time"

type BillingSchedule struct {
	Id                uint64       `json:"id" gorm:"column:id;primaryKey"`
	CreatedAt         time.Time    `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt         time.Time    `gorm:"column:updated_at" json:"updatedAt"`
	LoanAccountId     uint64       `gorm:"column:loan_account_id" json:"loanAccountId"`
	LoanAccount       *LoanAccount `gorm:"foreignKey:LoanAccountId" json:"loanAccount"`
	StartDate         time.Time    `gorm:"column:start_date" json:"startDate"`
	EndDate           time.Time    `gorm:"column:end_date" json:"endDate"`
	WeekCount         int          `gorm:"column:week_count" json:"weekCount"`
	InstallmentAmount string       `gorm:"column:installment_amount" json:"installmentAmount"`
}

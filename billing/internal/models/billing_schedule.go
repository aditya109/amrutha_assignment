package models

import "time"

type BillingSchedule struct {
	Id                uint64    `json:"id" gorm:"column:id;primaryKey"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt         time.Time `gorm:"column:updated_at" json:"updatedAt"`
	LoanId            int       `gorm:"column:loan_id" json:"loanId"`
	Loan              Loan      `gorm:"foreignKey:LoanId" json:"loan"`
	StartDate         time.Time `gorm:"column:start_date" json:"startDate"`
	EndDate           time.Time `gorm:"column:end_date" json:"endDate"`
	WeekCount         int       `gorm:"column:week_count" json:"weekCount"`
	InstallmentAmount string    `gorm:"column:installment_amount" json:"installmentAmount"`
}

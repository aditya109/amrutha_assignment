package models

import "time"

type LoanAccount struct {
	Id                     uint64     `json:"-" gorm:"column:id;primaryKey"`
	CreatedAt              *time.Time `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt              *time.Time `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	LoanId                 int        `gorm:"column:loan_id" json:"loanId,omitempty"`
	Loan                   *Loan      `gorm:"foreignKey:loan_id" json:"-"`
	PayablePrincipalAmount string     `gorm:"column:payable_principal_amount" json:"payablePrincipalAmount"`
	AccruedInterest        string     `gorm:"column:accrued_interest" json:"accruedInterest"`
	TotalPayableAmount     string     `gorm:"column:total_payable_amount" json:"totalPayableAmount"`
	TotalPaidAmount        string     `gorm:"column:total_paid_amount" json:"totalPaidAmount"`
	OutstandingAmount      string     `gorm:"column:outstanding_amount" json:"outstandingAmount"`
	InstallmentAmount      string     `gorm:"column:installment_amount" json:"installmentAmount"`
	DisplayId              string     `json:"displayId" gorm:"column:display_id"`
}

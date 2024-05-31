package models

import "time"

type Payment struct {
	Id            uint64       `json:"id" gorm:"column:id;primaryKey"`
	CreatedAt     time.Time    `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt     time.Time    `gorm:"column:updated_at" json:"updatedAt"`
	CustomerId    uint64       `gorm:"column:customer_id" json:"customerId"`
	LoanAccountId uint64       `gorm:"column:loan_account_id" json:"loanAccountId"`
	LoanAccount   *LoanAccount `gorm:"foreignKey:LoanAccountId" json:"loanAccount"`
	Customer      *Customer    `gorm:"foreignKey:CustomerId"`
	Amount        float64      `gorm:"column:amount" json:"amount"`
	IsAccepted    bool         `gorm:"column:is_accepted" json:"isAccepted"`
}

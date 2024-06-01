package models

import "time"

type Payment struct {
	Id                           uint64       `json:"-" gorm:"column:id;primaryKey"`
	CreatedAt                    *time.Time   `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt                    *time.Time   `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	CustomerId                   uint64       `gorm:"column:customer_id" json:"-"`
	LoanAccountId                uint64       `gorm:"column:loan_account_id" json:"-"`
	LoanAccount                  *LoanAccount `gorm:"foreignKey:LoanAccountId" json:"-"`
	Customer                     *Customer    `gorm:"foreignKey:CustomerId" json:"-"`
	PaidAmount                   string       `gorm:"column:paid_amount" json:"paidAmount"`
	PaymentDisplayId             string       `gorm:"column:payment_display_id" json:"paymentDisplayId"`
	ClientTransactionReferenceId string       `gorm:"column:client_transaction_reference_id" json:"clientTransactionReferenceId"`
	DateOfTransaction            time.Time    `gorm:"column:date_of_transaction" json:"dateOfTransaction"`
	IsAccepted                   bool         `gorm:"column:is_accepted" json:"-"`
	ScheduleId                   uint64       `gorm:"column:schedule_id" json:"-"`
}

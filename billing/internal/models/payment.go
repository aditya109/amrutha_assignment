package models

import "time"

type Payment struct {
	Id         uint64    `json:"id" gorm:"column:id;primaryKey"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updatedAt"`
	Loan       Loan
	Customer   Customer
	Amount     float64 `gorm:"column:amount" json:"amount"`
	IsAccepted bool    `gorm:"column:is_accepted" json:"isAccepted"`
}

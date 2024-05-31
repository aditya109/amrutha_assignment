package models

import (
	"time"
)

type Customer struct {
	ID        uint          `gorm:"column:id;primaryKey" json:"id,omitempty" `
	Name      string        `gorm:"column:name" json:"name"`
	Address   string        `gorm:"column:address" json:"address"`
	CreatedAt *time.Time    `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt *time.Time    `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	IsActive  bool          `gorm:"column:is_active" json:"isActive"`
	Type      CustomerState `gorm:"column:typ" json:"type,omitempty"`
	DisplayId string        `json:"displayId" gorm:"column:display_id"`
}

type CustomerState string

const (
	DelinquentCustomerState CustomerState = "DELINQUENT"
	RegularCustomerState    CustomerState = "REGULAR"
)

package models

import "time"

type Customer struct {
	Id        uint      `gorm:"column:id,primaryKey"`
	Name      string    `gorm:"column:name"`
	Address   string    `gorm:"column:address"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	IsActive  bool      `gorm:"column:is_active"`
}

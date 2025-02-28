package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID     uint    `gorm:"not null"`
	ProductID  uint    `gorm:"not null"`
	Quantity   int     `gorm:"not null"`
	TotalPrice float64 `gorm:"type:decimal(10,2);not null"`
	Status     string  `gorm:"type:enum('pending', 'completed', 'cancelled');default:'pending'"`
}

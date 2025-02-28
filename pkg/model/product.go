package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"type:decimal(10,2);not null"`
	Stock       int     `gorm:"not null"`
}

package model

import "gorm.io/gorm"

type Payment struct {
    gorm.Model
    OrderID uint    `gorm:"not null"`
    Amount  float64 `gorm:"type:decimal(10,2);not null"`
    Status  string  `gorm:"type:enum('pending', 'completed', 'failed');default:'pending'"`
}
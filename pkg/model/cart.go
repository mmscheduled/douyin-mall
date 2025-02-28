package model

import "gorm.io/gorm"

type CartItem struct {
    gorm.Model
    UserID    uint `gorm:"not null"`
    ProductID uint `gorm:"not null"`
    Quantity  int  `gorm:"not null"`
}
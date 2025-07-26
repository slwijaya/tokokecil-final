package model

import "time"

type Product struct {
	ID        uint    `gorm:"primaryKey"`
	Name      string  `gorm:"type:varchar(100);not null"`
	Price     float64 `gorm:"type:numeric(12,2);not null"`
	Stock     int     `gorm:"default:100"` // <-- TAMBAH FIELD STOCK
	CreatedAt time.Time
	UpdatedAt time.Time
}

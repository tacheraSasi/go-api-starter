package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"not null" json:"name"`
	Email       string         `gorm:"not null;uniqueIndex" json:"email"`
	Phone       string         `json:"phone"`
	Address     string         `gorm:"type:text" json:"address"`
	City        string         `json:"city"`
	State       string         `json:"state"`
	ZipCode     string         `json:"zip_code"`
	Country     string         `json:"country"`
	Invoices    []Invoice      `gorm:"foreignKey:CustomerID" json:"invoices,omitempty"`
}
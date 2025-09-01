package models

import (
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	InvoiceNumber string        `gorm:"uniqueIndex;not null" json:"invoice_number"`
	IssueDate    time.Time      `gorm:"not null" json:"issue_date"`
	DueDate      time.Time      `gorm:"not null" json:"due_date"`
	Status       string         `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft, sent, paid, cancelled
	CustomerID   uint           `gorm:"not null" json:"customer_id"`
	Customer     Customer       `json:"customer"`
	Items        []InvoiceItem  `gorm:"foreignKey:InvoiceID" json:"items"`
	Subtotal     float64        `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	TaxAmount    float64        `gorm:"type:decimal(10,2);default:0" json:"tax_amount"`
	Total        float64        `gorm:"type:decimal(10,2);not null" json:"total"`
	Notes        string         `gorm:"type:text" json:"notes"`
}

type InvoiceItem struct {
	ID          uint    `gorm:"primarykey" json:"id"`
	InvoiceID   uint    `gorm:"not null" json:"invoice_id"`
	Description string  `gorm:"not null" json:"description"`
	Quantity    int     `gorm:"not null" json:"quantity"`
	UnitPrice   float64 `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Total       float64 `gorm:"type:decimal(10,2);not null" json:"total"`
}
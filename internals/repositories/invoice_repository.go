package repositories

import (
	"github.com/tachRoutine/invoice-creator-api/internals/models"
	"gorm.io/gorm"
)

type InvoiceRepository interface {
	Create(invoice *models.Invoice) error
	FindByID(id uint) (*models.Invoice, error)
	FindAll(page, limit int) ([]models.Invoice, int64, error)
	Update(invoice *models.Invoice) error
	Delete(id uint) error
	FindByCustomerID(customerID uint, page, limit int) ([]models.Invoice, int64, error)
	FindByInvoiceNumber(invoiceNumber string) (*models.Invoice, error)
}

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) Create(invoice *models.Invoice) error {
	return r.db.Create(invoice).Error
}

func (r *invoiceRepository) FindByID(id uint) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.db.Preload("Customer").Preload("Items").First(&invoice, id).Error
	return &invoice, err
}

func (r *invoiceRepository) FindAll(page, limit int) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	offset := (page - 1) * limit

	err := r.db.Preload("Customer").
		Model(&models.Invoice{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&invoices).Error

	return invoices, total, err
}

func (r *invoiceRepository) Update(invoice *models.Invoice) error {
	return r.db.Save(invoice).Error
}

func (r *invoiceRepository) Delete(id uint) error {
	return r.db.Delete(&models.Invoice{}, id).Error
}

func (r *invoiceRepository) FindByCustomerID(customerID uint, page, limit int) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	offset := (page - 1) * limit

	err := r.db.Preload("Customer").
		Where("customer_id = ?", customerID).
		Model(&models.Invoice{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&invoices).Error

	return invoices, total, err
}

func (r *invoiceRepository) FindByInvoiceNumber(invoiceNumber string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.db.Preload("Customer").Preload("Items").
		Where("invoice_number = ?", invoiceNumber).
		First(&invoice).Error
	return &invoice, err
}
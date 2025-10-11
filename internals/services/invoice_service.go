package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/tacheraSasi/go-api-starter/internals/models"
	"github.com/tacheraSasi/go-api-starter/internals/repositories"
	"github.com/tacheraSasi/go-api-starter/internals/utils"
)

type InvoiceService interface {
	CreateInvoice(invoice *models.Invoice) error
	GetInvoiceByID(id uint) (*models.Invoice, error)
	GetAllInvoices(page, limit int) ([]models.Invoice, *utils.Pagination, error)
	UpdateInvoice(id uint, updatedInvoice *models.Invoice) error
	DeleteInvoice(id uint) error
	GetInvoicesByCustomerID(customerID uint, page, limit int) ([]models.Invoice, *utils.Pagination, error)
	GenerateInvoiceNumber() (string, error)
}

type invoiceService struct {
	repo repositories.InvoiceRepository
}

func NewInvoiceService(repo repositories.InvoiceRepository) InvoiceService {
	return &invoiceService{repo: repo}
}

func (s *invoiceService) CreateInvoice(invoice *models.Invoice) error {
	// Generate invoice number if not provided
	if invoice.InvoiceNumber == "" {
		invoiceNumber, err := s.GenerateInvoiceNumber()
		if err != nil {
			return err
		}
		invoice.InvoiceNumber = invoiceNumber
	}

	// Calculate totals
	s.calculateInvoiceTotals(invoice)

	// Set default status if not provided
	if invoice.Status == "" {
		invoice.Status = "draft"
	}

	return s.repo.Create(invoice)
}

func (s *invoiceService) GetInvoiceByID(id uint) (*models.Invoice, error) {
	return s.repo.FindByID(id)
}

func (s *invoiceService) GetAllInvoices(page, limit int) ([]models.Invoice, *utils.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	invoices, total, err := s.repo.FindAll(page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := utils.NewPagination(page, limit, total)
	return invoices, pagination, nil
}

func (s *invoiceService) UpdateInvoice(id uint, updatedInvoice *models.Invoice) error {
	existingInvoice, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("invoice not found")
	}

	// Update fields
	existingInvoice.IssueDate = updatedInvoice.IssueDate
	existingInvoice.DueDate = updatedInvoice.DueDate
	existingInvoice.Status = updatedInvoice.Status
	existingInvoice.CustomerID = updatedInvoice.CustomerID
	existingInvoice.Items = updatedInvoice.Items
	existingInvoice.Notes = updatedInvoice.Notes

	// Recalculate totals
	s.calculateInvoiceTotals(existingInvoice)

	return s.repo.Update(existingInvoice)
}

func (s *invoiceService) DeleteInvoice(id uint) error {
	return s.repo.Delete(id)
}

func (s *invoiceService) GetInvoicesByCustomerID(customerID uint, page, limit int) ([]models.Invoice, *utils.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	invoices, total, err := s.repo.FindByCustomerID(customerID, page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := utils.NewPagination(page, limit, total)
	return invoices, pagination, nil
}

func (s *invoiceService) GenerateInvoiceNumber() (string, error) {
	year := time.Now().Year()
	month := time.Now().Format("01")
	
	// In a real application, you might want to check the last invoice number
	// and increment it. For simplicity, we're using timestamp here.
	timestamp := time.Now().Unix()
	
	return fmt.Sprintf("INV-%d-%s-%d", year, month, timestamp), nil
}

func (s *invoiceService) calculateInvoiceTotals(invoice *models.Invoice) {
	subtotal := 0.0
	for i := range invoice.Items {
		item := &invoice.Items[i]
		item.Total = item.UnitPrice * float64(item.Quantity)
		subtotal += item.Total
	}
	
	invoice.Subtotal = subtotal
	// For simplicity, tax is 10% of subtotal
	invoice.TaxAmount = subtotal * 0.1
	invoice.Total = subtotal + invoice.TaxAmount
}
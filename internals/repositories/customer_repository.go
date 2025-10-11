package repositories

import (
	"github.com/tacheraSasi/go-api-starter/internals/models"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(customer *models.Customer) error
	FindByID(id uint) (*models.Customer, error)
	FindAll(page, limit int) ([]models.Customer, int64, error)
	Update(customer *models.Customer) error
	Delete(id uint) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) FindByID(id uint) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.First(&customer, id).Error
	return &customer, err
}

func (r *customerRepository) FindAll(page, limit int) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&models.Customer{}).Count(&total).Limit(limit).Offset(offset).Order("created_at DESC").Find(&customers).Error

	return customers, total, err
}

func (r *customerRepository) Update(customer *models.Customer) error {
	return r.db.Save(customer).Error
}

func (r *customerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Customer{}, id).Error
}

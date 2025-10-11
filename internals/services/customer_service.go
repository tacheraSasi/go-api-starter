package services

import (
	"errors"

	"github.com/tacheraSasi/go-api-starter/internals/models"
	"github.com/tacheraSasi/go-api-starter/internals/repositories"
	"github.com/tacheraSasi/go-api-starter/internals/utils"
)

type CustomerService interface {
	CreateCustomer(customer *models.Customer) error
	GetCustomerByID(id uint) (*models.Customer, error)
	GetAllCustomers(page, limit int) ([]models.Customer, *utils.Pagination, error)
	UpdateCustomer(id uint, updatedCustomer *models.Customer) error
	DeleteCustomer(id uint) error
}

type customerService struct {
	repo repositories.CustomerRepository
}

func NewCustomerService(repo repositories.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

func (s *customerService) CreateCustomer(customer *models.Customer) error {
	return s.repo.Create(customer)
}

func (s *customerService) GetCustomerByID(id uint) (*models.Customer, error) {
	return s.repo.FindByID(id)
}

func (s *customerService) GetAllCustomers(page, limit int) ([]models.Customer, *utils.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	customers, total, err := s.repo.FindAll(page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := utils.NewPagination(page, limit, total)
	return customers, pagination, nil
}

func (s *customerService) UpdateCustomer(id uint, updatedCustomer *models.Customer) error {
	existingCustomer, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("customer not found")
	}

	existingCustomer.Name = updatedCustomer.Name
	existingCustomer.Email = updatedCustomer.Email
	existingCustomer.Phone = updatedCustomer.Phone
	existingCustomer.Address = updatedCustomer.Address

	return s.repo.Update(existingCustomer)
}

func (s *customerService) DeleteCustomer(id uint) error {
	return s.repo.Delete(id)
}

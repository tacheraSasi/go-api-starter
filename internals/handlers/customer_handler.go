package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tacheraSasi/go-api-starter/internals/models"
	"github.com/tacheraSasi/go-api-starter/internals/services"
	"github.com/tacheraSasi/go-api-starter/internals/utils"
)

type CustomerHandler struct {
	service services.CustomerService
}

func NewCustomerHandler(service services.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	if err := h.service.CreateCustomer(&customer); err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Failed to create customer: "+err.Error())
		return
	}

	utils.APISuccess(c, http.StatusCreated, customer)
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	customer, err := h.service.GetCustomerByID(uint(id))
	if err != nil {
		utils.APIError(c, http.StatusNotFound, "Customer not found")
		return
	}

	utils.APISuccess(c, http.StatusOK, customer)
}

func (h *CustomerHandler) ListCustomers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	customers, pagination, err := h.service.GetAllCustomers(page, limit)
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Failed to fetch customers")
		return
	}

	response := map[string]interface{}{
		"customers":  customers,
		"pagination": pagination,
	}

	utils.APISuccess(c, http.StatusOK, response)
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	if err := h.service.UpdateCustomer(uint(id), &customer); err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Failed to update customer: "+err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, customer)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	if err := h.service.DeleteCustomer(uint(id)); err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Failed to delete customer: "+err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

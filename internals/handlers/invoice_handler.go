package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tacheraSasi/go-api-starter/internals/models"
	"github.com/tacheraSasi/go-api-starter/internals/services"
	"github.com/tacheraSasi/go-api-starter/internals/utils"
)

type InvoiceHandler struct {
	service services.InvoiceService
}

func NewInvoiceHandler(service services.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{service: service}
}

func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	var invoice models.Invoice
	if err := c.ShouldBindJSON(&invoice); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	if err := h.service.CreateInvoice(&invoice); err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Failed to create invoice: "+err.Error())
		return
	}

	utils.APISuccess(c, http.StatusCreated, invoice)
}

func (h *InvoiceHandler) GetInvoice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	invoice, err := h.service.GetInvoiceByID(uint(id))
	if err != nil {
		utils.APIError(c, http.StatusNotFound, "Invoice not found")
		return
	}

	utils.APISuccess(c, http.StatusOK, invoice)
}

func (h *InvoiceHandler) ListInvoices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	invoices, pagination, err := h.service.GetAllInvoices(page, limit)
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Failed to fetch invoices")
		return
	}

	response := map[string]interface{}{
		"invoices":   invoices,
		"pagination": pagination,
	}

	utils.APISuccess(c, http.StatusOK, response)
}

func (h *InvoiceHandler) UpdateInvoice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	var invoice models.Invoice
	if err := c.ShouldBindJSON(&invoice); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	if err := h.service.UpdateInvoice(uint(id), &invoice); err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Failed to update invoice: "+err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, invoice)
}

func (h *InvoiceHandler) DeleteInvoice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	if err := h.service.DeleteInvoice(uint(id)); err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Failed to delete invoice: "+err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
}
package handlers

import (
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tachRoutine/invoice-creator-api/internals/dtos"
	"github.com/tachRoutine/invoice-creator-api/internals/services"
	"github.com/tachRoutine/invoice-creator-api/pkg/styles"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// ValidateRequest validates the request body against the provided struct
func (h *AuthHandler) ValidateRequest(c *gin.Context, obj any) {
	dtos.Validate(c, obj)
}

func (h *AuthHandler) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var requestBody = c.Request.Body
	defer requestBody.Close()
	h.ValidateRequest(c, requestBody)
	bodyBytes, err := io.ReadAll(requestBody)
	if err != nil {
		log.Println("Failed to read request body:", err)
		return
	}
	log.Println(styles.Request.Render(string(bodyBytes)))

}

func (h *AuthHandler) Login(c *gin.Context) {

}

func (h *AuthHandler) Logout(c *gin.Context) {

}
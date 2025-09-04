package handlers

import (
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tachRoutine/invoice-creator-api/internals/dtos"
	"github.com/tachRoutine/invoice-creator-api/internals/models"
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
	var reqDto dtos.RegisterRequest
	var requestBody = c.Request.Body
	defer requestBody.Close()
	h.ValidateRequest(c, &reqDto)
	if c.IsAborted() {
		return
	}
	bodyBytes, err := io.ReadAll(requestBody)
	if err != nil {
		log.Println("Failed to read request body:", err)
		return
	}
	log.Println(styles.Request.Render(string(bodyBytes)))
	err = h.service.Register(&models.User{
		Email:    reqDto.Email,
		Password: reqDto.Password,
		Name:     reqDto.Name,
	})
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "Registration successful",
		"user":    user,
	})

}

func (h *AuthHandler) Login(c *gin.Context) {
	var reqDto dtos.LoginRequest
	var requestBody = c.Request.Body
	defer requestBody.Close()
	h.ValidateRequest(c, &reqDto)
	if c.IsAborted(){
		return
	}
	bodyBytes, err := io.ReadAll(requestBody)
	if err != nil {
		log.Println("Failed to read request body:", err)
		return
	}
	log.Println(styles.Request.Render(string(bodyBytes)))

	user, err := h.service.Login(reqDto.Email, reqDto.Password)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Login successful",
		"user":    user,
	})

}

func (h *AuthHandler) Logout(c *gin.Context) {

}

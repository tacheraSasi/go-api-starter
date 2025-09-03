package handlers

import (
	"io"
	"log"

	"github.com/gin-gonic/gin"
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

func (h *AuthHandler) Register(c *gin.Context) {
	var requestBody = c.Request.Body
	defer requestBody.Close()
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
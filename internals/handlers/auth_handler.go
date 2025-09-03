package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tachRoutine/invoice-creator-api/internals/services"
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

}

func (h *AuthHandler) Login(c *gin.Context) {

}

func (h *AuthHandler) Logout(c *gin.Context) {

}
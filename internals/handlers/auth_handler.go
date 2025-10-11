package handlers

import (
	"io"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tacheraSasi/go-api-starter/internals/config"
	"github.com/tacheraSasi/go-api-starter/internals/dtos"
	"github.com/tacheraSasi/go-api-starter/internals/models"
	"github.com/tacheraSasi/go-api-starter/internals/services"
	"github.com/tacheraSasi/go-api-starter/pkg/jwt"
	"github.com/tacheraSasi/go-api-starter/pkg/styles"
)

type AuthHandler struct {
	service services.AuthService
	cfg     *config.Config
}

func NewAuthHandler(service services.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		service: service,
		cfg:     cfg,
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
	user, err := h.service.GetUserByEmail(reqDto.Email)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to retrieve user after registration",
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
	token, err := jwt.GenerateToken(user, []byte(h.cfg.JWTSecret), h.cfg.JWTExpiresIn)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Login successful",
		"user":    user,
		"token":   token,
	})

}

func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Authorization header is missing"})
		return
	}

	tokenString := strings.Split(authHeader, " ")[1]

	claims, err := jwt.ValidateToken(tokenString, []byte(h.cfg.JWTSecret))
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid token"})
		return
	}

	expiresAt := claims.ExpiresAt.Time

	if err := h.service.Logout(tokenString, expiresAt); err != nil {
		c.JSON(500, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(200, gin.H{"message": "Logout successful"})
}

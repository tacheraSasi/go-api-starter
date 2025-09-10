package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tachRoutine/invoice-creator-api/internals/config"
	"github.com/tachRoutine/invoice-creator-api/internals/handlers"
	"github.com/tachRoutine/invoice-creator-api/internals/middlewares"
	"github.com/tachRoutine/invoice-creator-api/internals/models"
	"github.com/tachRoutine/invoice-creator-api/internals/repositories"
	"github.com/tachRoutine/invoice-creator-api/internals/services"
	"github.com/tachRoutine/invoice-creator-api/pkg/database"
	"github.com/tachRoutine/invoice-creator-api/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	logger, logErr := logger.NewLogger(cfg.LogFilePath)
	if logErr != nil {
		log.Fatal("Failed to initialize logger:", logErr)
	}

	// Connect to database 
	err := database.Connect(
		database.DBConfig{
			Type:	 cfg.DBType,
			FilePath:	 cfg.DBPath,
		},
	)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Auto migrate models
	err = database.AutoMigrate(
		&models.User{},
		&models.Customer{},
		&models.Invoice{},
		&models.InvoiceItem{},
		&models.BlacklistedToken{},
	)
	if err != nil {
		log.Fatal("Auto migration failed:", err)
	}

	// repositories
	userRepo := repositories.NewUserRepository(database.GetDB())
	customerRepo := repositories.NewCustomerRepository(database.GetDB())
	invoiceRepo := repositories.NewInvoiceRepository(database.GetDB())
	tokenRepo := repositories.NewTokenRepository(database.GetDB())

	// services
	tokenService := services.NewTokenService(tokenRepo)
	authService := services.NewAuthService(userRepo, tokenService)
	customerService := services.NewCustomerService(customerRepo)
	invoiceService := services.NewInvoiceService(invoiceRepo)

	// handlers
	authHandler := handlers.NewAuthHandler(authService, cfg)
	customerHandler := handlers.NewCustomerHandler(customerService)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService)

	// Setup router
	r := gin.Default()

	// Global middlewares
	r.Use(middlewares.LoggingMiddleware(logger.Logger))
	r.Use(middlewares.CORSMiddleware("*"))

	// Public routes
	public := r.Group("/api/v1")
	{
		public.POST("/login", authHandler.Login)
		public.POST("/register", authHandler.Register)
	}

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middlewares.AuthMiddleware(tokenService, []byte(cfg.JWTSecret)))
	{
		protected.POST("/logout", authHandler.Logout)
		// Customer routes
		protected.GET("/customers", customerHandler.ListCustomers)
		protected.GET("/customers/:id", customerHandler.GetCustomer)
		protected.POST("/customers", customerHandler.CreateCustomer)
		protected.PUT("/customers/:id", customerHandler.UpdateCustomer)
		protected.DELETE("/customers/:id", customerHandler.DeleteCustomer)

		// Invoice routes
		protected.GET("/invoices", invoiceHandler.ListInvoices)
		protected.GET("/invoices/:id", invoiceHandler.GetInvoice)
		protected.POST("/invoices", invoiceHandler.CreateInvoice)
		protected.PUT("/invoices/:id", invoiceHandler.UpdateInvoice)
		protected.DELETE("/invoices/:id", invoiceHandler.DeleteInvoice)
	}

	// Admin routes
	admin := r.Group("/api/v1/admin")
	admin.Use(middlewares.AuthMiddleware(tokenService, []byte(cfg.JWTSecret)), middlewares.AdminMiddleware())
	{
		// TODO: Add admin specific routes here
	}

	// Start server
	log.Printf("Server starting on :%s", cfg.ServerPort)
	log.Fatal(r.Run(":" + cfg.ServerPort))
}
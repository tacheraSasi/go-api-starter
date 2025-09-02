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
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	logger.InitLogger()

	// Connect to database
	err := database.Connect(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
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
	)
	if err != nil {
		log.Fatal("Auto migration failed:", err)
	}

	// repositories
	userRepo := repositories.NewUserRepository(database.GetDB())
	customerRepo := repositories.NewCustomerRepository(database.GetDB())
	invoiceRepo := repositories.NewInvoiceRepository(database.GetDB())

	// services
	authService := services.NewAuthService(userRepo)
	customerService := services.NewCustomerService(customerRepo)
	invoiceService := services.NewInvoiceService(invoiceRepo)

	// handlers
	authHandler := handlers.NewAuthHandler(authService)
	customerHandler := handlers.NewCustomerHandler(customerService)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService)

	// Setup router
	r := gin.Default()

	// Global middlewares
	r.Use(middlewares.LoggingMiddleware())
	r.Use(middlewares.CORSMiddleware())

	// Public routes
	public := r.Group("/api/v1")
	{
		public.POST("/login", authHandler.Login)
		public.POST("/register", authHandler.Register)
	}

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middlewares.AuthMiddleware())
	{
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
	admin.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		// Add admin-specific routes here
	}

	// Start server
	log.Printf("Server starting on :%s", cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}
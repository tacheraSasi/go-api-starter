package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/tacheraSasi/go-api-starter/internals/config"
	"github.com/tacheraSasi/go-api-starter/internals/handlers"
	"github.com/tacheraSasi/go-api-starter/internals/middlewares"
	"github.com/tacheraSasi/go-api-starter/internals/models"
	"github.com/tacheraSasi/go-api-starter/internals/repositories"
	"github.com/tacheraSasi/go-api-starter/internals/services"
	"github.com/tacheraSasi/go-api-starter/pkg/database"
	"github.com/tacheraSasi/go-api-starter/pkg/logger"
	"github.com/tacheraSasi/go-api-starter/views"
)

func main() {
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatal("Invalid configuration:", err)
	}

	gin.SetMode(cfg.GINMode)

	logger, logErr := logger.NewLogger(cfg.LogFilePath)
	if logErr != nil {
		log.Fatal("Failed to initialize logger:", logErr)
	}
	if cfg.JWTSecret == "secret" {
		logger.Logger.Warn("Using default JWT secret. Set JWT_SECRET in production environments")
	}

	// Connect to database
	err := database.Connect(
		database.DBConfig{
			Type:     cfg.DBType,
			Host:     cfg.DBHost,
			Port:     cfg.DBPort,
			User:     cfg.DBUser,
			Password: cfg.DBPassword,
			DBName:   cfg.DBName,
			SSLMode:  "disable",
			FilePath: cfg.DBPath,
		},
	)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer func() {
		if closeErr := database.Close(); closeErr != nil {
			logger.Logger.WithError(closeErr).Error("failed to close database")
		}
	}()

	// Auto migrate models
	err = database.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
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
	roleRepo := repositories.NewRoleRepository(database.GetDB())
	permissionRepo := repositories.NewPermissionRepository(database.GetDB())
	customerRepo := repositories.NewCustomerRepository(database.GetDB())
	invoiceRepo := repositories.NewInvoiceRepository(database.GetDB())
	tokenRepo := repositories.NewTokenRepository(database.GetDB())

	// services
	permissionService := services.NewPermissionService(permissionRepo)
	roleService := services.NewRoleService(roleRepo, permissionRepo)
	userService := services.NewUserService(userRepo, roleRepo)
	tokenService := services.NewTokenService(tokenRepo)
	authService := services.NewAuthService(userRepo, tokenService)
	customerService := services.NewCustomerService(customerRepo)
	invoiceService := services.NewInvoiceService(invoiceRepo)

	// Initialize default roles and permissions
	if err := permissionService.InitializeDefaultPermissions(); err != nil {
		log.Printf("Warning: Failed to initialize default permissions: %v", err)
	}
	if err := roleService.InitializeDefaultRoles(); err != nil {
		log.Printf("Warning: Failed to initialize default roles: %v", err)
	}

	// handlers
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService, cfg)
	userHandler := handlers.NewUserHandler(userService)
	roleHandler := handlers.NewRoleHandler(roleService)
	permissionHandler := handlers.NewPermissionHandler(permissionService)
	customerHandler := handlers.NewCustomerHandler(customerService)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService)

	// Setup router
	r := gin.New()
	r.Use(gin.Recovery())
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal("Failed to set trusted proxies:", err)
	}

	// Global middlewares
	r.Use(middlewares.LoggingMiddleware(logger.Logger))
	r.Use(middlewares.CORSMiddleware(cfg.CORSOrigins...))

	r.GET("/health", healthHandler.HealthCheck)
	r.GET("/health/ready", healthHandler.ReadinessCheck)

	//Home page
	r.GET("/", func(c *gin.Context) {
		templ.Handler(views.Home(views.HomeProps{AppName: "Go API Starter"})).ServeHTTP(c.Writer, c.Request)
	})

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

		// User routes
		protected.GET("/users/:id", userHandler.GetUser)
		protected.PUT("/users/:id", userHandler.UpdateUser)
		protected.PUT("/users/:id/password", userHandler.UpdateUserPassword)
		protected.GET("/users/:id/roles", userHandler.GetUserRoles)
		protected.GET("/users/:id/permissions/:resource/:action", userHandler.CheckUserPermission)

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
		// User management
		admin.GET("/users", userHandler.ListUsers)
		admin.DELETE("/users/:id", userHandler.DeleteUser)
		admin.POST("/users/:id/roles/:roleId", userHandler.AddRoleToUser)
		admin.DELETE("/users/:id/roles/:roleId", userHandler.RemoveRoleFromUser)

		// Role management
		admin.POST("/roles", roleHandler.CreateRole)
		admin.GET("/roles", roleHandler.ListRoles)
		admin.GET("/roles/:id", roleHandler.GetRole)
		admin.PUT("/roles/:id", roleHandler.UpdateRole)
		admin.DELETE("/roles/:id", roleHandler.DeleteRole)
		admin.POST("/roles/:id/permissions/:permissionId", roleHandler.AddPermissionToRole)
		admin.DELETE("/roles/:id/permissions/:permissionId", roleHandler.RemovePermissionFromRole)

		// Permission management
		admin.POST("/permissions", permissionHandler.CreatePermission)
		admin.GET("/permissions", permissionHandler.ListPermissions)
		admin.GET("/permissions/:id", permissionHandler.GetPermission)
		admin.PUT("/permissions/:id", permissionHandler.UpdatePermission)
		admin.DELETE("/permissions/:id", permissionHandler.DeletePermission)
		admin.GET("/permissions/resources", permissionHandler.GetAllResources)
		admin.GET("/permissions/resources/:resource/actions", permissionHandler.GetResourceActions)
	}

	// Start server
	server := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on :%s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}

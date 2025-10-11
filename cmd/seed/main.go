package main

import (
	"fmt"
	"log"

	"github.com/tacheraSasi/go-api-starter/internals/config"
	"github.com/tacheraSasi/go-api-starter/internals/models"
	"github.com/tacheraSasi/go-api-starter/internals/repositories"
	"github.com/tacheraSasi/go-api-starter/internals/services"
	"github.com/tacheraSasi/go-api-starter/pkg/database"
)

func main() {
	fmt.Println("🌱 Seeding database...")

	// Load config
	cfg := config.LoadConfig()

	// Connect to database
	err := database.Connect(database.DBConfig{
		Type:     cfg.DBType,
		FilePath: cfg.DBPath,
	})
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

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

	// Initialize repositories and services
	permissionRepo := repositories.NewPermissionRepository(database.GetDB())
	roleRepo := repositories.NewRoleRepository(database.GetDB())
	userRepo := repositories.NewUserRepository(database.GetDB())

	permissionService := services.NewPermissionService(permissionRepo)
	roleService := services.NewRoleService(roleRepo, permissionRepo)
	userService := services.NewUserService(userRepo, roleRepo)

	// Seed permissions
	fmt.Println("📋 Creating default permissions...")
	if err := permissionService.InitializeDefaultPermissions(); err != nil {
		log.Printf("Warning: Failed to create permissions: %v", err)
	} else {
		fmt.Println("✅ Default permissions created")
	}

	// Seed roles
	fmt.Println("👥 Creating default roles...")
	if err := roleService.InitializeDefaultRoles(); err != nil {
		log.Printf("Warning: Failed to create roles: %v", err)
	} else {
		fmt.Println("✅ Default roles created")
	}

	// Assign permissions to admin role
	fmt.Println("🔐 Assigning permissions to admin role...")
	adminRole, err := roleService.GetRoleByName(models.RoleAdmin)
	if err != nil {
		log.Printf("Warning: Could not find admin role: %v", err)
	} else {
		// Get all permissions
		permissions, err := permissionService.ListPermissions(1000, 0, "")
		if err != nil {
			log.Printf("Warning: Could not get permissions: %v", err)
		} else {
			for _, permission := range permissions {
				_ = roleService.AddPermissionToRole(adminRole.ID, permission.ID)
			}
			fmt.Println("✅ Admin permissions assigned")
		}
	}

	// Assign basic permissions to user role
	fmt.Println("📖 Assigning basic permissions to user role...")
	userRole, err := roleService.GetRoleByName(models.RoleUser)
	if err != nil {
		log.Printf("Warning: Could not find user role: %v", err)
	} else {
		// Basic read permissions
		basicPermissions := []struct {
			resource string
			action   string
		}{
			{models.ResourceCustomer, models.ActionRead},
			{models.ResourceCustomer, models.ActionList},
			{models.ResourceInvoice, models.ActionRead},
			{models.ResourceInvoice, models.ActionList},
			{models.ResourceUser, models.ActionRead}, // Own profile
		}

		for _, perm := range basicPermissions {
			permission, err := permissionRepo.GetByResourceAndAction(perm.resource, perm.action)
			if err == nil {
				_ = roleService.AddPermissionToRole(userRole.ID, permission.ID)
			}
		}
		fmt.Println("✅ User permissions assigned")
	}

	// Create admin user
	fmt.Println("🔑 Creating admin user...")
	adminUser, err := userService.CreateUser("Admin User", "admin@example.com", "admin123456")
	if err != nil {
		log.Printf("Warning: Could not create admin user: %v", err)
	} else {
		// Assign admin role
		if adminRole != nil {
			_ = userService.AddRoleToUser(fmt.Sprintf("%d", adminUser.ID), adminRole.ID)
		}
		fmt.Println("✅ Admin user created (email: admin@example.com, password: admin123456)")
	}

	// Create regular user
	fmt.Println("👤 Creating regular user...")
	regularUser, err := userService.CreateUser("John Doe", "user@example.com", "user123456")
	if err != nil {
		log.Printf("Warning: Could not create regular user: %v", err)
	} else {
		fmt.Println("✅ Regular user created (email: user@example.com, password: user123456)")
		_ = regularUser // Just to avoid unused variable warning
	}

	fmt.Println("\n🎉 Database seeding completed!")
	fmt.Println("\n📝 Summary:")
	fmt.Println("- Default permissions and roles created")
	fmt.Println("- Admin user: admin@example.com / admin123456")
	fmt.Println("- Regular user: user@example.com / user123456")
	fmt.Println("\n🚀 You can now start the API server with: make dev")
}
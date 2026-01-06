package main

import (
	"fmt"
	"log"

	"rires-be/config"
	_ "rires-be/docs" // Swagger docs
	"rires-be/internal/routes"
	"rires-be/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// @title Student RIRES Backend API
// @version 1.0
// @description API untuk sistem PKM (Program Kreativitas Mahasiswa) UMM
// @description Backend yang mendukung pengajuan, review, dan manajemen PKM

// @contact.name API Support
// @contact.email support@rires.com

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to main database rires
	if err := database.Connect(config.AppConfig.GetDSN()); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.CloseDB()

	// Connect to external databases (NEOMAA, NEOMAAREF, SIMPEG)
	if err := database.ConnectExternal(
		config.AppConfig.GetDSNNeomaa(),
		config.AppConfig.GetDSNNeomaaRef(),
		config.AppConfig.GetDSNSimpeg(),
	); err != nil {
		log.Fatal("Failed to connect to external databases:", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: config.AppConfig.AppName,
	})

	// Middleware
	app.Use(recover.New()) // Recover from panics
	app.Use(logger.New())  // Log requests
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Setup routes
	routes.Setup(app)

	// Start server
	port := fmt.Sprintf(":%s", config.AppConfig.AppPort)
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
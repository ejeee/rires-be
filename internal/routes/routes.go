package routes

import (
	"rires-be/internal/controllers"
	"rires-be/pkg/database"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Welcome route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"app":     "RIRES Backend API",
			"version": "1.0.0",
			"status":  "running",
		})
	})
	
	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		// Check database connection
		sqlDB, err := database.DB.DB()
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":   "error",
				"database": "disconnected",
				"message":  err.Error(),
			})
		}

		// Ping database
		if err := sqlDB.Ping(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":   "error",
				"database": "disconnected",
				"message":  err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status":   "ok",
			"database": "connected",
			"message":  "API is running",
		})
	})

	// API v1 routes
	api := app.Group("/api/v1")

	// Auth routes
	authController := controllers.NewAuthController()
	auth := api.Group("/auth")
	{
		auth.Post("/login/admin", authController.LoginAdmin)
		auth.Post("/login/mahasiswa", authController.LoginMahasiswa)
		auth.Post("/login/pegawai", authController.LoginPegawai)
	}

	// Protected routes (akan ditambahkan nanti dengan middleware JWT)
	// protected := api.Group("/", middleware.JWTProtected())
	// {
	//     protected.Get("/profile", controllers.GetProfile)
	// }
}
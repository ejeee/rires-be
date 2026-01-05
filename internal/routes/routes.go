package routes

import (
	"rires-be/internal/controllers"
	"rires-be/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func Setup(app *fiber.App) {
	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

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
		auth.Get("/me", authController.GetCurrentUser)

		// GET endpoints for browser testing (DEVELOPMENT ONLY - NOT SECURE!)
		auth.Get("/test/pegawai/:username/:password", func(c *fiber.Ctx) error {
			username := c.Params("username")
			password := c.Params("password")
			
			return c.JSON(fiber.Map{
				"warning": "This is for testing only - NOT SECURE!",
				"username": username,
				"password": password,
				"message": "Use POST /api/v1/auth/login/pegawai with JSON body for actual login",
			})
		})
	}

	// User Level routes (Master Data)
	userLevelController := controllers.NewUserLevelController()
	userLevels := api.Group("/user-levels")
	{
		userLevels.Get("/", userLevelController.GetList)
		userLevels.Get("/:id", userLevelController.GetByID)
		userLevels.Post("/", userLevelController.Create)
		userLevels.Put("/:id", userLevelController.Update)
		userLevels.Delete("/:id", userLevelController.Delete)
	}

		// Menu routes (Master Data)
	menuController := controllers.NewMenuController()
	menus := api.Group("/menus")
	{
		menus.Get("/", menuController.GetList)           // List flat
		menus.Get("/tree", menuController.GetTree)       // Tree structure
		menus.Get("/:id", menuController.GetByID)
		menus.Post("/", menuController.Create)
		menus.Put("/:id", menuController.Update)
		menus.Delete("/:id", menuController.Delete)
	}

	// Protected routes (akan ditambahkan nanti dengan middleware JWT)
	// protected := api.Group("/", middleware.JWTProtected())
	// {
	//     protected.Get("/profile", controllers.GetProfile)
	// }
}
package routes

import (
	"rires-be/internal/controllers"
	"rires-be/internal/middleware"
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

	// Database connection test endpoint
	// app.Get("/test-db", func(c *fiber.Ctx) error {
	// 	results := make(map[string]string)

	// 	// Test Main DB
	// 	if sqlDB, err := database.DB.DB(); err == nil {
	// 		if err := sqlDB.Ping(); err == nil {
	// 			results["main_db"] = "✅ Connected"
	// 		} else {
	// 			results["main_db"] = "❌ Ping failed: " + err.Error()
	// 		}
	// 	} else {
	// 		results["main_db"] = "❌ Connection failed: " + err.Error()
	// 	}

	// 	// Test NEOMAA
	// 	if database.DBNeomaa != nil {
	// 		if sqlDB, err := database.DBNeomaa.DB(); err == nil {
	// 			if err := sqlDB.Ping(); err == nil {
	// 				results["neomaa_db"] = "✅ Connected"
	// 			} else {
	// 				results["neomaa_db"] = "❌ Ping failed: " + err.Error()
	// 			}
	// 		} else {
	// 			results["neomaa_db"] = "❌ Connection failed: " + err.Error()
	// 		}
	// 	} else {
	// 		results["neomaa_db"] = "❌ Not initialized"
	// 	}

	// 	// Test NEOMAAREF
	// 	if database.DBNeomaaRef != nil {
	// 		if sqlDB, err := database.DBNeomaaRef.DB(); err == nil {
	// 			if err := sqlDB.Ping(); err == nil {
	// 				results["neomaaref_db"] = "✅ Connected"
	// 			} else {
	// 				results["neomaaref_db"] = "❌ Ping failed: " + err.Error()
	// 			}
	// 		} else {
	// 			results["neomaaref_db"] = "❌ Connection failed: " + err.Error()
	// 		}
	// 	} else {
	// 		results["neomaaref_db"] = "❌ Not initialized"
	// 	}

	// 	// Test SIMPEG
	// 	if database.DBSimpeg != nil {
	// 		if sqlDB, err := database.DBSimpeg.DB(); err == nil {
	// 			if err := sqlDB.Ping(); err == nil {
	// 				results["simpeg_db"] = "✅ Connected"
	// 			} else {
	// 				results["simpeg_db"] = "❌ Ping failed: " + err.Error()
	// 			}
	// 		} else {
	// 			results["simpeg_db"] = "❌ Connection failed: " + err.Error()
	// 		}
	// 	} else {
	// 		results["simpeg_db"] = "❌ Not initialized"
	// 	}

	// 	return c.JSON(fiber.Map{
	// 		"message": "Database Connection Test",
	// 		"results": results,
	// 	})
	// })

	// Test Helper endpoints (for development/testing)
	testHelperController := controllers.NewTestHelperController()
	test := app.Group("/test")
	{
		test.Get("/external-data", testHelperController.TestExternalData)
		test.Get("/code-generator", testHelperController.TestCodeGenerator)
		test.Get("/status-validator", testHelperController.TestStatusValidator)
		test.Post("/file-upload", testHelperController.TestFileUpload)
	}

	// API v1 routes
	api := app.Group("/api/v1")

	// ============================================
	// PUBLIC ROUTES (No authentication required)
	// ============================================

	// Auth routes (public)
	authController := controllers.NewAuthController()
	auth := api.Group("/auth")
	{
		auth.Post("/login/admin", authController.LoginAdmin)
		auth.Post("/login/mahasiswa", authController.LoginMahasiswa)
		auth.Post("/login/pegawai", authController.LoginPegawai)
	}

	// ============================================
	// PROTECTED ROUTES (JWT required)
	// ============================================
	protected := api.Group("/", middleware.JWTAuth())

	// Auth - Get current user (protected)
	authProtected := protected.Group("/auth")
	{
		authProtected.Get("/me", authController.GetCurrentUser)
	}

	// User Level routes (Admin only)
	userLevelController := controllers.NewUserLevelController()
	userLevels := protected.Group("/user-levels", middleware.RequireAdmin())
	{
		userLevels.Get("/", userLevelController.GetList)
		userLevels.Get("/:id", userLevelController.GetByID)
		userLevels.Post("/", userLevelController.Create)
		userLevels.Put("/:id", userLevelController.Update)
		userLevels.Delete("/:id", userLevelController.Delete)
	}

	// Menu routes (Admin only for CUD, All for Read)
	menuController := controllers.NewMenuController()
	menusPublic := protected.Group("/menus")
	{
		menusPublic.Get("/", menuController.GetList)       // All users can read
		menusPublic.Get("/tree", menuController.GetTree)   // All users can read
		menusPublic.Get("/:id", menuController.GetByID)    // All users can read
	}
	menusAdmin := protected.Group("/menus", middleware.RequireAdmin())
	{
		menusAdmin.Post("/", menuController.Create)
		menusAdmin.Put("/:id", menuController.Update)
		menusAdmin.Delete("/:id", menuController.Delete)
	}

	// Kategori PKM routes (Admin only for CUD, All for Read)
	kategoriPKMController := controllers.NewKategoriPKMController()
	kategoriPublic := protected.Group("/kategori-pkm")
	{
		kategoriPublic.Get("/", kategoriPKMController.GetList)      // All users can read
		kategoriPublic.Get("/:id", kategoriPKMController.GetByID)   // All users can read
	}
	kategoriAdmin := protected.Group("/kategori-pkm", middleware.RequireAdmin())
	{
		kategoriAdmin.Post("/", kategoriPKMController.Create)
		kategoriAdmin.Put("/:id", kategoriPKMController.Update)
		kategoriAdmin.Delete("/:id", kategoriPKMController.Delete)
	}

	// Status Review routes (Admin only for CUD, All for Read)
	statusReviewController := controllers.NewStatusReviewController()
	statusPublic := protected.Group("/status-review")
	{
		statusPublic.Get("/", statusReviewController.GetList)       // All users can read
		statusPublic.Get("/:id", statusReviewController.GetByID)    // All users can read
	}
	statusAdmin := protected.Group("/status-review", middleware.RequireAdmin())
	{
		statusAdmin.Post("/", statusReviewController.Create)
		statusAdmin.Put("/:id", statusReviewController.Update)
		statusAdmin.Delete("/:id", statusReviewController.Delete)
	}

	// Parameter Form routes (Admin only for CUD, All for Read)
	parameterFormController := controllers.NewParameterFormController()
	paramPublic := protected.Group("/parameter-form")
	{
		paramPublic.Get("/", parameterFormController.GetList)
		paramPublic.Get("/kategori/:kategori_id", parameterFormController.GetByKategori) // Important for mahasiswa
		paramPublic.Get("/:id", parameterFormController.GetByID)
	}
	paramAdmin := protected.Group("/parameter-form", middleware.RequireAdmin())
	{
		paramAdmin.Post("/", parameterFormController.Create)
		paramAdmin.Put("/:id", parameterFormController.Update)
		paramAdmin.Delete("/:id", parameterFormController.Delete)
	}

	// User Management routes (Admin only)
	userManagementController := controllers.NewUserManagementController()
	users := protected.Group("/users", middleware.RequireAdmin())
	{
		users.Get("/", userManagementController.GetList)
		users.Get("/:id", userManagementController.GetByID)
		users.Post("/", userManagementController.Create)
		users.Put("/:id", userManagementController.Update)
		users.Post("/:id/reset-password", userManagementController.ResetPassword)
		users.Delete("/:id", userManagementController.Delete)
	}

	// Tanggal Pendaftaran routes
	tanggalPendaftaranController := controllers.NewTanggalPendaftaranController()
	// Public endpoint - check if registration is open
	tanggalPublic := protected.Group("/tanggal-pendaftaran")
	{
		tanggalPublic.Get("/active", tanggalPendaftaranController.GetActive) // All authenticated users can check
	}
	// Admin only endpoints
	tanggalAdmin := protected.Group("/tanggal-pendaftaran", middleware.RequireAdmin())
	{
		tanggalAdmin.Get("/", tanggalPendaftaranController.GetList)
		tanggalAdmin.Get("/:id", tanggalPendaftaranController.GetByID)
		tanggalAdmin.Post("/", tanggalPendaftaranController.Create)
		tanggalAdmin.Put("/:id", tanggalPendaftaranController.Update)
		tanggalAdmin.Delete("/:id", tanggalPendaftaranController.Delete)
	}
}
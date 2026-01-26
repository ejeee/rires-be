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
	app.Get("/test-db", func(c *fiber.Ctx) error {
		results := make(map[string]string)

		// Test Main DB
		if sqlDB, err := database.DB.DB(); err == nil {
			if err := sqlDB.Ping(); err == nil {
				results["main_db"] = "✅ Connected"
			} else {
				results["main_db"] = "❌ Ping failed: " + err.Error()
			}
		} else {
			results["main_db"] = "❌ Connection failed: " + err.Error()
		}

		// Test NEOMAA
		if database.DBNeomaa != nil {
			if sqlDB, err := database.DBNeomaa.DB(); err == nil {
				if err := sqlDB.Ping(); err == nil {
					results["neomaa_db"] = "✅ Connected"
				} else {
					results["neomaa_db"] = "❌ Ping failed: " + err.Error()
				}
			} else {
				results["neomaa_db"] = "❌ Connection failed: " + err.Error()
			}
		} else {
			results["neomaa_db"] = "❌ Not initialized"
		}

		// Test NEOMAAREF
		if database.DBNeomaaRef != nil {
			if sqlDB, err := database.DBNeomaaRef.DB(); err == nil {
				if err := sqlDB.Ping(); err == nil {
					results["neomaaref_db"] = "✅ Connected"
				} else {
					results["neomaaref_db"] = "❌ Ping failed: " + err.Error()
				}
			} else {
				results["neomaaref_db"] = "❌ Connection failed: " + err.Error()
			}
		} else {
			results["neomaaref_db"] = "❌ Not initialized"
		}

		// Test SIMPEG
		if database.DBSimpeg != nil {
			if sqlDB, err := database.DBSimpeg.DB(); err == nil {
				if err := sqlDB.Ping(); err == nil {
					results["simpeg_db"] = "✅ Connected"
				} else {
					results["simpeg_db"] = "❌ Ping failed: " + err.Error()
				}
			} else {
				results["simpeg_db"] = "❌ Connection failed: " + err.Error()
			}
		} else {
			results["simpeg_db"] = "❌ Not initialized"
		}

		return c.JSON(fiber.Map{
			"message": "Database Connection Test",
			"results": results,
		})
	})

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

	// Reference Data routes (Fakultas & Prodi from NEOMAAREF)
	referenceController := controllers.NewReferenceController()
	reference := protected.Group("/reference")
	{
		reference.Get("/fakultas", referenceController.GetAllFakultas)
		reference.Get("/prodi", referenceController.GetAllProdi)
		reference.Get("/prodi/fakultas/:kode", referenceController.GetProdiByFakultas)
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
		menusPublic.Get("/", menuController.GetList)              // All users can read
		menusPublic.Get("/tree", menuController.GetTree)          // All users can read (all menus)
		menusPublic.Get("/my-tree", menuController.GetMyMenuTree) // Filtered by user level
		menusPublic.Get("/:id", menuController.GetByID)           // All users can read
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
		kategoriPublic.Get("/", kategoriPKMController.GetList)    // All users can read
		kategoriPublic.Get("/:id", kategoriPKMController.GetByID) // All users can read
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
		statusPublic.Get("/", statusReviewController.GetList)    // All users can read
		statusPublic.Get("/:id", statusReviewController.GetByID) // All users can read
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
		paramPublic.Get("/kategori/:id_kategori", parameterFormController.GetByKategori) // Important for mahasiswa
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

	// Tanggal Setting routes
	tglSettingController := controllers.NewTglSettingController()
	// Public endpoint - check if registration is open
	tglSettingPublic := protected.Group("/tgl-setting")
	{
		tglSettingPublic.Get("/active", tglSettingController.GetActive) // All authenticated users can check
	}
	// Admin only endpoints
	tglSettingAdmin := protected.Group("/tgl-setting", middleware.RequireAdmin())
	{
		tglSettingAdmin.Get("/", tglSettingController.GetList)
		tglSettingAdmin.Get("/:id", tglSettingController.GetByID)
		tglSettingAdmin.Post("/", tglSettingController.Create)
		tglSettingAdmin.Put("/:id", tglSettingController.Update)
		tglSettingAdmin.Delete("/:id", tglSettingController.Delete)
	}

	// pengajuan pkm - mahasiswa endpoints
	PengajuanController := controllers.NewPengajuanController()

	// Announcements (accessible to all authenticated users)
	protected.Get("/pengajuan/announcements", PengajuanController.GetAnnouncements)

	pengajuanMhs := protected.Group("/pengajuan", middleware.RequireMahasiswa())
	{
		// Judul PKM
		pengajuanMhs.Post("/judul", PengajuanController.CreateJudulPKM)
		pengajuanMhs.Put("/judul/:id", PengajuanController.UpdateJudul)

		// Proposal
		pengajuanMhs.Post("/:id/proposal", PengajuanController.UploadProposal)
		pengajuanMhs.Put("/:id/proposal", PengajuanController.ReviseProposal)

		// List & Detail
		pengajuanMhs.Get("/my-submissions", PengajuanController.GetMySubmissions)
		pengajuanMhs.Get("/:id", PengajuanController.GetPengajuanDetail)
	}

	// pengajuan pkm - admin endpoints
	pengajuanAdminController := controllers.NewPengajuanAdminController()
	pengajuanAdmin := protected.Group("/admin/pengajuan")
	{
		// List & Detail - Accessible by Admin and Reviewer
		pengajuanAdmin.Get("/", middleware.RequireAdminOrReviewer(), pengajuanAdminController.GetAllPengajuan)
		pengajuanAdmin.Get("/:id", middleware.RequireAdminOrReviewer(), pengajuanAdminController.GetPengajuanDetail)

		// Assign Reviewer - Strictly Admin only
		pengajuanAdmin.Post("/:id/assign-reviewer-judul", middleware.RequireAdmin(), pengajuanAdminController.AssignReviewerJudul)
		pengajuanAdmin.Post("/:id/assign-reviewer-proposal", middleware.RequireAdmin(), pengajuanAdminController.AssignReviewerProposal)

		// Cancel Plotting - Strictly Admin only
		pengajuanAdmin.Post("/:id/cancel-plotting-judul", middleware.RequireAdmin(), pengajuanAdminController.CancelPlottingJudul)
		pengajuanAdmin.Post("/:id/cancel-plotting-proposal", middleware.RequireAdmin(), pengajuanAdminController.CancelPlottingProposal)

		// Upload Proposal (admin can upload on behalf of mahasiswa) - Strictly Admin only
		pengajuanAdmin.Post("/:id/proposal", middleware.RequireAdmin(), pengajuanAdminController.UploadProposal)

		// Announce Final Result - Strictly Admin only
		pengajuanAdmin.Post("/:id/announce", middleware.RequireAdmin(), pengajuanAdminController.AnnounceFinalResult)
	}

	// reviewer management - admin endpoints
	reviewerController := controllers.NewReviewerController()
	reviewerAdmin := protected.Group("/admin/reviewers", middleware.RequireAdmin())
	{
		reviewerAdmin.Get("/", reviewerController.GetAllReviewers)
		reviewerAdmin.Get("/available", reviewerController.GetAvailablePegawai)
		reviewerAdmin.Post("/", reviewerController.ActivateReviewer)
		reviewerAdmin.Put("/:id", reviewerController.UpdateReviewer)
		reviewerAdmin.Delete("/:id", reviewerController.DeleteReviewer)
	}

	// pengajuan pkm - reviewer endpoints
	pengajuanReviewerController := controllers.NewPengajuanReviewerController()
	pengajuanReviewer := protected.Group("/reviewer", middleware.RequireReviewer())
	{
		// My Assignments
		pengajuanReviewer.Get("/my-assignments", pengajuanReviewerController.GetMyAssignments)

		// Detail
		pengajuanReviewer.Get("/pengajuan/:id", pengajuanReviewerController.GetPengajuanDetail)

		// Review
		pengajuanReviewer.Post("/judul/:id/review", pengajuanReviewerController.ReviewJudul)
		pengajuanReviewer.Post("/proposal/:id/review", pengajuanReviewerController.ReviewProposal)

		// Cancel Review
		pengajuanReviewer.Post("/judul/:id/cancel-review", pengajuanReviewerController.CancelReviewJudul)
		pengajuanReviewer.Post("/proposal/:id/cancel-review", pengajuanReviewerController.CancelReviewProposal)
	}

	//user akses management routes
	userAksesController := controllers.NewUserAksesController()
	userAksesAdmin := protected.Group("/admin/user-akses", middleware.RequireAdmin())
	{
		userAksesAdmin.Get("/", userAksesController.GetAllAccesses)
		userAksesAdmin.Get("/by-level/:id_user_level", userAksesController.GetAccessesByUserLevel)
		userAksesAdmin.Get("/:id", userAksesController.GetAccessDetail)
		userAksesAdmin.Post("/", userAksesController.CreateAccess)
		userAksesAdmin.Post("/bulk", userAksesController.BulkCreateAccess)
		userAksesAdmin.Put("/:id", userAksesController.UpdateAccess)
		userAksesAdmin.Delete("/:id", userAksesController.DeleteAccess)
	}
}

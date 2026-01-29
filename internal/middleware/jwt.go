package middleware

import (
	"strings"

	"rires-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// JWTAuth adalah middleware untuk validasi JWT token
func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		// ==================================
		// BYPASS PUBLIC ROUTES
		// ==================================
		if strings.HasPrefix(path, "/api/v1/auth/login") ||
			strings.HasPrefix(path, "/swagger") ||
			path == "/" ||
			path == "/health" {
			return c.Next()
		}

		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Missing authorization header",
			})
		}

		// Check Bearer format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid authorization header format. Use: Bearer <token>",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid or expired token",
			})
		}

		// Set user info to context
		c.Locals("id_user", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("email", claims.Email)
		c.Locals("user_type", claims.UserType)
		c.Locals("id_user_level", claims.IDUserLevel)
		c.Locals("user_data", claims.UserData)
		c.Locals("claims", claims)

		// Continue to next handler
		return c.Next()
	}
}

// RequireAdmin adalah middleware untuk memastikan user adalah admin
func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type")
		if userType != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Access denied. Admin only.",
			})
		}
		return c.Next()
	}
}

// RequireMahasiswa adalah middleware untuk memastikan user adalah mahasiswa atau admin
func RequireMahasiswa() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type")
		// Admin dapat mengakses semua endpoint
		if userType == "admin" {
			return c.Next()
		}
		if userType != "mahasiswa" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Access denied. Mahasiswa only.",
			})
		}
		return c.Next()
	}
}

// RequireReviewer adalah middleware untuk memastikan user adalah reviewer (pegawai) atau admin
func RequireReviewer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type")
		// Admin dapat mengakses semua endpoint
		if userType == "admin" {
			return c.Next()
		}
		if userType != "pegawai" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Access denied. Reviewer only.",
			})
		}
		return c.Next()
	}
}

// RequireAdminOrReviewer untuk route yang bisa diakses admin atau reviewer
func RequireAdminOrReviewer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type")
		if userType != "admin" && userType != "pegawai" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Access denied. Admin or Reviewer only.",
			})
		}
		return c.Next()
	}
}

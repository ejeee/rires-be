package utils

import (
	"github.com/gofiber/fiber/v2"
)

// GetCurrentUserID mengambil user ID dari context (setelah melalui JWT middleware)
func GetCurrentUserID(c *fiber.Ctx) uint {
	userID := c.Locals("id_user")
	if userID == nil {
		return 0
	}

	// Handle both uint and float64 (from JWT claims)
	switch v := userID.(type) {
	case uint:
		return v
	case float64:
		return uint(v)
	case int:
		return uint(v)
	default:
		return 0
	}
}

// GetCurrentUsername mengambil username dari context
func GetCurrentUsername(c *fiber.Ctx) string {
	username := c.Locals("username")
	if username == nil {
		return ""
	}
	return username.(string)
}

// GetCurrentUserType mengambil user type dari context (admin, mahasiswa, pegawai)
func GetCurrentUserType(c *fiber.Ctx) string {
	userType := c.Locals("user_type")
	if userType == nil {
		return ""
	}
	return userType.(string)
}

// GetCurrentUserEmail mengambil email dari context
func GetCurrentUserEmail(c *fiber.Ctx) string {
	email := c.Locals("email")
	if email == nil {
		return ""
	}
	return email.(string)
}

// GetCurrentUserData mengambil user data map dari context
func GetCurrentUserData(c *fiber.Ctx) map[string]string {
	userData := c.Locals("user_data")
	if userData == nil {
		return make(map[string]string)
	}
	return userData.(map[string]string)
}

// IsAdmin memeriksa apakah user adalah admin
func IsAdmin(c *fiber.Ctx) bool {
	return GetCurrentUserType(c) == "admin"
}

// IsMahasiswa memeriksa apakah user adalah mahasiswa
func IsMahasiswa(c *fiber.Ctx) bool {
	return GetCurrentUserType(c) == "mahasiswa"
}

// IsReviewer memeriksa apakah user adalah reviewer (pegawai)
func IsReviewer(c *fiber.Ctx) bool {
	return GetCurrentUserType(c) == "pegawai"
}

// GetCurrentUserLevel mengambil id_user_level dari context
// Returns: 1=superadmin, 2=admin, 3=mahasiswa, 4=reviewer, 0=not found
func GetCurrentUserLevel(c *fiber.Ctx) int {
	level := c.Locals("id_user_level")
	if level == nil {
		return 0
	}

	switch v := level.(type) {
	case int:
		return v
	case float64:
		return int(v)
	default:
		return 0
	}
}

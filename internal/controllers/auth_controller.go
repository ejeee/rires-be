package controllers

import (
	"crypto/sha1"
	"encoding/hex"
	"rires-be/internal/dto/request"
	"rires-be/internal/dto/response"
	"rires-be/internal/models"
	"rires-be/pkg/database"
	"rires-be/pkg/services"
	"rires-be/pkg/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	apiService *services.APIService
}

func NewAuthController() *AuthController {
	return &AuthController{
		apiService: services.NewAPIService(),
	}
}

// LoginAdmin godoc
// @Summary Login Admin
// @Description Login untuk administrator dari database lokal
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body object{username=string,password=string} true "Login Credentials"
// @Success 200 {object} object{success=bool,message=string,data=object{token=string,user_type=string,expires_in=int,user=object}}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 401 {object} object{success=bool,message=string}
// @Router /auth/login/admin [post]
func (ctrl *AuthController) LoginAdmin(c *fiber.Ctx) error {
	var req request.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		return utils.BadRequestResponse(c, "Username and password are required")
	}

	// Find user in database
	var user models.User
	result := database.DB.Where("username = ? AND hapus = 0", req.Username).First(&user)
	if result.Error != nil {
		return utils.UnauthorizedResponse(c, "Invalid username or password")
	}

	// Check if user is active
	if user.Status != 1 {
		return utils.UnauthorizedResponse(c, "User account is inactive")
	}

	// Verify password
	// Check if password is hashed with MySQL PASSWORD() function
	if len(user.Password) > 0 && user.Password[0] == '*' {
		// MySQL PASSWORD() hash format: *HEXSTRING
		hashedInput := hashMySQLPassword(req.Password)
		if user.Password != hashedInput {
			return utils.UnauthorizedResponse(c, "Invalid username or password")
		}
	} else if len(user.Password) >= 4 && (user.Password[0:4] == "$2a$" || user.Password[0:4] == "$2y$") {
		// Bcrypt hash
		if err := utils.VerifyPassword(user.Password, req.Password); err != nil {
			return utils.UnauthorizedResponse(c, "Invalid username or password")
		}
	} else {
		// Plain text password
		if user.Password != req.Password {
			return utils.UnauthorizedResponse(c, "Invalid username or password")
		}
	}

	// Generate JWT token with user data
	token, err := utils.GenerateTokenWithClaims(
		uint(user.ID),
		user.Username,
		"",
		"admin",
		map[string]string{
			"nama_user":  user.NamaUser,
			"level_user": strconv.Itoa(user.LevelUser),
		},
	)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to generate token")
	}

	// Parse JWT expired hours
	expiredHours, _ := strconv.Atoi("24")

	// Prepare response
	loginResponse := response.LoginResponse{
		Token:     token,
		UserType:  "admin",
		ExpiresIn: expiredHours,
		User: response.AdminLoginResponse{
			ID:        user.ID,
			NamaUser:  user.NamaUser,
			Username:  user.Username,
			LevelUser: user.LevelUser,
			Status:    user.Status,
		},
	}

	return utils.SuccessResponse(c, "Login successful", loginResponse)
}

// LoginMahasiswa godoc
// @Summary Login Mahasiswa
// @Description Login mahasiswa menggunakan NIM via API external
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body object{username=string,password=string} true "NIM and Password"
// @Success 200 {object} object{success=bool,message=string,data=object{token=string,user_type=string,expires_in=int,user=object}}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 401 {object} object{success=bool,message=string}
// @Router /auth/login/mahasiswa [post]
func (ctrl *AuthController) LoginMahasiswa(c *fiber.Ctx) error {
	var req request.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		return utils.BadRequestResponse(c, "Username and password are required")
	}

	// Call external API
	mahasiswa, err := ctrl.apiService.MahasiswaLogin(req.Username, req.Password)
	if err != nil {
		return utils.UnauthorizedResponse(c, err.Error())
	}

	// Generate JWT token with mahasiswa data
	token, err := utils.GenerateTokenWithClaims(
		0, // ID = 0 karena bukan dari DB lokal
		mahasiswa.NIM,
		"",
		"mahasiswa",
		map[string]string{
			"nama":     mahasiswa.Nama,
			"prodi":    mahasiswa.Prodi,
			"fakultas": mahasiswa.Fakultas,
		},
	)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to generate token")
	}

	// Parse JWT expired hours
	expiredHours, _ := strconv.Atoi("24")

	// Prepare response
	loginResponse := response.LoginResponse{
		Token:     token,
		UserType:  "mahasiswa",
		ExpiresIn: expiredHours,
		User:      mahasiswa,
	}

	return utils.SuccessResponse(c, "Login successful", loginResponse)
}

// LoginPegawai godoc
// @Summary Login Pegawai
// @Description Login pegawai/reviewer menggunakan email via API external
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body object{username=string,password=string} true "Email and Password"
// @Success 200 {object} object{success=bool,message=string,data=object{token=string,user_type=string,expires_in=int,user=object}}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 401 {object} object{success=bool,message=string}
// @Router /auth/login/pegawai [post]
func (ctrl *AuthController) LoginPegawai(c *fiber.Ctx) error {
	var req request.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		return utils.BadRequestResponse(c, "Username and password are required")
	}

	// Call external API
	pegawai, err := ctrl.apiService.PegawaiLogin(req.Username, req.Password)
	if err != nil {
		return utils.UnauthorizedResponse(c, err.Error())
	}

	// Generate JWT token with pegawai data
	token, err := utils.GenerateTokenWithClaims(
		0, // ID = 0 karena bukan dari DB lokal
		pegawai.NIP,
		pegawai.Email,
		"pegawai",
		map[string]string{
			"nama":    pegawai.Nama,
			"jabatan": pegawai.Jabatan,
			"unit":    pegawai.Unit,
		},
	)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to generate token")
	}

	// Parse JWT expired hours
	expiredHours, _ := strconv.Atoi("24")

	// Prepare response
	loginResponse := response.LoginResponse{
		Token:     token,
		UserType:  "pegawai",
		ExpiresIn: expiredHours,
		User:      pegawai,
	}

	return utils.SuccessResponse(c, "Login successful", loginResponse)
}

// hashMySQLPassword creates MySQL PASSWORD() compatible hash
// MySQL PASSWORD() uses double SHA1: *UPPERCASE_HEX(SHA1(SHA1(password)))
func hashMySQLPassword(password string) string {
	// First SHA1
	firstHash := sha1.Sum([]byte(password))
	
	// Second SHA1
	secondHash := sha1.Sum(firstHash[:])
	
	// Convert to uppercase hex with * prefix
	result := "*" + strings.ToUpper(hex.EncodeToString(secondHash[:]))
	
	return result
}

// GetCurrentUser godoc
// @Summary Get Current User
// @Description Get currently logged in user information from JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} object{success=bool,message=string,data=object}
// @Failure 401 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /auth/me [get]
func (ctrl *AuthController) GetCurrentUser(c *fiber.Ctx) error {
	// Get Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return utils.UnauthorizedResponse(c, "Missing authorization header")
	}

	// Check if it starts with "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return utils.UnauthorizedResponse(c, "Invalid authorization header format")
	}

	// Extract token
	tokenString := authHeader[7:]

	// Parse & validate token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid or expired token")
	}

	// Get user type from claims struct (NOT map!)
	userType := claims.UserType
	userID := int(claims.UserID)

	// Fetch user data based on type
	var userData interface{}

	switch userType {
	case "admin":
		// Get admin from database
		var user models.User
		if err := database.DB.Where("id = ? AND hapus = ?", userID, 0).First(&user).Error; err != nil {
			return utils.NotFoundResponse(c, "User not found")
		}

		userData = fiber.Map{
			"id":         user.ID,
			"username":   user.Username,
			"nama_user":  user.NamaUser,
			"level_user": user.LevelUser,
			"status":     user.Status,
		}

	case "mahasiswa":
		// Get from token claims (data dari API external)
		userData = fiber.Map{
			"nim":      claims.Username,
			"nama":     claims.UserData["nama"],
			"prodi":    claims.UserData["prodi"],
			"fakultas": claims.UserData["fakultas"],
		}

	case "pegawai":
		// Get from token claims (data dari API external)
		userData = fiber.Map{
			"nip":     claims.Username,
			"nama":    claims.UserData["nama"],
			"jabatan": claims.UserData["jabatan"],
			"unit":    claims.UserData["unit"],
			"email":   claims.Email,
		}

	default:
		return utils.BadRequestResponse(c, "Unknown user type")
	}

	response := fiber.Map{
		"user_type": userType,
		"user":      userData,
	}

	return utils.SuccessResponse(c, "User data retrieved successfully", response)
}
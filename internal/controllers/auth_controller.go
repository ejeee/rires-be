package controllers

import (
	"crypto/sha1"
	"encoding/hex"
	request "rires-be/internal/dto/request"
	response "rires-be/internal/dto/response"
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

	// Generate JWT token
	token, err := utils.GenerateToken(uint(user.ID), user.Username)
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

	// Generate JWT token
	// Gunakan NIM sebagai unique identifier
	token, err := utils.GenerateToken(0, mahasiswa.NIM) // ID = 0 karena bukan dari DB lokal
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

	// Generate JWT token
	// Gunakan NIP sebagai unique identifier
	token, err := utils.GenerateToken(0, pegawai.NIP) // ID = 0 karena bukan dari DB lokal
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
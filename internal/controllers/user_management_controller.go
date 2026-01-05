package controllers

import (
	"math"
	"rires-be/internal/dto/request"
	"rires-be/internal/dto/response"
	"rires-be/internal/models"
	"rires-be/pkg/database"
	"rires-be/pkg/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserManagementController struct{}

func NewUserManagementController() *UserManagementController {
	return &UserManagementController{}
}

// GetList godoc
// @Summary List Users
// @Description Get list of users with pagination
// @Tags User Management
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search by nama_user or username"
// @Param level_user query int false "Filter by level_user"
// @Success 200 {object} object{success=bool,message=string,data=response.UserListResponse}
// @Failure 500 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /users [get]
func (ctrl *UserManagementController) GetList(c *fiber.Ctx) error {
	// Pagination params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	search := c.Query("search", "")
	levelUser := c.Query("level_user", "")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// Query builder
	query := database.DB.Model(&models.User{}).Where("hapus = ?", 0)

	// Filter by level
	if levelUser != "" {
		query = query.Where("level_user = ?", levelUser)
	}

	// Search
	if search != "" {
		query = query.Where("nama_user LIKE ? OR username LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to count data")
	}

	// Get data with level
	var users []models.User
	if err := query.Preload("Level").Order("id DESC").Offset(offset).Limit(perPage).Find(&users).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch data")
	}

	// Transform to response
	var data []response.UserResponse
	for _, user := range users {
		statusText := "Aktif"
		if user.Status == 2 {
			statusText = "Tidak Aktif"
		}

		namaLevel := ""
		if user.Level != nil {
			namaLevel = user.Level.NamaLevel
		}

		data = append(data, response.UserResponse{
			ID:         user.ID,
			NamaUser:   user.NamaUser,
			Username:   user.Username,
			LevelUser:  user.LevelUser,
			NamaLevel:  namaLevel,
			Status:     user.Status,
			StatusText: statusText,
			TglInsert:  user.TglInsert,
			TglUpdate:  user.TglUpdate,
			UserUpdate: user.UserUpdate,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	result := response.UserListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}

	return utils.SuccessResponse(c, "Data retrieved successfully", result)
}

// GetByID godoc
// @Summary Get User by ID
// @Description Get user detail by ID
// @Tags User Management
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} object{success=bool,message=string,data=response.UserResponse}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /users/{id} [get]
func (ctrl *UserManagementController) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var user models.User
	if err := database.DB.Preload("Level").Where("id = ? AND hapus = ?", id, 0).First(&user).Error; err != nil {
		return utils.NotFoundResponse(c, "User not found")
	}

	statusText := "Aktif"
	if user.Status == 2 {
		statusText = "Tidak Aktif"
	}

	namaLevel := ""
	if user.Level != nil {
		namaLevel = user.Level.NamaLevel
	}

	result := response.UserResponse{
		ID:         user.ID,
		NamaUser:   user.NamaUser,
		Username:   user.Username,
		LevelUser:  user.LevelUser,
		NamaLevel:  namaLevel,
		Status:     user.Status,
		StatusText: statusText,
		TglInsert:  user.TglInsert,
		TglUpdate:  user.TglUpdate,
		UserUpdate: user.UserUpdate,
	}

	return utils.SuccessResponse(c, "Data retrieved successfully", result)
}

// Create godoc
// @Summary Create User
// @Description Create new user with hashed password
// @Tags User Management
// @Accept json
// @Produce json
// @Param body body request.CreateUserRequest true "User Data"
// @Success 201 {object} object{success=bool,message=string,data=response.UserResponse}
// @Failure 400 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /users [post]
func (ctrl *UserManagementController) Create(c *fiber.Ctx) error {
	var req request.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate
	if req.NamaUser == "" || req.Username == "" || req.Password == "" {
		return utils.BadRequestResponse(c, "Required fields are missing")
	}
	if req.Status != 1 && req.Status != 2 {
		return utils.BadRequestResponse(c, "Status must be 1 (active) or 2 (inactive)")
	}

	// Check if level exists
	var level models.UserLevel
	if err := database.DB.Where("id = ? AND hapus = ?", req.LevelUser, 0).First(&level).Error; err != nil {
		return utils.BadRequestResponse(c, "User level not found")
	}

	// Check duplicate username
	var count int64
	database.DB.Model(&models.User{}).Where("username = ? AND hapus = ?", req.Username, 0).Count(&count)
	if count > 0 {
		return utils.BadRequestResponse(c, "Username already exists")
	}

	// Hash password using MySQL PASSWORD() function
	hashedPassword := hashMySQLPassword(req.Password)

	// Create
	now := time.Now()
	user := models.User{
		NamaUser:   req.NamaUser,
		Username:   req.Username,
		Password:   hashedPassword,
		LevelUser:  req.LevelUser,
		Status:     req.Status,
		Hapus:      0,
		TglInsert:  &now,
		UserUpdate: "1", // TODO: Get from JWT token
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create user")
	}

	statusText := "Aktif"
	if user.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.UserResponse{
		ID:         user.ID,
		NamaUser:   user.NamaUser,
		Username:   user.Username,
		LevelUser:  user.LevelUser,
		NamaLevel:  level.NamaLevel,
		Status:     user.Status,
		StatusText: statusText,
		TglInsert:  user.TglInsert,
		TglUpdate:  user.TglUpdate,
		UserUpdate: user.UserUpdate,
	}

	return utils.CreatedResponse(c, "User created successfully", result)
}

// Update godoc
// @Summary Update User
// @Description Update existing user (without password)
// @Tags User Management
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body request.UpdateUserRequest true "User Data"
// @Success 200 {object} object{success=bool,message=string,data=response.UserResponse}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /users/{id} [put]
func (ctrl *UserManagementController) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var req request.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate
	if req.NamaUser == "" || req.Username == "" {
		return utils.BadRequestResponse(c, "Required fields are missing")
	}
	if req.Status != 1 && req.Status != 2 {
		return utils.BadRequestResponse(c, "Status must be 1 (active) or 2 (inactive)")
	}

	// Find existing
	var user models.User
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&user).Error; err != nil {
		return utils.NotFoundResponse(c, "User not found")
	}

	// Check if level exists
	var level models.UserLevel
	if err := database.DB.Where("id = ? AND hapus = ?", req.LevelUser, 0).First(&level).Error; err != nil {
		return utils.BadRequestResponse(c, "User level not found")
	}

	// Check duplicate username (exclude current)
	var count int64
	database.DB.Model(&models.User{}).
		Where("username = ? AND id != ? AND hapus = ?", req.Username, id, 0).
		Count(&count)
	if count > 0 {
		return utils.BadRequestResponse(c, "Username already exists")
	}

	// Update (password tidak diubah di sini)
	user.NamaUser = req.NamaUser
	user.Username = req.Username
	user.LevelUser = req.LevelUser
	user.Status = req.Status
	user.UserUpdate = "1" // TODO: Get from JWT token

	if err := database.DB.Save(&user).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to update user")
	}

	statusText := "Aktif"
	if user.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.UserResponse{
		ID:         user.ID,
		NamaUser:   user.NamaUser,
		Username:   user.Username,
		LevelUser:  user.LevelUser,
		NamaLevel:  level.NamaLevel,
		Status:     user.Status,
		StatusText: statusText,
		TglInsert:  user.TglInsert,
		TglUpdate:  user.TglUpdate,
		UserUpdate: user.UserUpdate,
	}

	return utils.SuccessResponse(c, "User updated successfully", result)
}

// ResetPassword godoc
// @Summary Reset User Password
// @Description Reset user password (admin only)
// @Tags User Management
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body request.ResetPasswordRequest true "New Password"
// @Success 200 {object} object{success=bool,message=string}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /users/{id}/reset-password [post]
func (ctrl *UserManagementController) ResetPassword(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var req request.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if req.NewPassword == "" || len(req.NewPassword) < 6 {
		return utils.BadRequestResponse(c, "Password must be at least 6 characters")
	}

	// Find existing
	var user models.User
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&user).Error; err != nil {
		return utils.NotFoundResponse(c, "User not found")
	}

	// Hash new password
	hashedPassword := hashMySQLPassword(req.NewPassword)

	// Update password
	user.Password = hashedPassword
	user.UserUpdate = "1" // TODO: Get from JWT token

	if err := database.DB.Save(&user).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to reset password")
	}

	return utils.SuccessResponse(c, "Password reset successfully", nil)
}

// Delete godoc
// @Summary Delete User
// @Description Soft delete user
// @Tags User Management
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} object{success=bool,message=string}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /users/{id} [delete]
func (ctrl *UserManagementController) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	// Find existing
	var user models.User
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&user).Error; err != nil {
		return utils.NotFoundResponse(c, "User not found")
	}

	// Prevent deleting yourself (TODO: get current user from JWT)
	// if currentUserID == id {
	//     return utils.BadRequestResponse(c, "Cannot delete your own account")
	// }

	// Soft delete
	user.Hapus = 1
	user.UserUpdate = "1" // TODO: Get from JWT token

	if err := database.DB.Save(&user).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to delete user")
	}

	return utils.SuccessResponse(c, "User deleted successfully", nil)
}
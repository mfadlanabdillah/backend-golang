package controllers

import (
	"fadlan/backend-api/database"
	"fadlan/backend-api/helpers"
	"fadlan/backend-api/models"
	"fadlan/backend-api/structs"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func FindUsers(c *gin.Context) {
	// Default values
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// Inisialisasi slice untuk menampung data user
	var users []models.User
	var total int64

	// Hitung total data
	if err := database.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to count users",
		})
		return
	}

	// Ambil data user dari database dengan pagination
	if err := database.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch users",
		})
		return
	}

	// Map users to response DTO
	var userResponses []structs.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, structs.UserResponse{
			Id:        user.Id,
			BaseUser: structs.BaseUser{
				Name:     user.Name,
				Username: user.Username,
				Email:    user.Email,
			},
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	// Kirimkan response sukses dengan data user (tanpa password)
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List Data Users",
		Data: gin.H{
			"data":  userResponses,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func CreateUser(c *gin.Context) {
	// Inisialisasi struct request
	var req structs.UserCreateRequest

	// Validasi input JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Validation error",
			Error:   helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Validasi kekuatan password
	if err := helpers.ValidatePasswordStrength(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Normalisasi input
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	req.Name = strings.TrimSpace(req.Name)

	// Buat user baru
	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: helpers.HashPassword(req.Password),
	}

	// Simpan ke database
	if err := database.DB.Create(&user).Error; err != nil {
		// Cek error duplikat
		if helpers.IsDuplicateEntryError(err) {
			c.JSON(http.StatusConflict, structs.ErrorResponse{
				Success: false,
				Message: "Username or email already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create user",
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "User created successfully",
		Data: structs.UserResponse{
			Id:        user.Id,
			BaseUser: structs.BaseUser{
				Name:     user.Name,
				Username: user.Username,
				Email:    user.Email,
			},
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})

}
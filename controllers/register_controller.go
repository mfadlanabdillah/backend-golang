package controllers

import (
	"fadlan/backend-api/database"
	"fadlan/backend-api/helpers"
	"fadlan/backend-api/models"
	"fadlan/backend-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register menangani proses registrasi user baru
func Register(c *gin.Context) {
	// Inisialisasi struct untuk menangkap data request
	var req = structs.UserCreateRequest{}

	// Validasi request JSON menggunakan binding dari Gin
	if err := c.ShouldBindJSON(&req); err != nil {
		// Jika validasi gagal, kirimkan response error
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validasi Error",
			Error:   helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Buat data user baru dengan password yang sudah di-hash
	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: helpers.HashPassword(req.Password),
	}

	// Simpan data user ke database
	if err := database.DB.Create(&user).Error; err != nil {
		// Cek apakah error karena data duplikat (misalnya username/email sudah terdaftar)
		if helpers.IsDuplicateEntryError(err) {
			// Jika duplikat, kirimkan response 409 Conflict
			c.JSON(http.StatusConflict, structs.ErrorResponse{
				Success: false,
				Message: "Duplicate entry error",
				Error:   helpers.TranslateErrorMessage(err),
			})
		} else {
			// Jika error lain, kirimkan response 500 Internal Server Error
			c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
				Success: false,
				Message: "Failed to create user",
				Error:   helpers.TranslateErrorMessage(err),
			})
		}
		return
	}

	// Jika berhasil, kirimkan response sukses
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
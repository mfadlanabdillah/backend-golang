package controllers

import (
	"errors"
	"fadlan/backend-api/database"
	"fadlan/backend-api/helpers"
	"fadlan/backend-api/models"
	"fadlan/backend-api/structs"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// validatePasswordStrength memeriksa kekuatan password
func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`).MatchString(password)

	switch {
	case !hasUpper:
		return errors.New("password must contain at least one uppercase letter")
	case !hasLower:
		return errors.New("password must contain at least one lowercase letter")
	case !hasNumber:
		return errors.New("password must contain at least one number")
	case !hasSpecial:
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// Register menangani proses registrasi user baru
func Register(c *gin.Context) {
	// Inisialisasi struct untuk menangkap data request
	var req = structs.UserCreateRequest{}

	// Validasi request JSON menggunakan binding dari Gin
	if err := c.ShouldBindJSON(&req); err != nil {
		// Jika validasi gagal, kirimkan response error
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation error",
			Error:   helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Validasi tambahan untuk kekuatan password
	if err := validatePasswordStrength(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Normalisasi email (lowercase)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))

	// Buat data user baru dengan password yang sudah di-hash
	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: helpers.HashPassword(req.Password),
	}

	// Simpan data user ke database
	if err := database.DB.Create(&user).Error; err != nil {
			// Cek apakah error karena data duplikat
		if helpers.IsDuplicateEntryError(err) {
			// Gunakan pesan error yang umum untuk menghindari user enumeration
			c.JSON(http.StatusConflict, structs.ErrorResponse{
				Success: false,
				Message: "Could not create user. The username or email may already be registered.",
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
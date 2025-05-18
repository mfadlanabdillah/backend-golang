package controllers

import (
	"fadlan/backend-api/database"
	"fadlan/backend-api/helpers"
	"fadlan/backend-api/models"
	"fadlan/backend-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {

	// Inisialisasi struct untuk menampung data dari request
	var req = structs.UserLoginRequest{}
	var user = models.User{}

	// Validasi input dari request body menggunakan ShouldBindJSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Error:   helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Cari user berdasarkan username yang diberikan di database
	// Jika tidak ditemukan, kirimkan respons error Unauthorized
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "User Not Found",
			Error:   helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Bandingkan password yang dimasukkan dengan password yang sudah di-hash di database
	// Jika tidak cocok, kirimkan respons error Unauthorized
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "Invalid Password",
			Error:   helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Jika login berhasil, generate token untuk user
	token, err := helpers.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to generate token",
			Error:   helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirimkan response sukses dengan status OK dan data user serta token
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Login Success",
		Data: structs.UserResponse{
			Id:        user.Id,
			BaseUser: structs.BaseUser{
				Name:     user.Name,
				Username: user.Username,
				Email:    user.Email,
			},
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Token:     &token,
		},
	})
}
package handlers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"candyzflorist-go/config"
	"candyzflorist-go/helpers"
	"candyzflorist-go/middleware"
	"candyzflorist-go/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Email dan password wajib diisi")
		return
	}

	var user models.User
	// Cari berdasarkan email atau username
	if err := config.DB.Where("email = ? OR username = ?", input.Email, input.Email).First(&user).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusUnauthorized, "Email atau password salah")
		return
	}

	// Verifikasi password dengan bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		helpers.ErrorResponse(c, http.StatusUnauthorized, "Email atau password salah")
		return
	}

	// Generate JWT Token
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "candyzflorist-super-secret-key-2026"
	}

	expiryHoursStr := os.Getenv("JWT_EXPIRY_HOURS")
	expiryHours, err := strconv.Atoi(expiryHoursStr)
	if err != nil || expiryHours <= 0 {
		expiryHours = 72
	}

	expirationTime := time.Now().Add(time.Duration(expiryHours) * time.Hour)
	claims := &middleware.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "candyzflorist-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat token autentikasi")
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Login berhasil", gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		helpers.ErrorResponse(c, http.StatusUnauthorized, "Sesi tidak ditemukan")
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusNotFound, "Pengguna tidak ditemukan")
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Profil pengguna", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

func ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		helpers.ErrorResponse(c, http.StatusUnauthorized, "Sesi tidak ditemukan")
		return
	}

	var input ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Data tidak valid: password baru minimal 8 karakter")
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusNotFound, "Pengguna tidak ditemukan")
		return
	}

	// Cek password lama
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Password lama tidak sesuai")
		return
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengenkripsi password baru")
		return
	}

	user.Password = string(hashedPassword)
	config.DB.Save(&user)

	helpers.SuccessResponse(c, http.StatusOK, "Password berhasil diperbarui", nil)
}

func Logout(c *gin.Context) {
	helpers.SuccessResponse(c, http.StatusOK, "Logout berhasil", nil)
}

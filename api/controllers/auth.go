package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"sipecut/config"
	"sipecut/middleware"
	"sipecut/models"
)

type LoginInput struct {
	NIP      string `json:"nip" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type ResetPasswordInput struct {
	NIP         string `json:"nip" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// Login handles user authentication
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()
	var user models.User
	if err := db.Where("nip = ?", input.NIP).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "NIP atau Password salah"})
		return
	}

	if user.Status != "aktif" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun tidak aktif"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "NIP atau Password salah"})
		return
	}

	// Generate JWT
	token, refreshToken, err := middleware.GenerateToken(user.NIP, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// Get Employee details if exists
	var employee models.Employee
	hasProfile := false
	if user.Role == "employee" {
		if err := db.Where("nip = ?", user.NIP).First(&employee).Error; err == nil {
			hasProfile = true
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"token":         token,
		"refresh_token": refreshToken,
		"role":          user.Role,
		"nip":           user.NIP,
		"no_hp":         user.NoHP,
		"has_profile":   hasProfile,
		"employee":      employee,
	})
}

// RefreshToken handles JWT token refresh requests
func RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := middleware.ValidateToken(body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token penyegar tidak valid"})
		return
	}

	nip := claims["nip"].(string)

	db := config.GetDB()
	var user models.User
	if err := db.Where("nip = ?", nip).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	token, newRefresh, err := middleware.GenerateToken(user.NIP, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":         token,
		"refresh_token": newRefresh,
	})
}

// ChangePassword allows users to change their own password
func ChangePassword(c *gin.Context) {
	var input ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nip, _ := c.Get("nip")
	db := config.GetDB()

	var user models.User
	if err := db.Where("nip = ?", nip).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password lama salah"})
		return
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password baru"})
		return
	}

	user.PasswordHash = string(newHash)
	db.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}

// ResetPassword allows admin to reset any user's password
func ResetPassword(c *gin.Context) {
	var input ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()
	var user models.User
	if err := db.Where("nip = ?", input.NIP).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pegawai dengan NIP tersebut tidak ditemukan"})
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mereset password"})
		return
	}

	user.PasswordHash = string(newHash)
	db.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Password untuk NIP " + input.NIP + " berhasil direset"})
}

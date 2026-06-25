package controllers

import (
	"net/http"
<<<<<<< HEAD
	"time"
=======
>>>>>>> 603353f54c6625439da1b7cf09eb935c784c51b4

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

<<<<<<< HEAD
type CreatePegawaiInput struct {
	NIP    string `json:"nip" binding:"required"`
	Nama   string `json:"nama" binding:"required"`
	RoleID uint   `json:"role_id"`
	NoHP   string `json:"no_hp"`
}

type UpdateAksesInput struct {
	RoleID  uint   `json:"role_id" binding:"required"`
	MenuIDs []uint `json:"menu_ids"`
}

// Login handles user authentication with NIP
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIP dan password wajib diisi"})
=======
// Login handles user authentication
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
>>>>>>> 603353f54c6625439da1b7cf09eb935c784c51b4
		return
	}

	db := config.GetDB()
	var user models.User
<<<<<<< HEAD
	if err := db.Preload("Menus").Preload("RoleRel").Where("nip = ?", input.NIP).First(&user).Error; err != nil {
		// Log failed login
		db.Create(&models.RequestHistory{
			RequestType: "login",
			RequestID:   0,
			StatusBaru:  "Failed",
			Catatan:     "Login failed: NIP not found. IP: " + c.ClientIP(),
			ChangedBy:   input.NIP,
			ChangedAt:   time.Now(),
		})
=======
	if err := db.Where("nip = ?", input.NIP).First(&user).Error; err != nil {
>>>>>>> 603353f54c6625439da1b7cf09eb935c784c51b4
		c.JSON(http.StatusUnauthorized, gin.H{"error": "NIP atau Password salah"})
		return
	}

	if user.Status != "aktif" {
<<<<<<< HEAD
		db.Create(&models.RequestHistory{
			RequestType: "login",
			RequestID:   user.ID,
			StatusBaru:  "Failed",
			Catatan:     "Login failed: Account inactive. IP: " + c.ClientIP(),
			ChangedBy:   input.NIP,
			ChangedAt:   time.Now(),
		})
=======
>>>>>>> 603353f54c6625439da1b7cf09eb935c784c51b4
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun tidak aktif"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
<<<<<<< HEAD
		db.Create(&models.RequestHistory{
			RequestType: "login",
			RequestID:   user.ID,
			StatusBaru:  "Failed",
			Catatan:     "Login failed: Password incorrect. IP: " + c.ClientIP(),
			ChangedBy:   input.NIP,
			ChangedAt:   time.Now(),
		})
=======
>>>>>>> 603353f54c6625439da1b7cf09eb935c784c51b4
		c.JSON(http.StatusUnauthorized, gin.H{"error": "NIP atau Password salah"})
		return
	}

<<<<<<< HEAD
	// Get menu codes
	var menuCodes []string
	for _, m := range user.Menus {
		menuCodes = append(menuCodes, m.KodeMenu)
	}

	// Generate JWT
	token, refreshToken, err := middleware.GenerateToken(user.NIP, user.Role, menuCodes)
=======
	// Generate JWT
	token, refreshToken, err := middleware.GenerateToken(user.NIP, user.Role)
>>>>>>> 603353f54c6625439da1b7cf09eb935c784c51b4
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

<<<<<<< HEAD
	// Log successful login
	db.Create(&models.RequestHistory{
		RequestType: "login",
		RequestID:   user.ID,
		StatusBaru:  "Success",
		Catatan:     "Login successful. IP: " + c.ClientIP(),
		ChangedBy:   user.NIP,
		ChangedAt:   time.Now(),
	})

	// Get Employee details if exists
	var employee models.Employee
	hasProfile := false
	if err := db.Where("nip = ?", user.NIP).First(&employee).Error; err == nil {
		hasProfile = true
	}

	c.JSON(http.StatusOK, gin.H{
		"token":               token,
		"refresh_token":       refreshToken,
		"role":                user.Role,
		"role_id":             user.RoleID,
		"nip":                 user.NIP,
		"no_hp":               user.NoHP,
		"is_password_default": user.IsPasswordDefault,
		"has_profile":         hasProfile,
		"employee":            employee,
		"menus":               user.Menus,
=======
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
>>>>>>> 603353f54c6625439da1b7cf09eb935c784c51b4
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
<<<<<<< HEAD
	if err := db.Preload("Menus").Preload("RoleRel").Where("nip = ?", nip).First(&user).Error; err != nil {
=======
	if err := db.Where("nip = ?", nip).First(&user).Error; err != nil {
>>>>>>> 603353f54c6625439da1b7cf09eb935c784c51b4
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

<<<<<<< HEAD
	if user.Status != "aktif" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun tidak aktif"})
		return
	}

	var menuCodes []string
	for _, m := range user.Menus {
		menuCodes = append(menuCodes, m.KodeMenu)
	}

	token, newRefresh, err := middleware.GenerateToken(user.NIP, user.Role, menuCodes)
=======
	token, newRefresh, err := middleware.GenerateToken(user.NIP, user.Role)
>>>>>>> 603353f54c6625439da1b7cf09eb935c784c51b4
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

<<<<<<< HEAD
	// Validate new password must not be equal to NIP for security
	if input.NewPassword == user.NIP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password baru tidak boleh sama dengan NIP"})
		return
	}

=======
>>>>>>> 603353f54c6625439da1b7cf09eb935c784c51b4
	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password baru"})
		return
	}

	user.PasswordHash = string(newHash)
<<<<<<< HEAD
	user.IsPasswordDefault = false
	db.Save(&user)

	// Log audit trail
	db.Create(&models.RequestHistory{
		RequestType: "change_password",
		RequestID:   user.ID,
		StatusBaru:  "Success",
		Catatan:     "Password changed by user. IP: " + c.ClientIP(),
		ChangedBy:   user.NIP,
		ChangedAt:   time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}

// ResetPassword allows admin to reset user's password to custom password (kept for backward compatibility)
=======
	db.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}

// ResetPassword allows admin to reset any user's password
>>>>>>> 603353f54c6625439da1b7cf09eb935c784c51b4
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
<<<<<<< HEAD
	user.IsPasswordDefault = false
	db.Save(&user)

	// Log audit trail
	adminNip, _ := c.Get("nip")
	db.Create(&models.RequestHistory{
		RequestType: "admin_reset_password_custom",
		RequestID:   user.ID,
		StatusBaru:  "Success",
		Catatan:     "Password for NIP " + input.NIP + " reset to custom password by Admin",
		ChangedBy:   adminNip.(string),
		ChangedAt:   time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Password untuk NIP " + input.NIP + " berhasil direset"})
}

// GetMe returns current logged-in employee details and menus
func GetMe(c *gin.Context) {
	nip, _ := c.Get("nip")
	db := config.GetDB()

	var user models.User
	if err := db.Preload("Menus").Preload("RoleRel").Where("nip = ?", nip).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pegawai tidak ditemukan"})
		return
	}

	var employee models.Employee
	hasProfile := false
	if err := db.Where("nip = ?", user.NIP).First(&employee).Error; err == nil {
		hasProfile = true
	}

	c.JSON(http.StatusOK, gin.H{
		"id":                  user.ID,
		"nip":                 user.NIP,
		"no_hp":               user.NoHP,
		"status":              user.Status,
		"role_id":             user.RoleID,
		"role":                user.Role,
		"role_rel":            user.RoleRel,
		"is_password_default": user.IsPasswordDefault,
		"menus":               user.Menus,
		"has_profile":         hasProfile,
		"employee":            employee,
	})
}

// AdminGetPegawai lists all employee users
func AdminGetPegawai(c *gin.Context) {
	db := config.GetDB()
	var users []models.User
	if err := db.Preload("RoleRel").Preload("Menus").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pegawai"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// AdminCreatePegawai creates a new employee with password default = NIP
func AdminCreatePegawai(c *gin.Context) {
	var input CreatePegawaiInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate NIP format (only numeric digits, 8-20 length)
	isNumeric := true
	for _, char := range input.NIP {
		if char < '0' || char > '9' {
			isNumeric = false
			break
		}
	}
	if !isNumeric || len(input.NIP) < 8 || len(input.NIP) > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format NIP tidak valid. Harus berupa angka 8-20 digit."})
		return
	}

	db := config.GetDB()

	var existingUser models.User
	if err := db.Where("nip = ?", input.NIP).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIP sudah terdaftar sebagai akun pengguna"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.NIP), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-hash password default"})
		return
	}

	roleID := input.RoleID
	if roleID == 0 {
		roleID = 2 // default employee
	}

	var role models.Role
	if err := db.First(&role, roleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak valid"})
		return
	}

	var menus []models.Menu
	if role.NamaRole == "admin" {
		db.Find(&menus)
	} else {
		db.Where("kode_menu IN ?", []string{"self_profile", "leave_request", "data_change", "berita_acara"}).Find(&menus)
	}

	user := models.User{
		NIP:               input.NIP,
		PasswordHash:      string(hash),
		Role:              role.NamaRole,
		NoHP:              input.NoHP,
		Status:            "aktif",
		RoleID:            role.ID,
		IsPasswordDefault: true,
		Menus:             menus,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun pegawai"})
		return
	}

	var existingEmployee models.Employee
	if err := db.Where("nip = ?", input.NIP).First(&existingEmployee).Error; err != nil {
		employee := models.Employee{
			NIP:                              input.NIP,
			Nama:                             input.Nama,
			StatusKepegawaian:                "Aktif",
			TanggalLahir:                     time.Now(),
			TanggalKgbTerakhir:               time.Now(),
			TanggalKgbBerikutnya:             time.Now().AddDate(2, 0, 0),
			TanggalKenaikanPangkatTerakhir:   time.Now(),
			TanggalKenaikanPangkatBerikutnya: time.Now().AddDate(4, 0, 0),
			TanggalPensiun:                   time.Now().AddDate(58, 0, 0),
		}
		db.Create(&employee)
	}

	adminNip, _ := c.Get("nip")
	db.Create(&models.RequestHistory{
		RequestType: "admin_create_pegawai",
		RequestID:   user.ID,
		StatusBaru:  "Success",
		Catatan:     "Created employee account: " + input.NIP + " (" + input.Nama + ")",
		ChangedBy:   adminNip.(string),
		ChangedAt:   time.Now(),
	})

	c.JSON(http.StatusCreated, gin.H{"message": "Pegawai berhasil dibuat dengan password default (NIP)"})
}

// AdminUpdateAkses updates pegawai role and menu access
func AdminUpdateAkses(c *gin.Context) {
	nip := c.Param("nip")
	var input UpdateAksesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()
	var user models.User
	if err := db.Where("nip = ?", nip).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pegawai tidak ditemukan"})
		return
	}

	var role models.Role
	if err := db.First(&role, input.RoleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak valid"})
		return
	}

	var menus []models.Menu
	if len(input.MenuIDs) > 0 {
		if err := db.Find(&menus, input.MenuIDs).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Daftar menu tidak valid"})
			return
		}
	}

	user.RoleID = role.ID
	user.Role = role.NamaRole
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui role"})
		return
	}

	if err := db.Model(&user).Association("Menus").Replace(&menus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui menu akses"})
		return
	}

	adminNip, _ := c.Get("nip")
	db.Create(&models.RequestHistory{
		RequestType: "admin_update_akses",
		RequestID:   user.ID,
		StatusBaru:  "Success",
		Catatan:     "Updated role/menu access for pegawai: " + nip,
		ChangedBy:   adminNip.(string),
		ChangedAt:   time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Akses role dan menu berhasil diperbarui"})
}

// AdminResetPassword resets a user's password back to their NIP
func AdminResetPassword(c *gin.Context) {
	nip := c.Param("nip")
	db := config.GetDB()

	var user models.User
	if err := db.Where("nip = ?", nip).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pegawai tidak ditemukan"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.NIP), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-hash password"})
		return
	}

	user.PasswordHash = string(hash)
	user.IsPasswordDefault = true
	db.Save(&user)

	adminNip, _ := c.Get("nip")
	db.Create(&models.RequestHistory{
		RequestType: "admin_reset_password",
		RequestID:   user.ID,
		StatusBaru:  "Success",
		Catatan:     "Reset password to NIP default for pegawai: " + nip,
		ChangedBy:   adminNip.(string),
		ChangedAt:   time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil di-reset kembali ke NIP default"})
}

// AdminGetMenu lists all available menus
func AdminGetMenu(c *gin.Context) {
	db := config.GetDB()
	var menus []models.Menu
	if err := db.Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data menu"})
		return
	}
	c.JSON(http.StatusOK, menus)
}
=======
	db.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Password untuk NIP " + input.NIP + " berhasil direset"})
}
>>>>>>> 603353f54c6625439da1b7cf09eb935c784c51b4

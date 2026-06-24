package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"sipecut/config"
	"sipecut/models"
	"sipecut/services"
)

// Pension Rules API
func GetPensionRules(c *gin.Context) {
	db := config.GetDB()
	var rules []models.PensionRule
	db.Order("jenis_jabatan asc, jabatan asc").Find(&rules)
	c.JSON(http.StatusOK, gin.H{"data": rules})
}

func UpdatePensionRule(c *gin.Context) {
	db := config.GetDB()
	var input struct {
		ID               uint   `json:"id"`
		JenisJabatan     string `json:"jenis_jabatan" binding:"required"`
		Jabatan          string `json:"jabatan" binding:"required"`
		BatasUsiaPensiun int    `json:"batas_usia_pensiun" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var rule models.PensionRule
	if input.ID > 0 {
		if err := db.First(&rule, input.ID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Aturan tidak ditemukan"})
			return
		}
	}

	rule.JenisJabatan = input.JenisJabatan
	rule.Jabatan = input.Jabatan
	rule.BatasUsiaPensiun = input.BatasUsiaPensiun
	rule.UpdatedAt = time.Now()

	db.Save(&rule)
	_ = services.RecalculateAllEmployees(db)
	c.JSON(http.StatusOK, gin.H{"message": "Aturan pensiun berhasil disimpan", "rule": rule})
}

func DeletePensionRule(c *gin.Context) {
	db := config.GetDB()
	idVal := c.Param("id")
	id, _ := strconv.Atoi(idVal)

	if err := db.Delete(&models.PensionRule{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus aturan pensiun: " + err.Error()})
		return
	}
	if err := services.RecalculateAllEmployees(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui kalkulasi pegawai: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Aturan pensiun berhasil dihapus"})
}

// KGB Cycle Rules API
func GetKgbRules(c *gin.Context) {
	db := config.GetDB()
	var rules []models.KgbCycleRule
	db.Order("jenis_jabatan asc, jabatan asc").Find(&rules)
	c.JSON(http.StatusOK, gin.H{"data": rules})
}

func UpdateKgbRule(c *gin.Context) {
	db := config.GetDB()
	var input struct {
		ID           uint   `json:"id"`
		JenisJabatan string `json:"jenis_jabatan" binding:"required"`
		Jabatan      string `json:"jabatan" binding:"required"`
		SiklusTahun  int    `json:"siklus_tahun" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var rule models.KgbCycleRule
	if input.ID > 0 {
		if err := db.First(&rule, input.ID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Aturan tidak ditemukan"})
			return
		}
	}

	rule.JenisJabatan = input.JenisJabatan
	rule.Jabatan = input.Jabatan
	rule.SiklusTahun = input.SiklusTahun
	rule.UpdatedAt = time.Now()

	db.Save(&rule)
	_ = services.RecalculateAllEmployees(db)
	c.JSON(http.StatusOK, gin.H{"message": "Aturan siklus KGB berhasil disimpan", "rule": rule})
}

func DeleteKgbRule(c *gin.Context) {
	db := config.GetDB()
	idVal := c.Param("id")
	id, _ := strconv.Atoi(idVal)

	if err := db.Delete(&models.KgbCycleRule{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus aturan KGB: " + err.Error()})
		return
	}
	if err := services.RecalculateAllEmployees(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui kalkulasi pegawai: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Aturan siklus KGB berhasil dihapus"})
}

// Pangkat Cycle Rules API
func GetPangkatRules(c *gin.Context) {
	db := config.GetDB()
	var rules []models.PangkatCycleRule
	db.Order("jenis_jabatan asc, jabatan asc").Find(&rules)
	c.JSON(http.StatusOK, gin.H{"data": rules})
}

func UpdatePangkatRule(c *gin.Context) {
	db := config.GetDB()
	var input struct {
		ID           uint   `json:"id"`
		JenisJabatan string `json:"jenis_jabatan" binding:"required"`
		Jabatan      string `json:"jabatan" binding:"required"`
		SiklusTahun  int    `json:"siklus_tahun" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var rule models.PangkatCycleRule
	if input.ID > 0 {
		if err := db.First(&rule, input.ID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Aturan tidak ditemukan"})
			return
		}
	}

	rule.JenisJabatan = input.JenisJabatan
	rule.Jabatan = input.Jabatan
	rule.SiklusTahun = input.SiklusTahun
	rule.UpdatedAt = time.Now()

	db.Save(&rule)
	_ = services.RecalculateAllEmployees(db)
	c.JSON(http.StatusOK, gin.H{"message": "Aturan kenaikan pangkat berhasil disimpan", "rule": rule})
}

func DeletePangkatRule(c *gin.Context) {
	db := config.GetDB()
	idVal := c.Param("id")
	id, _ := strconv.Atoi(idVal)

	if err := db.Delete(&models.PangkatCycleRule{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus aturan pangkat: " + err.Error()})
		return
	}
	if err := services.RecalculateAllEmployees(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui kalkulasi pegawai: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Aturan kenaikan pangkat berhasil dihapus"})
}

// General Settings API
func GetSettings(c *gin.Context) {
	db := config.GetDB()
	var settings []models.AppSetting
	db.Find(&settings)

	// Format as key-value JSON map for easier frontend consumption
	settingsMap := make(map[string]string)
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}

	c.JSON(http.StatusOK, gin.H{"data": settingsMap})
}

func UpdateSettings(c *gin.Context) {
	db := config.GetDB()
	var input map[string]string
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for k, v := range input {
		var setting models.AppSetting
		if err := db.Where("key = ?", k).First(&setting).Error; err == nil {
			setting.Value = v
			db.Save(&setting)
		} else {
			db.Create(&models.AppSetting{Key: k, Value: v})
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengaturan berhasil diperbarui"})
}

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

// SubmitBeritaAcara registers a new document upload (BA / ST / SI / SKS)
func SubmitBeritaAcara(c *gin.Context) {
	db := config.GetDB()
	nip, _ := c.Get("nip")

	var emp models.Employee
	if err := db.Where("nip = ?", nip).First(&emp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profil pegawai tidak ditemukan"})
		return
	}

	jenis := c.PostForm("jenis")
	if jenis != "BA" && jenis != "ST" && jenis != "SI" && jenis != "SKS" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jenis Berita Acara tidak valid. Harus berupa BA, ST, SI, atau SKS."})
		return
	}

	fileURL := handleFileUpload(c, "file_path")
	if fileURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File dokumen wajib diunggah"})
		return
	}

	ba := models.BeritaAcara{
		EmployeeID: emp.ID,
		Jenis:      jenis,
		FilePath:   fileURL,
		Status:     "Diajukan",
	}

	if err := db.Create(&ba).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengajuan dokumen: " + err.Error()})
		return
	}

	// Log History
	db.Create(&models.RequestHistory{
		RequestType: "berita_acara",
		RequestID:   ba.ID,
		StatusLama:  "",
		StatusBaru:  "Diajukan",
		Catatan:     "Mengunggah dokumen Berita Acara baru.",
		ChangedBy:   emp.NIP,
		ChangedAt:   time.Now(),
	})

	// Notify Admin
	var admins []models.User
	if err := db.Where("role = ?", "admin").Find(&admins).Error; err == nil {
		for _, admin := range admins {
			placeholders := map[string]string{
				"NAMA":  emp.Nama,
				"JENIS": jenis,
			}
			_ = services.SendWhatsAppNotification(db, &admin, "whatsapp_template_ba_new", placeholders)
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Dokumen Berita Acara berhasil dikirim", "berita_acara_id": ba.ID})
}

// GetBeritaAcara lists uploaded documents (supports Server-side DataTables parameters)
func GetBeritaAcara(c *gin.Context) {
	db := config.GetDB()
	role, _ := c.Get("role")
	nip, _ := c.Get("nip")

	drawVal := c.Query("draw")
	startVal := c.Query("start")
	lengthVal := c.Query("length")
	searchVal := c.Query("search[value]")
	statusFilter := c.Query("status")
	jenisFilter := c.Query("jenis")

	draw, _ := strconv.Atoi(drawVal)
	start, _ := strconv.Atoi(startVal)
	length, _ := strconv.Atoi(lengthVal)
	if length <= 0 {
		length = 10
	}

	var totalRecords int64
	var filteredRecords int64

	query := db.Model(&models.BeritaAcara{}).Preload("Employee")

	if role == "employee" {
		query = query.Joins("JOIN employees ON employees.id = berita_acara.employee_id").Where("employees.nip = ?", nip)
	}

	// Count Total
	query.Count(&totalRecords)

	if statusFilter != "" {
		query = query.Where("berita_acara.status = ?", statusFilter)
	}
	if jenisFilter != "" {
		query = query.Where("berita_acara.jenis = ?", jenisFilter)
	}

	if searchVal != "" {
		searchLike := "%" + searchVal + "%"
		query = query.Joins("LEFT JOIN employees emp_search ON emp_search.id = berita_acara.employee_id").
			Where("emp_search.nama LIKE ? OR emp_search.nip LIKE ? OR berita_acara.jenis LIKE ? OR berita_acara.status LIKE ?", searchLike, searchLike, searchLike, searchLike)
	}

	// Count Filtered
	query.Count(&filteredRecords)

	// Order
	query = query.Order("berita_acara.created_at DESC")

	var baList []models.BeritaAcara
	query.Limit(length).Offset(start).Find(&baList)

	c.JSON(http.StatusOK, gin.H{
		"draw":            draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            baList,
	})
}

// UpdateBeritaAcaraStatus handles admin approval or rejection (admin only)
func UpdateBeritaAcaraStatus(c *gin.Context) {
	db := config.GetDB()
	idVal := c.Param("id")
	id, _ := strconv.Atoi(idVal)

	var input struct {
		Status       string `json:"status" binding:"required"` // "Disetujui", "Dikembalikan", "Ditolak"
		CatatanAdmin string `json:"catatan_admin"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ba models.BeritaAcara
	if err := db.Preload("Employee").First(&ba, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dokumen Berita Acara tidak ditemukan"})
		return
	}

	oldStatus := ba.Status
	ba.Status = input.Status
	ba.CatatanAdmin = input.CatatanAdmin
	ba.UpdatedAt = time.Now()

	db.Save(&ba)

	adminNip, _ := c.Get("nip")
	db.Create(&models.RequestHistory{
		RequestType: "berita_acara",
		RequestID:   ba.ID,
		StatusLama:  oldStatus,
		StatusBaru:  input.Status,
		Catatan:     input.CatatanAdmin,
		ChangedBy:   adminNip.(string),
		ChangedAt:   time.Now(),
	})

	// Notify Employee via WhatsApp
	var user models.User
	if err := db.Where("nip = ?", ba.Employee.NIP).First(&user).Error; err == nil {
		placeholders := map[string]string{
			"NAMA":    ba.Employee.Nama,
			"JENIS":   ba.Jenis,
			"STATUS":  ba.Status,
			"CATATAN": ba.CatatanAdmin,
		}
		_ = services.SendWhatsAppNotification(db, &user, "whatsapp_template_ba_update", placeholders)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status dokumen Berita Acara berhasil diperbarui", "berita_acara": ba})
}

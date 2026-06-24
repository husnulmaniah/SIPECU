package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"sipecut/config"
	"sipecut/models"
	"sipecut/services"
)

// SubmitDataChangeRequest registers a self-service request to update details
func SubmitDataChangeRequest(c *gin.Context) {
	db := config.GetDB()
	nip, _ := c.Get("nip")

	var emp models.Employee
	if err := db.Where("nip = ?", nip).First(&emp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profil pegawai tidak ditemukan"})
		return
	}

	dataJSON := c.PostForm("data_json") // JSON format containing changed fields
	if dataJSON == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field data_json wajib diisi"})
		return
	}

	// Validate JSON structure
	var checkMap map[string]interface{}
	if err := json.Unmarshal([]byte(dataJSON), &checkMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data_json tidak valid"})
		return
	}

	skTerakhirURL := handleFileUpload(c, "sk_terakhir_file")
	if skTerakhirURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File bukti SK Terakhir wajib diunggah"})
		return
	}

	request := models.DataChangeRequest{
		EmployeeID:     emp.ID,
		DataJSON:       dataJSON,
		SkTerakhirFile: skTerakhirURL,
		Status:         "Diajukan",
	}

	if err := db.Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengajuan perubahan data: " + err.Error()})
		return
	}

	// Log History
	db.Create(&models.RequestHistory{
		RequestType: "change",
		RequestID:   request.ID,
		StatusLama:  "",
		StatusBaru:  "Diajukan",
		Catatan:     "Pengajuan perubahan data diajukan dengan bukti dokumen SK.",
		ChangedBy:   emp.NIP,
		ChangedAt:   time.Now(),
	})

	// Notify Admin
	var admins []models.User
	if err := db.Where("role = ?", "admin").Find(&admins).Error; err == nil {
		for _, admin := range admins {
			placeholders := map[string]string{
				"NAMA": emp.Nama,
				"NIP":  emp.NIP,
			}
			_ = services.SendWhatsAppNotification(db, &admin, "whatsapp_template_change_new", placeholders)
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pengajuan perubahan data berhasil dikirim", "request_id": request.ID})
}

// GetDataChangeRequests lists change requests (supports Server-side DataTables parameters)
func GetDataChangeRequests(c *gin.Context) {
	db := config.GetDB()
	role, _ := c.Get("role")
	nip, _ := c.Get("nip")

	drawVal := c.Query("draw")
	startVal := c.Query("start")
	lengthVal := c.Query("length")
	searchVal := c.Query("search[value]")
	statusFilter := c.Query("status")

	draw, _ := strconv.Atoi(drawVal)
	start, _ := strconv.Atoi(startVal)
	length, _ := strconv.Atoi(lengthVal)
	if length <= 0 {
		length = 10
	}

	var totalRecords int64
	var filteredRecords int64

	query := db.Model(&models.DataChangeRequest{}).Preload("Employee")

	if role == "employee" {
		query = query.Joins("JOIN employees ON employees.id = data_change_requests.employee_id").Where("employees.nip = ?", nip)
	}

	// Count Total
	query.Count(&totalRecords)

	if statusFilter != "" {
		query = query.Where("data_change_requests.status = ?", statusFilter)
	}

	if searchVal != "" {
		searchLike := "%" + searchVal + "%"
		query = query.Joins("LEFT JOIN employees emp_search ON emp_search.id = data_change_requests.employee_id").
			Where("emp_search.nama LIKE ? OR emp_search.nip LIKE ? OR data_change_requests.status LIKE ?", searchLike, searchLike, searchLike)
	}

	// Count Filtered
	query.Count(&filteredRecords)

	// Order
	query = query.Order("data_change_requests.created_at DESC")

	var requests []models.DataChangeRequest
	query.Limit(length).Offset(start).Find(&requests)

	c.JSON(http.StatusOK, gin.H{
		"draw":            draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            requests,
	})
}

// UpdateDataChangeRequestStatus processes approval or rejection (admin only)
func UpdateDataChangeRequestStatus(c *gin.Context) {
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

	var req models.DataChangeRequest
	if err := db.Preload("Employee").First(&req, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengajuan perubahan data tidak ditemukan"})
		return
	}

	oldStatus := req.Status
	req.Status = input.Status
	req.CatatanAdmin = input.CatatanAdmin
	req.UpdatedAt = time.Now()

	adminNip, _ := c.Get("nip")

	// If approved, apply the modifications to the employee profile!
	if input.Status == "Disetujui" {
		var changes map[string]interface{}
		if err := json.Unmarshal([]byte(req.DataJSON), &changes); err == nil {
			var emp models.Employee
			if err := db.First(&emp, req.EmployeeID).Error; err == nil {
				// Safely apply changes
				if val, ok := changes["nama"]; ok {
					emp.Nama = val.(string)
				}
				if val, ok := changes["jenis_jabatan"]; ok {
					emp.JenisJabatan = val.(string)
				}
				if val, ok := changes["jabatan"]; ok {
					emp.Jabatan = val.(string)
				}
				if val, ok := changes["tempat_lahir"]; ok {
					emp.TempatLahir = val.(string)
				}
				if val, ok := changes["tanggal_lahir"]; ok {
					emp.TanggalLahir = parseDateString(val.(string))
				}
				if val, ok := changes["tempat_tugas"]; ok {
					emp.TempatTugas = val.(string)
				}
				if val, ok := changes["jenis_tempat"]; ok {
					emp.JenisTempat = val.(string)
				}
				if val, ok := changes["pengangkatan"]; ok {
					emp.Pengangkatan = val.(string)
				}
				if val, ok := changes["jenis_pengangkatan"]; ok {
					emp.JenisPengangkatan = val.(string)
				}

				// If the change includes updating KGB or rank dates, also record history!
				if val, ok := changes["tanggal_kgb_terakhir"]; ok {
					oldKgb := emp.TanggalKgbTerakhir
					emp.TanggalKgbTerakhir = parseDateString(val.(string))
					emp.SkKgbFile = req.SkTerakhirFile
					
					// Save history
					db.Create(&models.EmployeeKgbHistory{
						EmployeeID: emp.ID,
						TanggalKgb: oldKgb,
						FileSkKgb:  emp.SkKgbFile,
					})
				}

				if val, ok := changes["tanggal_kenaikan_pangkat_terakhir"]; ok {
					oldPangkat := emp.TanggalKenaikanPangkatTerakhir
					emp.TanggalKenaikanPangkatTerakhir = parseDateString(val.(string))
					emp.SkPangkatFile = req.SkTerakhirFile

					// Save history
					db.Create(&models.EmployeePangkatHistory{
						EmployeeID:            emp.ID,
						TanggalKenaikanPangkat: oldPangkat,
						FileSkPangkat:         emp.SkPangkatFile,
					})
				}

				// Recalculate employee dates
				_ = services.RecalculateEmployeeDates(db, &emp)
				db.Save(&emp)
			}
		}
	}

	db.Save(&req)

	// Log History
	db.Create(&models.RequestHistory{
		RequestType: "change",
		RequestID:   req.ID,
		StatusLama:  oldStatus,
		StatusBaru:  input.Status,
		Catatan:     input.CatatanAdmin,
		ChangedBy:   adminNip.(string),
		ChangedAt:   time.Now(),
	})

	// Notify Employee via WhatsApp
	var user models.User
	if err := db.Where("nip = ?", req.Employee.NIP).First(&user).Error; err == nil {
		placeholders := map[string]string{
			"NAMA":    req.Employee.Nama,
			"STATUS":  req.Status,
			"CATATAN": req.CatatanAdmin,
		}
		_ = services.SendWhatsAppNotification(db, &user, "whatsapp_template_change_update", placeholders)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status pengajuan perubahan data berhasil diperbarui", "request": req})
}

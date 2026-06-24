package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"sipecut/config"
	"sipecut/models"
	"sipecut/services"
	"sipecut/storage"
)

// SubmitLeaveRequest creates a new leave application with dynamic file validation
// Admin can submit on behalf of an employee by providing employee_nip in the form.
func SubmitLeaveRequest(c *gin.Context) {
	db := config.GetDB()
	callerNip, _ := c.Get("nip")
	callerRole, _ := c.Get("role")

	// Determine which employee the leave is for
	targetNip := callerNip.(string)
	if callerRole == "admin" {
		if empNip := c.PostForm("employee_nip"); empNip != "" {
			targetNip = empNip
		}
	}

	var emp models.Employee
	if err := db.Where("nip = ?", targetNip).First(&emp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profil pegawai tidak ditemukan untuk NIP: " + targetNip})
		return
	}

	jenisCuti := c.PostForm("jenis_cuti")
	tanggalMulaiStr := c.PostForm("tanggal_mulai")
	tanggalSelesaiStr := c.PostForm("tanggal_selesai")

	if jenisCuti == "" || tanggalMulaiStr == "" || tanggalSelesaiStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jenis Cuti, Tanggal Mulai, dan Tanggal Selesai wajib diisi"})
		return
	}

	tanggalMulai := parseDateString(tanggalMulaiStr)
	tanggalSelesai := parseDateString(tanggalSelesaiStr)

	if tanggalMulai.After(tanggalSelesai) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal mulai tidak boleh melebihi tanggal selesai"})
		return
	}

	// Create Request
	request := models.LeaveRequest{
		EmployeeID:     emp.ID,
		JenisCuti:      jenisCuti,
		TanggalMulai:   tanggalMulai,
		TanggalSelesai: tanggalSelesai,
		Status:         "Diajukan",
	}

	// Begin GORM Transaction
	tx := db.Begin()


	if err := tx.Create(&request).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pengajuan cuti: " + err.Error()})
		return
	}

	// Helper to validate and upload
	uploadAttachment := func(fieldName string, isRequired bool, docLabel string) bool {
		fileURL := handleFileUpload(c, fieldName)
		if fileURL == "" {
			if isRequired {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Berkas '%s' wajib diunggah untuk jenis cuti %s", docLabel, jenisCuti)})
				c.Abort()
				return false
			}
			return true // optional not provided, valid
		}

		attachment := models.LeaveAttachment{
			LeaveRequestID: request.ID,
			JenisDokumen:   docLabel,
			FilePath:       fileURL,
		}
		if err := tx.Create(&attachment).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan lampiran berkas: " + err.Error()})
			c.Abort()
			return false
		}
		return true
	}

	// Dynamic validation based on Leave Type
	// Cuti Tahunan Biasa: Surat Rekomendasi Kepsek (wajib), SK Terakhir (wajib)
	// Cuti Melahirkan: Surat Rekomendasi Kepsek (wajib), SK Terakhir (wajib), HPL (wajib), Buku KIA (wajib), USG (opsional)
	// Cuti Tahunan Umroh: Surat Rekomendasi Kepsek (wajib), SK Terakhir (wajib), Ket Travel (wajib)
	// Cuti Sakit: SK Terakhir (wajib), Surat Rujukan (wajib), Rawat Inap (wajib)
	// Cuti Alasan Penting: SK Terakhir (wajib), Surat Rekomendasi Kepsek (wajib), Dokumen Pendukung (opsional/fleksibel)
	var valid bool
	switch jenisCuti {
	case "Cuti Tahunan Biasa":
		valid = uploadAttachment("surat_rekomendasi_kepsek", true, "Surat rekomendasi Kepala Sekolah") &&
			uploadAttachment("sk_terakhir", true, "SK terakhir")
	case "Cuti Melahirkan":
		valid = uploadAttachment("surat_rekomendasi_kepsek", true, "Surat rekomendasi Kepala Sekolah") &&
			uploadAttachment("sk_terakhir", true, "SK terakhir") &&
			uploadAttachment("surat_hpl", true, "Surat keterangan HPL") &&
			uploadAttachment("buku_kia", true, "Buku KIA") &&
			uploadAttachment("hasil_usg", false, "Hasil USG")
	case "Cuti Tahunan untuk Umroh":
		valid = uploadAttachment("surat_rekomendasi_kepsek", true, "Surat rekomendasi Kepala Sekolah") &&
			uploadAttachment("sk_terakhir", true, "SK terakhir") &&
			uploadAttachment("surat_travel", true, "Surat keterangan dari travel pemberangkatan")
	case "Cuti Sakit":
		valid = uploadAttachment("sk_terakhir", true, "SK terakhir") &&
			uploadAttachment("surat_rujukan", true, "Surat rujukan") &&
			uploadAttachment("surat_rawat_inap", true, "Surat keterangan rawat inap")
	case "Cuti Alasan Penting":
		valid = uploadAttachment("sk_terakhir", true, "SK terakhir") &&
			uploadAttachment("surat_rekomendasi_kepsek", true, "Surat rekomendasi Kepala Sekolah") &&
			uploadAttachment("dokumen_pendukung", true, "Dokumen pendukung")
	default:
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jenis cuti tidak valid"})
		return
	}

	if c.IsAborted() || !valid {
		return
	}

	// Commit Transaction
	tx.Commit()

	// Log History
	db.Create(&models.RequestHistory{
		RequestType: "leave",
		RequestID:   request.ID,
		StatusLama:  "",
		StatusBaru:  "Diajukan",
		Catatan:     "Pengajuan cuti baru diajukan oleh pegawai.",
		ChangedBy:   emp.NIP,
		ChangedAt:   time.Now(),
	})

	// Notify Admin (fetch admin accounts)
	var admins []models.User
	if err := db.Where("role = ?", "admin").Find(&admins).Error; err == nil {
		for _, admin := range admins {
			placeholders := map[string]string{
				"NAMA":       emp.Nama,
				"NIP":        emp.NIP,
				"JENIS_CUTI": jenisCuti,
			}
			_ = services.SendWhatsAppNotification(db, &admin, "whatsapp_template_leave_new", placeholders)
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pengajuan cuti berhasil dikirim", "request_id": request.ID})
}

// GetLeaveRequests lists leave requests with Server-Side DataTables support
func GetLeaveRequests(c *gin.Context) {
	db := config.GetDB()
	role, _ := c.Get("role")
	nip, _ := c.Get("nip")

	drawVal := c.Query("draw")
	startVal := c.Query("start")
	lengthVal := c.Query("length")
	searchVal := c.Query("search[value]")
	statusFilter := c.Query("status")
	jenisCutiFilter := c.Query("jenis_cuti")

	draw, _ := strconv.Atoi(drawVal)
	start, _ := strconv.Atoi(startVal)
	length, _ := strconv.Atoi(lengthVal)
	if length <= 0 {
		length = 10
	}

	var totalRecords int64
	var filteredRecords int64

	// Base query
	query := db.Model(&models.LeaveRequest{}).Preload("Employee").Preload("Attachments").Preload("Letters")

	// Filter by employee NIP if role is employee
	if role == "employee" {
		query = query.Joins("JOIN employees ON employees.id = leave_requests.employee_id").Where("employees.nip = ?", nip)
	}

	// Filter by NIP query if requested (admin only)
	nipQuery := c.Query("nip")
	if role == "admin" && nipQuery != "" {
		query = query.Joins("JOIN employees ON employees.id = leave_requests.employee_id").Where("employees.nip = ?", nipQuery)
	}

	// Count Total
	query.Count(&totalRecords)

	// Apply Filter Status & Jenis Cuti
	if statusFilter != "" {
		query = query.Where("leave_requests.status = ?", statusFilter)
	}
	if jenisCutiFilter != "" {
		query = query.Where("leave_requests.jenis_cuti = ?", jenisCutiFilter)
	}

	// Apply Search
	if searchVal != "" {
		searchLike := "%" + searchVal + "%"
		// If employee table isn't joined yet, join it for search
		query = query.Joins("LEFT JOIN employees emp_search ON emp_search.id = leave_requests.employee_id").
			Where("emp_search.nama LIKE ? OR emp_search.nip LIKE ? OR leave_requests.jenis_cuti LIKE ? OR leave_requests.status LIKE ?", searchLike, searchLike, searchLike, searchLike)
	}

	// Count Filtered
	query.Count(&filteredRecords)

	// Sort
	query = query.Order("leave_requests.created_at DESC")

	// Paginate
	var requests []models.LeaveRequest
	query.Limit(length).Offset(start).Find(&requests)

	c.JSON(http.StatusOK, gin.H{
		"draw":            draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            requests,
	})
}

// UpdateLeaveStatus processes leave request approval (admin only)
func UpdateLeaveStatus(c *gin.Context) {
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

	var req models.LeaveRequest
	if err := db.Preload("Employee").First(&req, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengajuan cuti tidak ditemukan"})
		return
	}

	oldStatus := req.Status
	req.Status = input.Status
	req.CatatanAdmin = input.CatatanAdmin
	req.UpdatedAt = time.Now()

	adminNip, _ := c.Get("nip")

	// If approved, automatically generate DOCX and PDF drafts for both documents
	if input.Status == "Disetujui" {
		type docJob struct {
			jenis    string
			nameBase string
			genDocx  func(*models.LeaveRequest) ([]byte, error)
			genPdf   func(*models.LeaveRequest) ([]byte, error)
		}

		jobs := []docJob{
			{
				jenis:    "Rekomendasi",
				nameBase: fmt.Sprintf("Rekomendasi_Cuti_%d", req.ID),
				genDocx:  services.GenerateLeaveDocx,
				genPdf:   services.GenerateLeavePdf,
			},
			{
				jenis:    "Formulir",
				nameBase: fmt.Sprintf("Formulir_Cuti_%d", req.ID),
				genDocx:  services.GenerateFormulirDocx,
				genPdf:   services.GenerateFormulirCutiPdf,
			},
		}

		for _, job := range jobs {
			// Rekomendasi: generate both DOCX + PDF
			// Formulir: generate PDF only (no DOCX per requirement)
			var docxURL string
			if job.jenis == "Rekomendasi" {
				docxBytes, err := job.genDocx(&req)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate Word " + job.jenis + ": " + err.Error()})
					return
				}
				if docxBytes != nil {
					docxURL, err = storage.CurrentProvider.UploadFile(job.nameBase+".docx", bytes.NewReader(docxBytes))
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload Word " + job.jenis + ": " + err.Error()})
						return
					}
				}
			}

			pdfBytes, err := job.genPdf(&req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate PDF " + job.jenis + ": " + err.Error()})
				return
			}

			pdfURL, err := storage.CurrentProvider.UploadFile(job.nameBase+".pdf", bytes.NewReader(pdfBytes))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload PDF " + job.jenis + ": " + err.Error()})
				return
			}

			// Upsert LeaveLetter per jenis surat
			var existing models.LeaveLetter
			if err := db.Where("leave_request_id = ? AND jenis_surat = ?", req.ID, job.jenis).First(&existing).Error; err == nil {
				existing.FileDocx = docxURL
				existing.FilePdf = pdfURL
				db.Save(&existing)
			} else {
				db.Create(&models.LeaveLetter{
					LeaveRequestID: req.ID,
					JenisSurat:     job.jenis,
					FileDocx:       docxURL,
					FilePdf:        pdfURL,
				})
			}
		}
	}

	db.Save(&req)

	// Log History
	db.Create(&models.RequestHistory{
		RequestType: "leave",
		RequestID:   req.ID,
		StatusLama:  oldStatus,
		StatusBaru:  input.Status,
		Catatan:     input.CatatanAdmin,
		ChangedBy:   adminNip.(string),
		ChangedAt:   time.Now(),
	})

	// Get Employee User for Notification
	var user models.User
	if err := db.Where("nip = ?", req.Employee.NIP).First(&user).Error; err == nil {
		placeholders := map[string]string{
			"NAMA":       req.Employee.Nama,
			"JENIS_CUTI": req.JenisCuti,
			"STATUS":     req.Status,
			"CATATAN":     req.CatatanAdmin,
		}
		_ = services.SendWhatsAppNotification(db, &user, "whatsapp_template_leave_update", placeholders)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status pengajuan cuti berhasil diperbarui", "request": req})
}

// UploadSignedLetter uploads signed PDF document by Admin
func UploadSignedLetter(c *gin.Context) {
	db := config.GetDB()
	idVal := c.Param("id")
	id, _ := strconv.Atoi(idVal)

	var req models.LeaveRequest
	if err := db.Preload("Employee").First(&req, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengajuan cuti tidak ditemukan"})
		return
	}

	if req.Status != "Disetujui" && req.Status != "Surat Terunggah" && req.Status != "Selesai" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Surat rekomendasi hanya dapat diunggah setelah pengajuan disetujui"})
		return
	}

	signedURL := handleFileUpload(c, "file_signed_pdf")
	if signedURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File surat rekomendasi bertanda tangan wajib diunggah"})
		return
	}

	var letter models.LeaveLetter
	if err := db.Where("leave_request_id = ?", req.ID).First(&letter).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Draft surat rekomendasi belum dibuat. Silakan setujui pengajuan terlebih dahulu."})
		return
	}

	letter.FileSignedPdf = signedURL
	letter.UploadedByAdminAt = time.Now()
	db.Save(&letter)

	oldStatus := req.Status
	req.Status = "Surat Terunggah" // Set status to Surat Terunggah (final)
	db.Save(&req)

	adminNip, _ := c.Get("nip")
	db.Create(&models.RequestHistory{
		RequestType: "leave",
		RequestID:   req.ID,
		StatusLama:  oldStatus,
		StatusBaru:  "Surat Terunggah",
		Catatan:     "Surat rekomendasi cuti final bertanda tangan telah diunggah oleh Admin.",
		ChangedBy:   adminNip.(string),
		ChangedAt:   time.Now(),
	})

	// Notify Employee via WhatsApp
	var user models.User
	if err := db.Where("nip = ?", req.Employee.NIP).First(&user).Error; err == nil {
		messageText := fmt.Sprintf("Halo %s, Surat Rekomendasi Cuti final bertanda tangan untuk pengajuan %s Anda telah diunggah oleh Admin. Silakan unduh melalui akun SIPECUT Anda.", req.Employee.Nama, req.JenisCuti)
		
		// Create simulation / send direct
		placeholders := map[string]string{
			"NAMA":       req.Employee.Nama,
			"JENIS_CUTI": req.JenisCuti,
			"STATUS":     "Selesai (Surat Rekomendasi Terunggah)",
			"CATATAN":     "Silakan unduh dokumen Anda di menu Cuti.",
		}
		_ = services.SendWhatsAppNotification(db, &user, "whatsapp_template_leave_update", placeholders)
		
		// Overwrite message to specific details in console if simulated
		_ = messageText
	}

	c.JSON(http.StatusOK, gin.H{"message": "Surat rekomendasi ACC berhasil diunggah", "letter": letter})
}

// RegenerateLeaveDocuments regenerates both Formulir and Rekomendasi documents for an approved leave request
func RegenerateLeaveDocuments(c *gin.Context) {
	db := config.GetDB()
	idVal := c.Param("id")
	id, _ := strconv.Atoi(idVal)

	var req models.LeaveRequest
	if err := db.Preload("Employee").First(&req, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengajuan cuti tidak ditemukan"})
		return
	}

	if req.Status != "Disetujui" && req.Status != "Surat Terunggah" && req.Status != "Selesai" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dokumen hanya dapat di-regenerasi setelah pengajuan disetujui"})
		return
	}

	type docJob struct {
		jenis    string
		nameBase string
		genDocx  func(*models.LeaveRequest) ([]byte, error)
		genPdf   func(*models.LeaveRequest) ([]byte, error)
	}

	jobs := []docJob{
		{
			jenis:    "Rekomendasi",
			nameBase: fmt.Sprintf("Rekomendasi_Cuti_%d", req.ID),
			genDocx:  services.GenerateLeaveDocx,
			genPdf:   services.GenerateLeavePdf,
		},
		{
			jenis:    "Formulir",
			nameBase: fmt.Sprintf("Formulir_Cuti_%d", req.ID),
			genDocx:  services.GenerateFormulirDocx,
			genPdf:   services.GenerateFormulirCutiPdf,
		},
	}

	for _, job := range jobs {
		docxBytes, err := job.genDocx(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate DOCX " + job.jenis + ": " + err.Error()})
			return
		}

		pdfBytes, err := job.genPdf(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate PDF " + job.jenis + ": " + err.Error()})
			return
		}

		docxURL, err := storage.CurrentProvider.UploadFile(job.nameBase+".docx", bytes.NewReader(docxBytes))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload DOCX " + job.jenis + ": " + err.Error()})
			return
		}

		pdfURL, err := storage.CurrentProvider.UploadFile(job.nameBase+".pdf", bytes.NewReader(pdfBytes))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload PDF " + job.jenis + ": " + err.Error()})
			return
		}

		var existing models.LeaveLetter
		if err := db.Where("leave_request_id = ? AND jenis_surat = ?", req.ID, job.jenis).First(&existing).Error; err == nil {
			existing.FileDocx = docxURL
			existing.FilePdf = pdfURL
			db.Save(&existing)
		} else {
			db.Create(&models.LeaveLetter{
				LeaveRequestID: req.ID,
				JenisSurat:     job.jenis,
				FileDocx:       docxURL,
				FilePdf:        pdfURL,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dokumen berhasil di-regenerasi"})
}


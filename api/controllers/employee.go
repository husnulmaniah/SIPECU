package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"sipecut/config"
	"sipecut/models"
	"sipecut/services"
	"sipecut/storage"
)

// Helper to parse dates
func parseDateString(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}
	}
	return t
}

// GetEmployees retrieves active employees with Server-side DataTables or single user details
func GetEmployees(c *gin.Context) {
	db := config.GetDB()
	role, _ := c.Get("role")
	nip, _ := c.Get("nip")

	// If user is employee, only return their own data
	if role == "employee" {
		var emp models.Employee
		if err := db.Where("nip = ?", nip).First(&emp).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profil pegawai tidak ditemukan"})
			return
		}

		// Calculate age dynamically
		years, months := services.CalculateAge(emp.TanggalLahir)
		ageStr := fmt.Sprintf("%d Tahun %d Bulan", years, months)

		c.JSON(http.StatusOK, gin.H{
			"data": []interface{}{
				gin.H{
					"employee": emp,
					"umur":     ageStr,
				},
			},
		})
		return
	}

	// Dropdown or autocomplete request (no DataTables pagination requested)
	if c.Query("all") == "true" || c.Query("draw") == "" {
		var employees []models.Employee
		if err := db.Where("status_kepegawaian = ?", "Aktif").Order("nama ASC").Find(&employees).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pegawai: " + err.Error()})
			return
		}
		
		// Format response data with calculated age
		type EmployeeWithAge struct {
			models.Employee
			Umur string `json:"umur"`
		}

		resultData := make([]EmployeeWithAge, len(employees))
		for i, emp := range employees {
			years, months := services.CalculateAge(emp.TanggalLahir)
			resultData[i] = EmployeeWithAge{
				Employee: emp,
				Umur:     fmt.Sprintf("%d Tahun %d Bulan", years, months),
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"data": resultData,
		})
		return
	}

	// Admin: Server-Side DataTables processing
	drawVal := c.Query("draw")
	startVal := c.Query("start")
	lengthVal := c.Query("length")
	searchVal := c.Query("search[value]")
	orderColumnVal := c.Query("order[0][column]")
	orderDirVal := c.Query("order[0][dir]")

	draw, _ := strconv.Atoi(drawVal)
	start, _ := strconv.Atoi(startVal)
	length, _ := strconv.Atoi(lengthVal)
	if length <= 0 {
		length = 10
	}

	var totalRecords int64
	var filteredRecords int64

	// Base query for Active employees
	query := db.Model(&models.Employee{}).Where("status_kepegawaian = ?", "Aktif")

	// Count Total
	query.Count(&totalRecords)

	// Apply Search
	if searchVal != "" {
		searchLike := "%" + searchVal + "%"
		query = query.Where("nip LIKE ? OR nama LIKE ? OR jabatan LIKE ? OR tempat_tugas LIKE ?", searchLike, searchLike, searchLike, searchLike)
	}

	// Count Filtered
	query.Count(&filteredRecords)

	// Apply Order
	columns := []string{"nip", "nama", "jenis_jabatan", "jabatan", "tempat_tugas", "status_kepegawaian"}
	orderColIndex, err := strconv.Atoi(orderColumnVal)
	if err == nil && orderColIndex >= 0 && orderColIndex < len(columns) {
		orderBy := columns[orderColIndex]
		if orderDirVal == "desc" {
			query = query.Order(orderBy + " DESC")
		} else {
			query = query.Order(orderBy + " ASC")
		}
	} else {
		query = query.Order("nama ASC") // default order
	}

	// Apply Pagination
	var employees []models.Employee
	query.Limit(length).Offset(start).Find(&employees)

	// Format response data with calculated age
	type EmployeeWithAge struct {
		models.Employee
		Umur string `json:"umur"`
	}

	resultData := make([]EmployeeWithAge, len(employees))
	for i, emp := range employees {
		years, months := services.CalculateAge(emp.TanggalLahir)
		resultData[i] = EmployeeWithAge{
			Employee: emp,
			Umur:     fmt.Sprintf("%d Tahun %d Bulan", years, months),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"draw":            draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            resultData,
	})
}

// GetRetiredEmployees retrieves retired employees with Server-side DataTables
func GetRetiredEmployees(c *gin.Context) {
	db := config.GetDB()

	drawVal := c.Query("draw")
	startVal := c.Query("start")
	lengthVal := c.Query("length")
	searchVal := c.Query("search[value]")
	orderColumnVal := c.Query("order[0][column]")
	orderDirVal := c.Query("order[0][dir]")

	draw, _ := strconv.Atoi(drawVal)
	start, _ := strconv.Atoi(startVal)
	length, _ := strconv.Atoi(lengthVal)
	if length <= 0 {
		length = 10
	}

	var totalRecords int64
	var filteredRecords int64

	// Base query for Retired employees
	query := db.Model(&models.Employee{}).Where("status_kepegawaian = ?", "Pensiun")

	// Count Total
	query.Count(&totalRecords)

	// Apply Search
	if searchVal != "" {
		searchLike := "%" + searchVal + "%"
		query = query.Where("nip LIKE ? OR nama LIKE ? OR jabatan LIKE ? OR tempat_tugas LIKE ?", searchLike, searchLike, searchLike, searchLike)
	}

	// Count Filtered
	query.Count(&filteredRecords)

	// Apply Order
	columns := []string{"nip", "nama", "jenis_jabatan", "jabatan", "tempat_tugas", "status_kepegawaian"}
	orderColIndex, err := strconv.Atoi(orderColumnVal)
	if err == nil && orderColIndex >= 0 && orderColIndex < len(columns) {
		orderBy := columns[orderColIndex]
		if orderDirVal == "desc" {
			query = query.Order(orderBy + " DESC")
		} else {
			query = query.Order(orderBy + " ASC")
		}
	} else {
		query = query.Order("nama ASC")
	}

	// Apply Pagination
	var employees []models.Employee
	query.Limit(length).Offset(start).Find(&employees)

	// Format response
	type EmployeeWithAge struct {
		models.Employee
		Umur string `json:"umur"`
	}

	resultData := make([]EmployeeWithAge, len(employees))
	for i, emp := range employees {
		years, months := services.CalculateAge(emp.TanggalLahir)
		resultData[i] = EmployeeWithAge{
			Employee: emp,
			Umur:     fmt.Sprintf("%d Tahun %d Bulan", years, months),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"draw":            draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            resultData,
	})
}

// CreateEmployee registers a new employee (admin only)
func CreateEmployee(c *gin.Context) {
	db := config.GetDB()

	// Parse fields
	nip := c.PostForm("nip")
	nama := c.PostForm("nama")
	jenisJabatan := services.NormalizeJenisJabatan(c.PostForm("jenis_jabatan"))
	jabatan := c.PostForm("jabatan")
	tempatLahir := c.PostForm("tempat_lahir")
	tanggalLahirStr := c.PostForm("tanggal_lahir")
	tempatTugas := c.PostForm("tempat_tugas")
	jenisTempat := c.PostForm("jenis_tempat")
	pengangkatan := c.PostForm("pengangkatan")
	tanggalKgbStr := c.PostForm("tanggal_kgb_terakhir")
	tanggalPangkatStr := c.PostForm("tanggal_kenaikan_pangkat_terakhir")
	jenisPengangkatan := c.PostForm("jenis_pengangkatan")
	noHP := c.PostForm("no_hp")

	if nip == "" || nama == "" || jenisJabatan == "" || jabatan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIP, Nama, Jenis Jabatan, dan Jabatan wajib diisi"})
		return
	}

	// Check if already exists
	var count int64
	db.Model(&models.User{}).Where("nip = ?", nip).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIP sudah terdaftar di sistem"})
		return
	}

	// Create User Account with default password "pegawai123"
	passHash, _ := bcrypt.GenerateFromPassword([]byte("pegawai123"), bcrypt.DefaultCost)
	newUser := models.User{
		NIP:          nip,
		PasswordHash: string(passHash),
		Role:         "employee",
		NoHP:         noHP,
		Status:       "aktif",
	}
	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun pengguna: " + err.Error()})
		return
	}

	// Handle uploads
	fotoURL := handleFileUpload(c, "foto_profil")
	skCpnsURL := handleFileUpload(c, "sk_cpns_pppk_file")
	skPnsURL := handleFileUpload(c, "sk_pns_file")
	skKgbURL := handleFileUpload(c, "sk_kgb_file")
	skPangkatURL := handleFileUpload(c, "sk_pangkat_file")

	employee := models.Employee{
		NIP:                            nip,
		Nama:                           nama,
		JenisJabatan:                   jenisJabatan,
		Jabatan:                        jabatan,
		TempatLahir:                    tempatLahir,
		TanggalLahir:                   parseDateString(tanggalLahirStr),
		TempatTugas:                    tempatTugas,
		JenisTempat:                    jenisTempat,
		Pengangkatan:                   pengangkatan,
		TanggalKgbTerakhir:             parseDateString(tanggalKgbStr),
		TanggalKenaikanPangkatTerakhir: parseDateString(tanggalPangkatStr),
		JenisPengangkatan:              jenisPengangkatan,
		FotoProfil:                     fotoURL,
		SkCpnsPppkFile:                 skCpnsURL,
		SkPnsFile:                      skPnsURL,
		SkKgbFile:                      skKgbURL,
		SkPangkatFile:                  skPangkatURL,
	}

	// Apply calculations
	_ = services.RecalculateEmployeeDates(db, &employee)

	if err := db.Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat profil pegawai: " + err.Error()})
		return
	}

	// Insert initial histories
	if !employee.TanggalKgbTerakhir.IsZero() {
		db.Create(&models.EmployeeKgbHistory{
			EmployeeID: employee.ID,
			TanggalKgb: employee.TanggalKgbTerakhir,
			FileSkKgb:  employee.SkKgbFile,
		})
	}
	if !employee.TanggalKenaikanPangkatTerakhir.IsZero() {
		db.Create(&models.EmployeePangkatHistory{
			EmployeeID:            employee.ID,
			TanggalKenaikanPangkat: employee.TanggalKenaikanPangkatTerakhir,
			FileSkPangkat:         employee.SkPangkatFile,
		})
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pegawai berhasil ditambahkan", "employee": employee})
}

// GetEmployeeByNIP gets details for a single employee
func GetEmployeeByNIP(c *gin.Context) {
	db := config.GetDB()
	nip := c.Param("nip")

	// RBAC Check
	userRole, _ := c.Get("role")
	userNip, _ := c.Get("nip")
	if userRole == "employee" && userNip != nip {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak diizinkan mengakses data pegawai lain"})
		return
	}

	var emp models.Employee
	if err := db.Where("nip = ?", nip).First(&emp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pegawai tidak ditemukan"})
		return
	}

	years, months := services.CalculateAge(emp.TanggalLahir)
	ageStr := fmt.Sprintf("%d Tahun %d Bulan", years, months)

	c.JSON(http.StatusOK, gin.H{
		"employee": emp,
		"umur":     ageStr,
	})
}

// UpdateEmployee updates employee data directly (admin only)
func UpdateEmployee(c *gin.Context) {
	db := config.GetDB()
	nip := c.Param("nip")

	var emp models.Employee
	if err := db.Where("nip = ?", nip).First(&emp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pegawai tidak ditemukan"})
		return
	}

	// Parse text fields
	emp.Nama = c.PostForm("nama")
	emp.JenisJabatan = services.NormalizeJenisJabatan(c.PostForm("jenis_jabatan"))
	emp.Jabatan = c.PostForm("jabatan")
	emp.TempatLahir = c.PostForm("tempat_lahir")
	emp.TanggalLahir = parseDateString(c.PostForm("tanggal_lahir"))
	emp.TempatTugas = c.PostForm("tempat_tugas")
	emp.JenisTempat = c.PostForm("jenis_tempat")
	emp.Pengangkatan = c.PostForm("pengangkatan")
	emp.JenisPengangkatan = c.PostForm("jenis_pengangkatan")

	// Parse tanggal KGB & Pangkat — CRITICAL: harus diparse dari form
	tanggalKgbStr := c.PostForm("tanggal_kgb_terakhir")
	if tanggalKgbStr != "" {
		emp.TanggalKgbTerakhir = parseDateString(tanggalKgbStr)
	}

	tanggalPangkatStr := c.PostForm("tanggal_kenaikan_pangkat_terakhir")
	if tanggalPangkatStr != "" {
		emp.TanggalKenaikanPangkatTerakhir = parseDateString(tanggalPangkatStr)
	}

	// Update phone inside user account
	noHP := c.PostForm("no_hp")
	if noHP != "" {
		db.Model(&models.User{}).Where("nip = ?", nip).Update("no_hp", noHP)
	}

	// Handle uploads (if provided, overwrite existing)
	if foto := handleFileUpload(c, "foto_profil"); foto != "" {
		emp.FotoProfil = foto
	}
	if skCpns := handleFileUpload(c, "sk_cpns_pppk_file"); skCpns != "" {
		emp.SkCpnsPppkFile = skCpns
	}
	if skPns := handleFileUpload(c, "sk_pns_file"); skPns != "" {
		emp.SkPnsFile = skPns
	}
	if skKgb := handleFileUpload(c, "sk_kgb_file"); skKgb != "" {
		emp.SkKgbFile = skKgb
	}
	if skPangkat := handleFileUpload(c, "sk_pangkat_file"); skPangkat != "" {
		emp.SkPangkatFile = skPangkat
	}
	if skPensiun := handleFileUpload(c, "sk_pensiun_file"); skPensiun != "" {
		emp.SkPensiunFile = skPensiun
	}
	if docPemberhentian := handleFileUpload(c, "dokumen_pemberhentian_pembayaran"); docPemberhentian != "" {
		emp.DokumenPemberhentianPembayaran = docPemberhentian
	}

	// Recalculate KGB Berikutnya, Pangkat Berikutnya, Pensiun berdasarkan aturan master
	_ = services.RecalculateEmployeeDates(db, &emp)

	db.Save(&emp)

	c.JSON(http.StatusOK, gin.H{"message": "Data pegawai berhasil diperbarui", "employee": emp})
}

// AddKgbHistory handles manual update/upload of new KGB (admin only)
func AddKgbHistory(c *gin.Context) {
	db := config.GetDB()
	nip := c.Param("nip")

	var emp models.Employee
	if err := db.Where("nip = ?", nip).First(&emp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pegawai tidak ditemukan"})
		return
	}

	tanggalKgbStr := c.PostForm("tanggal_kgb")
	if tanggalKgbStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal KGB wajib diisi"})
		return
	}

	skKgbFile := handleFileUpload(c, "file_sk_kgb")
	if skKgbFile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File SK KGB wajib diunggah"})
		return
	}

	tanggalKgb := parseDateString(tanggalKgbStr)

	// Create History
	history := models.EmployeeKgbHistory{
		EmployeeID: emp.ID,
		TanggalKgb: tanggalKgb,
		FileSkKgb:  skKgbFile,
	}
	db.Create(&history)

	// Update Employee Profile
	emp.TanggalKgbTerakhir = tanggalKgb
	emp.SkKgbFile = skKgbFile
	_ = services.RecalculateEmployeeDates(db, &emp)
	db.Save(&emp)

	c.JSON(http.StatusOK, gin.H{"message": "Histori KGB berhasil ditambahkan", "employee": emp})
}

// GetKgbHistory retrieves KGB history list for an employee
func GetKgbHistory(c *gin.Context) {
	db := config.GetDB()
	nip := c.Param("nip")

	var emp models.Employee
	if err := db.Where("nip = ?", nip).First(&emp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pegawai tidak ditemukan"})
		return
	}

	var history []models.EmployeeKgbHistory
	db.Where("employee_id = ?", emp.ID).Order("tanggal_kgb desc").Find(&history)

	c.JSON(http.StatusOK, gin.H{"data": history})
}

// AddPangkatHistory handles manual update/upload of new rank promotion (admin only)
func AddPangkatHistory(c *gin.Context) {
	db := config.GetDB()
	nip := c.Param("nip")

	var emp models.Employee
	if err := db.Where("nip = ?", nip).First(&emp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pegawai tidak ditemukan"})
		return
	}

	tanggalPangkatStr := c.PostForm("tanggal_kenaikan_pangkat")
	if tanggalPangkatStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal kenaikan pangkat wajib diisi"})
		return
	}

	skPangkatFile := handleFileUpload(c, "file_sk_pangkat")
	if skPangkatFile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File SK pangkat wajib diunggah"})
		return
	}

	tanggalPangkat := parseDateString(tanggalPangkatStr)

	// Create History
	history := models.EmployeePangkatHistory{
		EmployeeID:            emp.ID,
		TanggalKenaikanPangkat: tanggalPangkat,
		FileSkPangkat:         skPangkatFile,
	}
	db.Create(&history)

	// Update Employee Profile
	emp.TanggalKenaikanPangkatTerakhir = tanggalPangkat
	emp.SkPangkatFile = skPangkatFile
	_ = services.RecalculateEmployeeDates(db, &emp)
	db.Save(&emp)

	c.JSON(http.StatusOK, gin.H{"message": "Histori kenaikan pangkat berhasil ditambahkan", "employee": emp})
}

// GetPangkatHistory retrieves Rank history list for an employee
func GetPangkatHistory(c *gin.Context) {
	db := config.GetDB()
	nip := c.Param("nip")

	var emp models.Employee
	if err := db.Where("nip = ?", nip).First(&emp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pegawai tidak ditemukan"})
		return
	}

	var history []models.EmployeePangkatHistory
	db.Where("employee_id = ?", emp.ID).Order("tanggal_kenaikan_pangkat desc").Find(&history)

	c.JSON(http.StatusOK, gin.H{"data": history})
}

// Internal helper to process file upload
func handleFileUpload(c *gin.Context, fieldName string) string {
	file, header, err := c.Request.FormFile(fieldName)
	if err != nil {
		return "" // File not provided
	}
	defer file.Close()

	// Validate size (max 5MB)
	if header.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file " + fieldName + " melebihi batas 5MB"})
		c.Abort()
		return ""
	}

	// Validate extension
	ext := header.Filename
	allowed := false
	exts := []string{".pdf", ".jpg", ".jpeg", ".png"}
	for _, e := range exts {
		if hasSuffixCaseInsensitive(ext, e) {
			allowed = true
			break
		}
	}

	if !allowed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format file tidak didukung untuk " + fieldName + ". Harus berupa PDF, JPG, atau PNG."})
		c.Abort()
		return ""
	}

	// Upload using Storage provider
	url, err := storage.CurrentProvider.UploadFile(header.Filename, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan berkas: " + err.Error()})
		c.Abort()
		return ""
	}

	return url
}

func hasSuffixCaseInsensitive(str, suffix string) bool {
	return len(str) >= len(suffix) && stringsEqualCaseInsensitive(str[len(str)-len(suffix):], suffix)
}

func stringsEqualCaseInsensitive(a, b string) bool {
	return len(a) == len(b) && (a == b || (toLower(a) == toLower(b)))
}

func toLower(s string) string {
	// A simple mapping for basic extensions
	var b []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
		}
		b = append(b, c)
	}
	return string(b)
}

// DeleteAllEmployees deletes all employees and their associated accounts, histories, requests, etc. (admin only)
func DeleteAllEmployees(c *gin.Context) {
	db := config.GetDB()

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Delete EmployeeKgbHistory
	if err := tx.Exec("DELETE FROM employee_kgb_histories").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus riwayat KGB: " + err.Error()})
		return
	}

	// 2. Delete EmployeePangkatHistory
	if err := tx.Exec("DELETE FROM employee_pangkat_histories").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus riwayat Pangkat: " + err.Error()})
		return
	}

	// 3. Delete LeaveAttachments and LeaveLetters
	if err := tx.Exec("DELETE FROM leave_attachments").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus lampiran cuti: " + err.Error()})
		return
	}
	if err := tx.Exec("DELETE FROM leave_letters").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus surat cuti: " + err.Error()})
		return
	}

	// 4. Delete LeaveRequests
	if err := tx.Exec("DELETE FROM leave_requests").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pengajuan cuti: " + err.Error()})
		return
	}

	// 5. Delete DataChangeRequests
	if err := tx.Exec("DELETE FROM data_change_requests").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pengajuan ubah data: " + err.Error()})
		return
	}

	// 6. Delete BeritaAcaras
	if err := tx.Exec("DELETE FROM berita_acaras").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus berita acara: " + err.Error()})
		return
	}

	// 7. Delete User accounts where role is employee
	if err := tx.Exec("DELETE FROM users WHERE role = 'employee'").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus akun pegawai: " + err.Error()})
		return
	}

	// 8. Delete all Employees (active and retired)
	if err := tx.Exec("DELETE FROM employees").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data pegawai: " + err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal melakukan commit penghapusan: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Semua data pegawai dan akun terkait berhasil dihapus"})
}

package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"

	"sipecut/config"
	"sipecut/models"
	"sipecut/services"
)

// excelHeaders defines column mapping for import/export
var excelHeaders = []string{
	"NIP", "Nama Lengkap", "Jenis Jabatan", "Jabatan",
	"Tempat Lahir", "Tanggal Lahir (YYYY-MM-DD)",
	"Tempat Tugas", "Jenis Tempat (Dinas/Sekolah)",
	"Pengangkatan", "Jenis Kepegawaian (PNS/PPPK)",
	"No WhatsApp", "Tanggal KGB Terakhir (YYYY-MM-DD)",
	"Tanggal Pangkat Terakhir (YYYY-MM-DD)",
}

// DownloadImportTemplate generates and returns an empty Excel template for bulk import
func DownloadImportTemplate(c *gin.Context) {
	f := excelize.NewFile()
	sheet := "Data Pegawai"
	f.SetSheetName("Sheet1", sheet)

	// Set header row with style
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4F46E5"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "bottom", Color: "CCCCCC", Style: 1},
		},
	})

	for i, h := range excelHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, style)
		f.SetColWidth(sheet, string(rune('A'+i)), string(rune('A'+i)), 25)
	}

	// Add sample row
	sampleRow := []interface{}{
		"19900101001", "Ahmad Setiawan, S.Kom", "Fungsional", "Pranata Komputer Ahli Pertama",
		"Jakarta", "1990-01-01",
		"Dinas Pendidikan", "Dinas",
		"CPNS 2018", "PNS",
		"081234567890", "2024-01-01",
		"2022-01-01",
	}
	for i, v := range sampleRow {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheet, cell, v)
	}

	// Add notes sheet
	f.NewSheet("Petunjuk")
	notes := []string{
		"Petunjuk Pengisian:",
		"1. Isi data pegawai mulai dari baris ke-2 (baris pertama adalah header)",
		"2. Kolom NIP wajib diisi dan harus unik",
		"3. Jenis Jabatan: isi dengan 'Fungsional', 'Pelaksana', atau 'Struktural'",
		"4. Jenis Tempat: isi dengan 'Dinas' atau 'Sekolah'",
		"5. Jenis Kepegawaian: isi dengan 'PNS' atau 'PPPK'",
		"6. Format tanggal: YYYY-MM-DD (contoh: 1990-01-15)",
		"7. Password default akun: nip123 (dapat diubah setelah login)",
	}
	for i, note := range notes {
		f.SetCellValue("Petunjuk", fmt.Sprintf("A%d", i+1), note)
	}
	f.SetColWidth("Petunjuk", "A", "A", 80)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=template_import_pegawai.xlsx")
	c.Header("Cache-Control", "no-cache")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat file template"})
	}
}

// ExportEmployees generates and returns an Excel file of all active employees
func ExportEmployees(c *gin.Context) {
	db := config.GetDB()
	var employees []models.Employee
	if err := db.Where("status_kepegawaian = ?", "Aktif").Order("nama").Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pegawai"})
		return
	}

	f := excelize.NewFile()
	sheet := "Data Pegawai Aktif"
	f.SetSheetName("Sheet1", sheet)

	// Header
	headers := []string{
		"No", "NIP", "Nama Lengkap", "Jenis Jabatan", "Jabatan",
		"Tempat Lahir", "Tanggal Lahir", "Umur",
		"Tempat Tugas", "Jenis Tempat",
		"Pengangkatan", "Jenis Kepegawaian",
		"Tgl KGB Terakhir", "Tgl KGB Berikutnya",
		"Tgl Pangkat Terakhir", "Tgl Pangkat Berikutnya",
		"Tanggal Pensiun", "Status",
	}
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4F46E5"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// Zebra coloring
	evenStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"F1F5F9"}, Pattern: 1},
	})

	fmtDate := func(t time.Time) string {
		if t.IsZero() || t.Year() <= 1 {
			return "-"
		}
		return t.Format("02/01/2006")
	}

	for i, emp := range employees {
		row := i + 2
		years, months := services.CalculateAge(emp.TanggalLahir)
		data := []interface{}{
			i + 1, emp.NIP, emp.Nama, emp.JenisJabatan, emp.Jabatan,
			emp.TempatLahir, fmtDate(emp.TanggalLahir), fmt.Sprintf("%d tahun %d bulan", years, months),
			emp.TempatTugas, emp.JenisTempat,
			emp.Pengangkatan, emp.JenisPengangkatan,
			fmtDate(emp.TanggalKgbTerakhir), fmtDate(emp.TanggalKgbBerikutnya),
			fmtDate(emp.TanggalKenaikanPangkatTerakhir), fmtDate(emp.TanggalKenaikanPangkatBerikutnya),
			fmtDate(emp.TanggalPensiun), emp.StatusKepegawaian,
		}
		for j, v := range data {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheet, cell, v)
			if i%2 == 1 {
				f.SetCellStyle(sheet, cell, cell, evenStyle)
			}
		}
	}

	// Auto-width approximation
	colWidths := []float64{5, 18, 30, 14, 30, 15, 13, 18, 25, 12, 15, 16, 14, 14, 16, 16, 14, 10}
	for i, w := range colWidths {
		col := string(rune('A' + i))
		f.SetColWidth(sheet, col, col, w)
	}

	filename := fmt.Sprintf("data_pegawai_aktif_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Cache-Control", "no-cache")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat file export"})
	}
}

// ImportEmployees bulk-imports employees from an uploaded Excel file
func ImportEmployees(c *gin.Context) {
	db := config.GetDB()

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File Excel wajib diunggah (field: file)"})
		return
	}

	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".xlsx") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format file harus .xlsx"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka file"})
		return
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca file Excel: " + err.Error()})
		return
	}

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil || len(rows) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File Excel kosong atau format tidak valid"})
		return
	}

	type importResult struct {
		Row     int    `json:"row"`
		NIP     string `json:"nip"`
		Nama    string `json:"nama"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	results := []importResult{}
	successCount := 0
	failCount := 0

	// Pre-load reference rules to avoid thousands of DB queries
	var pensionRules []models.PensionRule
	db.Find(&pensionRules)

	var kgbRules []models.KgbCycleRule
	db.Find(&kgbRules)

	var pangkatRules []models.PangkatCycleRule
	db.Find(&pangkatRules)

	var setting models.AppSetting
	criteria := "dokumen_pemberhentian"
	if err := db.Where("key = ?", "kriteria_pensiun_lengkap").First(&setting).Error; err == nil {
		criteria = setting.Value
	}

	// Pre-load existing employees to check NIP existence in memory
	var existingEmployees []models.Employee
	if err := db.Find(&existingEmployees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data pegawai eksisting"})
		return
	}
	existingMap := make(map[string]models.Employee)
	for _, e := range existingEmployees {
		existingMap[e.NIP] = e
	}

	findKgbRule := func(jenisJabatan, jabatan string) int {
		normJenis := services.NormalizeJenisJabatan(jenisJabatan)
		for _, r := range kgbRules {
			if services.NormalizeJenisJabatan(r.JenisJabatan) == normJenis && r.Jabatan == jabatan {
				return r.SiklusTahun
			}
		}
		for _, r := range kgbRules {
			if services.NormalizeJenisJabatan(r.JenisJabatan) == normJenis && r.Jabatan == "*" {
				return r.SiklusTahun
			}
		}
		return 2
	}

	findPangkatRule := func(jenisJabatan, jabatan string) int {
		normJenis := services.NormalizeJenisJabatan(jenisJabatan)
		for _, r := range pangkatRules {
			if services.NormalizeJenisJabatan(r.JenisJabatan) == normJenis && r.Jabatan == jabatan {
				return r.SiklusTahun
			}
		}
		for _, r := range pangkatRules {
			if services.NormalizeJenisJabatan(r.JenisJabatan) == normJenis && r.Jabatan == "*" {
				return r.SiklusTahun
			}
		}
		return 4
	}

	findPensionRule := func(jenisJabatan, jabatan string) int {
		normJenis := services.NormalizeJenisJabatan(jenisJabatan)
		for _, r := range pensionRules {
			if services.NormalizeJenisJabatan(r.JenisJabatan) == normJenis && r.Jabatan == jabatan {
				return r.BatasUsiaPensiun
			}
		}
		for _, r := range pensionRules {
			if services.NormalizeJenisJabatan(r.JenisJabatan) == normJenis && r.Jabatan == "*" {
				return r.BatasUsiaPensiun
			}
		}
		return 58
	}

	parseDate := func(s string) time.Time {
		s = strings.TrimSpace(s)
		if s == "" || s == "-" {
			return time.Time{}
		}

		// 1. Try parsing as numeric Excel date
		if floatVal, err := strconv.ParseFloat(s, 64); err == nil {
			if floatVal > 1000 {
				if t, err := excelize.ExcelDateToTime(floatVal, false); err == nil {
					return t
				}
			}
		}

		// 2. Try parsing standard date string formats
		layouts := []string{
			"2006-01-02",
			"02/01/2006",
			"01/02/2006",
			"02-01-2006",
			"2006/01/02",
			"02/01/06",
			"02-01-06",
			"2006-1-2",
			"2-1-2006",
			"2/1/2006",
		}
		for _, layout := range layouts {
			if t, err := time.Parse(layout, s); err == nil {
				return t
			}
		}
		return time.Time{}
	}

	getCell := func(row []string, idx int) string {
		if idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	// Skip header row (index 0), process from row 1
	for rowIdx, row := range rows[1:] {
		actualRow := rowIdx + 2
		nip := getCell(row, 0)
		nama := getCell(row, 1)

		if nip == "" || nama == "" {
			results = append(results, importResult{Row: actualRow, NIP: nip, Nama: nama, Status: "Skip", Message: "NIP atau Nama kosong, baris dilewati"})
			continue
		}

		// Check memory map instead of calling DB
		existing, nipExists := existingMap[nip]

		jenisJabatan := services.NormalizeJenisJabatan(getCell(row, 2))
		if jenisJabatan == "" {
			jenisJabatan = "Fungsional"
		}
		jabatan := getCell(row, 3)
		tempatLahir := getCell(row, 4)
		tanggalLahir := parseDate(getCell(row, 5))
		tempatTugas := getCell(row, 6)
		jenisTempat := getCell(row, 7)
		if jenisTempat == "" {
			jenisTempat = "Dinas"
		}
		pengangkatan := getCell(row, 8)
		jenisPengangkatan := getCell(row, 9)
		if jenisPengangkatan == "" {
			jenisPengangkatan = "PNS"
		}
		noHP := getCell(row, 10)
		tglKgb := parseDate(getCell(row, 11))
		tglPangkat := parseDate(getCell(row, 12))

		emp := models.Employee{
			NIP:                            nip,
			Nama:                           nama,
			JenisJabatan:                   jenisJabatan,
			Jabatan:                        jabatan,
			TempatLahir:                    tempatLahir,
			TanggalLahir:                   tanggalLahir,
			TempatTugas:                    tempatTugas,
			JenisTempat:                    jenisTempat,
			Pengangkatan:                   pengangkatan,
			JenisPengangkatan:              jenisPengangkatan,
			TanggalKgbTerakhir:             tglKgb,
			TanggalKenaikanPangkatTerakhir: tglPangkat,
			StatusKepegawaian:              "Aktif",
		}

		// Recalculate dates in-memory to avoid DB queries
		zeroTime := time.Time{}
		if emp.TanggalKgbTerakhir.IsZero() || emp.TanggalKgbTerakhir.Year() <= 1 {
			emp.TanggalKgbBerikutnya = zeroTime
		} else {
			emp.TanggalKgbBerikutnya = emp.TanggalKgbTerakhir.AddDate(findKgbRule(emp.JenisJabatan, emp.Jabatan), 0, 0)
		}

		if emp.TanggalKenaikanPangkatTerakhir.IsZero() || emp.TanggalKenaikanPangkatTerakhir.Year() <= 1 {
			emp.TanggalKenaikanPangkatBerikutnya = zeroTime
		} else {
			emp.TanggalKenaikanPangkatBerikutnya = emp.TanggalKenaikanPangkatTerakhir.AddDate(findPangkatRule(emp.JenisJabatan, emp.Jabatan), 0, 0)
		}

		if !emp.TanggalLahir.IsZero() && emp.TanggalLahir.Year() > 1 {
			emp.TanggalPensiun = emp.TanggalLahir.AddDate(findPensionRule(emp.JenisJabatan, emp.Jabatan), 0, 0)
		}

		hasPassedPensionDate := !emp.TanggalPensiun.IsZero() && time.Now().After(emp.TanggalPensiun)
		if hasPassedPensionDate {
			if criteria == "tanggal_saja" {
				emp.StatusKepegawaian = "Pensiun"
			} else if criteria == "dokumen_pemberhentian" {
				if emp.DokumenPemberhentianPembayaran != "" {
					emp.StatusKepegawaian = "Pensiun"
				} else {
					emp.StatusKepegawaian = "Aktif"
				}
			}
		} else {
			emp.StatusKepegawaian = "Aktif"
		}

		tx := db.Begin()
		if nipExists {
			// Update existing
			if err := tx.Model(&existing).Updates(&emp).Error; err != nil {
				tx.Rollback()
				failCount++
				results = append(results, importResult{Row: actualRow, NIP: nip, Nama: nama, Status: "Gagal", Message: "Update gagal: " + err.Error()})
				continue
			}
			emp.ID = existing.ID
		} else {
			// Create new employee
			if err := tx.Create(&emp).Error; err != nil {
				tx.Rollback()
				failCount++
				results = append(results, importResult{Row: actualRow, NIP: nip, Nama: nama, Status: "Gagal", Message: "Simpan gagal: " + err.Error()})
				continue
			}
			// Create user account with default password = nip + "123"
			// Use bcrypt.MinCost to speed up bulk imports significantly
			defaultPass := nip + "123"
			hash, _ := bcrypt.GenerateFromPassword([]byte(defaultPass), bcrypt.MinCost)
			user := models.User{
				NIP:          nip,
				PasswordHash: string(hash),
				Role:         "employee",
				NoHP:         noHP,
				Status:       "aktif",
			}
			tx.Create(&user)
		}
		tx.Commit()

		// Update our in-memory map to reflect the new state so subsequent duplicates update instead of insert
		existingMap[nip] = emp

		successCount++
		action := "Dibuat"
		if nipExists {
			action = "Diperbarui"
		}
		results = append(results, importResult{Row: actualRow, NIP: nip, Nama: nama, Status: action, Message: "Berhasil"})
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       fmt.Sprintf("Import selesai: %d berhasil, %d gagal dari %d baris data", successCount, failCount, len(rows)-1),
		"success_count": successCount,
		"fail_count":    failCount,
		"details":       results,
	})
}

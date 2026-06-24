package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"sipecut/config"
	"sipecut/models"
	"sipecut/services"
)

// ---- MASTER JABATAN ----

// GetMasterJabatan returns list of all jabatan, optionally filtered by jenis_jabatan
func GetMasterJabatan(c *gin.Context) {
	db := config.GetDB()
	var jabatans []models.MasterJabatan

	q := db.Order("jenis_jabatan, nama_jabatan")
	if jj := c.Query("jenis_jabatan"); jj != "" {
		q = q.Where("jenis_jabatan = ?", jj)
	}

	if err := q.Find(&jabatans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data jabatan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": jabatans})
}

// CreateMasterJabatan adds a new job title to the master list
func CreateMasterJabatan(c *gin.Context) {
	var input models.MasterJabatan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.NamaJabatan == "" || input.JenisJabatan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama jabatan dan jenis jabatan wajib diisi"})
		return
	}

	db := config.GetDB()
	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan jabatan"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": input, "message": "Jabatan berhasil ditambahkan"})
}

// UpdateMasterJabatan updates an existing jabatan record
func UpdateMasterJabatan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	db := config.GetDB()
	var jabatan models.MasterJabatan
	if err := db.First(&jabatan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data jabatan tidak ditemukan"})
		return
	}
	var input models.MasterJabatan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	jabatan.NamaJabatan = input.NamaJabatan
	jabatan.JenisJabatan = input.JenisJabatan
	db.Save(&jabatan)
	c.JSON(http.StatusOK, gin.H{"data": jabatan, "message": "Jabatan berhasil diperbarui"})
}

// DeleteMasterJabatan removes a jabatan from the master list
func DeleteMasterJabatan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	db := config.GetDB()
	if err := db.Delete(&models.MasterJabatan{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus jabatan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Jabatan berhasil dihapus"})
}

// ---- MASTER TEMPAT TUGAS ----

// GetMasterTempatTugas returns list of all tempat tugas, optionally filtered by jenis_tempat
func GetMasterTempatTugas(c *gin.Context) {
	db := config.GetDB()
	var places []models.MasterTempatTugas

	q := db.Order("jenis_tempat, nama_tempat")
	if jt := c.Query("jenis_tempat"); jt != "" {
		q = q.Where("jenis_tempat = ?", jt)
	}

	if err := q.Find(&places).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data tempat tugas"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": places})
}

// CreateMasterTempatTugas adds a new work unit to the master list
func CreateMasterTempatTugas(c *gin.Context) {
	var input models.MasterTempatTugas
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.NamaTempat == "" || input.JenisTempat == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama tempat dan jenis tempat wajib diisi"})
		return
	}

	db := config.GetDB()

	// Check if db is nil (shouldn't happen but guard anyway)
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Koneksi database tidak tersedia"})
		return
	}

	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan tempat tugas"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": input, "message": "Tempat tugas berhasil ditambahkan"})
}

// UpdateMasterTempatTugas updates an existing tempat tugas record
func UpdateMasterTempatTugas(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	db := config.GetDB()
	var place models.MasterTempatTugas
	if err := db.First(&place, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tempat tugas tidak ditemukan"})
		return
	}
	var input models.MasterTempatTugas
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	place.NamaTempat = input.NamaTempat
	place.JenisTempat = input.JenisTempat
	db.Save(&place)
	c.JSON(http.StatusOK, gin.H{"data": place, "message": "Tempat tugas berhasil diperbarui"})
}

// DeleteMasterTempatTugas removes a tempat tugas from the master list
func DeleteMasterTempatTugas(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	db := config.GetDB()
	if err := db.Delete(&models.MasterTempatTugas{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus tempat tugas"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tempat tugas berhasil dihapus"})
}

// helper to avoid unused import warning
var _ = (*gorm.DB)(nil)

// DownloadMasterTemplate generates and returns an Excel template for master data import
func DownloadMasterTemplate(c *gin.Context) {
	f := excelize.NewFile()

	// Sheet 1: Master Pangkat (Aturan Pangkat)
	sheetPangkat := "Master Pangkat"
	f.SetSheetName("Sheet1", sheetPangkat)

	// Sheet 2: Master Jabatan
	sheetJabatan := "Master Jabatan"
	f.NewSheet(sheetJabatan)

	// Sheet 3: Master Tempat Tugas
	sheetTempat := "Master Tempat Tugas"
	f.NewSheet(sheetTempat)

	// Header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"1A365D"}, Pattern: 1}, // Navy style matching theme
		Alignment: &excelize.Alignment{Horizontal: "center", WrapText: true},
	})

	// 1. Setup Master Pangkat
	pangkatHeaders := []string{"Jenis Jabatan (Fungsional/Pelaksana/Struktural)", "Nama Jabatan (* untuk Semua)", "Siklus Kenaikan (Tahun)"}
	for i, h := range pangkatHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetPangkat, cell, h)
		f.SetCellStyle(sheetPangkat, cell, cell, headerStyle)
	}
	f.SetCellValue(sheetPangkat, "A2", "Fungsional")
	f.SetCellValue(sheetPangkat, "B2", "*")
	f.SetCellValue(sheetPangkat, "C2", 4)
	f.SetColWidth(sheetPangkat, "A", "C", 30)

	// 2. Setup Master Jabatan
	jabatanHeaders := []string{"Nama Jabatan", "Jenis Jabatan (Fungsional/Pelaksana/Struktural)"}
	for i, h := range jabatanHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetJabatan, cell, h)
		f.SetCellStyle(sheetJabatan, cell, cell, headerStyle)
	}
	f.SetCellValue(sheetJabatan, "A2", "Guru Ahli Pertama")
	f.SetCellValue(sheetJabatan, "B2", "Fungsional")
	f.SetColWidth(sheetJabatan, "A", "B", 30)

	// 3. Setup Master Tempat Tugas
	tempatHeaders := []string{"Nama Tempat Tugas", "Jenis Tempat (Dinas/Sekolah)"}
	for i, h := range tempatHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetTempat, cell, h)
		f.SetCellStyle(sheetTempat, cell, cell, headerStyle)
	}
	f.SetCellValue(sheetTempat, "A2", "SDN 01 Lembo")
	f.SetCellValue(sheetTempat, "B2", "Sekolah")
	f.SetColWidth(sheetTempat, "A", "B", 30)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=template_import_master.xlsx")
	c.Header("Cache-Control", "no-cache")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat file template master"})
	}
}

// ImportMasterData bulk imports master pangkat, jabatan, and tempat tugas from Excel
func ImportMasterData(c *gin.Context) {
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

	sheets := f.GetSheetList()
	
	pangkatImported := 0
	jabatanImported := 0
	tempatImported := 0
	
	getCell := func(row []string, idx int) string {
		if idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	// 1. Process "Master Pangkat" sheet
	sheetPangkat := getActualSheetName(sheets, "Master Pangkat")
	if sheetPangkat != "" {
		rows, err := f.GetRows(sheetPangkat)
		if err != nil {
			log.Printf("[Import] Error reading Master Pangkat rows: %v", err)
		} else if len(rows) > 1 {
			for _, row := range rows[1:] {
				jenisJabatan := services.NormalizeJenisJabatan(getCell(row, 0))
				jabatan := getCell(row, 1)
				siklusStr := getCell(row, 2)
				
				if jenisJabatan == "" || jabatan == "" || siklusStr == "" {
					continue
				}
				
				siklus, err := strconv.Atoi(siklusStr)
				if err != nil || siklus <= 0 {
					continue
				}

				// Find existing Pangkat rule
				var rule models.PangkatCycleRule
				errFind := db.Where("jenis_jabatan = ? AND jabatan = ?", jenisJabatan, jabatan).First(&rule).Error
				if errFind == nil {
					rule.SiklusTahun = siklus
					if errSave := db.Save(&rule).Error; errSave != nil {
						log.Printf("[Import] Error saving Pangkat cycle rule: %v", errSave)
					}
				} else {
					if errCreate := db.Create(&models.PangkatCycleRule{
						JenisJabatan: jenisJabatan,
						Jabatan:      jabatan,
						SiklusTahun:  siklus,
					}).Error; errCreate != nil {
						log.Printf("[Import] Error creating Pangkat cycle rule: %v", errCreate)
					}
				}
				pangkatImported++
			}
		}
	}

	// 2. Process "Master Jabatan" sheet
	sheetJabatan := getActualSheetName(sheets, "Master Jabatan")
	if sheetJabatan != "" {
		rows, err := f.GetRows(sheetJabatan)
		if err != nil {
			log.Printf("[Import] Error reading Master Jabatan rows: %v", err)
		} else if len(rows) > 1 {
			for _, row := range rows[1:] {
				namaJabatan := getCell(row, 0)
				jenisJabatanInput := getCell(row, 1)
				
				if namaJabatan == "" {
					continue
				}
				
				jenisJabatan := services.NormalizeJenisJabatan(jenisJabatanInput)
				if jenisJabatan == "" {
					// Smart autodetection based on keywords
					lowerNama := strings.ToLower(namaJabatan)
					if strings.Contains(lowerNama, "guru") || 
						strings.Contains(lowerNama, "kepala sekolah") || 
						strings.Contains(lowerNama, "pranata") || 
						strings.Contains(lowerNama, "analis") || 
						strings.Contains(lowerNama, "dokter") || 
						strings.Contains(lowerNama, "perawat") || 
						strings.Contains(lowerNama, "bidan") || 
						strings.Contains(lowerNama, "penyuluh") || 
						strings.Contains(lowerNama, "pengawas") {
						jenisJabatan = "Fungsional"
					} else {
						jenisJabatan = "Pelaksana"
					}
				}

				var count int64
				db.Model(&models.MasterJabatan{}).Where("nama_jabatan = ? AND jenis_jabatan = ?", namaJabatan, jenisJabatan).Count(&count)
				if count == 0 {
					if errCreate := db.Create(&models.MasterJabatan{
						NamaJabatan:  namaJabatan,
						JenisJabatan: jenisJabatan,
					}).Error; errCreate != nil {
						log.Printf("[Import] Error creating Master Jabatan: %v", errCreate)
					}
				}
				jabatanImported++
			}
		}
	}

	// 3. Process "Master Tempat Tugas" sheet
	sheetTempat := getActualSheetName(sheets, "Master Tempat Tugas")
	if sheetTempat != "" {
		rows, err := f.GetRows(sheetTempat)
		if err != nil {
			log.Printf("[Import] Error reading Master Tempat Tugas rows: %v", err)
		} else if len(rows) > 1 {
			for _, row := range rows[1:] {
				namaTempat := getCell(row, 0)
				jenisTempatInput := getCell(row, 1)
				
				if namaTempat == "" {
					continue
				}
				
				jenisTempat := strings.TrimSpace(jenisTempatInput)
				if jenisTempat == "" {
					// Smart autodetection based on keywords
					lowerTempat := strings.ToLower(namaTempat)
					if strings.Contains(lowerTempat, "sd") || 
						strings.Contains(lowerTempat, "smp") || 
						strings.Contains(lowerTempat, "sma") || 
						strings.Contains(lowerTempat, "smk") || 
						strings.Contains(lowerTempat, "sekolah") || 
						strings.Contains(lowerTempat, "tk ") || 
						strings.Contains(lowerTempat, "paud") {
						jenisTempat = "Sekolah"
					} else {
						jenisTempat = "Dinas"
					}
				} else {
					// Normalize first letter capitalization
					lower := strings.ToLower(jenisTempat)
					if strings.Contains(lower, "sekolah") {
						jenisTempat = "Sekolah"
					} else if strings.Contains(lower, "dinas") {
						jenisTempat = "Dinas"
					} else {
						// Custom title case
						if len(lower) > 0 {
							jenisTempat = strings.Title(lower)
						}
					}
				}

				var count int64
				db.Model(&models.MasterTempatTugas{}).Where("nama_tempat = ? AND jenis_tempat = ?", namaTempat, jenisTempat).Count(&count)
				if count == 0 {
					if errCreate := db.Create(&models.MasterTempatTugas{
						NamaTempat:  namaTempat,
						JenisTempat: jenisTempat,
					}).Error; errCreate != nil {
						log.Printf("[Import] Error creating Master Tempat Tugas: %v", errCreate)
					}
				}
				tempatImported++
			}
		}
	}

	// Recalculate employee dates in case master rules changed
	if pangkatImported > 0 {
		_ = services.RecalculateAllEmployees(db)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Berhasil mengimpor master: %d aturan pangkat, %d jabatan, %d tempat tugas.", pangkatImported, jabatanImported, tempatImported),
		"pangkat_count": pangkatImported,
		"jabatan_count": jabatanImported,
		"tempat_count": tempatImported,
	})
}

func getActualSheetName(sheets []string, target string) string {
	for _, s := range sheets {
		if strings.EqualFold(strings.TrimSpace(s), strings.TrimSpace(target)) {
			return s
		}
	}
	return ""
}

func hasSheet(sheets []string, target string) bool {
	return getActualSheetName(sheets, target) != ""
}

// DeleteAllMasterJabatan deletes all master jabatan items (admin only)
func DeleteAllMasterJabatan(c *gin.Context) {
	db := config.GetDB()
	if err := db.Exec("DELETE FROM master_jabatans").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus semua data jabatan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Semua data master jabatan berhasil dihapus"})
}

// DeleteAllMasterTempatTugas deletes all master tempat tugas items (admin only)
func DeleteAllMasterTempatTugas(c *gin.Context) {
	db := config.GetDB()
	if err := db.Exec("DELETE FROM master_tempat_tugas").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus semua data tempat tugas"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Semua data master tempat tugas berhasil dihapus"})
}

package services

import (
<<<<<<< HEAD
	"os"
=======
>>>>>>> 603353f54c6625439da1b7cf09eb935c784c51b4
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"sipecut/models"
)

func TestCalculateAge(t *testing.T) {
	// Mock current time relative testing
	birth := time.Now().AddDate(-35, -5, 0) // 35 years and 5 months ago
	years, months := CalculateAge(birth)

	if years != 35 {
		t.Errorf("Expected years to be 35, got %d", years)
	}
	if months != 5 {
		t.Errorf("Expected months to be 5, got %d", months)
	}

	birth2 := time.Now().AddDate(-10, 2, 0) // 10 years minus 2 months = 9 years and 10 months ago
	years2, months2 := CalculateAge(birth2)

	if years2 != 9 {
		t.Errorf("Expected years to be 9, got %d", years2)
	}
	if months2 != 10 {
		t.Errorf("Expected months to be 10, got %d", months2)
	}
}

func TestRulesCalculations(t *testing.T) {
	// Initialize In-Memory SQLite DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open in-memory sqlite: %v", err)
	}

	// Auto-Migrate
	err = db.AutoMigrate(&models.PensionRule{}, &models.KgbCycleRule{}, &models.PangkatCycleRule{}, &models.Employee{}, &models.AppSetting{})
	if err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	// Seed Master Data
	db.Create(&models.KgbCycleRule{JenisJabatan: "Fungsional", Jabatan: "*", SiklusTahun: 2})
	db.Create(&models.KgbCycleRule{JenisJabatan: "Pelaksana", Jabatan: "*", SiklusTahun: 3}) // specific test cycle
	db.Create(&models.PangkatCycleRule{JenisJabatan: "Fungsional", Jabatan: "*", SiklusTahun: 4})
	db.Create(&models.PensionRule{JenisJabatan: "Fungsional", Jabatan: "*", BatasUsiaPensiun: 60})

	// Test Dates Calculations
	lastKgb := time.Date(2025, 6, 17, 0, 0, 0, 0, time.UTC)
	nextKgbFungsional := CalculateNextKgbDate(db, "Fungsional", "Guru", lastKgb)
	if nextKgbFungsional.Year() != 2027 {
		t.Errorf("Expected next KGB year for Fungsional to be 2027, got %d", nextKgbFungsional.Year())
	}

	nextKgbPelaksana := CalculateNextKgbDate(db, "Pelaksana", "Staf", lastKgb)
	if nextKgbPelaksana.Year() != 2028 { // 3 years cycle
		t.Errorf("Expected next KGB year for Pelaksana to be 2028, got %d", nextKgbPelaksana.Year())
	}

	birthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	pensionDate := CalculatePensionDate(db, "Fungsional", "Guru", birthDate)
	if pensionDate.Year() != 2040 { // 1980 + 60 = 2040
		t.Errorf("Expected pension year for Fungsional to be 2040, got %d", pensionDate.Year())
	}
}

func TestRecalculateEmployeeDates(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open db: %v", err)
	}
	db.AutoMigrate(&models.PensionRule{}, &models.KgbCycleRule{}, &models.PangkatCycleRule{}, &models.Employee{}, &models.AppSetting{})

	// Seed criteria
	db.Create(&models.AppSetting{Key: "kriteria_pensiun_lengkap", Value: "dokumen_pemberhentian"})
	db.Create(&models.PensionRule{JenisJabatan: "Fungsional", Jabatan: "*", BatasUsiaPensiun: 60})
	db.Create(&models.KgbCycleRule{JenisJabatan: "Fungsional", Jabatan: "*", SiklusTahun: 2})
	db.Create(&models.PangkatCycleRule{JenisJabatan: "Fungsional", Jabatan: "*", SiklusTahun: 4})

	// Case 1: Pension date has passed but NO document uploaded (Criteria: dokumen_pemberhentian)
	// Born 70 years ago, so pension date passed
	birth := time.Now().AddDate(-70, 0, 0)
	emp := models.Employee{
		NIP:                            "11111",
		Nama:                           "Pegawai Tua",
		JenisJabatan:                   "Fungsional",
		Jabatan:                        "Guru",
		TanggalLahir:                   birth,
		TanggalKgbTerakhir:             time.Now().AddDate(-2, 0, 0),
		TanggalKenaikanPangkatTerakhir: time.Now().AddDate(-4, 0, 0),
		StatusKepegawaian:              "Aktif",
	}

	err = RecalculateEmployeeDates(db, &emp)
	if err != nil {
		t.Errorf("Recalculate failed: %v", err)
	}

	if emp.StatusKepegawaian != "Aktif" {
		t.Errorf("Expected status to remain Aktif since document is not uploaded, got %s", emp.StatusKepegawaian)
	}

	// Case 2: Document uploaded, should transition to Pensiun
	emp.DokumenPemberhentianPembayaran = "http://storage.com/sk_pensiun.pdf"
	err = RecalculateEmployeeDates(db, &emp)
	if err != nil {
		t.Errorf("Recalculate failed: %v", err)
	}

	if emp.StatusKepegawaian != "Pensiun" {
		t.Errorf("Expected status to transition to Pensiun, got %s", emp.StatusKepegawaian)
	}
}
<<<<<<< HEAD

func TestGenerateDocx(t *testing.T) {
	req := &models.LeaveRequest{
		Employee: models.Employee{
			Nama: "MOHAMAD RIDWAN DAENG MALURENG, S.AG",
			NIP:  "199005022025212051",
		},
		JenisCuti:      "Cuti Tahunan Biasa",
		TanggalMulai:   time.Now(),
		TanggalSelesai: time.Now().AddDate(0, 0, 5),
	}
	bytes, err := GenerateLeaveDocx(req)
	if err != nil {
		t.Fatalf("GenerateLeaveDocx failed: %v", err)
	}
	err = os.WriteFile("test.docx", bytes, 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
}
=======
>>>>>>> 603353f54c6625439da1b7cf09eb935c784c51b4

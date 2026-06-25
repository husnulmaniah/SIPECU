package config

import (
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"sipecut/models"
)

var DBCon *gorm.DB

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DBCon
}

// ConnectDatabase initializes database connection
func ConnectDatabase() {
	var db *gorm.DB
	var err error

	dbURL := os.Getenv("DATABASE_URL")

	// Choose SQLite for local development, Postgres for production
	if dbURL == "" {
		log.Println("DATABASE_URL not found, using SQLite local database (sipecut.db)")
		db, err = gorm.Open(sqlite.Open("sipecut.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			if sqlDB, errDB := db.DB(); errDB == nil {
				_, _ = sqlDB.Exec("PRAGMA journal_mode=WAL;")
				_, _ = sqlDB.Exec("PRAGMA busy_timeout=30000;")
				log.Println("SQLite WAL mode and busy_timeout configured.")
			}
		}
	} else {
		log.Println("Connecting to PostgreSQL database...")
		db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DBCon = db

	// Auto Migration
	log.Println("Running database auto-migrations...")
	err = db.AutoMigrate(
		&models.Role{},
		&models.Menu{},
		&models.User{},
		&models.Employee{},
		&models.PensionRule{},
		&models.KgbCycleRule{},
		&models.PangkatCycleRule{},
		&models.EmployeeKgbHistory{},
		&models.EmployeePangkatHistory{},
		&models.LeaveRequest{},
		&models.LeaveAttachment{},
		&models.LeaveLetter{},
		&models.DataChangeRequest{},
		&models.BeritaAcara{},
		&models.RequestHistory{},
		&models.NotificationLog{},
		&models.AppSetting{},
		&models.MasterJabatan{},
		&models.MasterTempatTugas{},
	)
	if err != nil {
		log.Fatalf("Auto-migration failed: %v", err)
	}

	// Seeding initial data
	SeedData(db)
}

// SeedData inserts master rules and default accounts if they do not exist
func SeedData(db *gorm.DB) {
	// 1. Seed Rules
	var count int64

	// Roles
	var roleCount int64
	db.Model(&models.Role{}).Count(&roleCount)
	if roleCount == 0 {
		log.Println("Seeding default roles...")
		roles := []models.Role{
			{ID: 1, NamaRole: "admin"},
			{ID: 2, NamaRole: "employee"},
		}
		db.Create(&roles)
	}

	// Menus
	var menuCount int64
	db.Model(&models.Menu{}).Count(&menuCount)
	if menuCount == 0 {
		log.Println("Seeding default menus...")
		menus := []models.Menu{
			{ID: 1, NamaMenu: "Dashboard", KodeMenu: "dashboard"},
			{ID: 2, NamaMenu: "Profil Saya", KodeMenu: "self_profile"},
			{ID: 3, NamaMenu: "Pengajuan Cuti", KodeMenu: "leave_request"},
			{ID: 4, NamaMenu: "Perubahan Data", KodeMenu: "data_change"},
			{ID: 5, NamaMenu: "Berita Acara", KodeMenu: "berita_acara"},
			{ID: 6, NamaMenu: "Manajemen Pegawai", KodeMenu: "admin_employee"},
			{ID: 7, NamaMenu: "Aturan Pensiun", KodeMenu: "admin_rules"},
			{ID: 8, NamaMenu: "Riwayat Audit", KodeMenu: "admin_audit"},
			{ID: 9, NamaMenu: "Notifikasi WA", KodeMenu: "admin_wa"},
		}
		db.Create(&menus)
	}

	// Pension Rules
	db.Model(&models.PensionRule{}).Count(&count)
	if count == 0 {
		log.Println("Seeding default pension rules...")
		rules := []models.PensionRule{
			{JenisJabatan: "Fungsional", Jabatan: "*", BatasUsiaPensiun: 60}, // Standard Fungsional (e.g. Guru)
			{JenisJabatan: "Pelaksana", Jabatan: "*", BatasUsiaPensiun: 58},  // Standard Pelaksana
		}
		db.Create(&rules)
	}

	// KGB Cycle Rules (Salary increments)
	db.Model(&models.KgbCycleRule{}).Count(&count)
	if count == 0 {
		log.Println("Seeding default KGB cycle rules...")
		rules := []models.KgbCycleRule{
			{JenisJabatan: "Fungsional", Jabatan: "*", SiklusTahun: 2},
			{JenisJabatan: "Pelaksana", Jabatan: "*", SiklusTahun: 2},
		}
		db.Create(&rules)
	}

	// Pangkat Cycle Rules (Rank progression)
	db.Model(&models.PangkatCycleRule{}).Count(&count)
	if count == 0 {
		log.Println("Seeding default pangkat cycle rules...")
		rules := []models.PangkatCycleRule{
			{JenisJabatan: "Fungsional", Jabatan: "*", SiklusTahun: 4},
			{JenisJabatan: "Pelaksana", Jabatan: "*", SiklusTahun: 4},
		}
		db.Create(&rules)
	}

	// App Settings
	db.Model(&models.AppSetting{}).Count(&count)
	if count == 0 {
		log.Println("Seeding default app settings...")
		settings := []models.AppSetting{
			{Key: "kriteria_pensiun_lengkap", Value: "dokumen_pemberhentian"}, // Options: "tanggal_saja" or "dokumen_pemberhentian"
			{Key: "whatsapp_template_leave_new", Value: "Halo Admin, ada pengajuan cuti baru dari {{NAMA}} (NIP: {{NIP}}) jenis {{JENIS_CUTI}}. Silakan tinjau di dashboard."},
			{Key: "whatsapp_template_leave_update", Value: "Halo {{NAMA}}, pengajuan cuti Anda ({{JENIS_CUTI}}) telah diupdate dengan status: {{STATUS}}. Catatan: {{CATATAN}}."},
			{Key: "whatsapp_template_change_new", Value: "Halo Admin, ada pengajuan perubahan data baru dari {{NAMA}} (NIP: {{NIP}})."},
			{Key: "whatsapp_template_change_update", Value: "Halo {{NAMA}}, pengajuan perubahan data Anda telah {{STATUS}}. Catatan: {{CATATAN}}."},
			{Key: "whatsapp_template_ba_new", Value: "Halo Admin, ada dokumen Berita Acara {{JENIS}} baru diunggah oleh {{NAMA}}."},
			{Key: "whatsapp_template_ba_update", Value: "Halo {{NAMA}}, dokumen Berita Acara {{JENIS}} Anda telah {{STATUS}}."},
		}
		db.Create(&settings)
	}

	// 2. Seed Users
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		log.Println("Seeding default users (Admin & Employee sample)...")

		// Load menus
		var allMenus []models.Menu
		db.Find(&allMenus)

		var employeeMenus []models.Menu
		for _, m := range allMenus {
			if m.KodeMenu == "self_profile" || m.KodeMenu == "leave_request" || m.KodeMenu == "data_change" || m.KodeMenu == "berita_acara" {
				employeeMenus = append(employeeMenus, m)
			}
		}

		// Create Admin
		adminHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		adminUser := models.User{
			NIP:               "admin",
			PasswordHash:      string(adminHash),
			Role:              "admin",
			NoHP:              "081234567890",
			Status:            "aktif",
			RoleID:            1,
			IsPasswordDefault: false,
			Menus:             allMenus,
		}
		db.Create(&adminUser)

		// Create Employee User
		empHash, _ := bcrypt.GenerateFromPassword([]byte("pegawai123"), bcrypt.DefaultCost)
		empUser := models.User{
			NIP:               "19900101001",
			PasswordHash:      string(empHash),
			Role:              "employee",
			NoHP:              "089876543210",
			Status:            "aktif",
			RoleID:            2,
			IsPasswordDefault: false,
			Menus:             employeeMenus,
		}
		db.Create(&empUser)

		// Create Employee Profile
		birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
		kgbDate := time.Now().AddDate(-1, 0, 0) // 1 year ago
		pangkatDate := time.Now().AddDate(-2, 0, 0) // 2 years ago

		// Automatic calculation placeholders (will be recalculate on logic run)
		employee := models.Employee{
			NIP:                              "19900101001",
			Nama:                             "Ahmad Setiawan, S.Kom",
			JenisJabatan:                     "Fungsional",
			Jabatan:                          "Pranata Komputer Ahli Pertama",
			TempatLahir:                      "Jakarta",
			TanggalLahir:                     birthDate,
			TempatTugas:                      "Dinas Kominfo",
			JenisTempat:                      "Dinas",
			Pengangkatan:                     "CPNS 2018",
			TanggalKgbTerakhir:               kgbDate,
			TanggalKgbBerikutnya:             kgbDate.AddDate(2, 0, 0),
			TanggalKenaikanPangkatTerakhir:   pangkatDate,
			TanggalKenaikanPangkatBerikutnya: pangkatDate.AddDate(4, 0, 0),
			TanggalPensiun:                   birthDate.AddDate(60, 0, 0),
			StatusKepegawaian:                "Aktif",
			JenisPengangkatan:                "PNS",
			FotoProfil:                       "",
			SkCpnsPppkFile:                   "",
			SkPnsFile:                        "",
			SkKgbFile:                        "",
			SkPangkatFile:                    "",
		}
		db.Create(&employee)
		
		// Insert initial histories
		db.Create(&models.EmployeeKgbHistory{
			EmployeeID: employee.ID,
			TanggalKgb: kgbDate,
			FileSkKgb:  "",
		})
		db.Create(&models.EmployeePangkatHistory{
			EmployeeID:            employee.ID,
			TanggalKenaikanPangkat: pangkatDate,
			FileSkPangkat:         "",
		})
	}

	// 3. Seed Master Jabatan
	db.Model(&models.MasterJabatan{}).Count(&count)
	if count == 0 {
		log.Println("Seeding default master jabatan...")
		jabatans := []models.MasterJabatan{
			{NamaJabatan: "Guru Ahli Pertama", JenisJabatan: "Fungsional"},
			{NamaJabatan: "Guru Ahli Muda", JenisJabatan: "Fungsional"},
			{NamaJabatan: "Guru Madya", JenisJabatan: "Fungsional"},
			{NamaJabatan: "Guru Utama", JenisJabatan: "Fungsional"},
			{NamaJabatan: "Kepala Sekolah", JenisJabatan: "Fungsional"},
			{NamaJabatan: "Pranata Komputer Ahli Pertama", JenisJabatan: "Fungsional"},
			{NamaJabatan: "Analis Kepegawaian Ahli Pertama", JenisJabatan: "Fungsional"},
			{NamaJabatan: "Staf Administrasi", JenisJabatan: "Pelaksana"},
			{NamaJabatan: "Pengadministrasi Umum", JenisJabatan: "Pelaksana"},
			{NamaJabatan: "Pengelola Keuangan", JenisJabatan: "Pelaksana"},
			{NamaJabatan: "Operator Komputer", JenisJabatan: "Pelaksana"},
		}
		db.Create(&jabatans)
	}

	// 4. Seed Master Tempat Tugas
	db.Model(&models.MasterTempatTugas{}).Count(&count)
	if count == 0 {
		log.Println("Seeding default master tempat tugas...")
		places := []models.MasterTempatTugas{
			{NamaTempat: "Dinas Pendidikan dan Kebudayaan Daerah Morowali Utara", JenisTempat: "Dinas"},
			{NamaTempat: "Sub Bagian Umum dan Kepegawaian", JenisTempat: "Dinas"},
			{NamaTempat: "Bidang Pendidikan Dasar", JenisTempat: "Dinas"},
			{NamaTempat: "Bidang Pendidikan Menengah", JenisTempat: "Dinas"},
			{NamaTempat: "SDN 01 Lembo", JenisTempat: "Sekolah"},
			{NamaTempat: "SMPN 01 Lembo", JenisTempat: "Sekolah"},
			{NamaTempat: "SMAN 01 Petasia", JenisTempat: "Sekolah"},
			{NamaTempat: "SMKN 01 Morowali Utara", JenisTempat: "Sekolah"},
		}
		db.Create(&places)
	}
}


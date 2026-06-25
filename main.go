package main

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
    
    // Sesuaikan ini dengan struktur folder Anda
	"sipecut/api/models"
)

func main() {
	db, err := gorm.Open(sqlite.Open("sipecut.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Gagal koneksi DB: %v", err)
	}

	// Reset password admin
	adminHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	result := db.Model(&models.User{}).Where("nip = ?", "admin").Updates(map[string]interface{}{
		"password_hash": string(adminHash),
		"status":        "aktif",
	})
	if result.RowsAffected == 0 {
		// User belum ada, buat baru
		adminUser := models.User{
			NIP:          "admin",
			PasswordHash: string(adminHash),
			Role:         "admin",
			NoHP:         "081234567890",
			Status:       "aktif",
		}
		db.Create(&adminUser)
		fmt.Println("✅ Akun admin berhasil dibuat dengan password: admin123")
	} else {
		fmt.Println("✅ Password admin berhasil direset menjadi: admin123")
	}

	// Verifikasi
	var user models.User
	if err := db.Where("nip = ?", "admin").First(&user).Error; err != nil {
		log.Fatalf("❌ Gagal memverifikasi user admin: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("admin123")); err != nil {
		log.Fatalf("❌ Verifikasi password gagal: %v", err)
	}
	fmt.Printf("✅ Verifikasi sukses — NIP: %s | Role: %s | Status: %s\n", user.NIP, user.Role, user.Status)
}

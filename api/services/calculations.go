package services

import (
	"strings"
	"time"

	"gorm.io/gorm"

	"sipecut/models"
)

// CalculateAge returns age in years and months
func CalculateAge(birthDate time.Time) (int, int) {
	now := time.Now()
	years := now.Year() - birthDate.Year()
	months := int(now.Month()) - int(birthDate.Month())
	days := now.Day() - birthDate.Day()

	if months < 0 || (months == 0 && days < 0) {
		years--
		months = 12 + months
	}
	
	if days < 0 {
		months--
		if months < 0 {
			months = 11
		}
	}

	return years, months
}

// NormalizeJenisJabatan converts any variations like "Jabatan Fungsional" or "Fungsional" into standard "Fungsional", etc.
func NormalizeJenisJabatan(s string) string {
	s = strings.TrimSpace(s)
	lower := strings.ToLower(s)
	if strings.Contains(lower, "fungsional") {
		return "Fungsional"
	}
	if strings.Contains(lower, "pelaksana") {
		return "Pelaksana"
	}
	if strings.Contains(lower, "struktural") {
		return "Struktural"
	}
	return s
}

// CalculateNextKgbDate computes next KGB date based on master rules
func CalculateNextKgbDate(db *gorm.DB, jenisJabatan, jabatan string, terakhir time.Time) time.Time {
	normJenis := NormalizeJenisJabatan(jenisJabatan)
	var rule models.KgbCycleRule
	// Look up specific rule first
	err := db.Where("jenis_jabatan = ? AND jabatan = ?", normJenis, jabatan).First(&rule).Error
	if err != nil {
		// Fallback to wildcard rule
		err = db.Where("jenis_jabatan = ? AND jabatan = ?", normJenis, "*").First(&rule).Error
		if err != nil {
			// Default fallback: 2 years
			return terakhir.AddDate(2, 0, 0)
		}
	}
	return terakhir.AddDate(rule.SiklusTahun, 0, 0)
}

// CalculateNextPangkatDate computes next Rank date based on master rules
func CalculateNextPangkatDate(db *gorm.DB, jenisJabatan, jabatan string, terakhir time.Time) time.Time {
	normJenis := NormalizeJenisJabatan(jenisJabatan)
	var rule models.PangkatCycleRule
	// Look up specific rule first
	err := db.Where("jenis_jabatan = ? AND jabatan = ?", normJenis, jabatan).First(&rule).Error
	if err != nil {
		// Fallback to wildcard rule
		err = db.Where("jenis_jabatan = ? AND jabatan = ?", normJenis, "*").First(&rule).Error
		if err != nil {
			// Default fallback: 4 years
			return terakhir.AddDate(4, 0, 0)
		}
	}
	return terakhir.AddDate(rule.SiklusTahun, 0, 0)
}

// CalculatePensionDate computes retirement date based on master rules
func CalculatePensionDate(db *gorm.DB, jenisJabatan, jabatan string, lahir time.Time) time.Time {
	normJenis := NormalizeJenisJabatan(jenisJabatan)
	var rule models.PensionRule
	// Look up specific rule first
	err := db.Where("jenis_jabatan = ? AND jabatan = ?", normJenis, jabatan).First(&rule).Error
	if err != nil {
		// Fallback to wildcard rule
		err = db.Where("jenis_jabatan = ? AND jabatan = ?", normJenis, "*").First(&rule).Error
		if err != nil {
			// Default fallback: 58 years
			return lahir.AddDate(58, 0, 0)
		}
	}
	return lahir.AddDate(rule.BatasUsiaPensiun, 0, 0)
}

// CalculateNextKgbDateInMem computes next KGB date based on in-memory rules
func CalculateNextKgbDateInMem(rules []models.KgbCycleRule, jenisJabatan, jabatan string, terakhir time.Time) time.Time {
	normJenis := NormalizeJenisJabatan(jenisJabatan)
	// Specific lookup
	for _, r := range rules {
		if NormalizeJenisJabatan(r.JenisJabatan) == normJenis && strings.EqualFold(r.Jabatan, jabatan) {
			return terakhir.AddDate(r.SiklusTahun, 0, 0)
		}
	}
	// Wildcard lookup
	for _, r := range rules {
		if NormalizeJenisJabatan(r.JenisJabatan) == normJenis && r.Jabatan == "*" {
			return terakhir.AddDate(r.SiklusTahun, 0, 0)
		}
	}
	return terakhir.AddDate(2, 0, 0)
}

// CalculateNextPangkatDateInMem computes next Rank date based on in-memory rules
func CalculateNextPangkatDateInMem(rules []models.PangkatCycleRule, jenisJabatan, jabatan string, terakhir time.Time) time.Time {
	normJenis := NormalizeJenisJabatan(jenisJabatan)
	// Specific lookup
	for _, r := range rules {
		if NormalizeJenisJabatan(r.JenisJabatan) == normJenis && strings.EqualFold(r.Jabatan, jabatan) {
			return terakhir.AddDate(r.SiklusTahun, 0, 0)
		}
	}
	// Wildcard lookup
	for _, r := range rules {
		if NormalizeJenisJabatan(r.JenisJabatan) == normJenis && r.Jabatan == "*" {
			return terakhir.AddDate(r.SiklusTahun, 0, 0)
		}
	}
	return terakhir.AddDate(4, 0, 0)
}

// CalculatePensionDateInMem computes retirement date based on in-memory rules
func CalculatePensionDateInMem(rules []models.PensionRule, jenisJabatan, jabatan string, lahir time.Time) time.Time {
	normJenis := NormalizeJenisJabatan(jenisJabatan)
	// Specific lookup
	for _, r := range rules {
		if NormalizeJenisJabatan(r.JenisJabatan) == normJenis && strings.EqualFold(r.Jabatan, jabatan) {
			return lahir.AddDate(r.BatasUsiaPensiun, 0, 0)
		}
	}
	// Wildcard lookup
	for _, r := range rules {
		if NormalizeJenisJabatan(r.JenisJabatan) == normJenis && r.Jabatan == "*" {
			return lahir.AddDate(r.BatasUsiaPensiun, 0, 0)
		}
	}
	return lahir.AddDate(58, 0, 0)
}

// RecalculateEmployeeDatesInMem computes dates/status for an employee without executing DB calls
func RecalculateEmployeeDatesInMem(
	emp *models.Employee,
	pensionRules []models.PensionRule,
	kgbRules []models.KgbCycleRule,
	pangkatRules []models.PangkatCycleRule,
	criteria string,
) {
	zeroTime := time.Time{}

	// Normalise JenisJabatan on employee record itself for DB consistency
	emp.JenisJabatan = NormalizeJenisJabatan(emp.JenisJabatan)

	// KGB
	if emp.TanggalKgbTerakhir.IsZero() || emp.TanggalKgbTerakhir.Year() <= 1 {
		emp.TanggalKgbBerikutnya = zeroTime
	} else {
		emp.TanggalKgbBerikutnya = CalculateNextKgbDateInMem(kgbRules, emp.JenisJabatan, emp.Jabatan, emp.TanggalKgbTerakhir)
	}

	// Pangkat
	if emp.TanggalKenaikanPangkatTerakhir.IsZero() || emp.TanggalKenaikanPangkatTerakhir.Year() <= 1 {
		emp.TanggalKenaikanPangkatBerikutnya = zeroTime
	} else {
		emp.TanggalKenaikanPangkatBerikutnya = CalculateNextPangkatDateInMem(pangkatRules, emp.JenisJabatan, emp.Jabatan, emp.TanggalKenaikanPangkatTerakhir)
	}

	// Pension
	if !emp.TanggalLahir.IsZero() && emp.TanggalLahir.Year() > 1 {
		emp.TanggalPensiun = CalculatePensionDateInMem(pensionRules, emp.JenisJabatan, emp.Jabatan, emp.TanggalLahir)
	}

	// Status update based on criteria
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
}

// RecalculateEmployeeDates updates next KGB, next rank, pension date, and checks retirement status (loads rules on demand)
func RecalculateEmployeeDates(db *gorm.DB, emp *models.Employee) error {
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

	RecalculateEmployeeDatesInMem(emp, pensionRules, kgbRules, pangkatRules, criteria)
	return nil
}

// RecalculateAllEmployees recalculates dates/status for all employees (fully cached and optimized in one transaction)
func RecalculateAllEmployees(db *gorm.DB) error {
	var pensionRules []models.PensionRule
	if err := db.Find(&pensionRules).Error; err != nil {
		return err
	}

	var kgbRules []models.KgbCycleRule
	if err := db.Find(&kgbRules).Error; err != nil {
		return err
	}

	var pangkatRules []models.PangkatCycleRule
	if err := db.Find(&pangkatRules).Error; err != nil {
		return err
	}

	var setting models.AppSetting
	criteria := "dokumen_pemberhentian"
	if err := db.Where("key = ?", "kriteria_pensiun_lengkap").First(&setting).Error; err == nil {
		criteria = setting.Value
	}

	var employees []models.Employee
	if err := db.Find(&employees).Error; err != nil {
		return err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i := range employees {
		emp := &employees[i]
		
		oldKgbNext := emp.TanggalKgbBerikutnya
		oldPangkatNext := emp.TanggalKenaikanPangkatBerikutnya
		oldPension := emp.TanggalPensiun
		oldStatus := emp.StatusKepegawaian
		oldJenis := emp.JenisJabatan

		RecalculateEmployeeDatesInMem(emp, pensionRules, kgbRules, pangkatRules, criteria)

		changed := !emp.TanggalKgbBerikutnya.Equal(oldKgbNext) ||
			!emp.TanggalKenaikanPangkatBerikutnya.Equal(oldPangkatNext) ||
			!emp.TanggalPensiun.Equal(oldPension) ||
			emp.StatusKepegawaian != oldStatus ||
			emp.JenisJabatan != oldJenis

		if changed {
			if err := tx.Save(emp).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}

// CheckAllEmployeesRetirementStatus runs through all active employees and checks if their pension status should update
func CheckAllEmployeesRetirementStatus(db *gorm.DB) error {
	var pensionRules []models.PensionRule
	if err := db.Find(&pensionRules).Error; err != nil {
		return err
	}

	var kgbRules []models.KgbCycleRule
	if err := db.Find(&kgbRules).Error; err != nil {
		return err
	}

	var pangkatRules []models.PangkatCycleRule
	if err := db.Find(&pangkatRules).Error; err != nil {
		return err
	}

	var setting models.AppSetting
	criteria := "dokumen_pemberhentian"
	if err := db.Where("key = ?", "kriteria_pensiun_lengkap").First(&setting).Error; err == nil {
		criteria = setting.Value
	}

	var employees []models.Employee
	// Find active employees whose pension date is in the past
	err := db.Where("status_kepegawaian = ? AND tanggal_pensiun < ?", "Aktif", time.Now()).Find(&employees).Error
	if err != nil {
		return err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i := range employees {
		emp := &employees[i]
		oldStatus := emp.StatusKepegawaian
		
		RecalculateEmployeeDatesInMem(emp, pensionRules, kgbRules, pangkatRules, criteria)

		if emp.StatusKepegawaian != oldStatus {
			if err := tx.Save(emp).Error; err != nil {
				tx.Rollback()
				return err
			}
			
			// Log history
			tx.Create(&models.RequestHistory{
				RequestType: "employee_status",
				RequestID:   emp.ID,
				StatusLama:  oldStatus,
				StatusBaru:  emp.StatusKepegawaian,
				Catatan:     "Status kepegawaian otomatis diubah ke Pensiun karena tanggal pensiun terlewati dan dokumen pemberhentian pembayaran (jika disyaratkan) lengkap.",
				ChangedBy:   "system",
				ChangedAt:   time.Now(),
			})
		}
	}

	return tx.Commit().Error
}


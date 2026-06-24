package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"sipecut/config"
	"sipecut/models"
)

// GetDashboardSummary aggregates statistics for Admin and Employee dashboard
func GetDashboardSummary(c *gin.Context) {
	db := config.GetDB()
	role, _ := c.Get("role")
	nip, _ := c.Get("nip")

	// 1. Employee Dashboard Summary
	if role == "employee" {
		var emp models.Employee
		if err := db.Where("nip = ?", nip).First(&emp).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pegawai tidak ditemukan"})
			return
		}

		// Count employee's own requests
		var totalLeaves int64
		var pendingLeaves int64
		var activeChanges int64
		var totalBA int64

		db.Model(&models.LeaveRequest{}).Where("employee_id = ?", emp.ID).Count(&totalLeaves)
		db.Model(&models.LeaveRequest{}).Where("employee_id = ? AND status = ?", emp.ID, "Diajukan").Count(&pendingLeaves)
		db.Model(&models.DataChangeRequest{}).Where("employee_id = ? AND status = ?", emp.ID, "Diajukan").Count(&activeChanges)
		db.Model(&models.BeritaAcara{}).Where("employee_id = ?", emp.ID).Count(&totalBA)

		c.JSON(http.StatusOK, gin.H{
			"role":            "employee",
			"employee_name":   emp.Nama,
			"nip":             emp.NIP,
			"next_kgb":        emp.TanggalKgbBerikutnya,
			"next_pangkat":    emp.TanggalKenaikanPangkatBerikutnya,
			"next_pension":    emp.TanggalPensiun,
			"total_leaves":    totalLeaves,
			"pending_leaves":  pendingLeaves,
			"active_changes":  activeChanges,
			"total_ba":        totalBA,
		})
		return
	}

	// 2. Admin Dashboard Summary
	targetRangeMonths := 3
	var setting models.AppSetting
	if err := db.Where("key = ?", "dashboard_alert_months").First(&setting).Error; err == nil {
		if val, err := strconv.Atoi(setting.Value); err == nil && val > 0 {
			targetRangeMonths = val
		}
	}

	now := time.Now()
	futureLimit := now.AddDate(0, targetRangeMonths, 0)

	var totalActive int64
	var totalRetired int64
	var akanKgb int64
	var akanPangkat int64
	var akanPensiun int64

	// Count active and retired
	db.Model(&models.Employee{}).Where("status_kepegawaian = ?", "Aktif").Count(&totalActive)
	db.Model(&models.Employee{}).Where("status_kepegawaian = ?", "Pensiun").Count(&totalRetired)

	// Count alert dates (only count active employees)
	db.Model(&models.Employee{}).Where("status_kepegawaian = ? AND tanggal_kgb_berikutnya >= ? AND tanggal_kgb_berikutnya <= ?", "Aktif", now, futureLimit).Count(&akanKgb)
	db.Model(&models.Employee{}).Where("status_kepegawaian = ? AND tanggal_kenaikan_pangkat_berikutnya >= ? AND tanggal_kenaikan_pangkat_berikutnya <= ?", "Aktif", now, futureLimit).Count(&akanPangkat)
	db.Model(&models.Employee{}).Where("status_kepegawaian = ? AND tanggal_pensiun >= ? AND tanggal_pensiun <= ?", "Aktif", now, futureLimit).Count(&akanPensiun)

	// Counter Badges
	var pendingLeaves int64
	var pendingChanges int64
	var pendingBA int64

	db.Model(&models.LeaveRequest{}).Where("status = ?", "Diajukan").Count(&pendingLeaves)
	db.Model(&models.DataChangeRequest{}).Where("status = ?", "Diajukan").Count(&pendingChanges)
	db.Model(&models.BeritaAcara{}).Where("status = ?", "Diajukan").Count(&pendingBA)

	c.JSON(http.StatusOK, gin.H{
		"role":            "admin",
		"total_active":    totalActive,
		"total_retired":   totalRetired,
		"akan_kgb":        akanKgb,
		"akan_pangkat":    akanPangkat,
		"akan_pensiun":    akanPensiun,
		"pending_leaves":  pendingLeaves,
		"pending_changes": pendingChanges,
		"pending_ba":      pendingBA,
		"config_months":   targetRangeMonths,
	})
}

// GetAuditHistory retrieves request_history audit logs (admin only)
func GetAuditHistory(c *gin.Context) {
	db := config.GetDB()

	drawVal := c.Query("draw")
	startVal := c.Query("start")
	lengthVal := c.Query("length")
	searchVal := c.Query("search[value]")

	draw, _ := strconv.Atoi(drawVal)
	start, _ := strconv.Atoi(startVal)
	length, _ := strconv.Atoi(lengthVal)
	if length <= 0 {
		length = 10
	}

	var totalRecords int64
	var filteredRecords int64

	query := db.Model(&models.RequestHistory{})
	query.Count(&totalRecords)

	if searchVal != "" {
		searchLike := "%" + searchVal + "%"
		query = query.Where("request_type LIKE ? OR status_baru LIKE ? OR catatan LIKE ? OR changed_by LIKE ?", searchLike, searchLike, searchLike, searchLike)
	}

	query.Count(&filteredRecords)
	query.Order("changed_at desc")

	var history []models.RequestHistory
	query.Limit(length).Offset(start).Find(&history)

	c.JSON(http.StatusOK, gin.H{
		"draw":            draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            history,
	})
}

// GetNotificationLogs retrieves whatsapp notification logs (admin only)
func GetNotificationLogs(c *gin.Context) {
	db := config.GetDB()

	drawVal := c.Query("draw")
	startVal := c.Query("start")
	lengthVal := c.Query("length")
	searchVal := c.Query("search[value]")

	draw, _ := strconv.Atoi(drawVal)
	start, _ := strconv.Atoi(startVal)
	length, _ := strconv.Atoi(lengthVal)
	if length <= 0 {
		length = 10
	}

	var totalRecords int64
	var filteredRecords int64

	query := db.Model(&models.NotificationLog{})
	query.Count(&totalRecords)

	if searchVal != "" {
		searchLike := "%" + searchVal + "%"
		query = query.Where("message LIKE ? OR status LIKE ?", searchLike, searchLike)
	}

	query.Count(&filteredRecords)
	query.Order("sent_at desc")

	var logs []models.NotificationLog
	query.Limit(length).Offset(start).Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"draw":            draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            logs,
	})
}

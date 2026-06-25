package routes

import (
	"github.com/gin-gonic/gin"

	"sipecut/controllers"
	"sipecut/middleware"
)

// CORSMiddleware allows cross-origin requests for local development
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// SetupRouter registers all routing groups and controllers
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())

	// Serve Uploaded Files statically
	r.Static("/api/uploads", "./uploads")

	api := r.Group("/api")
	{
		// 1. Anonymous routes
		api.POST("/login", middleware.RateLimiter(), controllers.Login)
		api.POST("/refresh-token", controllers.RefreshToken)

		auth := api.Group("/auth")
		{
			auth.POST("/login", middleware.RateLimiter(), controllers.Login)
			auth.POST("/refresh", controllers.RefreshToken)
		}

		// 2. Authenticated routes (all roles)
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			authorized.GET("/dashboard/summary", controllers.GetDashboardSummary)
			authorized.PUT("/auth/change-password", controllers.ChangePassword)
			authorized.POST("/change-password", controllers.ChangePassword)
			authorized.GET("/me", controllers.GetMe)

			// Employees
			authorized.GET("/employees", controllers.GetEmployees)
			authorized.GET("/employees/:nip", controllers.GetEmployeeByNIP)
			authorized.GET("/employees/:nip/kgb-history", controllers.GetKgbHistory)
			authorized.GET("/employees/:nip/pangkat-history", controllers.GetPangkatHistory)

			// Data Change Requests
			authorized.POST("/data-change-requests", controllers.SubmitDataChangeRequest)
			authorized.GET("/data-change-requests", controllers.GetDataChangeRequests)

			// Leave Requests — admin can also submit on behalf of employee
			authorized.POST("/leave-requests", controllers.SubmitLeaveRequest)
			authorized.GET("/leave-requests", controllers.GetLeaveRequests)

			// Berita Acara
			authorized.POST("/berita-acara", controllers.SubmitBeritaAcara)
			authorized.GET("/berita-acara", controllers.GetBeritaAcara)

			// Master reference data — read accessible to all authenticated users (for autocomplete)
			authorized.GET("/master/jabatan", controllers.GetMasterJabatan)
			authorized.GET("/master/tempat-tugas", controllers.GetMasterTempatTugas)
		}

		// 3. Admin-Only routes
		admin := api.Group("")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			admin.POST("/auth/reset-password", controllers.ResetPassword)
			admin.GET("/admin/pegawai", controllers.AdminGetPegawai)
			admin.POST("/admin/pegawai", controllers.AdminCreatePegawai)
			admin.PUT("/admin/pegawai/:nip/akses", controllers.AdminUpdateAkses)
			admin.PUT("/admin/pegawai/:nip/reset-password", controllers.AdminResetPassword)
			admin.GET("/admin/menu", controllers.AdminGetMenu)

			// Retired Employees
			admin.GET("/employees/retired", controllers.GetRetiredEmployees)

			// Employee CRUD
			admin.POST("/employees", controllers.CreateEmployee)
			admin.PUT("/employees/:nip", controllers.UpdateEmployee)

			// Excel Import / Export / Template
			admin.POST("/employees/import", controllers.ImportEmployees)
			admin.GET("/employees/export", controllers.ExportEmployees)
			admin.GET("/employees/import-template", controllers.DownloadImportTemplate)
			admin.DELETE("/employees/all", controllers.DeleteAllEmployees)

			// KGB/Pangkat manual additions
			admin.POST("/employees/:nip/kgb-history", controllers.AddKgbHistory)
			admin.POST("/employees/:nip/pangkat-history", controllers.AddPangkatHistory)

			// Approval workflows
			admin.PUT("/leave-requests/:id/status", controllers.UpdateLeaveStatus)
			admin.POST("/leave-requests/:id/letters", controllers.UploadSignedLetter)
			admin.POST("/leave-requests/:id/regenerate-docs", controllers.RegenerateLeaveDocuments)
			admin.PUT("/data-change-requests/:id/status", controllers.UpdateDataChangeRequestStatus)
			admin.PUT("/berita-acara/:id/status", controllers.UpdateBeritaAcaraStatus)

			// Master Rules Settings
			admin.GET("/master/pension-rules", controllers.GetPensionRules)
			admin.PUT("/master/pension-rules", controllers.UpdatePensionRule)
			admin.DELETE("/master/pension-rules/:id", controllers.DeletePensionRule)

			admin.GET("/master/kgb-cycle-rules", controllers.GetKgbRules)
			admin.PUT("/master/kgb-cycle-rules", controllers.UpdateKgbRule)
			admin.DELETE("/master/kgb-cycle-rules/:id", controllers.DeleteKgbRule)

			admin.GET("/master/pangkat-cycle-rules", controllers.GetPangkatRules)
			admin.PUT("/master/pangkat-cycle-rules", controllers.UpdatePangkatRule)
			admin.DELETE("/master/pangkat-cycle-rules/:id", controllers.DeletePangkatRule)

			admin.GET("/master/settings", controllers.GetSettings)
			admin.PUT("/master/settings", controllers.UpdateSettings)

			// Master Jabatan & Tempat Tugas CRUD
			admin.POST("/master/jabatan", controllers.CreateMasterJabatan)
			admin.PUT("/master/jabatan/:id", controllers.UpdateMasterJabatan)
			admin.DELETE("/master/jabatan/all", controllers.DeleteAllMasterJabatan)
			admin.DELETE("/master/jabatan/:id", controllers.DeleteMasterJabatan)

			admin.POST("/master/tempat-tugas", controllers.CreateMasterTempatTugas)
			admin.PUT("/master/tempat-tugas/:id", controllers.UpdateMasterTempatTugas)
			admin.DELETE("/master/tempat-tugas/all", controllers.DeleteAllMasterTempatTugas)
			admin.DELETE("/master/tempat-tugas/:id", controllers.DeleteMasterTempatTugas)

			// Master Bulk Import & Template Download
			admin.GET("/master/import-template", controllers.DownloadMasterTemplate)
			admin.POST("/master/import", controllers.ImportMasterData)

			// Audit & Notification Logs
			admin.GET("/dashboard/audit-history", controllers.GetAuditHistory)
			admin.GET("/dashboard/notification-logs", controllers.GetNotificationLogs)
		}
	}

	return r
}

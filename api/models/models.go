package models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents user roles
type Role struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	NamaRole string `gorm:"type:varchar(50);uniqueIndex;not null" json:"nama_role"`
}

// Menu represents system modules
type Menu struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	NamaMenu string `gorm:"type:varchar(100);not null" json:"nama_menu"`
	KodeMenu string `gorm:"type:varchar(50);uniqueIndex;not null" json:"kode_menu"`
	ParentID *uint  `json:"parent_id,omitempty"`
	Parent   *Menu  `gorm:"foreignKey:ParentID" json:"-"`
}

// User represents credentials and roles
type User struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	NIP               string         `gorm:"column:nip;type:varchar(50);uniqueIndex;not null" json:"nip"`
	PasswordHash      string         `gorm:"type:varchar(255);not null" json:"-"`
	Role              string         `gorm:"type:varchar(20);not null" json:"role"` // "admin" or "employee"
	NoHP              string         `gorm:"type:varchar(20)" json:"no_hp"`
	Status            string         `gorm:"type:varchar(20);default:'aktif'" json:"status"` // "aktif" or "nonaktif"
	RoleID            uint           `gorm:"default:2" json:"role_id"`
	RoleRel           Role           `gorm:"foreignKey:RoleID;references:ID" json:"role_rel"`
	IsPasswordDefault bool           `gorm:"default:true" json:"is_password_default"`
	Menus             []Menu         `gorm:"many2many:pegawai_menu;foreignKey:ID;joinForeignKey:user_id;References:ID;joinReferences:menu_id" json:"menus,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// Employee represents the employee profile
type Employee struct {
	ID                               uint           `gorm:"primaryKey" json:"id"`
	NIP                              string         `gorm:"column:nip;type:varchar(50);uniqueIndex;not null" json:"nip"`
	Nama                             string         `gorm:"type:varchar(100);not null" json:"nama"`
	JenisJabatan                     string         `gorm:"type:varchar(50);not null" json:"jenis_jabatan"` // "Fungsional" or "Pelaksana"
	Jabatan                          string         `gorm:"type:varchar(100);not null" json:"jabatan"`
	TempatLahir                      string         `gorm:"type:varchar(100)" json:"tempat_lahir"`
	TanggalLahir                     time.Time      `json:"tanggal_lahir"`
	TempatTugas                      string         `gorm:"type:varchar(100)" json:"tempat_tugas"`
	JenisTempat                      string         `gorm:"type:varchar(50)" json:"jenis_tempat"` // "Dinas" or "Sekolah"
	Pengangkatan                     string         `gorm:"type:varchar(100)" json:"pengangkatan"`
	TanggalKgbTerakhir               time.Time      `json:"tanggal_kgb_terakhir"`
	TanggalKgbBerikutnya             time.Time      `json:"tanggal_kgb_berikutnya"`
	TanggalKenaikanPangkatTerakhir   time.Time      `json:"tanggal_kenaikan_pangkat_terakhir"`
	TanggalKenaikanPangkatBerikutnya time.Time      `json:"tanggal_kenaikan_pangkat_berikutnya"`
	TanggalPensiun                   time.Time      `json:"tanggal_pensiun"`
	StatusKepegawaian                string         `gorm:"type:varchar(20);default:'Aktif'" json:"status_kepegawaian"` // "Aktif" or "Pensiun"
	JenisPengangkatan                string         `gorm:"type:varchar(20)" json:"jenis_pengangkatan"`                 // "PNS" or "PPPK"
	FotoProfil                       string         `gorm:"type:text" json:"foto_profil"`
	SkCpnsPppkFile                   string         `gorm:"type:text" json:"sk_cpns_pppk_file"`
	SkPnsFile                        string         `gorm:"type:text" json:"sk_pns_file"`
	SkKgbFile                        string         `gorm:"type:text" json:"sk_kgb_file"`
	SkPangkatFile                    string         `gorm:"type:text" json:"sk_pangkat_file"`
	SkPensiunFile                    string         `gorm:"type:text" json:"sk_pensiun_file"`
	DokumenPemberhentianPembayaran   string         `gorm:"type:text" json:"dokumen_pemberhentian_pembayaran"`
	CreatedAt                        time.Time      `json:"created_at"`
	UpdatedAt                        time.Time      `json:"updated_at"`
	DeletedAt                        gorm.DeletedAt `gorm:"index" json:"-"`
}

// PensionRule defines the retirement age rules configured by Admin
type PensionRule struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	JenisJabatan      string    `gorm:"type:varchar(50);not null" json:"jenis_jabatan"` // "Fungsional" or "Pelaksana"
	Jabatan           string    `gorm:"type:varchar(100);default:'*'" json:"jabatan"`    // specific job title, or "*" for all in this group
	BatasUsiaPensiun  int       `json:"batas_usia_pensiun"`                             // e.g. 58 or 60
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// KgbCycleRule defines the KGB cycles (salary increment)
type KgbCycleRule struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	JenisJabatan string    `gorm:"type:varchar(50);not null" json:"jenis_jabatan"`
	Jabatan      string    `gorm:"type:varchar(100);default:'*'" json:"jabatan"`
	SiklusTahun  int       `json:"siklus_tahun"` // e.g. 2
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PangkatCycleRule defines rank promotion cycles
type PangkatCycleRule struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	JenisJabatan string    `gorm:"type:varchar(50);not null" json:"jenis_jabatan"`
	Jabatan      string    `gorm:"type:varchar(100);default:'*'" json:"jabatan"`
	SiklusTahun  int       `json:"siklus_tahun"` // e.g. 4
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// EmployeeKgbHistory stores previous KGB records
type EmployeeKgbHistory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	EmployeeID uint      `gorm:"not null;index" json:"employee_id"`
	TanggalKgb time.Time `json:"tanggal_kgb"`
	FileSkKgb  string    `gorm:"type:text" json:"file_sk_kgb"`
	CreatedAt  time.Time `json:"created_at"`
}

// EmployeePangkatHistory stores previous pangkat records
type EmployeePangkatHistory struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	EmployeeID            uint      `gorm:"not null;index" json:"employee_id"`
	TanggalKenaikanPangkat time.Time `json:"tanggal_kenaikan_pangkat"`
	FileSkPangkat         string    `gorm:"type:text" json:"file_sk_pangkat"`
	CreatedAt             time.Time `json:"created_at"`
}

// LeaveRequest represents a leave application
type LeaveRequest struct {
	ID            uint              `gorm:"primaryKey" json:"id"`
	EmployeeID    uint              `gorm:"not null;index" json:"employee_id"`
	Employee      Employee          `gorm:"foreignKey:EmployeeID" json:"employee"`
	JenisCuti     string            `gorm:"type:varchar(50);not null" json:"jenis_cuti"` // "Cuti Tahunan Biasa", "Cuti Melahirkan", etc.
	TanggalMulai  time.Time         `json:"tanggal_mulai"`
	TanggalSelesai time.Time        `json:"tanggal_selesai"`
	Status        string            `gorm:"type:varchar(30);default:'Diajukan'" json:"status"` // Diajukan -> Sedang Diproses -> Disetujui/Dikembalikan/Ditolak -> Surat Terunggah -> Selesai
	CatatanAdmin  string            `gorm:"type:text" json:"catatan_admin"`
	Attachments   []LeaveAttachment `gorm:"foreignKey:LeaveRequestID;constraint:OnDelete:CASCADE" json:"attachments"`
	Letters       []LeaveLetter     `gorm:"foreignKey:LeaveRequestID;constraint:OnDelete:CASCADE" json:"letters"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// LeaveAttachment holds uploaded files for a leave request
type LeaveAttachment struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	LeaveRequestID uint   `gorm:"not null;index" json:"leave_request_id"`
	JenisDokumen   string `gorm:"type:varchar(100);not null" json:"jenis_dokumen"` // "SK Terakhir", "Buku KIA", etc.
	FilePath       string `gorm:"type:text;not null" json:"file_path"`
}

// LeaveLetter holds reference to generated and uploaded approved documents
type LeaveLetter struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	LeaveRequestID    uint      `gorm:"not null;index" json:"leave_request_id"`
	JenisSurat        string    `gorm:"type:varchar(20)" json:"jenis_surat"` // "Rekomendasi" or "Formulir"
	FileDocx          string    `gorm:"type:text" json:"file_docx"`          // URL of template docx generated
	FilePdf           string    `gorm:"type:text" json:"file_pdf"`           // URL of template pdf generated
	FileSignedPdf     string    `gorm:"type:text" json:"file_signed_pdf"`    // URL of signed PDF uploaded by admin
	UploadedByAdminAt time.Time `json:"uploaded_by_admin_at"`
}

// DataChangeRequest represents a self-service profile update request
type DataChangeRequest struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	EmployeeID    uint      `gorm:"not null;index" json:"employee_id"`
	Employee      Employee  `gorm:"foreignKey:EmployeeID" json:"employee"`
	DataJSON      string    `gorm:"type:text;not null" json:"data_json"` // JSON string containing field->value changes
	SkTerakhirFile string    `gorm:"type:text;not null" json:"sk_terakhir_file"`
	Status        string    `gorm:"type:varchar(30);default:'Diajukan'" json:"status"` // Diajukan, Sedang Diproses, Disetujui, Dikembalikan, Ditolak
	CatatanAdmin  string    `gorm:"type:text" json:"catatan_admin"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BeritaAcara represents document uploads like BA, ST, SI, SKS
type BeritaAcara struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	EmployeeID   uint      `gorm:"not null;index" json:"employee_id"`
	Employee     Employee  `gorm:"foreignKey:EmployeeID" json:"employee"`
	Jenis        string    `gorm:"type:varchar(20);not null" json:"jenis"` // BA / ST / SI / SKS
	FilePath     string    `gorm:"type:text;not null" json:"file_path"`
	Status       string    `gorm:"type:varchar(30);default:'Diajukan'" json:"status"` // Diajukan, Disetujui, Dikembalikan, Ditolak
	CatatanAdmin string    `gorm:"type:text" json:"catatan_admin"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RequestHistory logs request state transitions
type RequestHistory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RequestType string    `gorm:"type:varchar(50);not null" json:"request_type"` // "leave", "change", "berita_acara"
	RequestID   uint      `gorm:"not null;index" json:"request_id"`
	StatusLama  string    `gorm:"type:varchar(30)" json:"status_lama"`
	StatusBaru  string    `gorm:"type:varchar(30);not null" json:"status_baru"`
	Catatan     string    `gorm:"type:text" json:"catatan"`
	ChangedBy   string    `gorm:"type:varchar(50);not null" json:"changed_by"` // NIP of admin/employee
	ChangedAt   time.Time `json:"changed_at"`
}

// NotificationLog keeps logs of WhatsApp notifications sent
type NotificationLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Channel   string    `gorm:"type:varchar(20);default:'WhatsApp'" json:"channel"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Status    string    `gorm:"type:varchar(20);not null" json:"status"` // "Simulated", "Sent", "Failed"
	SentAt    time.Time `json:"sent_at"`
}

// AppSetting represents general configuration
type AppSetting struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Key       string `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
	Value     string `gorm:"type:text" json:"value"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MasterJabatan is the reference table for job titles used in autocomplete
type MasterJabatan struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	NamaJabatan  string    `gorm:"type:varchar(150);not null" json:"nama_jabatan"`
	JenisJabatan string    `gorm:"type:varchar(50);not null" json:"jenis_jabatan"` // "Fungsional" or "Pelaksana"
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// MasterTempatTugas is the reference table for work units used in autocomplete
type MasterTempatTugas struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	NamaTempat string    `gorm:"type:varchar(150);not null" json:"nama_tempat"`
	JenisTempat string   `gorm:"type:varchar(50);not null" json:"jenis_tempat"` // "Dinas" or "Sekolah"
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}


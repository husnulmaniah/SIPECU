package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"

	"sipecut/config"
	"sipecut/models"
)

// indonesianMonths maps month numbers to Indonesian month names
var indonesianMonths = []string{
	"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

func formatTanggalID(t time.Time) string {
	return fmt.Sprintf("%02d %s %d", t.Day(), indonesianMonths[int(t.Month())], t.Year())
}

func formatTanggalRangeID(start, end time.Time) string {
	if start.Year() == end.Year() {
		if start.Month() == end.Month() {
			return fmt.Sprintf("%02d-%02d %s %d", start.Day(), end.Day(), indonesianMonths[int(start.Month())], start.Year())
		}
		return fmt.Sprintf("%02d %s s/d %02d %s %d", start.Day(), indonesianMonths[int(start.Month())], end.Day(), indonesianMonths[int(end.Month())], start.Year())
	}
	return fmt.Sprintf("%02d %s %d s/d %02d %s %d", start.Day(), indonesianMonths[int(start.Month())], start.Year(), end.Day(), indonesianMonths[int(end.Month())], end.Year())
}

func getKepalaDinasInfo() (string, string) {
	nama := "MOH. RIDWAN DM, S.Ag" // fallback
	nip := "19740111 199803 1 004" // fallback

	db := config.GetDB()
	if db != nil {
		var emp models.Employee
		if err := db.Where("LOWER(jabatan) = ? AND status_kepegawaian = ?", "kepala dinas", "Aktif").First(&emp).Error; err == nil {
			nama = emp.Nama
			nip = emp.NIP
		}
	}
	return nama, nip
}

func findLogoPath() string {
	paths := []string{
		"storage/logo_morowali.png",
		"api/storage/logo_morowali.png",
		"../api/storage/logo_morowali.png",
		"./api/storage/logo_morowali.png",
		"/var/task/api/storage/logo_morowali.png",
		"/var/task/storage/logo_morowali.png",
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "storage/logo_morowali.png" // fallback
}

func buildDocxSignatureXML(now time.Time) string {
	kdNama, kdNip := getKepalaDinasInfo()
	return `<w:tbl>
      <w:tblPr>
        <w:tblW w:w="5000" w:type="pct"/>
        <w:tblBorders>
          <w:top w:val="none"/><w:left w:val="none"/><w:bottom w:val="none"/><w:right w:val="none"/>
          <w:insideH w:val="none"/><w:insideV w:val="none"/>
        </w:tblBorders>
      </w:tblPr>
      <w:tr>
        <w:tc>
          <w:tcPr>
            <w:tcW w:w="3500" w:type="dxa"/>
          </w:tcPr>
          <w:p/>
        </w:tc>
        <w:tc>
          <w:tcPr>
            <w:tcW w:w="6500" w:type="dxa"/>
          </w:tcPr>
          ` + wTimes("Kolonodale, "+formatTanggalID(now), "center", "", "24") + `
          ` + wTimes("Kepala Dinas,", "center", "", "24") + `
          <w:p/><w:p/><w:p/>
          ` + wTimes(kdNama, "center", "BU", "24") + `
          ` + wTimes("NIP: "+kdNip, "center", "", "24") + `
        </w:tc>
      </w:tr>
    </w:tbl>`
}

func hitungMasaKerja(tanggalLahir time.Time, now time.Time) (int, int) {
	startYear := tanggalLahir.Year() + 23
	start := time.Date(startYear, tanggalLahir.Month(), tanggalLahir.Day(), 0, 0, 0, 0, time.UTC)
	diff := now.Sub(start)
	totalMonths := int(diff.Hours() / 24 / 30.44)
	if totalMonths < 0 {
		totalMonths = 0
	}
	return totalMonths / 12, totalMonths % 12
}

// â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
// FORMULIR PERMINTAAN DAN PEMBERIAN CUTI  (PDF only, Times New Roman)
// â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

func GenerateFormulirCutiPdf(req *models.LeaveRequest) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 10, 15)
	pdf.AddPage()

	now := time.Now()
	pageW := 180.0

	// â”€â”€ KOP SURAT â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
	logoPath := findLogoPath()
	if _, err := os.Stat(logoPath); err == nil {
		pdf.Image(logoPath, 16, 10, 18, 22, false, "", 0, "")
	}

	pdf.SetFont("Times", "B", 12)
	pdf.SetX(35)
	pdf.CellFormat(pageW - 20, 5, "PEMERINTAH KABUPATEN MOROWALI UTARA", "", 1, "C", false, 0, "")
	pdf.SetX(35)
	pdf.CellFormat(pageW - 20, 5, "DINAS PENDIDIKAN DAN KEBUDAYAAN DAERAH", "", 1, "C", false, 0, "")
	pdf.SetFont("Times", "", 9)
	pdf.SetX(35)
	pdf.CellFormat(pageW - 20, 4, "Alamat : Jln. Bumi Nangka Kompleks Perkantoran Kode Pos (94971)", "", 1, "C", false, 0, "")
	pdf.SetFont("Times", "B", 12)
	pdf.SetX(35)
	pdf.CellFormat(pageW - 20, 5, "KOLONODALE", "", 1, "C", false, 0, "")

	pdf.SetLineWidth(0.8)
	pdf.Line(15, pdf.GetY()+1, 195, pdf.GetY()+1)
	pdf.SetLineWidth(0.3)
	pdf.Line(15, pdf.GetY()+3, 195, pdf.GetY()+3)
	pdf.Ln(5)

	// Recipient Block
	pdf.SetFont("Times", "", 9)
	pdf.SetX(100)
	pdf.CellFormat(80, 4, fmt.Sprintf("Kolonodale, %s", formatTanggalID(now)), "", 1, "L", false, 0, "")
	pdf.SetX(100)
	pdf.CellFormat(80, 4, "Kepada", "", 1, "L", false, 0, "")
	pdf.SetX(100)
	pdf.SetFont("Times", "B", 9)
	pdf.CellFormat(80, 4, "Yth. Bupati Morowali Utara", "", 1, "L", false, 0, "")
	pdf.SetX(100)
	pdf.CellFormat(80, 4, "Cq. Kepala Badan Kepegawaian dan Pengembangan SDM", "", 1, "L", false, 0, "")
	pdf.SetX(100)
	pdf.SetFont("Times", "", 9)
	pdf.CellFormat(80, 4, "Di -", "", 1, "L", false, 0, "")
	pdf.SetX(110)
	pdf.CellFormat(80, 4, "Kolonodale", "", 1, "L", false, 0, "")
	pdf.Ln(4)

	// â”€â”€ JUDUL â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
	pdf.SetFont("Times", "B", 11)
	pdf.CellFormat(pageW, 6, "FORMULIR PERMINTAAN DAN PEMBERIAN CUTI", "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// â”€â”€ SECTION I: DATA PEGAWAI â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
	sectionHeader(pdf, pageW, "I", "Data Pegawai")

	masaKerjaTahun, masaKerjaBulan := hitungMasaKerja(req.Employee.TanggalLahir, now)

	pdf.SetFont("Times", "", 9)
	// Nama & NIP
	pdf.CellFormat(8, 5, "", "LR", 0, "C", false, 0, "")
	pdf.CellFormat(15, 5, "Nama", "", 0, "L", false, 0, "")
	pdf.CellFormat(3, 5, ":", "", 0, "L", false, 0, "")
	pdf.SetFont("Times", "B", 9)
	pdf.CellFormat(77, 5, req.Employee.Nama, "R", 0, "L", false, 0, "")
	pdf.SetFont("Times", "", 9)
	pdf.CellFormat(15, 5, "NIP", "", 0, "L", false, 0, "")
	pdf.CellFormat(3, 5, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(54, 5, req.Employee.NIP, "R", 1, "L", false, 0, "")

	// Jabatan & Masa Kerja
	pdf.CellFormat(8, 5, "", "LR", 0, "C", false, 0, "")
	pdf.CellFormat(15, 5, "Jabatan", "", 0, "L", false, 0, "")
	pdf.CellFormat(3, 5, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(77, 5, req.Employee.Jabatan, "R", 0, "L", false, 0, "")
	pdf.CellFormat(15, 5, "Masa Kerja", "", 0, "L", false, 0, "")
	pdf.CellFormat(3, 5, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(54, 5, fmt.Sprintf("%d Tahun %d Bulan", masaKerjaTahun, masaKerjaBulan), "R", 1, "L", false, 0, "")

	// Unit Kerja
	pdf.CellFormat(8, 5, "", "LRB", 0, "C", false, 0, "")
	pdf.CellFormat(15, 5, "Unit Kerja", "B", 0, "L", false, 0, "")
	pdf.CellFormat(3, 5, ":", "B", 0, "L", false, 0, "")
	pdf.CellFormat(154, 5, req.Employee.TempatTugas, "RB", 1, "L", false, 0, "")

	// â”€â”€ SECTION II: JENIS CUTI â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
	sectionHeader(pdf, pageW, "II", "Jenis Cuti yang diambil")

	check1, check2, check3, check4, check5, check6 := "", "", "", "", "", ""
	switch req.JenisCuti {
	case "Cuti Tahunan Biasa", "Cuti Tahunan Umroh", "Cuti Tahunan untuk Umroh":
		check1 = "v"
	case "Cuti Sakit":
		check2 = "v"
	case "Cuti Alasan Penting":
		check3 = "v"
	case "Cuti Besar":
		check4 = "v"
	case "Cuti Melahirkan":
		check5 = "v"
	case "Cuti di Luar Tanggungan Negara":
		check6 = "v"
	}

	pdf.CellFormat(8, 5, "", "LR", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "1", "1", 0, "C", false, 0, "")
	pdf.CellFormat(71, 5, "Cuti Tahunan", "1", 0, "L", false, 0, "")
	pdf.CellFormat(10, 5, check1, "1", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "4", "1", 0, "C", false, 0, "")
	pdf.CellFormat(71, 5, "Cuti Besar", "1", 0, "L", false, 0, "")
	pdf.CellFormat(10, 5, check4, "1", 1, "C", false, 0, "")

	pdf.CellFormat(8, 5, "", "LR", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "2", "1", 0, "C", false, 0, "")
	pdf.CellFormat(71, 5, "Cuti Sakit", "1", 0, "L", false, 0, "")
	pdf.CellFormat(10, 5, check2, "1", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "5", "1", 0, "C", false, 0, "")
	pdf.CellFormat(71, 5, "Cuti Melahirkan", "1", 0, "L", false, 0, "")
	pdf.CellFormat(10, 5, check5, "1", 1, "C", false, 0, "")

	pdf.CellFormat(8, 5, "", "LRB", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "3", "1", 0, "C", false, 0, "")
	pdf.CellFormat(71, 5, "Cuti Karena Alasan Penting", "1", 0, "L", false, 0, "")
	pdf.CellFormat(10, 5, check3, "1", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "6", "1", 0, "C", false, 0, "")
	pdf.CellFormat(71, 5, "Cuti di Luar Tanggungan Negara", "1", 0, "L", false, 0, "")
	pdf.CellFormat(10, 5, check6, "1", 1, "C", false, 0, "")

	// â”€â”€ SECTION III: ALASAN CUTI â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
	sectionHeader(pdf, pageW, "III", "Alasan Cuti")

	alasan := req.JenisCuti
	switch req.JenisCuti {
	case "Cuti Tahunan Biasa":
		alasan = "Urusan Keluarga"
	case "Cuti Tahunan Umroh":
		alasan = "Menunaikan Ibadah Umroh"
	case "Cuti Melahirkan":
		alasan = "Melahirkan"
	case "Cuti Sakit":
		alasan = "Sakit / Rawat Inap"
	case "Cuti Alasan Penting":
		alasan = "Alasan Penting / Mendesak"
	}

	pdf.CellFormat(8, 7, "", "LRB", 0, "C", false, 0, "")
	pdf.SetFont("Times", "", 9)
	pdf.CellFormat(172, 7, alasan, "RB", 1, "L", false, 0, "")

	// â”€â”€ SECTION IV: LAMA CUTI â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
	sectionHeader(pdf, pageW, "IV", "Lama Cuti")

	lamaCuti := int(math.Round(req.TanggalSelesai.Sub(req.TanggalMulai).Hours()/24)) + 1
	tanggalRange := formatTanggalRangeID(req.TanggalMulai, req.TanggalSelesai)

	pdf.CellFormat(8, 5, "", "LRB", 0, "C", false, 0, "")
	pdf.CellFormat(20, 5, "Selama", "1", 0, "L", false, 0, "")
	pdf.CellFormat(30, 5, fmt.Sprintf("%d Hari", lamaCuti), "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 5, "Mulai Tanggal", "1", 0, "L", false, 0, "")
	pdf.CellFormat(92, 5, tanggalRange, "1", 1, "C", false, 0, "")

	// â”€â”€ SECTION V: CATATAN CUTI â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
	sectionHeader(pdf, pageW, "V", "Catatan Cuti")

	// Row 1
	pdf.CellFormat(8, 5, "", "LR", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "1", "1", 0, "C", false, 0, "")
	pdf.CellFormat(81, 5, "Cuti Tahunan", "1", 0, "L", false, 0, "")
	pdf.CellFormat(5, 5, "1", "1", 0, "C", false, 0, "")
	pdf.CellFormat(66, 5, "Cuti Tahunan", "1", 0, "L", false, 0, "")
	pdf.CellFormat(15, 5, "", "1", 1, "C", false, 0, "")

	// Row 2
	pdf.CellFormat(8, 5, "", "LR", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 5, "Tahun", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 5, "Sisa", "1", 0, "C", false, 0, "")
	pdf.CellFormat(31, 5, "Keterangan", "1", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "2", "1", 0, "C", false, 0, "")
	pdf.CellFormat(66, 5, "Cuti Sakit", "1", 0, "L", false, 0, "")
	pdf.CellFormat(15, 5, "", "1", 1, "C", false, 0, "")

	// Row 3
	pdf.CellFormat(8, 5, "", "LR", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 5, "N-1", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(31, 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "3", "1", 0, "C", false, 0, "")
	pdf.CellFormat(66, 5, "Cuti Karena Alasan Penting", "1", 0, "L", false, 0, "")
	pdf.CellFormat(15, 5, "", "1", 1, "C", false, 0, "")

	// Row 4
	pdf.CellFormat(8, 5, "", "LR", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 5, "N-2", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(31, 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "4", "1", 0, "C", false, 0, "")
	pdf.CellFormat(66, 5, "Cuti Besar", "1", 0, "L", false, 0, "")
	pdf.CellFormat(15, 5, "", "1", 1, "C", false, 0, "")

	// Row 5
	pdf.CellFormat(8, 5, "", "LR", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 5, "N", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(31, 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "5", "1", 0, "C", false, 0, "")
	pdf.CellFormat(66, 5, "Cuti Melahirkan", "1", 0, "L", false, 0, "")
	pdf.CellFormat(15, 5, "", "1", 1, "C", false, 0, "")

	// Row 6
	pdf.CellFormat(8, 5, "", "LRB", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(31, 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(5, 5, "6", "1", 0, "C", false, 0, "")
	pdf.CellFormat(66, 5, "Cuti di Luar Tanggungan Negara", "1", 0, "L", false, 0, "")
	pdf.CellFormat(15, 5, "", "1", 1, "C", false, 0, "")

	// â”€â”€ SECTION VI: ALAMAT â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
	sectionHeader(pdf, pageW, "VI", "Alamat Selama Menjalankan Cuti")

	phoneNo := "....................."
	db := config.GetDB()
	if db != nil {
		var u models.User
		if err := db.Where("nip = ?", req.Employee.NIP).First(&u).Error; err == nil && u.NoHP != "" {
			phoneNo = u.NoHP
		}
	}

	startY := pdf.GetY()
	pdf.CellFormat(8, 32, "", "1", 0, "C", false, 0, "")

	pdf.SetXY(23, startY)
	pdf.CellFormat(86, 32, "", "1", 0, "L", false, 0, "")
	pdf.SetXY(25, startY+2)
	pdf.SetFont("Times", "", 9)
	alamatText := req.CatatanAdmin
	if alamatText == "" {
		alamatText = "..................................................\n.................................................."
	}
	pdf.MultiCell(82, 4, alamatText, "", "L", false)

	pdf.SetXY(109, startY)
	pdf.CellFormat(86, 5, "Telp: "+phoneNo, "1", 1, "L", false, 0, "")
	pdf.SetX(109)
	pdf.CellFormat(86, 5, "Hormat Saya", "1", 1, "C", false, 0, "")
	pdf.SetX(109)
	pdf.CellFormat(86, 12, "", "1", 1, "C", false, 0, "")
	pdf.SetX(109)
	pdf.SetFont("Times", "BU", 9)
	pdf.CellFormat(86, 5, req.Employee.Nama, "1", 1, "C", false, 0, "")
	pdf.SetX(109)
	pdf.SetFont("Times", "", 9)
	pdf.CellFormat(86, 5, "NIP: "+req.Employee.NIP, "1", 1, "C", false, 0, "")

	// â”€â”€ SECTION VII: PERTIMBANGAN ATASAN â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
	pdf.SetXY(15, startY+32)
	sectionHeader(pdf, pageW, "VII", "Pertimbangan Atasan Langsung")

	startY = pdf.GetY()
	pdf.CellFormat(8, 32, "", "1", 0, "C", false, 0, "")

	pdf.SetXY(23, startY)
	pdf.SetFont("Times", "", 7.5)
	pdf.CellFormat(21.5, 5, "Disetujui", "1", 0, "C", false, 0, "")
	pdf.CellFormat(21.5, 5, "Perubahan", "1", 0, "C", false, 0, "")
	pdf.CellFormat(21.5, 5, "Ditangguhkan", "1", 0, "C", false, 0, "")
	pdf.CellFormat(21.5, 5, "Tdk Disetujui", "1", 1, "C", false, 0, "")

	pdf.SetX(23)
	pdf.CellFormat(21.5, 27, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(21.5, 27, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(21.5, 27, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(21.5, 27, "", "1", 1, "C", false, 0, "")

	kdNama, kdNip := getKepalaDinasInfo()

	pdf.SetXY(109, startY)
	pdf.SetFont("Times", "", 9)
	pdf.CellFormat(86, 5, "Kepala Dinas", "1", 1, "C", false, 0, "")
	pdf.SetX(109)
	pdf.CellFormat(86, 17, "", "1", 1, "C", false, 0, "")
	pdf.SetX(109)
	pdf.SetFont("Times", "BU", 9)
	pdf.CellFormat(86, 5, kdNama, "1", 1, "C", false, 0, "")
	pdf.SetX(109)
	pdf.SetFont("Times", "", 9)
	pdf.CellFormat(86, 5, "NIP: "+kdNip, "1", 1, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
// SURAT REKOMENDASI IZIN CUTI  (Word + PDF, Times New Roman, Justify body)
// â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

func GenerateRekomendasiCutiPdf(req *models.LeaveRequest) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(25, 15, 25)
	pdf.AddPage()

	now := time.Now()
	pageW := 160.0

	// â”€â”€ KOP SURAT â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
	logoPath := findLogoPath()
	if _, err := os.Stat(logoPath); err == nil {
		pdf.Image(logoPath, 25, 12, 18, 22, false, "", 0, "")
	}

	pdf.SetFont("Times", "B", 13)
	pdf.SetX(45)
	pdf.CellFormat(pageW-20, 6, "PEMERINTAH KABUPATEN MOROWALI UTARA", "", 1, "C", false, 0, "")
	pdf.SetX(45)
	pdf.CellFormat(pageW-20, 6, "DINAS PENDIDIKAN DAN KEBUDAYAAN DAERAH", "", 1, "C", false, 0, "")
	pdf.SetFont("Times", "", 10)
	pdf.SetX(45)
	pdf.CellFormat(pageW-20, 5, "Alamat : Jln. Bumi Nangka Kompleks Perkantoran Kode Pos (94971)", "", 1, "C", false, 0, "")
	pdf.SetFont("Times", "B", 12)
	pdf.SetX(45)
	pdf.CellFormat(pageW-20, 5, "KOLONODALE", "", 1, "C", false, 0, "")

	pdf.SetLineWidth(0.8)
	pdf.Line(25, pdf.GetY()+1, 185, pdf.GetY()+1)
	pdf.SetLineWidth(0.3)
	pdf.Line(25, pdf.GetY()+3, 185, pdf.GetY()+3)
	pdf.Ln(8)

	// â”€â”€ NOMOR / LAMPIRAN / PERIHAL â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
	pdf.SetFont("Times", "", 12)
	labelW := 25.0
	colonW := 5.0
	valW := pageW - labelW - colonW

	romawi := toRomawi(int(now.Month()))
	nomorSurat := fmt.Sprintf("800.1.11.4/     /Disdikbud/%s/%d", romawi, now.Year())

	drawSuratRow(pdf, labelW, colonW, valW, "Nomor", nomorSurat)
	drawSuratRow(pdf, labelW, colonW, valW, "Lampiran", "Satu Berkas")
	drawSuratRow(pdf, labelW, colonW, valW, "Perihal", "Rekomendasi Izin Cuti")
	pdf.Ln(8)

	// â”€â”€ TUJUAN â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
	pdf.SetFont("Times", "", 12)
	pdf.CellFormat(pageW, 6, "Yth. Bupati Morowali Utara", "", 1, "L", false, 0, "")
	pdf.CellFormat(pageW, 6, "Cq Kepala Badan Kepegawaian dan", "", 1, "L", false, 0, "")
	pdf.CellFormat(pageW, 6, "Pengembangan SDM", "", 1, "L", false, 0, "")
	pdf.CellFormat(pageW, 6, "Di -", "", 1, "L", false, 0, "")
	pdf.SetX(45)
	pdf.CellFormat(pageW-20, 6, "Tempat", "", 1, "L", false, 0, "")
	pdf.Ln(8)

	// â”€â”€ ISI SURAT (Times New Roman, Justify) â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
	pdf.SetFont("Times", "", 12)

	tanggalRange := formatTanggalRangeID(req.TanggalMulai, req.TanggalSelesai)
	jenisCutiDisplay := normalisasiJenisCuti(req.JenisCuti)

	para1 := fmt.Sprintf(
		"Menindak lanjuti surat permohonan %s atas nama %s; Tanggal %s dengan ini kami tidak keberatan dan menyetujui permohonan tersebut kami teruskan kepada Bapak untuk ditindaklanjuti (Permohonan Terlampir).",
		jenisCutiDisplay, req.Employee.Nama, tanggalRange,
	)

	// "J" = justify (rata kiri-kanan)
	pdf.MultiCell(pageW, 7, para1, "", "J", false)
	pdf.Ln(6)

	pdf.MultiCell(pageW, 7,
		"Demikian Surat Permohonan Cuti ini kami teruskan kepada Bapak, atas pertimbangan Bapak kami ucapkan terima kasih.",
		"", "J", false)
	pdf.Ln(15)

	kdNama, kdNip := getKepalaDinasInfo()

	// â”€â”€ TANDA TANGAN â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
	sigOffsetX := 25 + pageW*0.40
	pdf.SetX(sigOffsetX)
	pdf.SetFont("Times", "", 12)
	pdf.CellFormat(pageW*0.60, 6, fmt.Sprintf("Kolonodale, %s", formatTanggalID(now)), "", 1, "C", false, 0, "")
	pdf.SetX(sigOffsetX)
	pdf.CellFormat(pageW*0.60, 6, "Kepala Dinas,", "", 1, "C", false, 0, "")
	pdf.Ln(20)

	pdf.SetX(sigOffsetX)
	pdf.SetFont("Times", "BU", 12)
	pdf.CellFormat(pageW*0.60, 6, kdNama, "", 1, "C", false, 0, "")
	pdf.SetX(sigOffsetX)
	pdf.SetFont("Times", "", 11)
	pdf.CellFormat(pageW*0.60, 5, "NIP: "+kdNip, "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
// DOCX â€” Surat Rekomendasi (Times New Roman, Justify body)
// â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

func GenerateLeaveDocx(req *models.LeaveRequest) ([]byte, error) {
	return buildDocx(buildRekomendasiXML(req))
}

// GenerateFormulirDocx is kept for backward-compatibility but returns empty
// because Formulir only provides PDF (no Word file per user requirement)
func GenerateFormulirDocx(req *models.LeaveRequest) ([]byte, error) {
	return nil, nil // intentionally no DOCX for Formulir
}

// GenerateLeavePdf alias for Rekomendasi PDF
func GenerateLeavePdf(req *models.LeaveRequest) ([]byte, error) {
	return GenerateRekomendasiCutiPdf(req)
}

func buildDocx(xmlBody string) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	var logoBytes []byte
	logoPath := findLogoPath()
	var readErr error
	if bytesVal, err := os.ReadFile(logoPath); err == nil {
		logoBytes = bytesVal
	} else {
		readErr = err
		// Fallback to a 1x1 transparent PNG so the ZIP relationship is valid and NOT corrupt
		logoBytes = []byte{
			0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
			0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
			0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, 0x00, 0x00, 0x00,
			0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
			0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49,
			0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
		}
	}
	fmt.Printf("LOGO DIAGNOSTIC: path=%s, size=%d, err=%v\n", logoPath, len(logoBytes), readErr)

	contentTypesXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>`
	if len(logoBytes) > 0 {
		contentTypesXML += `<Default Extension="png" ContentType="image/png"/>`
	}
	contentTypesXML += `</Types>`

	contentTypes, _ := zipWriter.Create("[Content_Types].xml")
	_, _ = io.WriteString(contentTypes, contentTypesXML)

	rels, _ := zipWriter.Create("_rels/.rels")
	_, _ = io.WriteString(rels, `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`)

	docRelsXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`
	if len(logoBytes) > 0 {
		docRelsXML += `<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/image1.png"/>`
	}
	docRelsXML += `</Relationships>`

	docRels, _ := zipWriter.Create("word/_rels/document.xml.rels")
	_, _ = io.WriteString(docRels, docRelsXML)

	if len(logoBytes) > 0 {
		mediaFile, _ := zipWriter.Create("word/media/image1.png")
		_, _ = mediaFile.Write(logoBytes)
	}

	documentXML, _ := zipWriter.Create("word/document.xml")
	_, _ = io.WriteString(documentXML, xmlBody)

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func buildRekomendasiXML(req *models.LeaveRequest) string {
	now := time.Now()
	jenisCutiDisplay := normalisasiJenisCuti(req.JenisCuti)
	tanggalRange := formatTanggalRangeID(req.TanggalMulai, req.TanggalSelesai)
	romawi := toRomawi(int(now.Month()))
	nomorSurat := fmt.Sprintf("800.1.11.4/     /Disdikbud/%s/%d", romawi, now.Year())

	para1 := fmt.Sprintf(
		"Menindak lanjuti surat permohonan %s atas nama %s; Tanggal %s dengan ini kami tidak keberatan dan menyetujui permohonan tersebut kami teruskan kepada Bapak untuk ditindaklanjuti (Permohonan Terlampir).",
		jenisCutiDisplay, req.Employee.Nama, tanggalRange,
	)
	para2 := "Demikian Surat Permohonan Cuti ini kami teruskan kepada Bapak, atas pertimbangan Bapak kami ucapkan terima kasih."

	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
            xmlns:wp="http://schemas.openxmlformats.org/wordprocessingdrawingml/2006/main"
            xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
            xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture"
            xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <w:body>
    <w:tbl>
      <w:tblPr>
        <w:tblW w:w="5000" w:type="pct"/>
        <w:tblBorders>
          <w:top w:val="none"/><w:left w:val="none"/><w:bottom w:val="none"/><w:right w:val="none"/>
          <w:insideH w:val="none"/><w:insideV w:val="none"/>
        </w:tblBorders>
        <w:tblCellMar>
          <w:top w:w="0" w:type="dxa"/>
          <w:left w:w="100" w:type="dxa"/>
          <w:bottom w:w="0" w:type="dxa"/>
          <w:right w:w="100" w:type="dxa"/>
        </w:tblCellMar>
      </w:tblPr>
      <w:tr>
        <w:tc>
          <w:tcPr>
            <w:tcW w:w="1400" w:type="dxa"/>
            <w:vAlign w:val="center"/>
          </w:tcPr>
          <w:p>
            <w:pPr><w:jc w:val="center"/></w:pPr>
            <w:r>
              <w:drawing>
                <wp:inline distT="0" distB="0" distL="0" distR="0" simplePos="0" relativeHeight="251658240" behindDoc="0" locked="0" layoutInCell="1" allowOverlap="1">
                  <wp:extent cx="762000" cy="952500"/>
                  <wp:effectExtent l="0" t="0" r="0" b="0"/>
                  <wp:docPr id="1" name="Logo"/>
                  <wp:cNvGraphicFramePr>
                    <a:graphicFrameLocks xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" noChangeAspect="1"/>
                  </wp:cNvGraphicFramePr>
                  <a:graphic xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
                    <a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing">
                      <pic:pic xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">
                        <pic:nvPicPr>
                          <pic:cNvPr id="1" name="logo_morowali.png"/>
                          <pic:cNvPicPr/>
                        </pic:nvPicPr>
                        <pic:blipFill>
                          <a:blip r:embed="rId1" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/>
                          <a:stretch><a:fillRect/></a:stretch>
                        </pic:blipFill>
                        <pic:spPr>
                          <a:xfrm>
                            <a:off x="0" y="0"/>
                            <a:ext cx="762000" cy="952500"/>
                          </a:xfrm>
                          <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
                        </pic:spPr>
                      </pic:pic>
                    </a:graphicData>
                  </a:graphic>
                </wp:inline>
              </w:drawing>
            </w:r>
          </w:p>
        </w:tc>
        <w:tc>
          <w:tcPr>
            <w:tcW w:w="8600" w:type="dxa"/>
            <w:vAlign w:val="center"/>
          </w:tcPr>
          ` + wTimes("PEMERINTAH KABUPATEN MOROWALI UTARA", "center", "B", "26") + `
          ` + wTimes("DINAS PENDIDIKAN DAN KEBUDAYAAN DAERAH", "center", "B", "26") + `
          ` + wTimes("Alamat : Jln. Bumi Nangka Kompleks Perkantoran Kode Pos (94971)", "center", "", "20") + `
          ` + wTimes("KOLONODALE", "center", "B", "24") + `
        </w:tc>
      </w:tr>
    </w:tbl>
    <w:p><w:pPr><w:pBdr><w:bottom w:val="double" w:sz="6" w:space="1" w:color="000000"/></w:pBdr></w:pPr></w:p>
    <w:p/>
    ` + `<w:tbl>
      <w:tblPr>
        <w:tblW w:w="5000" w:type="pct"/>
        <w:tblBorders>
          <w:top w:val="none"/><w:left w:val="none"/><w:bottom w:val="none"/><w:right w:val="none"/>
          <w:insideH w:val="none"/><w:insideV w:val="none"/>
        </w:tblBorders>
        <w:tblCellMar>
          <w:top w:w="0" w:type="dxa"/>
          <w:left w:w="0" w:type="dxa"/>
          <w:bottom w:w="0" w:type="dxa"/>
          <w:right w:w="0" w:type="dxa"/>
        </w:tblCellMar>
      </w:tblPr>
      <w:tr>
        <w:tc>
          <w:tcPr><w:tcW w:w="1200" w:type="dxa"/></w:tcPr>
          ` + wTimes("Nomor", "left", "", "24") + `
        </w:tc>
        <w:tc>
          <w:tcPr><w:tcW w:w="300" w:type="dxa"/></w:tcPr>
          ` + wTimes(":", "left", "", "24") + `
        </w:tc>
        <w:tc>
          <w:tcPr><w:tcW w:w="8500" w:type="dxa"/></w:tcPr>
          ` + wTimes(nomorSurat, "left", "", "24") + `
        </w:tc>
      </w:tr>
      <w:tr>
        <w:tc>
          <w:tcPr><w:tcW w:w="1200" w:type="dxa"/></w:tcPr>
          ` + wTimes("Lampiran", "left", "", "24") + `
        </w:tc>
        <w:tc>
          <w:tcPr><w:tcW w:w="300" w:type="dxa"/></w:tcPr>
          ` + wTimes(":", "left", "", "24") + `
        </w:tc>
        <w:tc>
          <w:tcPr><w:tcW w:w="8500" w:type="dxa"/></w:tcPr>
          ` + wTimes("Satu Berkas", "left", "", "24") + `
        </w:tc>
      </w:tr>
      <w:tr>
        <w:tc>
          <w:tcPr><w:tcW w:w="1200" w:type="dxa"/></w:tcPr>
          ` + wTimes("Perihal", "left", "", "24") + `
        </w:tc>
        <w:tc>
          <w:tcPr><w:tcW w:w="300" w:type="dxa"/></w:tcPr>
          ` + wTimes(":", "left", "", "24") + `
        </w:tc>
        <w:tc>
          <w:tcPr><w:tcW w:w="8500" w:type="dxa"/></w:tcPr>
          ` + wTimes("Rekomendasi Izin Cuti", "left", "", "24") + `
        </w:tc>
      </w:tr>
    </w:tbl>` + `
    <w:p/>
    ` + wTimes("Yth. Bupati Morowali Utara", "left", "", "24") + `
    ` + wTimes("Cq Kepala Badan Kepegawaian dan", "left", "", "24") + `
    ` + wTimes("Pengembangan SDM", "left", "", "24") + `
    ` + wTimes("Di -", "left", "", "24") + `
    ` + wTimes("          Tempat", "left", "", "24") + `
    <w:p/>
    ` + wTimesJustify(para1, "24") + `
    <w:p/>
    ` + wTimesJustify(para2, "24") + `
    <w:p/>
    <w:p/>
    ` + buildDocxSignatureXML(now) + `
  </w:body>
</w:document>`
}

// â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
// PDF HELPER FUNCTIONS
// â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

func sectionHeader(pdf *gofpdf.Fpdf, pageW float64, num, title string) {
	pdf.SetFont("Times", "B", 9)
	pdf.SetFillColor(220, 220, 220)
	pdf.CellFormat(pageW, 5, fmt.Sprintf("  %s.  %s", num, title), "1", 1, "L", true, 0, "")
}

func drawLabelValue(pdf *gofpdf.Fpdf, pageW float64, label1, val1, label2, val2 string) {
	pdf.SetFont("Times", "", 9)
	half := pageW / 2
	l1W := half * 0.3
	v1W := half * 0.7
	l2W := half * 0.3
	v2W := half * 0.7

	pdf.CellFormat(l1W, 6, label1, "LT", 0, "L", false, 0, "")
	pdf.CellFormat(v1W, 6, ": "+val1, "TR", 0, "L", false, 0, "")
	pdf.CellFormat(l2W, 6, label2, "LT", 0, "L", false, 0, "")
	pdf.CellFormat(v2W, 6, ": "+val2, "TR", 1, "L", false, 0, "")

	pdf.CellFormat(l1W, 0, "", "LB", 0, "L", false, 0, "")
	pdf.CellFormat(v1W, 0, "", "RB", 0, "L", false, 0, "")
	pdf.CellFormat(l2W, 0, "", "LB", 0, "L", false, 0, "")
	pdf.CellFormat(v2W, 0, "", "RB", 1, "L", false, 0, "")
}

func drawLabelFull(pdf *gofpdf.Fpdf, pageW float64, label, val string) {
	pdf.SetFont("Times", "", 9)
	lW := pageW * 0.15
	pdf.CellFormat(lW, 6, label, "1", 0, "L", false, 0, "")
	pdf.CellFormat(pageW-lW, 6, ": "+val, "1", 1, "L", false, 0, "")
}

func drawCheckboxRow(pdf *gofpdf.Fpdf, pageW float64, opts []string, selected string) {
	pdf.SetFont("Times", "", 8)
	colW := pageW / 3
	for i, opt := range opts {
		check := "[ ]"
		if strings.Contains(strings.ToLower(selected), strings.ToLower(opt)) {
			check = "[v]"
		}
		pdf.CellFormat(colW, 5, fmt.Sprintf("  %s %s", check, opt), "1", 0, "L", false, 0, "")
		if (i+1)%3 == 0 {
			pdf.Ln(5)
		}
	}
	if len(opts)%3 != 0 {
		pdf.Ln(5)
	}
}

func drawSuratRow(pdf *gofpdf.Fpdf, labelW, colonW, valW float64, label, val string) {
	pdf.CellFormat(labelW, 6, label, "", 0, "L", false, 0, "")
	pdf.CellFormat(colonW, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(valW, 6, val, "", 1, "L", false, 0, "")
}

// â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
// DOCX XML HELPER FUNCTIONS â€” Times New Roman
// â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

// wTimes creates a paragraph with Times New Roman font
func wTimes(text, align, style, sz string) string {
	var rPrParts string
	rPrParts += `<w:rFonts w:ascii="Times New Roman" w:hAnsi="Times New Roman" w:cs="Times New Roman"/>`
	if strings.Contains(style, "B") {
		rPrParts += "<w:b/>"
	}
	if strings.Contains(style, "U") {
		rPrParts += `<w:u w:val="single"/>`
	}
	rPrParts += fmt.Sprintf(`<w:sz w:val="%s"/>`, sz)

	alignTag := ""
	if align == "center" {
		alignTag = `<w:pPr><w:jc w:val="center"/></w:pPr>`
	} else if align == "right" {
		alignTag = `<w:pPr><w:jc w:val="right"/></w:pPr>`
	}

	return fmt.Sprintf(`<w:p>%s<w:r><w:rPr>%s</w:rPr><w:t xml:space="preserve">%s</w:t></w:r></w:p>`,
		alignTag, rPrParts, xmlEscape(text))
}

// wTimesJustify creates a justified paragraph with Times New Roman font (for body text)
func wTimesJustify(text, sz string) string {
	rPrParts := fmt.Sprintf(`<w:rFonts w:ascii="Times New Roman" w:hAnsi="Times New Roman" w:cs="Times New Roman"/><w:sz w:val="%s"/>`, sz)
	alignTag := `<w:pPr><w:jc w:val="both"/></w:pPr>` // "both" = justify in OOXML
	return fmt.Sprintf(`<w:p>%s<w:r><w:rPr>%s</w:rPr><w:t xml:space="preserve">%s</w:t></w:r></w:p>`,
		alignTag, rPrParts, xmlEscape(text))
}

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

func normalisasiJenisCuti(jenis string) string {
	switch jenis {
	case "Cuti Tahunan Biasa":
		return "Cuti Tahunan"
	case "Cuti Tahunan Umroh":
		return "Cuti Tahunan (Umroh)"
	case "Cuti Alasan Penting":
		return "Cuti Karena Alasan Penting"
	default:
		return jenis
	}
}

func toRomawi(n int) string {
	romawi := []string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII"}
	if n >= 1 && n <= 12 {
		return romawi[n]
	}
	return fmt.Sprintf("%d", n)
}


package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"math"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"

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

// ─────────────────────────────────────────────────────────────────────────────
// FORMULIR PERMINTAAN DAN PEMBERIAN CUTI  (PDF only, Times New Roman)
// ─────────────────────────────────────────────────────────────────────────────

func GenerateFormulirCutiPdf(req *models.LeaveRequest) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 10, 15)
	pdf.AddPage()

	now := time.Now()
	pageW := 180.0

	// ── KOP SURAT ──────────────────────────────────────────────────────────
	pdf.SetFont("Times", "B", 12)
	pdf.CellFormat(pageW, 5, "PEMERINTAH KABUPATEN MOROWALI UTARA", "", 1, "C", false, 0, "")
	pdf.CellFormat(pageW, 5, "DINAS PENDIDIKAN DAN KEBUDAYAAN DAERAH", "", 1, "C", false, 0, "")
	pdf.SetFont("Times", "", 9)
	pdf.CellFormat(pageW, 4, "Alamat : Jln. Bumi Nangka Kompleks Perkantoran Kode Pos (94971)", "", 1, "C", false, 0, "")
	pdf.SetFont("Times", "B", 12)
	pdf.CellFormat(pageW, 5, "KOLONODALE", "", 1, "C", false, 0, "")

	pdf.SetLineWidth(0.8)
	pdf.Line(15, pdf.GetY()+1, 195, pdf.GetY()+1)
	pdf.SetLineWidth(0.3)
	pdf.Line(15, pdf.GetY()+3, 195, pdf.GetY()+3)
	pdf.Ln(5)

	// ── JUDUL ──────────────────────────────────────────────────────────────
	pdf.SetFont("Times", "BU", 11)
	pdf.CellFormat(pageW, 7, "FORMULIR PERMINTAAN DAN PEMBERIAN CUTI", "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// ── SECTION I: DATA PEGAWAI ────────────────────────────────────────────
	sectionHeader(pdf, pageW, "I", "Data Pegawai")

	masaKerjaTahun, masaKerjaBulan := hitungMasaKerja(req.Employee.TanggalLahir, now)

	drawLabelValue(pdf, pageW, "Nama", req.Employee.Nama, "NIP", req.Employee.NIP)
	drawLabelValue(pdf, pageW, "Jabatan", req.Employee.Jabatan, "Masa Kerja",
		fmt.Sprintf("%d Tahun %d Bulan", masaKerjaTahun, masaKerjaBulan))
	drawLabelFull(pdf, pageW, "Unit Kerja", req.Employee.TempatTugas)

	// ── SECTION II: JENIS CUTI ─────────────────────────────────────────────
	sectionHeader(pdf, pageW, "II", "Jenis Cuti yang diminta")

	jenisOpts := []string{
		"Cuti Tahunan",
		"Cuti Besar",
		"Cuti Sakit",
		"Cuti Melahirkan",
		"Cuti Karena Alasan Penting",
		"Cuti di Luar Tanggungan Negara",
	}
	drawCheckboxRow(pdf, pageW, jenisOpts, req.JenisCuti)

	// ── SECTION III: ALASAN CUTI ───────────────────────────────────────────
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

	pdf.SetFont("Times", "", 9)
	pdf.CellFormat(pageW, 7, alasan, "1", 1, "L", false, 0, "")

	// ── SECTION IV: LAMA CUTI ──────────────────────────────────────────────
	sectionHeader(pdf, pageW, "IV", "Lama Cuti")

	lamaCuti := int(math.Round(req.TanggalSelesai.Sub(req.TanggalMulai).Hours()/24)) + 1
	tanggalRange := fmt.Sprintf("%s s/d %s", formatTanggalID(req.TanggalMulai), formatTanggalID(req.TanggalSelesai))

	pdf.SetFont("Times", "", 9)
	col1W := pageW * 0.4
	col2W := pageW * 0.2
	col3W := pageW * 0.4
	pdf.CellFormat(col1W, 7, fmt.Sprintf("%d Hari", lamaCuti), "1", 0, "C", false, 0, "")
	pdf.CellFormat(col2W, 7, "Tanggal", "1", 0, "C", false, 0, "")
	pdf.CellFormat(col3W, 7, tanggalRange, "1", 1, "C", false, 0, "")

	// ── SECTION V: CATATAN CUTI ────────────────────────────────────────────
	sectionHeader(pdf, pageW, "V", "Catatan Cuti")

	pdf.SetFont("Times", "B", 8)
	hCol := []float64{pageW * 0.18, pageW * 0.08, pageW * 0.22, pageW * 0.04, pageW * 0.48}
	pdf.CellFormat(hCol[0], 5, "Cuti Tahunan", "1", 0, "C", false, 0, "")
	pdf.CellFormat(hCol[1], 5, "Sisa", "1", 0, "C", false, 0, "")
	pdf.CellFormat(hCol[2], 5, "Keterangan", "1", 0, "C", false, 0, "")
	pdf.CellFormat(hCol[3], 5, "", "1", 0, "C", false, 0, "")
	pdf.CellFormat(hCol[4], 5, "Jenis Cuti Lainnya", "1", 1, "C", false, 0, "")

	rightTypes := []string{"Cuti Tahunan", "Cuti Sakit", "Cuti Karena Alasan Penting", "Cuti Besar", "Cuti Melahirkan", "Cuti di Luar Tanggungan Negara"}
	rightColW := hCol[4] / 2
	years := []string{"N", "N-1", "N-2", "N-3"}
	pdf.SetFont("Times", "", 8)
	for i, yr := range years {
		pdf.CellFormat(hCol[0], 5, yr, "1", 0, "C", false, 0, "")
		pdf.CellFormat(hCol[1], 5, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(hCol[2], 5, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(hCol[3], 5, "", "1", 0, "C", false, 0, "")
		if i < len(rightTypes) {
			pdf.CellFormat(rightColW, 5, rightTypes[i], "1", 0, "L", false, 0, "")
			pdf.CellFormat(rightColW, 5, "", "1", 1, "L", false, 0, "")
		} else {
			pdf.CellFormat(hCol[4], 5, "", "1", 1, "C", false, 0, "")
		}
	}

	// ── SECTION VI: ALAMAT ─────────────────────────────────────────────────
	sectionHeader(pdf, pageW, "VI", "Alamat Selama Menjalankan Cuti")

	pdf.SetFont("Times", "", 9)
	leftW := pageW * 0.55
	rightW := pageW * 0.45

	pdf.CellFormat(leftW, 6, "Alamat : ................................", "1", 0, "L", false, 0, "")
	pdf.CellFormat(rightW, 6, fmt.Sprintf("Kolonodale, %s", formatTanggalID(now)), "1", 1, "L", false, 0, "")

	pdf.CellFormat(leftW, 6, "", "1", 0, "L", false, 0, "")
	pdf.CellFormat(rightW, 6, "Hormat Saya,", "1", 1, "L", false, 0, "")

	pdf.CellFormat(leftW, 6, fmt.Sprintf("Telp : %s", "....................."), "1", 0, "L", false, 0, "")
	pdf.CellFormat(rightW, 18, "", "1", 1, "L", false, 0, "")

	pdf.CellFormat(leftW, 6, "", "LRB", 0, "L", false, 0, "")
	pdf.SetFont("Times", "B", 9)
	pdf.CellFormat(rightW, 6, req.Employee.Nama, "1", 1, "C", false, 0, "")

	pdf.SetFont("Times", "", 8)
	pdf.CellFormat(leftW, 5, "", "LRB", 0, "L", false, 0, "")
	pdf.CellFormat(rightW, 5, fmt.Sprintf("NIP %s", req.Employee.NIP), "1", 1, "C", false, 0, "")

	// ── SECTION VII: PERTIMBANGAN ATASAN ──────────────────────────────────
	sectionHeader(pdf, pageW, "VII", "Pertimbangan Atasan Langsung")

	pdf.SetFont("Times", "", 8)
	colW4 := pageW / 4
	opts := []string{"Disetujui", "Perubahan", "Ditangguhkan", "Tidak Disetujui"}
	for _, o := range opts {
		pdf.CellFormat(colW4, 5, o, "LTR", 0, "C", false, 0, "")
	}
	pdf.Ln(5)
	for range opts {
		pdf.CellFormat(colW4, 12, "", "LBR", 0, "C", false, 0, "")
	}
	pdf.Ln(12)

	pdf.SetFont("Times", "", 9)
	pdf.CellFormat(pageW*0.6, 5, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(pageW*0.4, 5, "Kepala Dinas", "", 1, "C", false, 0, "")
	pdf.Ln(16)
	pdf.SetFont("Times", "BU", 10)
	pdf.CellFormat(pageW*0.6, 5, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(pageW*0.4, 5, "MOH. RIDWAN DM, S.Ag", "", 1, "C", false, 0, "")
	pdf.SetFont("Times", "", 9)
	pdf.CellFormat(pageW*0.6, 4, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(pageW*0.4, 4, "NIP. 19740111 199803 1 004", "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ─────────────────────────────────────────────────────────────────────────────
// SURAT REKOMENDASI IZIN CUTI  (Word + PDF, Times New Roman, Justify body)
// ─────────────────────────────────────────────────────────────────────────────

func GenerateRekomendasiCutiPdf(req *models.LeaveRequest) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(25, 15, 25)
	pdf.AddPage()

	now := time.Now()
	pageW := 160.0

	// ── KOP SURAT ──────────────────────────────────────────────────────────
	pdf.SetFont("Times", "B", 13)
	pdf.CellFormat(pageW, 6, "PEMERINTAH KABUPATEN MOROWALI UTARA", "", 1, "C", false, 0, "")
	pdf.CellFormat(pageW, 6, "DINAS PENDIDIKAN DAN KEBUDAYAAN DAERAH", "", 1, "C", false, 0, "")
	pdf.SetFont("Times", "", 10)
	pdf.CellFormat(pageW, 5, "Alamat : Jln. Bumi Nangka Kompleks Perkantoran Kode Pos (94971)", "", 1, "C", false, 0, "")
	pdf.SetFont("Times", "B", 12)
	pdf.CellFormat(pageW, 5, "KOLONODALE", "", 1, "C", false, 0, "")

	pdf.SetLineWidth(0.8)
	pdf.Line(25, pdf.GetY()+1, 185, pdf.GetY()+1)
	pdf.SetLineWidth(0.3)
	pdf.Line(25, pdf.GetY()+3, 185, pdf.GetY()+3)
	pdf.Ln(8)

	// ── NOMOR / LAMPIRAN / PERIHAL ─────────────────────────────────────────
	pdf.SetFont("Times", "", 12)
	labelW := 25.0
	colonW := 5.0
	valW := pageW - labelW - colonW

	romawi := toRomawi(int(now.Month()))
	nomorSurat := fmt.Sprintf("800.1.11.4/     /Disdikbud /%s/%d", romawi, now.Year())

	drawSuratRow(pdf, labelW, colonW, valW, "Nomor", nomorSurat)
	drawSuratRow(pdf, labelW, colonW, valW, "Lampiran", "Satu Berkas")
	drawSuratRow(pdf, labelW, colonW, valW, "Perihal", "Rekomendasi Izin Cuti")
	pdf.Ln(8)

	// ── TUJUAN ─────────────────────────────────────────────────────────────
	pdf.SetFont("Times", "", 12)
	pdf.CellFormat(pageW, 6, "Yth. Bupati Morowali Utara", "", 1, "L", false, 0, "")
	pdf.CellFormat(pageW, 6, "Cq Kepala Badan Kepegawaian dan", "", 1, "L", false, 0, "")
	pdf.CellFormat(pageW, 6, "Pengembangan SDM", "", 1, "L", false, 0, "")
	pdf.CellFormat(pageW, 6, "Di -", "", 1, "L", false, 0, "")
	pdf.SetX(45)
	pdf.CellFormat(pageW-20, 6, "Tempat", "", 1, "L", false, 0, "")
	pdf.Ln(8)

	// ── ISI SURAT (Times New Roman, Justify) ───────────────────────────────
	pdf.SetFont("Times", "", 12)

	namaJabatan := fmt.Sprintf("%s, %s", req.Employee.Nama, req.Employee.Jabatan)
	tanggalRange := fmt.Sprintf("%s s/d %s", formatTanggalID(req.TanggalMulai), formatTanggalID(req.TanggalSelesai))
	jenisCutiDisplay := normalisasiJenisCuti(req.JenisCuti)

	para1 := fmt.Sprintf(
		"Menindak lanjuti surat permohonan %s atas nama %s; Tanggal %s dengan ini kami tidak keberatan dan menyetujui permohonan tersebut kami teruskan kepada Bapak untuk ditindaklanjuti (Permohonan Terlampir).",
		jenisCutiDisplay, namaJabatan, tanggalRange,
	)

	// "J" = justify (rata kiri-kanan)
	pdf.MultiCell(pageW, 7, para1, "", "J", false)
	pdf.Ln(6)

	pdf.MultiCell(pageW, 7,
		"Demikian Surat Permohonan Cuti ini kami teruskan kepada Bapak, atas pertimbangan Bapak kami ucapkan terima kasih.",
		"", "J", false)
	pdf.Ln(15)

	// ── TANDA TANGAN ───────────────────────────────────────────────────────
	sigOffsetX := 25 + pageW*0.55
	pdf.SetX(sigOffsetX)
	pdf.SetFont("Times", "", 12)
	pdf.CellFormat(pageW*0.45, 6, fmt.Sprintf("Kolonodale, %s", formatTanggalID(now)), "", 1, "L", false, 0, "")
	pdf.SetX(sigOffsetX)
	pdf.CellFormat(pageW*0.45, 6, "Kepala Dinas", "", 1, "L", false, 0, "")
	pdf.Ln(20)

	pdf.SetX(sigOffsetX)
	pdf.SetFont("Times", "BU", 12)
	pdf.CellFormat(pageW*0.45, 6, "MOH. RIDWAN DM, S.Ag", "", 1, "L", false, 0, "")
	pdf.SetX(sigOffsetX)
	pdf.SetFont("Times", "", 11)
	pdf.CellFormat(pageW*0.45, 5, "NIP: 19740111 199803 1 004", "", 1, "L", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ─────────────────────────────────────────────────────────────────────────────
// DOCX — Surat Rekomendasi (Times New Roman, Justify body)
// ─────────────────────────────────────────────────────────────────────────────

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

	contentTypes, _ := zipWriter.Create("[Content_Types].xml")
	_, _ = io.WriteString(contentTypes, `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`)

	rels, _ := zipWriter.Create("_rels/.rels")
	_, _ = io.WriteString(rels, `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`)

	docRels, _ := zipWriter.Create("word/_rels/document.xml.rels")
	_, _ = io.WriteString(docRels, `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"/>`)

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
	namaJabatan := req.Employee.Nama + ", " + req.Employee.Jabatan
	tanggalRange := formatTanggalID(req.TanggalMulai) + " s/d " + formatTanggalID(req.TanggalSelesai)
	romawi := toRomawi(int(now.Month()))
	nomorSurat := fmt.Sprintf("800.1.11.4/     /Disdikbud /%s/%d", romawi, now.Year())

	para1 := fmt.Sprintf(
		"Menindak lanjuti surat permohonan %s atas nama %s; Tanggal %s dengan ini kami tidak keberatan dan menyetujui permohonan tersebut kami teruskan kepada Bapak untuk ditindaklanjuti (Permohonan Terlampir).",
		jenisCutiDisplay, namaJabatan, tanggalRange,
	)
	para2 := "Demikian Surat Permohonan Cuti ini kami teruskan kepada Bapak, atas pertimbangan Bapak kami ucapkan terima kasih."

	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    ` + wTimes("PEMERINTAH KABUPATEN MOROWALI UTARA", "center", "B", "26") + `
    ` + wTimes("DINAS PENDIDIKAN DAN KEBUDAYAAN DAERAH", "center", "B", "26") + `
    ` + wTimes("Alamat : Jln. Bumi Nangka Kompleks Perkantoran Kode Pos (94971)", "center", "", "20") + `
    ` + wTimes("KOLONODALE", "center", "B", "24") + `
    <w:p><w:pPr><w:pBdr><w:bottom w:val="double" w:sz="6" w:space="1" w:color="000000"/></w:pBdr></w:pPr></w:p>
    <w:p/>
    ` + wTimes("Nomor    :  "+nomorSurat, "left", "", "24") + `
    ` + wTimes("Lampiran :  Satu Berkas", "left", "", "24") + `
    ` + wTimes("Perihal  :  Rekomendasi Izin Cuti", "left", "", "24") + `
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
    ` + wTimes("                                                     Kolonodale, "+formatTanggalID(now), "left", "", "24") + `
    ` + wTimes("                                                     Kepala Dinas", "left", "", "24") + `
    <w:p/><w:p/><w:p/>
    ` + wTimes("                                                     MOH. RIDWAN DM, S.Ag", "left", "BU", "24") + `
    ` + wTimes("                                                     NIP: 19740111 199803 1 004", "left", "", "24") + `
  </w:body>
</w:document>`
}

// ─────────────────────────────────────────────────────────────────────────────
// PDF HELPER FUNCTIONS
// ─────────────────────────────────────────────────────────────────────────────

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

// ─────────────────────────────────────────────────────────────────────────────
// DOCX XML HELPER FUNCTIONS — Times New Roman
// ─────────────────────────────────────────────────────────────────────────────

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


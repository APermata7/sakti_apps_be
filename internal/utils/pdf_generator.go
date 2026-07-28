package utils

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf/v2"
)

type PDFData struct {
	CompanyName    string
	CompanyAddress string
	CompanyPhone   string
	CompanyEmail   string
	LogoPath       string

	TTDKaryawanURL string
	TTDAtasanURL   string
	TTDHRDURL      string

	Jenis string

	Nama        string
	Divisi      string
	Unit        string
	Jabatan     string
	LamaCuti    int
	TanggalCuti string
	AlasanCuti  string

	JumlahCutiTahun      int
	CutiDilaksanakan     int
	CutiAkanDilaksanakan int
	SisaCuti             int

	DisetujuiSelama int
	MulaiTanggal    string
	Tahun           int
	TanggalSekarang time.Time

	NamaPemohon string
	NamaAtasan  string
	NamaHRD     string

	Status  string
	Catatan string

	JumlahTTD int
}

func formatTanggalIndonesia(t time.Time) string {
	bulanIndonesia := map[time.Month]string{
		time.January:   "Januari",
		time.February:  "Februari",
		time.March:     "Maret",
		time.April:     "April",
		time.May:       "Mei",
		time.June:      "Juni",
		time.July:      "Juli",
		time.August:    "Agustus",
		time.September: "September",
		time.October:   "Oktober",
		time.November:  "November",
		time.December:  "Desember",
	}
	return fmt.Sprintf("%d %s %d", t.Day(), bulanIndonesia[t.Month()], t.Year())
}

func formatTanggalCuti(tanggalCuti string) string {
	if tanggalCuti == "" {
		return ""
	}

	parts := strings.Split(tanggalCuti, " s.d ")
	if len(parts) == 2 {
		tglMulai := formatTanggalSingle(parts[0])
		tglSelesai := formatTanggalSingle(parts[1])
		return tglMulai + " s.d " + tglSelesai
	}

	return formatTanggalSingle(tanggalCuti)
}

func formatTanggalSingle(tanggal string) string {
	if tanggal == "" {
		return ""
	}

	t, err := time.Parse("2006-01-02", tanggal)
	if err != nil {
		return tanggal
	}

	return formatTanggalIndonesia(t)
}

func cleanAlasanCuti(alasan string) string {
	if strings.Contains(alasan, "| Catatan HRD:") {
		parts := strings.Split(alasan, "| Catatan HRD:")
		return strings.TrimSpace(parts[0])
	}
	return alasan
}

func addHeader(pdf *gofpdf.Fpdf, data PDFData) {
	logoSize := 24.0
	startY := 5.0

	if data.LogoPath != "" {
		pdf.Image(data.LogoPath, 18, startY, logoSize, logoSize, false, "", 0, "")
	}

	textStartY := startY + 4.0

	textX := 15.0 + logoSize + 7

	pdf.SetXY(textX, textStartY)
	pdf.SetFont("Helvetica", "B", 24)
	pdf.SetTextColor(0, 0, 100)
	pdf.CellFormat(105, 9, "KOPEGTEL MALANG", "", 0, "L", false, 0, "")

	rightX := 150.0
	pdf.SetXY(rightX, textStartY)
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(0, 0, 100)
	pdf.CellFormat(45, 4.5, "Badan Hukum :", "", 1, "L", false, 0, "")
	pdf.SetXY(rightX, textStartY+4.5)
	pdf.CellFormat(45, 4.5, "5538/BH/II/1983", "", 1, "L", false, 0, "")
	pdf.SetXY(rightX, textStartY+9)
	pdf.CellFormat(45, 4.5, "Tgl. 28 September 1983", "", 1, "L", false, 0, "")

	pdf.SetXY(textX, textStartY+10)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.SetTextColor(0, 0, 100)
	pdf.CellFormat(105, 8, "Koperasi Pegawai PT. Telkom Malang", "", 1, "L", false, 0, "")

	pdf.SetY(34)
	pdf.SetTextColor(0, 0, 0)

	pdf.SetDrawColor(255, 165, 0)
	pdf.SetLineWidth(0.8)
	pdf.Line(18, 34, 192, 34)
	pdf.SetLineWidth(0.3)

	pdf.Ln(8)
}

func addFooter(pdf *gofpdf.Fpdf) {
	pdf.SetY(-32)

	pdf.SetDrawColor(255, 165, 0)
	pdf.SetLineWidth(0.8)
	pdf.Line(18, pdf.GetY(), 192, pdf.GetY())
	pdf.SetLineWidth(0.3)

	pdf.SetY(-30)
	pdf.SetFont("Helvetica", "B", 11)
	pdf.SetTextColor(0, 0, 100)
	pdf.CellFormat(0, 4, "Alamat: Jl. Ahmad Yani 11 Malang - Jawa Timur | Telepon: 0341 479881 | Faximile: 0341 499553", "", 1, "C", false, 0, "")
	pdf.SetY(-25)
	pdf.SetFont("Helvetica", "B", 11)
	pdf.CellFormat(0, 4, "Email: kopegtel1malang@gmail.com", "", 1, "C", false, 0, "")
}

func generateBodyCuti(pdf *gofpdf.Fpdf, data PDFData) {
	pdf.SetY(42)
	pdf.SetFont("Helvetica", "B", 15)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 8, "PERMOHONAN / LAPORAN CUTI TAHUNAN", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 0.5, "", "", 1, "", false, 0, "")

	pdf.Ln(6)

	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(0, 6.5, "DATA PEGAWAI", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	leftLabel := 30.0
	leftValue := 55.0
	lineHeight := 5.5
	indent := 8.0

	alasanBersih := cleanAlasanCuti(data.AlasanCuti)
	tanggalCutiFormatted := formatTanggalCuti(data.TanggalCuti)

	pdf.SetFont("Helvetica", "", 12)
	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "1) Nama")
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(leftValue, lineHeight, ": "+data.Nama)
	pdf.Ln(lineHeight)

	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "2) Divisi")
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(leftValue, lineHeight, ": "+data.Divisi)
	pdf.Ln(lineHeight)

	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "3) Unit")
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(leftValue, lineHeight, ": "+data.Unit)
	pdf.Ln(lineHeight)

	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "4) Jabatan")
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(leftValue, lineHeight, ": "+data.Jabatan)
	pdf.Ln(lineHeight)

	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "5) Lama Cuti")
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(leftValue, lineHeight, fmt.Sprintf(": %d hari", data.LamaCuti))
	pdf.Ln(lineHeight)

	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "6) Tanggal Cuti")
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(leftValue, lineHeight, ": "+tanggalCutiFormatted)
	pdf.Ln(lineHeight)

	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "7) Alasan Cuti")
	pdf.SetFont("Helvetica", "", 12)
	pdf.MultiCell(0, lineHeight, ": "+alasanBersih, "", "L", false)
	pdf.Ln(6)

	if data.Jenis == "CUTI" {
		pdf.SetFont("Helvetica", "B", 12)
		pdf.CellFormat(0, 6.5, "CATATAN SDM KOPEGTEL", "", 1, "L", false, 0, "")
		pdf.Ln(2)

		pdf.SetFont("Helvetica", "", 12)
		pdf.SetX(18)

		indent := 8.0
		labelWidth := 85.0
		valueWidth := 70.0

		pdf.SetX(18 + indent)
		pdf.Cell(labelWidth, lineHeight, fmt.Sprintf("- Jumlah cuti pada tahun %d", data.Tahun))
		pdf.SetX(18 + indent + labelWidth)
		pdf.CellFormat(valueWidth, lineHeight, fmt.Sprintf(": %d hari", data.JumlahCutiTahun), "", 1, "L", false, 0, "")

		pdf.SetX(18 + indent)
		pdf.Cell(labelWidth, lineHeight, fmt.Sprintf("- Cuti yang telah dilaksanakan tahun %d", data.Tahun))
		pdf.SetX(18 + indent + labelWidth)
		pdf.CellFormat(valueWidth, lineHeight, fmt.Sprintf(": %d hari", data.CutiDilaksanakan), "", 1, "L", false, 0, "")

		pdf.SetX(18 + indent)
		pdf.Cell(labelWidth, lineHeight, fmt.Sprintf("- Cuti yang akan dilaksanakan tahun %d", data.Tahun))
		pdf.SetX(18 + indent + labelWidth)
		pdf.CellFormat(valueWidth, lineHeight, fmt.Sprintf(": %d hari", data.CutiAkanDilaksanakan), "", 1, "L", false, 0, "")

		pdf.SetX(18 + indent)
		pdf.Cell(labelWidth, lineHeight, fmt.Sprintf("- Sisa cuti tahun %d", data.Tahun))
		pdf.SetX(18 + indent + labelWidth)
		pdf.CellFormat(valueWidth, lineHeight, fmt.Sprintf(": %d hari", data.SisaCuti), "", 1, "L", false, 0, "")

		pdf.Ln(6)
	}
}

func generateBodyDispen(pdf *gofpdf.Fpdf, data PDFData) {
	pdf.SetY(42)
	pdf.SetFont("Helvetica", "B", 15)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 8, "PERMOHONAN / LAPORAN DISPENSASI", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 0.5, "", "", 1, "", false, 0, "")

	pdf.Ln(6)

	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(0, 6.5, "DATA PEGAWAI", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	leftLabel := 30.0
	leftValue := 55.0
	lineHeight := 6.0
	indent := 8.0

	alasanBersih := cleanAlasanCuti(data.AlasanCuti)
	tanggalCutiFormatted := formatTanggalCuti(data.TanggalCuti)

	pdf.SetFont("Helvetica", "", 12)
	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "1) Nama")
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(leftValue, lineHeight, ": "+data.Nama)
	pdf.Ln(lineHeight)

	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "2) Divisi")
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(leftValue, lineHeight, ": "+data.Divisi)
	pdf.Ln(lineHeight)

	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "3) Unit")
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(leftValue, lineHeight, ": "+data.Unit)
	pdf.Ln(lineHeight)

	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "4) Jabatan")
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(leftValue, lineHeight, ": "+data.Jabatan)
	pdf.Ln(lineHeight)

	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "5) Lama Cuti")
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(leftValue, lineHeight, fmt.Sprintf(": %d hari", data.LamaCuti))
	pdf.Ln(lineHeight)

	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "6) Tanggal Cuti")
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(leftValue, lineHeight, ": "+tanggalCutiFormatted)
	pdf.Ln(lineHeight)

	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "7) Alasan Cuti")
	pdf.SetFont("Helvetica", "", 12)
	pdf.MultiCell(0, lineHeight, ": "+alasanBersih, "", "L", false)
	pdf.Ln(6)
}

func addKeputusan2TTD(pdf *gofpdf.Fpdf, data PDFData) {
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(0, 6.5, "KEPUTUSAN PEJABAT YANG BERWENANG", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	leftLabel := 35.0
	leftValue := 55.0
	lineHeight := 6.0
	indent := 8.0

	mulaiTanggalFormatted := formatTanggalSingle(data.MulaiTanggal)

	pdf.SetFont("Helvetica", "", 12)
	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "Disetujui selama")
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(leftValue, lineHeight, fmt.Sprintf(": %d hari kerja", data.DisetujuiSelama))
	pdf.Ln(lineHeight)

	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "Mulai tanggal")
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(leftValue, lineHeight, ": "+mulaiTanggalFormatted)
	pdf.Ln(lineHeight + 12)

	pdf.SetX(145)
	pdf.CellFormat(0, 6.5, fmt.Sprintf("Malang, %s", formatTanggalIndonesia(data.TanggalSekarang)), "", 1, "L", false, 0, "")
	pdf.Ln(4)

	lebarKolom := 55.0
	jarakKolom := 45.0
	marginKanan := 25.0
	totalLebar := (lebarKolom * 2) + jarakKolom
	xAwal := 210.0 - marginKanan - totalLebar

	x := xAwal
	y := pdf.GetY()

	pdf.SetDrawColor(0, 0, 0)

	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetXY(x, y)
	pdf.CellFormat(lebarKolom, 6.5, "Mengetahui HRD", "", 0, "C", false, 0, "")
	pdf.SetXY(x, y+18)
	pdf.CellFormat(lebarKolom, 20, "", "B", 0, "C", false, 0, "")
	if data.TTDHRDURL != "" {
		pdf.Image(data.TTDHRDURL, x+2, y+15, lebarKolom-4, 20, false, "", 0, "")
	}
	pdf.SetXY(x, y+40)
	pdf.SetFont("Helvetica", "", 12)
	pdf.CellFormat(lebarKolom, 6, data.NamaHRD, "", 0, "C", false, 0, "")

	x += lebarKolom + jarakKolom
	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetXY(x, y)
	pdf.CellFormat(lebarKolom, 6.5, "Pemohon", "", 0, "C", false, 0, "")
	pdf.SetXY(x, y+18)
	pdf.CellFormat(lebarKolom, 20, "", "B", 0, "C", false, 0, "")
	if data.TTDKaryawanURL != "" {
		pdf.Image(data.TTDKaryawanURL, x+2, y+15, lebarKolom-4, 20, false, "", 0, "")
	}
	pdf.SetXY(x, y+40)
	pdf.SetFont("Helvetica", "", 12)
	pdf.CellFormat(lebarKolom, 6, data.NamaPemohon, "", 0, "C", false, 0, "")

	pdf.SetY(y + 52)
}

func addKeputusan3TTD(pdf *gofpdf.Fpdf, data PDFData) {
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(0, 6.5, "KEPUTUSAN PEJABAT YANG BERWENANG", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	leftLabel := 35.0
	leftValue := 55.0
	lineHeight := 6.0
	indent := 8.0

	mulaiTanggalFormatted := formatTanggalSingle(data.MulaiTanggal)

	pdf.SetFont("Helvetica", "", 12)
	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "Disetujui selama")
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(leftValue, lineHeight, fmt.Sprintf(": %d hari kerja", data.DisetujuiSelama))
	pdf.Ln(lineHeight)

	pdf.SetXY(18+indent, pdf.GetY())
	pdf.Cell(leftLabel, lineHeight, "Mulai tanggal")
	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(leftValue, lineHeight, ": "+mulaiTanggalFormatted)
	pdf.Ln(lineHeight + 12)

	pdf.SetX(145)
	pdf.CellFormat(0, 6.5, fmt.Sprintf("Malang, %s", formatTanggalIndonesia(data.TanggalSekarang)), "", 1, "L", false, 0, "")
	pdf.Ln(4)

	lebarKolom := 50.0
	jarakKolom := 30.0
	totalLebar := (lebarKolom * 3) + (jarakKolom * 2)
	xAwal := (210.0 - totalLebar) / 2

	x := xAwal
	y := pdf.GetY()

	pdf.SetDrawColor(0, 0, 0)

	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetXY(x, y)
	pdf.CellFormat(lebarKolom, 6.5, "Mengetahui HRD", "", 0, "C", false, 0, "")
	pdf.SetXY(x, y+18)
	pdf.CellFormat(lebarKolom, 20, "", "B", 0, "C", false, 0, "")
	if data.TTDHRDURL != "" {
		pdf.Image(data.TTDHRDURL, x+2, y+15, lebarKolom-4, 20, false, "", 0, "")
	}
	pdf.SetXY(x, y+40)
	pdf.SetFont("Helvetica", "", 12)
	pdf.CellFormat(lebarKolom, 6, data.NamaHRD, "", 0, "C", false, 0, "")

	x += lebarKolom + jarakKolom
	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetXY(x, y)
	pdf.CellFormat(lebarKolom, 6.5, "Menyetujui Atasan", "", 0, "C", false, 0, "")
	pdf.SetXY(x, y+18)
	pdf.CellFormat(lebarKolom, 20, "", "B", 0, "C", false, 0, "")
	if data.TTDAtasanURL != "" {
		pdf.Image(data.TTDAtasanURL, x+2, y+15, lebarKolom-4, 20, false, "", 0, "")
	}
	pdf.SetXY(x, y+40)
	pdf.SetFont("Helvetica", "", 12)
	pdf.CellFormat(lebarKolom, 6, data.NamaAtasan, "", 0, "C", false, 0, "")

	x += lebarKolom + jarakKolom
	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetXY(x, y)
	pdf.CellFormat(lebarKolom, 6.5, "Pemohon", "", 0, "C", false, 0, "")
	pdf.SetXY(x, y+18)
	pdf.CellFormat(lebarKolom, 20, "", "B", 0, "C", false, 0, "")
	if data.TTDKaryawanURL != "" {
		pdf.Image(data.TTDKaryawanURL, x+2, y+15, lebarKolom-4, 20, false, "", 0, "")
	}
	pdf.SetXY(x, y+40)
	pdf.SetFont("Helvetica", "", 12)
	pdf.CellFormat(lebarKolom, 6, data.NamaPemohon, "", 0, "C", false, 0, "")

	pdf.SetY(y + 52)
}

func addCatatan(pdf *gofpdf.Fpdf, data PDFData) {
	pdf.SetFont("Helvetica", "I", 10)
	pdf.CellFormat(0, 5, "Catatan:", "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "I", 8)
	if data.Jenis == "CUTI" {
		pdf.CellFormat(0, 4.5, "1. Form cuti bisa dicetak di masing-masing unit.", "", 1, "L", false, 0, "")
		pdf.CellFormat(0, 4.5, "2. Cuti yang tidak diketahui HRD KOPEGTEL Malang maka dianggap mangkir.", "", 1, "L", false, 0, "")
	} else {
		pdf.CellFormat(0, 4.5, "1. Form dispensasi bisa dicetak di masing-masing unit.", "", 1, "L", false, 0, "")
		pdf.CellFormat(0, 4.5, "2. Dispensasi yang tidak diketahui HRD KOPEGTEL Malang maka dianggap mangkir.", "", 1, "L", false, 0, "")
	}
}

func GeneratePDFCuti(data PDFData) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Surat Cuti Tahunan", true)
	pdf.SetAuthor("Sakti Apps - KOPEGTEL", true)
	pdf.SetMargins(18, 18, 18)
	pdf.AddPage()

	addHeader(pdf, data)
	generateBodyCuti(pdf, data)
	if data.JumlahTTD == 2 {
		addKeputusan2TTD(pdf, data)
	} else {
		addKeputusan3TTD(pdf, data)
	}
	addCatatan(pdf, data)
	addFooter(pdf)

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GeneratePDFDispensasi(data PDFData) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Surat Dispensasi", true)
	pdf.SetAuthor("Sakti Apps - KOPEGTEL", true)
	pdf.SetMargins(18, 18, 18)
	pdf.AddPage()

	addHeader(pdf, data)
	generateBodyDispen(pdf, data)
	if data.JumlahTTD == 2 {
		addKeputusan2TTD(pdf, data)
	} else {
		addKeputusan3TTD(pdf, data)
	}
	addCatatan(pdf, data)
	addFooter(pdf)

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GeneratePDF(data PDFData) ([]byte, error) {
	switch data.Jenis {
	case "DISPENSASI":
		return GeneratePDFDispensasi(data)
	default:
		return GeneratePDFCuti(data)
	}
}
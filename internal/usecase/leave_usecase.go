package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
	"sakti_apps_be/internal/utils"
)

type LeaveUsecase struct {
	LeaveRepo    *repository.LeaveRepo
	KaryawanRepo *repository.KaryawanRepo
	TTDRepo      *repository.TTDRepo
	ConfigRepo   *repository.KonfigurasiRepo
}

func NewLeaveUsecase(
	leaveRepo *repository.LeaveRepo,
	karyawanRepo *repository.KaryawanRepo,
	ttdRepo *repository.TTDRepo,
	configRepo *repository.KonfigurasiRepo,
) *LeaveUsecase {
	return &LeaveUsecase{
		LeaveRepo:    leaveRepo,
		KaryawanRepo: karyawanRepo,
		TTDRepo:      ttdRepo,
		ConfigRepo:   configRepo,
	}
}

func (u *LeaveUsecase) CreateLeave(ctx context.Context, karyawanID string, req domain.CreateCutiRequest) (*domain.PengajuanCuti, error) {
	karyawan, err := u.KaryawanRepo.GetByID(ctx, karyawanID)
	if err != nil || karyawan == nil {
		return nil, errors.New("karyawan tidak ditemukan")
	}

	if req.SubTipe != "izin" && req.SubTipe != "sakit" && req.SubTipe != "dispensasi" {
		return nil, errors.New("sub tipe harus 'izin', 'sakit', atau 'dispensasi'")
	}

	start, err := time.Parse("2006-01-02", req.TanggalMulai)
	if err != nil {
		return nil, errors.New("format tanggal mulai tidak valid (YYYY-MM-DD)")
	}
	end, err := time.Parse("2006-01-02", req.TanggalSelesai)
	if err != nil {
		return nil, errors.New("format tanggal selesai tidak valid (YYYY-MM-DD)")
	}
	totalHari := int(end.Sub(start).Hours()/24) + 1

	if totalHari <= 0 {
		return nil, errors.New("tanggal tidak valid")
	}

	if req.SubTipe == "dispensasi" && totalHari > 2 {
		return nil, errors.New("dispensasi maksimal 2 hari")
	}

	existingLeaves, err := u.LeaveRepo.GetActiveLeaves(ctx, karyawanID)
	if err != nil {
		return nil, err
	}

	for _, existing := range existingLeaves {
		existingStart := existing.TanggalMulai
		existingEnd := existing.TanggalSelesai

		if !(end.Before(existingStart) || start.After(existingEnd)) {
			return nil, errors.New("terdapat pengajuan cuti yang overlap pada tanggal tersebut")
		}
	}

	balance, err := u.LeaveRepo.GetBalance(ctx, karyawanID, time.Now().Year())
	if err != nil {
		return nil, errors.New("gagal mendapatkan kuota cuti")
	}

	if balance == nil {
		balance = &domain.SisaCuti{
			JumlahCuti:        12,
			TelahDilaksanakan: 0,
			AkanDilaksanakan:  0,
			SisaCuti:          12,
		}
	}

	mengurangiCuti := true
	if req.SubTipe == "dispensasi" {
		mengurangiCuti = false
	}

	if mengurangiCuti && totalHari > balance.SisaCuti {
		return nil, errors.New("kuota cuti tidak mencukupi")
	}

	langsungApprove := req.LangsungApprove
	langsungFinal := false

	if req.SubTipe == "dispensasi" {
		langsungApprove = true
		langsungFinal = true
	}

	judulDokumen := "PERMOHONAN/LAPORAN CUTI TAHUNAN"
	if req.SubTipe == "dispensasi" {
		judulDokumen = "PERMOHONAN/LAPORAN DISPENSASI"
	}

	status := "menunggu"
	if langsungApprove {
		status = "disetujui"
	}

	leave := &domain.PengajuanCuti{
		KaryawanID:        karyawanID,
		SubTipe:           req.SubTipe,
		TanggalMulai:      start,
		TanggalSelesai:    end,
		TotalHari:         totalHari,
		Alasan:            req.Alasan,
		Status:            status,
		BackDate:          req.BackDate,
		MengurangiCuti:    mengurangiCuti,
		LangsungApprove:   langsungApprove,
		LangsungFinal:     langsungFinal,
		JudulDokumen:      judulDokumen,
	}

	if req.SubTipe == "dispensasi" {
		hrd, err := u.KaryawanRepo.GetByRole(ctx, "hrd")
		if err == nil && hrd != nil {
			leave.DifinalisasiOleh = &hrd.ID
			now := time.Now()
			leave.TanggalDifinalisasi = &now
		}
	}

	if err := u.LeaveRepo.Create(ctx, leave); err != nil {
		return nil, err
	}

	return u.LeaveRepo.GetByID(ctx, leave.ID)
}

func (u *LeaveUsecase) GetStatus(ctx context.Context, karyawanID, status string, limit, page int) ([]domain.PengajuanCuti, int, error) {
	offset := (page - 1) * limit
	return u.LeaveRepo.GetByKaryawanID(ctx, karyawanID, status, limit, offset)
}

func (u *LeaveUsecase) CancelLeave(ctx context.Context, leaveID, karyawanID string) (*domain.PengajuanCuti, error) {
	leave, err := u.LeaveRepo.GetByID(ctx, leaveID)
	if err != nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if leave == nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if leave.KaryawanID != karyawanID {
		return nil, errors.New("anda tidak memiliki akses")
	}

	if leave.SubTipe != "dispensasi" {
		if leave.DifinalisasiOleh == nil || leave.TanggalDifinalisasi == nil {
			return nil, errors.New("pengajuan belum difinalisasi HRD, tidak dapat dibatalkan")
		}
	} else {
		if leave.TanggalDifinalisasi == nil {
			return nil, errors.New("pengajuan belum difinalisasi, tidak dapat dibatalkan")
		}
	}

	if leave.Status != "disetujui" {
		return nil, errors.New("pengajuan tidak bisa dibatalkan")
	}

	if err := u.LeaveRepo.UpdateStatus(ctx, leaveID, "dibatalkan"); err != nil {
		return nil, err
	}

	return u.LeaveRepo.GetByID(ctx, leaveID)
}

func (u *LeaveUsecase) ApproveLeave(ctx context.Context, leaveID, managerID string) (*domain.PengajuanCuti, error) {
	leave, err := u.LeaveRepo.GetByID(ctx, leaveID)
	if err != nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if leave == nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if leave.Status != "menunggu" {
		return nil, errors.New("pengajuan sudah diproses")
	}

	if err := u.LeaveRepo.Approve(ctx, leaveID, managerID); err != nil {
		return nil, err
	}

	return u.LeaveRepo.GetByID(ctx, leaveID)
}

func (u *LeaveUsecase) RejectLeave(ctx context.Context, leaveID, managerID, alasan string) (*domain.PengajuanCuti, error) {
	leave, err := u.LeaveRepo.GetByID(ctx, leaveID)
	if err != nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if leave == nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if leave.Status != "menunggu" {
		return nil, errors.New("pengajuan sudah diproses")
	}

	if err := u.LeaveRepo.Reject(ctx, leaveID, managerID, alasan); err != nil {
		return nil, err
	}

	return u.LeaveRepo.GetByID(ctx, leaveID)
}

func (u *LeaveUsecase) FinalizeLeave(ctx context.Context, leaveID, hrdID, catatan string) (*domain.PengajuanCuti, error) {
	leave, err := u.LeaveRepo.GetByID(ctx, leaveID)
	if err != nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if leave == nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if leave.Status != "disetujui" {
		return nil, errors.New("pengajuan belum disetujui atasan")
	}

	if err := u.LeaveRepo.Finalize(ctx, leaveID, hrdID, catatan); err != nil {
		return nil, err
	}

	return u.LeaveRepo.GetByID(ctx, leaveID)
}

func (u *LeaveUsecase) DownloadSuratCuti(ctx context.Context, leaveID string) ([]byte, string, error) {
	leave, err := u.LeaveRepo.GetByID(ctx, leaveID)
	if err != nil {
		return nil, "", err
	}
	if leave == nil {
		return nil, "", errors.New("pengajuan tidak ditemukan")
	}

	if leave.Status != "disetujui" && leave.Status != "dibatalkan" {
		return nil, "", errors.New("pengajuan belum disetujui")
	}

	karyawan, err := u.KaryawanRepo.GetByID(ctx, leave.KaryawanID)
	if err != nil {
		return nil, "", err
	}
	if karyawan == nil {
		return nil, "", errors.New("karyawan tidak ditemukan")
	}

	var atasan *domain.Karyawan
	if leave.DisetujuiOleh != nil {
		atasan, _ = u.KaryawanRepo.GetByID(ctx, *leave.DisetujuiOleh)
	}

	if leave.SubTipe == "dispensasi" && atasan == nil {
		if karyawan.AtasanLangsungID != nil {
			atasan, _ = u.KaryawanRepo.GetByID(ctx, *karyawan.AtasanLangsungID)
		}
	}

	if leave.SubTipe == "dispensasi" && atasan == nil {
		atasanList, _, err := u.KaryawanRepo.GetAll(ctx, 10, 0, "", "atasan", "aktif")
		if err == nil && len(atasanList) > 0 {
			for _, calonAtasan := range atasanList {
				if calonAtasan.ID != leave.KaryawanID {
					atasan = &calonAtasan
					break
				}
			}
		}
	}

	if leave.SubTipe == "dispensasi" && atasan == nil {
		allKaryawan, _, err := u.KaryawanRepo.GetAll(ctx, 100, 0, "", "", "")
		if err == nil && len(allKaryawan) > 0 {
			for _, calonAtasan := range allKaryawan {
				if calonAtasan.ID != leave.KaryawanID && calonAtasan.Role != "hrd" && calonAtasan.Role != "admin" {
					atasan = &calonAtasan
					break
				}
			}
		}
	}

	var hrd *domain.Karyawan
	if leave.DifinalisasiOleh != nil {
		hrd, _ = u.KaryawanRepo.GetByID(ctx, *leave.DifinalisasiOleh)
	}

	if hrd == nil {
		hrdList, _, err := u.KaryawanRepo.GetAll(ctx, 10, 0, "", "hrd", "aktif")
		if err == nil && len(hrdList) > 0 {
			hrd = &hrdList[0]
		}
	}

	var ttdKaryawanURL string
	ttdKaryawan, err := u.TTDRepo.GetByKaryawanID(ctx, leave.KaryawanID)
	if err == nil && ttdKaryawan != nil {
		ttdKaryawanURL = ttdKaryawan.URLTandaTangan
	}

	var ttdAtasanURL string
	if leave.DisetujuiOleh != nil {
		ttdAtasan, err := u.TTDRepo.GetByKaryawanID(ctx, *leave.DisetujuiOleh)
		if err == nil && ttdAtasan != nil {
			ttdAtasanURL = ttdAtasan.URLTandaTangan
		}
	}

	if leave.SubTipe == "dispensasi" && ttdAtasanURL == "" {
		if karyawan.AtasanLangsungID != nil {
			ttdAtasan, err := u.TTDRepo.GetByKaryawanID(ctx, *karyawan.AtasanLangsungID)
			if err == nil && ttdAtasan != nil {
				ttdAtasanURL = ttdAtasan.URLTandaTangan
			}
		}
	}

	if leave.SubTipe == "dispensasi" && ttdAtasanURL == "" {
		if atasan != nil {
			ttdAtasan, err := u.TTDRepo.GetByKaryawanID(ctx, atasan.ID)
			if err == nil && ttdAtasan != nil {
				ttdAtasanURL = ttdAtasan.URLTandaTangan
			}
		}
	}

	var ttdHRDURL string
	if leave.DifinalisasiOleh != nil {
		ttdHRD, err := u.TTDRepo.GetByKaryawanID(ctx, *leave.DifinalisasiOleh)
		if err == nil && ttdHRD != nil {
			ttdHRDURL = ttdHRD.URLTandaTangan
		}
	}

	if ttdHRDURL == "" && hrd != nil {
		ttdHRD, err := u.TTDRepo.GetByKaryawanID(ctx, hrd.ID)
		if err == nil && ttdHRD != nil {
			ttdHRDURL = ttdHRD.URLTandaTangan
		}
	}

	config, err := u.ConfigRepo.GetActive(ctx)
	if err != nil {
		config = &domain.KonfigurasiKerja{}
	}
	logoURL := ""
	if config.LogoKantor != nil {
		logoURL = *config.LogoKantor
	}

	balance, err := u.LeaveRepo.GetBalance(ctx, leave.KaryawanID, time.Now().Year())
	if err != nil {
		balance = &domain.SisaCuti{
			JumlahCuti:        12,
			TelahDilaksanakan: 0,
			AkanDilaksanakan:  0,
			SisaCuti:          12,
		}
	}

	jenisPDF := "CUTI"
	if leave.SubTipe == "dispensasi" {
		jenisPDF = "DISPENSASI"
	}

	namaAtasan := ""
	if atasan != nil {
		namaAtasan = atasan.NamaLengkap
	}
	namaHRD := ""
	if hrd != nil {
		namaHRD = hrd.NamaLengkap
	}

	divisi := ""
	if karyawan.Divisi != nil {
		divisi = *karyawan.Divisi
	}
	unit := ""
	if karyawan.Unit != nil {
		unit = *karyawan.Unit
	}
	levelJabatan := ""
	if karyawan.LevelJabatan != nil {
		levelJabatan = *karyawan.LevelJabatan
	}

	cutiDilaksanakan := balance.TelahDilaksanakan - leave.TotalHari
	if cutiDilaksanakan < 0 {
		cutiDilaksanakan = 0
	}

	jumlahTTD := 3
	if hrd == nil && atasan == nil {
		jumlahTTD = 2
	} else if hrd != nil && atasan == nil {
		jumlahTTD = 2
	}

	tanggalMulai := leave.TanggalMulai.Format("2006-01-02")
	tanggalSelesai := leave.TanggalSelesai.Format("2006-01-02")

	logoPath, err := utils.DownloadImage(logoURL)
	if err != nil {
		log.Printf("Gagal download logo: %v", err)
		logoPath = ""
	}
	defer utils.CleanupTempFiles(logoPath)

	ttdKaryawanPath, err := utils.DownloadImage(ttdKaryawanURL)
	if err != nil {
		log.Printf("Gagal download TTD karyawan: %v", err)
		ttdKaryawanPath = ""
	}
	defer utils.CleanupTempFiles(ttdKaryawanPath)

	ttdAtasanPath, err := utils.DownloadImage(ttdAtasanURL)
	if err != nil {
		log.Printf("Gagal download TTD atasan: %v", err)
		ttdAtasanPath = ""
	}
	defer utils.CleanupTempFiles(ttdAtasanPath)

	ttdHRDPath, err := utils.DownloadImage(ttdHRDURL)
	if err != nil {
		log.Printf("Gagal download TTD HRD: %v", err)
		ttdHRDPath = ""
	}
	defer utils.CleanupTempFiles(ttdHRDPath)

	pdfData := utils.PDFData{
		CompanyName:    "KOPEGTEL MALANG",
		CompanyAddress: "Jl. Ahmad Yani No.11, Blimbing, Kota Malang",
		LogoPath:       logoPath,
		TTDKaryawanURL: ttdKaryawanPath,
		TTDAtasanURL:   ttdAtasanPath,
		TTDHRDURL:      ttdHRDPath,

		Jenis:       jenisPDF,
		Nama:        karyawan.NamaLengkap,
		Divisi:      divisi,
		Unit:        unit,
		Jabatan:     levelJabatan,
		LamaCuti:    leave.TotalHari,
		TanggalCuti: tanggalMulai + " s.d " + tanggalSelesai,
		AlasanCuti:  leave.Alasan,

		JumlahCutiTahun:      12,
		CutiDilaksanakan:     cutiDilaksanakan,
		CutiAkanDilaksanakan: leave.TotalHari,
		SisaCuti:             balance.SisaCuti,

		DisetujuiSelama: leave.TotalHari,
		MulaiTanggal:    tanggalMulai,
		Tahun:           time.Now().Year(),
		TanggalSekarang: time.Now(),

		NamaPemohon: karyawan.NamaLengkap,
		NamaAtasan:  namaAtasan,
		NamaHRD:     namaHRD,

		JumlahTTD: jumlahTTD,
	}

	pdfBytes, err := utils.GeneratePDF(pdfData)
	if err != nil {
		return nil, "", err
	}

	filename := "surat_" + jenisPDF + "_" + karyawan.NamaLengkap + "_" + time.Now().Format("20060102") + ".pdf"

	pdfURL, err := utils.UploadPDF(pdfBytes, filename)
	if err != nil {
		log.Printf("Gagal upload PDF ke Cloudinary: %v", err)
	} else {
		err = u.LeaveRepo.UpdatePDFURL(ctx, leaveID, pdfURL)
		if err != nil {
			log.Printf("Gagal simpan URL PDF: %v", err)
		}
	}

	return pdfBytes, filename, nil
}
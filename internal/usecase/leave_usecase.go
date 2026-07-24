package usecase

import (
	"context"
	"errors"
	"log"
	"strconv"
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
	RiwayatRepo  *repository.RiwayatRepo
}

type ApprovalFlow struct {
	LangsungApprove bool
	LangsungFinal   bool
	ButuhAtasan     bool
	ButuhHRD        bool
	JumlahTTD       int
	AtasanID        *string
	HRDID           *string
}

func NewLeaveUsecase(
	leaveRepo *repository.LeaveRepo,
	karyawanRepo *repository.KaryawanRepo,
	ttdRepo *repository.TTDRepo,
	configRepo *repository.KonfigurasiRepo,
	riwayatRepo *repository.RiwayatRepo,
) *LeaveUsecase {
	return &LeaveUsecase{
		LeaveRepo:    leaveRepo,
		KaryawanRepo: karyawanRepo,
		TTDRepo:      ttdRepo,
		ConfigRepo:   configRepo,
		RiwayatRepo:  riwayatRepo,
	}
}

func (u *LeaveUsecase) DetermineApprovalFlow(ctx context.Context, karyawan *domain.Karyawan) (*ApprovalFlow, error) {
	flow := &ApprovalFlow{
		LangsungApprove: false,
		LangsungFinal:   false,
		ButuhAtasan:     true,
		ButuhHRD:        true,
		JumlahTTD:       3,
		AtasanID:        nil,
		HRDID:           nil,
	}

	switch karyawan.Role {
	case "admin":
		flow.LangsungApprove = true
		flow.LangsungFinal = true
		flow.ButuhAtasan = false
		flow.ButuhHRD = false
		flow.JumlahTTD = 1

	case "hrd":
		flow.LangsungApprove = true
		flow.LangsungFinal = true
		flow.ButuhAtasan = true
		flow.ButuhHRD = false
		flow.JumlahTTD = 2

		atasanList, _, err := u.KaryawanRepo.GetAll(ctx, 10, 0, "", "atasan", "aktif")
		if err == nil && len(atasanList) > 0 {
			for _, atasan := range atasanList {
				if atasan.ID != karyawan.ID {
					flow.AtasanID = &atasan.ID
					break
				}
			}
		}

	case "atasan", "manager":
		flow.LangsungApprove = true
		flow.LangsungFinal = true
		flow.ButuhAtasan = false
		flow.ButuhHRD = true
		flow.JumlahTTD = 2

		hrdList, _, err := u.KaryawanRepo.GetAll(ctx, 10, 0, "", "hrd", "aktif")
		if err == nil && len(hrdList) > 0 {
			flow.HRDID = &hrdList[0].ID
		}

	default:
		flow.LangsungApprove = false
		flow.LangsungFinal = false
		flow.ButuhAtasan = true
		flow.ButuhHRD = true
		flow.JumlahTTD = 3

		if karyawan.AtasanLangsungID != nil {
			flow.AtasanID = karyawan.AtasanLangsungID
		}
	}

	return flow, nil
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

	balance, err := u.LeaveRepo.GetBalanceWithCarryOver(ctx, karyawanID, time.Now().Year())
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

	flow, err := u.DetermineApprovalFlow(ctx, karyawan)
	if err != nil {
		return nil, err
	}

	judulDokumen := "PERMOHONAN/LAPORAN CUTI TAHUNAN"
	if req.SubTipe == "dispensasi" {
		judulDokumen = "PERMOHONAN/LAPORAN DISPENSASI"
	}

	status := "menunggu"
	if flow.LangsungApprove {
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
		LangsungApprove:   flow.LangsungApprove,
		LangsungFinal:     flow.LangsungFinal,
		JudulDokumen:      judulDokumen,
	}

	if flow.LangsungFinal {
		if flow.HRDID != nil {
			leave.DifinalisasiOleh = flow.HRDID
		} else {
			hrd, err := u.KaryawanRepo.GetByRole(ctx, "hrd")
			if err == nil && hrd != nil {
				leave.DifinalisasiOleh = &hrd.ID
			}
		}
		if leave.DifinalisasiOleh != nil {
			now := time.Now()
			leave.TanggalDifinalisasi = &now
		}
	}

	if err := u.LeaveRepo.Create(ctx, leave); err != nil {
		return nil, err
	}

	if leave.MengurangiCuti {
		tahun := leave.TanggalMulai.Year()
		u.LeaveRepo.UpdateAkanDilaksanakan(ctx, karyawanID, tahun, totalHari)
	}

	if u.RiwayatRepo != nil {
		detail := "Pengajuan cuti " + req.SubTipe + " " + strconv.Itoa(totalHari) + " hari dari " + req.TanggalMulai + " sampai " + req.TanggalSelesai
		u.RiwayatRepo.CreateRiwayat(ctx, karyawanID, "cuti_diajukan", detail)
	}

	telegram := utils.NewTelegramBot()
	if telegram != nil && karyawan.AtasanLangsungID != nil {
		atasan, _ := u.KaryawanRepo.GetByID(ctx, *karyawan.AtasanLangsungID)
		if atasan != nil && atasan.TelegramChatID != nil && *atasan.TelegramChatID != "" {
			go telegram.SendCreateLeaveNotification(*atasan.TelegramChatID, karyawan.NamaLengkap, strconv.Itoa(totalHari), req.Alasan)
		}
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

	sekarang := time.Now()
	tanggalMulai := leave.TanggalMulai
	batasBatal := time.Date(tanggalMulai.Year(), tanggalMulai.Month(), tanggalMulai.Day(), 0, 0, 0, 0, time.Local).Add(-24 * time.Hour)
	batasBatal = time.Date(batasBatal.Year(), batasBatal.Month(), batasBatal.Day(), 23, 59, 59, 0, time.Local)

	if sekarang.After(batasBatal) {
		return nil, errors.New("pembatalan cuti hanya dapat dilakukan maksimal H-24 jam sebelum tanggal mulai cuti")
	}

	if err := u.LeaveRepo.UpdateStatus(ctx, leaveID, "dibatalkan"); err != nil {
		return nil, err
	}

	if leave.MengurangiCuti {
		tahun := leave.TanggalMulai.Year()
		u.LeaveRepo.UpdateBalance(ctx, leave.KaryawanID, tahun)
	}

	if u.RiwayatRepo != nil {
		detail := "Cuti dibatalkan oleh karyawan"
		u.RiwayatRepo.CreateRiwayat(ctx, karyawanID, "cuti_dibatalkan", detail)
	}

	karyawan, _ := u.KaryawanRepo.GetByID(ctx, leave.KaryawanID)
	if karyawan != nil {
		telegram := utils.NewTelegramBot()
		if telegram != nil && leave.DisetujuiOleh != nil {
			atasan, _ := u.KaryawanRepo.GetByID(ctx, *leave.DisetujuiOleh)
			if atasan != nil && atasan.TelegramChatID != nil && *atasan.TelegramChatID != "" {
				go telegram.SendCancelLeaveNotification(*atasan.TelegramChatID, karyawan.NamaLengkap, "Dibatalkan oleh karyawan")
			}
		}
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

	if leave.MengurangiCuti {
		tahun := leave.TanggalMulai.Year()
		u.LeaveRepo.UpdateBalance(ctx, leave.KaryawanID, tahun)
	}

	if u.RiwayatRepo != nil {
		karyawan, _ := u.KaryawanRepo.GetByID(ctx, managerID)
		namaAtasan := "atasan"
		if karyawan != nil {
			namaAtasan = karyawan.NamaLengkap
		}
		detail := "Cuti disetujui oleh " + namaAtasan
		u.RiwayatRepo.CreateRiwayat(ctx, leave.KaryawanID, "cuti_disetujui", detail)
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

	if leave.MengurangiCuti {
		tahun := leave.TanggalMulai.Year()
		u.LeaveRepo.UpdateBalance(ctx, leave.KaryawanID, tahun)
	}

	if u.RiwayatRepo != nil {
		karyawan, _ := u.KaryawanRepo.GetByID(ctx, managerID)
		namaAtasan := "atasan"
		if karyawan != nil {
			namaAtasan = karyawan.NamaLengkap
		}
		detail := "Cuti ditolak oleh " + namaAtasan + " dengan alasan: " + alasan
		u.RiwayatRepo.CreateRiwayat(ctx, leave.KaryawanID, "cuti_ditolak", detail)
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

	if leave.MengurangiCuti {
		tahun := leave.TanggalMulai.Year()
		u.LeaveRepo.UpdateBalance(ctx, leave.KaryawanID, tahun)
	}

	if u.RiwayatRepo != nil {
		karyawan, _ := u.KaryawanRepo.GetByID(ctx, hrdID)
		namaHRD := "HRD"
		if karyawan != nil {
			namaHRD = karyawan.NamaLengkap
		}
		detail := "Cuti difinalisasi oleh " + namaHRD
		u.RiwayatRepo.CreateRiwayat(ctx, leave.KaryawanID, "cuti_difinalisasi", detail)
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

	flow, err := u.DetermineApprovalFlow(ctx, karyawan)
	if err != nil {
		return nil, "", err
	}

	var atasan *domain.Karyawan
	if flow.ButuhAtasan {
		if leave.DisetujuiOleh != nil {
			atasan, _ = u.KaryawanRepo.GetByID(ctx, *leave.DisetujuiOleh)
		}
		if atasan == nil && flow.AtasanID != nil {
			atasan, _ = u.KaryawanRepo.GetByID(ctx, *flow.AtasanID)
		}
		if atasan == nil && karyawan.AtasanLangsungID != nil {
			atasan, _ = u.KaryawanRepo.GetByID(ctx, *karyawan.AtasanLangsungID)
		}
		if atasan == nil {
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
	}

	var hrd *domain.Karyawan
	if flow.ButuhHRD {
		if leave.DifinalisasiOleh != nil {
			hrd, _ = u.KaryawanRepo.GetByID(ctx, *leave.DifinalisasiOleh)
		}
		if hrd == nil && flow.HRDID != nil {
			hrd, _ = u.KaryawanRepo.GetByID(ctx, *flow.HRDID)
		}
		if hrd == nil {
			hrdList, _, err := u.KaryawanRepo.GetAll(ctx, 10, 0, "", "hrd", "aktif")
			if err == nil && len(hrdList) > 0 {
				hrd = &hrdList[0]
			}
		}
	}

	var ttdKaryawanURL string
	ttdKaryawan, err := u.TTDRepo.GetByKaryawanID(ctx, leave.KaryawanID)
	if err == nil && ttdKaryawan != nil {
		ttdKaryawanURL = ttdKaryawan.URLTandaTangan
	}

	var ttdAtasanURL string
	if atasan != nil {
		ttdAtasan, err := u.TTDRepo.GetByKaryawanID(ctx, atasan.ID)
		if err == nil && ttdAtasan != nil {
			ttdAtasanURL = ttdAtasan.URLTandaTangan
		}
	}

	var ttdHRDURL string
	if hrd != nil {
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

	jumlahTTD := flow.JumlahTTD

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

func (u *LeaveUsecase) GetBalance(ctx context.Context, karyawanID string, year int) (*domain.BalanceResponse, error) {
	var tahun int
	if year > 0 {
		tahun = year
	} else {
		tahun = time.Now().Year()
	}

	balance, err := u.LeaveRepo.GetBalanceWithCarryOver(ctx, karyawanID, tahun)
	if err != nil {
		return nil, err
	}

	var sisaCutiTahunLalu int
	var berlakuSampai string
	sekarang := time.Now()
	batasAkhirMaret := time.Date(sekarang.Year(), 3, 31, 23, 59, 59, 0, time.Local)

	if sekarang.Before(batasAkhirMaret) || sekarang.Equal(batasAkhirMaret) {
		balanceTahunLalu, err := u.LeaveRepo.GetBalance(ctx, karyawanID, tahun-1)
		if err == nil && balanceTahunLalu != nil && balanceTahunLalu.SisaCuti > 0 {
			sisaCutiTahunLalu = balanceTahunLalu.SisaCuti
			berlakuSampai = "31 Maret " + strconv.Itoa(tahun)
		}
	}

	if balance == nil {
		return &domain.BalanceResponse{
			Tahun:                        tahun,
			JumlahCuti:                   12,
			TelahDilaksanakan:            0,
			AkanDilaksanakan:             0,
			SisaCuti:                     12,
			SisaCutiTahunLalu:            sisaCutiTahunLalu,
			SisaCutiTahunLaluBerlakuSampai: berlakuSampai,
		}, nil
	}

	return &domain.BalanceResponse{
		Tahun:                        balance.Tahun,
		JumlahCuti:                   balance.JumlahCuti,
		TelahDilaksanakan:            balance.TelahDilaksanakan,
		AkanDilaksanakan:             balance.AkanDilaksanakan,
		SisaCuti:                     balance.SisaCuti,
		SisaCutiTahunLalu:            sisaCutiTahunLalu,
		SisaCutiTahunLaluBerlakuSampai: berlakuSampai,
	}, nil
}

func (u *LeaveUsecase) GetAllLeaves(ctx context.Context, userID string, role string, req domain.LeaveFilterRequest) ([]domain.LeaveWithKaryawanResponse, int, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	var atasanID string
	if role == "atasan" {
		karyawan, err := u.KaryawanRepo.GetByID(ctx, userID)
		if err != nil || karyawan == nil {
			return nil, 0, errors.New("karyawan tidak ditemukan")
		}
		atasanID = userID
	}

	return u.LeaveRepo.GetAllLeaves(ctx, atasanID, role, req.Status, req.SubTipe, req.StartDate, req.EndDate, limit, offset)
}

func (u *LeaveUsecase) GetApprovalList(ctx context.Context, atasanID string, limit int, page int) ([]domain.LeaveWithKaryawanResponse, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	return u.LeaveRepo.GetPendingLeavesByAtasan(ctx, atasanID, limit, offset)
}

func (u *LeaveUsecase) GetFinalizationList(ctx context.Context, limit int, page int) ([]domain.LeaveWithKaryawanResponse, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	return u.LeaveRepo.GetAllLeaves(ctx, "", "hrd", "disetujui", "", "", "", limit, offset)
}
package usecase

import (
	"context"
	"errors"
	"time"
	
	"sakti_apps_be/internal/repository"
	"sakti_apps_be/internal/utils"
)

type LeaveUsecase struct {
	LeaveRepo    *repository.LeaveRepo
	KaryawanRepo *repository.KaryawanRepo
}

func NewLeaveUsecase(leaveRepo *repository.LeaveRepo, karyawanRepo *repository.KaryawanRepo) *LeaveUsecase {
	return &LeaveUsecase{
		LeaveRepo:    leaveRepo,
		KaryawanRepo: karyawanRepo,
	}
}

func (u *LeaveUsecase) CreateLeave(ctx context.Context, karyawanID string, req struct {
	TipePengajuan  string `json:"tipe_pengajuan"`
	TanggalMulai   string `json:"tanggal_mulai"`
	TanggalSelesai string `json:"tanggal_selesai"`
	Alasan         string `json:"alasan"`
}) error {
	karyawan, err := u.KaryawanRepo.GetByID(ctx, karyawanID)
	if err != nil || karyawan == nil {
		return errors.New("karyawan tidak ditemukan")
	}

	start, _ := time.Parse("2006-01-02", req.TanggalMulai)
	end, _ := time.Parse("2006-01-02", req.TanggalSelesai)
	totalHari := int(end.Sub(start).Hours()/24) + 1

	if totalHari <= 0 {
		return errors.New("tanggal tidak valid")
	}

	balance, _ := u.LeaveRepo.GetBalance(ctx, karyawanID, time.Now().Year())
	if req.TipePengajuan == "cuti" && totalHari > balance.Sisa {
		return errors.New("kuota cuti tidak mencukupi")
	}

	leave := &repository.Leave{
		KaryawanID:     karyawanID,
		TipePengajuan:  req.TipePengajuan,
		TanggalMulai:   req.TanggalMulai,
		TanggalSelesai: req.TanggalSelesai,
		TotalHari:      totalHari,
		Alasan:         req.Alasan,
		Status:         "menunggu",
	}

	return u.LeaveRepo.Create(ctx, leave)
}

func (u *LeaveUsecase) GetStatus(ctx context.Context, karyawanID string) ([]repository.Leave, error) {
	return u.LeaveRepo.GetByKaryawanID(ctx, karyawanID)
}

func (u *LeaveUsecase) CancelLeave(ctx context.Context, leaveID, karyawanID string) error {
	leave, err := u.LeaveRepo.GetByID(ctx, leaveID)
	if err != nil {
		return errors.New("pengajuan tidak ditemukan")
	}
	if leave.KaryawanID != karyawanID {
		return errors.New("anda tidak memiliki akses")
	}
	if leave.Status != "menunggu" && leave.Status != "disetujui" {
		return errors.New("pengajuan tidak bisa dibatalkan")
	}

	return u.LeaveRepo.UpdateStatus(ctx, leaveID, "dibatalkan")
}

func (u *LeaveUsecase) ApproveLeave(ctx context.Context, leaveID, managerID string) error {
	leave, err := u.LeaveRepo.GetByID(ctx, leaveID)
	if err != nil {
		return errors.New("pengajuan tidak ditemukan")
	}
	if leave.Status != "menunggu" {
		return errors.New("pengajuan sudah diproses")
	}

	return u.LeaveRepo.Approve(ctx, leaveID, managerID)
}

func (u *LeaveUsecase) RejectLeave(ctx context.Context, leaveID, managerID, alasan string) error {
	leave, err := u.LeaveRepo.GetByID(ctx, leaveID)
	if err != nil {
		return errors.New("pengajuan tidak ditemukan")
	}
	if leave.Status != "menunggu" {
		return errors.New("pengajuan sudah diproses")
	}

	return u.LeaveRepo.Reject(ctx, leaveID, managerID, alasan)
}

func (u *LeaveUsecase) FinalizeLeave(ctx context.Context, leaveID, hrdID string) error {
	leave, err := u.LeaveRepo.GetByID(ctx, leaveID)
	if err != nil {
		return errors.New("pengajuan tidak ditemukan")
	}
	if leave.Status != "disetujui" {
		return errors.New("pengajuan belum disetujui atasan")
	}

	return u.LeaveRepo.Finalize(ctx, leaveID, hrdID)
}

func (u *LeaveUsecase) GenerateSuratCuti(ctx context.Context, leaveID string) ([]byte, error) {
	leave, err := u.LeaveRepo.GetByID(ctx, leaveID)
	if err != nil {
		return nil, err
	}

	if leave.Status != "disetujui" && leave.Status != "ditolak" {
		return nil, errors.New("pengajuan belum difinalisasi")
	}

	karyawan, err := u.KaryawanRepo.GetByID(ctx, leave.KaryawanID)
	if err != nil {
		return nil, err
	}

	atasan, _ := u.KaryawanRepo.GetByID(ctx, *leave.DisetujuiOleh)
	hrd, _ := u.KaryawanRepo.GetByID(ctx, *leave.DifinalisasiOleh)

	balance, _ := u.LeaveRepo.GetBalance(ctx, leave.KaryawanID, time.Now().Year())

	jenis := "CUTI"
	if leave.TipePengajuan == "dispen" {
		jenis = "DISPENSASI"
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

	pdfData := utils.PDFData{
		CompanyName:    "KOPEGTEL MALANG",
		CompanyAddress: "Jl. Ahmad Yani No.11, Blimbing, Kota Malang",

		Jenis:       jenis,
		Nama:        karyawan.NamaLengkap,
		Divisi:      divisi,
		Unit:        unit,
		Jabatan:     karyawan.LevelJabatan,
		LamaCuti:    leave.TotalHari,
		TanggalCuti: leave.TanggalMulai + " s.d " + leave.TanggalSelesai,
		AlasanCuti:  leave.Alasan,

		JumlahCutiTahun:      12,
		CutiDilaksanakan:     balance.Digunakan - leave.TotalHari,
		CutiAkanDilaksanakan: leave.TotalHari,
		SisaCuti:             balance.Sisa,

		DisetujuiSelama: leave.TotalHari,
		MulaiTanggal:    leave.TanggalMulai,
		Tahun:           time.Now().Year(),
		TanggalSekarang: time.Now(),

		NamaPemohon: karyawan.NamaLengkap,
		NamaAtasan:  namaAtasan,
		NamaHRD:     namaHRD,
	}

	return utils.GeneratePDF(pdfData)
}
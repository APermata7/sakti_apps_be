package usecase

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminUsecase struct {
	DB *pgxpool.Pool
}

func NewAdminUsecase(db *pgxpool.Pool) *AdminUsecase {
	return &AdminUsecase{DB: db}
}

type DashboardStats struct {
	TotalKaryawan      int `json:"total_karyawan"`
	KaryawanAktif      int `json:"karyawan_aktif"`
	PresensiHariIni    int `json:"presensi_hari_ini"`
	PresensiTerlambat  int `json:"presensi_terlambat"`
	CutiPending        int `json:"cuti_pending"`
	TotalCutiTahun     int `json:"total_cuti_tahun"`
}

func (u *AdminUsecase) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	var stats DashboardStats

	queryTotal := `SELECT COUNT(*) FROM karyawan`
	u.DB.QueryRow(ctx, queryTotal).Scan(&stats.TotalKaryawan)

	queryAktif := `SELECT COUNT(*) FROM karyawan WHERE status_karyawan = 'aktif'`
	u.DB.QueryRow(ctx, queryAktif).Scan(&stats.KaryawanAktif)

	queryPresensi := `SELECT COUNT(*) FROM presensi WHERE tanggal = CURRENT_DATE`
	u.DB.QueryRow(ctx, queryPresensi).Scan(&stats.PresensiHariIni)

	queryTerlambat := `SELECT COUNT(*) FROM presensi WHERE tanggal = CURRENT_DATE AND status = 'terlambat'`
	u.DB.QueryRow(ctx, queryTerlambat).Scan(&stats.PresensiTerlambat)

	queryPending := `SELECT COUNT(*) FROM pengajuan_cuti WHERE status = 'menunggu'`
	u.DB.QueryRow(ctx, queryPending).Scan(&stats.CutiPending)

	queryTotalCuti := `SELECT COALESCE(SUM(total_hari), 0) FROM pengajuan_cuti WHERE status = 'disetujui' AND EXTRACT(YEAR FROM dibuat_pada) = EXTRACT(YEAR FROM CURRENT_DATE)`
	u.DB.QueryRow(ctx, queryTotalCuti).Scan(&stats.TotalCutiTahun)

	return &stats, nil
}
package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
)

type AdminUsecase struct {
	DB           *pgxpool.Pool
	KaryawanRepo *repository.KaryawanRepo
	SupabaseURL  string
	AnonKey      string
}

func NewAdminUsecase(db *pgxpool.Pool, karyawanRepo *repository.KaryawanRepo) *AdminUsecase {
	return &AdminUsecase{
		DB:           db,
		KaryawanRepo: karyawanRepo,
		SupabaseURL:  os.Getenv("SUPABASE_URL"),
		AnonKey:      os.Getenv("SUPABASE_ANON_KEY"),
	}
}

type DashboardStats struct {
	TotalKaryawan     int `json:"total_karyawan"`
	KaryawanAktif     int `json:"karyawan_aktif"`
	PresensiHariIni   int `json:"presensi_hari_ini"`
	PresensiTerlambat int `json:"presensi_terlambat"`
	CutiPending       int `json:"cuti_pending"`
	TotalCutiTahun    int `json:"total_cuti_tahun"`
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

func (u *AdminUsecase) CreateKaryawan(ctx context.Context, req domain.CreateKaryawanRequest) (*domain.Karyawan, error) {
	existing, _ := u.KaryawanRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	userMetadata := map[string]interface{}{
		"nama_lengkap": req.NamaLengkap,
		"role":         req.Role,
	}

	if req.LevelJabatan != nil {
		userMetadata["level_jabatan"] = *req.LevelJabatan
	} else {
		userMetadata["level_jabatan"] = nil
	}

	supabaseReq := map[string]interface{}{
		"email":          req.Email,
		"password":       req.Password,
		"user_metadata":  userMetadata,
	}

	jsonBody, _ := json.Marshal(supabaseReq)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", u.SupabaseURL+"/auth/v1/admin/users", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("apikey", u.AnonKey)
	httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_SERVICE_KEY"))
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("gagal membuat akun di Supabase")
	}

	var authResp struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	karyawan := &domain.Karyawan{
		ID:             authResp.ID,
		NamaLengkap:    req.NamaLengkap,
		Email:          req.Email,
		Role:           req.Role,
		LevelJabatan:   req.LevelJabatan,
		StatusKaryawan: "aktif",
	}

	if req.NomorTelepon != nil && *req.NomorTelepon != "" {
		karyawan.NomorTelepon = req.NomorTelepon
	}
	if req.FotoURL != nil && *req.FotoURL != "" {
		karyawan.FotoURL = req.FotoURL
	}
	if req.Divisi != nil && *req.Divisi != "" {
		karyawan.Divisi = req.Divisi
	}
	if req.Unit != nil && *req.Unit != "" {
		karyawan.Unit = req.Unit
	}
	if req.AtasanLangsungID != nil && *req.AtasanLangsungID != "" {
		karyawan.AtasanLangsungID = req.AtasanLangsungID
	}

	if err := u.KaryawanRepo.Create(ctx, karyawan); err != nil {
		return nil, err
	}

	return karyawan, nil
}

func (u *AdminUsecase) GetAllKaryawan(ctx context.Context, page, limit int, search, role, status string) ([]domain.Karyawan, int, error) {
	offset := (page - 1) * limit
	return u.KaryawanRepo.GetAll(ctx, limit, offset, search, role, status)
}

func (u *AdminUsecase) GetKaryawanByID(ctx context.Context, id string) (*domain.Karyawan, error) {
	return u.KaryawanRepo.GetByID(ctx, id)
}

func (u *AdminUsecase) UpdateKaryawan(ctx context.Context, id string, req domain.UpdateKaryawanRequest) (*domain.Karyawan, error) {
	existing, err := u.KaryawanRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("karyawan tidak ditemukan")
	}

	if req.NamaLengkap != nil && *req.NamaLengkap != "" {
		existing.NamaLengkap = *req.NamaLengkap
	}
	if req.NomorTelepon != nil {
		if *req.NomorTelepon != "" {
			existing.NomorTelepon = req.NomorTelepon
		} else {
			existing.NomorTelepon = nil
		}
	}
	if req.FotoURL != nil {
		if *req.FotoURL != "" {
			existing.FotoURL = req.FotoURL
		} else {
			existing.FotoURL = nil
		}
	}
	if req.Role != nil && *req.Role != "" {
		existing.Role = *req.Role
	}
	if req.LevelJabatan != nil {
		if *req.LevelJabatan != "" {
			existing.LevelJabatan = req.LevelJabatan
		} else {
			existing.LevelJabatan = nil
		}
	}
	if req.AtasanLangsungID != nil {
		if *req.AtasanLangsungID != "" {
			existing.AtasanLangsungID = req.AtasanLangsungID
		} else {
			existing.AtasanLangsungID = nil
		}
	}
	if req.Divisi != nil {
		if *req.Divisi != "" {
			existing.Divisi = req.Divisi
		} else {
			existing.Divisi = nil
		}
	}
	if req.Unit != nil {
		if *req.Unit != "" {
			existing.Unit = req.Unit
		} else {
			existing.Unit = nil
		}
	}
	if req.StatusKaryawan != nil && *req.StatusKaryawan != "" {
		existing.StatusKaryawan = *req.StatusKaryawan
	}

	if err := u.KaryawanRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (u *AdminUsecase) DeleteKaryawan(ctx context.Context, id string) error {
	existing, err := u.KaryawanRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("karyawan tidak ditemukan")
	}

	if existing.Role == "admin" {
		var count int
		query := `SELECT COUNT(*) FROM karyawan WHERE role = 'admin' AND status_karyawan = 'aktif'`
		u.DB.QueryRow(ctx, query).Scan(&count)
		if count <= 1 {
			return errors.New("tidak dapat menghapus admin terakhir")
		}
	}

	return u.KaryawanRepo.Delete(ctx, id)
}
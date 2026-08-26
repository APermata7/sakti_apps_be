package usecase

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
)

type AdminUsecase struct {
	DB              *pgxpool.Pool
	KaryawanRepo    *repository.KaryawanRepo
	PresensiRepo    *repository.PresensiRepo
	LeaveRepo       *repository.LeaveRepo
	LiburUsecase    *LiburUsecase
	KonfigurasiRepo *repository.KonfigurasiRepo
	TTDRepo         *repository.TTDRepo
	LogUsecase      *LogUsecase
	SupabaseURL     string
	AnonKey         string
}

func NewAdminUsecase(
	db *pgxpool.Pool,
	karyawanRepo *repository.KaryawanRepo,
	presensiRepo *repository.PresensiRepo,
	leaveRepo *repository.LeaveRepo,
	liburUsecase *LiburUsecase,
	konfigurasiRepo *repository.KonfigurasiRepo,
	ttdRepo *repository.TTDRepo,
	logUsecase *LogUsecase,
) *AdminUsecase {
	return &AdminUsecase{
		DB:              db,
		KaryawanRepo:    karyawanRepo,
		PresensiRepo:    presensiRepo,
		LeaveRepo:       leaveRepo,
		LiburUsecase:    liburUsecase,
		KonfigurasiRepo: konfigurasiRepo,
		TTDRepo:         ttdRepo,
		LogUsecase:      logUsecase,
		SupabaseURL:     os.Getenv("SUPABASE_URL"),
		AnonKey:         os.Getenv("SUPABASE_ANON_KEY"),
	}
}

type MonthlyChartData struct {
	Tanggal    string `json:"tanggal"`
	TepatWaktu int    `json:"tepat_waktu"`
	Terlambat  int    `json:"terlambat"`
}

type DashboardStats struct {
	TotalKaryawan      int         `json:"total_karyawan"`
	KaryawanAktif      int         `json:"karyawan_aktif"`
	TotalTerlambat     int         `json:"total_terlambat"`
	TotalLembur        int         `json:"total_lembur"`
	TotalCutiDisetujui int         `json:"total_cuti_disetujui"`
	KaryawanPerDept    []DeptStat  `json:"karyawan_per_dept"`
	PresensiMasuk      []StatusStat `json:"presensi_masuk"`
	PresensiKeluar     []StatusStat `json:"presensi_keluar"`
	TotalPengajuanCuti []CutiStat  `json:"total_pengajuan_cuti"`
}

type DeptStat struct {
	Departemen string `json:"departemen"`
	Total      int    `json:"total"`
}

type StatusStat struct {
	Status string `json:"status"`
	Total  int    `json:"total"`
}

type CutiStat struct {
	Status string `json:"status"`
	Total  int    `json:"total"`
}

type PresensiReportItem struct {
	ID                   string  `json:"id"`
	KaryawanNama         string  `json:"karyawan_nama"`
	Tanggal              string  `json:"tanggal"`
	JamMasuk             string  `json:"jam_masuk"`
	StatusMasuk          string  `json:"status_masuk"`
	JamKeluar            string  `json:"jam_keluar"`
	StatusKeluar         string  `json:"status_keluar"`
	JenisCuti            string  `json:"jenis_cuti"`
	LocationStatusMasuk  string  `json:"location_status_masuk"`
	LocationStatusKeluar string  `json:"location_status_keluar"`
}

type CutiReportItem struct {
	ID             string `json:"id"`
	KaryawanNama   string `json:"karyawan_nama"`
	Divisi         string `json:"divisi"`
	SubTipe        string `json:"sub_tipe"`
	Status         string `json:"status"`
	TanggalMulai   string `json:"tanggal_mulai"`
	TanggalSelesai string `json:"tanggal_selesai"`
	TotalHari      int    `json:"total_hari"`
	SisaCuti       int    `json:"sisa_cuti"`
}

func (u *AdminUsecase) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	var stats DashboardStats

	queryTotal := `SELECT COUNT(*) FROM karyawan WHERE role != 'admin'`
	u.DB.QueryRow(ctx, queryTotal).Scan(&stats.TotalKaryawan)

	queryAktif := `SELECT COUNT(*) FROM karyawan WHERE status_karyawan = 'aktif' AND role != 'admin'`
	u.DB.QueryRow(ctx, queryAktif).Scan(&stats.KaryawanAktif)

	monthStart := time.Now().AddDate(0, 0, -30)
	queryTerlambat := `SELECT COUNT(*) FROM presensi WHERE tanggal >= $1 AND status = 'terlambat'`
	u.DB.QueryRow(ctx, queryTerlambat, monthStart).Scan(&stats.TotalTerlambat)

	queryLembur := `SELECT COUNT(*) FROM presensi WHERE tanggal >= $1 AND lembur = true`
	u.DB.QueryRow(ctx, queryLembur, monthStart).Scan(&stats.TotalLembur)

	queryCutiDisetujui := `SELECT COUNT(*) FROM pengajuan_cuti WHERE status = 'disetujui' AND EXTRACT(MONTH FROM dibuat_pada) = EXTRACT(MONTH FROM CURRENT_DATE)`
	u.DB.QueryRow(ctx, queryCutiDisetujui).Scan(&stats.TotalCutiDisetujui)

	deptQuery := `
		SELECT 
			CASE 
				WHEN divisi IS NOT NULL AND unit IS NOT NULL THEN divisi || ' - ' || unit 
				WHEN divisi IS NOT NULL THEN divisi 
				ELSE 'Tidak Ada' 
			END as departemen,
			COUNT(*) as total
		FROM karyawan 
		WHERE status_karyawan = 'aktif' AND role != 'admin'
		GROUP BY departemen
		ORDER BY total DESC
	`
	rows, err := u.DB.Query(ctx, deptQuery)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var dept DeptStat
			rows.Scan(&dept.Departemen, &dept.Total)
			stats.KaryawanPerDept = append(stats.KaryawanPerDept, dept)
		}
	}

	masukQuery := `
		SELECT 
			CASE 
				WHEN jam_masuk IS NULL THEN 'Belum Presensi'
				WHEN status = 'terlambat' THEN 'Masuk Terlambat'
				ELSE 'Masuk Tepat Waktu'
			END as status,
			COUNT(*) as total
		FROM presensi 
		WHERE tanggal = CURRENT_DATE
		GROUP BY status
	`
	rows, err = u.DB.Query(ctx, masukQuery)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var stat StatusStat
			rows.Scan(&stat.Status, &stat.Total)
			stats.PresensiMasuk = append(stats.PresensiMasuk, stat)
		}
	}

	keluarQuery := `
		SELECT 
			CASE 
				WHEN jam_keluar IS NULL THEN 'Belum Presensi'
				WHEN lembur = true THEN 'Presensi Lembur'
				ELSE 'Presensi Keluar'
			END as status,
			COUNT(*) as total
		FROM presensi 
		WHERE tanggal = CURRENT_DATE
		GROUP BY status
	`
	rows, err = u.DB.Query(ctx, keluarQuery)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var stat StatusStat
			rows.Scan(&stat.Status, &stat.Total)
			stats.PresensiKeluar = append(stats.PresensiKeluar, stat)
		}
	}

	cutiQuery := `
		SELECT status, COUNT(*) as total
		FROM pengajuan_cuti
		GROUP BY status
	`
	rows, err = u.DB.Query(ctx, cutiQuery)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var stat CutiStat
			rows.Scan(&stat.Status, &stat.Total)
			stats.TotalPengajuanCuti = append(stats.TotalPengajuanCuti, stat)
		}
	}

	return &stats, nil
}

func (u *AdminUsecase) GetMonthlyChart(ctx context.Context, bulan, tahun int) ([]MonthlyChartData, error) {
	var results []MonthlyChartData
	query := `
		SELECT 
			TO_CHAR(tanggal, 'YYYY-MM-DD') as tanggal,
			COUNT(CASE WHEN status = 'tepat_waktu' THEN 1 END) as tepat_waktu,
			COUNT(CASE WHEN status = 'terlambat' THEN 1 END) as terlambat
		FROM presensi
		WHERE EXTRACT(YEAR FROM tanggal) = $1 AND EXTRACT(MONTH FROM tanggal) = $2
		GROUP BY tanggal
		ORDER BY tanggal ASC
	`
	rows, err := u.DB.Query(ctx, query, tahun, bulan)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d MonthlyChartData
		if err := rows.Scan(&d.Tanggal, &d.TepatWaktu, &d.Terlambat); err != nil {
			return nil, err
		}
		results = append(results, d)
	}
	return results, nil
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
		"email_confirm":  true,
	}

	jsonBody, _ := json.Marshal(supabaseReq)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", u.SupabaseURL+"/auth/v1/admin/users", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("apikey", u.AnonKey)
	httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, errors.New("gagal membuat akun di Supabase: " + resp.Status)
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

	if err := u.LeaveRepo.UpdateBalance(ctx, karyawan.ID, time.Now().Year()); err != nil {
		log.Printf("Gagal inisialisasi sisa cuti untuk karyawan baru: %v", err)
	}

	if u.LogUsecase != nil {
		detail := "Menambahkan karyawan: " + req.Email + " (" + req.NamaLengkap + ")"
		u.LogUsecase.CreateLog(ctx, karyawan.ID, "create_karyawan", detail)
	}

	return karyawan, nil
}

func (u *AdminUsecase) GetAllKaryawan(ctx context.Context, page, limit int, search, role, levelJabatan, divisi, unit, status string) ([]domain.Karyawan, int, error) {
	offset := (page - 1) * limit

	karyawanList, total, err := u.KaryawanRepo.GetAll(ctx, limit, offset, search, role, levelJabatan, divisi, unit, status)
	if err != nil {
		return nil, 0, err
	}

	for i := range karyawanList {
		ttd, err := u.TTDRepo.GetByKaryawanID(ctx, karyawanList[i].ID)
		if err == nil && ttd != nil {
			karyawanList[i].TandaTangan = &domain.TandaTangan{
				ID:             ttd.ID,
				KaryawanID:     ttd.KaryawanID,
				URLTandaTangan: ttd.URLTandaTangan,
				DiunggahPada:   ttd.DiunggahPada,
				DiperbaruiPada: ttd.DiperbaruiPada,
			}
		} else {
			karyawanList[i].TandaTangan = nil
		}
	}

	return karyawanList, total, nil
}

func (u *AdminUsecase) GetKaryawanByID(ctx context.Context, id string) (*domain.Karyawan, error) {
	karyawan, err := u.KaryawanRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if karyawan == nil {
		return nil, nil
	}

	ttd, err := u.TTDRepo.GetByKaryawanID(ctx, id)
	if err == nil && ttd != nil {
		karyawan.TandaTangan = &domain.TandaTangan{
			ID:             ttd.ID,
			KaryawanID:     ttd.KaryawanID,
			URLTandaTangan: ttd.URLTandaTangan,
			DiunggahPada:   ttd.DiunggahPada,
			DiperbaruiPada: ttd.DiperbaruiPada,
		}
	} else {
		karyawan.TandaTangan = nil
	}

	return karyawan, nil
}

func (u *AdminUsecase) UpdateKaryawan(ctx context.Context, id string, req domain.UpdateKaryawanRequest) (*domain.Karyawan, error) {
	existing, err := u.KaryawanRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("karyawan tidak ditemukan")
	}

	if req.StatusKaryawan != nil && *req.StatusKaryawan != "" {
		existing.StatusKaryawan = *req.StatusKaryawan
	}

	if req.Role != nil && *req.Role != "" {
		existing.Role = *req.Role
	}

	if req.AtasanLangsungID != nil {
		if *req.AtasanLangsungID != "" {
			existing.AtasanLangsungID = req.AtasanLangsungID
		} else {
			existing.AtasanLangsungID = nil
		}
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
	if req.LevelJabatan != nil {
		if *req.LevelJabatan != "" {
			existing.LevelJabatan = req.LevelJabatan
		} else {
			existing.LevelJabatan = nil
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

	if err := u.KaryawanRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	if u.LogUsecase != nil {
		detail := "Memperbarui karyawan: " + existing.Email + " (" + existing.NamaLengkap + ")"
		u.LogUsecase.CreateLog(ctx, existing.ID, "update_karyawan", detail)
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

	if err := u.KaryawanRepo.Delete(ctx, id); err != nil {
		return err
	}

	if u.LogUsecase != nil {
		detail := "Menonaktifkan karyawan: " + existing.Email + " (" + existing.NamaLengkap + ")"
		u.LogUsecase.CreateLog(ctx, existing.ID, "delete_karyawan", detail)
	}

	return nil
}

func (u *AdminUsecase) ActivateKaryawan(ctx context.Context, id string) error {
	existing, err := u.KaryawanRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("karyawan tidak ditemukan")
	}
	existing.StatusKaryawan = "aktif"
	if err := u.KaryawanRepo.Update(ctx, existing); err != nil {
		return err
	}

	if u.LogUsecase != nil {
		detail := "Mengaktifkan kembali karyawan: " + existing.Email + " (" + existing.NamaLengkap + ")"
		u.LogUsecase.CreateLog(ctx, existing.ID, "activate_karyawan", detail)
	}

	return nil
}

func (u *AdminUsecase) GetPresensiReport(ctx context.Context, startDate, endDate, status, search string, limit, offset int) ([]PresensiReportItem, int, error) {
	var items []PresensiReportItem
	var total int

	if startDate == "" {
		startDate = "1970-01-01"
	}
	if endDate == "" {
		endDate = "2999-12-31"
	}

	query := `
		FROM presensi p
		JOIN karyawan k ON p.karyawan_id = k.id
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		query += fmt.Sprintf(` AND (k.nama_lengkap ILIKE $%d OR k.email ILIKE $%d)`, argIdx, argIdx+1)
		args = append(args, "%"+search+"%", "%"+search+"%")
		argIdx += 2
	}

	if startDate != "" {
		query += fmt.Sprintf(` AND p.tanggal >= $%d`, argIdx)
		args = append(args, startDate)
		argIdx++
	}
	if endDate != "" {
		query += fmt.Sprintf(` AND p.tanggal <= $%d`, argIdx)
		args = append(args, endDate)
		argIdx++
	}
	if status != "" {
		query += fmt.Sprintf(` AND p.status = $%d`, argIdx)
		args = append(args, status)
		argIdx++
	}

	countQuery := `SELECT COUNT(*) ` + query
	err := u.DB.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	dataQuery := `
		SELECT p.id, k.nama_lengkap, 
		       p.tanggal::TEXT, 
		       COALESCE(p.jam_masuk::TEXT, '') as jam_masuk,
		       COALESCE(p.status, '') as status_masuk,
		       COALESCE(p.jam_keluar::TEXT, '') as jam_keluar,
		       CASE 
		           WHEN p.jam_keluar IS NULL THEN 'Belum Presensi'
		           WHEN p.lembur = true THEN 'Presensi Lembur'
		           ELSE 'Presensi Keluar'
		       END as status_keluar,
		       '' as jenis_cuti,
		       COALESCE(p.location_status_masuk, '') as location_status_masuk,
		       COALESCE(p.location_status_keluar, '') as location_status_keluar
	` + query + ` ORDER BY p.tanggal DESC LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)

	finalArgs := append(args, limit, offset)
	rows, err := u.DB.Query(ctx, dataQuery, finalArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var item PresensiReportItem
		err := rows.Scan(
			&item.ID, &item.KaryawanNama, &item.Tanggal,
			&item.JamMasuk, &item.StatusMasuk,
			&item.JamKeluar, &item.StatusKeluar,
			&item.JenisCuti, &item.LocationStatusMasuk, &item.LocationStatusKeluar,
		)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	return items, total, nil
}

func (u *AdminUsecase) ExportPresensiCSV(ctx context.Context, startDate, endDate, status, search string) ([]byte, error) {
	items, _, err := u.GetPresensiReport(ctx, startDate, endDate, status, search, 10000, 0)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	buf.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(&buf)
	writer.Comma = '\t'

	writer.Write([]string{"ID", "Nama Karyawan", "Tanggal", "Jam Masuk", "Status Masuk", "Jam Keluar", "Status Keluar", "Jenis Cuti", "Status Lokasi Masuk", "Status Lokasi Keluar"})

	for _, item := range items {
		row := []string{
			item.ID,
			item.KaryawanNama,
			item.Tanggal,
			item.JamMasuk,
			item.StatusMasuk,
			item.JamKeluar,
			item.StatusKeluar,
			item.JenisCuti,
			item.LocationStatusMasuk,
			item.LocationStatusKeluar,
		}
		writer.Write(row)
	}
	writer.Flush()
	return buf.Bytes(), nil
}

func (u *AdminUsecase) GetCutiReport(ctx context.Context, startDate, endDate, status, search string, limit, offset int) ([]CutiReportItem, int, error) {
	var items []CutiReportItem
	var total int

	if startDate == "" {
		startDate = "1970-01-01"
	}
	if endDate == "" {
		endDate = "2999-12-31"
	}

	query := `
		FROM pengajuan_cuti pc
		JOIN karyawan k ON pc.karyawan_id = k.id
		LEFT JOIN sisa_cuti sc ON pc.karyawan_id = sc.karyawan_id AND sc.tahun = EXTRACT(YEAR FROM NOW())
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		query += fmt.Sprintf(` AND (k.nama_lengkap ILIKE $%d OR k.email ILIKE $%d)`, argIdx, argIdx+1)
		args = append(args, "%"+search+"%", "%"+search+"%")
		argIdx += 2
	}

	if startDate != "" {
		query += fmt.Sprintf(` AND pc.tanggal_mulai >= $%d`, argIdx)
		args = append(args, startDate)
		argIdx++
	}
	if endDate != "" {
		query += fmt.Sprintf(` AND pc.tanggal_selesai <= $%d`, argIdx)
		args = append(args, endDate)
		argIdx++
	}
	if status != "" {
		query += fmt.Sprintf(` AND pc.status = $%d`, argIdx)
		args = append(args, status)
		argIdx++
	}

	countQuery := `SELECT COUNT(*) ` + query
	err := u.DB.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	dataQuery := `
		SELECT pc.id, k.nama_lengkap, COALESCE(k.divisi, '') as divisi, pc.sub_tipe, pc.status,
		       pc.tanggal_mulai::TEXT, pc.tanggal_selesai::TEXT, pc.total_hari,
		       COALESCE(sc.sisa_cuti, 12) as sisa_cuti
	` + query + ` ORDER BY pc.dibuat_pada DESC LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)

	finalArgs := append(args, limit, offset)
	rows, err := u.DB.Query(ctx, dataQuery, finalArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var item CutiReportItem
		err := rows.Scan(
			&item.ID, &item.KaryawanNama, &item.Divisi, &item.SubTipe, &item.Status,
			&item.TanggalMulai, &item.TanggalSelesai, &item.TotalHari,
			&item.SisaCuti,
		)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	return items, total, nil
}

func (u *AdminUsecase) ExportCutiCSV(ctx context.Context, startDate, endDate, status, search string) ([]byte, error) {
	items, _, err := u.GetCutiReport(ctx, startDate, endDate, status, search, 10000, 0)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	buf.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(&buf)
	writer.Comma = '\t'

	writer.Write([]string{"ID", "Nama Karyawan", "Divisi", "Jenis Cuti", "Status", "Tanggal Mulai", "Tanggal Selesai", "Jumlah Hari", "Kuota Tersedia"})

	for _, item := range items {
		row := []string{
			item.ID,
			item.KaryawanNama,
			item.Divisi,
			item.SubTipe,
			item.Status,
			item.TanggalMulai,
			item.TanggalSelesai,
			strconv.Itoa(item.TotalHari),
			strconv.Itoa(item.SisaCuti),
		}
		writer.Write(row)
	}
	writer.Flush()
	return buf.Bytes(), nil
}

func (u *AdminUsecase) CreateLibur(ctx context.Context, userID string, req domain.CreateLiburRequest) (*domain.Libur, error) {
	libur, err := u.LiburUsecase.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	if u.LogUsecase != nil {
		detail := "Menambahkan hari libur: " + req.Nama + " (" + req.Tanggal + ")"
		u.LogUsecase.CreateLog(ctx, userID, "create_libur", detail)
	}

	return libur, nil
}

func (u *AdminUsecase) UpdateLibur(ctx context.Context, userID string, id string, req domain.UpdateLiburRequest) (*domain.Libur, error) {
	libur, err := u.LiburUsecase.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	if u.LogUsecase != nil {
		detail := "Memperbarui hari libur ID: " + id
		u.LogUsecase.CreateLog(ctx, userID, "update_libur", detail)
	}

	return libur, nil
}

func (u *AdminUsecase) DeleteLibur(ctx context.Context, userID string, id string) error {
	if err := u.LiburUsecase.Delete(ctx, id); err != nil {
		return err
	}

	if u.LogUsecase != nil {
		detail := "Menghapus hari libur ID: " + id
		u.LogUsecase.CreateLog(ctx, userID, "delete_libur", detail)
	}

	return nil
}

func (u *AdminUsecase) ToggleLibur(ctx context.Context, userID string, id string) (bool, error) {
	aktif, err := u.LiburUsecase.Toggle(ctx, id)
	if err != nil {
		return false, err
	}

	if u.LogUsecase != nil {
		statusText := "nonaktif"
		if aktif {
			statusText = "aktif"
		}
		detail := "Mengubah status hari libur ID: " + id + " menjadi " + statusText
		u.LogUsecase.CreateLog(ctx, userID, "toggle_libur", detail)
	}

	return aktif, nil
}

func (u *AdminUsecase) GetAllLibur(ctx context.Context, tahun int, jenis, sumber string, aktif *bool, limit, page int) ([]domain.Libur, int, error) {
	return u.LiburUsecase.GetAll(ctx, tahun, jenis, sumber, aktif, limit, page)
}

func (u *AdminUsecase) GetKonfigurasi(ctx context.Context) (*domain.KonfigurasiKerja, error) {
	config, err := u.KonfigurasiRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("konfigurasi kerja tidak ditemukan")
	}
	return config, nil
}

func (u *AdminUsecase) UpdateKonfigurasi(ctx context.Context, userID string, req domain.UpdateKonfigurasiRequest) (*domain.KonfigurasiKerja, error) {
	config, err := u.KonfigurasiRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("konfigurasi kerja tidak ditemukan")
	}

	if req.NamaKantor != "" {
		config.NamaKantor = req.NamaKantor
	}
	if req.LatKantor != 0 {
		config.LatKantor = req.LatKantor
	}
	if req.LongKantor != 0 {
		config.LongKantor = req.LongKantor
	}
	if req.JamMasuk != "" {
		config.JamMasuk = req.JamMasuk
	}
	if req.JamMinimalMasuk != "" {
		config.JamMinimalMasuk = req.JamMinimalMasuk
	}
	if req.JamPulang != "" {
		config.JamPulang = req.JamPulang
	}
	if req.JamMinimalPulang != "" {
		config.JamMinimalPulang = req.JamMinimalPulang
	}
	if req.RadiusKantor != 0 {
		config.RadiusKantor = req.RadiusKantor
	}
	if req.LogoKantor != nil {
		config.LogoKantor = req.LogoKantor
	}

	config.DiperbaruiOleh = &userID

	if err := u.KonfigurasiRepo.Update(ctx, config); err != nil {
		return nil, err
	}

	if u.LogUsecase != nil {
		detail := "Memperbarui konfigurasi kerja oleh user ID: " + userID
		u.LogUsecase.CreateLog(ctx, userID, "update_konfigurasi", detail)
	}

	return config, nil
}

func (u *AdminUsecase) UpdateCutiBalance(ctx context.Context, adminID, karyawanID string, tahun int, sisaCutiBaru int) error {
	balance, err := u.LeaveRepo.GetBalance(ctx, karyawanID, tahun)
	if err != nil {
		return err
	}
	if balance == nil {
		balance = &domain.SisaCuti{
			Tahun:            tahun,
			JumlahCuti:       12,
			TelahDilaksanakan: 0,
			AkanDilaksanakan:  0,
			SisaCuti:          sisaCutiBaru,
		}
	}

	balance.SisaCuti = sisaCutiBaru

	if err := u.LeaveRepo.UpdateBalance(ctx, karyawanID, tahun); err != nil {
		return err
	}

	if u.LogUsecase != nil {
		detail := "Mengubah sisa cuti karyawan " + karyawanID + " tahun " + strconv.Itoa(tahun) + " menjadi " + strconv.Itoa(sisaCutiBaru)
		u.LogUsecase.CreateLog(ctx, adminID, "update_cuti_balance", detail)
	}

	return nil
}

func (u *AdminUsecase) ExportKaryawanCSV(ctx context.Context, search, role, status string) ([]byte, error) {
	items, _, err := u.GetAllKaryawan(ctx, 1, 10000, search, role, "", "", "", status)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	buf.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(&buf)
	writer.Comma = '\t'

	writer.Write([]string{"ID", "Nama Lengkap", "Email", "Nomor Telepon", "Role", "Jabatan", "Divisi", "Unit", "Status"})

	for _, item := range items {
		jabatan := ""
		if item.LevelJabatan != nil {
			jabatan = *item.LevelJabatan
		}
		divisi := ""
		if item.Divisi != nil {
			divisi = *item.Divisi
		}
		unit := ""
		if item.Unit != nil {
			unit = *item.Unit
		}
		noTelp := ""
		if item.NomorTelepon != nil {
			noTelp = *item.NomorTelepon
		}

		row := []string{
			item.ID,
			item.NamaLengkap,
			item.Email,
			noTelp,
			item.Role,
			jabatan,
			divisi,
			unit,
			item.StatusKaryawan,
		}
		writer.Write(row)
	}
	writer.Flush()
	return buf.Bytes(), nil
}

func (u *AdminUsecase) GetLogs(ctx context.Context, page, limit int) ([]map[string]interface{}, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}
	offset := (page - 1) * limit
	return u.LogUsecase.GetLogs(ctx, limit, offset)
}
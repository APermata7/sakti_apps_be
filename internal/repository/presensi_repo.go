package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"sakti_apps_be/internal/domain"
)

type PresensiRepo struct {
	DB *pgxpool.Pool
}

func NewPresensiRepo(db *pgxpool.Pool) *PresensiRepo {
	return &PresensiRepo{DB: db}
}

func (r *PresensiRepo) Create(ctx context.Context, p *domain.Presensi) error {
	query := `
		INSERT INTO presensi (
			karyawan_id, tanggal, jam_masuk, status,
			lintang_masuk, bujur_masuk,
			validasi_wajah, url_foto, alasan_terlambat,
			distance_meter, is_outside_radius, location_status,
			dibuat_pada, diperbarui_pada
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
		RETURNING id
	`

	var id string
	err := r.DB.QueryRow(ctx, query,
		p.KaryawanID,
		p.Tanggal,
		p.JamMasuk,
		p.Status,
		p.LintangMasuk,
		p.BujurMasuk,
		p.ValidasiWajah,
		p.URLFoto,
		p.AlasanTerlambat,
		p.DistanceMeter,
		p.IsOutsideRadius,
		p.LocationStatus,
	).Scan(&id)

	if err != nil {
		return err
	}

	p.ID = id
	return nil
}

func (r *PresensiRepo) GetToday(ctx context.Context, karyawanID string) (*domain.Presensi, error) {
	query := `
		SELECT id, karyawan_id, tanggal, jam_masuk, jam_keluar, status,
		       lintang_masuk, bujur_masuk, lintang_keluar, bujur_keluar,
		       validasi_wajah, url_foto, alasan_terlambat, lembur, jam_lembur,
		       dibuat_pada, diperbarui_pada
		FROM presensi
		WHERE karyawan_id = $1 AND tanggal = CURRENT_DATE
	`

	var p domain.Presensi
	err := r.DB.QueryRow(ctx, query, karyawanID).Scan(
		&p.ID, &p.KaryawanID, &p.Tanggal, &p.JamMasuk, &p.JamKeluar,
		&p.Status, &p.LintangMasuk, &p.BujurMasuk, &p.LintangKeluar,
		&p.BujurKeluar, &p.ValidasiWajah, &p.URLFoto, &p.AlasanTerlambat,
		&p.Lembur, &p.JamLembur, &p.DibuatPada, &p.DiperbaruiPada,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PresensiRepo) UpdateCheckOut(ctx context.Context, id string, jamKeluar string, lembur bool, jamLembur float64, lat, lon float64) error {
	query := `
		UPDATE presensi 
		SET jam_keluar = $2, lembur = $3, jam_lembur = $4,
		    lintang_keluar = $5, bujur_keluar = $6,
		    diperbarui_pada = NOW()
		WHERE id = $1
	`

	_, err := r.DB.Exec(ctx, query, id, jamKeluar, lembur, jamLembur, lat, lon)
	return err
}

func (r *PresensiRepo) UpdateAlasanTerlambat(ctx context.Context, id string, alasan string) error {
	query := `UPDATE presensi SET alasan_terlambat = $1, diperbarui_pada = NOW() WHERE id = $2`
	_, err := r.DB.Exec(ctx, query, alasan, id)
	return err
}

func (r *PresensiRepo) GetHistory(ctx context.Context, karyawanID, startDate, endDate, status string, limit, offset int) ([]domain.Presensi, int, error) {
	var items []domain.Presensi
	var total int

	baseQuery := `
		FROM presensi
		WHERE karyawan_id = $1
	`

	args := []interface{}{karyawanID}
	argIdx := 2

	if startDate != "" {
		baseQuery += ` AND tanggal >= $` + string(rune(argIdx+'0'))
		args = append(args, startDate)
		argIdx++
	}
	if endDate != "" {
		baseQuery += ` AND tanggal <= $` + string(rune(argIdx+'0'))
		args = append(args, endDate)
		argIdx++
	}
	if status != "" {
		baseQuery += ` AND status = $` + string(rune(argIdx+'0'))
		args = append(args, status)
		argIdx++
	}

	countQuery := `SELECT COUNT(*) ` + baseQuery
	err := r.DB.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	dataQuery := `
		SELECT id, karyawan_id, tanggal, jam_masuk, jam_keluar, status,
		       lintang_masuk, bujur_masuk, lintang_keluar, bujur_keluar,
		       validasi_wajah, url_foto, alasan_terlambat, lembur, jam_lembur,
		       dibuat_pada, diperbarui_pada
	` + baseQuery + ` ORDER BY tanggal DESC LIMIT $` + string(rune(argIdx+'0')) + ` OFFSET $` + string(rune(argIdx+1+'0'))

	finalArgs := append(args, limit, offset)
	rows, err := r.DB.Query(ctx, dataQuery, finalArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.Presensi
		err := rows.Scan(
			&p.ID, &p.KaryawanID, &p.Tanggal, &p.JamMasuk, &p.JamKeluar,
			&p.Status, &p.LintangMasuk, &p.BujurMasuk, &p.LintangKeluar,
			&p.BujurKeluar, &p.ValidasiWajah, &p.URLFoto, &p.AlasanTerlambat,
			&p.Lembur, &p.JamLembur, &p.DibuatPada, &p.DiperbaruiPada,
		)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, p)
	}

	return items, total, nil
}
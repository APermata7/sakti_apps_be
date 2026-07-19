package repository

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5"
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
	log.Printf("Create dimulai untuk karyawanID: %s", p.KaryawanID)

	query := `
		INSERT INTO presensi (
			karyawan_id, kantor_id, tanggal, jam_masuk, status,
			lintang_masuk, bujur_masuk,
			validasi_wajah, face_similarity, url_foto, alasan_terlambat,
			distance_meter, is_outside_radius, location_status,
			dibuat_pada, diperbarui_pada
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW(), NOW())
		RETURNING id
	`

	var id string
	err := r.DB.QueryRow(ctx, query,
		p.KaryawanID,
		p.KantorID,
		p.Tanggal,
		p.JamMasuk,
		p.Status,
		p.LintangMasuk,
		p.BujurMasuk,
		p.ValidasiWajah,
		p.FaceSimilarity,
		p.URLFoto,
		p.AlasanTerlambat,
		p.DistanceMeter,
		p.IsOutsideRadius,
		p.LocationStatus,
	).Scan(&id)

	if err != nil {
		log.Printf("Error create: %v", err)
		return err
	}

	p.ID = id
	log.Printf("Create berhasil dengan ID: %s", id)
	return nil
}

func (r *PresensiRepo) GetToday(ctx context.Context, karyawanID string) (*domain.Presensi, error) {
	log.Printf("GetToday dimulai untuk karyawanID: %s", karyawanID)

	query := `
		SELECT id, karyawan_id, kantor_id, tanggal, jam_masuk, jam_keluar, status,
		       lintang_masuk, bujur_masuk, lintang_keluar, bujur_keluar,
		       validasi_wajah, face_similarity, url_foto, alasan_terlambat,
		       lembur, jam_lembur, distance_meter, is_outside_radius, location_status,
		       dibuat_pada, diperbarui_pada
		FROM presensi
		WHERE karyawan_id = $1 AND tanggal = DATE(NOW() AT TIME ZONE 'Asia/Jakarta')
		LIMIT 1
	`

	var p domain.Presensi
	err := r.DB.QueryRow(ctx, query, karyawanID).Scan(
		&p.ID, &p.KaryawanID, &p.KantorID, &p.Tanggal, &p.JamMasuk, &p.JamKeluar,
		&p.Status, &p.LintangMasuk, &p.BujurMasuk, &p.LintangKeluar,
		&p.BujurKeluar, &p.ValidasiWajah, &p.FaceSimilarity, &p.URLFoto,
		&p.AlasanTerlambat, &p.Lembur, &p.JamLembur, &p.DistanceMeter,
		&p.IsOutsideRadius, &p.LocationStatus, &p.DibuatPada, &p.DiperbaruiPada,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("GetToday tidak ada record untuk karyawanID: %s", karyawanID)
			return nil, nil
		}
		log.Printf("Error GetToday: %v", err)
		return nil, err
	}

	log.Printf("GetToday berhasil untuk karyawanID: %s, ID: %s", karyawanID, p.ID)
	return &p, nil
}

func (r *PresensiRepo) UpdateCheckOut(ctx context.Context, id string, jamKeluar string, lembur bool, jamLembur float64, lat, lon float64, selfieURL string) error {
	log.Printf("UpdateCheckOut dimulai untuk ID: %s", id)

	query := `
		UPDATE presensi 
		SET jam_keluar = $2, 
		    lembur = $3, 
		    jam_lembur = $4,
		    lintang_keluar = $5, 
		    bujur_keluar = $6,
		    url_foto = $7,
		    diperbarui_pada = NOW()
		WHERE id = $1
	`

	result, err := r.DB.Exec(ctx, query, id, jamKeluar, lembur, jamLembur, lat, lon, selfieURL)
	if err != nil {
		log.Printf("Error UpdateCheckOut: %v", err)
		return err
	}

	rowsAffected := result.RowsAffected()
	log.Printf("UpdateCheckOut baris terpengaruh: %d", rowsAffected)

	if rowsAffected == 0 {
		return errors.New("tidak ada data presensi yang ditemukan untuk diupdate")
	}

	log.Printf("UpdateCheckOut berhasil untuk ID: %s", id)
	return nil
}

func (r *PresensiRepo) UpdateAlasanTerlambat(ctx context.Context, id string, alasan string) error {
	log.Printf("UpdateAlasanTerlambat dimulai untuk ID: %s", id)

	query := `UPDATE presensi SET alasan_terlambat = $1, diperbarui_pada = NOW() WHERE id = $2`

	result, err := r.DB.Exec(ctx, query, alasan, id)
	if err != nil {
		log.Printf("Error UpdateAlasanTerlambat: %v", err)
		return err
	}

	rowsAffected := result.RowsAffected()
	log.Printf("UpdateAlasanTerlambat baris terpengaruh: %d", rowsAffected)

	if rowsAffected == 0 {
		return errors.New("tidak ada data presensi yang ditemukan untuk diupdate")
	}

	log.Printf("UpdateAlasanTerlambat berhasil untuk ID: %s", id)
	return nil
}

func (r *PresensiRepo) GetHistory(ctx context.Context, karyawanID, startDate, endDate, status string, limit, offset int) ([]domain.Presensi, int, error) {
	log.Printf("GetHistory dimulai untuk karyawanID: %s", karyawanID)

	var items []domain.Presensi
	var total int

	baseQuery := `
		FROM presensi
		WHERE karyawan_id = $1
	`

	args := []interface{}{karyawanID}
	argIdx := 2

	if startDate != "" {
		baseQuery += ` AND tanggal >= $` + strconv.Itoa(argIdx)
		args = append(args, startDate)
		argIdx++
	}
	if endDate != "" {
		baseQuery += ` AND tanggal <= $` + strconv.Itoa(argIdx)
		args = append(args, endDate)
		argIdx++
	}
	if status != "" {
		baseQuery += ` AND status = $` + strconv.Itoa(argIdx)
		args = append(args, status)
		argIdx++
	}

	countQuery := `SELECT COUNT(*) ` + baseQuery
	err := r.DB.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		log.Printf("Error GetHistory count: %v", err)
		return nil, 0, err
	}
	log.Printf("GetHistory total data: %d", total)

	dataQuery := `
		SELECT id, karyawan_id, kantor_id, tanggal, jam_masuk, jam_keluar, status,
		       lintang_masuk, bujur_masuk, lintang_keluar, bujur_keluar,
		       validasi_wajah, face_similarity, url_foto, alasan_terlambat,
		       lembur, jam_lembur, distance_meter, is_outside_radius, location_status,
		       dibuat_pada, diperbarui_pada
	` + baseQuery + ` ORDER BY tanggal DESC LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)

	finalArgs := append(args, limit, offset)
	rows, err := r.DB.Query(ctx, dataQuery, finalArgs...)
	if err != nil {
		log.Printf("Error GetHistory query: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.Presensi
		err := rows.Scan(
			&p.ID, &p.KaryawanID, &p.KantorID, &p.Tanggal, &p.JamMasuk, &p.JamKeluar,
			&p.Status, &p.LintangMasuk, &p.BujurMasuk, &p.LintangKeluar,
			&p.BujurKeluar, &p.ValidasiWajah, &p.FaceSimilarity, &p.URLFoto,
			&p.AlasanTerlambat, &p.Lembur, &p.JamLembur, &p.DistanceMeter,
			&p.IsOutsideRadius, &p.LocationStatus, &p.DibuatPada, &p.DiperbaruiPada,
		)
		if err != nil {
			log.Printf("Error GetHistory scan: %v", err)
			return nil, 0, err
		}
		items = append(items, p)
	}

	log.Printf("GetHistory berhasil, items: %d", len(items))
	return items, total, nil
}
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"sakti_apps_be/internal/domain"
)

type RiwayatRepo struct {
	DB *pgxpool.Pool
}

func NewRiwayatRepo(db *pgxpool.Pool) *RiwayatRepo {
	return &RiwayatRepo{DB: db}
}

func (r *RiwayatRepo) CreateRiwayat(ctx context.Context, karyawanID, action, detail string) error {
	query := `
		INSERT INTO riwayat_user (karyawan_id, action, detail, dibuat_pada)
		VALUES ($1, $2, $3, NOW())
	`
	_, err := r.DB.Exec(ctx, query, karyawanID, action, detail)
	return err
}

func (r *RiwayatRepo) GetRiwayat(ctx context.Context, karyawanID string, limit, offset int) ([]domain.RiwayatItem, int, error) {
	var items []domain.RiwayatItem
	var total int

	presensiQuery := `
		SELECT 
			'presensi' AS tipe,
			id,
			tanggal::text,
			jam_masuk,
			jam_keluar,
			status,
			0 AS total_hari,
			lembur,
			jam_lembur,
			url_foto,
			dibuat_pada::text
		FROM presensi
		WHERE karyawan_id = $1
	`
	cutiQuery := `
		SELECT 
			'cuti' AS tipe,
			id,
			tanggal_mulai::text,
			NULL AS jam_masuk,
			NULL AS jam_keluar,
			status,
			total_hari,
			false AS lembur,
			0 AS jam_lembur,
			'' AS url_foto,
			dibuat_pada::text
		FROM pengajuan_cuti
		WHERE karyawan_id = $1
	`
	riwayatUserQuery := `
		SELECT 
			'aktivitas' AS tipe,
			id,
			dibuat_pada::text AS tanggal,
			NULL AS jam_masuk,
			NULL AS jam_keluar,
			action AS status,
			0 AS total_hari,
			false AS lembur,
			0 AS jam_lembur,
			'' AS url_foto,
			dibuat_pada::text
		FROM riwayat_user
		WHERE karyawan_id = $1
	`
	finalQuery := `
		SELECT * FROM (
			` + presensiQuery + `
			UNION ALL
			` + cutiQuery + `
			UNION ALL
			` + riwayatUserQuery + `
		) AS riwayat
		ORDER BY tanggal DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.Query(ctx, finalQuery, karyawanID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.RiwayatItem
		err := rows.Scan(
			&item.Tipe,
			&item.ID,
			&item.Tanggal,
			&item.JamMasuk,
			&item.JamKeluar,
			&item.Status,
			&item.TotalHari,
			&item.Lembur,
			&item.JamLembur,
			&item.URLFoto,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	countQuery := `
		SELECT COUNT(*) FROM (
			SELECT 1 FROM presensi WHERE karyawan_id = $1
			UNION ALL
			SELECT 1 FROM pengajuan_cuti WHERE karyawan_id = $1
			UNION ALL
			SELECT 1 FROM riwayat_user WHERE karyawan_id = $1
		) AS total
	`
	err = r.DB.QueryRow(ctx, countQuery, karyawanID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"sakti_apps_be/internal/domain"
)

type NotifikasiRepo struct {
	DB *pgxpool.Pool
}

func NewNotifikasiRepo(db *pgxpool.Pool) *NotifikasiRepo {
	return &NotifikasiRepo{DB: db}
}

func (r *NotifikasiRepo) Create(ctx context.Context, n *domain.Notifikasi) error {
	query := `
		INSERT INTO notifikasi (
			karyawan_id, jenis, channel, judul, pesan,
			dibaca, referensi_id, referensi_tipe, dibuat_pada
		) VALUES ($1, $2, $3, $4, $5, false, $6, $7, NOW())
		RETURNING id
	`
	err := r.DB.QueryRow(ctx, query,
		n.KaryawanID,
		n.Jenis,
		n.Channel,
		n.Judul,
		n.Pesan,
		n.ReferensiID,
		n.ReferensiTipe,
	).Scan(&n.ID)
	return err
}

func (r *NotifikasiRepo) GetByKaryawanID(ctx context.Context, karyawanID string, limit, offset int) ([]domain.Notifikasi, int, error) {
	var items []domain.Notifikasi
	var total int

	query := `
		SELECT id, karyawan_id, jenis, channel, judul, pesan,
		       dibaca, dibaca_pada, referensi_id, referensi_tipe, dibuat_pada
		FROM notifikasi
		WHERE karyawan_id = $1
		ORDER BY dibuat_pada DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.Query(ctx, query, karyawanID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var n domain.Notifikasi
		var dibacaPada *time.Time
		err := rows.Scan(
			&n.ID,
			&n.KaryawanID,
			&n.Jenis,
			&n.Channel,
			&n.Judul,
			&n.Pesan,
			&n.Dibaca,
			&dibacaPada,
			&n.ReferensiID,
			&n.ReferensiTipe,
			&n.DibuatPada,
		)
		if err != nil {
			return nil, 0, err
		}
		if dibacaPada != nil {
			n.DibacaPada = dibacaPada
		}
		items = append(items, n)
	}

	countQuery := `SELECT COUNT(*) FROM notifikasi WHERE karyawan_id = $1`
	err = r.DB.QueryRow(ctx, countQuery, karyawanID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *NotifikasiRepo) MarkAsRead(ctx context.Context, id, karyawanID string) error {
	query := `
		UPDATE notifikasi 
		SET dibaca = true, dibaca_pada = NOW()
		WHERE id = $1 AND karyawan_id = $2
	`
	_, err := r.DB.Exec(ctx, query, id, karyawanID)
	return err
}

func (r *NotifikasiRepo) MarkAllAsRead(ctx context.Context, karyawanID string) error {
	query := `
		UPDATE notifikasi 
		SET dibaca = true, dibaca_pada = NOW()
		WHERE karyawan_id = $1 AND dibaca = false
	`
	_, err := r.DB.Exec(ctx, query, karyawanID)
	return err
}

func (r *NotifikasiRepo) GetUnreadCount(ctx context.Context, karyawanID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifikasi WHERE karyawan_id = $1 AND dibaca = false`
	err := r.DB.QueryRow(ctx, query, karyawanID).Scan(&count)
	return count, err
}

func (r *NotifikasiRepo) GetByReferensi(ctx context.Context, referensiID, referensiTipe string) ([]domain.Notifikasi, error) {
	query := `
		SELECT id, karyawan_id, jenis, channel, judul, pesan,
		       dibaca, dibaca_pada, referensi_id, referensi_tipe, dibuat_pada
		FROM notifikasi
		WHERE referensi_id = $1 AND referensi_tipe = $2
		ORDER BY dibuat_pada DESC
	`
	rows, err := r.DB.Query(ctx, query, referensiID, referensiTipe)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Notifikasi
	for rows.Next() {
		var n domain.Notifikasi
		var dibacaPada *time.Time
		err := rows.Scan(
			&n.ID,
			&n.KaryawanID,
			&n.Jenis,
			&n.Channel,
			&n.Judul,
			&n.Pesan,
			&n.Dibaca,
			&dibacaPada,
			&n.ReferensiID,
			&n.ReferensiTipe,
			&n.DibuatPada,
		)
		if err != nil {
			return nil, err
		}
		if dibacaPada != nil {
			n.DibacaPada = dibacaPada
		}
		items = append(items, n)
	}
	return items, nil
}
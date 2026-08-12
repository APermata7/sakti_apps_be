package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sakti_apps_be/internal/domain"
)

type TTDRepo struct {
	DB *pgxpool.Pool
}

func NewTTDRepo(db *pgxpool.Pool) *TTDRepo {
	return &TTDRepo{DB: db}
}

func (r *TTDRepo) Create(ctx context.Context, ttd *domain.TandaTangan) error {
	query := `
		INSERT INTO tanda_tangan (karyawan_id, url_tanda_tangan, hash_tanda_tangan)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.DB.QueryRow(ctx, query, ttd.KaryawanID, ttd.URLTandaTangan, ttd.HashTandaTangan).Scan(&ttd.ID)
	return err
}

func (r *TTDRepo) GetByKaryawanID(ctx context.Context, karyawanID string) (*domain.TandaTangan, error) {
	query := `
		SELECT id, karyawan_id, url_tanda_tangan, hash_tanda_tangan, diunggah_pada, diperbarui_pada
		FROM tanda_tangan
		WHERE karyawan_id = $1
	`
	var ttd domain.TandaTangan
	err := r.DB.QueryRow(ctx, query, karyawanID).Scan(
		&ttd.ID,
		&ttd.KaryawanID,
		&ttd.URLTandaTangan,
		&ttd.HashTandaTangan,
		&ttd.DiunggahPada,
		&ttd.DiperbaruiPada,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ttd, nil
}

func (r *TTDRepo) Update(ctx context.Context, ttd *domain.TandaTangan) error {
	query := `
		UPDATE tanda_tangan
		SET url_tanda_tangan = $2, hash_tanda_tangan = $3, diperbarui_pada = NOW()
		WHERE karyawan_id = $1
	`
	_, err := r.DB.Exec(ctx, query, ttd.KaryawanID, ttd.URLTandaTangan, ttd.HashTandaTangan)
	return err
}

func (r *TTDRepo) Delete(ctx context.Context, karyawanID string) error {
	query := `DELETE FROM tanda_tangan WHERE karyawan_id = $1`
	_, err := r.DB.Exec(ctx, query, karyawanID)
	return err
}
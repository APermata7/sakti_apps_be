package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sakti_apps_be/internal/domain"
)

type KantorRepo struct {
	DB *pgxpool.Pool
}

func NewKantorRepo(db *pgxpool.Pool) *KantorRepo {
	return &KantorRepo{DB: db}
}

func (r *KantorRepo) GetByID(ctx context.Context, id string) (*domain.Kantor, error) {
	query := `
		SELECT id, nama, lintang, bujur, radius, alamat, dibuat_pada, diperbarui_pada
		FROM kantor
		WHERE id = $1
	`

	var k domain.Kantor
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&k.ID, &k.Nama, &k.Lintang, &k.Bujur, &k.Radius, &k.Alamat,
		&k.DibuatPada, &k.DiperbaruiPada,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &k, nil
}

func (r *KantorRepo) GetByNama(ctx context.Context, nama string) (*domain.Kantor, error) {
	query := `
		SELECT id, nama, lintang, bujur, radius, alamat, dibuat_pada, diperbarui_pada
		FROM kantor
		WHERE nama = $1
		LIMIT 1
	`

	var k domain.Kantor
	err := r.DB.QueryRow(ctx, query, nama).Scan(
		&k.ID, &k.Nama, &k.Lintang, &k.Bujur, &k.Radius, &k.Alamat,
		&k.DibuatPada, &k.DiperbaruiPada,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &k, nil
}
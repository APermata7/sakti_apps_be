package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sakti_apps_be/internal/domain"
)

type LiburRepo struct {
	DB *pgxpool.Pool
}

func NewLiburRepo(db *pgxpool.Pool) *LiburRepo {
	return &LiburRepo{DB: db}
}

func (r *LiburRepo) IsHoliday(ctx context.Context, tanggal string) (bool, error) {
	query := `SELECT COUNT(*) FROM libur WHERE tanggal = $1 AND aktif = true`
	var count int
	err := r.DB.QueryRow(ctx, query, tanggal).Scan(&count)
	return count > 0, err
}

func (r *LiburRepo) Create(ctx context.Context, libur *domain.Libur) error {
	query := `
		INSERT INTO libur (tanggal, nama, jenis, sumber, aktif, dibuat_pada, diperbarui_pada)
		VALUES ($1, $2, $3, 'manual', true, NOW(), NOW())
		RETURNING id
	`
	err := r.DB.QueryRow(ctx, query, libur.Tanggal, libur.Nama, libur.Jenis).Scan(&libur.ID)
	return err
}

func (r *LiburRepo) GetByID(ctx context.Context, id string) (*domain.Libur, error) {
	query := `
		SELECT id, tanggal, nama, jenis, aktif, sumber, dibuat_pada, diperbarui_pada
		FROM libur
		WHERE id = $1
	`
	var l domain.Libur
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&l.ID, &l.Tanggal, &l.Nama, &l.Jenis, &l.Aktif,
		&l.Sumber, &l.DibuatPada, &l.DiperbaruiPada,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &l, nil
}

func (r *LiburRepo) GetAll(ctx context.Context, tahun string) ([]domain.Libur, error) {
	query := `
		SELECT id, tanggal, nama, jenis, aktif, sumber, dibuat_pada, diperbarui_pada
		FROM libur
		WHERE EXTRACT(YEAR FROM tanggal) = $1
		ORDER BY tanggal ASC
	`
	rows, err := r.DB.Query(ctx, query, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Libur
	for rows.Next() {
		var l domain.Libur
		err := rows.Scan(
			&l.ID, &l.Tanggal, &l.Nama, &l.Jenis, &l.Aktif,
			&l.Sumber, &l.DibuatPada, &l.DiperbaruiPada,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, l)
	}
	return items, nil
}

func (r *LiburRepo) Update(ctx context.Context, libur *domain.Libur) error {
	query := `
		UPDATE libur 
		SET nama = $2, jenis = $3, aktif = $4, diperbarui_pada = NOW()
		WHERE id = $1
	`
	_, err := r.DB.Exec(ctx, query, libur.ID, libur.Nama, libur.Jenis, libur.Aktif)
	return err
}

func (r *LiburRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM libur WHERE id = $1`
	_, err := r.DB.Exec(ctx, query, id)
	return err
}
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
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
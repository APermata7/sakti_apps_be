package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SupabaseDB struct {
	Pool *pgxpool.Pool
}

func NewSupabaseDB() (*SupabaseDB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL tidak ditemukan di .env")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("gagal parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("gagal koneksi: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("gagal ping database: %w", err)
	}

	log.Println("Berhasil terhubung ke Supabase")
	return &SupabaseDB{Pool: pool}, nil
}

func (db *SupabaseDB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		log.Println("Koneksi database ditutup")
	}
}

func (db *SupabaseDB) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

func (db *SupabaseDB) GetPool() *pgxpool.Pool {
	return db.Pool
}
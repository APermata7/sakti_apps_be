package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sakti_apps_be/internal/domain"
)

type KaryawanRepo struct {
	DB *pgxpool.Pool
}

func NewKaryawanRepo(db *pgxpool.Pool) *KaryawanRepo {
	return &KaryawanRepo{DB: db}
}

func (r *KaryawanRepo) GetByID(ctx context.Context, id string) (*domain.Karyawan, error) {
	query := `
		SELECT id, nama_lengkap, email, nomor_telepon, foto_url, 
		       peran, level_jabatan, atasan_langsung_id, 
		       divisi, unit, status_karyawan, dibuat_pada, diperbarui_pada
		FROM karyawan
		WHERE id = $1
	`

	var k domain.Karyawan
	var atasanID *string

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&k.ID,
		&k.NamaLengkap,
		&k.Email,
		&k.NomorTelepon,
		&k.FotoURL,
		&k.Peran,
		&k.LevelJabatan,
		&atasanID,
		&k.Divisi,
		&k.Unit,
		&k.StatusKaryawan,
		&k.DibuatPada,
		&k.DiperbaruiPada,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	k.AtasanLangsungID = atasanID
	return &k, nil
}

func (r *KaryawanRepo) GetByEmail(ctx context.Context, email string) (*domain.Karyawan, error) {
	query := `
		SELECT id, nama_lengkap, email, nomor_telepon, foto_url, 
		       peran, level_jabatan, atasan_langsung_id, 
		       divisi, unit, status_karyawan, dibuat_pada, diperbarui_pada
		FROM karyawan
		WHERE email = $1
	`

	var k domain.Karyawan
	var atasanID *string

	err := r.DB.QueryRow(ctx, query, email).Scan(
		&k.ID,
		&k.NamaLengkap,
		&k.Email,
		&k.NomorTelepon,
		&k.FotoURL,
		&k.Peran,
		&k.LevelJabatan,
		&atasanID,
		&k.Divisi,
		&k.Unit,
		&k.StatusKaryawan,
		&k.DibuatPada,
		&k.DiperbaruiPada,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	k.AtasanLangsungID = atasanID
	return &k, nil
}

func (r *KaryawanRepo) Create(ctx context.Context, k *domain.Karyawan) error {
	query := `
		INSERT INTO karyawan (id, nama_lengkap, email, nomor_telepon, foto_url, 
		                     peran, level_jabatan, atasan_langsung_id, 
		                     divisi, unit, status_karyawan, dibuat_pada, diperbarui_pada)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
	`

	_, err := r.DB.Exec(ctx, query,
		k.ID,
		k.NamaLengkap,
		k.Email,
		k.NomorTelepon,
		k.FotoURL,
		k.Peran,
		k.LevelJabatan,
		k.AtasanLangsungID,
		k.Divisi,
		k.Unit,
		k.StatusKaryawan,
	)

	return err
}

func (r *KaryawanRepo) Update(ctx context.Context, k *domain.Karyawan) error {
	query := `
		UPDATE karyawan 
		SET nama_lengkap = $2, nomor_telepon = $3, foto_url = $4,
		    peran = $5, level_jabatan = $6, atasan_langsung_id = $7,
		    divisi = $8, unit = $9, status_karyawan = $10, diperbarui_pada = NOW()
		WHERE id = $1
	`

	_, err := r.DB.Exec(ctx, query,
		k.ID,
		k.NamaLengkap,
		k.NomorTelepon,
		k.FotoURL,
		k.Peran,
		k.LevelJabatan,
		k.AtasanLangsungID,
		k.Divisi,
		k.Unit,
		k.StatusKaryawan,
	)

	return err
}

func (r *KaryawanRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE karyawan SET status_karyawan = 'nonaktif', diperbarui_pada = NOW() WHERE id = $1`
	_, err := r.DB.Exec(ctx, query, id)
	return err
}
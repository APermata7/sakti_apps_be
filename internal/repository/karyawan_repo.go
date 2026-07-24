package repository

import (
	"context"
	"errors"
	"fmt"

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

func getDepartemen(divisi *string, unit *string) string {
	divisiVal := ""
	unitVal := ""
	if divisi != nil {
		divisiVal = *divisi
	}
	if unit != nil {
		unitVal = *unit
	}

	if divisiVal != "" && unitVal != "" && divisiVal != unitVal {
		return divisiVal + " - " + unitVal
	} else if divisiVal != "" {
		return divisiVal
	} else if unitVal != "" {
		return unitVal
	}
	return "-"
}

func (r *KaryawanRepo) GetByID(ctx context.Context, id string) (*domain.Karyawan, error) {
	query := `
		SELECT id, nama_lengkap, email, nomor_telepon, foto_url, 
		       role, level_jabatan, atasan_langsung_id, 
		       divisi, unit, status_karyawan, dibuat_pada, diperbarui_pada
		FROM karyawan
		WHERE id = $1
	`

	var k domain.Karyawan

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&k.ID,
		&k.NamaLengkap,
		&k.Email,
		&k.NomorTelepon,
		&k.FotoURL,
		&k.Role,
		&k.LevelJabatan,
		&k.AtasanLangsungID,
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

	k.Departemen = getDepartemen(k.Divisi, k.Unit)

	return &k, nil
}

func (r *KaryawanRepo) GetByEmail(ctx context.Context, email string) (*domain.Karyawan, error) {
	query := `
		SELECT id, nama_lengkap, email, nomor_telepon, foto_url, 
		       role, level_jabatan, atasan_langsung_id, 
		       divisi, unit, status_karyawan, dibuat_pada, diperbarui_pada
		FROM karyawan
		WHERE email = $1
	`

	var k domain.Karyawan

	err := r.DB.QueryRow(ctx, query, email).Scan(
		&k.ID,
		&k.NamaLengkap,
		&k.Email,
		&k.NomorTelepon,
		&k.FotoURL,
		&k.Role,
		&k.LevelJabatan,
		&k.AtasanLangsungID,
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

	k.Departemen = getDepartemen(k.Divisi, k.Unit)

	return &k, nil
}

func (r *KaryawanRepo) GetByRole(ctx context.Context, role string) (*domain.Karyawan, error) {
	query := `
		SELECT id, nama_lengkap, email, nomor_telepon, foto_url, 
		       role, level_jabatan, atasan_langsung_id, 
		       divisi, unit, status_karyawan, dibuat_pada, diperbarui_pada
		FROM karyawan
		WHERE role = $1 AND status_karyawan = 'aktif'
		LIMIT 1
	`

	var k domain.Karyawan

	err := r.DB.QueryRow(ctx, query, role).Scan(
		&k.ID,
		&k.NamaLengkap,
		&k.Email,
		&k.NomorTelepon,
		&k.FotoURL,
		&k.Role,
		&k.LevelJabatan,
		&k.AtasanLangsungID,
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

	k.Departemen = getDepartemen(k.Divisi, k.Unit)

	return &k, nil
}

func (r *KaryawanRepo) Create(ctx context.Context, k *domain.Karyawan) error {
	query := `
		INSERT INTO karyawan (
			id, nama_lengkap, email, nomor_telepon, foto_url, 
			role, level_jabatan, atasan_langsung_id, 
			divisi, unit, status_karyawan, dibuat_pada, diperbarui_pada
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET
			nama_lengkap = EXCLUDED.nama_lengkap,
			email = EXCLUDED.email,
			nomor_telepon = EXCLUDED.nomor_telepon,
			foto_url = EXCLUDED.foto_url,
			role = EXCLUDED.role,
			level_jabatan = EXCLUDED.level_jabatan,
			atasan_langsung_id = EXCLUDED.atasan_langsung_id,
			divisi = EXCLUDED.divisi,
			unit = EXCLUDED.unit,
			status_karyawan = EXCLUDED.status_karyawan,
			diperbarui_pada = NOW()
	`

	_, err := r.DB.Exec(ctx, query,
		k.ID,
		k.NamaLengkap,
		k.Email,
		k.NomorTelepon,
		k.FotoURL,
		k.Role,
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
		    role = $5, level_jabatan = $6, atasan_langsung_id = $7,
		    divisi = $8, unit = $9, status_karyawan = $10, diperbarui_pada = NOW()
		WHERE id = $1
	`

	_, err := r.DB.Exec(ctx, query,
		k.ID,
		k.NamaLengkap,
		k.NomorTelepon,
		k.FotoURL,
		k.Role,
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

func (r *KaryawanRepo) GetAll(ctx context.Context, limit, offset int, search, role, status string) ([]domain.Karyawan, int, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		whereClause += fmt.Sprintf(" AND (nama_lengkap ILIKE $%d OR email ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	if role != "" {
		whereClause += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, role)
		argIdx++
	}

	if status != "" {
		whereClause += fmt.Sprintf(" AND status_karyawan = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	countQuery := "SELECT COUNT(*) FROM karyawan " + whereClause
	var total int
	err := r.DB.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, nama_lengkap, email, nomor_telepon, foto_url, 
		       role, level_jabatan, atasan_langsung_id, 
		       divisi, unit, status_karyawan, dibuat_pada, diperbarui_pada
		FROM karyawan
	` + whereClause + fmt.Sprintf(" ORDER BY dibuat_pada DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)

	finalArgs := append(args, limit, offset)
	rows, err := r.DB.Query(ctx, query, finalArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var karyawanList []domain.Karyawan
	for rows.Next() {
		var k domain.Karyawan
		err := rows.Scan(
			&k.ID,
			&k.NamaLengkap,
			&k.Email,
			&k.NomorTelepon,
			&k.FotoURL,
			&k.Role,
			&k.LevelJabatan,
			&k.AtasanLangsungID,
			&k.Divisi,
			&k.Unit,
			&k.StatusKaryawan,
			&k.DibuatPada,
			&k.DiperbaruiPada,
		)
		if err != nil {
			return nil, 0, err
		}
		k.Departemen = getDepartemen(k.Divisi, k.Unit)
		karyawanList = append(karyawanList, k)
	}

	return karyawanList, total, nil
}
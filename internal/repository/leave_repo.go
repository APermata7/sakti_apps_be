package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LeaveRepo struct {
	DB *pgxpool.Pool
}

type Leave struct {
	ID               string  `json:"id"`
	KaryawanID       string  `json:"karyawan_id"`
	TipePengajuan    string  `json:"tipe_pengajuan"`
	TanggalMulai     string  `json:"tanggal_mulai"`
	TanggalSelesai   string  `json:"tanggal_selesai"`
	TotalHari        int     `json:"total_hari"`
	Alasan           string  `json:"alasan"`
	Status           string  `json:"status"`
	DisetujuiOleh    *string `json:"disetujui_oleh"`
	DifinalisasiOleh *string `json:"difinalisasi_oleh"`
}

type LeaveBalance struct {
	KuotaTotal int `json:"kuota_total"`
	Digunakan  int `json:"digunakan"`
	Sisa       int `json:"sisa"`
}

func NewLeaveRepo(db *pgxpool.Pool) *LeaveRepo {
	return &LeaveRepo{DB: db}
}

func (r *LeaveRepo) Create(ctx context.Context, leave *Leave) error {
	query := `
		INSERT INTO pengajuan_cuti (
			karyawan_id, tipe_pengajuan, tanggal_mulai, tanggal_selesai,
			total_hari, alasan, status, dibuat_pada, diperbarui_pada
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id
	`

	err := r.DB.QueryRow(ctx, query,
		leave.KaryawanID,
		leave.TipePengajuan,
		leave.TanggalMulai,
		leave.TanggalSelesai,
		leave.TotalHari,
		leave.Alasan,
		leave.Status,
	).Scan(&leave.ID)

	return err
}

func (r *LeaveRepo) GetByID(ctx context.Context, id string) (*Leave, error) {
	query := `
		SELECT id, karyawan_id, tipe_pengajuan, tanggal_mulai, tanggal_selesai,
		       total_hari, alasan, status, disetujui_oleh, difinalisasi_oleh
		FROM pengajuan_cuti
		WHERE id = $1
	`

	var l Leave
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&l.ID, &l.KaryawanID, &l.TipePengajuan, &l.TanggalMulai, &l.TanggalSelesai,
		&l.TotalHari, &l.Alasan, &l.Status, &l.DisetujuiOleh, &l.DifinalisasiOleh,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("pengajuan tidak ditemukan")
		}
		return nil, err
	}
	return &l, nil
}

func (r *LeaveRepo) GetByKaryawanID(ctx context.Context, karyawanID string) ([]Leave, error) {
	query := `
		SELECT id, karyawan_id, tipe_pengajuan, tanggal_mulai, tanggal_selesai,
		       total_hari, alasan, status, disetujui_oleh, difinalisasi_oleh
		FROM pengajuan_cuti
		WHERE karyawan_id = $1
		ORDER BY dibuat_pada DESC
	`

	rows, err := r.DB.Query(ctx, query, karyawanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaves []Leave
	for rows.Next() {
		var l Leave
		err := rows.Scan(
			&l.ID, &l.KaryawanID, &l.TipePengajuan, &l.TanggalMulai, &l.TanggalSelesai,
			&l.TotalHari, &l.Alasan, &l.Status, &l.DisetujuiOleh, &l.DifinalisasiOleh,
		)
		if err != nil {
			return nil, err
		}
		leaves = append(leaves, l)
	}
	return leaves, nil
}

func (r *LeaveRepo) UpdateStatus(ctx context.Context, id, status string) error {
	query := `UPDATE pengajuan_cuti SET status = $1, diperbarui_pada = NOW() WHERE id = $2`
	_, err := r.DB.Exec(ctx, query, status, id)
	return err
}

func (r *LeaveRepo) Approve(ctx context.Context, id, managerID string) error {
	query := `
		UPDATE pengajuan_cuti 
		SET status = 'disetujui', disetujui_oleh = $1, tanggal_disetujui = NOW(),
		    diperbarui_pada = NOW()
		WHERE id = $2
	`
	_, err := r.DB.Exec(ctx, query, managerID, id)
	return err
}

func (r *LeaveRepo) Reject(ctx context.Context, id, managerID, alasan string) error {
	query := `
		UPDATE pengajuan_cuti 
		SET status = 'ditolak', disetujui_oleh = $1, tanggal_disetujui = NOW(),
		    alasan = $2, diperbarui_pada = NOW()
		WHERE id = $3
	`
	_, err := r.DB.Exec(ctx, query, managerID, alasan, id)
	return err
}

func (r *LeaveRepo) Finalize(ctx context.Context, id, hrdID string) error {
	query := `
		UPDATE pengajuan_cuti 
		SET status = 'disetujui', difinalisasi_oleh = $1, tanggal_difinalisasi = NOW(),
		    diperbarui_pada = NOW()
		WHERE id = $2 AND status = 'disetujui'
	`
	_, err := r.DB.Exec(ctx, query, hrdID, id)
	return err
}

func (r *LeaveRepo) GetBalance(ctx context.Context, karyawanID string, year int) (*LeaveBalance, error) {
	query := `
		SELECT kuota_total, digunakan, sisa
		FROM sisa_cuti
		WHERE karyawan_id = $1 AND tahun = $2
	`

	var b LeaveBalance
	err := r.DB.QueryRow(ctx, query, karyawanID, year).Scan(
		&b.KuotaTotal, &b.Digunakan, &b.Sisa,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &LeaveBalance{KuotaTotal: 12, Digunakan: 0, Sisa: 12}, nil
		}
		return nil, err
	}
	return &b, nil
}
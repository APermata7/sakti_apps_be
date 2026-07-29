package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"sakti_apps_be/internal/domain"
)

type FCMTokenRepo struct {
	DB *pgxpool.Pool
}

func NewFCMTokenRepo(db *pgxpool.Pool) *FCMTokenRepo {
	return &FCMTokenRepo{DB: db}
}

func (r *FCMTokenRepo) SaveToken(ctx context.Context, token *domain.FCMToken) error {
	query := `
		INSERT INTO token_fcm (karyawan_id, fcm_token, is_active, dibuat_pada, diperbarui_pada)
		VALUES ($1, $2, true, NOW(), NOW())
		ON CONFLICT (fcm_token) 
		DO UPDATE SET karyawan_id = $1, is_active = true, diperbarui_pada = NOW()
	`
	_, err := r.DB.Exec(ctx, query, token.KaryawanID, token.FCMToken)
	return err
}

func (r *FCMTokenRepo) GetTokensByKaryawanID(ctx context.Context, karyawanID string) ([]string, error) {
	query := `SELECT fcm_token FROM token_fcm WHERE karyawan_id = $1 AND is_active = true`
	rows, err := r.DB.Query(ctx, query, karyawanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			log.Printf("Error scan token: %v", err)
			continue
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (r *FCMTokenRepo) GetTokensByAtasan(ctx context.Context, karyawanID string) ([]string, error) {
	query := `
		SELECT ft.fcm_token 
		FROM token_fcm ft
		JOIN karyawan k ON k.id = ft.karyawan_id
		WHERE k.id = (SELECT atasan_langsung_id FROM karyawan WHERE id = $1)
		AND ft.is_active = true
	`
	rows, err := r.DB.Query(ctx, query, karyawanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			continue
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (r *FCMTokenRepo) GetTokensByHRD(ctx context.Context) ([]string, error) {
	query := `
		SELECT ft.fcm_token 
		FROM token_fcm ft
		JOIN karyawan k ON k.id = ft.karyawan_id
		WHERE k.role = 'hrd' AND ft.is_active = true
	`
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			continue
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (r *FCMTokenRepo) GetTokensByRole(ctx context.Context, role string) ([]string, error) {
	query := `
		SELECT ft.fcm_token 
		FROM token_fcm ft
		JOIN karyawan k ON k.id = ft.karyawan_id
		WHERE k.role = $1 AND ft.is_active = true
	`
	rows, err := r.DB.Query(ctx, query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			continue
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (r *FCMTokenRepo) DeactivateToken(ctx context.Context, token string) error {
	query := `UPDATE token_fcm SET is_active = false, diperbarui_pada = NOW() WHERE fcm_token = $1`
	_, err := r.DB.Exec(ctx, query, token)
	return err
}

func (r *FCMTokenRepo) DeactivateTokensByKaryawan(ctx context.Context, karyawanID string) error {
	query := `UPDATE token_fcm SET is_active = false, diperbarui_pada = NOW() WHERE karyawan_id = $1`
	_, err := r.DB.Exec(ctx, query, karyawanID)
	return err
}
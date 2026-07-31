package repository

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sakti_apps_be/internal/domain"
)

type TelegramRepo struct {
	DB *pgxpool.Pool
}

func NewTelegramRepo(db *pgxpool.Pool) *TelegramRepo {
	return &TelegramRepo{DB: db}
}

func (r *TelegramRepo) SaveVerification(ctx context.Context, code, chatID, username string) error {
	query := `
		INSERT INTO telegram_verification (code, chat_id, username, created_at, expired_at, is_used)
		VALUES ($1, $2, $3, NOW(), NOW() + INTERVAL '5 minutes', false)
	`
	_, err := r.DB.Exec(ctx, query, code, chatID, username)
	if err != nil {
		log.Printf("Failed to save verification code: %v", err)
		return err
	}
	return nil
}

func (r *TelegramRepo) GetVerificationByCode(ctx context.Context, code string) (*domain.TelegramVerification, error) {
	query := `
		SELECT id, code, chat_id, username, karyawan_id, created_at, expired_at, is_used, used_at
		FROM telegram_verification
		WHERE code = $1
	`

	var v domain.TelegramVerification
	err := r.DB.QueryRow(ctx, query, code).Scan(
		&v.ID,
		&v.Code,
		&v.ChatID,
		&v.Username,
		&v.KaryawanID,
		&v.CreatedAt,
		&v.ExpiredAt,
		&v.IsUsed,
		&v.UsedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("Failed to get verification code: %v", err)
		return nil, err
	}
	return &v, nil
}

func (r *TelegramRepo) MarkVerificationAsUsed(ctx context.Context, code, karyawanID string) error {
	query := `
		UPDATE telegram_verification
		SET is_used = true, used_at = NOW(), karyawan_id = $2
		WHERE code = $1 AND is_used = false
	`
	_, err := r.DB.Exec(ctx, query, code, karyawanID)
	if err != nil {
		log.Printf("Failed to mark verification as used: %v", err)
		return err
	}
	return nil
}

func (r *TelegramRepo) DeleteExpiredVerifications(ctx context.Context) error {
	query := `DELETE FROM telegram_verification WHERE expired_at < NOW() OR is_used = true`
	_, err := r.DB.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to delete expired verifications: %v", err)
		return err
	}
	return nil
}
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TelegramRepo struct {
	DB *pgxpool.Pool
}

func NewTelegramRepo(db *pgxpool.Pool) *TelegramRepo {
	return &TelegramRepo{DB: db}
}

func (r *TelegramRepo) SaveVerificationCode(ctx context.Context, chatID, username, code string) error {
	query := `
		INSERT INTO telegram_verification (chat_id, username, code, expired_at, created_at)
		VALUES ($1, $2, $3, NOW() + INTERVAL '5 minutes', NOW())
		ON CONFLICT (code) DO UPDATE
		SET chat_id = $1, username = $2, expired_at = NOW() + INTERVAL '5 minutes', created_at = NOW(), is_used = false
	`
	_, err := r.DB.Exec(ctx, query, chatID, username, code)
	return err
}

func (r *TelegramRepo) DeleteExpiredCodes(ctx context.Context) error {
	query := `DELETE FROM telegram_verification WHERE expired_at < NOW()`
	_, err := r.DB.Exec(ctx, query)
	return err
}

func (r *TelegramRepo) GetVerificationCode(ctx context.Context, chatID string) (string, error) {
	var code string
	query := `SELECT code FROM telegram_verification WHERE chat_id = $1 AND expired_at > NOW() AND is_used = false ORDER BY created_at DESC LIMIT 1`
	err := r.DB.QueryRow(ctx, query, chatID).Scan(&code)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (r *TelegramRepo) DeleteVerificationCode(ctx context.Context, chatID string) error {
	query := `DELETE FROM telegram_verification WHERE chat_id = $1`
	_, err := r.DB.Exec(ctx, query, chatID)
	return err
}
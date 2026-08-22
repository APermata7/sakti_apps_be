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
		INSERT INTO telegram_verification (chat_id, username, code, expires_at, created_at)
		VALUES ($1, $2, $3, NOW() + INTERVAL '5 minutes', NOW())
		ON CONFLICT (chat_id) DO UPDATE
		SET code = $3, expires_at = NOW() + INTERVAL '5 minutes', created_at = NOW()
	`
	_, err := r.DB.Exec(ctx, query, chatID, username, code)
	return err
}

func (r *TelegramRepo) DeleteExpiredCodes(ctx context.Context) error {
	query := `DELETE FROM telegram_verification WHERE expires_at < NOW()`
	_, err := r.DB.Exec(ctx, query)
	return err
}

func (r *TelegramRepo) GetVerificationCode(ctx context.Context, chatID string) (string, error) {
	var code string
	query := `SELECT code FROM telegram_verification WHERE chat_id = $1 AND expires_at > NOW() ORDER BY created_at DESC LIMIT 1`
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
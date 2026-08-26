package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TelegramRepo struct {
	DB *pgxpool.Pool
}

func NewTelegramRepo(db *pgxpool.Pool) *TelegramRepo {
	return &TelegramRepo{DB: db}
}

func (r *TelegramRepo) SaveVerificationCode(ctx context.Context, chatID, username, code string, now, expiredAt time.Time) error {
	query := `
		INSERT INTO telegram_verification (chat_id, username, code, created_at, expired_at, is_used)
		VALUES ($1, $2, $3, $4, $5, false)
		ON CONFLICT (code) DO UPDATE
		SET chat_id = $1, username = $2, created_at = $4, expired_at = $5, is_used = false
	`
	_, err := r.DB.Exec(ctx, query, chatID, username, code, now, expiredAt)
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

func (r *TelegramRepo) GetChatIDByCode(ctx context.Context, code string) (string, error) {
	var chatID string
	now := time.Now().UTC()
	query := `
		SELECT chat_id
		FROM telegram_verification
		WHERE LOWER(code) = LOWER($1)
		  AND is_used = false
		  AND expired_at > $2
	`
	err := r.DB.QueryRow(ctx, query, code, now).Scan(&chatID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("[GetChatIDByCode] code not found or expired: %s", code)
			return "", errors.New("kode tidak ditemukan atau sudah kadaluarsa")
		}
		return "", err
	}
	log.Printf("[GetChatIDByCode] found chat_id=%s for code=%s", chatID, code)
	return chatID, nil
}

func (r *TelegramRepo) MarkCodeAsUsed(ctx context.Context, code string) error {
	query := `
		UPDATE telegram_verification
		SET is_used = true, used_at = $1
		WHERE LOWER(code) = LOWER($2)
		  AND is_used = false
		  AND expired_at > $1
	`
	now := time.Now().UTC()
	result, err := r.DB.Exec(ctx, query, now, code)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("kode verifikasi tidak valid, sudah digunakan, atau sudah kadaluarsa")
	}
	return nil
}

func (r *TelegramRepo) MarkCodeAsUsedWithTx(ctx context.Context, tx pgx.Tx, code string) error {
	query := `
		UPDATE telegram_verification
		SET is_used = true, used_at = $1
		WHERE LOWER(code) = LOWER($2)
		  AND is_used = false
		  AND expired_at > $1
	`
	now := time.Now().UTC()
	result, err := tx.Exec(ctx, query, now, code)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("kode verifikasi tidak valid, sudah digunakan, atau sudah kadaluarsa")
	}
	return nil
}

func (r *TelegramRepo) DeleteVerificationCode(ctx context.Context, chatID string) error {
	query := `DELETE FROM telegram_verification WHERE chat_id = $1`
	_, err := r.DB.Exec(ctx, query, chatID)
	return err
}

func (r *TelegramRepo) SendMessageByKaryawanID(ctx context.Context, karyawanID, message string) error {
	if karyawanID == "" || message == "" {
		return nil
	}

	query := `
		SELECT telegram_chat_id
		FROM karyawan
		WHERE id = $1
		  AND telegram_chat_id IS NOT NULL
		  AND telegram_chat_id != ''
		  AND status_karyawan = 'aktif'
	`
	var chatID string
	err := r.DB.QueryRow(ctx, query, karyawanID).Scan(&chatID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	}

	if chatID == "" {
		return nil
	}

	return r.SendMessage(ctx, chatID, message)
}

func (r *TelegramRepo) SendMessage(ctx context.Context, chatID, message string) error {
	if chatID == "" || message == "" {
		return nil
	}

	query := `
		INSERT INTO telegram_messages (chat_id, message, sent_at)
		VALUES ($1, $2, NOW())
	`
	_, err := r.DB.Exec(ctx, query, chatID, message)
	if err != nil {
		log.Printf("[SendMessage] error: %v", err)
		return err
	}
	log.Printf("[SendMessage] sent to chat_id=%s", chatID)
	return nil
}
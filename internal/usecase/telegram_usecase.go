package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
)

type TelegramUsecase struct {
	KaryawanRepo *repository.KaryawanRepo
	TelegramRepo *repository.TelegramRepo
}

func NewTelegramUsecase(
	karyawanRepo *repository.KaryawanRepo,
	telegramRepo *repository.TelegramRepo,
	db *pgxpool.Pool,
) *TelegramUsecase {
	return &TelegramUsecase{
		KaryawanRepo: karyawanRepo,
		TelegramRepo: telegramRepo,
	}
}

func (u *TelegramUsecase) GenerateVerificationCode(ctx context.Context, chatID, username string) (string, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	code := "SAKTI-" + hex.EncodeToString(bytes)

	if err := u.TelegramRepo.SaveVerification(ctx, code, chatID, username); err != nil {
		log.Printf("Failed to save verification code: %v", err)
		return "", err
	}

	log.Printf("Verification code saved to database: %s", code)
	return code, nil
}

func (u *TelegramUsecase) VerifyCode(ctx context.Context, code string) (string, error) {
	code = strings.ToLower(code)

	if code == "" {
		return "", errors.New("kode verifikasi tidak valid")
	}

	verif, err := u.TelegramRepo.GetVerificationByCode(ctx, code)
	if err != nil {
		return "", err
	}
	if verif == nil {
		return "", errors.New("kode verifikasi tidak valid")
	}

	if verif.IsUsed {
		return "", errors.New("kode verifikasi sudah digunakan")
	}

	if time.Now().After(verif.ExpiredAt) {
		return "", errors.New("kode verifikasi sudah kadaluarsa")
	}

	return verif.ChatID, nil
}

func (u *TelegramUsecase) ConnectTelegram(ctx context.Context, karyawanID, code string) error {
	code = strings.ToLower(code)

	log.Printf("[Telegram] Connect started: karyawanID=%s, code=%s", karyawanID, code)

	if code == "" {
		log.Println("[Telegram] Code is empty")
		return errors.New("kode verifikasi wajib diisi")
	}

	chatID, err := u.VerifyCode(ctx, code)
	if err != nil {
		log.Printf("[Telegram] VerifyCode failed: %v", err)
		return err
	}

	log.Printf("[Telegram] ChatID found: %s", chatID)

	if err := u.KaryawanRepo.UpdateTelegramChatID(ctx, karyawanID, chatID); err != nil {
		log.Printf("[Telegram] UpdateTelegramChatID failed: %v", err)
		return err
	}

	if err := u.TelegramRepo.MarkVerificationAsUsed(ctx, code, karyawanID); err != nil {
		log.Printf("[Telegram] MarkVerificationAsUsed failed: %v", err)
	}

	log.Printf("[Telegram] Connect success: karyawanID=%s, chatID=%s", karyawanID, chatID)
	return nil
}

func (u *TelegramUsecase) DisconnectTelegram(ctx context.Context, karyawanID string) error {
	log.Printf("[Telegram] Disconnect started: karyawanID=%s", karyawanID)

	if err := u.KaryawanRepo.ClearTelegramChatID(ctx, karyawanID); err != nil {
		log.Printf("[Telegram] ClearTelegramChatID failed: %v", err)
		return err
	}

	log.Printf("[Telegram] Disconnect success: karyawanID=%s", karyawanID)
	return nil
}

func (u *TelegramUsecase) GetTelegramStatus(ctx context.Context, karyawanID string) (*domain.TelegramStatusResponse, error) {
	log.Printf("[Telegram] GetStatus started: karyawanID=%s", karyawanID)

	chatID, err := u.KaryawanRepo.GetTelegramStatus(ctx, karyawanID)
	if err != nil {
		log.Printf("[Telegram] GetTelegramStatus error: %v", err)
		return nil, err
	}

	log.Printf("[Telegram] GetStatus result: chatID=%s, connected=%v", chatID, chatID != "")
	return &domain.TelegramStatusResponse{
		IsConnected:    chatID != "",
		TelegramChatID: chatID,
	}, nil
}

func (u *TelegramUsecase) CleanupExpiredCodes(ctx context.Context) error {
	log.Println("[Telegram] CleanupExpiredCodes started")
	if err := u.TelegramRepo.DeleteExpiredVerifications(ctx); err != nil {
		log.Printf("[Telegram] CleanupExpiredCodes error: %v", err)
		return err
	}
	log.Println("[Telegram] CleanupExpiredCodes success")
	return nil
}
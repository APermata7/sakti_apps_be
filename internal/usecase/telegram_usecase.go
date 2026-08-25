package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"

	"sakti_apps_be/internal/repository"
)

type TelegramUsecase struct {
	KaryawanRepo *repository.KaryawanRepo
	TelegramRepo *repository.TelegramRepo
}

func NewTelegramUsecase(
	karyawanRepo *repository.KaryawanRepo,
	telegramRepo *repository.TelegramRepo,
	db interface{},
) *TelegramUsecase {
	return &TelegramUsecase{
		KaryawanRepo: karyawanRepo,
		TelegramRepo: telegramRepo,
	}
}

func (u *TelegramUsecase) GetTelegramStatus(ctx context.Context, karyawanID string) (map[string]interface{}, error) {
	karyawan, err := u.KaryawanRepo.GetByID(ctx, karyawanID)
	if err != nil || karyawan == nil {
		return nil, errors.New("karyawan tidak ditemukan")
	}

	connected := false
	chatID := ""
	if karyawan.TelegramChatID != nil && *karyawan.TelegramChatID != "" {
		connected = true
		chatID = *karyawan.TelegramChatID
	}

	return map[string]interface{}{
		"connected": connected,
		"chat_id":   chatID,
	}, nil
}

func (u *TelegramUsecase) ConnectTelegram(ctx context.Context, karyawanID string, verificationCode string) error {
	chatID, err := u.TelegramRepo.GetChatIDByCode(ctx, verificationCode)
	if err != nil {
		return errors.New("kode verifikasi tidak valid, sudah digunakan, atau sudah kadaluarsa")
	}

	latestCode, err := u.TelegramRepo.GetVerificationCode(ctx, chatID)
	if err != nil {
		return errors.New("kode verifikasi tidak valid, sudah digunakan, atau sudah kadaluarsa")
	}

	if latestCode != verificationCode {
		return errors.New("kode verifikasi sudah tidak berlaku, silakan minta kode baru dengan /start")
	}

	karyawan, err := u.KaryawanRepo.GetByID(ctx, karyawanID)
	if err != nil || karyawan == nil {
		return errors.New("karyawan tidak ditemukan")
	}

	tx, err := u.KaryawanRepo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	karyawan.TelegramChatID = &chatID
	if err := u.KaryawanRepo.UpdateWithTx(ctx, tx, karyawan); err != nil {
		log.Printf("UpdateChatID error: %v", err)
		return errors.New("gagal menyimpan koneksi Telegram")
	}

	if err := u.TelegramRepo.MarkCodeAsUsedWithTx(ctx, tx, verificationCode); err != nil {
		log.Printf("MarkCodeAsUsed error: %v", err)
		return errors.New("gagal menyelesaikan verifikasi Telegram")
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (u *TelegramUsecase) UpdateChatID(ctx context.Context, karyawanID, chatID string) error {
	karyawan, err := u.KaryawanRepo.GetByID(ctx, karyawanID)
	if err != nil || karyawan == nil {
		return errors.New("karyawan tidak ditemukan")
	}

	karyawan.TelegramChatID = &chatID
	if err := u.KaryawanRepo.Update(ctx, karyawan); err != nil {
		log.Printf("UpdateChatID error: %v", err)
		return errors.New("gagal menyimpan chat ID")
	}

	log.Printf("Telegram chat ID updated for karyawan: %s", karyawanID)
	return nil
}

func (u *TelegramUsecase) ClearChatID(ctx context.Context, karyawanID string) error {
	karyawan, err := u.KaryawanRepo.GetByID(ctx, karyawanID)
	if err != nil || karyawan == nil {
		return errors.New("karyawan tidak ditemukan")
	}

	karyawan.TelegramChatID = nil
	if err := u.KaryawanRepo.Update(ctx, karyawan); err != nil {
		log.Printf("ClearChatID error: %v", err)
		return errors.New("gagal menghapus chat ID")
	}

	log.Printf("Telegram chat ID cleared for karyawan: %s", karyawanID)
	return nil
}

func (u *TelegramUsecase) GenerateVerificationCode(ctx context.Context, chatID, username string) (string, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	code := "SAKTI-" + hex.EncodeToString(bytes)

	if err := u.TelegramRepo.SaveVerificationCode(ctx, chatID, username, code); err != nil {
		log.Printf("SaveVerificationCode error: %v", err)
		return "", err
	}

	log.Printf("Verification code generated for chat_id: %s, code: %s", chatID, code)
	return code, nil
}

func (u *TelegramUsecase) CleanupExpiredCodes(ctx context.Context) error {
	log.Println("CleanupExpiredCodes started")
	if err := u.TelegramRepo.DeleteExpiredCodes(ctx); err != nil {
		log.Printf("CleanupExpiredCodes error: %v", err)
		return err
	}
	log.Println("CleanupExpiredCodes success")
	return nil
}

func (u *TelegramUsecase) VerifyCode(ctx context.Context, chatID, code string) (bool, error) {
	expected, err := u.TelegramRepo.GetVerificationCode(ctx, chatID)
	if err != nil {
		return false, err
	}
	if expected != code {
		return false, nil
	}
	return true, nil
}
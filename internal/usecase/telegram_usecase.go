package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
)

type TelegramUsecase struct {
	KaryawanRepo      *repository.KaryawanRepo
	verificationCache map[string]*domain.TelegramVerificationCode
	mu                sync.RWMutex
}

func NewTelegramUsecase(karyawanRepo *repository.KaryawanRepo, db *pgxpool.Pool) *TelegramUsecase {
	return &TelegramUsecase{
		KaryawanRepo:      karyawanRepo,
		verificationCache: make(map[string]*domain.TelegramVerificationCode),
	}
}

func (u *TelegramUsecase) GenerateVerificationCode(chatID, username string) (string, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	code := "SAKTI-" + hex.EncodeToString(bytes)

	u.mu.Lock()
	u.verificationCache[code] = &domain.TelegramVerificationCode{
		Code:      code,
		ChatID:    chatID,
		Username:  username,
		CreatedAt: time.Now().Unix(),
	}
	u.mu.Unlock()

	return code, nil
}

func (u *TelegramUsecase) VerifyCode(code string) (string, error) {
	u.mu.RLock()
	verif, exists := u.verificationCache[code]
	u.mu.RUnlock()

	if !exists {
		return "", errors.New("kode verifikasi tidak valid")
	}

	if time.Now().Unix()-verif.CreatedAt > 300 {
		u.mu.Lock()
		delete(u.verificationCache, code)
		u.mu.Unlock()
		return "", errors.New("kode verifikasi sudah kadaluarsa")
	}

	return verif.ChatID, nil
}

func (u *TelegramUsecase) ConnectTelegram(ctx context.Context, karyawanID, code string) error {
	if code == "" {
		return errors.New("kode verifikasi wajib diisi")
	}

	chatID, err := u.VerifyCode(code)
	if err != nil {
		return err
	}

	if err := u.KaryawanRepo.UpdateTelegramChatID(ctx, karyawanID, chatID); err != nil {
		return err
	}

	u.mu.Lock()
	delete(u.verificationCache, code)
	u.mu.Unlock()

	return nil
}

func (u *TelegramUsecase) DisconnectTelegram(ctx context.Context, karyawanID string) error {
	return u.KaryawanRepo.ClearTelegramChatID(ctx, karyawanID)
}

func (u *TelegramUsecase) GetTelegramStatus(ctx context.Context, karyawanID string) (*domain.TelegramStatusResponse, error) {
	chatID, err := u.KaryawanRepo.GetTelegramStatus(ctx, karyawanID)
	if err != nil {
		return nil, err
	}

	return &domain.TelegramStatusResponse{
		IsConnected:    chatID != "",
		TelegramChatID: chatID,
	}, nil
}
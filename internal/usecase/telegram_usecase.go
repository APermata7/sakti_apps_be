package usecase

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "log"
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
    if code == "" {
        return errors.New("kode verifikasi wajib diisi")
    }

    chatID, err := u.VerifyCode(ctx, code)
    if err != nil {
        return err
    }

    if err := u.KaryawanRepo.UpdateTelegramChatID(ctx, karyawanID, chatID); err != nil {
        return err
    }

    if err := u.TelegramRepo.MarkVerificationAsUsed(ctx, code, karyawanID); err != nil {
        log.Printf("Failed to mark verification as used: %v", err)
    }

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

func (u *TelegramUsecase) CleanupExpiredCodes(ctx context.Context) error {
    return u.TelegramRepo.DeleteExpiredVerifications(ctx)
}
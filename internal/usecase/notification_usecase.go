package usecase

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
	"sakti_apps_be/internal/utils"
)

type NotificationUsecase struct {
	FCMTokenRepo   *repository.FCMTokenRepo
	NotifikasiRepo *repository.NotifikasiRepo
	KaryawanRepo   *repository.KaryawanRepo
	TelegramBot    *utils.TelegramBot
}

func NewNotificationUsecase(
	fcmTokenRepo *repository.FCMTokenRepo,
	notifikasiRepo *repository.NotifikasiRepo,
	karyawanRepo *repository.KaryawanRepo,
	db *pgxpool.Pool,
) *NotificationUsecase {
	return &NotificationUsecase{
		FCMTokenRepo:   fcmTokenRepo,
		NotifikasiRepo: notifikasiRepo,
		KaryawanRepo:   karyawanRepo,
		TelegramBot:    utils.NewTelegramBot(db),
	}
}

func (u *NotificationUsecase) RegisterFCMToken(ctx context.Context, karyawanID, fcmToken string) error {
	token := &domain.FCMToken{
		KaryawanID: karyawanID,
		FCMToken:   fcmToken,
		IsActive:   true,
	}
	return u.FCMTokenRepo.SaveToken(ctx, token)
}

func (u *NotificationUsecase) DeactivateFCMToken(ctx context.Context, karyawanID, fcmToken string) error {
	tokens, err := u.FCMTokenRepo.GetTokensByKaryawanID(ctx, karyawanID)
	if err != nil {
		return err
	}

	for _, token := range tokens {
		if token == fcmToken {
			return u.FCMTokenRepo.DeactivateToken(ctx, fcmToken)
		}
	}

	return nil
}

func (u *NotificationUsecase) KirimNotifikasi(ctx context.Context, req domain.KirimNotifikasiRequest, chatID string) error {
	notif := &domain.Notifikasi{
		KaryawanID:    req.KaryawanID,
		Jenis:         req.Jenis,
		Channel:       "inapp",
		Judul:         req.Judul,
		Pesan:         req.Pesan,
		Dibaca:        false,
		ReferensiID:   req.ReferensiID,
		ReferensiTipe: req.ReferensiTipe,
	}

	if err := u.NotifikasiRepo.Create(ctx, notif); err != nil {
		return err
	}

	tokens, _ := u.FCMTokenRepo.GetTokensByKaryawanID(ctx, req.KaryawanID)
	if len(tokens) > 0 {
		go utils.SendMulticast(tokens, req.Judul, req.Pesan)
	}

	if chatID != "" && u.TelegramBot != nil {
		go func() {
			err := u.TelegramBot.SendNotification(chatID, req.Judul, req.Pesan)
			if err != nil {
				log.Printf("Failed to send Telegram notification: %v", err)
			}
		}()
	}

	return nil
}

func (u *NotificationUsecase) KirimInApp(ctx context.Context, req domain.KirimNotifikasiRequest) error {
	notif := &domain.Notifikasi{
		KaryawanID:    req.KaryawanID,
		Jenis:         req.Jenis,
		Channel:       "inapp",
		Judul:         req.Judul,
		Pesan:         req.Pesan,
		Dibaca:        false,
		ReferensiID:   req.ReferensiID,
		ReferensiTipe: req.ReferensiTipe,
	}

	if err := u.NotifikasiRepo.Create(ctx, notif); err != nil {
		return err
	}

	tokens, _ := u.FCMTokenRepo.GetTokensByKaryawanID(ctx, req.KaryawanID)
	if len(tokens) > 0 {
		go utils.SendMulticast(tokens, req.Judul, req.Pesan)
	}

	return nil
}

func (u *NotificationUsecase) SendTelegramNotification(chatID, title, message string) error {
	if u.TelegramBot == nil {
		log.Println("Telegram bot not initialized")
		return nil
	}
	return u.TelegramBot.SendNotification(chatID, title, message)
}

func (u *NotificationUsecase) SendLeaveNotification(chatID, karyawanNama, status, tanggal string) error {
	if u.TelegramBot == nil {
		return nil
	}
	return u.TelegramBot.SendLeaveNotification(chatID, karyawanNama, status, tanggal)
}

func (u *NotificationUsecase) SendApprovalNotification(chatID, karyawanNama, totalHari, alasan string) error {
	if u.TelegramBot == nil {
		return nil
	}
	return u.TelegramBot.SendApprovalNotification(chatID, karyawanNama, totalHari, alasan)
}

func (u *NotificationUsecase) GetNotifikasi(ctx context.Context, karyawanID string, page, limit int) ([]domain.Notifikasi, int, error) {
	offset := (page - 1) * limit
	return u.NotifikasiRepo.GetByKaryawanID(ctx, karyawanID, limit, offset)
}

func (u *NotificationUsecase) GetUnreadCount(ctx context.Context, karyawanID string) (int, error) {
	return u.NotifikasiRepo.GetUnreadCount(ctx, karyawanID)
}

func (u *NotificationUsecase) MarkAsRead(ctx context.Context, id, karyawanID string) error {
	return u.NotifikasiRepo.MarkAsRead(ctx, id, karyawanID)
}

func (u *NotificationUsecase) MarkAllAsRead(ctx context.Context, karyawanID string) error {
	return u.NotifikasiRepo.MarkAllAsRead(ctx, karyawanID)
}
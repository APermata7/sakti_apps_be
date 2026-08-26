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
	TelegramRepo   *repository.TelegramRepo
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
		TelegramRepo:   repository.NewTelegramRepo(db),
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

func (u *NotificationUsecase) sendTelegramNotification(ctx context.Context, karyawanID, message string) {
	if u.TelegramRepo == nil || karyawanID == "" || message == "" {
		return
	}

	go func() {
		err := u.TelegramRepo.SendMessageByKaryawanID(ctx, karyawanID, message)
		if err != nil {
			log.Printf("[sendTelegramNotification] error: %v", err)
		}
	}()
}

func (u *NotificationUsecase) KirimNotifikasi(ctx context.Context, req domain.KirimNotifikasiRequest) error {
	karyawan, err := u.KaryawanRepo.GetByID(ctx, req.KaryawanID)
	if err != nil || karyawan == nil {
		log.Printf("[KirimNotifikasi] karyawan tidak ditemukan: %s", req.KaryawanID)
		return nil
	}

	if karyawan.StatusKaryawan != "aktif" {
		log.Printf("[KirimNotifikasi] karyawan %s nonaktif, skip notifikasi", req.KaryawanID)
		return nil
	}

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
		log.Printf("[KirimNotifikasi] Error create notification: %v", err)
		return err
	}

	tokens, err := u.FCMTokenRepo.GetTokensByKaryawanID(ctx, req.KaryawanID)
	if err != nil {
		log.Printf("[KirimNotifikasi] Error get tokens: %v", err)
	}

	if len(tokens) > 0 {
		go func() {
			err := utils.SendMulticast(tokens, req.Judul, req.Pesan)
			if err != nil {
				log.Printf("[KirimNotifikasi] SendMulticast error: %v", err)
			}
		}()
	}

	if u.TelegramRepo != nil {
		u.sendTelegramNotification(ctx, req.KaryawanID, req.Judul+"\n\n"+req.Pesan)
	}

	return nil
}

func (u *NotificationUsecase) KirimInApp(ctx context.Context, req domain.KirimNotifikasiRequest) error {
	karyawan, err := u.KaryawanRepo.GetByID(ctx, req.KaryawanID)
	if err != nil || karyawan == nil {
		log.Printf("[KirimInApp] karyawan tidak ditemukan: %s", req.KaryawanID)
		return nil
	}

	if karyawan.StatusKaryawan != "aktif" {
		log.Printf("[KirimInApp] karyawan %s nonaktif, skip notifikasi", req.KaryawanID)
		return nil
	}

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
		log.Printf("[KirimInApp] Error create notification: %v", err)
		return err
	}

	log.Printf("[KirimInApp] Notification saved: id=%s", notif.ID)

	tokens, err := u.FCMTokenRepo.GetTokensByKaryawanID(ctx, req.KaryawanID)
	if err != nil {
		log.Printf("[KirimInApp] Error get tokens: %v", err)
	}

	if len(tokens) > 0 {
		go func() {
			err := utils.SendMulticast(tokens, req.Judul, req.Pesan)
			if err != nil {
				log.Printf("[KirimInApp] SendMulticast error: %v", err)
			}
		}()
	}

	return nil
}

func (u *NotificationUsecase) KirimInAppWithTelegram(ctx context.Context, req domain.KirimNotifikasiRequest) error {
	err := u.KirimInApp(ctx, req)
	if err != nil {
		return err
	}

	if u.TelegramRepo != nil {
		u.sendTelegramNotification(ctx, req.KaryawanID, req.Judul+"\n\n"+req.Pesan)
	}

	return nil
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
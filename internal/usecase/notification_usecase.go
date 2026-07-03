package usecase

import (
	"context"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
	"sakti_apps_be/internal/utils"
)

type NotificationUsecase struct {
	FCMTokenRepo   *repository.FCMTokenRepo
	NotifikasiRepo *repository.NotifikasiRepo
	KaryawanRepo   *repository.KaryawanRepo
}

func NewNotificationUsecase(
	fcmTokenRepo *repository.FCMTokenRepo,
	notifikasiRepo *repository.NotifikasiRepo,
	karyawanRepo *repository.KaryawanRepo,
) *NotificationUsecase {
	return &NotificationUsecase{
		FCMTokenRepo:   fcmTokenRepo,
		NotifikasiRepo: notifikasiRepo,
		KaryawanRepo:   karyawanRepo,
	}
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

func (u *NotificationUsecase) KirimWhatsApp(ctx context.Context, karyawanID, jenis, tanggal string) error {
	karyawan, err := u.KaryawanRepo.GetByID(ctx, karyawanID)
	if err != nil || karyawan == nil {
		return err
	}

	if karyawan.NomorTelepon == nil || *karyawan.NomorTelepon == "" {
		return nil
	}

	nomor := *karyawan.NomorTelepon
	if len(nomor) > 0 && nomor[0:1] == "0" {
		nomor = "62" + nomor[1:]
	}

	go utils.SendWhatsAppNotification(nomor, karyawan.NamaLengkap, "pegawai", jenis, tanggal)
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
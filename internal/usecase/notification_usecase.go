package usecase

import (
	"context"
	"log"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
	"sakti_apps_be/internal/utils"
)

type NotificationUsecase struct {
	FCMTokenRepo *repository.FCMTokenRepo
}

func NewNotificationUsecase(fcmTokenRepo *repository.FCMTokenRepo) *NotificationUsecase {
	return &NotificationUsecase{FCMTokenRepo: fcmTokenRepo}
}

func (u *NotificationUsecase) SendToKaryawan(ctx context.Context, karyawanID, title, body string) error {
	tokens, err := u.FCMTokenRepo.GetTokensByKaryawanID(ctx, karyawanID)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		log.Printf("No FCM token for karyawan %s", karyawanID)
		return nil
	}
	return utils.SendMulticast(tokens, title, body)
}

func (u *NotificationUsecase) SendToKaryawanWithData(ctx context.Context, karyawanID, title, body string, data map[string]string) error {
	tokens, err := u.FCMTokenRepo.GetTokensByKaryawanID(ctx, karyawanID)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		log.Printf("No FCM token for karyawan %s", karyawanID)
		return nil
	}
	return utils.SendMulticastWithData(tokens, title, body, data)
}

func (u *NotificationUsecase) SendToAtasan(ctx context.Context, karyawanID, title, body string) error {
	tokens, err := u.FCMTokenRepo.GetTokensByAtasan(ctx, karyawanID)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		log.Printf("No FCM token for atasan of karyawan %s", karyawanID)
		return nil
	}
	return utils.SendMulticast(tokens, title, body)
}

func (u *NotificationUsecase) SendToAtasanWithData(ctx context.Context, karyawanID, title, body string, data map[string]string) error {
	tokens, err := u.FCMTokenRepo.GetTokensByAtasan(ctx, karyawanID)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		log.Printf("No FCM token for atasan of karyawan %s", karyawanID)
		return nil
	}
	return utils.SendMulticastWithData(tokens, title, body, data)
}

func (u *NotificationUsecase) SendToHRD(ctx context.Context, title, body string) error {
	tokens, err := u.FCMTokenRepo.GetTokensByHRD(ctx)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		log.Println("No FCM token for HRD")
		return nil
	}
	return utils.SendMulticast(tokens, title, body)
}

func (u *NotificationUsecase) SendToHRDWithData(ctx context.Context, title, body string, data map[string]string) error {
	tokens, err := u.FCMTokenRepo.GetTokensByHRD(ctx)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		log.Println("No FCM token for HRD")
		return nil
	}
	return utils.SendMulticastWithData(tokens, title, body, data)
}

func (u *NotificationUsecase) SendToRole(ctx context.Context, role, title, body string) error {
	tokens, err := u.FCMTokenRepo.GetTokensByRole(ctx, role)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		log.Printf("No FCM token for role %s", role)
		return nil
	}
	return utils.SendMulticast(tokens, title, body)
}

func (u *NotificationUsecase) SendToRoleWithData(ctx context.Context, role, title, body string, data map[string]string) error {
	tokens, err := u.FCMTokenRepo.GetTokensByRole(ctx, role)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		log.Printf("No FCM token for role %s", role)
		return nil
	}
	return utils.SendMulticastWithData(tokens, title, body, data)
}

func (u *NotificationUsecase) SaveToken(ctx context.Context, token *domain.FCMToken) error {
	return u.FCMTokenRepo.SaveToken(ctx, token)
}

func (u *NotificationUsecase) DeleteToken(ctx context.Context, token string) error {
	return u.FCMTokenRepo.DeactivateToken(ctx, token)
}

func (u *NotificationUsecase) DeleteAllTokens(ctx context.Context, karyawanID string) error {
	return u.FCMTokenRepo.DeactivateTokensByKaryawan(ctx, karyawanID)
}
package usecase

import (
	"context"
	"errors"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
)

type KonfigurasiUsecase struct {
	KonfigurasiRepo *repository.KonfigurasiRepo
}

func NewKonfigurasiUsecase(konfigurasiRepo *repository.KonfigurasiRepo) *KonfigurasiUsecase {
	return &KonfigurasiUsecase{KonfigurasiRepo: konfigurasiRepo}
}

func (u *KonfigurasiUsecase) GetConfig(ctx context.Context) (*domain.KonfigurasiKerja, error) {
	config, err := u.KonfigurasiRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("konfigurasi tidak ditemukan")
	}
	return config, nil
}

func (u *KonfigurasiUsecase) UpdateConfig(ctx context.Context, userID string, req domain.UpdateKonfigurasiRequest) (*domain.KonfigurasiKerja, error) {
	config, err := u.KonfigurasiRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("konfigurasi tidak ditemukan")
	}

	if req.NamaKantor != "" {
		config.NamaKantor = req.NamaKantor
	}
	if req.LatKantor != 0 {
		config.LatKantor = req.LatKantor
	}
	if req.LongKantor != 0 {
		config.LongKantor = req.LongKantor
	}
	if req.LogoKantor != "" {
		config.LogoKantor = req.LogoKantor
	}
	if req.JamMasuk != "" {
		config.JamMasuk = req.JamMasuk
	}
	if req.JamMinimalMasuk != "" {
		config.JamMinimalMasuk = req.JamMinimalMasuk
	}
	if req.JamPulang != "" {
		config.JamPulang = req.JamPulang
	}
	if req.JamMinimalPulang != "" {
		config.JamMinimalPulang = req.JamMinimalPulang
	}
	if req.RadiusKantor != 0 {
		config.RadiusKantor = req.RadiusKantor
	}
	config.DiperbaruiOleh = userID

	if err := u.KonfigurasiRepo.Update(ctx, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (u *KonfigurasiUsecase) UpdateLogo(ctx context.Context, userID, logoURL string) (*domain.KonfigurasiKerja, error) {
	config, err := u.KonfigurasiRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("konfigurasi tidak ditemukan")
	}

	config.LogoKantor = logoURL
	config.DiperbaruiOleh = userID

	if err := u.KonfigurasiRepo.Update(ctx, config); err != nil {
		return nil, err
	}
	return config, nil
}
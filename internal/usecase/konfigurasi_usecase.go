package usecase

import (
	"context"
	"errors"
	"mime/multipart"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
	"sakti_apps_be/internal/utils"
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

func (u *KonfigurasiUsecase) GetWorkConfig(ctx context.Context) (*domain.WorkConfigResponse, error) {
	config, err := u.KonfigurasiRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("konfigurasi tidak ditemukan")
	}

	return &domain.WorkConfigResponse{
		JamMasuk:         config.JamMasuk,
		JamMinimalMasuk:  config.JamMinimalMasuk,
		JamPulang:        config.JamPulang,
		JamMinimalPulang: config.JamMinimalPulang,
		RadiusKantor:     config.RadiusKantor,
	}, nil
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
	if req.LogoKantor != nil {
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

	if userID != "" {
		config.DiperbaruiOleh = &userID
	} else {
		config.DiperbaruiOleh = nil
	}

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

	logo := logoURL
	config.LogoKantor = &logo

	if userID != "" {
		config.DiperbaruiOleh = &userID
	} else {
		config.DiperbaruiOleh = nil
	}

	if err := u.KonfigurasiRepo.Update(ctx, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (u *KonfigurasiUsecase) UploadLogo(ctx context.Context, userID string, file multipart.File, filename string) (*domain.KonfigurasiKerja, error) {
	logoURL, err := utils.UploadImage(file, filename)
	if err != nil {
		return nil, errors.New("gagal upload logo ke Cloudinary: " + err.Error())
	}

	if logoURL == "" {
		return nil, errors.New("Cloudinary tidak tersedia")
	}

	return u.UpdateLogo(ctx, userID, logoURL)
}
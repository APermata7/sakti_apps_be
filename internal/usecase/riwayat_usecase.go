package usecase

import (
	"context"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
)

type RiwayatUsecase struct {
	RiwayatRepo *repository.RiwayatRepo
}

func NewRiwayatUsecase(riwayatRepo *repository.RiwayatRepo) *RiwayatUsecase {
	return &RiwayatUsecase{
		RiwayatRepo: riwayatRepo,
	}
}

func (u *RiwayatUsecase) GetRiwayat(ctx context.Context, karyawanID string, page, limit int) (*domain.RiwayatResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	offset := (page - 1) * limit

	items, total, err := u.RiwayatRepo.GetRiwayat(ctx, karyawanID, limit, offset)
	if err != nil {
		return nil, err
	}

	totalPages := (total + limit - 1) / limit

	return &domain.RiwayatResponse{
		Items: items,
		Meta: domain.MetaPagination{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}, nil
}
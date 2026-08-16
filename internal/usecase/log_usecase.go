package usecase

import (
    "context"
    "sakti_apps_be/internal/repository"
)

type LogUsecase struct {
    LogRepo *repository.LogRepo
}

func NewLogUsecase(logRepo *repository.LogRepo) *LogUsecase {
    return &LogUsecase{LogRepo: logRepo}
}

func (u *LogUsecase) CreateLog(ctx context.Context, karyawanID, action, detail string) error {
    return u.LogRepo.CreateLog(ctx, karyawanID, action, detail)
}

func (u *LogUsecase) GetLogs(ctx context.Context, limit, offset int) ([]map[string]interface{}, int, error) {
    return u.LogRepo.GetLogs(ctx, limit, offset)
}
package usecase

import (
	"context"
	"errors"
	"time"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
	"sakti_apps_be/internal/utils"
)

type PresensiUsecase struct {
	PresensiRepo *repository.PresensiRepo
	KaryawanRepo *repository.KaryawanRepo
	ConfigRepo   *repository.KonfigurasiRepo
}

func NewPresensiUsecase(presensiRepo *repository.PresensiRepo, karyawanRepo *repository.KaryawanRepo, configRepo *repository.KonfigurasiRepo) *PresensiUsecase {
	return &PresensiUsecase{
		PresensiRepo: presensiRepo,
		KaryawanRepo: karyawanRepo,
		ConfigRepo:   configRepo,
	}
}

func (u *PresensiUsecase) CheckIn(ctx context.Context, karyawanID string, req domain.CheckInRequest) (*domain.CheckInResponse, error) {
	existing, _ := u.PresensiRepo.GetToday(ctx, karyawanID)
	if existing != nil && existing.JamMasuk != "" {
		return nil, errors.New("anda sudah melakukan presensi masuk hari ini")
	}

	karyawan, err := u.KaryawanRepo.GetByID(ctx, karyawanID)
	if err != nil || karyawan == nil {
		return nil, errors.New("karyawan tidak ditemukan")
	}

	if karyawan.StatusKaryawan != "aktif" {
		return nil, errors.New("akun tidak aktif")
	}

	config, err := u.ConfigRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	faceMatch, similarity, err := utils.VerifyFaceWithRepo(ctx, req.SelfieURL, karyawanID, u.KaryawanRepo)
	if err != nil {
		return nil, err
	}
	if !faceMatch {
		return nil, errors.New("wajah tidak dikenali")
	}

	distance := utils.Haversine(config.LatKantor, config.LongKantor, req.Latitude, req.Longitude)
	isOutside := distance > float64(config.RadiusKantor)

	locationStatus := "di_dalam_radius"
	if isOutside {
		locationStatus = "di_luar_radius"
	}

	jamMasuk := time.Now()
	jamMasukStr := jamMasuk.Format("15:04:05")

	batasToleransi, _ := time.Parse("15:04:05", config.JamMinimalMasuk)
	status := "tepat_waktu"
	if jamMasuk.After(batasToleransi) {
		status = "terlambat"
	}

	presensi := &domain.Presensi{
		KaryawanID:      karyawanID,
		Tanggal:         jamMasuk.Format("2006-01-02"),
		JamMasuk:        jamMasukStr,
		Status:          status,
		LintangMasuk:    req.Latitude,
		BujurMasuk:      req.Longitude,
		ValidasiWajah:   faceMatch,
		URLFoto:         req.SelfieURL,
		AlasanTerlambat: req.AlasanTerlambat,
		DistanceMeter:   distance,
		IsOutsideRadius: isOutside,
		LocationStatus:  locationStatus,
	}

	if err := u.PresensiRepo.Create(ctx, presensi); err != nil {
		return nil, err
	}

	return &domain.CheckInResponse{
		ID:              presensi.ID,
		KaryawanID:      presensi.KaryawanID,
		Tanggal:         presensi.Tanggal,
		JamMasuk:        presensi.JamMasuk,
		Status:          presensi.Status,
		ValidasiWajah:   presensi.ValidasiWajah,
		FaceSimilarity:  similarity,
		URLFoto:         presensi.URLFoto,
		DistanceMeter:   presensi.DistanceMeter,
		IsOutsideRadius: presensi.IsOutsideRadius,
		LocationStatus:  presensi.LocationStatus,
		OfficeLatitude:  config.LatKantor,
		OfficeLongitude: config.LongKantor,
		OfficeRadius:    config.RadiusKantor,
	}, nil
}

func (u *PresensiUsecase) CheckOut(ctx context.Context, karyawanID string, req domain.CheckOutRequest) (*domain.CheckOutResponse, error) {
	presensi, err := u.PresensiRepo.GetToday(ctx, karyawanID)
	if err != nil {
		return nil, errors.New("anda belum melakukan presensi masuk hari ini")
	}
	if presensi == nil || presensi.JamMasuk == "" {
		return nil, errors.New("anda belum melakukan presensi masuk hari ini")
	}
	if presensi.JamKeluar != "" {
		return nil, errors.New("anda sudah melakukan presensi keluar hari ini")
	}

	config, err := u.ConfigRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	jamKeluar := time.Now()
	jamKeluarStr := jamKeluar.Format("15:04:05")

	jamPulang, _ := time.Parse("15:04:05", config.JamPulang)
	lembur := jamKeluar.After(jamPulang)

	var jamLembur float64
	if lembur {
		selisih := jamKeluar.Sub(jamPulang)
		jamLembur = selisih.Hours()
	}

	distance := utils.Haversine(config.LatKantor, config.LongKantor, req.Latitude, req.Longitude)
	isOutside := distance > float64(config.RadiusKantor)

	locationStatus := "di_dalam_radius"
	if isOutside {
		locationStatus = "di_luar_radius"
	}

	if err := u.PresensiRepo.UpdateCheckOut(ctx, presensi.ID, jamKeluarStr, lembur, jamLembur, req.Latitude, req.Longitude); err != nil {
		return nil, err
	}

	return &domain.CheckOutResponse{
		ID:              presensi.ID,
		KaryawanID:      presensi.KaryawanID,
		Tanggal:         presensi.Tanggal,
		JamMasuk:        presensi.JamMasuk,
		JamKeluar:       jamKeluarStr,
		Lembur:          lembur,
		JamLembur:       jamLembur,
		DistanceMeter:   distance,
		IsOutsideRadius: isOutside,
		LocationStatus:  locationStatus,
	}, nil
}

func (u *PresensiUsecase) GetToday(ctx context.Context, karyawanID string) (*domain.TodayResponse, error) {
	presensi, err := u.PresensiRepo.GetToday(ctx, karyawanID)
	if err != nil {
		return &domain.TodayResponse{
			HasCheckedIn:  false,
			HasCheckedOut: false,
			Tanggal:       time.Now().Format("2006-01-02"),
			KaryawanID:    karyawanID,
		}, nil
	}

	if presensi == nil {
		return &domain.TodayResponse{
			HasCheckedIn:  false,
			HasCheckedOut: false,
			Tanggal:       time.Now().Format("2006-01-02"),
			KaryawanID:    karyawanID,
		}, nil
	}

	return &domain.TodayResponse{
		HasCheckedIn:  presensi.JamMasuk != "",
		HasCheckedOut: presensi.JamKeluar != "",
		CheckInTime:   presensi.JamMasuk,
		CheckInStatus: presensi.Status,
		Tanggal:       presensi.Tanggal,
		KaryawanID:    presensi.KaryawanID,
	}, nil
}

func (u *PresensiUsecase) GetHistory(ctx context.Context, karyawanID, startDate, endDate, status string, limit, page int) ([]domain.Presensi, int, error) {
	offset := (page - 1) * limit
	return u.PresensiRepo.GetHistory(ctx, karyawanID, startDate, endDate, status, limit, offset)
}

func (u *PresensiUsecase) UpdateAlasanTerlambat(ctx context.Context, karyawanID, alasan string) error {
	presensi, err := u.PresensiRepo.GetToday(ctx, karyawanID)
	if err != nil || presensi == nil {
		return errors.New("data presensi hari ini tidak ditemukan")
	}
	if presensi.Status != "terlambat" {
		return errors.New("anda tidak terlambat")
	}
	return u.PresensiRepo.UpdateAlasanTerlambat(ctx, presensi.ID, alasan)
}

func (u *PresensiUsecase) AutoClockOut(ctx context.Context) error {
	jamKeluar := "23:59:59"
	lat := 0.0
	lon := 0.0

	query := `
		UPDATE presensi 
		SET jam_keluar = $1, lembur = false, jam_lembur = 0,
		    lintang_keluar = $2, bujur_keluar = $3,
		    diperbarui_pada = NOW()
		WHERE tanggal = CURRENT_DATE AND (jam_keluar IS NULL OR jam_keluar = '')
	`

	_, err := u.PresensiRepo.DB.Exec(ctx, query, jamKeluar, lat, lon)
	return err
}
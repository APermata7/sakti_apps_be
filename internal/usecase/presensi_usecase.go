package usecase

import (
	"context"
	"errors"
	"log"
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

	log.Printf("config.JamMasuk: %s", config.JamMasuk)

	faceMatch, similarity, err := utils.VerifyFaceWithRepo(ctx, req.SelfieURL, karyawanID, u.KaryawanRepo)
	if err != nil {
		return nil, err
	}
	if !faceMatch {
		return nil, errors.New("wajah tidak dikenali")
	}

	distance := utils.Haversine(config.LatKantor, config.LongKantor, req.Latitude, req.Longitude)
	isOutside := distance > float64(config.RadiusKantor)

	if isOutside {
		return nil, errors.New("Anda berada di luar radius kantor")
	}

	locationStatus := "di_dalam_radius"
	if isOutside {
		locationStatus = "di_luar_radius"
	}

	jamMasuk := time.Now()
	jamMasukStr := jamMasuk.Format("15:04:05")

	wib := time.FixedZone("WIB", 7*60*60)
	jam, _ := time.Parse("15:04:05", config.JamMasuk)
	jamMasukResmi := time.Date(jamMasuk.Year(), jamMasuk.Month(), jamMasuk.Day(), jam.Hour(), jam.Minute(), jam.Second(), 0, wib)

	log.Printf("jamMasuk: %v", jamMasuk)
	log.Printf("jamMasukResmi: %v", jamMasukResmi)
	log.Printf("jamMasuk.After(jamMasukResmi): %v", jamMasuk.After(jamMasukResmi))

	status := "tepat_waktu"
	if jamMasuk.After(jamMasukResmi) {
		status = "terlambat"
	}
	log.Printf("status: %s", status)

	presensi := &domain.Presensi{
		KaryawanID:      karyawanID,
		KantorID:        config.KantorID,
		Tanggal:         jamMasuk,
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
		Tanggal:         presensi.Tanggal.Format("2006-01-02"),
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
	if req.SelfieURL == "" {
		return nil, errors.New("selfie URL wajib diisi untuk check-out")
	}

	if req.Latitude == 0 || req.Longitude == 0 {
		return nil, errors.New("latitude dan longitude wajib diisi")
	}

	presensi, err := u.PresensiRepo.GetToday(ctx, karyawanID)
	if err != nil {
		return nil, errors.New("anda belum melakukan presensi masuk hari ini")
	}
	if presensi == nil || presensi.JamMasuk == "" {
		return nil, errors.New("anda belum melakukan presensi masuk hari ini")
	}
	if presensi.JamKeluar != nil && *presensi.JamKeluar != "" {
		return nil, errors.New("anda sudah melakukan presensi keluar hari ini")
	}

	config, err := u.ConfigRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	faceMatch, _, err := utils.VerifyFaceWithRepo(ctx, req.SelfieURL, karyawanID, u.KaryawanRepo)
	if err != nil {
		return nil, err
	}
	if !faceMatch {
		return nil, errors.New("wajah tidak dikenali")
	}

	jamKeluar := time.Now()
	jamKeluarStr := jamKeluar.Format("15:04:05")

	wib := time.FixedZone("WIB", 7*60*60)
	jam, _ := time.Parse("15:04:05", config.JamPulang)
	jamPulang := time.Date(jamKeluar.Year(), jamKeluar.Month(), jamKeluar.Day(), jam.Hour(), jam.Minute(), jam.Second(), 0, wib)
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

	if err := u.PresensiRepo.UpdateCheckOut(ctx, presensi.ID, jamKeluarStr, lembur, jamLembur, req.Latitude, req.Longitude, req.SelfieURL); err != nil {
		return nil, err
	}

	return &domain.CheckOutResponse{
		ID:              presensi.ID,
		KaryawanID:      presensi.KaryawanID,
		Tanggal:         presensi.Tanggal.Format("2006-01-02"),
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
	log.Printf("karyawanID: %s", karyawanID)

	presensi, err := u.PresensiRepo.GetToday(ctx, karyawanID)

	log.Printf("presensi: %+v", presensi)
	log.Printf("err: %v", err)

	if err != nil {
		log.Printf("Error GetToday: %v", err)
		return &domain.TodayResponse{
			HasCheckedIn:  false,
			HasCheckedOut: false,
			Tanggal:       time.Now().Format("2006-01-02"),
			KaryawanID:    karyawanID,
		}, nil
	}

	if presensi == nil {
		log.Printf("presensi is nil")
		return &domain.TodayResponse{
			HasCheckedIn:  false,
			HasCheckedOut: false,
			Tanggal:       time.Now().Format("2006-01-02"),
			KaryawanID:    karyawanID,
		}, nil
	}

	log.Printf("presensi.JamMasuk: %s", presensi.JamMasuk)

	hasCheckedOut := false
	if presensi.JamKeluar != nil && *presensi.JamKeluar != "" {
		hasCheckedOut = true
	}
	log.Printf("presensi.JamKeluar: %v", presensi.JamKeluar)

	return &domain.TodayResponse{
		HasCheckedIn:  presensi.JamMasuk != "",
		HasCheckedOut: hasCheckedOut,
		CheckInTime:   presensi.JamMasuk,
		CheckInStatus: presensi.Status,
		Tanggal:       presensi.Tanggal.Format("2006-01-02"),
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
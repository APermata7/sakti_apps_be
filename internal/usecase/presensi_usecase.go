package usecase

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"sakti_apps_be/internal/domain"
	"sakti_apps_be/internal/repository"
	"sakti_apps_be/internal/utils"
)

type PresensiUsecase struct {
	PresensiRepo   *repository.PresensiRepo
	KaryawanRepo   *repository.KaryawanRepo
	ConfigRepo     *repository.KonfigurasiRepo
	RiwayatRepo    *repository.RiwayatRepo
	LeaveRepo      *repository.LeaveRepo
	NotifikasiRepo *repository.NotifikasiRepo
	FCMTokenRepo   *repository.FCMTokenRepo
}

func NewPresensiUsecase(
	presensiRepo *repository.PresensiRepo,
	karyawanRepo *repository.KaryawanRepo,
	configRepo *repository.KonfigurasiRepo,
	riwayatRepo *repository.RiwayatRepo,
	leaveRepo *repository.LeaveRepo,
	notifikasiRepo *repository.NotifikasiRepo,
	fcmTokenRepo *repository.FCMTokenRepo,
) *PresensiUsecase {
	return &PresensiUsecase{
		PresensiRepo:   presensiRepo,
		KaryawanRepo:   karyawanRepo,
		ConfigRepo:     configRepo,
		RiwayatRepo:    riwayatRepo,
		LeaveRepo:      leaveRepo,
		NotifikasiRepo: notifikasiRepo,
		FCMTokenRepo:   fcmTokenRepo,
	}
}

func formatTanggalIndonesia(t time.Time) string {
	bulan := map[string]string{
		"January": "Januari", "February": "Februari", "March": "Maret",
		"April": "April", "May": "Mei", "June": "Juni",
		"July": "Juli", "August": "Agustus", "September": "September",
		"October": "Oktober", "November": "November", "December": "Desember",
	}
	bulanInggris := t.Format("January")
	bulanIndo := bulan[bulanInggris]
	return t.Format("02 " + bulanIndo + " 2006 | 15:04")
}

func extractPublicID(url string) string {
	parts := strings.Split(url, "/upload/")
	if len(parts) < 2 {
		return ""
	}
	path := parts[1]
	if idx := strings.Index(path, "/"); idx != -1 {
		path = path[idx+1:]
	}
	if idx := strings.LastIndex(path, "."); idx != -1 {
		path = path[:idx]
	}
	return path
}

func (u *PresensiUsecase) CheckIn(ctx context.Context, karyawanID string, req domain.CheckInRequest) (*domain.CheckInResponse, error) {
	log.Printf("CheckIn dimulai untuk karyawanID: %s", karyawanID)

	existing, _ := u.PresensiRepo.GetToday(ctx, karyawanID)
	if existing != nil && existing.JamMasuk != "" {
		log.Printf("Karyawan sudah melakukan check-in hari ini")
		return nil, errors.New("anda sudah melakukan presensi masuk hari ini")
	}

	karyawan, err := u.KaryawanRepo.GetByID(ctx, karyawanID)
	if err != nil || karyawan == nil {
		log.Printf("Karyawan tidak ditemukan: %v", err)
		return nil, errors.New("karyawan tidak ditemukan")
	}

	if karyawan.StatusKaryawan != "aktif" {
		log.Printf("Akun tidak aktif")
		return nil, errors.New("akun tidak aktif")
	}

	config, err := u.ConfigRepo.GetActive(ctx)
	if err != nil {
		log.Printf("Error mengambil konfigurasi: %v", err)
		return nil, err
	}

	log.Printf("Memulai verifikasi wajah")
	faceMatch, similarity, err := utils.VerifyFaceWithRepo(ctx, req.SelfieURL, karyawanID, u.KaryawanRepo)
	if err != nil {
		log.Printf("Error verifikasi wajah: %v", err)
		return nil, err
	}
	if !faceMatch {
		log.Printf("Wajah tidak dikenali")
		return nil, errors.New("wajah tidak dikenali")
	}
	log.Printf("Wajah cocok dengan similarity: %f", similarity)

	distance := utils.Haversine(config.LatKantor, config.LongKantor, req.Latitude, req.Longitude)
	isOutside := distance > float64(config.RadiusKantor)
	log.Printf("Jarak: %f meter, Radius: %d meter, Di luar radius: %t", distance, config.RadiusKantor, isOutside)

	locationStatusMasuk := "di_dalam_radius"
	if isOutside {
		locationStatusMasuk = "di_luar_radius"
	}

	jamMasuk := time.Now()
	jamMasukStr := jamMasuk.Format("15:04:05")

	wib := time.FixedZone("WIB", 7*60*60)
	jam, _ := time.Parse("15:04:05", config.JamMasuk)
	jamMasukResmi := time.Date(jamMasuk.Year(), jamMasuk.Month(), jamMasuk.Day(), jam.Hour(), jam.Minute(), jam.Second(), 0, wib)

	status := "tepat_waktu"
	if jamMasuk.After(jamMasukResmi) {
		status = "terlambat"
	}
	log.Printf("Jam masuk: %s, Jam masuk resmi: %s, Status: %s", jamMasukStr, jamMasukResmi.Format("15:04:05"), status)

	presensi := &domain.Presensi{
		KaryawanID:           karyawanID,
		KantorID:             config.KantorID,
		Tanggal:              jamMasuk,
		JamMasuk:             jamMasukStr,
		Status:               status,
		LintangMasuk:         req.Latitude,
		BujurMasuk:           req.Longitude,
		ValidasiWajah:        faceMatch,
		FaceSimilarity:       similarity,
		URLFoto:              req.SelfieURL,
		AlasanTerlambat:      nil,
		DistanceMeter:        distance,
		IsOutsideRadius:      isOutside,
		LocationStatusMasuk:  locationStatusMasuk,
		LocationStatusKeluar: nil,
	}

	log.Printf("Menyimpan data presensi")
	if err := u.PresensiRepo.Create(ctx, presensi); err != nil {
		log.Printf("Error menyimpan presensi: %v", err)
		return nil, err
	}
	log.Printf("Presensi berhasil disimpan dengan ID: %s", presensi.ID)

	if u.RiwayatRepo != nil {
		detail := "Check-in pada " + jamMasukStr + " dengan status " + status + " di lokasi " + locationStatusMasuk
		u.RiwayatRepo.CreateRiwayat(ctx, karyawanID, "check_in", detail)
	}

	return &domain.CheckInResponse{
		ID:                   presensi.ID,
		KaryawanID:           presensi.KaryawanID,
		Tanggal:              presensi.Tanggal.Format("2006-01-02"),
		JamMasuk:             presensi.JamMasuk,
		Status:               presensi.Status,
		ValidasiWajah:        presensi.ValidasiWajah,
		FaceSimilarity:       similarity,
		URLFoto:              presensi.URLFoto,
		DistanceMeter:        presensi.DistanceMeter,
		IsOutsideRadius:      presensi.IsOutsideRadius,
		LocationStatusMasuk:  locationStatusMasuk,
		OfficeLatitude:       config.LatKantor,
		OfficeLongitude:      config.LongKantor,
		OfficeRadius:         config.RadiusKantor,
	}, nil
}

func (u *PresensiUsecase) CheckOut(ctx context.Context, karyawanID string, req domain.CheckOutRequest) (*domain.CheckOutResponse, error) {
	log.Printf("CheckOut dimulai untuk karyawanID: %s", karyawanID)

	if req.SelfieURL == "" {
		log.Printf("URL selfie kosong")
		return nil, errors.New("selfie URL wajib diisi untuk check-out")
	}

	if req.Latitude == 0 || req.Longitude == 0 {
		log.Printf("Koordinat tidak valid: Lat=%f, Lon=%f", req.Latitude, req.Longitude)
		return nil, errors.New("latitude dan longitude wajib diisi")
	}

	presensi, err := u.PresensiRepo.GetToday(ctx, karyawanID)
	if err != nil {
		log.Printf("Error mengambil data presensi hari ini: %v", err)
		return nil, errors.New("anda belum melakukan presensi masuk hari ini")
	}
	if presensi == nil || presensi.JamMasuk == "" {
		log.Printf("Tidak ada check-in untuk hari ini")
		return nil, errors.New("anda belum melakukan presensi masuk hari ini")
	}
	if presensi.JamKeluar != nil && *presensi.JamKeluar != "" {
		log.Printf("Karyawan sudah check-out")
		return nil, errors.New("anda sudah melakukan presensi keluar hari ini")
	}

	config, err := u.ConfigRepo.GetActive(ctx)
	if err != nil {
		log.Printf("Error mengambil konfigurasi: %v", err)
		return nil, err
	}

	log.Printf("Memulai verifikasi wajah untuk check-out")
	faceMatch, _, err := utils.VerifyFaceWithRepo(ctx, req.SelfieURL, karyawanID, u.KaryawanRepo)
	if err != nil {
		log.Printf("Error verifikasi wajah: %v", err)
		return nil, err
	}
	if !faceMatch {
		log.Printf("Wajah tidak dikenali untuk check-out")
		return nil, errors.New("wajah tidak dikenali")
	}

	jamKeluar := time.Now()
	jamKeluarStr := jamKeluar.Format("15:04:05")

	wib := time.FixedZone("WIB", 7*60*60)
	jam, _ := time.Parse("15:04:05", config.JamPulang)
	jamPulangTime := time.Date(jamKeluar.Year(), jamKeluar.Month(), jamKeluar.Day(), jam.Hour(), jam.Minute(), jam.Second(), 0, wib)

	var lembur bool
	var jamLembur float64

	if req.Lembur {
		lembur = true
		if jamKeluar.After(jamPulangTime) {
			selisih := jamKeluar.Sub(jamPulangTime)
			jamLembur = selisih.Hours()
		} else {
			jamLembur = 0
		}
	} else {
		lembur = false
		jamLembur = 0
	}
	log.Printf("Lembur: %t, Jam lembur: %f", lembur, jamLembur)

	distance := utils.Haversine(config.LatKantor, config.LongKantor, req.Latitude, req.Longitude)
	isOutside := distance > float64(config.RadiusKantor)
	log.Printf("Jarak: %f meter, Radius: %d meter, Di luar radius: %t", distance, config.RadiusKantor, isOutside)

	locationStatusKeluar := "di_dalam_radius"
	if isOutside {
		locationStatusKeluar = "di_luar_radius"
	}

	log.Printf("Memperbarui data check-out")
	if err := u.PresensiRepo.UpdateCheckOut(ctx, presensi.ID, jamKeluarStr, lembur, jamLembur, req.Latitude, req.Longitude, req.SelfieURL, distance, isOutside, &locationStatusKeluar); err != nil {
		log.Printf("Error update check-out: %v", err)
		return nil, err
	}
	log.Printf("Check-out berhasil diperbarui untuk presensi ID: %s", presensi.ID)

	if u.RiwayatRepo != nil {
		statusLembur := ""
		if lembur {
			statusLembur = " dengan lembur " + strconv.FormatFloat(jamLembur, 'f', 1, 64) + " jam"
		}
		detail := "Check-out pada " + jamKeluarStr + statusLembur + " di lokasi " + locationStatusKeluar
		u.RiwayatRepo.CreateRiwayat(ctx, karyawanID, "check_out", detail)
	}

	return &domain.CheckOutResponse{
		ID:                    presensi.ID,
		KaryawanID:            presensi.KaryawanID,
		Tanggal:               presensi.Tanggal.Format("2006-01-02"),
		JamMasuk:              presensi.JamMasuk,
		JamKeluar:             jamKeluarStr,
		Lembur:                lembur,
		JamLembur:             jamLembur,
		DistanceMeter:         distance,
		IsOutsideRadius:       isOutside,
		LocationStatusKeluar:  locationStatusKeluar,
	}, nil
}

func (u *PresensiUsecase) GetToday(ctx context.Context, karyawanID string) (*domain.TodayResponse, error) {
	log.Printf("GetToday dimulai untuk karyawanID: %s", karyawanID)

	presensi, err := u.PresensiRepo.GetToday(ctx, karyawanID)

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
		log.Printf("Tidak ada data presensi untuk hari ini")
		return &domain.TodayResponse{
			HasCheckedIn:  false,
			HasCheckedOut: false,
			Tanggal:       time.Now().Format("2006-01-02"),
			KaryawanID:    karyawanID,
		}, nil
	}

	hasCheckedOut := false
	if presensi.JamKeluar != nil && *presensi.JamKeluar != "" {
		hasCheckedOut = true
	}

	log.Printf("Hasil GetToday - CheckIn: %t, CheckOut: %t", presensi.JamMasuk != "", hasCheckedOut)

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
	log.Printf("GetHistory dimulai untuk karyawanID: %s", karyawanID)

	offset := (page - 1) * limit
	items, total, err := u.PresensiRepo.GetHistory(ctx, karyawanID, startDate, endDate, status, limit, offset)

	if err != nil {
		log.Printf("Error GetHistory: %v", err)
		return nil, 0, err
	}

	log.Printf("Hasil GetHistory - items: %d, total: %d", len(items), total)
	return items, total, nil
}

func (u *PresensiUsecase) UpdateAlasanTerlambat(ctx context.Context, karyawanID, alasan string) error {
	log.Printf("UpdateAlasanTerlambat dimulai untuk karyawanID: %s", karyawanID)

	presensi, err := u.PresensiRepo.GetToday(ctx, karyawanID)
	if err != nil || presensi == nil {
		log.Printf("Data presensi tidak ditemukan: %v", err)
		return errors.New("data presensi hari ini tidak ditemukan")
	}

	if presensi.Status != "terlambat" {
		log.Printf("Karyawan tidak terlambat, status: %s", presensi.Status)
		return errors.New("anda tidak terlambat")
	}

	log.Printf("Memperbarui alasan terlambat untuk presensi ID: %s", presensi.ID)
	return u.PresensiRepo.UpdateAlasanTerlambat(ctx, presensi.ID, alasan)
}

func (u *PresensiUsecase) AutoClockOut(ctx context.Context) error {
	log.Printf("AutoClockOut dimulai")

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

	result, err := u.PresensiRepo.DB.Exec(ctx, query, jamKeluar, lat, lon)
	if err != nil {
		log.Printf("Error AutoClockOut: %v", err)
		return err
	}

	rowsAffected := result.RowsAffected()
	log.Printf("AutoClockOut selesai, baris terpengaruh: %d", rowsAffected)
	return nil
}

func (u *PresensiUsecase) SendPresensiReminder(ctx context.Context) error {
	now := time.Now()
	weekday := now.Weekday()

	if weekday == time.Saturday || weekday == time.Sunday {
		return nil
	}

	libur, err := u.LeaveRepo.GetLiburByDate(ctx, now)
	if err == nil && libur != nil && libur.Aktif {
		return nil
	}

	config, err := u.ConfigRepo.GetActive(ctx)
	if err != nil {
		return err
	}

	karyawanCuti, err := u.LeaveRepo.GetActiveLeavesByDate(ctx, now)
	if err != nil {
		karyawanCuti = []string{}
	}
	karyawanCutiMap := make(map[string]bool)
	for _, id := range karyawanCuti {
		karyawanCutiMap[id] = true
	}

	jamSekarang := now.Format("15:04")

	semuaKaryawan, _, err := u.KaryawanRepo.GetAll(ctx, 1000, 0, "", "", "aktif")
	if err != nil {
		return err
	}

	tanggalJam := formatTanggalIndonesia(now)

	if jamSekarang >= config.JamMinimalMasuk {
		for _, karyawan := range semuaKaryawan {
			if karyawanCutiMap[karyawan.ID] {
				continue
			}
			alreadyCheckedIn, _ := u.PresensiRepo.AlreadyCheckedIn(ctx, karyawan.ID, now)
			if !alreadyCheckedIn {
				judul := "Presensi Masuk"
				pesan := "Segera lakukan presensi masuk. Jika Anda melakukan presensi masuk setelah pukul 08.30 WIB presensi akan dihitung sebagai terlambat.\n\n" + tanggalJam

				go u.KirimInApp(ctx, domain.KirimNotifikasiRequest{
					KaryawanID:    karyawan.ID,
					Jenis:         "reminder",
					Judul:         judul,
					Pesan:         pesan,
					ReferensiID:   "",
					ReferensiTipe: "presensi",
				})
			}
		}
	}

	if jamSekarang >= config.JamMinimalPulang {
		for _, karyawan := range semuaKaryawan {
			if karyawanCutiMap[karyawan.ID] {
				continue
			}
			alreadyCheckedIn, _ := u.PresensiRepo.AlreadyCheckedIn(ctx, karyawan.ID, now)
			alreadyCheckedOut, _ := u.PresensiRepo.AlreadyCheckedOut(ctx, karyawan.ID, now)
			if alreadyCheckedIn && !alreadyCheckedOut {
				judul := "Presensi Keluar"
				pesan := "Segera lakukan presensi keluar. Lakukan presensi keluar sebelum pukul 17.00 WIB. Jika sedang lembur, lakukan presensi sesuai waktu Anda pulang.\n\n" + tanggalJam

				go u.KirimInApp(ctx, domain.KirimNotifikasiRequest{
					KaryawanID:    karyawan.ID,
					Jenis:         "reminder",
					Judul:         judul,
					Pesan:         pesan,
					ReferensiID:   "",
					ReferensiTipe: "presensi",
				})
			}
		}
	}

	return nil
}

func (u *PresensiUsecase) KirimInApp(ctx context.Context, req domain.KirimNotifikasiRequest) error {
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

func (u *PresensiUsecase) CleanupOldPhotos(ctx context.Context) error {
	log.Println("CleanupOldPhotos dimulai")

	query := `
		SELECT url_foto FROM presensi 
		WHERE tanggal < NOW() - INTERVAL '2 months'
		AND url_foto IS NOT NULL AND url_foto != ''
	`
	rows, err := u.PresensiRepo.DB.Query(ctx, query)
	if err != nil {
		log.Printf("CleanupOldPhotos query error: %v", err)
		return err
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			continue
		}
		urls = append(urls, url)
	}

	log.Printf("CleanupOldPhotos found %d old photos", len(urls))

	for _, url := range urls {
		publicID := extractPublicID(url)
		if publicID != "" {
			if err := utils.DeleteFile(publicID); err != nil {
				log.Printf("CleanupOldPhotos delete file error: %v", err)
			}
		}
	}

	deleteQuery := `
		UPDATE presensi 
		SET url_foto = NULL 
		WHERE tanggal < NOW() - INTERVAL '2 months'
	`
	_, err = u.PresensiRepo.DB.Exec(ctx, deleteQuery)
	if err != nil {
		log.Printf("CleanupOldPhotos update error: %v", err)
		return err
	}

	log.Println("CleanupOldPhotos selesai")
	return nil
}
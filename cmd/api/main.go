package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"sakti_apps_be/internal/handler"
	"sakti_apps_be/internal/middleware"
	"sakti_apps_be/internal/repository"
	"sakti_apps_be/internal/usecase"
	"sakti_apps_be/internal/utils"
	"sakti_apps_be/pkg/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	log.Println("Menghubungkan ke database...")
	dbConn, err := db.NewSupabaseDB()
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}
	defer dbConn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dbConn.Ping(ctx); err != nil {
		log.Fatalf("Database tidak merespon: %v", err)
	}
	log.Println("Database connected")

	if err := utils.InitCloudinary(); err != nil {
		log.Printf("Cloudinary init error: %v", err)
	} else {
		log.Println("Cloudinary ready")
	}

	if err := utils.InitFCM(); err != nil {
		log.Printf("FCM init error: %v", err)
	} else {
		log.Println("FCM ready")
	}

	karyawanRepo := repository.NewKaryawanRepo(dbConn.Pool)
	presensiRepo := repository.NewPresensiRepo(dbConn.Pool)
	configRepo := repository.NewKonfigurasiRepo(dbConn.Pool)
	leaveRepo := repository.NewLeaveRepo(dbConn.Pool)
	riwayatRepo := repository.NewRiwayatRepo(dbConn.Pool)
	notifikasiRepo := repository.NewNotifikasiRepo(dbConn.Pool)
	fcmTokenRepo := repository.NewFCMTokenRepo(dbConn.Pool)
	liburRepo := repository.NewLiburRepo(dbConn.Pool)
	konfigurasiRepo := repository.NewKonfigurasiRepo(dbConn.Pool)
	ttdRepo := repository.NewTTDRepo(dbConn.Pool)
	telegramRepo := repository.NewTelegramRepo(dbConn.Pool)

	authUsecase := usecase.NewAuthUsecase(karyawanRepo, riwayatRepo)
	presensiUsecase := usecase.NewPresensiUsecase(
		presensiRepo,
		karyawanRepo,
		configRepo,
		riwayatRepo,
		leaveRepo,
		notifikasiRepo,
		fcmTokenRepo,
	)
	notificationUsecase := usecase.NewNotificationUsecase(fcmTokenRepo, notifikasiRepo, karyawanRepo, dbConn.Pool)
	leaveUsecase := usecase.NewLeaveUsecase(leaveRepo, karyawanRepo, ttdRepo, konfigurasiRepo, riwayatRepo, notificationUsecase)
	riwayatUsecase := usecase.NewRiwayatUsecase(riwayatRepo)
	liburUsecase := usecase.NewLiburUsecase(liburRepo)
	konfigurasiUsecase := usecase.NewKonfigurasiUsecase(konfigurasiRepo)
	ttdUsecase := usecase.NewTTDUsecase(ttdRepo)
	telegramUsecase := usecase.NewTelegramUsecase(karyawanRepo, telegramRepo, dbConn.Pool)

	adminUsecase := usecase.NewAdminUsecase(
		dbConn.Pool,
		karyawanRepo,
		presensiRepo,
		leaveRepo,
		liburUsecase,
		konfigurasiRepo,
	)

	authHandler := handler.NewAuthHandler(authUsecase)
	presensiHandler := handler.NewPresensiHandler(presensiUsecase)
	leaveHandler := handler.NewLeaveHandler(leaveUsecase)
	riwayatHandler := handler.NewRiwayatHandler(riwayatUsecase)
	notificationHandler := handler.NewNotificationHandler(notificationUsecase)
	liburHandler := handler.NewLiburHandler(liburUsecase)
	konfigurasiHandler := handler.NewKonfigurasiHandler(konfigurasiUsecase)
	adminHandler := handler.NewAdminHandler(adminUsecase)
	ttdHandler := handler.NewTTDHandler(ttdUsecase)
	telegramHandler := handler.NewTelegramHandler(telegramUsecase)
	telegramWebhookHandler := handler.NewTelegramWebhookHandler(telegramUsecase)

	app := fiber.New(fiber.Config{
		AppName: os.Getenv("APP_NAME"),
	})

	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": os.Getenv("APP_NAME"),
			"db":      "connected",
		})
	})

	app.Post("/upload/file", handler.UploadFile)
	app.Post("/upload/image", handler.UploadImage)
	app.Post("/upload/ttd", handler.UploadTTD)

	api := app.Group("/api")

	api.Use(middleware.LoginRateLimiter())

	api.Post("/auth/login", authHandler.Login)
	api.Post("/auth/forgot-password", authHandler.ForgotPassword)
	api.Post("/auth/reset-password", authHandler.ResetPassword)

	api.Get("/libur/check", liburHandler.IsHoliday)

	api.Post("/webhook/telegram", telegramWebhookHandler.Webhook)

	protected := api.Group("/", middleware.AuthMiddleware())
	protected.Get("/auth/me", authHandler.GetProfile)
	protected.Post("/auth/logout", authHandler.Logout)
	protected.Put("/auth/change-password", authHandler.ChangePassword)

	protected.Post("/attendance/check-in", presensiHandler.CheckIn)
	protected.Post("/attendance/check-out", presensiHandler.CheckOut)
	protected.Get("/attendance/today", presensiHandler.GetToday)
	protected.Get("/attendance/history", presensiHandler.GetHistory)
	protected.Put("/attendance/check-in/reason", presensiHandler.UpdateAlasanTerlambat)
	protected.Get("/attendance/work-config", konfigurasiHandler.GetWorkConfig)

	protected.Post("/leave/request", leaveHandler.CreateLeave)
	protected.Get("/leave/status", leaveHandler.GetStatus)
	protected.Get("/leave/balance", leaveHandler.GetBalance)
	protected.Get("/leave/all", leaveHandler.GetAllLeaves)
	protected.Get("/leave/approval/list", leaveHandler.GetApprovalList)
	protected.Get("/leave/finalization/list", leaveHandler.GetFinalizationList)
	protected.Get("/leave/:id/download", leaveHandler.DownloadSuratCuti)
	protected.Put("/leave/:id/cancel", leaveHandler.CancelLeave)
	protected.Put("/leave/:id/approve", leaveHandler.ApproveLeave)
	protected.Put("/leave/:id/reject", leaveHandler.RejectLeave)
	protected.Put("/leave/:id/finalize", leaveHandler.FinalizeLeave)

	protected.Get("/riwayat", riwayatHandler.GetRiwayat)

	protected.Get("/notifikasi", notificationHandler.GetNotifikasi)
	protected.Get("/notifikasi/unread", notificationHandler.GetUnreadCount)
	protected.Put("/notifikasi/:id/read", notificationHandler.MarkAsRead)
	protected.Put("/notifikasi/read-all", notificationHandler.MarkAllAsRead)

	protected.Post("/fcm/register", notificationHandler.RegisterFCMToken)
	protected.Delete("/fcm/token", notificationHandler.DeactivateFCMToken)

	protected.Get("/telegram/status", telegramHandler.GetStatus)
	protected.Post("/telegram/connect", telegramHandler.Connect)
	protected.Delete("/telegram/disconnect", telegramHandler.Disconnect)

	protected.Get("/libur", liburHandler.GetActiveList)

	protected.Post("/ttd/upload", ttdHandler.UploadTTD)
	protected.Get("/ttd", ttdHandler.GetTTD)
	protected.Put("/ttd", ttdHandler.UpdateTTD)
	protected.Delete("/ttd", ttdHandler.DeleteTTD)
	protected.Get("/ttd/verify", ttdHandler.VerifyTTD)

	adminGroup := protected.Group("/admin", middleware.RequireRole("admin"))
	adminGroup.Get("/dashboard", adminHandler.GetDashboard)
	adminGroup.Get("/dashboard/chart", adminHandler.GetDashboardChart)

	adminGroup.Post("/karyawan", adminHandler.CreateKaryawan)
	adminGroup.Get("/karyawan", adminHandler.GetAllKaryawan)
	adminGroup.Get("/karyawan/:id", adminHandler.GetKaryawan)
	adminGroup.Put("/karyawan/:id", adminHandler.UpdateKaryawan)
	adminGroup.Delete("/karyawan/:id", adminHandler.DeleteKaryawan)
	adminGroup.Put("/karyawan/:id/activate", adminHandler.ActivateKaryawan)

	adminGroup.Get("/presensi", adminHandler.GetPresensiReport)
	adminGroup.Get("/presensi/export", adminHandler.ExportPresensiCSV)
	adminGroup.Get("/cuti", adminHandler.GetCutiReport)
	adminGroup.Get("/cuti/export", adminHandler.ExportCutiCSV)

	adminGroup.Post("/libur", adminHandler.CreateLibur)
	adminGroup.Get("/libur", adminHandler.GetAllLibur)
	adminGroup.Put("/libur/:id", adminHandler.UpdateLibur)
	adminGroup.Delete("/libur/:id", adminHandler.DeleteLibur)
	adminGroup.Put("/libur/:id/toggle", adminHandler.ToggleLibur)

	adminGroup.Get("/konfigurasi-kerja", adminHandler.GetKonfigurasi)
	adminGroup.Put("/konfigurasi-kerja", adminHandler.UpdateKonfigurasi)

	adminGroup.Get("/karyawan/export", adminHandler.ExportKaryawanCSV)

	go func() {
		for {
			time.Sleep(1 * time.Hour)
			ctx := context.Background()
			if err := telegramUsecase.CleanupExpiredCodes(ctx); err != nil {
				log.Printf("Cleanup expired codes failed: %v", err)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(24 * time.Hour)
			ctx := context.Background()
			if err := presensiUsecase.CleanupOldPhotos(ctx); err != nil {
				log.Printf("Cleanup old photos failed: %v", err)
			}
		}
	}()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
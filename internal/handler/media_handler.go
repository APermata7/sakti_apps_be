package handler

import (
	"log"
	"mime/multipart"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"sakti_apps_be/internal/repository"
	"sakti_apps_be/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File not found",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to open file",
		})
	}
	defer src.Close()

	url, err := utils.UploadFile(src, file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"url":     url,
	})
}

func sanitizePublicID(name string) string {
	replacer := strings.NewReplacer(
		" ", "_",
		"-", "_",
		"\n", "",
		"\r", "",
		"\t", "",
		"(", "",
		")", "",
		"[", "",
		"]", "",
		"{", "",
		"}", "",
		"&", "",
		"@", "",
		"#", "",
		"$", "",
		"%", "",
		"^", "",
		"*", "",
		"+", "",
		"=", "",
		"?", "",
		"!", "",
		"'", "",
		`"`, "",
		":", "",
		";", "",
		"<", "",
		">", "",
		"/", "",
		"\\", "",
		"|", "",
		"`", "",
		"~", "",
	)
	result := replacer.Replace(name)
	for strings.Contains(result, "__") {
		result = strings.ReplaceAll(result, "__", "_")
	}
	return strings.Trim(result, "_")
}

func validateFile(c *fiber.Ctx, key string, maxSize int64) (*multipart.FileHeader, error) {
	file, err := c.FormFile(key)
	if err != nil {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File tidak ditemukan",
		})
	}

	if file.Size > maxSize {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Ukuran file maksimal " + strconv.FormatInt(maxSize/1024/1024, 10) + " MB",
		})
	}

	if file.Size == 0 {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File kosong",
		})
	}

	src, err := file.Open()
	if err != nil {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuka file",
		})
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membaca file",
		})
	}

	mimeType := http.DetectContentType(buffer)
	log.Printf("validateFile: mimeType=%s, filename=%s, size=%d", mimeType, file.Filename, file.Size)

	validTypes := []string{
		"image/jpeg",
		"image/png",
		"image/webp",
	}

	isValid := false
	for _, t := range validTypes {
		if mimeType == t {
			isValid = true
			break
		}
	}

	if !isValid {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format file tidak valid. Gunakan JPEG, PNG, atau WebP",
		})
	}

	return file, nil
}

func UploadImage(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("UploadImage panic: %v", r)
			log.Printf("Stack: %s", debug.Stack())
		}
	}()

	log.Println("UploadImage: start")

	role, ok := c.Locals("role").(string)
	if !ok || role == "" {
		log.Println("UploadImage: role not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		log.Println("UploadImage: user_id not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	karyawanID := c.FormValue("karyawan_id")
	log.Printf("UploadImage: role=%s, userID=%s, karyawanID=%s", role, userID, karyawanID)

	if !(role == "admin" && karyawanID != "") {
		karyawanID = userID
	}
	log.Printf("UploadImage: final karyawanID=%s", karyawanID)

	file, err := validateFile(c, "image", 5*1024*1024)
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to open image",
		})
	}
	defer src.Close()

	log.Println("UploadImage: uploading to Cloudinary")
	url, err := utils.UploadImage(src, file.Filename)
	if err != nil {
		log.Printf("UploadImage: Cloudinary error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal upload: " + err.Error(),
		})
	}
	log.Printf("UploadImage: Cloudinary URL: %s", url)

	dbPool, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok || dbPool == nil {
		log.Printf("UploadImage: Database connection not found")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Database connection not found",
		})
	}
	log.Println("UploadImage: database connection OK")

	karyawanRepo := repository.NewKaryawanRepo(dbPool)
	err = karyawanRepo.UpdateFotoURL(c.Context(), karyawanID, url)
	if err != nil {
		log.Printf("UploadImage: UpdateFotoURL error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal simpan ke database: " + err.Error(),
		})
	}
	log.Println("UploadImage: success")

	return c.JSON(fiber.Map{
		"success":     true,
		"url":         url,
		"karyawan_id": karyawanID,
	})
}

func UploadLogoHandler(c *fiber.Ctx) error {
	log.Println("UploadLogoHandler: start")

	file, err := validateFile(c, "image", 5*1024*1024)
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to open image",
		})
	}
	defer src.Close()

	log.Println("UploadLogoHandler: uploading to Cloudinary logos folder")
	url, err := utils.UploadLogo(src, file.Filename)
	if err != nil {
		log.Printf("UploadLogoHandler: Cloudinary error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal upload logo: " + err.Error(),
		})
	}
	log.Printf("UploadLogoHandler: Cloudinary URL: %s", url)

	return c.JSON(fiber.Map{
		"success": true,
		"url":     url,
	})
}

func UploadPresensi(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("UploadPresensi panic: %v", r)
			log.Printf("Stack: %s", debug.Stack())
		}
	}()

	log.Println("UploadPresensi: start")

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		log.Println("UploadPresensi: user_id not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	dbPool, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok || dbPool == nil {
		log.Println("UploadPresensi: database connection not found")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Database connection error",
		})
	}

	karyawanRepo := repository.NewKaryawanRepo(dbPool)
	karyawan, err := karyawanRepo.GetByID(c.Context(), userID)
	if err != nil || karyawan == nil {
		log.Printf("UploadPresensi: karyawan not found")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Karyawan tidak ditemukan",
		})
	}

	var req struct {
		Tipe string `json:"tipe"`
	}
	if err := c.BodyParser(&req); err != nil {
		req.Tipe = "in"
	}

	tipe := req.Tipe
	if tipe != "in" && tipe != "out" {
		tipe = "in"
	}

	file, err := validateFile(c, "image", 5*1024*1024)
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to open image",
		})
	}
	defer src.Close()

	log.Printf("UploadPresensi: uploading to Cloudinary presensi folder, tipe=%s, filename=%s, size=%d", tipe, file.Filename, file.Size)

	sanitizedName := sanitizePublicID(karyawan.NamaLengkap)
	url, err := utils.UploadPresensi(src, file.Filename, sanitizedName, tipe)
	if err != nil {
		log.Printf("UploadPresensi: Cloudinary error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal upload presensi: " + err.Error(),
		})
	}

	if url == "" {
		log.Printf("UploadPresensi: Cloudinary URL kosong")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal upload presensi, URL foto tidak ditemukan dari server",
		})
	}

	log.Printf("UploadPresensi: Cloudinary URL: %s", url)

	return c.JSON(fiber.Map{
		"success": true,
		"url":     url,
		"tipe":    tipe,
	})
}
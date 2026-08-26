package utils

import (
	"bytes"
	"context"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var Cld *cloudinary.Cloudinary

func InitCloudinary() error {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		log.Println("Cloudinary credentials not found, skipping init")
		return nil
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return err
	}

	cld.Config.URL.Secure = true
	Cld = cld
	log.Println("Cloudinary initialized successfully")
	return nil
}

func sanitizePublicID(name string) string {
	if name == "" {
		return "file"
	}

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

func sanitizeFilename(filename string) string {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.ReplaceAll(name, "(", "")
	name = strings.ReplaceAll(name, ")", "")
	name = strings.ReplaceAll(name, "[", "")
	name = strings.ReplaceAll(name, "]", "")
	name = strings.ReplaceAll(name, "{", "")
	name = strings.ReplaceAll(name, "}", "")
	name = strings.ToLower(name)
	if name == "" {
		name = "file"
	}
	return name
}

func UploadFile(file multipart.File, filename string) (string, error) {
	if Cld == nil {
		return "", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	folder := os.Getenv("CLOUDINARY_UPLOAD_FOLDER")
	if folder == "" {
		folder = "sakti-apps"
	}

	safeFilename := sanitizeFilename(filename)

	resp, err := Cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         folder,
		PublicID:       safeFilename,
		UseFilename:    api.Bool(false),
		UniqueFilename: api.Bool(true),
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

func UploadImage(file multipart.File, filename string) (string, error) {
	if Cld == nil {
		return "", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	folder := os.Getenv("CLOUDINARY_UPLOAD_FOLDER")
	if folder == "" {
		folder = "sakti-apps"
	}
	folder = folder + "/photos"

	safeFilename := sanitizeFilename(filename)

	resp, err := Cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         folder,
		PublicID:       safeFilename,
		UseFilename:    api.Bool(false),
		UniqueFilename: api.Bool(true),
		ResourceType:   "image",
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

func UploadImageWithCustomFolder(file multipart.File, folderPath, filename string) (string, error) {
	if Cld == nil {
		return "", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	baseFolder := os.Getenv("CLOUDINARY_UPLOAD_FOLDER")
	if baseFolder == "" {
		baseFolder = "sakti-apps"
	}

	folder := baseFolder
	if folderPath != "" {
		folder = baseFolder + "/" + folderPath
	}

	safeFilename := sanitizeFilename(filename)

	resp, err := Cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         folder,
		PublicID:       safeFilename,
		UseFilename:    api.Bool(false),
		UniqueFilename: api.Bool(true),
		ResourceType:   "image",
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

func UploadTTD(file multipart.File, filename string) (string, error) {
	if Cld == nil {
		return "", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	folder := os.Getenv("CLOUDINARY_UPLOAD_FOLDER")
	if folder == "" {
		folder = "sakti-apps"
	}
	folder = folder + "/ttd"

	safeFilename := sanitizeFilename(filename)

	resp, err := Cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         folder,
		PublicID:       safeFilename,
		UseFilename:    api.Bool(false),
		UniqueFilename: api.Bool(true),
		ResourceType:   "image",
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

func UploadLogo(file multipart.File, filename string) (string, error) {
	if Cld == nil {
		return "", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	folder := os.Getenv("CLOUDINARY_UPLOAD_FOLDER")
	if folder == "" {
		folder = "sakti-apps"
	}
	folder = folder + "/logos"

	safeFilename := sanitizeFilename(filename)

	resp, err := Cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         folder,
		PublicID:       safeFilename,
		UseFilename:    api.Bool(false),
		UniqueFilename: api.Bool(true),
		ResourceType:   "image",
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

func UploadPresensi(file multipart.File, filename string, karyawanNama, tipe string) (string, error) {
	if Cld == nil {
		log.Println("UploadPresensi: Cloudinary not initialized")
		return "", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if karyawanNama == "" {
		log.Println("UploadPresensi: karyawanNama is empty")
		return "", nil
	}

	now := time.Now()
	tanggal := now.Format("02012006")
	sanitizedName := sanitizePublicID(strings.ToLower(karyawanNama))
	publicID := sanitizedName + "_" + tipe + "_" + tanggal

	folder := os.Getenv("CLOUDINARY_UPLOAD_FOLDER")
	if folder == "" {
		folder = "sakti-apps"
	}
	folder = folder + "/presensi"

	log.Printf("UploadPresensi: folder=%s, publicID=%s, filename=%s", folder, publicID, filename)

	resp, err := Cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         folder,
		PublicID:       publicID,
		UseFilename:    api.Bool(false),
		UniqueFilename: api.Bool(false),
		Overwrite:      api.Bool(true),
		ResourceType:   "image",
	})
	if err != nil {
		log.Printf("UploadPresensi: Cloudinary upload error: %v", err)
		return "", err
	}

	if resp == nil {
		log.Println("UploadPresensi: Cloudinary response is nil")
		return "", nil
	}

	log.Printf("UploadPresensi: response PublicID=%s, SecureURL=%s", resp.PublicID, resp.SecureURL)

	if resp.SecureURL == "" {
		log.Printf("UploadPresensi: Cloudinary URL kosong")
		return "", nil
	}

	return resp.SecureURL, nil
}

func UploadPDF(data []byte, filename string) (string, error) {
	if Cld == nil {
		return "", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	folder := os.Getenv("CLOUDINARY_UPLOAD_FOLDER")
	if folder == "" {
		folder = "sakti-apps"
	}
	folder = folder + "/documents"

	safeFilename := sanitizeFilename(filename)

	reader := bytes.NewReader(data)

	resp, err := Cld.Upload.Upload(ctx, reader, uploader.UploadParams{
		Folder:         folder,
		PublicID:       safeFilename,
		UseFilename:    api.Bool(false),
		UniqueFilename: api.Bool(true),
		ResourceType:   "raw",
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

func UploadBytes(data []byte, filename string) (string, error) {
	if Cld == nil {
		return "", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	folder := os.Getenv("CLOUDINARY_UPLOAD_FOLDER")
	if folder == "" {
		folder = "sakti-apps"
	}

	safeFilename := sanitizeFilename(filename)

	reader := bytes.NewReader(data)

	resp, err := Cld.Upload.Upload(ctx, reader, uploader.UploadParams{
		Folder:         folder,
		PublicID:       safeFilename,
		UseFilename:    api.Bool(false),
		UniqueFilename: api.Bool(true),
		ResourceType:   "auto",
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

func DeleteFile(publicID string) error {
	if Cld == nil {
		return nil
	}

	if publicID == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := Cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}

func DeleteFileWithFolder(publicID string) error {
	if Cld == nil {
		return nil
	}

	if publicID == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := Cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
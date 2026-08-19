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

func sanitizeFilename(filename string) string {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.ReplaceAll(name, "(", "")
	name = strings.ReplaceAll(name, ")", "")
	name = strings.ToLower(name)
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
		return "", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	now := time.Now()
	tanggal := now.Format("02012006")
	nama := strings.ReplaceAll(strings.ToLower(karyawanNama), " ", "-")
	publicID := nama + "_" + tipe + "_" + tanggal

	folder := os.Getenv("CLOUDINARY_UPLOAD_FOLDER")
	if folder == "" {
		folder = "sakti-apps"
	}
	folder = folder + "/presensi"

	resp, err := Cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         folder,
		PublicID:       publicID,
		UseFilename:    api.Bool(false),
		UniqueFilename: api.Bool(true),
		ResourceType:   "image",
	})
	if err != nil {
		return "", err
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := Cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
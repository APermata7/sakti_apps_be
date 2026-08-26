package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func DownloadImage(url string) (string, error) {
	if url == "" {
		return "", nil
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("gagal membuat request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal download gambar: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gagal download gambar: status %d", resp.StatusCode)
	}

	ext := filepath.Ext(url)
	if ext == "" || ext == ".webp" {
		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "png") {
			ext = ".png"
		} else if strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg") {
			ext = ".jpg"
		} else if strings.Contains(contentType, "webp") {
			ext = ".webp"
		} else {
			ext = ".png"
		}
	}

	tmpFile, err := os.CreateTemp("", "image-*"+ext)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

func CleanupTempFiles(paths ...string) {
	for _, path := range paths {
		if path != "" {
			os.Remove(path)
		}
	}
}
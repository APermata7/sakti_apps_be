package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DownloadImage(url string) (string, error) {
	if url == "" {
		return "", nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("gagal download gambar: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gagal download gambar: status %d", resp.StatusCode)
	}

	ext := filepath.Ext(url)
	if ext == "" {
		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "png") {
			ext = ".png"
		} else if strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg") {
			ext = ".jpg"
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
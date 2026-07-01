package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type DayoffHoliday struct {
	Tanggal    string `json:"tanggal"`
	Keterangan string `json:"keterangan"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system env")
	}

	year := 2026
	if len(os.Args) > 1 {
		fmt.Sscanf(os.Args[1], "%d", &year)
	}

	url := fmt.Sprintf("https://raw.githubusercontent.com/gerinsp/dayoff-API/master/data/%d.json", year)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Gagal fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Data untuk tahun %d tidak ditemukan", year)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Gagal baca response: %v", err)
	}

	var holidays []DayoffHoliday
	if err := json.Unmarshal(body, &holidays); err != nil {
		log.Fatalf("Gagal parse JSON: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL tidak ditemukan di .env")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}
	defer pool.Close()

	inserted := 0
	skipped := 0

	for _, h := range holidays {
		jenis := "nasional"
		if strings.Contains(h.Keterangan, "Cuti Bersama") {
			jenis = "cuti_bersama"
		}

		query := `
			INSERT INTO libur (tanggal, nama, jenis, sumber, aktif)
			VALUES ($1, $2, $3, 'api', true)
			ON CONFLICT (tanggal) DO UPDATE SET 
				nama = EXCLUDED.nama,
				jenis = EXCLUDED.jenis,
				aktif = true,
				sumber = 'api',
				diperbarui_pada = NOW()
		`
		_, err := pool.Exec(context.Background(), query, h.Tanggal, h.Keterangan, jenis)
		if err != nil {
			skipped++
			continue
		}
		inserted++
	}

	fmt.Printf("Tahun: %d\n", year)
	fmt.Printf("Berhasil: %d\n", inserted)
	fmt.Printf("Gagal: %d\n", skipped)
	fmt.Printf("Total: %d\n", len(holidays))
}
<p align="center">
  <img src="https://drive.google.com/thumbnail?id=1EYcDhVXlj18zXUUOYEyinFvSgybJHQDO&sz=w2000" alt="Banner SAKTI Backend" width="100%">
</p>

<h1 align="center">
SAKTI Backend
</h1>

<p align="center">
Backend Service untuk Sistem Presensi dan Manajemen Kepegawaian KOPEGTEL Malang.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go" alt="Go">
  &nbsp;
  <img src="https://img.shields.io/badge/Fiber-v2-00ADD8?style=for-the-badge" alt="Fiber">
  &nbsp;
  <img src="https://img.shields.io/badge/Supabase-PostgreSQL-3ECF8E?style=for-the-badge&logo=supabase" alt="Supabase">
  &nbsp;
  <img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker" alt="Docker">
</p>

---

# 🎯 Tentang Proyek

**SAKTI (Sistem Absensi dan Kinerja Terintegrasi)** merupakan backend service yang dikembangkan untuk mendukung aplikasi presensi karyawan di **KOPEGTEL Malang**.

Backend ini dibangun menggunakan bahasa pemrograman **Go** dengan framework **Fiber** dan menggunakan **Supabase PostgreSQL** sebagai basis data utama.

Sistem dirancang untuk menggantikan sistem presensi sebelumnya sehingga proses pencatatan kehadiran, pengajuan cuti, persetujuan atasan, hingga pelaporan dapat dilakukan secara terintegrasi, cepat, dan mudah dipelihara.

---

# 🚀 Tujuan Pengembangan

Pengembangan backend ini bertujuan untuk:

- Meningkatkan efisiensi proses presensi karyawan.
- Mengintegrasikan seluruh proses administrasi kepegawaian dalam satu sistem.
- Mendukung validasi presensi menggunakan lokasi dan pengenalan wajah.
- Mempermudah proses persetujuan cuti secara bertingkat.
- Menyediakan laporan presensi yang akurat dan terdokumentasi.
- Mengurangi proses administrasi manual.

---

# ✨ Fitur Utama

## 🔐 Autentikasi

- Login menggunakan JWT
- Manajemen Session
- Role Based Access Control (RBAC)

Role pengguna:

- Administrator
- HRD
- Atasan / Manager
- Karyawan

---

## 👤 Manajemen Pengguna

- Data Karyawan
- Profil Pengguna
- Divisi
- Jabatan
- Hak Akses

---

## 📍 Presensi

Fitur presensi meliputi:

- Check In
- Check Out
- Validasi Geofencing
- Validasi Face Recognition
- Perhitungan keterlambatan
- Rekap jam kerja
- Status Kehadiran

Status presensi meliputi:

- Hadir
- Terlambat
- Izin
- Sakit
- Cuti
- Lembur

---

## 📋 Manajemen Pengajuan

Jenis pengajuan yang didukung:

### Cuti

- Mengurangi saldo cuti tahunan
- Maksimal pengajuan H-1
- Memerlukan persetujuan atasan dan HRD

### Dispensasi

- Tidak mengurangi saldo cuti
- Digunakan untuk keperluan tertentu
- Persetujuan langsung

### Cuti Darurat

- Dapat diajukan pada hari yang sama
- Digunakan untuk kondisi mendesak
- Memerlukan proses validasi HRD

---

## 📄 Generate Dokumen

Sistem dapat menghasilkan:

- Surat Cuti
- Surat Dispensasi
- Rekap Presensi
- Rekap Pengajuan

Dokumen dihasilkan dalam format PDF.

---

## 🔔 Notifikasi

Jenis notifikasi:

- In App Notification
- WhatsApp API

Notifikasi dikirim ketika:

- Pengajuan dibuat
- Pengajuan disetujui
- Pengajuan ditolak
- Finalisasi HRD

---

## 📊 Pelaporan

Laporan yang tersedia:

- Rekap Presensi
- Rekap Keterlambatan
- Rekap Cuti
- Rekap Lembur

Format ekspor:

- CSV
- PDF

---

## 📝 Audit Log

Seluruh aktivitas penting akan dicatat seperti:

- Login
- Perubahan data
- Approval
- Penghapusan data
- Generate laporan

---

# 🛠 Teknologi yang Digunakan

| Komponen | Teknologi |
|----------|------------|
| Bahasa Pemrograman | Go 1.25+ |
| Framework | Fiber |
| Database | PostgreSQL |
| Database Cloud | Supabase |
| Driver Database | pgx |
| Authentication | JWT |
| Configuration | Viper |
| Environment | Godotenv |
| Storage | Supabase Storage |
| Image Storage | Cloudinary |
| PDF Generator | gofpdf |
| Deployment | Docker |

---

# 🏛 Arsitektur Sistem

Backend menerapkan konsep **Clean Architecture** sehingga setiap lapisan memiliki tanggung jawab yang jelas.

```
Request

↓

Handler

↓

Usecase

↓

Repository

↓

Database

↓

Response
```

Keuntungan:

- Mudah dikembangkan
- Mudah diuji
- Mudah dipelihara
- Memisahkan Business Logic dari Database

---

# 📂 Struktur Folder

```text
sakti-backend
│
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── config/
│   ├── domain/
│   ├── handler/
│   ├── middleware/
│   ├── repository/
│   ├── usecase/
│   └── utils/
│
├── pkg/
│   ├── database/
│   └── logger/
│
├── migrations/
├── docs/
├── scripts/
├── .env.example
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

# ⚙ Instalasi

Clone repository

```bash
git clone https://github.com/APermata7/sakti-backend.git
```

Masuk ke direktori proyek

```bash
cd sakti-backend
```

Install dependency

```bash
go mod tidy
```

---

# 🔑 Konfigurasi Environment

Salin file contoh environment.

```bash
cp .env.example .env
```

---

# ▶ Menjalankan Aplikasi

Development

```bash
go run cmd/api/main.go
```

Aplikasi akan berjalan pada

```
http://localhost:8080
```

---

# 🐳 Menjalankan Menggunakan Docker

Build image

```bash
docker build -t sakti-backend .
```

Menjalankan container

```bash
docker run -p 8080:8080 sakti-backend
```

---

# 📡 Modul API

Backend menyediakan beberapa kelompok endpoint:

| Modul | Deskripsi |
|--------|-----------|
| Authentication | Login & Autentikasi |
| User | Data pengguna |
| Employee | Data karyawan |
| Attendance | Presensi |
| Leave | Pengajuan cuti |
| Approval | Persetujuan |
| Notification | Notifikasi |
| Dashboard | Ringkasan data |
| Report | Laporan |
---

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

## 🎯 Tentang Proyek

**SAKTI (Sistem Absensi dan Kinerja Terintegrasi)** merupakan backend service yang dikembangkan untuk mendukung aplikasi presensi karyawan di **KOPEGTEL Malang**.

Backend ini dibangun menggunakan bahasa pemrograman **Go** dengan framework **Fiber** dan menggunakan **Supabase PostgreSQL** sebagai basis data utama.

Sistem dirancang untuk menggantikan sistem presensi sebelumnya sehingga proses pencatatan kehadiran, pengajuan cuti, persetujuan atasan, hingga pelaporan dapat dilakukan secara terintegrasi, cepat, dan mudah dipelihara.

---

## 🚀 Tujuan Pengembangan

Pengembangan backend ini bertujuan untuk:

- Meningkatkan efisiensi proses presensi karyawan.
- Mengintegrasikan seluruh proses administrasi kepegawaian dalam satu sistem.
- Mendukung validasi presensi menggunakan lokasi (geofencing) dan pengenalan wajah (face recognition).
- Mempermudah proses persetujuan cuti secara bertingkat.
- Menyediakan laporan presensi yang akurat dan terdokumentasi.
- Mengurangi proses administrasi manual.

---

## ✨ Fitur Utama

### 🔐 Autentikasi & Otorisasi

- Login menggunakan JWT (JSON Web Token) dari Supabase Auth
- Manajemen Session
- Role Based Access Control (RBAC) dengan 4 role:

| Role | Hak Akses |
|------|-----------|
| **Admin** | Akses penuh ke semua fitur dan data |
| **Atasan** | Melihat dan menyetujui pengajuan bawahan |
| **HRD** | Mengelola data karyawan, cuti, dan finalisasi |
| **Karyawan** | Presensi, pengajuan cuti, dan data pribadi |

---

### 👤 Manajemen Pengguna

- CRUD Data Karyawan
- Profil Pengguna
- Divisi & Unit
- Level Jabatan (staff, officer, spv, ka_unit, manager, gm, hrd)
- Status Karyawan (aktif / nonaktif)
- Atasan Langsung

---

### 📍 Presensi

Fitur presensi meliputi:

- **Check In** dengan validasi:
  - Validasi Geofencing (radius kantor)
  - Validasi Face Recognition (FaceNet)
  - Deteksi keterlambatan berdasarkan jam kerja
- **Check Out** dengan validasi:
  - Validasi Geofencing
  - Perhitungan lembur otomatis
- Riwayat Presensi dengan filter
- Status Kehadiran hari ini

**Status Presensi:**
- Tepat Waktu
- Terlambat
- Izin
- Sakit
- Cuti
- Lembur

---

### 📋 Manajemen Pengajuan Cuti

Jenis pengajuan yang didukung:

| Jenis | Deskripsi | Proses |
|-------|-----------|--------|
| **Cuti** | Mengurangi saldo cuti tahunan | Atasan → HRD |
| **Dispensasi** | Tidak mengurangi saldo cuti | Atasan → HRD |
| **Cuti Darurat** | Dapat diajukan hari yang sama | HRD |

**Status Pengajuan:**
- Menunggu
- Disetujui (Atasan)
- Ditolak
- Dibatalkan
- Difinalisasi (HRD)

---

### 📄 Generate Dokumen

Sistem dapat menghasilkan dokumen PDF:

- **Surat Cuti** (2 atau 3 Tanda Tangan)
  - 3 TTD: Pemohon + Atasan + HRD
  - 2 TTD: Pemohon + HRD
- **Surat Dispensasi** (2 atau 3 Tanda Tangan)
- **Rekap Presensi**
- **Rekap Pengajuan**

---

### 🔔 Notifikasi

**Jenis Notifikasi:**
- In-App Notification
- WhatsApp API

**Trigger Notifikasi:**
- Pengajuan dibuat
- Pengajuan disetujui atasan
- Pengajuan ditolak
- Pengajuan difinalisasi HRD
- Reminder presensi

---

### 📊 Dashboard & Pelaporan

**Statistik Dashboard:**
- Total Karyawan
- Karyawan Aktif
- Presensi Hari Ini
- Presensi Terlambat
- Cuti Pending
- Total Cuti Tahun Berjalan

**Laporan Tersedia:**
- Rekap Presensi
- Rekap Keterlambatan
- Rekap Cuti
- Rekap Lembur

**Format Ekspor:**
- CSV
- PDF

---

### 📝 Audit Log

Seluruh aktivitas penting dicatat:

- Login / Logout
- Perubahan Data Karyawan
- Approval / Reject Pengajuan
- Finalisasi HRD
- Generate Laporan
- Hapus Data

---

## 🛠 Teknologi yang Digunakan

| Komponen | Teknologi |
|----------|------------|
| Bahasa Pemrograman | Go 1.25+ |
| Framework | Fiber v2 |
| Database | PostgreSQL |
| Database Cloud | Supabase |
| Driver Database | pgx/v5 |
| Authentication | Supabase Auth (JWT) |
| Environment | Godotenv |
| Storage | Supabase Storage |
| Image Storage | Cloudinary |
| PDF Generator | gofpdf/v2 |
| Face Recognition | FaceNet (TFLite/ONNX) |
| Notification | FCM + WhatsApp API |
| Deployment | Docker |

---

## 🏛 Arsitektur Sistem

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

## 📂 Struktur Folder

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

## ⚙ Instalasi

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

## 🔑 Konfigurasi Environment

Salin file contoh environment.

```bash
cp .env.example .env
```

---

## ▶ Menjalankan Aplikasi

Development

```bash
go run cmd/api/main.go
```

Aplikasi akan berjalan pada

```
http://localhost:8080
```

---

## 🐳 Menjalankan Menggunakan Docker

Build image

```bash
docker build -t sakti-backend .
```

Menjalankan container

```bash
docker run -p 8080:8080 sakti-backend
```

---

## 📡 Modul API

Backend menyediakan beberapa kelompok endpoint:

| Modul | Deskripsi | Endpoint Prefix |
|-------|-----------|-----------------|
| **Authentication** | Login, logout, profile, change password, forgot/reset password | `/api/auth/*` |
| **Attendance** | Check-in, check-out, today, history, update alasan terlambat | `/api/attendance/*` |
| **Leave** | Create, status, download surat, cancel, approve, reject, finalize | `/api/leave/*` |
| **Admin** | Dashboard, CRUD karyawan, CRUD libur, konfigurasi kerja | `/api/admin/*` |
| **Notification** | Get notifications, unread count, mark read, mark all read | `/api/notifikasi/*` |
| **Riwayat** | Riwayat aktivitas user | `/api/riwayat` |
| **Libur** | Cek hari libur (public) | `/api/libur/check` |
| **Upload** | Upload file, image, TTD | `/upload/*` |
---

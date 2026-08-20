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
  <img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go" alt="Go">
  &nbsp;
  <img src="https://img.shields.io/badge/Fiber-v2-00ADD8?style=for-the-badge" alt="Fiber">
  &nbsp;
  <img src="https://img.shields.io/badge/Supabase-PostgreSQL-3ECF8E?style=for-the-badge&logo=supabase" alt="Supabase">
  &nbsp;
  <img src="https://img.shields.io/badge/Cloudinary-Cloud-3448C5?style=for-the-badge&logo=cloudinary" alt="Cloudinary">
  &nbsp;
  <img src="https://img.shields.io/badge/Resend-Email-000000?style=for-the-badge&logo=resend" alt="Resend">
  &nbsp;
  <img src="https://img.shields.io/badge/FCM-Notifications-FFCA28?style=for-the-badge&logo=firebase" alt="FCM">
  &nbsp;
  <img src="https://img.shields.io/badge/Telegram-Bot-26A5E4?style=for-the-badge&logo=telegram" alt="Telegram Bot">
  &nbsp;
  <img src="https://img.shields.io/badge/cPanel-Deployment-FF6C2C?style=for-the-badge&logo=cpanel" alt="cPanel">
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
| **Atasan** | Melihat, menyetujui, dan menolak pengajuan bawahan |
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
- Tanda Tangan Digital (TTD)

---

### 📍 Presensi

Fitur presensi meliputi:

- **Check In** dengan validasi:
  - Validasi Geofencing (radius kantor)
  - Validasi Face Recognition (MobileFaceNet)
  - Deteksi keterlambatan berdasarkan jam kerja
  - Status lokasi (di_dalam_radius / di_luar_radius)
- **Check Out** dengan validasi:
  - Validasi Geofencing
  - Perhitungan lembur otomatis
  - Status lokasi keluar
- Riwayat Presensi dengan filter tanggal dan status
- Status Kehadiran hari ini

**Status Presensi:**
- Tepat Waktu
- Terlambat
- Belum Presensi

---

### 📋 Manajemen Pengajuan Cuti

Jenis pengajuan yang didukung:

| Jenis | Deskripsi | Proses |
|-------|-----------|--------|
| **Cuti** | Mengurangi saldo cuti tahunan | Atasan → HRD |
| **Dispensasi** | Tidak mengurangi saldo cuti, maksimal 2 hari | Langsung Final (Auto) |
| **Sakit** | Mengurangi saldo cuti tahunan | Atasan → HRD |

**Status Pengajuan:**
- Menunggu
- Disetujui (Atasan)
- Ditolak
- Dibatalkan
- Difinalisasi (HRD)

**Aturan Bisnis:**
- Cuti hanya bisa dibatalkan maksimal H-24 jam
- Kuota cuti = sisa cuti tahun ini + sisa cuti tahun lalu (berlaku sampai 31 Maret)
- Karyawan tanpa atasan langsung masuk ke HRD
- Karyawan dengan role atasan/manager langsung auto approve
- Karyawan dengan role HRD langsung auto final

---

### 📄 Generate Dokumen PDF

Sistem dapat menghasilkan dokumen PDF:

- **Surat Cuti** (2 atau 3 Tanda Tangan)
  - 3 TTD: Pemohon + Atasan + HRD (karyawan biasa)
  - 2 TTD: Pemohon + HRD (atasan/manager, tanpa atasan)
  - 2 TTD: Pemohon + Atasan (HRD)
- **Surat Dispensasi** (2 atau 3 Tanda Tangan)
- Format tanggal Indonesia (contoh: 24 Juli 2026)
- Logo Perusahaan dari konfigurasi

---

### 🔔 Notifikasi

**Jenis Notifikasi:**
- In-App Notification (FCM)
- Telegram Bot (untuk atasan)

**Trigger Notifikasi:**
- Pengajuan cuti baru → Atasan
- Pengajuan disetujui atasan → Karyawan
- Pengajuan ditolak → Karyawan
- Pengajuan difinalisasi HRD → Karyawan
- Pembatalan cuti → Atasan

---

### 📊 Dashboard & Pelaporan

**Dashboard Admin:**
- Total Karyawan
- Karyawan Aktif
- Total Terlambat (Bulan ini)
- Total Lembur (Bulan ini)
- Total Cuti Disetujui (Bulan ini)
- Grafik Karyawan per Departemen (Divisi-Unit)
- Grafik Presensi Masuk (Tepat Waktu, Terlambat, Belum Presensi)
- Grafik Presensi Keluar (Presensi Keluar, Presensi Lembur, Belum Presensi)
- Grafik Total Pengajuan Cuti (Disetujui, Ditolak, Menunggu, Dibatalkan)

**Laporan Admin:**
- Laporan Presensi (filter tanggal, status)
- Laporan Cuti (filter tanggal, status)
- Export ke CSV

---

### 📝 Riwayat & Log Aktivitas

Seluruh aktivitas penting dicatat:

- Login / Logout
- Check-in / Check-out
- Pengajuan Cuti (diajukan, disetujui, ditolak, dibatalkan, difinalisasi)
- Perubahan Password

---

## 🛠 Teknologi yang Digunakan

| Komponen | Teknologi |
|----------|------------|
| Bahasa Pemrograman | Go 1.21+ |
| Framework | Fiber v2 |
| Database | PostgreSQL |
| Database Cloud | Supabase |
| Driver Database | pgx/v5 |
| Authentication | Supabase Auth (JWT) |
| Email Service | Resend |
| Environment | Godotenv |
| Storage | Supabase Storage |
| Image Storage | Cloudinary |
| PDF Generator | gofpdf/v2 |
| Face Recognition | MobileFaceNet (InsightFace) |
| Notification | FCM + Telegram Bot |
| Deployment | cPanel |

---

## 🏛 Arsitektur Sistem

Backend menerapkan konsep **Clean Architecture** sehingga setiap lapisan memiliki tanggung jawab yang jelas.

```
Request

↓

Handler (HTTP Layer)

↓

Usecase (Business Logic)

↓

Repository (Data Access)

↓

Database (Supabase PostgreSQL)

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
│   ├── db/
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

## 📡 Modul API

Backend menyediakan beberapa kelompok endpoint:

| Modul | Deskripsi | Endpoint Prefix |
|-------|-----------|-----------------|
| **Authentication** | Login, logout, profile, change password, forgot/reset password | `/api/auth/*` |
| **Attendance** | Check-in, check-out, today, history, update alasan terlambat, work-config | `/api/attendance/*` |
| **Leave** | Create, status, balance, all, approval list, finalization list, download surat, cancel, approve, reject, finalize | `/api/leave/*` |
| **Admin** | Dashboard, CRUD karyawan, laporan presensi, laporan cuti, export CSV | `/api/admin/*` |
| **Config** | Konfigurasi kerja (jam kerja, radius, logo) | `/api/admin/konfigurasi/*` |
| **Holiday** | CRUD hari libur, cek hari libur (public) | `/api/admin/libur/*`, `/api/libur/*` |
| **Notification** | Get notifications, unread count, mark read, mark all read | `/api/notifikasi/*` |
| **TTD** | Upload, get, update, delete tanda tangan digital | `/api/ttd/*` |
| **Riwayat** | Riwayat aktivitas user | `/api/riwayat` |
| **Upload** | Upload file, image, TTD | `/upload/*` |

---

## 🌐 Deployment

### Production

Backend dideploy pada server cPanel dengan base URL:

```text
https://backendsakti.kopegtelmalang.co.id
```

### Build Application

Build binary untuk environment Linux:

```bash
GOOS=linux GOARCH=amd64 go build -o sakti-backend cmd/api/main.go
```

Setelah proses build:

1. Upload binary ke server.
2. Upload atau konfigurasi environment variables.
3. Pastikan permission executable telah sesuai.
4. Jalankan aplikasi melalui konfigurasi server atau process manager yang digunakan.

Contoh menjalankan binary:

```bash
chmod +x sakti-backend
./sakti-backend
```

---
CREATE TABLE IF NOT EXISTS public.pengajuan_cuti (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    karyawan_id UUID NOT NULL REFERENCES karyawan(id) ON DELETE CASCADE,
    sub_tipe VARCHAR(20) NOT NULL CHECK (sub_tipe IN ('izin', 'sakit', 'dispensasi')),
    tanggal_mulai DATE NOT NULL,
    tanggal_selesai DATE NOT NULL,
    total_hari INT NOT NULL,
    alasan TEXT,
    status VARCHAR(50) NOT NULL CHECK (status IN ('menunggu_atasan', 'menunggu_hrd', 'disetujui', 'ditolak', 'dibatalkan')),
    back_date BOOLEAN DEFAULT FALSE,
    mengurangi_cuti BOOLEAN DEFAULT TRUE,
    langsung_approve BOOLEAN DEFAULT FALSE,
    langsung_final BOOLEAN DEFAULT FALSE,
    judul_dokumen VARCHAR(100) DEFAULT 'PERMOHONAN/LAPORAN CUTI TAHUNAN',
    disetujui_oleh UUID REFERENCES karyawan(id) ON DELETE SET NULL,
    tanggal_disetujui TIMESTAMP,
    difinalisasi_oleh UUID REFERENCES karyawan(id) ON DELETE SET NULL,
    tanggal_difinalisasi TIMESTAMP,
    url_pdf VARCHAR(500),
    alasan_batal TEXT,
    tanggal_dibatalkan TIMESTAMP,
    alasan_ditolak TEXT,
    tanggal_ditolak TIMESTAMP,
    dibuat_pada TIMESTAMP DEFAULT NOW(),
    diperbarui_pada TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cuti_karyawan ON pengajuan_cuti(karyawan_id);
CREATE INDEX IF NOT EXISTS idx_cuti_status ON pengajuan_cuti(status);
CREATE INDEX IF NOT EXISTS idx_cuti_tanggal ON pengajuan_cuti(tanggal_mulai, tanggal_selesai);
CREATE INDEX IF NOT EXISTS idx_pengajuan_cuti_status_dibuat ON pengajuan_cuti (status, dibuat_pada DESC);
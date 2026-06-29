CREATE TABLE pengajuan_cuti (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    karyawan_id UUID NOT NULL REFERENCES karyawan(id) ON DELETE CASCADE,
    tipe_pengajuan VARCHAR(20) NOT NULL CHECK (tipe_pengajuan IN ('cuti', 'darurat')),
    sub_tipe VARCHAR(20) CHECK (sub_tipe IN ('izin', 'sakit', 'dispensasi', 'darurat')),
    tanggal_mulai DATE NOT NULL,
    tanggal_selesai DATE NOT NULL,
    total_hari INT NOT NULL,
    alasan TEXT,
    status VARCHAR(50) NOT NULL CHECK (status IN ('menunggu', 'disetujui', 'ditolak', 'dibatalkan')),
    back_date BOOLEAN DEFAULT FALSE,
    mengurangi_cuti BOOLEAN DEFAULT TRUE,
    langsung_approve BOOLEAN DEFAULT FALSE,
    judul_dokumen VARCHAR(100) DEFAULT 'PERMOHONAN/LAPORAN CUTI TAHUNAN',
    disetujui_oleh UUID REFERENCES karyawan(id) ON DELETE SET NULL,
    tanggal_disetujui TIMESTAMP,
    difinalisasi_oleh UUID REFERENCES karyawan(id) ON DELETE SET NULL,
    tanggal_difinalisasi TIMESTAMP,
    url_pdf VARCHAR(500),
    alasan_batal TEXT,
    tanggal_dibatalkan TIMESTAMP,
    dibuat_pada TIMESTAMP DEFAULT NOW(),
    diperbarui_pada TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_cuti_karyawan ON pengajuan_cuti(karyawan_id);
CREATE INDEX idx_cuti_status ON pengajuan_cuti(status);
CREATE INDEX idx_cuti_tanggal ON pengajuan_cuti(tanggal_mulai, tanggal_selesai);
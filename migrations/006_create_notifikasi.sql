CREATE TABLE notifikasi (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    karyawan_id UUID NOT NULL REFERENCES karyawan(id) ON DELETE CASCADE,
    jenis VARCHAR(50) NOT NULL CHECK (jenis IN ('pengajuan', 'persetujuan', 'penolakan', 'reminder')),
    channel VARCHAR(20) NOT NULL CHECK (channel IN ('inapp', 'whatsapp')),
    judul VARCHAR(255) NOT NULL,
    pesan TEXT NOT NULL,
    dibuat_pada TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_notifikasi_karyawan ON notifikasi(karyawan_id);
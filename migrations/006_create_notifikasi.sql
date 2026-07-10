CREATE TABLE notifikasi (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    karyawan_id UUID NOT NULL REFERENCES karyawan(id) ON DELETE CASCADE,
    jenis VARCHAR(50) NOT NULL CHECK (jenis IN ('pengajuan', 'persetujuan', 'penolakan', 'reminder')),
    channel VARCHAR(20) NOT NULL CHECK (channel IN ('inapp', 'telegram')),
    judul VARCHAR(255) NOT NULL,
    pesan TEXT NOT NULL,
    dibaca BOOLEAN NOT NULL DEFAULT FALSE,
    dibaca_pada TIMESTAMP,
    referensi_id UUID,
    referensi_tipe VARCHAR(50),
    dibuat_pada TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_notifikasi_karyawan ON notifikasi(karyawan_id);
CREATE INDEX idx_notifikasi_belum_dibaca ON notifikasi(karyawan_id, dibaca);
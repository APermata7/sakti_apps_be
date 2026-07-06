CREATE TABLE libur (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tanggal DATE NOT NULL UNIQUE,
    nama VARCHAR(255) NOT NULL,
    jenis VARCHAR(20) NOT NULL CHECK (jenis IN ('nasional', 'cuti_bersama')),
    aktif BOOLEAN NOT NULL DEFAULT TRUE,
    sumber VARCHAR(20) NOT NULL DEFAULT 'api' CHECK (sumber IN ('api', 'manual')),
    dibuat_pada TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    diperbarui_pada TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_libur_tanggal ON libur(tanggal);
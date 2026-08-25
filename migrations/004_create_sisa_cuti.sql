CREATE TABLE sisa_cuti (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    karyawan_id UUID NOT NULL REFERENCES karyawan(id) ON DELETE CASCADE,
    tahun INT NOT NULL,
    jumlah_cuti INT DEFAULT 12,
    telah_dilaksanakan INT DEFAULT 0,
    akan_dilaksanakan INT DEFAULT 0,
    sisa_cuti INT DEFAULT 12,
    dibuat_pada TIMESTAMPTZ DEFAULT NOW(),
    diperbarui_pada TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_sisa_cuti_karyawan_tahun ON sisa_cuti(karyawan_id, tahun);
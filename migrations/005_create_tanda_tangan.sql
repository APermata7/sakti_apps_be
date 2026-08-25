CREATE TABLE tanda_tangan (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    karyawan_id UUID NOT NULL REFERENCES karyawan(id) ON DELETE CASCADE,
    url_tanda_tangan VARCHAR(500) NOT NULL,
    hash_tanda_tangan VARCHAR(255),
    diunggah_pada TIMESTAMPTZ DEFAULT NOW(),
    diperbarui_pada TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_ttd_karyawan ON tanda_tangan(karyawan_id);
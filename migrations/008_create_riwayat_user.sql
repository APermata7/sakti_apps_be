CREATE TABLE riwayat_user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    karyawan_id UUID NOT NULL REFERENCES karyawan(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,
    detail TEXT,
    dibuat_pada TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_riwayat_karyawan ON riwayat_user(karyawan_id);
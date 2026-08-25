CREATE TABLE log_aktivitas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    karyawan_id UUID NOT NULL REFERENCES karyawan(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,
    detail TEXT,
    dibuat_pada TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_log_karyawan ON log_aktivitas(karyawan_id);
CREATE INDEX idx_log_action ON log_aktivitas(action);
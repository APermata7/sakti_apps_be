CREATE TABLE fcm_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    karyawan_id UUID NOT NULL REFERENCES karyawan(id) ON DELETE CASCADE,
    fcm_token VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    dibuat_pada TIMESTAMP DEFAULT NOW(),
    diperbarui_pada TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_fcm_token ON fcm_tokens(fcm_token);
CREATE INDEX idx_fcm_karyawan ON fcm_tokens(karyawan_id);
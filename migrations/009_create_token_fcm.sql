CREATE TABLE token_fcm (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    karyawan_id UUID NOT NULL REFERENCES karyawan(id) ON DELETE CASCADE,
    fcm_token VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    dibuat_pada TIMESTAMPTZ DEFAULT NOW(),
    diperbarui_pada TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_token_fcm ON token_fcm(fcm_token);
CREATE INDEX idx_token_fcm_karyawan ON token_fcm(karyawan_id);
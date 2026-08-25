CREATE TABLE presensi (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    karyawan_id UUID NOT NULL REFERENCES karyawan(id) ON DELETE CASCADE,
    kantor_id UUID REFERENCES kantor(id),
    tanggal DATE NOT NULL,
    jam_masuk TIME,
    jam_keluar TIME,
    status VARCHAR(50) NOT NULL CHECK (status IN ('tepat_waktu', 'terlambat')),
    lintang_masuk DECIMAL(10, 8),
    bujur_masuk DECIMAL(11, 8),
    lintang_keluar DECIMAL(10, 8),
    bujur_keluar DECIMAL(11, 8),
    validasi_wajah BOOLEAN DEFAULT FALSE,
    face_similarity DECIMAL(5, 4) DEFAULT 0,
    url_foto VARCHAR(500),
    alasan_terlambat TEXT,
    lembur BOOLEAN DEFAULT FALSE,
    jam_lembur DECIMAL(4, 2) DEFAULT 0,
    distance_meter DECIMAL(10, 2) DEFAULT 0,
    is_outside_radius BOOLEAN DEFAULT FALSE,
    location_status_masuk VARCHAR(20) DEFAULT NULL,
    location_status_keluar VARCHAR(20) DEFAULT NULL,
    dibuat_pada TIMESTAMPTZ DEFAULT NOW(),
    diperbarui_pada TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_presensi_karyawan ON presensi(karyawan_id);
CREATE INDEX idx_presensi_tanggal ON presensi(tanggal);
CREATE INDEX idx_presensi_status ON presensi(status);
CREATE INDEX idx_presensi_kantor ON presensi(kantor_id);
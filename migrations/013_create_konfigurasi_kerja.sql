CREATE TABLE konfigurasi_kerja (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama_kantor VARCHAR(255) NOT NULL DEFAULT 'KOPEGTEL MALANG',
    kantor_id UUID REFERENCES kantor(id),
    lat_kantor DECIMAL(10, 8) NOT NULL DEFAULT -7.942777,
    long_kantor DECIMAL(11, 8) NOT NULL DEFAULT 112.641110,
    logo_kantor VARCHAR(255),
    jam_masuk TIME NOT NULL DEFAULT '08:30',
    jam_minimal_masuk TIME NOT NULL DEFAULT '08:00',
    jam_pulang TIME NOT NULL DEFAULT '17:00',
    jam_minimal_pulang TIME NOT NULL DEFAULT '16:00',
    radius_kantor INTEGER NOT NULL DEFAULT 500,
    diperbarui_oleh UUID REFERENCES karyawan(id) ON DELETE SET NULL,
    diperbarui_pada TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_konfigurasi_diperbarui ON konfigurasi_kerja(diperbarui_pada);

INSERT INTO konfigurasi_kerja (jam_masuk, jam_minimal_masuk, jam_pulang, jam_minimal_pulang, radius_kantor)
VALUES ('08:30:00', '08:00:00', '17:00:00', '16:00:00', 500);
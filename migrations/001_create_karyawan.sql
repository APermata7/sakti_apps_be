CREATE TABLE karyawan (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama_lengkap VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    nomor_telepon VARCHAR(20),
    foto_url VARCHAR(500),
    peran VARCHAR(50) NOT NULL CHECK (peran IN ('admin', 'manager', 'hrd', 'karyawan', 'kepala_unit')),
    level_jabatan VARCHAR(50) NOT NULL CHECK (level_jabatan IN ('staff', 'spv', 'ka_unit', 'manager', 'gm', 'pengurus')),
    atasan_langsung_id UUID REFERENCES karyawan(id) ON DELETE SET NULL,
    divisi VARCHAR(100),
    unit VARCHAR(100),
    status_karyawan VARCHAR(20) DEFAULT 'aktif' CHECK (status_karyawan IN ('aktif', 'nonaktif')),
    kata_sandi_hash VARCHAR(255) NOT NULL,
    dibuat_pada TIMESTAMP DEFAULT NOW(),
    diperbarui_pada TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_karyawan_email ON karyawan(email);
CREATE INDEX idx_karyawan_peran ON karyawan(peran);
CREATE INDEX idx_karyawan_atasan ON karyawan(atasan_langsung_id);
CREATE TABLE IF NOT EXISTS public.karyawan (
    id UUID PRIMARY KEY,
    nama_lengkap VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    nomor_telepon VARCHAR(20),
    foto_url VARCHAR(500),
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'atasan', 'hrd', 'karyawan')),
    level_jabatan VARCHAR(50) NULL CHECK (level_jabatan IN ('staff', 'officer', 'spv', 'ka_unit', 'manager', 'gm', 'hrd')),
    atasan_langsung_id UUID REFERENCES karyawan(id) ON DELETE SET NULL,
    divisi VARCHAR(100),
    unit VARCHAR(100),
    status_karyawan VARCHAR(20) DEFAULT 'aktif' CHECK (status_karyawan IN ('aktif', 'nonaktif')),
    telegram_chat_id VARCHAR(100) DEFAULT NULL,
    dibuat_pada TIMESTAMP DEFAULT NOW(),
    diperbarui_pada TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_karyawan_email ON karyawan(email);
CREATE INDEX IF NOT EXISTS idx_karyawan_role ON karyawan(role);
CREATE INDEX IF NOT EXISTS idx_karyawan_atasan ON karyawan(atasan_langsung_id);
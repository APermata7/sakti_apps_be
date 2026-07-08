CREATE TABLE level_approval (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    level_jabatan VARCHAR(50) NOT NULL CHECK (level_jabatan IN ('staff', 'officer', 'spv', 'ka_unit', 'manager', 'gm', 'hrd')),
    jabatan_approver VARCHAR(50) NOT NULL CHECK (jabatan_approver IN ('ka_unit', 'manager', 'gm', 'ka_admin_support')),
    dibuat_pada TIMESTAMP DEFAULT NOW()
);

INSERT INTO level_approval (level_jabatan, jabatan_approver) VALUES
('staff', 'ka_unit'),
('officer', 'ka_unit'),
('spv', 'ka_unit'),
('ka_unit', 'manager'),
('manager', 'gm'),
('hrd', 'ka_admin_support');

CREATE INDEX idx_level_approval_jabatan ON level_approval(level_jabatan);
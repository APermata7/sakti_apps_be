CREATE TABLE level_approval (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    level_jabatan VARCHAR(50) NOT NULL CHECK (level_jabatan IN ('staff', 'spv', 'ka_unit', 'manager', 'gm', 'pengurus')),
    jabatan_approver VARCHAR(50) NOT NULL CHECK (jabatan_approver IN ('ka_unit', 'manager', 'gm', 'pengurus')),
    dibuat_pada TIMESTAMP DEFAULT NOW()
);

INSERT INTO level_approval (level_jabatan, jabatan_approver) VALUES
('staff', 'ka_unit'),
('spv', 'ka_unit'),
('ka_unit', 'manager'),
('manager', 'gm'),
('gm', 'pengurus');
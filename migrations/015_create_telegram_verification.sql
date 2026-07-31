CREATE TABLE telegram_verification (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL UNIQUE,
    chat_id VARCHAR(50) NOT NULL,
    username VARCHAR(100),
    karyawan_id UUID NULL REFERENCES karyawan(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    expired_at TIMESTAMP DEFAULT NOW() + INTERVAL '5 minutes',
    is_used BOOLEAN DEFAULT FALSE,
    used_at TIMESTAMP
);

CREATE INDEX idx_telegram_verification_code ON telegram_verification(code);
CREATE INDEX idx_telegram_verification_expired ON telegram_verification(expired_at);
CREATE INDEX idx_telegram_verification_used ON telegram_verification(is_used);
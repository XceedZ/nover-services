-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    -- Kolom identitas utama
    user_id BIGSERIAL PRIMARY KEY,
    user_code VARCHAR(255) UNIQUE NOT NULL, -- Dibuat lebih panjang untuk mengakomodasi google_id

    -- Kolom untuk login dan profil
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL DEFAULT '', 
    
    full_name VARCHAR(100) NOT NULL DEFAULT '',
    username VARCHAR(50) UNIQUE,
    avatar_url TEXT NOT NULL DEFAULT '',

    -- Metadata
    login_with VARCHAR(20) NOT NULL DEFAULT 'local',
    is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    create_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_datetime TIMESTAMPTZ
);

COMMENT ON TABLE users IS 'Menyimpan data pengguna, baik untuk login lokal maupun sosial media';
COMMENT ON COLUMN users.user_code IS 'Kode unik pengguna. Untuk login Google, diisi langsung dengan google_id. Untuk login lokal, diisi dengan MD5 dari kombinasi unik (misal: email+timestamp).';
COMMENT ON COLUMN users.password IS 'Hash dari password pengguna (misal: bcrypt). Dibiarkan string kosong jika login via sosial media.';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
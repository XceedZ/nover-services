-- Nama file yang disarankan: 00005_add_pen_name_to_users.sql
-- (Asumsi file migrasi sebelumnya adalah 00004)

-- +goose Up
-- +goose StatementBegin

-- Menambahkan kolom untuk nama pena ke tabel 'users'
ALTER TABLE users
    ADD COLUMN pen_name VARCHAR(100) UNIQUE DEFAULT '';

-- Menambahkan komentar untuk dokumentasi
COMMENT ON COLUMN users.pen_name IS 'Nama pena unik milik pengguna (pen name/alias).';

-- Menambahkan indeks untuk mempercepat pencarian berdasarkan nama pena
CREATE INDEX idx_users_pen_name ON users(pen_name);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Perintah ini akan membatalkan perubahan di atas
ALTER TABLE users
    DROP COLUMN IF EXISTS pen_name;

-- +goose StatementEnd

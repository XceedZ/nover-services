-- Nama file yang disarankan: 00004_add_profile_columns_to_users.sql
-- (Asumsi file migrasi 'banks' adalah 00003)

-- +goose Up
-- +goose StatementBegin

-- Menambahkan kolom baru ke tabel 'users' untuk data profil, bank, dan status penulis
ALTER TABLE users
    ADD COLUMN phone VARCHAR(20) UNIQUE DEFAULT '',
    ADD COLUMN instagram VARCHAR(100) DEFAULT '',
    ADD COLUMN bank_id BIGINT DEFAULT -99,
    ADD COLUMN account_number VARCHAR(50) DEFAULT '',
    ADD COLUMN flg_author CHAR(1) NOT NULL DEFAULT 'N' CHECK (flg_author IN ('Y', 'N'));

-- Menambahkan komentar untuk dokumentasi
COMMENT ON COLUMN users.phone IS 'Nomor telepon pengguna, harus unik.';
COMMENT ON COLUMN users.instagram IS 'Username akun Instagram pengguna.';
COMMENT ON COLUMN users.bank_id IS 'ID bank milik pengguna. Relasi ke tabel `banks` dijaga oleh aplikasi.';
COMMENT ON COLUMN users.account_number IS 'Nomor rekening bank milik pengguna.';
COMMENT ON COLUMN users.flg_author IS 'Flag untuk menandai apakah user adalah seorang penulis/author (Y/N).';

-- Menambahkan indeks untuk mempercepat pencarian
-- Indeks pada kolom yang sering dicari (WHERE) atau di-JOIN sangat disarankan.
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_bank_id ON users(bank_id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Perintah ini akan membatalkan semua perubahan di atas jika migrasi di-rollback
-- Menghapus kolom yang sudah ditambahkan pada blok 'Up'
ALTER TABLE users
    DROP COLUMN IF EXISTS phone,
    DROP COLUMN IF EXISTS instagram,
    DROP COLUMN IF EXISTS bank_id,
    DROP COLUMN IF EXISTS account_number,
    DROP COLUMN IF EXISTS flg_author;

-- +goose StatementEnd

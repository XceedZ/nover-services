-- +goose Up
-- +goose StatementBegin

-- Membuat tabel untuk mencatat genre yang disukai oleh pengguna
-- CATATAN: Foreign key tidak digunakan. Integritas data (memastikan user_id dan genre_id valid)
-- harus dikelola di level aplikasi.
CREATE TABLE user_genres (
    user_genre_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    genre_id BIGINT NOT NULL,
    create_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_datetime TIMESTAMPTZ,

    -- Constraint untuk mencegah data duplikat tetap dipertahankan
    CONSTRAINT user_loves_genre_once UNIQUE (user_id, genre_id)
);

-- Menambahkan indeks untuk mempercepat query berdasarkan user_id atau genre_id
CREATE INDEX idx_user_genres_user_id ON user_genres(user_id);
CREATE INDEX idx_user_genres_genre_id ON user_genres(genre_id);

-- Menambahkan komentar untuk dokumentasi database
COMMENT ON TABLE user_genres IS 'Tabel penghubung untuk mencatat genre yang disukai oleh setiap pengguna. Tidak menggunakan foreign key.';
COMMENT ON COLUMN user_genres.user_id IS 'ID pengguna yang menyukai genre ini. Relasi ke tabel `users` dijaga oleh aplikasi.';
COMMENT ON COLUMN user_genres.genre_id IS 'ID genre yang disukai. Relasi ke tabel `genres` dijaga oleh aplikasi.';
COMMENT ON COLUMN user_genres.create_datetime IS 'Waktu kapan pengguna pertama kali menyukai genre ini.';
COMMENT ON CONSTRAINT user_loves_genre_once ON user_genres IS 'Memastikan seorang pengguna hanya bisa menyukai satu genre spesifik sebanyak satu kali.';

-- Menerapkan trigger update_datetime otomatis ke tabel 'user_genres'
-- Asumsi fungsi trigger_set_timestamp() sudah ada dari migrasi sebelumnya
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON user_genres
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Perintah untuk membatalkan semua yang dilakukan di atas jika migrasi di-rollback
DROP TABLE IF EXISTS user_genres;

-- +goose StatementEnd

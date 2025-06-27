-- +goose Up
-- +goose StatementBegin

-- 1. Membuat struktur tabel utama untuk 'genres' dengan kolom baru
CREATE TABLE genres (
    genre_id BIGSERIAL PRIMARY KEY,
    genre_name VARCHAR(100) UNIQUE NOT NULL,
    genre_tl VARCHAR(100) UNIQUE NOT NULL, -- DITAMBAHKAN: Kolom untuk key terjemahan/kode
    remark TEXT,
    active_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    non_active_datetime TIMESTAMPTZ,
    create_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_datetime TIMESTAMPTZ
);

-- Menambahkan komentar untuk dokumentasi database
COMMENT ON TABLE genres IS 'Menyimpan daftar genre buku atau cerita.';
COMMENT ON COLUMN genres.genre_tl IS 'Key unik untuk keperluan kode dan terjemahan (i18n). Contoh: romance, fiksiSejarah.';
COMMENT ON COLUMN genres.active_datetime IS 'Waktu kapan genre ini mulai aktif dan bisa digunakan.';
COMMENT ON COLUMN genres.non_active_datetime IS 'Waktu kapan genre ini dinonaktifkan (untuk soft delete).';


-- 2. Menerapkan trigger update_datetime otomatis ke tabel 'genres'
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON genres
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


-- 3. Mengisi data awal (seeding) ke dalam tabel 'genres' termasuk kolom genre_tl
-- DIPERBARUI: Query INSERT sekarang menyertakan kolom 'genre_tl'
INSERT INTO genres (genre_name, genre_tl) VALUES
('Romantis', 'romance'),
('Fantasi', 'fantasy'),
('Fiksi Ilmiah', 'scienceFiction'),
('Misteri', 'mystery'),
('Horor', 'horror'),
('Thriller', 'thriller'),
('Aksi', 'action'),
('Petualangan', 'adventure'),
('Komedi', 'comedy'),
('Fiksi Sejarah', 'historicalFiction'),
('Fiksi Penggemar', 'fanfiction'),
('Remaja', 'youngAdult'),
('Dewasa', 'adult'),
('Spiritual', 'spiritual'),
('Slice of Life', 'sliceOfLife');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Perintah untuk membatalkan semua yang dilakukan di atas jika migrasi di-rollback
DROP TABLE IF EXISTS genres;

-- +goose StatementEnd

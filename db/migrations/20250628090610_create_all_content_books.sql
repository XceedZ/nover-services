-- +goose Up
-- +goose StatementBegin

-- Membuat tipe data ENUM untuk status buku dan chapter.
-- D=Draft, P=Published, C=Completed, H=On_Hold
CREATE TYPE book_status AS ENUM ('D', 'P', 'C', 'H');
CREATE TYPE chapter_status AS ENUM ('D', 'P');

-- 1. Tabel utama untuk BUKU
CREATE TABLE books (
    book_id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    cover_image_url TEXT,
    status book_status NOT NULL DEFAULT 'D',
    rating_average NUMERIC(3, 2) NOT NULL DEFAULT 0.00,
    total_views BIGINT NOT NULL DEFAULT 0,
    create_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_datetime TIMESTAMPTZ
);
COMMENT ON TABLE books IS 'Menyimpan informasi master untuk setiap buku/cerita.';

-- 2. Tabel untuk CHAPTER
CREATE TABLE chapters (
    chapter_id BIGSERIAL PRIMARY KEY,
    book_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    chapter_order INT NOT NULL,
    status chapter_status NOT NULL DEFAULT 'D',
    coin_cost INT NOT NULL DEFAULT 0, -- 0 berarti gratis
    total_views BIGINT NOT NULL DEFAULT 0,
    create_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_datetime TIMESTAMPTZ
);
COMMENT ON TABLE chapters IS 'Menyimpan konten untuk setiap chapter dari sebuah buku.';

-- 3. Tabel penghubung antara AUTHOR (User) dan BUKU
CREATE TABLE author_books (
    user_id BIGINT NOT NULL,
    book_id BIGINT NOT NULL
);
COMMENT ON TABLE author_books IS 'Tabel penghubung untuk hubungan many-to-many antara pengguna (penulis) dan buku.';

-- 4. Tabel penghubung antara BUKU dan GENRE
CREATE TABLE book_genres (
    book_id BIGINT NOT NULL,
    genre_id BIGINT NOT NULL
);
COMMENT ON TABLE book_genres IS 'Tabel penghubung untuk hubungan many-to-many antara buku dan genre.';

-- 5. Tabel untuk REVIEW dan RATING BUKU
CREATE TABLE reviews (
    review_id BIGSERIAL PRIMARY KEY,
    book_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    rating INT NOT NULL,
    review_text TEXT,
    create_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_datetime TIMESTAMPTZ,
    UNIQUE(book_id, user_id)
);
COMMENT ON TABLE reviews IS 'Menyimpan rating dan ulasan dari pengguna untuk sebuah buku.';

-- 6. Tabel untuk KOMENTAR pada REVIEW
CREATE TABLE review_comments (
    comment_id BIGSERIAL PRIMARY KEY,
    review_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    comment_text TEXT NOT NULL,
    parent_comment_id BIGINT,
    create_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_datetime TIMESTAMPTZ
);
COMMENT ON TABLE review_comments IS 'Menyimpan komentar atau balasan untuk sebuah review.';

-- 7. Tabel untuk KOMENTAR pada CHAPTER (BARU DITAMBAHKAN)
CREATE TABLE chapter_comments (
    comment_id BIGSERIAL PRIMARY KEY,
    chapter_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    comment_text TEXT NOT NULL,
    parent_comment_id BIGINT, -- Untuk sistem balasan komentar (nested)
    create_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_datetime TIMESTAMPTZ
);
COMMENT ON TABLE chapter_comments IS 'Menyimpan komentar atau balasan untuk sebuah chapter spesifik.';


-- 8. Menerapkan trigger update otomatis ke semua tabel yang relevan
CREATE TRIGGER set_timestamp BEFORE UPDATE ON books FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
CREATE TRIGGER set_timestamp BEFORE UPDATE ON chapters FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
CREATE TRIGGER set_timestamp BEFORE UPDATE ON reviews FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
CREATE TRIGGER set_timestamp BEFORE UPDATE ON review_comments FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
CREATE TRIGGER set_timestamp BEFORE UPDATE ON chapter_comments FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp(); -- Trigger untuk tabel baru

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Membatalkan semua perintah di atas dengan urutan terbalik
DROP TABLE IF EXISTS chapter_comments; -- Hapus tabel baru dulu
DROP TABLE IF EXISTS review_comments;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS book_genres;
DROP TABLE IF EXISTS author_books;
DROP TABLE IF EXISTS chapters;
DROP TABLE IF EXISTS books;
DROP TYPE IF EXISTS chapter_status;
DROP TYPE IF EXISTS book_status;

-- +goose StatementEnd

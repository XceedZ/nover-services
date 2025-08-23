-- +goose Up
-- +goose StatementBegin

-- Menambahkan nilai baru ke tipe ENUM yang sudah ada
-- IF NOT EXISTS mencegah error jika nilai sudah ditambahkan sebelumnya
ALTER TYPE related_entity ADD VALUE IF NOT EXISTS 'BOOK_COMMENT';
ALTER TYPE related_entity ADD VALUE IF NOT EXISTS 'CHAPTER_COMMENT';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Menghapus nilai dari ENUM tidak didukung secara langsung di banyak versi PostgreSQL
-- dan biasanya tidak diperlukan untuk proses development. Biarkan kosong.

-- +goose StatementEnd
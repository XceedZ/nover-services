-- +goose Up
-- +goose StatementBegin

CREATE TABLE book_comments (
    comment_id BIGSERIAL PRIMARY KEY,
    book_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    comment_text TEXT NOT NULL,
    parent_comment_id BIGINT, -- Untuk sistem balasan komentar (nested)
    create_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_datetime TIMESTAMPTZ,
    CONSTRAINT fk_book FOREIGN KEY(book_id) REFERENCES books(book_id) ON DELETE CASCADE
);

COMMENT ON TABLE book_comments IS 'Menyimpan komentar atau diskusi umum pada sebuah buku.';

CREATE TRIGGER set_timestamp 
BEFORE UPDATE ON book_comments 
FOR EACH ROW 
EXECUTE PROCEDURE trigger_set_timestamp();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS book_comments;

-- +goose StatementEnd
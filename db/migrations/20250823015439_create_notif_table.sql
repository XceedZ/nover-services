-- +goose Up
-- +goose StatementBegin

-- Membuat tipe ENUM untuk jenis notifikasi dan entitas terkait
CREATE TYPE notification_type AS ENUM (
    'NEW_RATING', 
    'NEW_COMMENT', 
    'COMMENT_REPLY', 
    'SYSTEM_ANNOUNCEMENT', 
    'NEW_CHAPTER', 
    'NEW_BOOK_BY_AUTHOR'
);

CREATE TYPE related_entity AS ENUM (
    'BOOK',
    'CHAPTER',
    'REVIEW',
    'COMMENT',
    'USER'
);

-- Membuat tabel utama untuk notifikasi sistem
CREATE TABLE system_notifications (
    notification_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL, -- Pengguna yang menerima notifikasi
    notification_type notification_type NOT NULL,
    content TEXT NOT NULL, -- Isi pesan notifikasi, e.g., "Seraphina membalas komentar Anda."
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Kolom ini membantu frontend untuk navigasi saat notifikasi diklik
    related_entity_type related_entity, 
    related_entity_id BIGINT,

    create_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_datetime TIMESTAMPTZ
);

COMMENT ON TABLE system_notifications IS 'Menyimpan semua notifikasi yang dikirim ke pengguna.';
COMMENT ON COLUMN system_notifications.user_id IS 'ID pengguna yang akan menerima notifikasi ini.';
COMMENT ON COLUMN system_notifications.content IS 'Teks yang akan ditampilkan dalam notifikasi.';
COMMENT ON COLUMN system_notifications.related_entity_type IS 'Jenis entitas yang terkait dengan notifikasi (misal: CHAPTER, BOOK).';
COMMENT ON COLUMN system_notifications.related_entity_id IS 'ID dari entitas terkait (misal: chapter_id, book_id).';

-- Menambahkan index untuk query yang efisien
CREATE INDEX idx_system_notifications_user_id_is_read ON system_notifications(user_id, is_read);

-- Menerapkan trigger update otomatis
CREATE TRIGGER set_timestamp 
BEFORE UPDATE ON system_notifications 
FOR EACH ROW 
EXECUTE PROCEDURE trigger_set_timestamp();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS system_notifications;
DROP TYPE IF EXISTS related_entity;
DROP TYPE IF EXISTS notification_type;

-- +goose StatementEnd
-- +goose Up
-- +goose StatementBegin

-- Menambahkan kolom actor_id yang bisa bernilai NULL
-- (NULL digunakan untuk notifikasi sistem yang tidak memiliki aktor spesifik)
ALTER TABLE system_notifications
ADD COLUMN actor_id BIGINT;

COMMENT ON COLUMN system_notifications.actor_id IS 'ID pengguna yang memicu notifikasi (misalnya, yang berkomentar).';

-- Menambahkan index pada actor_id untuk performa query
CREATE INDEX idx_system_notifications_actor_id ON system_notifications(actor_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE system_notifications
DROP COLUMN IF EXISTS actor_id;

-- +goose StatementEnd
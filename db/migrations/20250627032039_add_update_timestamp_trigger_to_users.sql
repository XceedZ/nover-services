-- +goose Up
-- +goose StatementBegin

-- Membuat sebuah fungsi yang akan mengembalikan waktu saat ini
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.update_datetime = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Membuat sebuah trigger yang akan memanggil fungsi di atas
-- SETIAP KALI ada operasi UPDATE pada tabel 'users'
CREATE OR REPLACE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Menghapus trigger dari tabel users secara spesifik
DROP TRIGGER IF EXISTS set_timestamp ON users;

-- Menghapus fungsi DAN semua objek yang bergantung padanya (seperti trigger lain)
DROP FUNCTION IF EXISTS trigger_set_timestamp() CASCADE;

-- +goose StatementEnd

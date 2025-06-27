-- Nama file yang disarankan: 00003_create_banks_table.sql
-- (Asumsi file migrasi sebelumnya adalah 00002)

-- +goose Up
-- +goose StatementBegin

-- 1. Membuat struktur tabel untuk 'banks'
CREATE TABLE banks (
    bank_id BIGSERIAL PRIMARY KEY,
    bank_name VARCHAR(100) UNIQUE NOT NULL,
    bank_code VARCHAR(10) UNIQUE,
    remark TEXT,
    active_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    non_active_datetime TIMESTAMPTZ,
    create_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_datetime TIMESTAMPTZ
);

-- Menambahkan komentar untuk dokumentasi
COMMENT ON TABLE banks IS 'Menyimpan daftar bank yang terdaftar di Indonesia.';
COMMENT ON COLUMN banks.bank_name IS 'Nama resmi bank.';
COMMENT ON COLUMN banks.bank_code IS 'Kode unik bank yang digunakan untuk transfer antarbank.';

-- 2. Menerapkan trigger update_datetime otomatis
-- Asumsi fungsi trigger_set_timestamp() sudah ada dari migrasi sebelumnya
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON banks
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- 3. Mengisi data awal (seeding) bank-bank umum di Indonesia
-- Data diambil dari sumber terpercaya OJK dan Wikipedia.
INSERT INTO banks (bank_name, bank_code) VALUES
('BANK BCA', '014'),
('BANK MANDIRI', '008'),
('BANK BNI', '009'),
('BANK BRI', '002'),
('BANK BTN', '200'),
('BANK CIMB NIAGA', '022'),
('BANK DANAMON', '011'),
('BANK PERMATA', '013'),
('BANK PANIN', '019'),
('BANK OCBC NISP', '028'),
('BANK MEGA', '426'),
('BANK SINARMAS', '153'),
('BANK UOB INDONESIA', '023'),
('BANK BTPN', '213'),
('BANK JAGO', '542'),
('ALLO BANK', '567'),
('BANK SYARIAH INDONESIA (BSI)', '451');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Perintah untuk membatalkan semua yang dilakukan di atas jika migrasi di-rollback
DROP TABLE IF EXISTS banks;

-- +goose StatementEnd

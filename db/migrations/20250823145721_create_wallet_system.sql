-- +goose Up
-- +goose StatementBegin

-- Tipe ENUM untuk membedakan Koin Biasa dan Koin Bonus
CREATE TYPE coin_type AS ENUM ('PAID', 'BONUS');

-- Tipe ENUM untuk berbagai jenis transaksi
CREATE TYPE transaction_type AS ENUM (
    'REGISTRATION', 
    'CHECK_IN', 
    'MISSION_REWARD', 
    'PURCHASE', 
    'UNLOCK_CHAPTER', 
    'GIFT', 
    'EXPIRATION_ADJUSTMENT'
);

-- Tabel untuk menyimpan saldo koin setiap pengguna
CREATE TABLE wallets (
    wallet_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    paid_coins BIGINT NOT NULL DEFAULT 0,
    bonus_coins BIGINT NOT NULL DEFAULT 0,
    update_datetime TIMESTAMPTZ,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
COMMENT ON TABLE wallets IS 'Menyimpan saldo Koin Biasa dan Koin Bonus untuk setiap pengguna.';

-- Tabel untuk mencatat semua riwayat transaksi koin
CREATE TABLE coin_transactions (
    transaction_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    transaction_type transaction_type NOT NULL,
    coin_type coin_type NOT NULL,
    amount INT NOT NULL, -- Positif untuk pendapatan, negatif untuk pengeluaran
    description TEXT, -- Contoh: "Bonus Check-in hari ke-23"
    related_entity_id BIGINT, -- Bisa diisi chapter_id, mission_id, dll.
    expiry_date DATE, -- Tanggal kedaluwarsa khusus untuk Koin Bonus
    create_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
COMMENT ON TABLE coin_transactions IS 'Mencatat semua histori pendapatan dan pengeluaran koin.';
COMMENT ON COLUMN coin_transactions.amount IS 'Jumlah koin. (+) untuk pendapatan, (-) untuk pengeluaran.';

-- Index untuk mempercepat query histori
CREATE INDEX idx_coin_transactions_user_id ON coin_transactions(user_id);

-- Trigger
CREATE TRIGGER set_timestamp BEFORE UPDATE ON wallets FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS coin_transactions;
DROP TABLE IF EXISTS wallets;
DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS coin_type;

-- +goose StatementEnd
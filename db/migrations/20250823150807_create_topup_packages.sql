-- +goose Up
-- +goose StatementBegin

CREATE TABLE topup_packages (
    package_id SERIAL PRIMARY KEY,
    package_name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(12, 2) NOT NULL,
    base_paid_coins INT NOT NULL,
    bonus_coins INT NOT NULL DEFAULT 0,
    sku_google_play VARCHAR(255) UNIQUE, -- ID Produk unik dari Google Play Console
    sku_app_store VARCHAR(255) UNIQUE, -- ID Produk unik dari App Store Connect
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    display_order INT -- Untuk mengurutkan tampilan paket di aplikasi
);
COMMENT ON TABLE topup_packages IS 'Menyimpan konfigurasi paket-paket top-up koin.';
COMMENT ON COLUMN topup_packages.sku_google_play IS 'Product ID dari Google Play Console.';
COMMENT ON COLUMN topup_packages.sku_app_store IS 'Product ID dari App Store Connect.';

-- Mengisi tabel dengan data paket yang sudah kita diskusikan
INSERT INTO topup_packages 
    (package_name, description, price, base_paid_coins, bonus_coins, sku_google_play, sku_app_store, display_order) 
VALUES
    ('Paket Hemat', '100 Koin', 10000.00, 100, 0, 'nover.coins.100', 'nover.coins.100', 10),
    ('Paket Standar', 'Paling Populer!', 25000.00, 250, 25, 'nover.coins.275', 'nover.coins.275', 20),
    ('Paket Super', 'Bonus Lebih Besar', 50000.00, 500, 75, 'nover.coins.575', 'nover.coins.575', 30),
    ('Paket Sultan', 'Untung Maksimal!', 100000.00, 1000, 200, 'nover.coins.1200', 'nover.coins.1200', 40);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS topup_packages;

-- +goose StatementEnd
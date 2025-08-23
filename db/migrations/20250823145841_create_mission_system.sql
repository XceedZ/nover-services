-- +goose Up
-- +goose StatementBegin

-- Tipe ENUM untuk berbagai jenis misi (REVISI dengan tambahan tipe baru)
CREATE TYPE mission_type AS ENUM (
    'READING_DURATION',
    'COMMENT',
    'GIVE_RATING'
);

-- Tabel master untuk semua misi yang tersedia
CREATE TABLE missions (
    mission_id BIGSERIAL PRIMARY KEY,
    mission_type mission_type NOT NULL,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);
COMMENT ON TABLE missions IS 'Menyimpan daftar semua misi yang ada di aplikasi.';

-- Contoh pengisian data untuk misi
INSERT INTO missions (mission_type, title, description) VALUES
('READING_DURATION', 'Membaca Karya', '+12 Koin Bonus | 0/40 menit'),
('COMMENT', 'Tinggalkan Komentar', 'Berikan komentar di buku apa saja untuk mendapatkan hadiah.'),
('GIVE_RATING', 'Beri Rating Buku', 'Beri rating pada sebuah buku setelah membacanya.');

-- Tabel untuk tingkatan/tier hadiah dalam sebuah misi
CREATE TABLE mission_tiers (
    tier_id BIGSERIAL PRIMARY KEY,
    mission_id BIGINT NOT NULL,
    threshold INT NOT NULL, -- Target yang harus dicapai (misal: 20 menit, 1 komentar)
    reward_amount INT NOT NULL,
    tier_order INT NOT NULL, -- Urutan tier (1, 2, 3, ...)
    CONSTRAINT fk_mission FOREIGN KEY(mission_id) REFERENCES missions(mission_id) ON DELETE CASCADE
);
COMMENT ON TABLE mission_tiers IS 'Menentukan tingkatan hadiah untuk setiap misi.';

-- Contoh pengisian data tier untuk semua misi
-- Tier Misi Membaca (ID 1)
INSERT INTO mission_tiers (mission_id, threshold, reward_amount, tier_order) VALUES
(1, 20, 5, 1),
(1, 40, 7, 2);

-- Tier Misi Komentar (ID 2)
INSERT INTO mission_tiers (mission_id, threshold, reward_amount, tier_order) VALUES
(2, 1, 5, 1); -- Cukup 1 komentar untuk dapat 5 koin

-- Tier Misi Rating (ID 3)
INSERT INTO mission_tiers (mission_id, threshold, reward_amount, tier_order) VALUES
(3, 1, 5, 1); -- Cukup 1 rating untuk dapat 5 koin

-- Tabel untuk melacak progres misi harian setiap pengguna
CREATE TABLE user_mission_progress (
    progress_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    mission_id BIGINT NOT NULL,
    progress_date DATE NOT NULL,
    current_value INT NOT NULL DEFAULT 0, -- Nilai progres saat ini (misal: total menit membaca)
    last_claimed_tier_id BIGINT,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_mission FOREIGN KEY(mission_id) REFERENCES missions(mission_id) ON DELETE CASCADE,
    CONSTRAINT fk_claimed_tier FOREIGN KEY(last_claimed_tier_id) REFERENCES mission_tiers(tier_id),
    UNIQUE(user_id, mission_id, progress_date)
);
COMMENT ON TABLE user_mission_progress IS 'Mencatat progres misi harian setiap pengguna.';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS user_mission_progress;
DROP TABLE IF EXISTS mission_tiers;
DROP TABLE IF EXISTS missions;
DROP TYPE IF EXISTS mission_type;

-- +goose StatementEnd
-- +goose Up
-- +goose StatementBegin

-- Tabel konfigurasi untuk hadiah check-in harian
CREATE TABLE daily_checkin_rewards (
    day_number INT PRIMARY KEY,
    reward_amount INT NOT NULL
);
COMMENT ON TABLE daily_checkin_rewards IS 'Konfigurasi hadiah untuk setiap hari check-in.';

-- Pengisian data reward yang lebih seimbang (Final)
INSERT INTO daily_checkin_rewards (day_number, reward_amount) VALUES
(1, 4), (2, 4), (3, 5), (4, 4), (5, 5), (6, 4), (7, 10),  -- Minggu 1
(8, 4), (9, 5), (10, 4), (11, 5), (12, 6), (13, 5), (14, 12), -- Minggu 2
(15, 5), (16, 6), (17, 5), (18, 6), (19, 7), (20, 6), (21, 15), -- Minggu 3
(22, 6), (23, 7), (24, 6), (25, 10), (26, 7), (27, 6), (28, 20), -- Minggu 4
(29, 8), (30, 10), (31, 40); -- Akhir Bulan

-- Tabel untuk mencatat riwayat check-in setiap pengguna
CREATE TABLE user_daily_checkins (
    checkin_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    checkin_date DATE NOT NULL,
    consecutive_streak INT NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    UNIQUE(user_id, checkin_date)
);
COMMENT ON TABLE user_daily_checkins IS 'Mencatat riwayat check-in harian setiap pengguna.';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS user_daily_checkins;
DROP TABLE IF EXISTS daily_checkin_rewards;

-- +goose StatementEnd
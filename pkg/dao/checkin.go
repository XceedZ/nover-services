package dao

import (
	"context"
	"fmt"
	"noversystem/pkg/tables"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CheckinDao struct {
	DB *pgxpool.Pool
}

func NewCheckinDao(db *pgxpool.Pool) *CheckinDao {
	return &CheckinDao{DB: db}
}

// GetAllCheckinRewards mengambil semua konfigurasi hadiah check-in.
func (d *CheckinDao) GetAllCheckinRewards(ctx context.Context) ([]tables.DailyCheckinReward, error) {
	var rewards []tables.DailyCheckinReward
	query := "SELECT day_number, reward_amount FROM daily_checkin_rewards ORDER BY day_number ASC"
	err := pgxscan.Select(ctx, d.DB, &rewards, query)
	return rewards, err
}

// GetUserCheckinsForMonth mengambil tanggal-tanggal user sudah check-in di bulan ini.
func (d *CheckinDao) GetUserCheckinsForMonth(ctx context.Context, userID int64, year int, month time.Month) ([]string, error) {
	var dates []string
	query := `
		SELECT TO_CHAR(checkin_date, 'YYYY-MM-DD') 
		FROM user_daily_checkins 
		WHERE user_id = $1 AND EXTRACT(YEAR FROM checkin_date) = $2 AND EXTRACT(MONTH FROM checkin_date) = $3
	`
	// ✨ PERBAIKAN: Kirim 'month' sebagai integer, bukan string
	err := pgxscan.Select(ctx, d.DB, &dates, query, userID, year, int(month))
	return dates, err
}

// PerformCheckin melakukan aksi check-in untuk user.
func (d *CheckinDao) PerformCheckin(ctx context.Context, userID int64) (int, error) {
	tx, err := d.DB.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	// 1. Cek apakah hari ini sudah check-in
	var exists bool
	today := time.Now()
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM user_daily_checkins WHERE user_id = $1 AND checkin_date = $2)", userID, today.Format("2006-01-02")).Scan(&exists)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, fmt.Errorf("sudah check-in hari ini")
	}

	// 2. Hitung total check-in di bulan ini untuk menentukan hadiah hari ke berapa
	var totalCheckinsSoFar int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*) FROM user_daily_checkins 
		WHERE user_id = $1 AND EXTRACT(YEAR FROM checkin_date) = $2 AND EXTRACT(MONTH FROM checkin_date) = $3
	`, userID, today.Year(), int(today.Month())).Scan(&totalCheckinsSoFar)
	if err != nil {
		return 0, err
	}

	rewardDay := totalCheckinsSoFar + 1

	// 3. Ambil hadiah untuk hari tersebut
	var rewardAmount int
	err = tx.QueryRow(ctx, "SELECT reward_amount FROM daily_checkin_rewards WHERE day_number = $1", rewardDay).Scan(&rewardAmount)
	if err != nil {
		return 0, fmt.Errorf("konfigurasi hadiah untuk hari ke-%d tidak ditemukan", rewardDay)
	}

	// 4. Masukkan record check-in baru
	// ✨ PERBAIKAN: Tambahkan kolom consecutive_streak dengan nilai default 0
	_, err = tx.Exec(ctx, "INSERT INTO user_daily_checkins (user_id, checkin_date, consecutive_streak) VALUES ($1, $2, $3)", userID, today.Format("2006-01-02"), 0)
	if err != nil {
		return 0, err
	}

	// 5. Update wallet dan catat transaksi (tidak ada perubahan)
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sql, args, _ := psql.Update("wallets").Set("bonus_coins", squirrel.Expr("bonus_coins + ?", rewardAmount)).Where(squirrel.Eq{"user_id": userID}).ToSql()
	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}

	desc := fmt.Sprintf("Bonus Check-in hari ke-%d", rewardDay)
	expiry := today.AddDate(0, 0, 7)
	sql, args, _ = psql.Insert("coin_transactions").
		Columns("user_id", "transaction_type", "coin_type", "amount", "description", "expiry_date").
		Values(userID, "CHECK_IN", "BONUS", rewardAmount, desc, expiry.Format("2006-01-02")).ToSql()
	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}

	return rewardAmount, tx.Commit(ctx)
}
